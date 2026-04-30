package canon

// CameraInfoDecode parses a Canon CameraInfo byte payload using model-specific
// byte offsets. The spec should be obtained from CameraInfoSpecLayout* variables.
func CameraInfoDecode(buf []byte, spec CameraInfoSpec) CameraInfo {
	dst := CameraInfo{}
	if len(buf) == 0 {
		return dst
	}
	if spec.FNumberOff >= 0 && spec.FNumberOff < len(buf) {
		dst.FNumber = ciFNumber(buf[spec.FNumberOff])
	}
	if spec.ExposureTimeOff >= 0 && spec.ExposureTimeOff < len(buf) {
		dst.ExposureTime = ciExposureTime(buf[spec.ExposureTimeOff])
	}
	if spec.ISOOff >= 0 && spec.ISOOff < len(buf) {
		dst.ISO = ciISO(buf[spec.ISOOff])
	}
	if spec.HighlightTonePriorityOff >= 0 && spec.HighlightTonePriorityOff < len(buf) {
		dst.HighlightTonePriority = int16(buf[spec.HighlightTonePriorityOff])
	}
	if spec.FlashMeteringModeOff >= 0 && spec.FlashMeteringModeOff < len(buf) {
		dst.FlashMeteringMode = int16(buf[spec.FlashMeteringModeOff])
	}
	if spec.MeasuredEV2Off >= 0 && spec.MeasuredEV2Off < len(buf) {
		dst.MeasuredEV2 = ciMeasuredEV2(buf[spec.MeasuredEV2Off])
	}
	if spec.CameraTemperatureOff >= 0 && spec.CameraTemperatureOff < len(buf) {
		dst.CameraTemperature = ciTemperature(buf[spec.CameraTemperatureOff])
	}
	if spec.MacroMagnificationOff >= 0 && spec.MacroMagnificationOff < len(buf) {
		dst.MacroMagnification = ciMacroMagnification(buf[spec.MacroMagnificationOff])
	}
	if spec.FocalLengthOff >= 0 && spec.FocalLengthOff+2 <= len(buf) {
		dst.FocalLength = ciFocalLength(u16BEAt(buf, spec.FocalLengthOff))
	}
	if spec.CameraOrientationOff >= 0 && spec.CameraOrientationOff < len(buf) {
		dst.CameraOrientation = buf[spec.CameraOrientationOff]
		if dst.CameraOrientation == 0 && spec.CameraOrientationOff == 0x36 {
			if alt := byteAt(buf, 0x3a); alt != 0 {
				dst.CameraOrientation = alt
			}
		}
	}
	if spec.WhiteBalanceOff >= 0 && spec.WhiteBalanceOff+2 <= len(buf) {
		dst.WhiteBalance = WhiteBalance(u16LEAt(buf, spec.WhiteBalanceOff))
	}
	if spec.ColorTemperatureOff >= 0 && spec.ColorTemperatureOff+2 <= len(buf) {
		dst.ColorTemperature = u16LEAt(buf, spec.ColorTemperatureOff)
	}
	if spec.LensTypeOff >= 0 && spec.LensTypeOff+2 <= len(buf) {
		dst.LensType = CanonLensType(u16BEAt(buf, spec.LensTypeOff))
	}
	if spec.MinFocalLengthOff >= 0 && spec.MinFocalLengthOff+2 <= len(buf) {
		dst.MinFocalLength = ciFocalLength(u16BEAt(buf, spec.MinFocalLengthOff))
	}
	if spec.MaxFocalLengthOff >= 0 && spec.MaxFocalLengthOff+2 <= len(buf) {
		dst.MaxFocalLength = ciFocalLength(u16BEAt(buf, spec.MaxFocalLengthOff))
	}
	if spec.JPEGQualityOff >= 0 && spec.JPEGQualityOff < len(buf) {
		dst.JPEGQuality = buf[spec.JPEGQualityOff]
	}
	if spec.PictureStyleOff >= 0 && spec.PictureStyleOff < len(buf) {
		dst.PictureStyle = int16(buf[spec.PictureStyleOff])
	}
	if spec.FirmwareVersionOff >= 0 && spec.FirmwareVersionLen > 0 && spec.FirmwareVersionOff+spec.FirmwareVersionLen <= len(buf) {
		dst.FirmwareVersion = ciASCIIBytes(buf, spec.FirmwareVersionOff, spec.FirmwareVersionLen)
	}
	if spec.FileIndexOff >= 0 && spec.FileIndexOff+4 <= len(buf) {
		if v := u32LEAt(buf, spec.FileIndexOff); v > 0 {
			dst.FileIndex = v + 1
		}
	}
	if spec.DirectoryIndexOff >= 0 && spec.DirectoryIndexOff+4 <= len(buf) {
		if v := u32LEAt(buf, spec.DirectoryIndexOff); v > 0 {
			dst.DirectoryIndex = v - 1
		}
	}
	return dst
}
