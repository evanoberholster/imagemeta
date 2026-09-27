// Package tags holds ExifTool-derived tag tables as build-time YAML data.
//
// The YAML files in this directory (exif.yaml, gps.yaml, canon.yaml,
// nikon.yaml, sony.yaml, panasonic.yaml, apple.yaml) are simplified
// projections of ExifTool %Main tables (see README.md for versions and
// regeneration). They are not loaded at runtime; a future cmd/taggen step
// will generate plain Go switch dispatch from them, preserving the
// zero-allocation hot path in meta/exif.
package tags
