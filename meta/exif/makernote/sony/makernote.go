package sony

// Sony contains the selected Sony maker-note fields currently decoded by
// imagemeta.
//
// The field set mirrors the subset of ExifTool's Image::ExifTool::Sony::Main
// and Sony::CameraSettings3 tables that imagemeta parses today.
type Sony struct {
	Rating                uint32
	Contrast              int32
	Saturation            int32
	Sharpness             int32
	CreativeStyle         string
	DynamicRangeOptimizer uint32
	ImageStabilization    uint32
	ColorMode             uint32
	Quality               uint32
	Quality2              [2]uint16
	WhiteBalance          uint32
	WhiteBalanceFineTune  int32
	FlashExposureComp     float64
	Teleconverter         uint32
	SonyModelID           uint16
	LensType              uint32
	CameraSettings3       CameraSettings3
}

// CameraSettings3 stores selected values from ExifTool's
// Image::ExifTool::Sony::CameraSettings3 table.
type CameraSettings3 struct {
	FocalLength         float64
	FocalLengthTeleZoom float64
	AFPointSelected     uint8
	FocusMode           uint8
	AFPoint             uint8
	FocusStatus         uint8
}
