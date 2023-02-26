package isobmff

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

var (
	// cr3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
	cr3MetaBoxUUID = meta.UUIDFromString("85c0b687-820f-11e0-8111-f4ce462b6a48")

	// cr3XPacketUUID is the uuid that corresponds with Canon CR3 xpacket data
	cr3XPacketUUID = meta.UUIDFromString("be7acfcb-97a9-42e8-9c71-999491e3afac")
)

func (r *Reader) readUUIDBox(b *box) error {
	if !b.isType(typeUUID) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	uuid, err := b.readUUID()
	if err != nil {
		return err
	}
	if logLevelInfo() {
		logInfoBox(b).Str("uuid", uuid.String()).Send()
	}
	switch uuid {
	case cr3XPacketUUID:
		if r.XMPReader != nil {
			if err = r.XMPReader(b); err != nil {
				b.close()
				return err
			}
		}
	case cr3MetaBoxUUID:
		if _, err = readCrxMoovBox(b, r.ExifReader); err != nil {
			return err
		}
	default:
		if logLevelDebug() {
			logDebug().Object("box", b).Send()
		}
	}
	return b.close()
}
