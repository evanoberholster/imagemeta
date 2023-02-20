package canon

// WIP: Work in Progress

// CameraModel is a Canon Camera Model found in Exif
type CameraModel uint32

// CameraModelFromString returns a canon camera model from the given string
func CameraModelFromString(str string) (CameraModel, bool) {
	if cm, ok := mapStringCameraModel[str]; ok {
		return cm, true
	}
	return CanonModelUnknown, false
}

var mapStringCameraModel = map[string]CameraModel{
	"Canon EOS-1D X Mark III": EOS1DXMarkIII,
	"Canon EOS 90D":           EOS90D,
	"Canon EOS M200":          EOSM200,
	"Canon EOS M50m2":         EOSM50MarkII,
	"Canon EOS M6 Mark II":    EOSM6MarkII,
	"Canon EOS R10":           EOSR10,
	"Canon EOS R3":            EOSR3,
	"Canon EOS R50":           EOSR50,
	"Canon EOS R5":            EOSR5,
	"Canon EOS R6m2":          EOSR6MarkII,
	"Canon EOS R6":            EOSR6,
	"Canon EOS R7":            EOSR7,
	"Canon EOS R8":            EOSR8,
	"Canon EOS RP":            EOSRP,
	"Canon EOS R":             EOSR,
	"Canon EOS Rebel SL3":     EOS250D,
	"Canon EOS 6D":            EOS6D,
}

var mapCameraModelString = map[CameraModel]string{
	EOS1DXMarkIII: "Canon EOS-1D X Mark III",
	EOS90D:        "Canon EOS 90D",
	EOSM200:       "Canon EOS M200",
	EOSM50MarkII:  "Canon EOS M50 Mark II",
	EOSM6MarkII:   "Canon EOS M6 Mark II",
	EOSR10:        "Canon EOS R10",
	EOSR3:         "Canon EOS R3",
	EOSR50:        "Canon EOS R50",
	EOSR5:         "Canon EOS R5",
	EOSR6MarkII:   "Canon EOS R6 Mark II",
	EOSR6:         "Canon EOS R6",
	EOSR7:         "Canon EOS R7",
	EOSR8:         "Canon EOS R8",
	EOSRP:         "Canon EOS RP",
	EOSR:          "Canon EOS R",
	EOS250D:       "Canon EOS SL3",
	EOS6D:         "Canon EOS 6D",
}

func (cm CameraModel) String() string {
	if str, ok := mapCameraModelString[cm]; ok {
		return str
	}
	return ""
}

// Inspired by Exiftool source code at https://github.com/exiftool/exiftool/blob/master/lib/Image/ExifTool/Canon.pm

