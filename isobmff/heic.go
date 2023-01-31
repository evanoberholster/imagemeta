package isobmff

type HeicMeta struct {
	pitm itemID
	idat idat
	exif item
	xml  item
	// irot
}

type item struct {
	id itemID
	ol offsetLength
}
