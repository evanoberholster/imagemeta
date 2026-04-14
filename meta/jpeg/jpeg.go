// Copyright (c) 2018-2023 Evan Oberholster. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// Package jpeg reads metadata information (Exif and XMP) from a JPEG Image.
package jpeg

import (
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

// Errors
var (
	ErrNoExif       = meta.ErrNoExif
	ErrNoJPEGMarker = errors.New("no JPEG Marker")
	ErrEndOfImage   = errors.New("end of Image")
)

// ScanJPEG scans a reader for JPEG Image markers. exifReader and xmpReader are run at their respective
// positions during the scan.
//
// Returns the error ErrNoJPEGMarker if a JPEG SOF was not found.
func ScanJPEG(r io.Reader, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	var readerAt io.ReaderAt
	if ra, ok := r.(io.ReaderAt); ok {
		readerAt = ra
	}
	return scanJPEG(r, readerAt, exifReader, xmpReader)
}

// ScanJPEGWithReaderAt scans JPEG markers using r as the forward stream and
// readerAt for independent segment reads. readerAt must use the same byte
// offsets as r. This lets metadata callbacks read from a section without moving
// the forward JPEG scanner.
func ScanJPEGWithReaderAt(r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	return scanJPEG(r, readerAt, exifReader, xmpReader)
}

// ScanMetadata scans a JPEG stream and returns metadata stored directly in JPEG
// marker segments, such as JFIF, CIFF, MPF, ICC, Photoshop/IPTC, Adobe APP14 and
// SOF image dimensions.
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
