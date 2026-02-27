package isobmff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

func TestReadMetadataCR3XPacketCallback(t *testing.T) {
	xpacketPayload := []byte("<?xpacket begin='\\ufeff' id='W5M0MpCehiHzreSzNTczkc9d'?>\n<x:xmpmeta></x:xmpmeta>")
	file := append(makeFTYP("crx "), makeUUIDBox(cr3XPacketUUID, xpacketPayload)...)

	var gotPayload []byte
	var gotHeader XPacketHeader

	r := NewReader(bytes.NewReader(file), nil, func(rr io.Reader, h XPacketHeader) error {
		var err error
		gotPayload, err = io.ReadAll(rr)
		gotHeader = h
		return err
	}, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error: %v", err)
	}
	if !bytes.Equal(gotPayload, xpacketPayload) {
		t.Fatalf("xpacket payload mismatch: got %q, want %q", string(gotPayload), string(xpacketPayload))
	}
	if !gotHeader.HasXPacketPI {
		t.Fatal("expected HasXPacketPI=true")
	}
	if !gotHeader.HasXMPMeta {
		t.Fatal("expected HasXMPMeta=true")
	}
	if gotHeader.Length != uint32(len(xpacketPayload)) {
		t.Fatalf("xpacket length = %d, want %d", gotHeader.Length, len(xpacketPayload))
	}
}

func TestReadMetadataCR3THMBPreviewCallback(t *testing.T) {
	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xD9}
	thmbPayload := make([]byte, 16+len(jpeg))
	binary.BigEndian.PutUint16(thmbPayload[4:6], 160)
	binary.BigEndian.PutUint16(thmbPayload[6:8], 120)
	binary.BigEndian.PutUint32(thmbPayload[8:12], uint32(len(jpeg)))
	copy(thmbPayload[16:], jpeg)
	thmb := makeBox("THMB", thmbPayload)

	cr3Meta := append(cr3MetaBoxUUID.Bytes(), thmb...)
	moov := makeBox("moov", makeBox("uuid", cr3Meta))
	file := append(makeFTYP("crx "), moov...)

	var gotHeader meta.PreviewHeader
	var gotPreview []byte

	r := NewReader(bytes.NewReader(file), nil, nil, func(rr io.Reader, h meta.PreviewHeader) error {
		gotHeader = h
		var err error
		gotPreview, err = io.ReadAll(rr)
		return err
	})
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error: %v", err)
	}
	if !bytes.Equal(gotPreview, jpeg) {
		t.Fatalf("preview payload mismatch: got=%x want=%x", gotPreview, jpeg)
	}
	if gotHeader.Source != meta.PreviewSourceTHMB {
		t.Fatalf("preview source = %s, want %s", gotHeader.Source, meta.PreviewSourceTHMB)
	}
	if gotHeader.ImageType != imagetype.ImageJPEG {
		t.Fatalf("preview image type = %v, want %v", gotHeader.ImageType, imagetype.ImageJPEG)
	}
	if gotHeader.Width != 160 || gotHeader.Height != 120 {
		t.Fatalf("preview dimensions = %dx%d, want 160x120", gotHeader.Width, gotHeader.Height)
	}
}

