package canon

// DecodeFaceDetect1 decodes a Canon FaceDetect1 payload (tag 0x0024).
func DecodeFaceDetect1(words []uint16) FaceDetectInfo {
	if len(words) < 5 {
		return FaceDetectInfo{}
	}
	dst := FaceDetectInfo{
		FacesDetected: words[2],
	}
	dst.FaceDetectFrameSize[0] = words[3]
	dst.FaceDetectFrameSize[1] = words[4]

	faceCount := int(dst.FacesDetected)
	if faceCount > len(dst.FacePositions) {
		faceCount = len(dst.FacePositions)
	}
	for i := 0; i < faceCount; i++ {
		start := 8 + i*2
		if start+1 >= len(words) {
			break
		}
		dst.FacePositions[i] = FacePosition{
			X: int16(words[start]),
			Y: int16(words[start+1]),
		}
	}
	return dst
}

// DecodeFaceDetect2 decodes a Canon FaceDetect2 payload (tag 0x0025).
func DecodeFaceDetect2(raw []byte) FaceDetectInfo {
	if len(raw) < 3 {
		return FaceDetectInfo{}
	}
	return FaceDetectInfo{
		FaceWidth:     raw[1],
		FacesDetected: uint16(raw[2]),
	}
}

// DecodeFaceDetect3 decodes a Canon FaceDetect3 payload (tag 0x002f).
func DecodeFaceDetect3(words []uint16) FaceDetectInfo {
	if len(words) < 4 {
		return FaceDetectInfo{}
	}
	return FaceDetectInfo{
		FacesDetected: words[3],
	}
}
