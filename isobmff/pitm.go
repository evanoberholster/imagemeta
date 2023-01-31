package isobmff

import "github.com/rs/zerolog"

// pitmID is a "pitm" box.
//
// Primary Item Reference pitm allows setting one image as the primary item.
// -1 represents not set.

func readPitm(b *box) (id itemID, err error) {
	buf, err := b.Peek(b.remain)
	if err != nil {
		return -1, err
	}
	b.readFlagsFromBuf(buf)
	id = itemID(bmffEndian.Uint16(buf[4:]))
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Uint16("ptim", uint16(id)).Send()
	}
	return id, b.close()
}

// itemID
type itemID int16
