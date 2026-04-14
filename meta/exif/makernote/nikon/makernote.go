package nikon

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// Nikon contains the selected Nikon maker-note fields currently decoded by
// imagemeta.
//
// The field set mirrors the subset of ExifTool's Image::ExifTool::Nikon::Main
// and related Nikon subdirectory tables that imagemeta parses today. Field
// names generally follow ExifTool's public tag names so side-by-side comparison
// with Nikon.pm and exiftool output is straightforward.
//
// This is intentionally not a complete representation of every Nikon maker-note
// tag. Fields are added here only when imagemeta has a parser path for them and
// when the exposed value has reasonably stable semantics across files.
type Nikon struct {
	// MakerNoteVersion is Nikon maker-note tag 0x0001.
	//
	// In Nikon.pm this is "MakerNoteVersion". ExifTool notes that older Nikon
	// cameras may store the 4-byte payload either as ASCII like "0211" or as
	// binary bytes that should be rendered as digits. imagemeta stores the
	// normalized 4-character version string.
	MakerNoteVersion string

	// ISO is Nikon maker-note tag 0x0002.
	//
	// In Nikon.pm this tag is named "ISO" and described as the ISO actually used
	// by the camera, which may differ from the EXIF ISO setting when Auto ISO is
	// enabled. Nikon stores this as a 2-word form where the leading word carries
	// mode information and the second word carries the visible ISO value. imagemeta
	// stores the effective numeric ISO value.
	ISO uint32

	// ISOSetting is Nikon maker-note tag 0x0013.
	//
	// This is a separate Nikon maker-note ISO field distinct from tag 0x0002.
	// ExifTool formats it similarly to Nikon ISO, and imagemeta stores the
	// effective numeric setting value.
	ISOSetting uint32

	// ColorMode is Nikon maker-note tag 0x0003.
	//
	// ExifTool treats this as a Nikon string value and applies a default Nikon
	// text formatter. imagemeta exposes the trimmed text payload.
	ColorMode string

	// Quality is Nikon maker-note tag 0x0004.
	//
	// Typical values include "RAW". Nikon frequently pads this field with spaces;
	// imagemeta trims trailing padding.
	Quality string

	// WhiteBalance is Nikon maker-note tag 0x0005.
	WhiteBalance string

	// Sharpness is Nikon maker-note tag 0x0006.
	//
	// Nikon.pm handles this as a string field. Some modern Nikon files store
	// values that are primarily meaningful when interpreted alongside Picture
	// Control state. imagemeta preserves the decoded string when present.
	Sharpness string

	// FocusMode is Nikon maker-note tag 0x0007.
	//
	// Nikon.pm records this before some AF-related conditional decode. imagemeta
	// currently exposes the raw trimmed text such as "AF-C".
	FocusMode string

	// FlashSetting is Nikon maker-note tag 0x0008.
	//
	// ExifTool notes this would be more precisely described as a flash sync mode.
	FlashSetting string

	// FlashType is Nikon maker-note tag 0x0009.
	//
	// ExifTool documentation notes that the observed value set varies by internal
	// versus optional flashes. imagemeta preserves the raw Nikon text.
	FlashType string

	// ISOSelection is Nikon maker-note tag 0x000f.
	ISOSelection string

	// SerialNumber is Nikon maker-note tags 0x001d and 0x00a0.
	//
	// Nikon.pm treats 0x001d as protected because it is also used as part of the
	// decryption key for some encrypted Nikon blocks. imagemeta stores the visible
	// serial string only; it does not expose decryption internals here.
	SerialNumber string

	// Lens is Nikon maker-note tag 0x0084.
	//
	// Nikon stores this as 4 rational values:
	//   short focal, long focal, max aperture at short focal, max aperture at long focal
	//
	// ExifTool formats this through Exif.pm PrintLensInfo, yielding strings like
	// "100 400 4.5 5.6". imagemeta stores the same human-readable normalized form.
	Lens string

	// PowerUpTime is Nikon maker-note tag 0x00b6.
	//
	// ExifTool describes this as the date/time when the camera was last powered
	// up. The exact meaning of "powered up" is Nikon-specific and may correspond
	// either to camera startup or power application. imagemeta stores the parsed
	// timestamp in UTC because the maker-note block itself does not encode a zone.
	PowerUpTime time.Time

	// ColorSpace is Nikon maker-note tag 0x001e.
	//
	// imagemeta keeps the raw Nikon code. ExifTool maps known values such as
	// 1=sRGB, 2=Adobe RGB, and 4=BT.2100.
	ColorSpace uint16

	// ActiveDLighting is Nikon maker-note tag 0x0022.
	//
	// Stored as the Nikon numeric code instead of ExifTool's printed label so the
	// value remains stable for JSON comparison and sorting.
	ActiveDLighting uint16

	// VignetteControl is Nikon maker-note tag 0x002a.
	VignetteControl uint16

	// ShutterMode is Nikon maker-note tag 0x0034.
	//
	// ExifTool prints this as labels such as Mechanical, Electronic, and
	// Electronic Front Curtain. imagemeta intentionally stores the raw Nikon code.
	ShutterMode uint16

	// ImageSizeRAW is Nikon maker-note tag 0x003e.
	//
	// ExifTool maps 1/2/3 to Large/Medium/Small. imagemeta stores the code.
	ImageSizeRAW uint16

	// ColorTemperatureAuto is Nikon maker-note tag 0x004f.
	//
	// This is the auto white-balance color temperature chosen by the camera in
	// Kelvin for files that contain it.
	ColorTemperatureAuto uint16

	// LensType is Nikon maker-note tag 0x0083.
	//
	// Nikon packs multiple lens capability bits into this byte. ExifTool contains
	// a PrintConv that expands the bitfield into labels such as MF, D, G, VR,
	// AF-P, E, and FT-1. imagemeta currently preserves the raw bitfield byte.
	LensType uint8

	// FlashMode is Nikon maker-note tag 0x0087.
	//
	// ExifTool maps values like 0, 1, 7, 8, 9, and 18 to descriptive flash mode
	// labels. imagemeta stores the underlying numeric code.
	FlashMode uint8

	// ShootingMode is Nikon maker-note tag 0x0089.
	//
	// Nikon uses this as a bitfield; ExifTool has extensive label logic around
	// release mode, bracketing, tethering, and Auto ISO. imagemeta currently
	// preserves the raw uint16 code.
	ShootingMode uint16

	// LensFStops is Nikon maker-note tag 0x008b.
	//
	// Nikon stores this as a 4-byte packed fractional encoding rather than a
	// standard EXIF rational. ExifTool converts it using the first 3 bytes:
	//   a * (b / c)
	//
	// imagemeta stores the converted numeric result.
	LensFStops float64

	// SilentPhotography is Nikon maker-note tag 0x00bf.
	//
	// ExifTool maps this as a simple Off/On boolean. imagemeta stores it as bool.
	SilentPhotography bool

	// MechanicalShutterCount is Nikon maker-note tag 0x0037.
	//
	// This is distinct from ShutterCount and records only mechanical actuation
	// count on cameras that expose both.
	MechanicalShutterCount uint32

	// ImageCount is Nikon maker-note tag 0x00a5.
	ImageCount uint32

	// DeletedImageCount is Nikon maker-note tag 0x00a6.
	DeletedImageCount uint32

	// ShutterCount is Nikon maker-note tag 0x00a7.
	//
	// Nikon.pm notes this value is also used in the keying of certain encrypted
	// maker-note records. imagemeta exposes only the visible count.
	ShutterCount uint32

	// ManualFocusDistance is Nikon maker-note tag 0x0085.
	//
	// Nikon stores this as a rational64u value. imagemeta preserves the exact
	// rational instead of formatting it to a display string.
	ManualFocusDistance tag.RationalU

	// DigitalZoom is Nikon maker-note tag 0x0086.
	//
	// Nikon stores this as a rational64u value. imagemeta preserves the exact
	// rational so callers can decide how to round or display it.
	DigitalZoom tag.RationalU

	// VRInfo is Nikon maker-note tag 0x001f interpreted using Nikon::VRInfo.
	VRInfo NikonVRInfo

	// WorldTime is Nikon maker-note tag 0x0024 interpreted using Nikon::WorldTime.
	WorldTime NikonWorldTime

	// ISOInfo is Nikon maker-note tag 0x0025 interpreted using Nikon::ISOInfo.
	ISOInfo NikonISOInfo

	// AFInfo is Nikon maker-note tag 0x0088 interpreted using Nikon::AFInfo.
	AFInfo NikonAFInfo

	// AFInfo2 is Nikon maker-note tag 0x00b7 interpreted using the appropriate
	// versioned AFInfo2 table in Nikon.pm.
	AFInfo2 NikonAFInfo2

	// FileInfo is Nikon maker-note tag 0x00b8 interpreted using Nikon::FileInfo.
	FileInfo NikonFileInfo

	// AFTune is Nikon maker-note tag 0x00b9 interpreted using Nikon::AFTune.
	AFTune NikonAFTune
}

