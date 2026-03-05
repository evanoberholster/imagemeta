package xmp

import "time"

// CRS is Camera Raw Settings. Photoshop Camera Raw namespace tags.
//
//	xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
//
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#crs
type CRS struct {
	// RawFileName is the sidecar target filename (for example IMG_0001.CR2).
	RawFileName string
	// WhiteBalance is the Camera Raw white balance mode (for example "As Shot").
	WhiteBalance string
	// ColorTemperature is the effective white-balance temperature in Kelvin.
	ColorTemperature int16
	// Contrast is the legacy Camera Raw contrast adjustment.
	Contrast int16
	// Saturation is the legacy Camera Raw saturation adjustment.
	Saturation int16
	// Sharpness is the legacy Camera Raw sharpening amount.
	Sharpness int16
	// AlreadyApplied indicates whether lens/profile corrections were already applied.
	AlreadyApplied bool
}

func (crs *CRS) parse(p property) (err error) {
	switch p.Name() {
	case RawFileName:
		crs.RawFileName = parseString(p.Value())
	case WhiteBalance:
		crs.WhiteBalance = parseString(p.Value())
	case Temperature:
		crs.ColorTemperature = parseInt16(p.Value())
	case Contrast:
		crs.Contrast = parseInt16(p.Value())
	case Saturation:
		crs.Saturation = parseInt16(p.Value())
	case Sharpness:
		crs.Sharpness = parseInt16(p.Value())
	case AlreadyApplied:
		crs.AlreadyApplied = parseBool(p.Value())
	// Accepted for compatibility but intentionally not materialized in CRS.
	// Return ErrPropertyNotSet so parseNamespace doesn't allocate CRS for ignored-only payloads.
	case HueAdjustmentRed, HueAdjustmentOrange, HueAdjustmentYellow, HueAdjustmentGreen, HueAdjustmentAqua, HueAdjustmentBlue, HueAdjustmentPurple, HueAdjustmentMagenta:
		return ErrPropertyNotSet
	case SaturationAdjustmentRed, SaturationAdjustmentOrange, SaturationAdjustmentYellow, SaturationAdjustmentGreen, SaturationAdjustmentAqua, SaturationAdjustmentBlue, SaturationAdjustmentPurple, SaturationAdjustmentMagenta:
		return ErrPropertyNotSet
	case LuminanceAdjustmentRed, LuminanceAdjustmentOrange, LuminanceAdjustmentYellow, LuminanceAdjustmentGreen, LuminanceAdjustmentAqua, LuminanceAdjustmentBlue, LuminanceAdjustmentPurple, LuminanceAdjustmentMagenta:
		return ErrPropertyNotSet
	// Accepted for compatibility but intentionally not materialized in CRS.
	case ToneCurve, ToneCurveRed, ToneCurveGreen, ToneCurveBlue:
		return ErrPropertyNotSet
	case ToneCurvePV2012, ToneCurvePV2012Red, ToneCurvePV2012Green, ToneCurvePV2012Blue:
		return ErrPropertyNotSet
	default:
		return ErrPropertyNotSet
	}
	return nil
}

// DynamicMedia stores tags from xmlns:xmpDM="http://ns.adobe.com/xmp/1.0/DynamicMedia/".
type DynamicMedia struct {
	// Pick is Lightroom's pick flag (0: none, 1: picked, -1: rejected).
	Pick int8
	// Good is Adobe's binary "good shot" marker.
	Good bool
}

func (dm *DynamicMedia) parse(prop property) error {
	switch prop.Name() {
	case Pick:
		dm.Pick = int8(parseInt16(prop.Value()))
	case Good:
		dm.Good = parseBool(prop.Value())
	default:
		return ErrPropertyNotSet
	}
	return nil
}

// Lightroom stores tags from xmlns:lr="http://ns.adobe.com/lightroom/1.0/".
type Lightroom struct {
	// HierarchicalSubject stores catalog subjects with hierarchy semantics.
	HierarchicalSubject []string
	// WeightedFlatSubject stores Lightroom's weighted flat subject list.
	WeightedFlatSubject []string
}

func (lr *Lightroom) parse(prop property) error {
	switch prop.Name() {
	case HierarchicalSubject:
		lr.HierarchicalSubject = append(lr.HierarchicalSubject, parseString(prop.Value()))
	case WeightedFlatSubject:
		lr.WeightedFlatSubject = append(lr.WeightedFlatSubject, parseString(prop.Value()))
	default:
		return ErrPropertyNotSet
	}
	return nil
}

// Photoshop stores tags from xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/".
type Photoshop struct {
	// DateCreated is the resource creation timestamp stored in Photoshop namespace.
	DateCreated time.Time
	// SidecarForExtension identifies the related raw extension (for example, "CR2").
	SidecarForExtension string
	// EmbeddedXMPDigest is Adobe's digest for embedded XMP reconciliation.
	EmbeddedXMPDigest string
	// LegacyIPTCDigest is Adobe's digest for IPTC legacy synchronization.
	LegacyIPTCDigest string
	// ColorMode is Photoshop's color-mode code.
	ColorMode uint16
	// ICCProfile is the embedded profile name.
	ICCProfile string
	// History stores Photoshop history text when present.
	History string
}

func (p *Photoshop) parse(prop property) (err error) {
	switch prop.Name() {
	case DateCreated:
		p.DateCreated, err = parseDate(prop.Value())
	case SidecarForExtension:
		p.SidecarForExtension = parseString(prop.Value())
	case EmbeddedXMPDigest:
		p.EmbeddedXMPDigest = parseString(prop.Value())
	case LegacyIPTCDigest:
		p.LegacyIPTCDigest = parseString(prop.Value())
	case ColorMode:
		p.ColorMode = parseUint16(prop.Value())
	case ICCProfile:
		p.ICCProfile = parseString(prop.Value())
	case HistoryTag:
		p.History = parseString(prop.Value())
	default:
		return ErrPropertyNotSet
	}
	return err
}
