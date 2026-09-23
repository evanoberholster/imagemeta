// Copyright (c) 2018-2023 Evan Oberholster. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// Package jpeg reads metadata information (Exif and XMP) from a JPEG Image.
//
// The scanner walks JPEG markers and dispatches APP segments: APP0 carries
// JFIF and Canon CIFF records, APP1 Exif/XMP (including extended XMP
// reassembly), APP2 ICC profiles and MPF, APP13 Photoshop/IPTC and APP14
// Adobe data. ScanMetadata additionally reports SOF image dimensions.
//
// Header scanning is bounded (4 MiB metadata budget, 64 MiB extended-XMP
// cap); malformed streams fail fast with an error, never hang. When only
// EXIF is requested the scan returns right after decoding it instead of
// walking trailing markers. Passing an io.ReaderAt (or using ScanBytes)
// lets segment payloads read independently with seekable callbacks.
package jpeg

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

// Errors
var (
	ErrNoExif            = meta.ErrNoExif
	ErrNoJPEGMarker      = errors.New("no JPEG Marker")
	ErrEndOfImage        = errors.New("end of Image")
	ErrInvalidMarkerSize = errors.New("invalid jpeg marker size")
	ErrMetadataScanLimit = errors.New("jpeg metadata scan limit exceeded")
)

// ScanJPEGContext scans a reader for JPEG Image markers with cancellation and
// scan limits. exifReader and xmpReader are run at their respective positions
// during the scan.
//
// Returns ErrNoJPEGMarker if a JPEG SOF was not found.
func ScanJPEGContext(ctx context.Context, r io.Reader, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	var readerAt io.ReaderAt
	if ra, ok := r.(io.ReaderAt); ok {
		readerAt = ra
	}
	return scanJPEG(ctx, r, readerAt, exifReader, xmpReader)
}

// ScanJPEG scans a reader for JPEG Image markers. exifReader and xmpReader are run at their respective
// positions during the scan.
//
// Returns the error ErrNoJPEGMarker if a JPEG SOF was not found.
func ScanJPEG(r io.Reader, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	return ScanJPEGContext(context.Background(), r, exifReader, xmpReader)
}

// ScanJPEGWithReaderAtContext scans JPEG markers using r as the forward stream and
// readerAt for independent segment reads. readerAt must use the same byte
// offsets as r. This lets metadata callbacks read from a section without moving
// the forward JPEG scanner.
func ScanJPEGWithReaderAtContext(ctx context.Context, r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	return scanJPEG(ctx, r, readerAt, exifReader, xmpReader)
}

// ScanJPEGWithReaderAt scans JPEG markers using r as the forward stream and
// readerAt for independent segment reads. readerAt must use the same byte
// offsets as r. This lets metadata callbacks read from a section without moving
// the forward JPEG scanner.
func ScanJPEGWithReaderAt(r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	return ScanJPEGWithReaderAtContext(context.Background(), r, readerAt, exifReader, xmpReader)
}

// ScanJPEGWithSourceContext scans JPEG markers from stream and uses source for
// independent segment reads when source implements io.ReaderAt.
func ScanJPEGWithSourceContext(ctx context.Context, stream io.Reader, source io.Reader, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) error {
	if ra, ok := source.(io.ReaderAt); ok {
		return ScanJPEGWithReaderAtContext(ctx, stream, ra, exifReader, xmpReader)
	}
	return ScanJPEGContext(ctx, stream, exifReader, xmpReader)
}

// ScanJPEGWithSource scans JPEG markers from stream and uses source for
// independent segment reads when source implements io.ReaderAt.
func ScanJPEGWithSource(stream io.Reader, source io.Reader, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) error {
	return ScanJPEGWithSourceContext(context.Background(), stream, source, exifReader, xmpReader)
}

// ScanBytesContext scans JPEG markers in an in-memory buffer. The buffer
// backs both the forward stream and the independent segment reads, so
// segment payloads never copy and callbacks receive seekable readers.
func ScanBytesContext(ctx context.Context, buf []byte, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) error {
	r := bytes.NewReader(buf)
	return scanJPEGWithMetadata(ctx, r, r, exifReader, xmpReader, nil)
}

// ScanBytes scans JPEG markers in an in-memory buffer. See ScanBytesContext.
func ScanBytes(buf []byte, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) error {
	return ScanBytesContext(context.Background(), buf, exifReader, xmpReader)
}

// ScanMetadata scans a JPEG stream and returns metadata stored directly in JPEG
// marker segments, such as JFIF, CIFF, MPF, ICC, Photoshop/IPTC, Adobe APP14 and
// SOF image dimensions. Photoshop resources parse selectively: only decoded
// fields materialize (thumbnail pixels contribute just their length), so
// large APP13 segments cost a fraction of their size.
func ScanMetadata(r io.Reader) (Metadata, error) {
	var readerAt io.ReaderAt
	if ra, ok := r.(io.ReaderAt); ok {
		readerAt = ra
	}
	return scanMetadata(r, readerAt)
}

// ScanMetadataWithReaderAt scans JPEG marker metadata using r as the forward
// stream and readerAt for independent segment reads.
func ScanMetadataWithReaderAt(r io.Reader, readerAt io.ReaderAt) (Metadata, error) {
	return scanMetadata(r, readerAt)
}
