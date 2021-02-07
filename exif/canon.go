package exif

import (
	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/meta/canon"
)

// CanonCameraSettings convenience func. "IFD/Exif/Makernotes.Canon" CanonCameraSettings
// Canon Camera Settings from the Makernote
func (e *ExifData) CanonCameraSettings() (canon.CameraSettings, error) {
	if e.make != "Canon" {
		return canon.CameraSettings{}, ErrEmptyTag
	}
	t, err := e.GetTag(ifds.MknoteIFD, 0, mknote.CanonCameraSettings)
	if err != nil {
		return canon.CameraSettings{}, err
	}
	ii, err := t.Uint16Values(e.exifReader)
	if len(ii) < 24 || err != nil {
		return canon.CameraSettings{}, err
	}
	return canon.CameraSettings{
		Macromode:         intToBool(ii[1]),
		SelfTimer:         intToBool(ii[2]),
		ContinuousDrive:   canon.ContinuousDrive(ii[5]),
		FocusMode:         canon.FocusMode(ii[7]),
		MeteringMode:      canon.MeteringMode(ii[17]),
		FocusRange:        canon.FocusRange(ii[18]),
		CanonExposureMode: canon.ExposureMode(ii[20]),
		AESetting:         canon.AESetting(ii[33]),
	}, nil
}

// CanonFileInfo convenience func. "IFD/Exif/Makernotes.Canon" CanonFileInfo
// Canon Camera File Info from the Makernote
func (e *ExifData) CanonFileInfo() (canon.FileInfo, error) {
	t, err := e.GetTag(ifds.MknoteIFD, 0, mknote.CanonFileInfo)
	if err != nil {
		return canon.FileInfo{}, err
	}
	ii, err := t.Uint16Values(e.exifReader)
	if len(ii) < 21 || err != nil {
		return canon.FileInfo{}, err
	}
	return canon.FileInfo{
		FocusDistance:     canon.NewFocusDistance(ii[20], ii[21]),
		BracketMode:       canon.BracketMode(ii[3]),
		BracketValue:      canon.Ev(int16(ii[4])),
		BracketShotNumber: int16(ii[5]),
		LiveViewShooting:  intToBool(ii[19]),
	}, nil
}

// CanonShotInfo convenience func. "IFD/Exif/Makernotes.Canon" CanonShotInfo
// Canon Camera Shot Info from the Makernote
func (e *ExifData) CanonShotInfo() (canon.ShotInfo, error) {
	t, err := e.GetTag(ifds.MknoteIFD, 0, mknote.CanonShotInfo)
	if err != nil {
		return canon.ShotInfo{}, err
	}
	si, err := t.Uint16Values(e.exifReader)
	if len(si) < 29 || err != nil {
		return canon.ShotInfo{}, err
	}

	return canon.ShotInfo{
		CameraTemperature:      canon.TempConv(si[12]),
		FlashExposureComp:      int16(si[15]),
		AutoExposureBracketing: int16(si[16]),
		AEBBracketValue:        canon.Ev(int16(si[17])),
		SelfTimer:              int16(si[29]) / 10,
		FocusDistance:          canon.NewFocusDistance(si[19], si[20]),
	}, nil
}

// CanonAFInfo -
// Canon Camera AutoFocus Information from the Makernote
func (e *ExifData) CanonAFInfo() (afInfo canon.AFInfo, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	t, err := e.GetTag(ifds.MknoteIFD, 0, mknote.CanonAFInfo2)
	if err != nil {
		return canon.AFInfo{}, err
	}
	af, err := t.Uint16Values(e.exifReader)
	if len(af) < 8 || err != nil {
		panic(ErrEmptyTag)
	}

	afInfo = canon.AFInfo{
		AFAreaMode:    canon.AFAreaMode(af[1]),
		NumAFPoints:   af[2],
		ValidAFPoints: af[3],
		AFPoints:      canon.ParseAFPoints(af),
	}

	if infocus, selected, err := canon.PointsInFocus(af); err == nil {
		afInfo.InFocus = infocus
		afInfo.Selected = selected
	} else {
		panic(err)
	}

	return afInfo, nil
}

func intToBool(i uint16) bool {
	return i == 1
}
