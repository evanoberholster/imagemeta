package isobmff

import "github.com/rs/zerolog"

// pitmID is a "pitm" box.
//
// Primary Item Reference pitm allows setting one image as the primary item.
// -1 represents not set.
type pitmID int16

func readPitm(b *box) (id pitmID, err error) {
	if err = b.readFlags(); err != nil {
		return -1, err
	}
	i, err := b.readUint16()
	if err != nil {
		return -1, err
	}
	if logLevelInfo() {
		logInfoBoxExt(b, zerolog.InfoLevel).Uint16("ptim", i).Send()
	}
	return pitmID(i), b.close()
}
