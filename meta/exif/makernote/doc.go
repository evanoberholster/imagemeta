// Package makernote defines maker-note tag IDs, camera-make identification,
// and typed containers for vendor-specific EXIF maker-note metadata.
//
// The package is intentionally conservative: only selected Apple, Canon, and
// Nikon fields are modeled in the hot parse path, while unknown/unsupported
// makes can still be identified via CameraMake.
package makernote
