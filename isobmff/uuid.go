package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta"
)

var (
	// cr3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
	cr3MetaBoxUUID = meta.UUIDFromString("85c0b687-820f-11e0-8111-f4ce462b6a48")

	// cr3XPacketUUID is the uuid that corresponds with Canon CR3 xpacket data
	cr3XPacketUUID = meta.UUIDFromString("be7acfcb-97a9-42e8-9c71-999491e3afac")

	// cr3PreviewUUID is the uuid that corresponds with Canon CR3 Preview Image.
	cr3PreviewUUID = meta.UUIDFromString("eaf42b5e-1c98-4b88-b9fb-b7dc406e4d16")
)

func (r *Reader) readUUIDBox(b *box) error {
	if !b.isType(typeUUID) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
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
		if r.xmpReader != nil {
			if err = r.xmpReader(b); err != nil {
				b.close()
				return err
			}
		}
	case cr3MetaBoxUUID:
		if _, err = readCrxMoovBox(b, r.exifReader); err != nil {
			return err
		}
	case cr3PreviewUUID:
		if err = r.readPreview(b); err != nil {
			return err
		}
	default:
		if logLevelDebug() {
			logDebug().Object("box", b).Send()
		}
	}
	return b.close()
}