// NikonVRInfo models ExifTool's Image::ExifTool::Nikon::VRInfo table.
//
// The table is a short binary subdirectory introduced for vibration reduction
// state and mode reporting. Nikon.pm warns that byte order must be set
// explicitly before adding multi-byte integers, but the fields imagemeta
// currently exposes are all single-byte status values except for the version
// string.
type NikonVRInfo struct {
	// VRInfoVersion is the 4-byte version field at offset 0.
	VRInfoVersion string

	// VibrationReduction is the raw byte at offset 4.
	//
	// ExifTool maps 0=n/a, 1=On, 2=Off.
	VibrationReduction uint8

	// VRMode is the raw byte at offset 6.
	//
	// The meaning varies slightly across Nikon generations. ExifTool uses
	// separate conversions for Z-series versus older models.
	VRMode uint8

	// VRType is the raw byte at offset 8.
	//
	// ExifTool maps values such as 2=In-body and 3=In-body + Lens.
	VRType uint8
}

// NikonWorldTime models ExifTool's Image::ExifTool::Nikon::WorldTime table.
//
// Nikon stores this as a tiny binary block containing timezone minutes and two
// one-byte state fields. Nikon software is known to sometimes flip the byte
// order of this structure, so imagemeta chooses byte order heuristically during
// parse and stores the resolved values here.
type NikonWorldTime struct {
	// TimeZone is the signed minute offset from UTC stored at offset 0.
	//
	// ExifTool formats this as "+/-HH:MM". imagemeta stores the canonical signed
	// minute count because it is lossless and easy to marshal differently later.
	TimeZone int16

	// DaylightSavings is the byte at offset 2.
	//
	// ExifTool treats 0 as No and 1 as Yes.
	DaylightSavings uint8

	// DateDisplayFormat is the byte at offset 3.
	//
	// ExifTool maps 0=Y/M/D, 1=M/D/Y, 2=D/M/Y.
	DateDisplayFormat uint8
}

