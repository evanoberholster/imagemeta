package makernote

import (
	"math/bits"

	metacanon "github.com/evanoberholster/imagemeta/meta/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

const (
	canonCameraSettingsFirstEntry = 1
	canonShotInfoFirstEntry       = 1
	canonFileInfoFirstEntry       = 1
)

// CanonValueDecoder provides value decoding primitives for Canon maker-note tags.
//
// Reader implementations in the exif package provide concrete decoders so
// Canon parser logic can live in this package.
type CanonValueDecoder interface {
	String(tag.Entry) string
	Uint32(tag.Entry) uint32
	Int16List(tag.Entry, []int16) int
	Uint16List(tag.Entry, []uint16) int
}

// Canon contains selected Canon maker-note fields.
type Canon struct {
	ImageType            string
	FirmwareVersion      string
	OwnerName            string
	LensModel            string
	InternalSerialNumber string

	FileNumber   uint32
	SerialNumber uint32

	// Structured Canon maker-note tables (ExifTool Canon.pm mappings).
	CameraSettings metacanon.CameraSettings
	ShotInfo       metacanon.ShotInfo
	FileInfo       metacanon.FileInfo
	AFInfo         metacanon.AFInfo
}

// ParseCanonTag parses a Canon maker-note tag into the destination model.
func ParseCanonTag(dst *Canon, t tag.Entry, dec CanonValueDecoder) bool {
	switch metacanon.MakerNoteTag(t.ID) {
	case metacanon.CanonImageType:
		dst.ImageType = dec.String(t)
	case metacanon.CanonFirmwareVersion:
		dst.FirmwareVersion = dec.String(t)
	case metacanon.FileNumber:
		dst.FileNumber = dec.Uint32(t)
	case metacanon.OwnerName:
		dst.OwnerName = dec.String(t)
	case metacanon.SerialNumber:
		dst.SerialNumber = dec.Uint32(t)
	case metacanon.LensModel:
		dst.LensModel = dec.String(t)
	case metacanon.CanonInternalSerialNumber:
		dst.InternalSerialNumber = dec.String(t)
	case metacanon.CanonCameraSettings:
		parseCanonCameraSettings(&dst.CameraSettings, t, dec)
	case metacanon.CanonShotInfo:
		parseCanonShotInfo(&dst.ShotInfo, t, dec)
	case metacanon.CanonFileInfo:
		parseCanonFileInfo(&dst.FileInfo, t, dec)
	case metacanon.CanonAFInfo:
		parseCanonAFInfo(&dst.AFInfo, t, dec)
	case metacanon.CanonAFInfo2, metacanon.AFInfo3:
		parseCanonAFInfo2(&dst.AFInfo, t, dec)
	default:
		return false
	}
	return true
}

// parseCanonCameraSettings parses tag 0x0001 (CanonCameraSettings).
func parseCanonCameraSettings(dst *metacanon.CameraSettings, t tag.Entry, dec CanonValueDecoder) {
	var rawSigned [64]int16
	var rawUnsigned [64]uint16
	n := parseCanonInt16Table(t, dec, rawSigned[:], rawUnsigned[:])
	if n == 0 {
		return
	}

	s := rawSigned[:n]
	u := rawUnsigned[:n]
	atS := func(tableIndex int) int16 {
		idx := tableIndex - canonCameraSettingsFirstEntry
		if uint(idx) >= uint(len(s)) {
			return 0
		}
		return s[idx]
	}
	atU := func(tableIndex int) uint16 {
		idx := tableIndex - canonCameraSettingsFirstEntry
		if uint(idx) >= uint(len(u)) {
			return 0
		}
		return u[idx]
	}

	dst.MacroMode = atS(1)
	dst.SelfTimer = atS(2)
	dst.Quality = atS(3)
	dst.CanonFlashMode = atS(4)
	dst.ContinuousDrive = metacanon.ContinuousDrive(atS(5))
	dst.FocusMode = metacanon.FocusMode(atS(7))
	dst.RecordMode = atS(9)
	dst.CanonImageSize = atS(10)
	dst.EasyMode = atS(11)
	dst.DigitalZoom = atS(12)
	dst.Contrast = atS(13)
	dst.Saturation = atS(14)
	dst.Sharpness = atS(15)
	dst.CameraISO = atS(16)
	dst.MeteringMode = metacanon.MeteringMode(atS(17))
	dst.FocusRange = metacanon.FocusRange(atS(18))
	dst.AFPoint = atU(19)
	dst.CanonExposureMode = metacanon.ExposureMode(atS(20))
	dst.LensType = atU(22)
	dst.MaxFocalLength = atU(23)
	dst.MinFocalLength = atU(24)
	dst.FocalUnits = atU(25)
	dst.MaxAperture = atS(26)
	dst.MinAperture = atS(27)
	dst.FlashActivity = atS(28)
	dst.FlashBits = atU(29)
	dst.FocusContinuous = atS(32)
	dst.AESetting = metacanon.AESetting(atS(33))
	dst.ImageStabilization = atS(34)
	dst.DisplayAperture = atU(35)
	dst.ZoomSourceWidth = atU(36)
	dst.ZoomTargetWidth = atU(37)
	dst.SpotMeteringMode = atS(39)
	dst.PhotoEffect = atS(40)
	dst.ManualFlashOutput = atS(41)
	dst.ColorTone = atS(42)
	dst.SRAWQuality = atS(46)
	dst.Clarity = atS(51)
}

// parseCanonShotInfo parses tag 0x0004 (CanonShotInfo).
func parseCanonShotInfo(dst *metacanon.ShotInfo, t tag.Entry, dec CanonValueDecoder) {
	var rawSigned [96]int16
	var rawUnsigned [96]uint16
	n := parseCanonInt16Table(t, dec, rawSigned[:], rawUnsigned[:])
	if n == 0 {
		return
	}

	const base = canonShotInfoFirstEntry
	dst.AutoISO = i16At(rawSigned[:], n, 1-base)
	dst.BaseISO = i16At(rawSigned[:], n, 2-base)
	dst.MeasuredEV = i16At(rawSigned[:], n, 3-base)
	dst.TargetAperture = i16At(rawSigned[:], n, 4-base)
	dst.TargetExposureTime = i16At(rawSigned[:], n, 5-base)
	dst.ExposureCompensation = i16At(rawSigned[:], n, 6-base)
	dst.WhiteBalance = i16At(rawSigned[:], n, 7-base)
	dst.SlowShutter = i16At(rawSigned[:], n, 8-base)
	dst.SequenceNumber = i16At(rawSigned[:], n, 9-base)
	dst.OpticalZoomCode = i16At(rawSigned[:], n, 10-base)
	dst.CameraTemperature = i16At(rawSigned[:], n, 12-base)
	dst.FlashGuideNumber = i16At(rawSigned[:], n, 13-base)
	dst.AFPointsInFocus = u16At(rawUnsigned[:], n, 14-base)
	dst.FlashExposureComp = i16At(rawSigned[:], n, 15-base)
	dst.AutoExposureBracketing = i16At(rawSigned[:], n, 16-base)
	dst.AEBBracketValue = i16At(rawSigned[:], n, 17-base)
	dst.ControlMode = i16At(rawSigned[:], n, 18-base)
	dst.FocusDistance = metacanon.NewFocusDistance(
		u16At(rawUnsigned[:], n, 19-base),
		u16At(rawUnsigned[:], n, 20-base),
	)
	dst.FNumber = i16At(rawSigned[:], n, 21-base)
	dst.ExposureTime = i16At(rawSigned[:], n, 22-base)
	dst.MeasuredEV2 = i16At(rawSigned[:], n, 23-base)
	dst.BulbDuration = i16At(rawSigned[:], n, 24-base)
	dst.CameraType = i16At(rawSigned[:], n, 26-base)
	dst.AutoRotate = i16At(rawSigned[:], n, 27-base)
	dst.NDFilter = i16At(rawSigned[:], n, 28-base)
	dst.SelfTimer2 = i16At(rawSigned[:], n, 29-base)
	dst.FlashOutput = i16At(rawSigned[:], n, 33-base)
}

// parseCanonFileInfo parses tag 0x0093 (CanonFileInfo).
func parseCanonFileInfo(dst *metacanon.FileInfo, t tag.Entry, dec CanonValueDecoder) {
	var rawSigned [128]int16
	var rawUnsigned [128]uint16
	n := parseCanonInt16Table(t, dec, rawSigned[:], rawUnsigned[:])
	if n == 0 {
		return
	}

	s := rawSigned[:n]
	u := rawUnsigned[:n]
	atS := func(tableIndex int) int16 {
		idx := tableIndex - canonFileInfoFirstEntry
		if uint(idx) >= uint(len(s)) {
			return 0
		}
		return s[idx]
	}
	atU := func(tableIndex int) uint16 {
		idx := tableIndex - canonFileInfoFirstEntry
		if uint(idx) >= uint(len(u)) {
			return 0
		}
		return u[idx]
	}

	if n > 1 {
		// Tag 0x0093 index 1 is model-dependent (FileNumber or ShutterCount).
		// Preserve raw 32-bit representation for both fields.
		raw32 := uint32(u[0]) | (uint32(u[1]) << 16)
		dst.FileNumber = raw32
		dst.ShutterCount = raw32
	} else {
		raw32 := uint32(u[0])
		dst.FileNumber = raw32
		dst.ShutterCount = raw32
	}
	dst.BracketMode = metacanon.BracketMode(atS(3))
	dst.BracketValue = atS(4)
	dst.BracketShotNumber = atS(5)
	dst.RawJpgQuality = atS(6)
	dst.RawJpgSize = atS(7)
	dst.LongExposureNoiseReduction2 = atS(8)
	dst.WBBracketMode = atS(9)
	dst.WBBracketValueAB = atS(12)
	dst.WBBracketValueGM = atS(13)
	dst.FilterEffect = atS(14)
	dst.ToningEffect = atS(15)
	dst.MacroMagnification = atS(16)
	dst.LiveViewShooting = atS(19) != 0
	dst.FocusDistance = metacanon.NewFocusDistance(
		atU(20),
		atU(21),
	)
	dst.ShutterMode = atS(23)
	dst.FlashExposureLock = atS(25) != 0
	dst.AntiFlicker = atS(32) != 0
	dst.RFLensType = atU(0x3d)
}

// parseCanonAFInfo parses tag 0x0012 (AFInfo).
func parseCanonAFInfo(dst *metacanon.AFInfo, t tag.Entry, dec CanonValueDecoder) {
	var words [2048]uint16
	n := dec.Uint16List(t, words[:])
	if n == 0 {
		return
	}

	dst.NumAFPoints = u16At(words[:], n, 0)
	dst.ValidAFPoints = u16At(words[:], n, 1)
	dst.CanonImageWidth = u16At(words[:], n, 2)
	dst.CanonImageHeight = u16At(words[:], n, 3)
	dst.AFImageWidth = u16At(words[:], n, 4)
	dst.AFImageHeight = u16At(words[:], n, 5)
	dst.AFAreaWidth = u16At(words[:], n, 6)
	dst.AFAreaHeight = u16At(words[:], n, 7)

	num := int(dst.NumAFPoints)
	if num <= 0 {
		dst.AFAreaXPositions = nil
		dst.AFAreaYPositions = nil
		dst.AFPointsInFocusBits = nil
		dst.AFPoints = nil
		dst.InFocus = nil
		dst.Selected = nil
		return
	}

	xStart := 8
	yStart := xStart + num
	xVals := signedRangeFromUint16(words[:], n, xStart, num)
	yVals := signedRangeFromUint16(words[:], n, yStart, num)
	dst.AFAreaXPositions = xVals
	dst.AFAreaYPositions = yVals

	bitWords := bitWordCount(num)
	inFocusStart := yStart + num
	inFocusWords := uint16Range(words[:], n, inFocusStart, bitWords)
	dst.AFPointsInFocusBits = decodeBitWords(inFocusWords, num)
	dst.AFPointsSelectedBits = nil
	dst.InFocus = dst.AFPointsInFocusBits
	dst.Selected = nil

	primaryIndex := inFocusStart + bitWords
	dst.PrimaryAFPoint = u16At(words[:], n, primaryIndex)

	pointCount := num
	if len(xVals) < pointCount {
		pointCount = len(xVals)
	}
	if len(yVals) < pointCount {
		pointCount = len(yVals)
	}
	if pointCount == 0 {
		dst.AFPoints = nil
		return
	}

	pts := make([]metacanon.AFPoint, pointCount)
	w := int16(dst.AFAreaWidth)
	h := int16(dst.AFAreaHeight)
	for i := 0; i < pointCount; i++ {
		pts[i] = metacanon.NewAFPoint(w, h, xVals[i], yVals[i])
	}
	dst.AFPoints = pts
}

// parseCanonAFInfo2 parses tags 0x0026 and 0x003c (AFInfo2/AFInfo3).
func parseCanonAFInfo2(dst *metacanon.AFInfo, t tag.Entry, dec CanonValueDecoder) {
	var words [2048]uint16
	n := dec.Uint16List(t, words[:])
	if n == 0 {
		return
	}

	dst.AFAreaMode = metacanon.AFAreaMode(u16At(words[:], n, 1))
	dst.NumAFPoints = u16At(words[:], n, 2)
	dst.ValidAFPoints = u16At(words[:], n, 3)
	dst.CanonImageWidth = u16At(words[:], n, 4)
	dst.CanonImageHeight = u16At(words[:], n, 5)
	dst.AFImageWidth = u16At(words[:], n, 6)
	dst.AFImageHeight = u16At(words[:], n, 7)

	num := int(dst.NumAFPoints)
	if num <= 0 {
		dst.AFAreaWidths = nil
		dst.AFAreaHeights = nil
		dst.AFAreaXPositions = nil
		dst.AFAreaYPositions = nil
		dst.AFPointsInFocusBits = nil
		dst.AFPointsSelectedBits = nil
		dst.AFPoints = nil
		dst.InFocus = nil
		dst.Selected = nil
		return
	}

	widthStart := 8
	heightStart := widthStart + num
	xStart := heightStart + num
	yStart := xStart + num
	bitsStart := yStart + num
	maskWordCount := bitWordCount(num)
	selectedStart := bitsStart + maskWordCount

	dst.AFAreaWidths = signedRangeFromUint16(words[:], n, widthStart, num)
	dst.AFAreaHeights = signedRangeFromUint16(words[:], n, heightStart, num)
	dst.AFAreaXPositions = signedRangeFromUint16(words[:], n, xStart, num)
	dst.AFAreaYPositions = signedRangeFromUint16(words[:], n, yStart, num)

	inFocusWords := uint16Range(words[:], n, bitsStart, maskWordCount)
	selectedWords := uint16Range(words[:], n, selectedStart, maskWordCount)
	dst.AFPointsInFocusBits = decodeBitWords(inFocusWords, num)
	dst.AFPointsSelectedBits = decodeBitWords(selectedWords, num)
	dst.InFocus = dst.AFPointsInFocusBits
	dst.Selected = dst.AFPointsSelectedBits
	dst.PrimaryAFPoint = u16At(words[:], n, selectedStart+maskWordCount)

	pointCount := num
	if len(dst.AFAreaWidths) < pointCount {
		pointCount = len(dst.AFAreaWidths)
	}
	if len(dst.AFAreaHeights) < pointCount {
		pointCount = len(dst.AFAreaHeights)
	}
	if len(dst.AFAreaXPositions) < pointCount {
		pointCount = len(dst.AFAreaXPositions)
	}
	if len(dst.AFAreaYPositions) < pointCount {
		pointCount = len(dst.AFAreaYPositions)
	}
	if pointCount == 0 {
		dst.AFPoints = nil
		return
	}

	pts := make([]metacanon.AFPoint, pointCount)
	xAdjust := int16(dst.CanonImageWidth / 2)
	yAdjust := int16(dst.CanonImageHeight / 2)
	for i := 0; i < pointCount; i++ {
		w := dst.AFAreaWidths[i]
		h := dst.AFAreaHeights[i]
		x := dst.AFAreaXPositions[i] + xAdjust - (w / 2)
		y := dst.AFAreaYPositions[i] + yAdjust - (h / 2)
		pts[i] = metacanon.NewAFPoint(w, h, x, y)
	}
	dst.AFPoints = pts
}

func parseCanonInt16Table(t tag.Entry, dec CanonValueDecoder, signed []int16, unsigned []uint16) int {
	if n := dec.Int16List(t, signed); n > 0 {
		for i := 0; i < n; i++ {
			unsigned[i] = uint16(signed[i])
		}
		return n
	}
	n := dec.Uint16List(t, unsigned)
	for i := 0; i < n; i++ {
		signed[i] = int16(unsigned[i])
	}
	return n
}

func i16At(vals []int16, n, idx int) int16 {
	if idx < 0 || idx >= n {
		return 0
	}
	return vals[idx]
}

func u16At(vals []uint16, n, idx int) uint16 {
	if idx < 0 || idx >= n {
		return 0
	}
	return vals[idx]
}

func bitWordCount(pointCount int) int {
	if pointCount <= 0 {
		return 0
	}
	return (pointCount + 15) / 16
}

func signedRangeFromUint16(vals []uint16, n, start, count int) []int16 {
	if count <= 0 || start < 0 || start >= n {
		return nil
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return nil
	}
	out := make([]int16, end-start)
	for i := 0; i < len(out); i++ {
		out[i] = int16(vals[start+i])
	}
	return out
}

func uint16Range(vals []uint16, n, start, count int) []uint16 {
	if count <= 0 || start < 0 || start >= n {
		return nil
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return nil
	}
	out := make([]uint16, end-start)
	copy(out, vals[start:end])
	return out
}

func decodeBitWords(words []uint16, limit int) []int {
	if len(words) == 0 || limit <= 0 {
		return nil
	}
	capHint := 0
	for _, w := range words {
		capHint += bits.OnesCount16(w)
	}
	if capHint > limit {
		capHint = limit
	}
	out := make([]int, 0, capHint)
	base := 0
	for _, word := range words {
		for bit := 0; bit < 16; bit++ {
			idx := base + bit
			if idx >= limit {
				return out
			}
			if word&(1<<bit) != 0 {
				out = append(out, idx)
			}
		}
		base += 16
	}
	return out
}
