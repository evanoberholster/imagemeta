package nikon

import "strconv"

// NikonCameraModelSourceFile is the source used to derive this table.
// "Image/ExifTool/Nikon.pm", "Image/ExifTool/NikonSettings.pm", "Image/ExifTool/NikonCustom.pm"
// "https://en.wikipedia.org/wiki/Category:Nikon_digital_cameras"
// "https://en.wikipedia.org/wiki/Category:Nikon_DSLR_cameras"
// "https://en.wikipedia.org/wiki/Category:Nikon_1_cameras"
// "https://en.wikipedia.org/wiki/Category:Nikon_Z-mount_cameras"
// NikonCameraModelSourceVersion is the ExifTool version used to generate this table. "13.50"
// Note: ExifTool does not expose a single Nikon model-ID table like Canon's %canonModelID.

type NikonCameraModel uint32

const (
	NikonModelUnknown NikonCameraModel = iota
	NikonModelD3
	NikonModelD300S
	NikonModelD3S
	NikonModelD3X
	NikonModelD4
	NikonModelD4S
	NikonModelD5
	NikonModelD6
	NikonModelD40
	NikonModelD80
	NikonModelD90
	NikonModelD500
	NikonModelD5000
	NikonModelD5100
	NikonModelD5200
	NikonModelD610
	NikonModelD700
	NikonModelD7000
	NikonModelD7500
	NikonModelD780
	NikonModelD800
	NikonModelD810
	NikonModelD850
	NikonModelDf
	NikonModelD3300
	NikonModelD5300
	NikonModelD7100
	NikonModelZ30
	NikonModelZ5
	NikonModelZ5_2
	NikonModelZ50
	NikonModelZ50_2
	NikonModelZ6
	NikonModelZ6_2
	NikonModelZ6_3
	NikonModelZ6III
	NikonModelZ7
	NikonModelZ7_2
	NikonModelZ7II
	NikonModelZ8
	NikonModelZ9
	NikonModelZf
	NikonModelZfc
	NikonModelCoolpixA
	NikonModel1AW1
	NikonModel1J1
	NikonModel1J2
	NikonModel1J3
	NikonModel1J4
	NikonModel1J5
	NikonModel1S1
	NikonModel1S2
	NikonModel1V1
	NikonModel1V2
	NikonModel1V3
	NikonModelD1
	NikonModelD1H
	NikonModelD1X
	NikonModelD2H
	NikonModelD2Hs
	NikonModelD2X
	NikonModelD2Xs
	NikonModelD40X
	NikonModelD50
	NikonModelD60
	NikonModelD70
	NikonModelD100
	NikonModelD200
	NikonModelD300
	NikonModelD3000
	NikonModelD3100
	NikonModelD3200
	NikonModelD3400
	NikonModelD3500
	NikonModelD5500
	NikonModelD5600
	NikonModelD600
	NikonModelD7200
	NikonModelD750
	NikonModelD800E
	NikonModelD810A
	NikonModelE2
	NikonModelE2N
	NikonModelE2NS
	NikonModelE2S
	NikonModelE3
	NikonModelE3S
	NikonModelZ5II
	NikonModelZ50II
	NikonModelZR
	NikonModelCoolpixA1000
	NikonModelCoolpixA300
	NikonModelCoolpixA900
	NikonModelCoolpixB500
	NikonModelCoolpixB600
	NikonModelCoolpixP1
	NikonModelCoolpixP100
	NikonModelCoolpixP300
	NikonModelCoolpixP310
	NikonModelCoolpixP340
	NikonModelCoolpixP500
	NikonModelCoolpixP5000
	NikonModelCoolpixP510
	NikonModelCoolpixP520
	NikonModelCoolpixP530
	NikonModelCoolpixP600
	NikonModelCoolpixP6000
	NikonModelCoolpixP7000
	NikonModelCoolpixP7100
	NikonModelCoolpixP7700
	NikonModelCoolpixP7800
	NikonModelCoolpixP900
	NikonModelCoolpixP950
	NikonModelCoolpixP1000
	NikonModelCoolpixS10
	NikonModelCoolpixS610
	NikonModelCoolpixS620
	NikonModelCoolpixS4000
	NikonModelCoolpixS7000
	NikonModelCoolpixS8000
	NikonModelCoolpixS8100
	NikonModelCoolpixS9100
	NikonModelCoolpixS9300
	NikonModelCoolpixS9700
	NikonModelCoolpixS9900
	NikonModelCoolpixW100
	NikonModelCoolpixW150
	NikonModelCoolpixW300
)

