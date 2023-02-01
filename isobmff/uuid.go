package isobmff

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	// CR3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
	CR3MetaBoxUUID = meta.UUIDFromString("85c0b687-820f-11e0-8111-f4ce462b6a48")

	// CR3XPacketUUID is the uuid that corresponds with Canon CR3 xpacket data
	CR3XPacketUUID = meta.UUIDFromString("be7acfcb-97a9-42e8-9c71-999491e3afac")
)

var uuidMap = map[meta.UUID]string{
	CR3MetaBoxUUID: "CR3MetaBoxUUID",
	CR3XPacketUUID: "CR3XPacketUUID",
}

func (r *Reader) readUUIDBox(b *box) error {
	if !b.isType(typeUUID) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	uuid, err := b.readUUID()
	if err != nil {
		return err
	}
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Str("uuid", uuid.String()).Send()
	}
	switch uuid {
	case CR3XPacketUUID:
		if r.XMPReader != nil {
			if err = r.XMPReader(b); err != nil {
				b.close()
				return err
			}
		}
	case CR3MetaBoxUUID:
		if _, err = readCrxMoovBox(b, r.ExifReader); err != nil {
			return err
		}
	default:
		if logLevelDebug() {
			logBoxExt(b, zerolog.DebugLevel).Send()
		}
	}
	return b.close()
}

func readUuid(b *box) (err error) {
	if !b.isType(typeUUID) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	uuid, err := b.readUUID()
	if err != nil {
		return err
	}
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Str("uuid", uuid.String()).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeCNCV:
			_, err = readCNCVBox(&inner)
		default:
			if logLevelDebug() {
				logBoxExt(b, zerolog.DebugLevel).Send()
			}
		}
		if err != nil && logLevelError() {
			logBoxExt(&inner, zerolog.ErrorLevel).Err(err).Send()
		}
		if err = inner.close(); err != nil {
			break
		}
	}
	return b.close()
}