// Canon Camera Models
const (
	CanonModelUnknown  CameraModel = iota + 0x10000
	PowerShotA30                   // 0x1010000
	PowerShotS300                  // 0x1040000
	PowerShotA20                   // 0x1060000
	PowerShotA10                   // 0x1080000
	PowerShotS110                  // 0x1090000
	PowerShotG2                    // 0x1100000
	PowerShotS40                   // 0x1110000
	PowerShotS30                   // 0x1120000
	PowerShotA40                   // 0x1130000
	EOSD30                         // 0x1140000
	PowerShotA100                  // 0x1150000
	PowerShotS200                  // 0x1160000
	PowerShotA200                  // 0x1170000
	PowerShotS330                  // 0x1180000
	PowerShotG3                    // 0x1190000
	PowerShotS45                   // 0x1210000
	PowerShotSD100                 // 0x1230000
	PowerShotS230                  // 0x1240000
	PowerShotA70                   // 0x1250000
	PowerShotA60                   // 0x1260000
	PowerShotS400                  // 0x1270000
	PowerShotG5                    // 0x1290000
	PowerShotA300                  // 0x1300000
	PowerShotS50                   // 0x1310000
	PowerShotA80                   // 0x1340000
	PowerShotSD10                  // 0x1350000
	PowerShotS1IS                  // 0x1360000
	PowerShotPro1                  // 0x1370000
	PowerShotS70                   // 0x1380000
	PowerShotS60                   // 0x1390000
	PowerShotG6                    // 0x1400000
	PowerShotS500                  // 0x1410000
	PowerShotA75                   // 0x1420000
	PowerShotSD110                 // 0x1440000
	PowerShotA400                  // 0x1450000
	PowerShotA310                  // 0x1470000
	PowerShotA85                   // 0x1490000
	PowerShotS410                  // 0x1520000
	PowerShotA95                   // 0x1530000
	PowerShotSD300                 // 0x1540000
	PowerShotSD200                 // 0x1550000
	PowerShotA520                  // 0x1560000
	PowerShotA510                  // 0x1570000
	PowerShotSD20                  // 0x1590000
	PowerShotS2IS                  // 0x1640000
	PowerShotSD430                 // 0x1650000
	PowerShotSD500                 // 0x1660000
	EOSD60                         // 0x1668000
	PowerShotSD30                  // 0x1700000
	PowerShotA430                  // 0x1740000
	PowerShotA410                  // 0x1750000
	PowerShotS80                   // 0x1760000
	PowerShotA620                  // 0x1780000
	PowerShotA610                  // 0x1790000
	PowerShotSD630                 // 0x1800000
	PowerShotSD450                 // 0x1810000
	PowerShotTX1                   // 0x1820000
	PowerShotSD400                 // 0x1870000
	PowerShotA420                  // 0x1880000
	PowerShotSD900                 // 0x1890000
	PowerShotSD550                 // 0x1900000
	PowerShotA700                  // 0x1920000
	PowerShotSD700IS               // 0x1940000
	PowerShotS3IS                  // 0x1950000
	PowerShotA540                  // 0x1960000
	PowerShotSD600                 // 0x1970000
	PowerShotG7                    // 0x1980000
	PowerShotA530                  // 0x1990000
	PowerShotSD800IS               // 0x2000000
	PowerShotSD40                  // 0x2010000
	PowerShotA710IS                // 0x2020000
	PowerShotA640                  // 0x2030000
	PowerShotA630                  // 0x2040000
	PowerShotS5IS                  // 0x2090000
	PowerShotA460                  // 0x2100000
	PowerShotSD850IS               // 0x2120000
	PowerShotA570IS                // 0x2130000
	PowerShotA560                  // 0x2140000
	PowerShotSD750                 // 0x2150000
	PowerShotSD1000                // 0x2160000
	PowerShotA550                  // 0x2180000
	PowerShotA450                  // 0x2190000
	PowerShotG9                    // 0x2230000
	PowerShotA650IS                // 0x2240000
	PowerShotA720IS                // 0x2260000
	PowerShotSX100IS               // 0x2290000
	PowerShotSD950IS               // 0x2300000
	PowerShotSD870IS               // 0x2310000
	PowerShotSD890IS               // 0x2320000
	PowerShotSD790IS               // 0x2360000
	PowerShotSD770IS               // 0x2370000
	PowerShotA590IS                // 0x2380000
	PowerShotA580                  // 0x2390000
	PowerShotA470                  // 0x2420000
	PowerShotSD1100IS              // 0x2430000
	PowerShotSX1IS                 // 0x2460000
	PowerShotSX10IS                // 0x2470000
	PowerShotA1000IS               // 0x2480000
	PowerShotG10                   // 0x2490000
	PowerShotA2000IS               // 0x2510000
	PowerShotSX110IS               // 0x2520000
	PowerShotSD990IS               // 0x2530000
	PowerShotSD880IS               // 0x2540000
	PowerShotE1                    // 0x2550000
	PowerShotD10                   // 0x2560000
	PowerShotSD960IS               // 0x2570000
	PowerShotA2100IS               // 0x2580000
	PowerShotA480                  // 0x2590000
	PowerShotSX200IS               // 0x2600000
	PowerShotSD970IS               // 0x2610000
	PowerShotSD780IS               // 0x2620000
	PowerShotA1100IS               // 0x2630000
	PowerShotSD1200IS              // 0x2640000
	PowerShotG11                   // 0x2700000
	PowerShotSX120IS               // 0x2710000
	PowerShotS90                   // 0x2720000
	PowerShotSX20IS                // 0x2750000
	PowerShotSD980IS               // 0x2760000
	PowerShotSD940IS               // 0x2770000
	PowerShotA495                  // 0x2800000
	PowerShotA490                  // 0x2810000
	PowerShotA3100IS               // (different cameras, same ID) // 0x2820000
	PowerShotA3150IS               // (different cameras, same ID) // 0x2820000
	PowerShotA3000IS               // 0x2830000
	PowerShotSD1400IS              // 0x2840000
	PowerShotSD1300IS              // 0x2850000
	PowerShotSD3500IS              // 0x2860000
	PowerShotSX210IS               // 0x2870000
	PowerShotSD4000IS              // 0x2880000
	PowerShotSD4500IS              // 0x2890000
	PowerShotG12                   // 0x2920000
	PowerShotSX30IS                // 0x2930000
	PowerShotSX130IS               // 0x2940000
	PowerShotS95                   // 0x2950000
	PowerShotA3300IS               // 0x2980000
	PowerShotA3200IS               // 0x2990000
	PowerShotELPH500HS             // 0x3000000
	PowerShotPro90IS               // 0x3010000
	PowerShotA800                  // 0x3010001
	PowerShotELPH100HS             // 0x3020000
	PowerShotSX230HS               // 0x3030000
	PowerShotELPH300HS             // 0x3040000
	PowerShotA2200                 // 0x3050000
	PowerShotA1200                 // 0x3060000
	PowerShotSX220HS               // 0x3070000
	PowerShotG1X                   // 0x3080000
	PowerShotSX150IS               // 0x3090000
	PowerShotELPH510HS             // 0x3100000
	PowerShotS100                  // 0x3110000
	PowerShotSX40HS                // 0x3130000
	PowerShotELPH310HS             // 0x3120000

	IXY32S             // 0x3140000 (PowerShot ELPH 500 HS / IXUS 320 HS ??)
	PowerShotA1300     // 0x3160000
	PowerShotA810      // 0x3170000
	PowerShotELPH320HS // 0x3180000
	PowerShotELPH110HS // 0x3190000
	PowerShotD20       // 0x3200000
	PowerShotA4000IS   // 0x3210000
	PowerShotSX260HS   // 0x3220000
	PowerShotSX240HS   // 0x3230000
	PowerShotELPH530HS // 0x3240000
	PowerShotELPH520HS // 0x3250000
	PowerShotA3400IS   // 0x3260000
	PowerShotA2400IS   // 0x3270000
	PowerShotA2300     // 0x3280000
	PowerShotS100V     // 0x3320000
	PowerShotG15       // 0x3330000
	PowerShotSX50HS    // 0x3340000
	PowerShotSX160IS   // 0x3350000
	//PowerShotS110       //(new) // 0x3360000
	PowerShotSX500IS   // 0x3370000
	PowerShotN         // 0x3380000
	IXUS245HS          //  (no PowerShot) 0x3390000
	PowerShotSX280HS   // 0x3400000
	PowerShotSX270HS   // 0x3410000
	PowerShotA3500IS   // 0x3420000
	PowerShotA2600     // 0x3430000
	PowerShotSX275HS   // 0x3440000
	PowerShotA1400     // 0x3450000
	PowerShotELPH130IS // 0x3460000
	PowerShotELPH115   // 0x3470000
	PowerShotELPH330HS // 0x3490000
	PowerShotA2500     // 0x3510000
	PowerShotG16       // 0x3540000
	PowerShotS120      // 0x3550000
	PowerShotSX170IS   // 0x3560000
	PowerShotSX510HS   // 0x3580000
	//PowerShotS200       // (new) 0x3590000
	IXY620F             // (no PowerShot or IXUS?) 0x3600000
	PowerShotN100       // 0x3610000
	PowerShotG1XMarkII  // 0x3640000
	PowerShotD30        // 0x3650000
	PowerShotSX700HS    // 0x3660000
	PowerShotSX600HS    // 0x3670000
	PowerShotELPH140IS  // 0x3680000
	PowerShotELPH135    // 0x3690000
	PowerShotELPH340HS  // 0x3700000
	PowerShotELPH150IS  // 0x3710000
	EOSM3               // 0x3740000
	PowerShotSX60HS     // 0x3750000
	PowerShotSX520HS    // 0x3760000
	PowerShotSX400IS    // 0x3770000
	PowerShotG7X        // 0x3780000
	PowerShotN2         // 0x3790000
	PowerShotSX530HS    // 0x3800000
	PowerShotSX710HS    // 0x3820000
	PowerShotSX610HS    // 0x3830000
	EOSM10              // 0x3840000
	PowerShotG3X        // 0x3850000
	PowerShotELPH165HS  // 0x3860000
	PowerShotELPH160    // 0x3870000
	PowerShotELPH350HS  // 0x3880000
	PowerShotELPH170IS  // 0x3890000
	PowerShotSX410IS    // 0x3910000
	PowerShotG9X        // 0x3930000
	EOSM5               // 0x3940000
	PowerShotG5X        // 0x3950000
	PowerShotG7XMarkII  // 0x3970000
	EOSM100             // 0x3980000
	PowerShotELPH360HS  // 0x3990000
	PowerShotSX540HS    // 0x4010000
	PowerShotSX420IS    // 0x4020000
	PowerShotELPH190IS  // 0x4030000
	PowerShotG1         // 0x4040000
	PowerShotELPH180IS  // 0x4040001
	PowerShotSX720HS    // 0x4050000
	PowerShotSX620HS    // 0x4060000
	EOSM6               // 0x4070000
	PowerShotG9XMarkII  // 0x4100000
	EOSM50              // 0x412
	PowerShotELPH185    // 0x4150000
	PowerShotSX430IS    // 0x4160000
	PowerShotSX730HS    // 0x4170000
	PowerShotG1XMarkIII // 0x4180000
	//PowerShotS100       // 0x6040000
	PowerShotSX740HS    // 0x801
	PowerShotG5XMarkII  // 0x804
	PowerShotSX70HS     // 0x805
	PowerShotG7XMarkIII // 0x808
	EOSM6MarkII         // 0x811
	EOSM200             // 0x812

	DC19    // (DC21/DC22) 0x4007d673
	XHA1    // 0x4007d674
	HV10    // 0x4007d675
	MD130   // (MD140/MD150/MD160/ZR850) 0x4007d676
	DC50    // 0x4007d777
	HV20    // 0x4007d778
	DC211   // 0x4007d779
	HG10    // 0x4007d77a
	HR10    // 0x4007d77b
	MD255   // (ZR950) 0x4007d77d
	HF11    // 0x4007d81c
	HV30    // 0x4007d878
	XHA1S   // 0x4007d87c
	DC301   // (DC310/DC311/DC320/DC330) 0x4007d87e
	FS100   // 0x4007d87f
	HF10    // 0x4007d880
	HG20    // (HG21) 0x4007d882
	HF21    // 0x4007d925
	HFS11   // 0x4007d926
	HV40    // 0x4007d978
	DC410   // (DC411/DC420) 0x4007d987
	FS19    // (FS20/FS21/FS22/FS200) 0x4007d988
	HF20    // (HF200) 0x4007d989
	HFS10   // (S100) 0x4007d98a
	HFR10   // (R16/R17/R18/R100/R106) 0x4007da8e
	HFM30   // (M31/M36/M300/M306) 0x4007da8f
	HFS20   // (S21/S200) 0x4007da90
	FS31    // (FS36/FS37/FS300/FS305/FS306/FS307) 0x4007da92
	EOSC300 // 0x4007dca0
	HFG25   // 0x4007dda9
	XC10    // 0x4007dfb4
	EOSC200 // 0x4007e1c3

	EOS1D         // 0x80000001
	EOS1DS        // 0x80000167
	EOS10D        // 0x80000168
	EOS1DMarkIII  // 0x80000169
	EOS300D       // 0x80000170
	EOS1DMarkII   // 0x80000174
	EOS20D        // 0x80000175
	EOS450D       // 0x80000176
	EOS1DsMarkII  // 0x80000188
	EOS350D       // 0x80000189
	EOS40D        // 0x80000190
	EOS5D         // 0x80000213
	EOS1DsMarkIII // 0x80000215
	EOS5DMarkII   // 0x80000218
	WFTE1         // 0x80000219
	EOS1DMarkIIN  // 0x80000232
	EOS30D        // 0x80000234
	EOS400D       // 0x80000236
	WFTE2         // 0x80000241
	WFTE3         // 0x80000246
	EOS7D         // 0x80000250
	EOS500D       // 0x80000252
	EOS1000D      // 0x80000254
	EOS50D        // 0x80000261
	EOS1DX        // 0x80000269
	EOS550D       // 0x80000270
	WFTE4         // 0x80000271
	WFTE5         // 0x80000273
	EOS1DMarkIV   // 0x80000281
	EOS5DMarkIII  // 0x80000285
	EOS600D       // 0x80000286
	EOS60D        // 0x80000287
	EOS1100D      // 0x80000288
	EOS7DMarkII   // 0x80000289
	WFTE2II       // 0x80000297
	WFTE4II       // 0x80000298
	EOS650D       // 0x80000301
	EOS6D         // 0x80000302
	EOS1DC        // 0x80000324
	EOS70D        // 0x80000325
	EOS700D       // 0x80000326
	EOS1200D      // 0x80000327
	EOS1DXMarkII  // 0x80000328
	EOSM          // 0x80000331
	EOS80D        // 0x80000350
	EOSM2         // 0x80000355
	EOS100D       // 0x80000346
	EOS760D       // 0x80000347
	EOS5DMarkIV   // 0x80000349
	EOS5DS        // 0x80000382
	EOS750D       // 0x80000393
	EOS5DSR       // 0x80000401
	EOS1300D      // 0x80000404
	EOS800D       // 0x80000405
	EOS6DMarkII   // 0x80000406
	EOS77D        // 0x80000408
	EOS200D       // 0x80000417
	EOSR5         // 0x80000421
	EOS4000D      // 0x80000422
	EOSR          // 0x80000424
	EOS1DXMarkIII // 0x80000428
	EOS2000D      // 0x80000432
	EOSRP         // 0x80000433
	EOS850D       // 0x80000435
	EOS250D       // 0x80000436 // SL3
	EOS90D        // 0x80000437
	EOSR3         // 0x80000450
	EOSR6         // 0x80000453
	EOSR7         // 0x80000464
	EOSR10        // 0x80000465
	PowerShotZOOM // 0x80000467
	EOSM50MarkII  // 0x80000468
	EOSR50        // 0x80000480
	EOSR6MarkII   // 0x80000481
	EOSR8         // 0x80000487
	EOSD2000C     // 0x80000520
	EOSD6000C     // 0x80000560

	// New additions

)