func TestReadMetadataCR3PRVWPreviewCallback(t *testing.T) {
	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xD9}
	prvwPayload := make([]byte, 16+len(jpeg))
	binary.BigEndian.PutUint32(prvwPayload[4:8], 0x00000140) // width=320 at bytes 6:8 in parse window
	binary.BigEndian.PutUint16(prvwPayload[8:10], 240)
	binary.BigEndian.PutUint16(prvwPayload[10:12], 2)
	binary.BigEndian.PutUint32(prvwPayload[12:16], uint32(len(jpeg)))
	copy(prvwPayload[16:], jpeg)
	prvw := makeBox("PRVW", prvwPayload)

	uuidPayload := append(cr3PreviewUUID.Bytes(), make([]byte, 8)...)
	uuidPayload[23] = 1
	uuidPayload = append(uuidPayload, prvw...)
	file := append(makeFTYP("crx "), makeBox("uuid", uuidPayload)...)

	var gotHeader meta.PreviewHeader
	var gotPreview []byte

	r := NewReader(bytes.NewReader(file), nil, nil, func(rr io.Reader, h meta.PreviewHeader) error {
		gotHeader = h
		var err error
		gotPreview, err = io.ReadAll(rr)
		return err
	})
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error: %v", err)
	}
	if !bytes.Equal(gotPreview, jpeg) {
		t.Fatalf("preview payload mismatch: got=%x want=%x", gotPreview, jpeg)
	}
	if gotHeader.Source != meta.PreviewSourcePRVW {
		t.Fatalf("preview source = %s, want %s", gotHeader.Source, meta.PreviewSourcePRVW)
	}
	if gotHeader.ImageType != imagetype.ImageJPEG {
		t.Fatalf("preview image type = %v, want %v", gotHeader.ImageType, imagetype.ImageJPEG)
	}
	if gotHeader.Width != 320 || gotHeader.Height != 240 {
		t.Fatalf("preview dimensions = %dx%d, want 320x240", gotHeader.Width, gotHeader.Height)
	}
}

func TestReadMetadataExifCallbackNonEOFErrorNonFatal(t *testing.T) {
	file := append(makeFTYP("avif"), makeBox("Exif", []byte{
		'I', 'I', 0x2A, 0x00,
		0x08, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})...)

	r := NewReader(bytes.NewReader(file), func(_ io.Reader, _ meta.ExifHeader) error {
		return errors.New("callback failed")
	}, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error: %v", err)
	}
	if r.hasHave(metadataKindExif) {
		t.Fatal("haveExif should remain false on non-EOF callback error")
	}
}

func TestReadMetadataExifCallbackEOFFatal(t *testing.T) {
	file := append(makeFTYP("avif"), makeBox("Exif", []byte{
		'I', 'I', 0x2A, 0x00,
		0x08, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})...)

	r := NewReader(bytes.NewReader(file), func(_ io.Reader, _ meta.ExifHeader) error {
		return io.EOF
	}, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	if err := r.ReadMetadata(); err != io.EOF {
		t.Fatalf("ReadMetadata() error: %v, want io.EOF", err)
	}
}

func TestReadMetadataPreviewCallbackUsesLimitedReader(t *testing.T) {
	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xD9}
	extra := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	thmbPayload := make([]byte, 16+len(jpeg)+len(extra))
	binary.BigEndian.PutUint16(thmbPayload[4:6], 160)
	binary.BigEndian.PutUint16(thmbPayload[6:8], 120)
	binary.BigEndian.PutUint32(thmbPayload[8:12], uint32(len(jpeg)))
	copy(thmbPayload[16:], jpeg)
	copy(thmbPayload[16+len(jpeg):], extra)
	thmb := makeBox("THMB", thmbPayload)

	cr3Meta := append(cr3MetaBoxUUID.Bytes(), thmb...)
	moov := makeBox("moov", makeBox("uuid", cr3Meta))
	file := append(makeFTYP("crx "), moov...)

	var gotPreview []byte
	r := NewReader(bytes.NewReader(file), nil, nil, func(rr io.Reader, _ meta.PreviewHeader) error {
		var err error
		gotPreview, err = io.ReadAll(rr)
		return err
	})
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	if err := r.ReadMetadata(); err != nil {
		t.Fatalf("ReadMetadata() error: %v", err)
	}
	if !bytes.Equal(gotPreview, jpeg) {
		t.Fatalf("preview payload mismatch: got=%x want=%x", gotPreview, jpeg)
	}
}

