package isobmff

import "fmt"

type block [64]uint32

var data []uint32
var offset uint32

func build(o map[string]boxType) {
	data = make([]uint32, 64)
	for str, _ := range o {
		v := hash(str[0])
		a := data[v]
		if a == 0 {
			data[v] = offset
			offset += 64
			a = offset
			fmt.Println(string(str[0]), offset)
		}
		//v = rel[str[1]]
		//data[a+v]
	}
	fmt.Println(offset)
}

func hash(b byte) uint8 {
	return hashMap[b]
}

var hashMap = [256]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0, 0, 0, 0, 0, 0,
	0, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
	26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 0, 0, 0, 0, 0,
	0, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 0, 0, 0, 0, 0,
}

func (b *box) peekRemainingBox() {
	buf, err := b.Peek(512)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf, len(buf))
	fmt.Println(string(buf))
}
