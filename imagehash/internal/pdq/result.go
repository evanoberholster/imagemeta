// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import "fmt"

type Result struct {
	Hash      Hash256
	Quality   int
	Dihedrals [8][16]uint16
}

func (r Result) String() string {
	return fmt.Sprintf("pdq{hash=%s quality=%d}", r.Hash.String(), r.Quality)
}