func TestReadMetadataHEIFXMPFromMdatCallback(t *testing.T) {
	xmpPayload := []byte("<?xpacket begin='\\ufeff'?><x:xmpmeta></x:xmpmeta>")

	infePayload := make([]byte, 0, 32)
	infePayload = append(infePayload, 0x02, 0x00, 0x00, 0x00) // version=2
	infePayload = append(infePayload, 0x00, 0x01)             // item_ID
	infePayload = append(infePayload, 0x00, 0x00)             // protection_index
	infePayload = append(infePayload, 'm', 'i', 'm', 'e')     // item_type
	infePayload = append(infePayload, 0x00)                   // item_name
	infePayload = append(infePayload, []byte("application/rdf+xml")...)
	infePayload = append(infePayload, 0x00) // content_type
	infe := makeBox("infe", infePayload)

	iinfPayload := make([]byte, 0, 8+len(infe))
	iinfPayload = append(iinfPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	iinfPayload = append(iinfPayload, 0x00, 0x01)             // entry_count
	iinfPayload = append(iinfPayload, infe...)
	iinf := makeBox("iinf", iinfPayload)

	ilocPayload := make([]byte, 0, 32)
	ilocPayload = append(ilocPayload, 0x00, 0x00, 0x00, 0x00)                         // version=0
	ilocPayload = append(ilocPayload, 0x44, 0x00)                                     // offset_size=4 length_size=4 base_offset_size=0
	ilocPayload = append(ilocPayload, 0x00, 0x01)                                     // item_count
	ilocPayload = append(ilocPayload, 0x00, 0x01)                                     // item_ID
	ilocPayload = append(ilocPayload, 0x00, 0x00)                                     // data_reference_index
	ilocPayload = append(ilocPayload, 0x00, 0x01)                                     // extent_count
	ilocPayload = append(ilocPayload, 0x00, 0x00, 0x00, 0x00)                         // extent_offset (relative to mdat payload start)
	ilocPayload = binary.BigEndian.AppendUint32(ilocPayload, uint32(len(xmpPayload))) // extent_length
	iloc := makeBox("iloc", ilocPayload)

	metaPayload := make([]byte, 0, 4+len(iinf)+len(iloc))
	metaPayload = append(metaPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	metaPayload = append(metaPayload, iinf...)
	metaPayload = append(metaPayload, iloc...)
	metaBox := makeBox("meta", metaPayload)

	file := append(makeFTYP("avif"), metaBox...)
	file = append(file, makeBox("mdat", xmpPayload)...)

	var gotPayload []byte
	var gotHeader XPacketHeader

	r := NewReader(bytes.NewReader(file), nil, func(rr io.Reader, h XPacketHeader) error {
		var err error
		gotPayload, err = io.ReadAll(rr)
		gotHeader = h
		return err
	}, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	for {
		err := r.ReadMetadata()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("ReadMetadata() error: %v", err)
		}
	}

	if !bytes.Equal(gotPayload, xmpPayload) {
		t.Fatalf("xmp payload mismatch: got=%q want=%q", gotPayload, xmpPayload)
	}
	if gotHeader.Length != uint32(len(xmpPayload)) {
		t.Fatalf("xmp length = %d, want %d", gotHeader.Length, len(xmpPayload))
	}
	if !r.hasHave(metadataKindXMP) {
		t.Fatal("expected haveXMP=true")
	}
}

func TestReadIrefParsesThmbReference(t *testing.T) {
	refPayload := []byte{
		0x00, 0x02, // from_item_id
		0x00, 0x01, // reference_count
		0x00, 0x01, // to_item_id
	}
	thmbRef := makeBox("thmb", refPayload)
	irefPayload := make([]byte, 0, 4+len(thmbRef))
	irefPayload = append(irefPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	irefPayload = append(irefPayload, thmbRef...)
	data := makeBox("iref", irefPayload)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error: %v", err)
	}
	if err := r.readIref(&b); err != nil {
		t.Fatalf("readIref() error: %v", err)
	}
	if len(r.heic.references) != 1 {
		t.Fatalf("reference count = %d, want 1", len(r.heic.references))
	}
	got := r.heic.references[0]
	if got.referenceType != typeThmb || got.fromID != 2 || got.toID != 1 {
		t.Fatalf("unexpected reference: %+v", got)
	}
}

func TestReadIprpParsesIspeAndIpma(t *testing.T) {
	ispePayload := make([]byte, 0, 12)
	ispePayload = append(ispePayload, 0x00, 0x00, 0x00, 0x00) // version=0
	ispePayload = binary.BigEndian.AppendUint32(ispePayload, 320)
	ispePayload = binary.BigEndian.AppendUint32(ispePayload, 240)
	ispe := makeBox("ispe", ispePayload)
	ipco := makeBox("ipco", ispe)

	ipmaPayload := make([]byte, 0, 16)
	ipmaPayload = append(ipmaPayload, 0x00, 0x00, 0x00, 0x00) // version=0, flags=0
	ipmaPayload = binary.BigEndian.AppendUint32(ipmaPayload, 1)
	ipmaPayload = append(ipmaPayload, 0x00, 0x07) // item_ID
	ipmaPayload = append(ipmaPayload, 0x01)       // association_count
	ipmaPayload = append(ipmaPayload, 0x01)       // property_index=1, essential=0
	ipma := makeBox("ipma", ipmaPayload)

	iprpPayload := append(ipco, ipma...)
	data := makeBox("iprp", iprpPayload)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	b, err := r.readBox()
	if err != nil {
		t.Fatalf("readBox() error: %v", err)
	}
	if err := r.readIprp(&b); err != nil {
		t.Fatalf("readIprp() error: %v", err)
	}
	if len(r.heic.properties) != 1 {
		t.Fatalf("property count = %d, want 1", len(r.heic.properties))
	}
	if got := r.heic.properties[0]; got.boxType != typeIspe || got.width != 320 || got.height != 240 {
		t.Fatalf("unexpected ispe property: %+v", got)
	}
	if len(r.heic.propertyLinks) != 1 {
		t.Fatalf("property link count = %d, want 1", len(r.heic.propertyLinks))
	}
	link := r.heic.propertyLinks[0]
	if link.itemID != 7 || link.propertyIndex != 1 || link.essential {
		t.Fatalf("unexpected property link: %+v", link)
	}
}

func TestReadMetadataHEIFPreviewFromThmbAndIpma(t *testing.T) {
	primaryData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	thumbData := []byte{0xAA, 0xBB, 0xCC}

	pitmPayload := []byte{
		0x00, 0x00, 0x00, 0x00, // version=0
		0x00, 0x01, // primary item id
	}
	pitm := makeBox("pitm", pitmPayload)

	infePrimary := makeBox("infe", []byte{
		0x02, 0x00, 0x00, 0x00, // version=2
		0x00, 0x01, // item_ID
		0x00, 0x00, // protection_index
		'h', 'v', 'c', '1', // item_type
		0x00, // item_name
	})
	infeThumb := makeBox("infe", []byte{
		0x02, 0x00, 0x00, 0x00, // version=2
		0x00, 0x02, // item_ID
		0x00, 0x00, // protection_index
		'h', 'v', 'c', '1', // item_type
		0x00, // item_name
	})
	iinfPayload := make([]byte, 0, 6+len(infePrimary)+len(infeThumb))
	iinfPayload = append(iinfPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	iinfPayload = append(iinfPayload, 0x00, 0x02)             // entry_count
	iinfPayload = append(iinfPayload, infePrimary...)
	iinfPayload = append(iinfPayload, infeThumb...)
	iinf := makeBox("iinf", iinfPayload)

	thmbRef := makeBox("thmb", []byte{
		0x00, 0x02, // from thumbnail
		0x00, 0x01, // ref count
		0x00, 0x01, // to primary
	})
	irefPayload := make([]byte, 0, 4+len(thmbRef))
	irefPayload = append(irefPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	irefPayload = append(irefPayload, thmbRef...)
	iref := makeBox("iref", irefPayload)

	ispePrimary := makeBox("ispe", binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32([]byte{
		0x00, 0x00, 0x00, 0x00, // version=0
	}, 4000), 3000))
	ispeThumb := makeBox("ispe", binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32([]byte{
		0x00, 0x00, 0x00, 0x00, // version=0
	}, 200), 120))
	ipco := makeBox("ipco", append(ispePrimary, ispeThumb...))

	ipmaPayload := make([]byte, 0, 32)
	ipmaPayload = append(ipmaPayload, 0x00, 0x00, 0x00, 0x00) // version=0 flags=0
	ipmaPayload = binary.BigEndian.AppendUint32(ipmaPayload, 2)
	ipmaPayload = append(ipmaPayload,
		0x00, 0x01, // item 1
		0x01,       // assoc count
		0x01,       // property 1
		0x00, 0x02, // item 2
		0x01, // assoc count
		0x02, // property 2
	)
	ipma := makeBox("ipma", ipmaPayload)
	iprp := makeBox("iprp", append(ipco, ipma...))

	ilocPayload := make([]byte, 0, 64)
	ilocPayload = append(ilocPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	ilocPayload = append(ilocPayload, 0x44, 0x00)             // offset=4 length=4 base=0
	ilocPayload = append(ilocPayload, 0x00, 0x02)             // item_count
	// item 1 (primary): offset after thumbnail data
	ilocPayload = append(ilocPayload, 0x00, 0x01) // item_ID
	ilocPayload = append(ilocPayload, 0x00, 0x00) // data_reference_index
	ilocPayload = append(ilocPayload, 0x00, 0x01) // extent_count
	ilocPayload = binary.BigEndian.AppendUint32(ilocPayload, uint32(len(thumbData)))
	ilocPayload = binary.BigEndian.AppendUint32(ilocPayload, uint32(len(primaryData)))
	// item 2 (thumbnail): starts at mdat payload beginning
	ilocPayload = append(ilocPayload, 0x00, 0x02) // item_ID
	ilocPayload = append(ilocPayload, 0x00, 0x00) // data_reference_index
	ilocPayload = append(ilocPayload, 0x00, 0x01) // extent_count
	ilocPayload = binary.BigEndian.AppendUint32(ilocPayload, 0)
	ilocPayload = binary.BigEndian.AppendUint32(ilocPayload, uint32(len(thumbData)))
	iloc := makeBox("iloc", ilocPayload)

	metaPayload := make([]byte, 0, 4+len(pitm)+len(iinf)+len(iref)+len(iprp)+len(iloc))
	metaPayload = append(metaPayload, 0x00, 0x00, 0x00, 0x00) // version=0
	metaPayload = append(metaPayload, pitm...)
	metaPayload = append(metaPayload, iinf...)
	metaPayload = append(metaPayload, iref...)
	metaPayload = append(metaPayload, iprp...)
	metaPayload = append(metaPayload, iloc...)
	metaBox := makeBox("meta", metaPayload)

	mdatPayload := append(append([]byte{}, thumbData...), primaryData...)
	file := append(makeFTYP("heic"), metaBox...)
	file = append(file, makeBox("mdat", mdatPayload)...)

	var gotHeader meta.PreviewHeader
	var gotPreview []byte
	r := NewReader(bytes.NewReader(file), nil, nil, func(rr io.Reader, h meta.PreviewHeader) error {
		gotHeader = h
		var err error
		gotPreview, err = io.ReadAll(rr)
		return err
	})
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error: %v", err)
	}
	for {
		err := r.ReadMetadata()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("ReadMetadata() error: %v", err)
		}
	}

	if !bytes.Equal(gotPreview, thumbData) {
		t.Fatalf("preview payload mismatch: got=%x want=%x", gotPreview, thumbData)
	}
	if gotHeader.Width != 200 || gotHeader.Height != 120 {
		t.Fatalf("preview dimensions = %dx%d, want 200x120", gotHeader.Width, gotHeader.Height)
	}
	if gotHeader.ImageType != imagetype.ImageHEIC {
		t.Fatalf("preview image type = %v, want %v", gotHeader.ImageType, imagetype.ImageHEIC)
	}
	if gotHeader.Source != meta.PreviewSourcePRVW {
		t.Fatalf("preview source = %s, want %s", gotHeader.Source, meta.PreviewSourcePRVW)
	}
}

func makeFTYP(major string) []byte {
	payload := make([]byte, 8)
	copy(payload[:4], []byte(major))
	copy(payload[4:8], []byte("0001"))
	return makeBox("ftyp", payload)
}

func makeUUIDBox(uuid meta.UUID, payload []byte) []byte {
	uuidPayload := append(uuid.Bytes(), payload...)
	return makeBox("uuid", uuidPayload)
}

func makeBox(boxType string, payload []byte) []byte {
	out := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(out[:4], uint32(len(out)))
	copy(out[4:8], []byte(boxType))
	copy(out[8:], payload)
	return out
}
