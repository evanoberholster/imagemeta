// Package makernote defines maker-note tag IDs, camera-make identification,
// and the shared EXIF maker-note container for vendor-specific metadata.
//
// Vendor-specific typed models live in subpackages:
//
//   - makernote/canon
//   - makernote/nikon
//   - makernote/panasonic
//   - makernote/sony
//
// The package is intentionally conservative: only selected fields are modeled
// in the hot parse path, while unknown or unsupported makes can still be
// identified via CameraMake.
package makernote