// NikonISOInfo models ExifTool's Image::ExifTool::Nikon::ISOInfo table.
//
// Nikon.pm interprets ISO and ISO2 using the Nikon-specific formula:
//
//	ISO = 100 * 2**(raw/12 - 5)
//
// imagemeta stores the converted floating-point ISO values and preserves the
// raw expansion codes separately.
type NikonISOInfo struct {
	// ISO is the converted value from offset 0.
	ISO float64

	// ISOExpansion is the raw expansion code at offset 4.
	//
	// ExifTool maps values such as 0x101="Hi 0.3" and 0x201="Lo 0.3".
	ISOExpansion uint16

	// ISO2 is the converted value from offset 6.
	//
	// Nikon.pm notes that bytes 6-11 often duplicate 0-5 in available samples.
	ISO2 float64

	// ISOExpansion2 is the raw expansion code at offset 10.
	ISOExpansion2 uint16
}

// NikonAFInfo models ExifTool's Image::ExifTool::Nikon::AFInfo table.
//
// This is the older Nikon AF metadata block used before AFInfo2. ExifTool
// exposes AFAreaMode, AFPoint, and AFPointsInFocus. imagemeta keeps both the
// raw focus bitmask and the expanded focus-point indices.
type NikonAFInfo struct {
	// AFAreaMode is the raw byte at offset 0.
	AFAreaMode uint8

	// AFPoint is the raw byte at offset 1.
	//
	// ExifTool warns that this value is not meaningful in some focus modes.
	AFPoint uint8

	// AFPointsInFocusMask is the uint16 bitmask at offset 2.
	AFPointsInFocusMask uint16

	// AFPointsInFocus is the expanded set of point indices derived from
	// AFPointsInFocusMask.
	AFPointsInFocus []int
}