// NikonCameraModelCount is the number of Nikon model entries in this table.
const NikonCameraModelCount = 128

const _NikonCameraModel_name = "UnknownNIKON D3NIKON D300SNIKON D3SNIKON D3XNIKON D4NIKON D4SNIKON D5NIKON D6NIKON D40NIKON D80NIKON D90NIKON D500NIKON D5000NIKON D5100NIKON D5200NIKON D610NIKON D700NIKON D7000NIKON D7500NIKON D780NIKON D800NIKON D810NIKON D850NIKON DfNIKON D3300NIKON D5300NIKON D7100NIKON Z 30NIKON Z 5NIKON Z5_2NIKON Z 50NIKON Z50_2NIKON Z 6NIKON Z 6_2NIKON Z6_3NIKON Z 6IIINIKON Z 7NIKON Z 7_2NIKON Z 7IINIKON Z 8NIKON Z 9NIKON Z fNIKON Z fcNIKON COOLPIX ANIKON 1 AW1NIKON 1 J1NIKON 1 J2NIKON 1 J3NIKON 1 J4NIKON 1 J5NIKON 1 S1NIKON 1 S2NIKON 1 V1NIKON 1 V2NIKON 1 V3NIKON D1NIKON D1HNIKON D1XNIKON D2HNIKON D2HsNIKON D2XNIKON D2XsNIKON D40XNIKON D50NIKON D60NIKON D70NIKON D100NIKON D200NIKON D300NIKON D3000NIKON D3100NIKON D3200NIKON D3400NIKON D3500NIKON D5500NIKON D5600NIKON D600NIKON D7200NIKON D750NIKON D800ENIKON D810ANIKON E2NIKON E2NNIKON E2NSNIKON E2SNIKON E3NIKON E3SNIKON Z5IINIKON Z50IINIKON ZRNIKON COOLPIX A1000NIKON COOLPIX A300NIKON COOLPIX A900NIKON COOLPIX B500NIKON COOLPIX B600NIKON COOLPIX P1NIKON COOLPIX P100NIKON COOLPIX P300NIKON COOLPIX P310NIKON COOLPIX P340NIKON COOLPIX P500NIKON COOLPIX P5000NIKON COOLPIX P510NIKON COOLPIX P520NIKON COOLPIX P530NIKON COOLPIX P600NIKON COOLPIX P6000NIKON COOLPIX P7000NIKON COOLPIX P7100NIKON COOLPIX P7700NIKON COOLPIX P7800NIKON COOLPIX P900NIKON COOLPIX P950NIKON COOLPIX P1000NIKON COOLPIX S10NIKON COOLPIX S610NIKON COOLPIX S620NIKON COOLPIX S4000NIKON COOLPIX S7000NIKON COOLPIX S8000NIKON COOLPIX S8100NIKON COOLPIX S9100NIKON COOLPIX S9300NIKON COOLPIX S9700NIKON COOLPIX S9900NIKON COOLPIX W100NIKON COOLPIX W150NIKON COOLPIX W300"

var _NikonCameraModel_index = [...]uint16{
	0, 7, 15, 26, 35, 44, 52, 61, 69, 77, 86, 95, 104, 114, 125, 136, 147, 157, 167,
	178, 189, 199, 209, 219, 229, 237, 248, 259, 270, 280, 289, 299, 309, 320, 329,
	340, 350, 362, 371, 382, 393, 402, 411, 420, 430, 445, 456, 466, 476, 486, 496,
	506, 516, 526, 536, 546, 556, 564, 573, 582, 591, 601, 610, 620, 630, 639, 648,
	657, 667, 677, 687, 698, 709, 720, 731, 742, 753, 764, 774, 785, 795, 806, 817,
	825, 834, 844, 853, 861, 870, 880, 891, 899, 918, 936, 954, 972, 990, 1006, 1024,
	1042, 1060, 1078, 1096, 1115, 1133, 1151, 1169, 1187, 1206, 1225, 1244, 1263,
	1282, 1300, 1318, 1337, 1354, 1372, 1390, 1409, 1428, 1447, 1466, 1485, 1504,
	1523, 1542, 1560, 1578, 1596,
}

func (m NikonCameraModel) String() string {
	i := int(m)
	if i < 0 || i >= len(_NikonCameraModel_index)-1 {
		return "NikonCameraModel(" + strconv.FormatUint(uint64(m), 10) + ")"
	}
	return _NikonCameraModel_name[_NikonCameraModel_index[i]:_NikonCameraModel_index[i+1]]
}
