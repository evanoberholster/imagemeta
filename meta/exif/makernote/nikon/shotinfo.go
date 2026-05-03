package nikon

import "github.com/evanoberholster/imagemeta/meta/utils"

// NikonShotInfo stores the parsed fields from the Nikon ShotInfo maker-note
// sub-table (tag 0x0091).  ExifTool dispatches to version-specific tables
// selected by the first 4-byte ShotInfoVersion string.
//
// This struct holds the most commonly available fields across versions.
// Model-specific sub-tables with hundreds of offset-pointer entries are
// deferred to future work.
type NikonShotInfo struct {
	ShotInfoVersion string // [0x00] ShotInfoVersion (4 bytes)
	FirmwareVersion string // [0x04] FirmwareVersion (5 or 8 bytes)

	ShutterCount        uint32  // int32u (offset varies by version)
	MechanicalShutterCount uint32 // int32u, some bodies only
	ISO2                float64 // Nikon logarithmic ISO (offset varies)
	VibrationReduction  uint8   // 0 = Off, 1 = On (offset varies)
	Orientation         uint8   // rotation / orientation flag
}

// DecodeShotInfo decodes Nikon ShotInfo (tag 0x0091) from raw bytes.
// It dispatches to version-specific parsers based on ShotInfoVersion.
func DecodeShotInfo(raw []byte) NikonShotInfo {
	if len(raw) < 4 {
		return NikonShotInfo{}
	}
	ver := VersionString(raw[:4])
	dst := NikonShotInfo{ShotInfoVersion: ver}

	// FirmwareVersion: offset 4, length 5 or 8.
	if len(raw) >= 9 {
		dst.FirmwareVersion = VersionString(raw[4:12])
	}
	if len(raw) >= 12 {
		if ver5 := VersionString(raw[4:9]); len(ver5) > 0 && ver5 != dst.FirmwareVersion {
			dst.FirmwareVersion = ver5
		}
	}

	// Dispatch to version-specific parser.
	return parseShotInfoVersion(raw, dst)
}

