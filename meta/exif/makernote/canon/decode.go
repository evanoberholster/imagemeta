package canon

// DecodeFileInfo decodes a Canon FileInfo payload (tag 0x0093).
// The words slice must be the FIRST_ENTRY-stripped data (size word removed).
func DecodeFileInfo(s Seq16, modelID CanonCameraModel) FileInfo {
	dst := FileInfo{
		FileNumber:                  uint32(s.U16(1)) | (uint32(s.U16(2)) << 16),
		BracketMode:                 BracketMode(s.I16(3)),
		BracketValue:                s.I16(4),
		BracketShotNumber:           s.I16(5),
		RawJpgQuality:               RawJpgQuality(s.U16(6)),
		RawJpgSize:                  RawJpgSize(s.U16(7)),
		LongExposureNoiseReduction2: OnOffAuto(s.U16(8)),
		WBBracketMode:               s.I16(9),
		WBBracketValueAB:            s.I16(12),
		WBBracketValueGM:            s.I16(13),
		FilterEffect:                FilterEffect(s.U16(14)),
		ToningEffect:                ToningEffect(s.U16(15)),
		MacroMagnification:          s.I16(16),
		LiveViewShooting:            OnOffAuto(s.U16(19)),
		FocusDistance:               NewFocusDistance(s.U16(20), s.U16(21)),
		ShutterMode:                 ShutterMode(s.U16(23)),
		FlashExposureLock:           OnOffAuto(s.U16(25)),
		AntiFlicker:                 OnOffAuto(s.U16(32)),
		RFLensType:                  CanonRFLensType(s.U16(61)),
	}
	if ModelUsesLegacyShutterCount(modelID) {
		dst.ShutterCount = uint32(s.U16(2))
		dst.FileNumber = 0
	}
	return dst
}

// DecodeCameraSettings decodes a Canon CameraSettings payload (tag 0x0001).
func DecodeCameraSettings(s Seq16) CameraSettings {
	return CameraSettings{
		MacroMode:          MacroMode(s.U16(1)),
		SelfTimer:          s.I16(2),
		Quality:            Quality(s.I16(3)),
		CanonFlashMode:     CanonFlashMode(s.I16(4)),
		ContinuousDrive:    ContinuousDrive(s.I16(5)),
		FocusMode:          FocusMode(s.I16(7)),
		RecordMode:         RecordMode(s.I16(9)),
		CanonImageSize:     CanonImageSize(s.I16(10)),
		EasyMode:           EasyMode(s.I16(11)),
		DigitalZoom:        DigitalZoom(s.I16(12)),
		Contrast:           CameraSettingValue(s.I16(13)),
		Saturation:         CameraSettingValue(s.I16(14)),
		Sharpness:          CameraSettingValue(s.I16(15)),
		CameraISO:          uint32(CameraISOValue(s.I16(16))),
		MeteringMode:       MeteringMode(s.I16(17)),
		FocusRange:         FocusRange(s.I16(18)),
		AFPoint:            s.U16(19),
		CanonExposureMode:  ExposureMode(s.I16(20)),
		LensType:           CanonLensType(s.U16(22)),
		MaxFocalLength:     s.U16(23),
		MinFocalLength:     s.U16(24),
		FocalUnits:         s.U16(25),
		MaxAperture:        MaxApertureFromCode(s.U16(26)),
		MinAperture:        MaxApertureFromCode(s.U16(27)),
		FlashModel:         FlashModel(s.I16(28)),
		FlashBits:          s.U16(29),
		FocusContinuous:    FocusContinuous(s.I16(32)),
		AESetting:          AESetting(s.I16(33)),
		ImageStabilization: ImageStabilization(s.I16(34)),
		DisplayAperture:    DisplayApertureFromCode(s.U16(35)),
		ZoomSourceWidth:    s.U16(36),
		ZoomTargetWidth:    s.U16(37),
		SpotMeteringMode:   SpotMeteringMode(s.I16(39)),
		PhotoEffect:        PhotoEffect(s.I16(40)),
		ManualFlashOutput:  ManualFlashOutput(s.I16(41)),
		ColorTone:          CameraSettingValue(s.I16(42)),
		SRAWQuality:        SRAWQuality(s.I16(46)),
		FocusBracketing:    FocusBracketing(s.I16(50)),
		Clarity:            CameraSettingValue(s.I16(51)),
		HDRPQ:              HDRPQ(s.U16(52)),
	}
}

