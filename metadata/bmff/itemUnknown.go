package bmff

import "fmt"

type UnknownBox struct {
	t BoxType
	s int64
}

func (ub UnknownBox) Type() BoxType {
	return ub.t
}

func (ub UnknownBox) Size() int64 {
	return ub.s
}

func (ub UnknownBox) String() string {
	return fmt.Sprintf(" Type: %s, Size: %d", ub.t, ub.s)
}