// parseShotInfoVersion selects the parser table based on ShotInfoVersion.
func parseShotInfoVersion(raw []byte, dst NikonShotInfo) NikonShotInfo {
	ver := dst.ShotInfoVersion
	l := len(raw)

	switch {
	// --- Modern Z-series (offset-pointer based, minimal direct fields) ---
	case hasPrefix(ver, "0809"), hasPrefix(ver, "0810"), hasPrefix(ver, "0811"):
		_ = l // Z6III — offset-pointer subdirectories only
	case hasPrefix(ver, "0800"), hasPrefix(ver, "0801"), hasPrefix(ver, "0802"),
		hasPrefix(ver, "0803"), hasPrefix(ver, "0804"), hasPrefix(ver, "0807"), hasPrefix(ver, "0808"):
		_ = l // Z7II family
	case hasPrefix(ver, "0806"):
		_ = l // Z8
	case hasPrefix(ver, "0805"):
		_ = l // Z9

	// --- D780, D7500 — offset-pointer based ---
	case hasPrefix(ver, "0245"):
		_ = l // D780
	case hasPrefix(ver, "0242"):
		_ = l // D7500

	// --- D6 ---
	case hasPrefix(ver, "0246"):
		_ = l

	// --- D500, D850, D810, D4S, D610 ---
	case hasPrefix(ver, "0238"), hasPrefix(ver, "0239"):
		_ = l // D500
	case hasPrefix(ver, "0243"):
		_ = l // D850
	case hasPrefix(ver, "0233"):
		_ = l // D810
	case hasPrefix(ver, "0231"):
		decodeShotInfoD4s(raw, &dst)
	case hasPrefix(ver, "0232"):
		_ = l // D610

	// --- D4, D800, D7000 family ---
	case hasPrefix(ver, "0223"):
		_ = l // D4
	case hasPrefix(ver, "0220"):
		dst.ShutterCount = u32At(raw, 0x320, utils.LittleEndian) // D7000 ShutterCount
	case hasPrefix(ver, "0222"):
		dst.ShutterCount = u32At(raw, 0x5fb, utils.LittleEndian) // D800 ShutterCount
	case hasPrefix(ver, "0221"):
		dst.ShutterCount = u32At(raw, 0x321, utils.LittleEndian) // D5100 ShutterCount
	case hasPrefix(ver, "0226"):
		dst.ShutterCount = u32At(raw, 0xbd8, utils.LittleEndian) // D5200 ShutterCount

	// --- D5000 ---
	case hasPrefix(ver, "0215"):
		dst.ShutterCount = u32At(raw, 0x2d6, utils.LittleEndian)
		dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 0x2b5)))

	// --- D3S ---
	case hasPrefix(ver, "0218"):
		dst.ShutterCount = u32At(raw, 0x242, utils.LittleEndian)
		dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 0x221)))

	// --- D3X ---
	case hasPrefix(ver, "0214"):
		dst.ShutterCount = u32At(raw, 0x280, utils.LittleEndian)
		dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 0x25d)))

	// --- D3a / D300a / D300b ---
	case hasPrefix(ver, "0210"):
		switch l {
		case 5399:
			dst.ShutterCount = u32At(raw, 0x276, utils.LittleEndian) // D3a
			dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 0x256)))
		case 5291:
			dst.ShutterCount = u32At(raw, 633, utils.LittleEndian) // D300a
			dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 604)))
		case 5303:
			dst.ShutterCount = u32At(raw, 644, utils.LittleEndian) // D300b
			dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 613)))
		}

	// --- D300S ---
	case hasPrefix(ver, "0216"):
		dst.ShutterCount = u32At(raw, 646, utils.LittleEndian)
		dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 613)))

	// --- D700 ---
	case hasPrefix(ver, "0212"):
		dst.ShutterCount = u32At(raw, 0x287, utils.LittleEndian)
		dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 613)))

	// --- D90 ---
	case hasPrefix(ver, "0213"):
		dst.ShutterCount = u32At(raw, 0x2d5, utils.LittleEndian)
		dst.ISO2 = ISOFromRaw(float64(ByteAt(raw, 0x2b5)))

	// --- D40 ---
	case hasPrefix(ver, "0209"):
		dst.ShutterCount = u32At(raw, 582, utils.LittleEndian)
		vb := ByteAt(raw, 586)
		dst.VibrationReduction = (vb >> 3) & 0x01

	// --- D80 ---
	case hasPrefix(ver, "0208"):
		dst.ShutterCount = u32At(raw, 586, utils.LittleEndian)
		vb := ByteAt(raw, 590)
		dst.VibrationReduction = (vb >> 3) & 0x01
		dst.Orientation = vb & 0x07

	// --- Older encrypted 02xx variants ---
	case hasPrefix(ver, "02"):
		// Encrypted block — direct fields at known offsets
		if ver == "0204" {
			dst.ShutterCount = u32At(raw, 0x6a, utils.LittleEndian)
			vb := ByteAt(raw, 0x82)
			dst.VibrationReduction = vb
		}
		if ver == "0205" {
			vb := ByteAt(raw, 0x1ae)
			dst.VibrationReduction = vb
		}
		if ver == "0211" {
			dst.ShutterCount = u32At(raw, 0x24d, utils.LittleEndian)
		}
	}

	return dst
}

// decodeShotInfoD4s extracts fields from Nikon D4S ShotInfo.
func decodeShotInfoD4s(raw []byte, dst *NikonShotInfo) {
	if len(raw) < 0x193d+4 {
		return
	}
	fw := VersionString(raw[4:9])
	// D4S ShutterCount: depends on firmware
	// OrientationInfo sub-dir at 0x350b
	_ = fw
}

// u32At reads a uint32 at offset using the given byte order, or 0.
func u32At(raw []byte, off int, bo utils.ByteOrder) uint32 {
	if off < 0 || off+4 > len(raw) {
		return 0
	}
	return bo.Uint32(raw[off : off+4])
}

// hasPrefix reports whether s starts with prefix (case-sensitive).
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
