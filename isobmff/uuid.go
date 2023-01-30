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

func (r *Reader) ReadUUIDBox() error {
	b, err := r.readBox()
	if err != nil {
		return errors.Wrapf(err, "ReadMOOVBox")
	}
	if !b.isType(typeUUID) {
		return errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	uuid, err := b.readUUID()
	if err != nil {
		return err
	}
	if logLevelInfo() {
		logInfoBoxExt(&b, zerolog.InfoLevel).Str("uuid", uuid.String()).Send()
	}
	switch uuid {
	case CR3XPacketUUID:
		if r.XMPReader != nil {
			if err = r.XMPReader(&b); err != nil {
				b.close()
				return err
			}
		}
	}
	b.close()
	return err
}
