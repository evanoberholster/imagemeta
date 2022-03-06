package exif

// TODO: Write tests for exifReader
//func TestExifReader(t *testing.T) {
//	exifOffset := uint32(0)
//	byteOrder := binary.BigEndian
//	reader := bytes.NewReader([]byte{0, 0, 0, 0})
//	header := meta.ExifHeader{ByteOrder: byteOrder, TiffHeaderOffset: exifOffset}
//
//	r := newReader(reader, header)
//
//	// Error ExifReader
//	tempbuf := make([]byte, 0)
//	if n, err := r.Read(tempbuf); err != nil && n != 0 {
//		t.Errorf("Wanted Exif Read Error %s", err)
//	}
//	if _, err := r.ReadAt(tempbuf, -1); err != ErrReadNegativeOffset {
//		t.Errorf("Error reader.ReadAt negative offset %s", err)
//	}
//
//	// ByteOrder
//	if r.byteOrder != binary.BigEndian {
//		t.Errorf("Error with ByteOrder")
//	}
//
//	// TODO: test Reader
//}