// NikonAFInfo2 models the subset of ExifTool's versioned Nikon AFInfo2 tables
// currently decoded by imagemeta.
//
// Nikon uses multiple AFInfo2 layouts selected by the first 4-byte version
// string. ExifTool dispatches to version-specific tables such as:
//   - Nikon::AFInfo2V0100
//   - Nikon::AFInfo2V0101
//   - Nikon::AFInfo2V0300
//   - Nikon::AFInfo2V0400
//
// imagemeta currently focuses on the fields required for parity across the
// local Nikon sample corpus, especially the version 0400/0401/0402 Expeed 7
// family and the legacy 0100/0101 DSLR family.
type NikonAFInfo2 struct {
	// AFInfo2Version is the 4-byte version field at offset 0.
	AFInfo2Version string

	// AFDetectionMethod is the byte at offset 4.
	//
	// ExifTool uses this to distinguish phase-detect versus contrast-detect AF.
	AFDetectionMethod uint8

	// AFAreaMode is the byte at offset 5.
	//
	// For modern Expeed 7 files, ExifTool notes that this is the active AF area
	// mode at shutter time rather than simply a menu-position value.
	AFAreaMode uint8

	// FocusPointSchema is the legacy schema byte at offset 6 in AFInfo2 v0100/
	// v0101 style blocks.
	//
	// ExifTool uses this to choose 11-point, 39-point, 51-point, or 153-point
	// point-grid decoding.
	FocusPointSchema uint8

	// AFCoordinatesAvailable is the Expeed 7 byte at offset 7.
	//
	// ExifTool documents 0 meaning AFPointsUsed is populated and 1 meaning
	// AFAreaXPosition/AFAreaYPosition are populated.
	AFCoordinatesAvailable uint8

	// PrimaryAFPoint is the main selected focus point for legacy AFInfo2 layouts.
	PrimaryAFPoint uint8

	// ContrastDetectAFInFocus is the v0100/v0101 contrast-detect flag at offset
	// 0x1c when phase-detect point maps are not being used.
	ContrastDetectAFInFocus bool

	// AFImageWidth and AFImageHeight are the pixel dimensions of the AF
	// coordinate space used by the camera-specific AF sub-block.
	AFImageWidth  uint16
	AFImageHeight uint16

	// AFAreaXPosition and AFAreaYPosition locate the AF area in AF-image
	// coordinates.
	AFAreaXPosition uint16
	AFAreaYPosition uint16

	// AFAreaWidth and AFAreaHeight encode the AF box size in AF-image
	// coordinates.
	AFAreaWidth  uint16
	AFAreaHeight uint16

	// FocusResult is the Expeed 7 byte at offset 0x4a.
	FocusResult uint8

	// AFPointsUsed is the decoded AF point-set reported by the AFInfo2 block.
	//
	// For legacy layouts this usually comes from a schema-dependent packed mask.
	// For Expeed 7 layouts it may come from camera-model-specific AF point maps.
	AFPointsUsed []int

	// AFPointsSelected is the decoded selected-point mask for layouts that expose
	// it, such as Nikon::AFInfo2V0101 153-point group/3D-tracking modes.
	AFPointsSelected []int

	// AFPointsInFocus is the decoded "in focus" point set for layouts that expose
	// it.
	AFPointsInFocus []int
}

// NikonFileInfo models ExifTool's Image::ExifTool::Nikon::FileInfo table.
//
// ExifTool notes that newer Nikon models may store this block little-endian,
// while other files remain big-endian, and Nikon desktop software may rewrite
// the surrounding maker-note byte order without rewriting this specific record.
// imagemeta resolves the byte order heuristically at parse time and stores the
// resulting values here.
type NikonFileInfo struct {
	// FileInfoVersion is the 4-byte version field at offset 0.
	FileInfoVersion string

	// MemoryCardNumber is the uint16 at index 2 in Nikon.pm's int16u table.
	MemoryCardNumber uint16

	// DirectoryNumber is the uint16 at index 3.
	//
	// ExifTool formats this as a zero-padded 3-digit string. imagemeta stores
	// the numeric value.
	DirectoryNumber uint16

	// FileNumber is the uint16 at index 4.
	//
	// ExifTool formats this as a zero-padded 4-digit string. imagemeta stores
	// the numeric value.
	FileNumber uint16
}

// NikonAFTune models ExifTool's Image::ExifTool::Nikon::AFTune table.
//
// This table carries Nikon AF fine-tune state for the current or remembered
// lens. ExifTool exposes these as small integer values with descriptive labels
// layered on top.
type NikonAFTune struct {
	// AFFineTune is the mode byte at offset 0.
	//
	// ExifTool maps this to Off, On (1), On (2), and On (Zoom).
	AFFineTune uint8

	// AFFineTuneIndex is the saved-lens index byte at offset 1.
	//
	// ExifTool treats 255 as "n/a".
	AFFineTuneIndex uint8

	// AFFineTuneAdj is the signed wide-end fine-tune adjustment at offset 2.
	AFFineTuneAdj int8

	// AFFineTuneAdjTele is the signed tele-end fine-tune adjustment at offset 3.
	//
	// ExifTool notes this is only valid for zoom lenses when AFFineTune indicates
	// zoom-lens tuning.
	AFFineTuneAdjTele int8
}
