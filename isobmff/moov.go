package isobmff

import (
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func (r *Reader) ReadMetadata() (err error) {
	for {
		b, readErr := r.readBox()
		if readErr != nil {
			if readErr == io.EOF {
				return io.EOF
			}
			return fmt.Errorf("ReadMetadata: %w", readErr)
		}
		switch b.boxType {
		case typeMdat:
			err = r.readMdat(&b)
		case typeExif:
			err = r.readExif(&b)
		case typeMeta:
			err = r.readMeta(&b)
		case typeMoov:
			err = r.readMoovBox(&b)
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
			err = b.close()
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
	if r.heic.exif.ol.offset == 0 || r.heic.exif.ol.length == 0 {
		return b.close()
	}
	inner, err := r.newExifBox(b)
	if err != nil {
		if logLevelError() {
			logError().Object("box", inner).Err(err).Send()
		}
		return b.close()
	}
	if err = seekExifTIFFHeader(&inner); err != nil {
		if logLevelDebug() {
			logDebug().Object("box", inner).Err(err).Msg("skip non-TIFF mdat exif candidate")
		}
		return b.close()
	}
	imageType := r.metadataImageType()
	header, err := readExifHeader(&inner, ifds.IFD0, imageType)
	if err != nil {
		if logLevelDebug() {
			logDebug().Object("box", inner).Err(err).Msg("skip invalid mdat exif header")
		}
		return b.close()
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
	if r.heic.exif.ol.length == 0 {
		return inner, ErrBufLength
	}

	currentOffset := b.offset + int(b.size) - b.remain
	targetOffset := r.heic.exif.ol.offset

	var discardBytes uint64
	if targetOffset >= uint64(currentOffset) {
		discardBytes = targetOffset - uint64(currentOffset)
	} else {
		// Some encoders store extent offsets relative to mdat payload start.
		discardBytes = targetOffset
	}
	if discardBytes > uint64(b.remain) {
		return inner, ErrRemainLengthInsufficient
	}
	if _, err = b.Discard(int(discardBytes)); err != nil {
		return inner, err
	}

	if r.heic.exif.ol.length > uint64(b.remain) {
		return inner, ErrRemainLengthInsufficient
	}
	maxInt := int64(^uint(0) >> 1)
	if r.heic.exif.ol.length > uint64(maxInt) {
		return inner, errLargeBox
	}

	size := int64(r.heic.exif.ol.length)

	inner = box{
		reader:  b.reader,
		outer:   b,
		boxType: typeExif,
		offset:  int(b.size) - b.remain + b.offset,
		size:    size,
		remain:  int(size),
	}
	return inner, nil
}

func readExifHeader(b *box, firstIfd ifds.IfdType, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		err = fmt.Errorf("readExifHeader: %w", err)
		return
	}
	endian := utils.BinaryOrder(buf[:4])
	if endian == utils.UnknownEndian {
		return header, ErrBufLength
	}
	header = meta.NewExifHeader(endian, endian.Uint32(buf[4:8]), 0, uint32(b.remain), it)
	header.FirstIfd = firstIfd
	if logLevelInfo() {
		logInfo().Object("box", b).Object("header", header).Send()
	}
	_, err = b.Discard(8)
	return header, err
}

func seekExifTIFFHeader(b *box) error {
	for {
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
			continue
		}

		// Some payloads are wrapped in a local Exif box header.
		if hasEmbeddedExifBoxHeader(buf) {
			boxSize := bmffEndian.Uint32(buf[:4])
			if boxSize < 8 || int(boxSize) > b.remain {
				return ErrBufLength
			}
			if _, err = b.Discard(8); err != nil {
				return err
			}
			continue
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
}

func hasEmbeddedExifBoxHeader(buf []byte) bool {
	return len(buf) >= 8 &&
		buf[4] == 'E' &&
		buf[5] == 'x' &&
		buf[6] == 'i' &&
		buf[7] == 'f'
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
		if err != nil {
			if logLevelError() {
				logError().Object("box", inner).Err(err).Send()
			}
			return err
		}

		if err = inner.close(); err != nil {
			if logLevelError() {
				logError().Object("box", inner).Err(err).Send()
			}
			return err
		}
	}
	if err != nil {
		return err
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
		if err != nil {
			if logLevelError() {
				logError().Object("box", inner).Err(err).Send()
			}
			return err
		}
		if err = inner.close(); err != nil {
			if logLevelError() {
				logError().Object("box", inner).Err(err).Send()
			}
			return err
		}
	}
	if err != nil {
		return err
	}
	return b.close()
}
