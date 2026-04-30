package canon

// BatteryTypePayloadSize is the fixed MakerNote BatteryType tag payload length.
const BatteryTypePayloadSize = 76

// BatteryTypeHeaderLen is the number of leading bytes to skip.
const BatteryTypeHeaderLen = 4

// ParseBatteryType extracts the NUL-terminated battery model string from a
// raw 72-byte CameraMaker:BatteryType payload (tag 0x0038).
//
// string(payload) does not allocate in modern Go when used only for switch
// comparison and the source is not subsequently modified.
func ParseBatteryType(payload []byte) string {
	if len(payload) < BatteryTypePayloadSize-BatteryTypeHeaderLen {
		return ""
	}

	// Find first NUL byte. The 72-byte payload contains the model string
	// followed by NUL padding and a trailing 0x01 marker.
	i := 0
	for i < len(payload) && payload[i] != 0 {
		i++
	}
	if i == 0 {
		return ""
	}

	switch string(payload[:i]) {
	case "LP-E6":
		return "LP-E6"
	case "LP-E6N":
		return "LP-E6N"
	case "LP-E6NH":
		return "LP-E6NH"
	case "LP-E6P":
		return "LP-E6P"
	case "LP-E12":
		return "LP-E12"
	case "LP-E17":
		return "LP-E17"
	case "LP-E19":
		return "LP-E19"
	case "NB-13L":
		return "NB-13L"
	default:
		return string(payload[:i])
	}
}
