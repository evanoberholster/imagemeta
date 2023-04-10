package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/pkg/errors"
)

func (r *Reader) ReadMetadata() (err error) {
	b, err := r.readBox()
	if err != nil {
		buf, err := r.br.Peek(128)
		fmt.Println(buf, err, len(buf))
		fmt.Println(string(buf))
		return errors.Wrapf(err, "ReadMetadata")
	}
	switch b.boxType {
	case typeMdat:
		err = r.readMdat(&b)
	case typeMeta:
		err = r.readMeta(&b)
		b.close()
	case typeMoov:
		err = r.readMoovBox(&b)
		b.close()
	case typeUUID:
		err = r.readUUIDBox(&b)
	default:
		if logLevelInfo() {
			logInfo().Object("box", b).Send()
		}
	}
	if err != nil && logLevelError() {
		logError().Object("box", b).Err(err).Send()
	}
	return err
}

func (r *Reader) readMdat(b *box) (err error) {
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	if r.heic.exif.ol.offset == 0 {
		return b.close()
	}
	inner, err := r.newExifBox(b)
	if err != nil {
		if logLevelError() {
			logError().Object("box", inner).Err(err).Send()
		}
		return
	}
	header, err := readExifHeader(&inner, ifds.IFD0, imagetype.ImageHEIF)
	if err != nil {
		panic(err)
	}

	if r.ExifReader != nil {
		if err = r.ExifReader(&inner, header); err != nil {
			if logLevelError() {
				logError().Object("box", inner).Err(err).Send()
			}
		}
	}

	if logLevelInfo() {
		logInfo().Object("box", inner).Int("remain", inner.remain).Send()
	}

	return b.close()
}

func (r *Reader) newExifBox(b *box) (inner box, err error) {
	if _, err = b.Discard(int(r.heic.exif.ol.offset) - b.offset - 16); err != nil {
		return
	}
	buf, err := b.Peek(16)
	if err != nil {
		return
	}
	var size int
	for i := 0; i < len(buf); i += 4 {
		if string(buf[i+4:i+4+4]) == "Exif" {
			size = int(bmffEndian.Uint32(buf[i:i+4])) + i
			break
		}
	}

	inner = box{
		reader:  b.reader,
		outer:   b,
		boxType: typeExif,
		offset:  int(b.size) - b.remain + b.offset,
		size:    int64(r.heic.exif.ol.length),
		remain:  int(r.heic.exif.ol.length),
	}

	_, err = inner.Discard(size + 4)
	return inner, err
}

func readExifHeader(b *box, firstIfd ifds.IfdType, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		err = errors.WithMessage(err, "readExifHeader")
		return
	}
	endian := utils.BinaryOrder(buf[:4])
	header = meta.NewExifHeader(endian, endian.Uint32(buf[4:8]), 0, uint32(b.remain), imagetype.ImageCR3)
	header.FirstIfd = firstIfd
	if logLevelInfo() {
		logInfo().Object("box", b).Object("header", header).Send()
	}
	_, err = b.Discard(8)
	return header, err
}

func (r *Reader) readMeta(b *box) (err error) {
	if !b.isType(typeMeta) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	if err = b.readFlags(); err != nil {
		return err
	}
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeUUID:
			err = r.readUUIDBox(&inner)
		case typeHdlr:
			_, err = readHdlr(&inner)
		case typePitm:
			r.heic.pitm, err = readPitm(&inner)
		case typeIinf:
			err = r.readIinf(&inner)
		case typeIref:
			err = readIref(&inner)
		case typeIprp:
			err = readIprp(&inner)
		case typeIdat:
			r.heic.idat, err = readIdat(&inner)
		case typeIloc:
			err = r.readIloc(&inner)
		default:
			if logLevelInfo() {
				logInfo().Object("box", inner).Send()
			}
		}
		if err != nil && logLevelError() {
			logError().Object("box", inner).Err(err).Send()
		}

		if err = inner.close(); err != nil {
			logError().Object("box", inner).Err(err).Send()
			break
		}
	}
	return b.close()
}

// ReadMOOV reads an 'moov' box from a BMFF file.
func (r *Reader) readMoovBox(b *box) (err error) {
	if !b.isType(typeMoov) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeUUID:
			err = r.readUUIDBox(&inner)
		case typeTrak:
			_, err = readCrxTrakBox(&inner)
		//case typeMvhd:
		default:
			if logLevelInfo() {
				logInfo().Object("box", inner).Send()
			}
		}
		if err != nil && logLevelError() {
			logError().Object("box", inner).Err(err).Send()
		}
		if err = inner.close(); err != nil {
			logError().Object("box", inner).Err(err).Send()
			break
		}
	}
	return b.close()
}
