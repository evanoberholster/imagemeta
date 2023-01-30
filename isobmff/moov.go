package isobmff

import (
	"fmt"

	"github.com/pkg/errors"
)

func (r *Reader) ReadMetadata() (err error) {
	b, err := r.readBox()
	if err != nil {
		return errors.Wrapf(err, "ReadMetadata")
	}
	switch b.boxType {
	case typeMdat:
		err = r.readMdat(&b)
	case typeMeta:
		err = r.readMeta(&b)
	case typeMoov:
		err = r.readMoovBox(&b)
	default:
		if logLevelInfo() {
			logInfoBox(b)
		}
	}
	return err
}

func (r *Reader) readMdat(b *box) (err error) {
	if logLevelInfo() {
		logInfoBox(*b)
	}
	err = b.readFlags()
	fmt.Println(b.flags, err)
	buf, err := b.Peek(128)
	fmt.Println(string(buf), buf, err)

	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		if logLevelInfo() {
			logInfoBox(inner)
		}
		b.close()
	}
	return b.close()
}

func (r *Reader) readMeta(b *box) (err error) {
	if !b.isType(typeMeta) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	if err = b.readFlags(); err != nil {
		return err
	}
	if logLevelInfo() {
		logInfoBox(*b)
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		//case typeUUID:
		case typeHdlr:
			_, err = readHdlr(&inner)
		case typePitm:
			_, err = readPitm(&inner)
		case typeIinf:
			_, err = readIinf(&inner)
		default:
			if logLevelInfo() {
				logInfoBox(inner)
			}
		}
		if err != nil {
			// log error
			return err
		}
		inner.close()
	}
	return b.close()
}

// ReadMOOV reads an 'moov' box from a BMFF file.
func (r *Reader) readMoovBox(b *box) (err error) {
	if !b.isType(typeMoov) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	if logLevelInfo() {
		logInfoBox(*b)
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		if logLevelInfo() {
			logInfoBox(inner)
		}
		switch inner.boxType {
		case typeUUID:
			uuid, err := inner.readUUID()
			if err != nil {
				return err
			}
			switch uuid {
			case CR3MetaBoxUUID:
				if _, err = readCrxMoovBox(&inner, r.ExifReader); err != nil {
					return err
				}
			default:
				inner.close()
			}
		case typeTrak:
			_, err = readCrxTrakBox(&inner)
		//case typeMvhd:
		default:
			inner.close()
		}
	}
	return b.close()
}
