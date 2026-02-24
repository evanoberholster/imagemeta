package isobmff

import "fmt"

// pitmID is a "pitm" box.
//
// Primary Item Reference pitm allows setting one image as the primary item.
// 0 represents not set.

func readPitm(b *box) (id itemID, err error) {
	if err = b.readFlags(); err != nil {
		return invalidItemID, err
	}
	switch b.flags.version() {
	case 0:
		v, readErr := b.readUint16()
		if readErr != nil {
			return invalidItemID, readErr
		}
		id = itemID(v)
	case 1:
		v, readErr := b.readUint32()
		if readErr != nil {
			return invalidItemID, readErr
		}
		id = itemID(v)
	default:
		return invalidItemID, fmt.Errorf("readPitm: unsupported version %d", b.flags.version())
	}
	if logLevelInfo() {
		logInfoBox(b).Uint32("ptim", uint32(id)).Send()
	}
	return id, b.close()
}

// itemID
type itemID uint32

const invalidItemID itemID = 0