// DecodeAFConfig decodes a Canon AFConfig payload (tag 0x4028).
func DecodeAFConfig(words []int32) AFConfig {
	s := Seq32(words)
	return AFConfig{
		AFConfigTool:              s.U32(1) + 1,
		AFTrackingSensitivity:     s.I32(2),
		AFAccelDecelTracking:      s.I32(3),
		AFPointSwitching:          s.I32(4),
		AIServoFirstImage:         s.I32(5),
		AIServoSecondImage:        s.I32(6),
		USMLensElectronicMF:       s.I32(7),
		AFAssistBeam:              s.I32(8),
		OneShotAFRelease:          s.I32(9),
		AutoAFPointSelEOSiTRAF:    s.I32(10),
		LensDriveWhenAFImpossible: s.I32(11),
		SelectAFAreaSelectionMode: s.U32(12),
		AFAreaSelectionMethod:     s.I32(13),
		OrientationLinkedAF:       s.I32(14),
		ManualAFPointSelPattern:   s.I32(15),
		AFPointDisplayDuringFocus: s.I32(16),
		VFDisplayIllumination:     s.I32(17),
		AFStatusViewfinder:        s.I32(18),
		InitialAFPointInServo:     s.I32(19),
		SubjectToDetect:           s.I32(20),
		EyeDetection:              s.I32(24),
	}
}

// DecodeTimeInfo decodes a Canon TimeInfo payload (tag 0x0035).
func DecodeTimeInfo(words []int32) CanonTimeInfo {
	s := Seq32(words)
	return CanonTimeInfo{
		TimeZone:        s.I32(1),
		TimeZoneCity:    TimeZoneCity(s.I32(2)),
		DaylightSavings: DaylightSavings(s.I32(3)),
	}
}

// DecodeLightingOpt decodes a Canon LightingOpt payload (tag 0x4018).
func DecodeLightingOpt(words []int32) LightingOptInfo {
	s := Seq32(words)
	return LightingOptInfo{
		PeripheralIlluminationCorr: s.I32(1),
		AutoLightingOptimizer:      s.I32(2),
		HighlightTonePriority:      s.I32(3),
		LongExposureNoiseReduction: s.I32(4),
		HighISONoiseReduction:      s.I32(5),
		DigitalLensOptimizer:       s.I32(10),
		DualPixelRaw:               s.I32(11),
	}
}

// DecodeMultiExp decodes a Canon MultiExp payload (tag 0x4021).
func DecodeMultiExp(words []int32) MultiExpInfo {
	s := Seq32(words)
	return MultiExpInfo{
		MultiExposure:        s.I32(1),
		MultiExposureControl: s.I32(2),
		MultiExposureShots:   s.I32(3),
	}
}

// DecodeHDRInfo decodes a Canon HDRInfo payload (tag 0x4025).
func DecodeHDRInfo(words []int32) HDRInfo {
	s := Seq32(words)
	return HDRInfo{
		HDR:       s.I32(1),
		HDREffect: s.I32(2),
	}
}

// DecodeAFMicroAdj decodes a Canon AFMicroAdj payload (tag 0x4013).
func DecodeAFMicroAdj(words []int32) AFMicroAdjInfo {
	s := Seq32(words)
	return AFMicroAdjInfo{
		Mode:             s.I32(1),
		ValueNumerator:   s.I32(2),
		ValueDenominator: s.I32(3),
	}
}

// DecodeSensorInfo decodes a Canon SensorInfo payload (tag 0x00e0).
func DecodeSensorInfo(s Seq16) SensorInfo {
	if len(s) < 13 {
		return SensorInfo{}
	}
	return SensorInfo{
		SensorWidth:           s.I16(1),
		SensorHeight:          s.I16(2),
		SensorLeftBorder:      s.I16(5),
		SensorTopBorder:       s.I16(6),
		SensorRightBorder:     s.I16(7),
		SensorBottomBorder:    s.I16(8),
		BlackMaskLeftBorder:   s.I16(9),
		BlackMaskTopBorder:    s.I16(10),
		BlackMaskRightBorder:  s.I16(11),
		BlackMaskBottomBorder: s.I16(12),
	}
}

