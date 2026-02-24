package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func (r *Reader) ReadMetadata() (err error) {
	for {
		b, readErr := r.readBox()
		if readErr != nil {
			return fmt.Errorf("ReadMetadata: %w", readErr)
		}
		switch b.boxType {
		case typeMdat:
			err = r.readMdat(&b)
		case typeExif:
			err = r.readExif(&b)
		case typeMeta:
			err = r.readMeta(&b)
			b.close()
		case typeMoov:
			err = r.readMoovBox(&b)
			b.close()
		case typeUUID:
			err = r.readUUIDBox(&b)
		case typeJXL, typeJumb, typeJxlc, typeJxll, typeJxlp:
			// JPEG XL container boxes can appear before metadata boxes.
			// Skip and continue scanning.
			err = b.close()
			if err == nil {
				continue
			}
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
	imageType := r.metadataImageType()
	header, err := readExifHeader(&inner, ifds.IFD0, imageType)
	if err != nil {
		return fmt.Errorf("readMdat: %w", err)
	}
	if r.exifReader != nil {
		if err = r.exifReader(&inner, header); err != nil {

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

func (r *Reader) readExif(b *box) (err error) {
	if !b.isType(typeExif) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}

	if err = seekExifTIFFHeader(b); err != nil {
		return fmt.Errorf("readExif: %w", err)
	}

	header, err := readExifHeader(b, ifds.IFD0, r.metadataImageType())
	if err != nil {
		return fmt.Errorf("readExif: %w", err)
	}

	if r.exifReader != nil {
		if err = r.exifReader(b, header); err != nil {
			if logLevelError() {
				logError().Object("box", b).Err(err).Send()
			}
		}
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
		err = fmt.Errorf("readExifHeader: %w", err)
		return
	}
	endian := utils.BinaryOrder(buf[:4])
	header = meta.NewExifHeader(endian, endian.Uint32(buf[4:8]), 0, uint32(b.remain), it)
	header.FirstIfd = firstIfd
	if logLevelInfo() {
		logInfo().Object("box", b).Object("header", header).Send()
	}
	_, err = b.Discard(8)
	return header, err
}

func seekExifTIFFHeader(b *box) error {
	if b.remain < 8 {
		return ErrBufLength
	}

	buf, err := b.Peek(8)
	if err != nil {
		return err
	}

	// Some Exif payloads include the APP1 prefix before TIFF data.
	if hasExifAPP1Prefix(buf) {
		if _, err = b.Discard(6); err != nil {
			return err
		}
		if b.remain < 8 {
			return ErrBufLength
		}
		buf, err = b.Peek(8)
		if err != nil {
			return err
		}
	}

	if hasTIFFHeader(buf[:4]) {
		return nil
	}

	// HEIF/JXL style Exif payloads can start with a 4-byte TIFF header offset.
	offset := int(bmffEndian.Uint32(buf[:4]))
	if offset < 0 || offset > b.remain-8 {
		return nil
	}
	_, err = b.Discard(4 + offset)
	return err
}

func hasExifAPP1Prefix(buf []byte) bool {
	return len(buf) >= 6 &&
		buf[0] == 'E' &&
		buf[1] == 'x' &&
		buf[2] == 'i' &&
		buf[3] == 'f' &&
		buf[4] == 0x00 &&
		buf[5] == 0x00
}

func hasTIFFHeader(buf []byte) bool {
	return len(buf) >= 4 &&
		((buf[0] == 'I' && buf[1] == 'I' && buf[2] == 0x2A && buf[3] == 0x00) ||
			(buf[0] == 'M' && buf[1] == 'M' && buf[2] == 0x00 && buf[3] == 0x2A))
}

func (r *Reader) metadataImageType() imagetype.ImageType {
	switch r.ftyp.MajorBrand {
	case brandJxl:
		return imagetype.ImageJXL
	case brandAvif, brandAvis:
		return imagetype.ImageAVIF
	case brandHeic, brandHeim, brandHeis, brandHeix, brandHevc, brandHevm, brandHevs, brandHevx:
		return imagetype.ImageHEIC
	case brandHeif, brandMiaf, brandMif1, brandMif2, brandMsf1:
		return imagetype.ImageHEIF
	case brandCrx:
		return imagetype.ImageCR3
	default:
		return imagetype.ImageHEIF
	}
}

func (r *Reader) readMeta(b *box) (err error) {
	if !b.isType(typeMeta) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
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
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
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
