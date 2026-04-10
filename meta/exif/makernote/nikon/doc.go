// Package nikon contains Nikon-specific metadata models shared by the EXIF
// decoder and downstream callers.
//
// The primary source for the maker-note structures in this package is the local
// ExifTool Nikon implementation:
//
//   - /usr/share/perl5/Image/ExifTool/Nikon.pm
//
// Nikon maker notes are vendor-defined TIFF-like records embedded inside EXIF
// MakerNote (tag 0x927c). Unlike baseline EXIF IFDs defined by the CIPA DC-008
// EXIF/TIFF specification, Nikon maker-note layouts are not standardized across
// vendors and often vary by camera generation, processor family, and even Nikon
// desktop software rewrites. ExifTool is therefore used here as the practical
// reference for field names, offsets, version switches, and conditional decode
// rules.
//
// This package intentionally models the typed Nikon values that imagemeta
// currently parses and exposes. It is not a complete re-implementation of
// Nikon.pm. In particular:
//
//   - The maker-note model is a selected subset of Nikon::Main plus a few
//     structured subdirectories such as VRInfo, WorldTime, ISOInfo, AFInfo,
//     AFInfo2, FileInfo, and AFTune.
//   - Unknown or currently unsupported Nikon tags remain outside this package.
//   - A number of values are preserved as raw Nikon numeric codes because that
//     is the most stable representation for JSON, sorting, and parity checks.
//     Human-readable label mapping can be layered on top by callers.
//   - Some values are normalized compared with ExifTool's raw storage:
//     for example Lens is exposed using an ExifTool-style PrintLensInfo text
//     form, ISOInfo ISO/ISO2 are converted from Nikon's raw logarithmic ISO
//     encoding, and AF point bitmasks are expanded into point-index slices.
//
// The Nikon camera model table in this package is derived from ExifTool's Nikon
// sources plus Nikon model naming conventions. ExifTool does not expose a
// single Nikon model-ID table comparable to Canon's canonModelID mappings, so
// this package keeps Nikon model names as string-oriented identifiers rather
// than a maker-note numeric model ID.
package nikon