// DecodeFocalLength decodes a Canon FocalLength payload (tag 0x0002).
func DecodeFocalLength(s Seq16) FocalLengthInfo {
	if len(s) < 4 {
		return FocalLengthInfo{}
	}
	return FocalLengthInfo{
		FocalType:       s[0],
		FocalLength:     s[1],
		FocalPlaneXSize: FocalPlaneSizeMM(s[2]),
		FocalPlaneYSize: FocalPlaneSizeMM(s[3]),
	}
}

// DecodeProcessingInfo decodes a Canon ProcessingInfo payload (tag 0x00a0).
// The payload uses FIRST_ENTRY=1 convention.
func DecodeProcessingInfo(s Seq16) ProcessingInfo {
	return ProcessingInfo{
		ToneCurve:            s.I16(1),
		Sharpness:            s.I16(2),
		SharpnessFrequency:   s.I16(3),
		SensorRedLevel:       s.I16(4),
		SensorBlueLevel:      s.I16(5),
		WhiteBalanceRed:      s.I16(6),
		WhiteBalanceBlue:     s.I16(7),
		WhiteBalance:         s.I16(8),
		ColorTemperature:     s.I16(9),
		PictureStyle:         s.I16(10),
		DigitalGain:          s.I16(11),
		WBShiftAB:            s.I16(12),
		WBShiftGM:            s.I16(13),
		UnsharpMaskFineness:  s.I16(14),
		UnsharpMaskThreshold: s.I16(15),
	}
}

// ShotInfoDecodeConfig holds options for decoding Canon ShotInfo data.
type ShotInfoDecodeConfig struct {
	LegacyExposureTime bool
	ModelID            CanonCameraModel
}

// DecodeShotInfo decodes a Canon ShotInfo payload (tag 0x0004).
func (s Seq16) DecodeShotInfo(cfg ShotInfoDecodeConfig) ShotInfo {
	dst := ShotInfo{
		AutoISO:                ShotISO(s.I16(1)),
		BaseISO:                ShotISO(s.I16(2)),
		MeasuredEV:             ShotMeasuredEV(s.I16(3)),
		TargetAperture:         ShotAperture(s.I16(4)),
		TargetExposureTime:     ShotExposureTime(s.I16(5), false),
		ExposureCompensation:   ShotExposureCompensation(s.I16(6)),
		WhiteBalance:           WhiteBalance(s.I16(7)),
		SlowShutter:            SlowShutter(s.I16(8)),
		SequenceNumber:         s.I16(9),
		OpticalZoomCode:        s.I16(10),
		CameraTemperature:      ShotCameraTemperature(s.I16(12), cfg.ModelID),
		FlashGuideNumber:       ShotFlashGuideNumber(s.I16(13)),
		AFPointsInFocus:        s.U16(14),
		FlashExposureComp:      s.I16(15),
		AutoExposureBracketing: s.I16(16),
		AEBBracketValue:        s.I16(17),
		ControlMode:            s.I16(18),
		FNumber:                ShotAperture(s.I16(21)),
		ExposureTime:           ShotExposureTime(s.I16(22), cfg.LegacyExposureTime),
		MeasuredEV2:            ShotMeasuredEV2(s.I16(23)),
		BulbDuration:           s.I16(24),
		CameraType:             CameraType(s.I16(26)),
		AutoRotate:             AutoRotate(s.I16(27)),
		NDFilter:               NDFilter(s.I16(28)),
		SelfTimer2:             s.I16(29),
		FlashOutput:            s.I16(33),
	}
	dst.ActualISO = ShotActualISO(dst.AutoISO, dst.BaseISO)
	upper := s.U16(19)
	lower := s.U16(20)
	if s.Present(20) && upper != 0 {
		dst.FocusDistance = NewFocusDistance(upper, lower)
	}
	return dst
}
