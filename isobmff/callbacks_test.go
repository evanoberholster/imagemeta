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
	if r.haveExif {
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
