package mknote

// Based on https://github.com/exiftool/exiftool/blob/master/lib/Image/ExifTool/Canon.pm
// Exiftool source code

// Canon Camera Models
const (
	UnknownCanonModel  CameraModel = iota + 0x10000
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
	EOS250D       // 0x80000436
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
)

var (
	mapCanonModelID = map[uint32]CameraModel{
		0x1010000: PowerShotA30,
		0x1040000: PowerShotS300,
		0x1060000: PowerShotA20,
		0x1080000: PowerShotA10,
		0x1090000: PowerShotS110,
		0x1100000: PowerShotG2,
		0x1110000: PowerShotS40,
		0x1120000: PowerShotS30,
		0x1130000: PowerShotA40,
		0x1140000: EOSD30,
		0x1150000: PowerShotA100,
		0x1160000: PowerShotS200,
		0x1170000: PowerShotA200,
		0x1180000: PowerShotS330,
		0x1190000: PowerShotG3,
		0x1210000: PowerShotS45,
		0x1230000: PowerShotSD100,
		0x1240000: PowerShotS230,
		0x1250000: PowerShotA70,
		0x1260000: PowerShotA60,
		0x1270000: PowerShotS400,
		0x1290000: PowerShotG5,
		0x1300000: PowerShotA300,
		0x1310000: PowerShotS50,
		0x1340000: PowerShotA80,
		0x1350000: PowerShotSD10,
		0x1360000: PowerShotS1IS,
		0x1370000: PowerShotPro1,
		0x1380000: PowerShotS70,
		0x1390000: PowerShotS60,
		0x1400000: PowerShotG6,
		0x1410000: PowerShotS500,
		0x1420000: PowerShotA75,
		0x1440000: PowerShotSD110,
		0x1450000: PowerShotA400,
		0x1470000: PowerShotA310,
		0x1490000: PowerShotA85,
		0x1520000: PowerShotS410,
		0x1530000: PowerShotA95,
		0x1540000: PowerShotSD300,
		0x1550000: PowerShotSD200,
		0x1560000: PowerShotA520,
		0x1570000: PowerShotA510,
		0x1590000: PowerShotSD20,
		0x1640000: PowerShotS2IS,
		0x1650000: PowerShotSD430,
		0x1660000: PowerShotSD500,
		0x1668000: EOSD60,
		0x1700000: PowerShotSD30,
		0x1740000: PowerShotA430,
		0x1750000: PowerShotA410,
		0x1760000: PowerShotS80,
		0x1780000: PowerShotA620,
		0x1790000: PowerShotA610,
		0x1800000: PowerShotSD630,
		0x1810000: PowerShotSD450,
		0x1820000: PowerShotTX1,
		0x1870000: PowerShotSD400,
		0x1880000: PowerShotA420,
		0x1890000: PowerShotSD900,
		0x1900000: PowerShotSD550,
		0x1920000: PowerShotA700,
		0x1940000: PowerShotSD700IS,
		0x1950000: PowerShotS3IS,
		0x1960000: PowerShotA540,
		0x1970000: PowerShotSD600,
		0x1980000: PowerShotG7,
		0x1990000: PowerShotA530,
		0x2000000: PowerShotSD800IS,
		0x2010000: PowerShotSD40,
		0x2020000: PowerShotA710IS,
		0x2030000: PowerShotA640,
		0x2040000: PowerShotA630,
		0x2090000: PowerShotS5IS,
		0x2100000: PowerShotA460,
		0x2120000: PowerShotSD850IS,
		0x2130000: PowerShotA570IS,
		0x2140000: PowerShotA560,
		0x2150000: PowerShotSD750,
		0x2160000: PowerShotSD1000,
		0x2180000: PowerShotA550,
		0x2190000: PowerShotA450,
		0x2230000: PowerShotG9,
		0x2240000: PowerShotA650IS,
		0x2260000: PowerShotA720IS,
		0x2290000: PowerShotSX100IS,
		0x2300000: PowerShotSD950IS,
		0x2310000: PowerShotSD870IS,
		0x2320000: PowerShotSD890IS,
		0x2360000: PowerShotSD790IS,
		0x2370000: PowerShotSD770IS,
		0x2380000: PowerShotA590IS,
		0x2390000: PowerShotA580,
		0x2420000: PowerShotA470,
		0x2430000: PowerShotSD1100IS,
		0x2460000: PowerShotSX1IS,
		0x2470000: PowerShotSX10IS,
		0x2480000: PowerShotA1000IS,
		0x2490000: PowerShotG10,
		0x2510000: PowerShotA2000IS,
		0x2520000: PowerShotSX110IS,
		0x2530000: PowerShotSD990IS,
		0x2540000: PowerShotSD880IS,
		0x2550000: PowerShotE1,
		0x2560000: PowerShotD10,
		0x2570000: PowerShotSD960IS,
		0x2580000: PowerShotA2100IS,
		0x2590000: PowerShotA480,
		0x2600000: PowerShotSX200IS,
		0x2610000: PowerShotSD970IS,
		0x2620000: PowerShotSD780IS,
		0x2630000: PowerShotA1100IS,
		0x2640000: PowerShotSD1200IS,
		0x2700000: PowerShotG11,
		0x2710000: PowerShotSX120IS,
		0x2720000: PowerShotS90,
		0x2750000: PowerShotSX20IS,
		0x2760000: PowerShotSD980IS,
		0x2770000: PowerShotSD940IS,
		0x2800000: PowerShotA495,
		0x2810000: PowerShotA490,
		0x2820000: PowerShotA3100IS,
		//0x2820000:  PowerShotA3150IS,
		0x2830000: PowerShotA3000IS,
		0x2840000: PowerShotSD1400IS,
		0x2850000: PowerShotSD1300IS,
		0x2860000: PowerShotSD3500IS,
		0x2870000: PowerShotSX210IS,
		0x2880000: PowerShotSD4000IS,
		0x2890000: PowerShotSD4500IS,
		0x2920000: PowerShotG12,
		0x2930000: PowerShotSX30IS,
		0x2940000: PowerShotSX130IS,
		0x2950000: PowerShotS95,
		0x2980000: PowerShotA3300IS,
		0x2990000: PowerShotA3200IS,
		0x3000000: PowerShotELPH500HS,
		0x3010000: PowerShotPro90IS,
		0x3010001: PowerShotA800,
		0x3020000: PowerShotELPH100HS,
		0x3030000: PowerShotSX230HS,
		0x3040000: PowerShotELPH300HS,
		0x3050000: PowerShotA2200,
		0x3060000: PowerShotA1200,
		0x3070000: PowerShotSX220HS,
		0x3080000: PowerShotG1X,
		0x3090000: PowerShotSX150IS,
		0x3100000: PowerShotELPH510HS,
		0x3110000: PowerShotS100,
		0x3130000: PowerShotSX40HS,
		0x3120000: PowerShotELPH310HS,

		0x3140000: IXY32S,
		0x3160000: PowerShotA1300,
		0x3170000: PowerShotA810,
		0x3180000: PowerShotELPH320HS,
		0x3190000: PowerShotELPH110HS,
		0x3200000: PowerShotD20,
		0x3210000: PowerShotA4000IS,
		0x3220000: PowerShotSX260HS,
		0x3230000: PowerShotSX240HS,
		0x3240000: PowerShotELPH530HS,
		0x3250000: PowerShotELPH520HS,
		0x3260000: PowerShotA3400IS,
		0x3270000: PowerShotA2400IS,
		0x3280000: PowerShotA2300,
		0x3320000: PowerShotS100V,
		0x3330000: PowerShotG15,
		0x3340000: PowerShotSX50HS,
		0x3350000: PowerShotSX160IS,
		0x3360000: PowerShotS110,
		0x3370000: PowerShotSX500IS,
		0x3380000: PowerShotN,
		0x3390000: IXUS245HS,
		0x3400000: PowerShotSX280HS,
		0x3410000: PowerShotSX270HS,
		0x3420000: PowerShotA3500IS,
		0x3430000: PowerShotA2600,
		0x3440000: PowerShotSX275HS,
		0x3450000: PowerShotA1400,
		0x3460000: PowerShotELPH130IS,
		0x3470000: PowerShotELPH115,
		0x3490000: PowerShotELPH330HS,
		0x3510000: PowerShotA2500,
		0x3540000: PowerShotG16,
		0x3550000: PowerShotS120,
		0x3560000: PowerShotSX170IS,
		0x3580000: PowerShotSX510HS,
		0x3590000: PowerShotS200,
		0x3600000: IXY620F,
		0x3610000: PowerShotN100,
		0x3640000: PowerShotG1XMarkII,
		0x3650000: PowerShotD30,
		0x3660000: PowerShotSX700HS,
		0x3670000: PowerShotSX600HS,
		0x3680000: PowerShotELPH140IS,
		0x3690000: PowerShotELPH135,
		0x3700000: PowerShotELPH340HS,
		0x3710000: PowerShotELPH150IS,
		0x3740000: EOSM3,
		0x3750000: PowerShotSX60HS,
		0x3760000: PowerShotSX520HS,
		0x3770000: PowerShotSX400IS,
		0x3780000: PowerShotG7X,
		0x3790000: PowerShotN2,
		0x3800000: PowerShotSX530HS,
		0x3820000: PowerShotSX710HS,
		0x3830000: PowerShotSX610HS,
		0x3840000: EOSM10,
		0x3850000: PowerShotG3X,
		0x3860000: PowerShotELPH165HS,
		0x3870000: PowerShotELPH160,
		0x3880000: PowerShotELPH350HS,
		0x3890000: PowerShotELPH170IS,
		0x3910000: PowerShotSX410IS,
		0x3930000: PowerShotG9X,
		0x3940000: EOSM5,
		0x3950000: PowerShotG5X,
		0x3970000: PowerShotG7XMarkII,
		0x3980000: EOSM100,
		0x3990000: PowerShotELPH360HS,
		0x4010000: PowerShotSX540HS,
		0x4020000: PowerShotSX420IS,
		0x4030000: PowerShotELPH190IS,
		0x4040000: PowerShotG1,
		0x4040001: PowerShotELPH180IS,
		0x4050000: PowerShotSX720HS,
		0x4060000: PowerShotSX620HS,
		0x4070000: EOSM6,
		0x4100000: PowerShotG9XMarkII,
		0x412:     EOSM50,
		0x4150000: PowerShotELPH185,
		0x4160000: PowerShotSX430IS,
		0x4170000: PowerShotSX730HS,
		0x4180000: PowerShotG1XMarkIII,
		0x6040000: PowerShotS100,
		0x801:     PowerShotSX740HS,
		0x804:     PowerShotG5XMarkII,
		0x805:     PowerShotSX70HS,
		0x808:     PowerShotG7XMarkIII,
		0x811:     EOSM6MarkII,
		0x812:     EOSM200,

		0x4007d673: DC19,
		0x4007d674: XHA1,
		0x4007d675: HV10,
		0x4007d676: MD130,
		0x4007d777: DC50,
		0x4007d778: HV20,
		0x4007d779: DC211,
		0x4007d77a: HG10,
		0x4007d77b: HR10,
		0x4007d77d: MD255,
		0x4007d81c: HF11,
		0x4007d878: HV30,
		0x4007d87c: XHA1S,
		0x4007d87e: DC301,
		0x4007d87f: FS100,
		0x4007d880: HF10,
		0x4007d882: HG20,
		0x4007d925: HF21,
		0x4007d926: HFS11,
		0x4007d978: HV40,
		0x4007d987: DC410,
		0x4007d988: FS19,
		0x4007d989: HF20,
		0x4007d98a: HFS10,
		0x4007da8e: HFR10,
		0x4007da8f: HFM30,
		0x4007da90: HFS20,
		0x4007da92: FS31,
		0x4007dca0: EOSC300,
		0x4007dda9: HFG25,
		0x4007dfb4: XC10,
		0x4007e1c3: EOSC200,

		0x80000001: EOS1D,
		0x80000167: EOS1DS,
		0x80000168: EOS10D,
		0x80000169: EOS1DMarkIII,
		0x80000170: EOS300D,
		0x80000174: EOS1DMarkII,
		0x80000175: EOS20D,
		0x80000176: EOS450D,
		0x80000188: EOS1DsMarkII,
		0x80000189: EOS350D,
		0x80000190: EOS40D,
		0x80000213: EOS5D,
		0x80000215: EOS1DsMarkIII,
		0x80000218: EOS5DMarkII,
		0x80000219: WFTE1,
		0x80000232: EOS1DMarkIIN,
		0x80000234: EOS30D,
		0x80000236: EOS400D,
		0x80000241: WFTE2,
		0x80000246: WFTE3,
		0x80000250: EOS7D,
		0x80000252: EOS500D,
		0x80000254: EOS1000D,
		0x80000261: EOS50D,
		0x80000269: EOS1DX,
		0x80000270: EOS550D,
		0x80000271: WFTE4,
		0x80000273: WFTE5,
		0x80000281: EOS1DMarkIV,
		0x80000285: EOS5DMarkIII,
		0x80000286: EOS600D,
		0x80000287: EOS60D,
		0x80000288: EOS1100D,
		0x80000289: EOS7DMarkII,
		0x80000297: WFTE2II,
		0x80000298: WFTE4II,
		0x80000301: EOS650D,
		0x80000302: EOS6D,
		0x80000324: EOS1DC,
		0x80000325: EOS70D,
		0x80000326: EOS700D,
		0x80000327: EOS1200D,
		0x80000328: EOS1DXMarkII,
		0x80000331: EOSM,
		0x80000350: EOS80D,
		0x80000355: EOSM2,
		0x80000346: EOS100D,
		0x80000347: EOS760D,
		0x80000349: EOS5DMarkIV,
		0x80000382: EOS5DS,
		0x80000393: EOS750D,
		0x80000401: EOS5DSR,
		0x80000404: EOS1300D,
		0x80000405: EOS800D,
		0x80000406: EOS6DMarkII,
		0x80000408: EOS77D,
		0x80000417: EOS200D,
		0x80000421: EOSR5,
		0x80000422: EOS4000D,
		0x80000424: EOSR,
		0x80000428: EOS1DXMarkIII,
		0x80000432: EOS2000D,
		0x80000433: EOSRP,
		0x80000435: EOS850D,
		0x80000436: EOS250D,
		0x80000437: EOS90D,
		0x80000450: EOSR3,
		0x80000453: EOSR6,
		0x80000464: EOSR7,
		0x80000465: EOSR10,
		0x80000467: PowerShotZOOM,
		0x80000468: EOSM50MarkII,
		0x80000480: EOSR50,
		0x80000481: EOSR6MarkII,
		0x80000487: EOSR8,
		0x80000520: EOSD2000C,
		0x80000560: EOSD6000C,
	}
)

func parseCanonModelID(id uint32) CameraModel {
	prefix := id >> 16
	suffix := id << 16
	switch prefix {
	case 0x4007:
		return CameraModel(suffix)
	case 0x8000:

		return CameraModel(id)

	}
	return UnknownCanonModel
}

var mapStringCanonCameraModel = map[string]CameraModel{
	//DC19/DC21/DC22',
	//XH A1',
	//HV10',
	//MD130/MD140/MD150/MD160/ZR850',
	//DC50', # (iVIS)
	//HV20', # (iVIS)
	//DC211', #29
	//HG10',
	//HR10', #29 (iVIS)
	//MD255/ZR950',
	//HF11',
	//HV30',
	//XH A1S',
	//DC301/DC310/DC311/DC320/DC330',
	//FS100',
	//HF10', #29 (iVIS/VIXIA)
	//HG20/HG21', # (VIXIA)
	//HF21', # (LEGRIA)
	//HF S11', # (LEGRIA)
	//HV40', # (LEGRIA)
	//DC410/DC411/DC420',
	//FS19/FS20/FS21/FS22/FS200', # (LEGRIA)
	//HF20/HF200', # (LEGRIA)
	//HF S10/S100', # (LEGRIA/VIXIA)
	//HF R10/R16/R17/R18/R100/R106', # (LEGRIA/VIXIA)
	//HF M30/M31/M36/M300/M306', # (LEGRIA/VIXIA)
	//HF S20/S21/S200', # (LEGRIA/VIXIA)
	//FS31/FS36/FS37/FS300/FS305/FS306/FS307',
	//EOS C300',
	//HF G25', # (LEGRIA)
	//XC10',
	//EOS C200',

	//EOS-1D',
	//EOS-1DS',
	//EOS 10D',
	//EOS-1D Mark III',
	//EOS Digital Rebel / 300D / Kiss Digital',
	//EOS-1D Mark II',
	//EOS 20D',
	//EOS Digital Rebel XSi / 450D / Kiss X2',
	//EOS-1Ds Mark II',
	//EOS Digital Rebel XT / 350D / Kiss Digital N',
	//EOS 40D',
	//EOS 5D',
	//EOS-1Ds Mark III',
	//EOS 5D Mark II',
	//WFT-E1',
	//EOS-1D Mark II N',
	//EOS 30D',
	//EOS Digital Rebel XTi / 400D / Kiss Digital X',
	//WFT-E2',
	//WFT-E3',
	//EOS 7D',
	//EOS Rebel T1i / 500D / Kiss X3',
	//EOS Rebel XS / 1000D / Kiss F',
	//EOS 50D',
	//EOS-1D X',
	//EOS Rebel T2i / 550D / Kiss X4',
	//WFT-E4',
	//WFT-E5',
	//EOS-1D Mark IV',
	//EOS 5D Mark III',
	//EOS Rebel T3i / 600D / Kiss X5',
	//EOS 60D',
	//EOS Rebel T3 / 1100D / Kiss X50',
	//EOS 7D Mark II', #IB
	//WFT-E2 II',
	//WFT-E4 II',
	//EOS Rebel T4i / 650D / Kiss X6i',
	//EOS 6D', #25
	//EOS-1D C', #(NC)
	//EOS 70D',
	//EOS Rebel T5i / 700D / Kiss X7i',
	//EOS Rebel T5 / 1200D / Kiss X70 / Hi',
	//EOS-1D X Mark II', #42
	//EOS M',
	//EOS 80D', #42
	//EOS M2',
	//EOS Rebel SL1 / 100D / Kiss X7',
	//EOS Rebel T6s / 760D / 8000D',
	//EOS 5D Mark IV', #42
	//EOS 5DS',
	//EOS Rebel T6i / 750D / Kiss X8i',
	//EOS 5DS R',
	//EOS Rebel T6 / 1300D / Kiss X80',
	//EOS Rebel T7i / 800D / Kiss X9i',
	//EOS 6D Mark II', #IB/42
	//EOS 77D / 9000D',
	//EOS Rebel SL2 / 200D / Kiss X9', #IB/42
	//EOS R5', #PH
	//EOS Rebel T100 / 4000D / 3000D', #IB (3000D in China; Kiss? - PH)
	//EOS R', #IB
	//EOS-1D X Mark III', #IB
	//EOS Rebel T7 / 2000D / 1500D / Kiss X90', #IB
	//EOS RP',
	//EOS Rebel T8i / 850D / X10i', #JR/PH
	//EOS SL3 / 250D / Kiss X10', #25
	//EOS 90D', #IB
	//EOS R3', #42
	//EOS R6', #PH
	//EOS R7', #42
	//EOS R10', #42
	//PowerShot ZOOM',
	//EOS M50 Mark II / Kiss M2', #IB
	//EOS R50', #42
	//EOS R6 Mark II', #42
	//EOS R8', #42
	//EOS D2000C', #IB
	//EOS D6000C', #PH (guess)
}

// 0x80000001 => 'EOS-1D',
// 0x80000167 => 'EOS-1DS',
// 0x80000168 => 'EOS 10D',
// 0x80000169 => 'EOS-1D Mark III',
// 0x80000170 => 'EOS Digital Rebel / 300D / Kiss Digital',
// 0x80000174 => 'EOS-1D Mark II',
// 0x80000175 => 'EOS 20D',
// 0x80000176 => 'EOS Digital Rebel XSi / 450D / Kiss X2',
// 0x80000188 => 'EOS-1Ds Mark II',
// 0x80000189 => 'EOS Digital Rebel XT / 350D / Kiss Digital N',
// 0x80000190 => 'EOS 40D',
// 0x80000213 => 'EOS 5D',
// 0x80000215 => 'EOS-1Ds Mark III',
// 0x80000218 => 'EOS 5D Mark II',
// 0x80000219 => 'WFT-E1',
// 0x80000232 => 'EOS-1D Mark II N',
// 0x80000234 => 'EOS 30D',
// 0x80000236 => 'EOS Digital Rebel XTi / 400D / Kiss Digital X',
// 0x80000241 => 'WFT-E2',
// 0x80000246 => 'WFT-E3',
// 0x80000250 => 'EOS 7D',
// 0x80000252 => 'EOS Rebel T1i / 500D / Kiss X3',
// 0x80000254 => 'EOS Rebel XS / 1000D / Kiss F',
// 0x80000261 => 'EOS 50D',
// 0x80000269 => 'EOS-1D X',
// 0x80000270 => 'EOS Rebel T2i / 550D / Kiss X4',
// 0x80000271 => 'WFT-E4',
// 0x80000273 => 'WFT-E5',
// 0x80000281 => 'EOS-1D Mark IV',
// 0x80000285 => 'EOS 5D Mark III',
// 0x80000286 => 'EOS Rebel T3i / 600D / Kiss X5',
// 0x80000287 => 'EOS 60D',
// 0x80000288 => 'EOS Rebel T3 / 1100D / Kiss X50',
// 0x80000289 => 'EOS 7D Mark II', #IB
// 0x80000297 => 'WFT-E2 II',
// 0x80000298 => 'WFT-E4 II',
// 0x80000301 => 'EOS Rebel T4i / 650D / Kiss X6i',
// 0x80000302 => 'EOS 6D', #25
// 0x80000324 => 'EOS-1D C', #(NC)
// 0x80000325 => 'EOS 70D',
// 0x80000326 => 'EOS Rebel T5i / 700D / Kiss X7i',
// 0x80000327 => 'EOS Rebel T5 / 1200D / Kiss X70 / Hi',
// 0x80000328 => 'EOS-1D X Mark II', #42
// 0x80000331 => 'EOS M',
// 0x80000350 => 'EOS 80D', #42
// 0x80000355 => 'EOS M2',
// 0x80000346 => 'EOS Rebel SL1 / 100D / Kiss X7',
// 0x80000347 => 'EOS Rebel T6s / 760D / 8000D',
// 0x80000349 => 'EOS 5D Mark IV', #42
// 0x80000382 => 'EOS 5DS',
// 0x80000393 => 'EOS Rebel T6i / 750D / Kiss X8i',
// 0x80000401 => 'EOS 5DS R',
// 0x80000404 => 'EOS Rebel T6 / 1300D / Kiss X80',
// 0x80000405 => 'EOS Rebel T7i / 800D / Kiss X9i',
// 0x80000406 => 'EOS 6D Mark II', #IB/42
// 0x80000408 => 'EOS 77D / 9000D',
// 0x80000417 => 'EOS Rebel SL2 / 200D / Kiss X9', #IB/42
// 0x80000421 => 'EOS R5', #PH
// 0x80000422 => 'EOS Rebel T100 / 4000D / 3000D', #IB (3000D in China; Kiss? - PH)
// 0x80000424 => 'EOS R', #IB
// 0x80000428 => 'EOS-1D X Mark III', #IB
// 0x80000432 => 'EOS Rebel T7 / 2000D / 1500D / Kiss X90', #IB
// 0x80000433 => 'EOS RP',
// 0x80000435 => 'EOS Rebel T8i / 850D / X10i', #JR/PH
// 0x80000436 => 'EOS SL3 / 250D / Kiss X10', #25
// 0x80000437 => 'EOS 90D', #IB
// 0x80000450 => 'EOS R3', #42
// 0x80000453 => 'EOS R6', #PH
// 0x80000464 => 'EOS R7', #42
// 0x80000465 => 'EOS R10', #42
// 0x80000467 => 'PowerShot ZOOM',
// 0x80000468 => 'EOS M50 Mark II / Kiss M2', #IB
// 0x80000480 => 'EOS R50', #42
// 0x80000481 => 'EOS R6 Mark II', #42
// 0x80000487 => 'EOS R8', #42
// 0x80000520 => 'EOS D2000C', #IB
// 0x80000560 => 'EOS D6000C', #PH (guess)
var mapCanonCameraModel = map[string]uint32{}

//0x1010000 => 'PowerShot A30',
//0x1040000 => 'PowerShot S300 / Digital IXUS 300 / IXY Digital 300',
//0x1060000 => 'PowerShot A20',
//0x1080000 => 'PowerShot A10',
//0x1090000 => 'PowerShot S110 / Digital IXUS v / IXY Digital 200',
//0x1100000 => 'PowerShot G2',
//0x1110000 => 'PowerShot S40',
//0x1120000 => 'PowerShot S30',
//0x1130000 => 'PowerShot A40',
//0x1140000 => 'EOS D30',

// Obtained from https://en.wikipedia.org/wiki/List_of_Canon_products obtained on (02/12/2023)
// EOS 300D/Digital Rebel/Kiss Digital (discontinued)
// EOS 350D/Digital Rebel XT/Kiss Digital N (discontinued)
// EOS 400D/Digital Rebel XTi/Kiss Digital X (discontinued)
// EOS 450D/Rebel XSi/Kiss X2 (discontinued)
// EOS 500D/Rebel T1i/Kiss X3 (discontinued)
// EOS 550D/Rebel T2i/Kiss X4 (discontinued)
// EOS 600D/Rebel T3i/Kiss X5 (discontinued)
// EOS 650D/Rebel T4i/Kiss X6i (discontinued)
// EOS 700D/Rebel T5i/Kiss X7i (discontinued)
// EOS 750D/Rebel T6i/Kiss X8i
// EOS 760D/Rebel T6s/8000D
// EOS 800D/Rebel T7i/Kiss X9i
// EOS 850D/Rebel T8i/Kiss X10i
// EOS 100D/Rebel SL1/Kiss X7 (discontinued)
// EOS 200D/Rebel SL2/Kiss X9
// EOS 250D/Rebel SL3/Kiss X10
// EOS 1000D/Rebel XS/Kiss F (discontinued)
// EOS 1100D/Rebel T3/Kiss X50 (discontinued)
// EOS 1200D/Rebel T5/Kiss X70 (discontinued)
// EOS 1300D/Rebel T6/Kiss X80 (discontinued)
// EOS 1500D/EOS 2000D/Rebel T7/Kiss X90
// EOS 3000D/EOS 4000D/Rebel T100
//
// Canon EOS 77D (EOS 9000D in Japan)
//
// EOS D30 (discontinued)
// EOS D60 (discontinued)
// EOS 10D (discontinued)
// EOS 20D (discontinued)
// EOS 20Da (discontinued) – designed for astrophotography
// EOS 30D (discontinued)
// EOS 40D (discontinued)
// EOS 50D (discontinued)
// EOS 60D (discontinued)
// EOS 60Da – designed for astrophotography
// EOS 70D (discontinued)
// EOS 77D
// EOS 80D
// EOS 90D
//
// APS-C sensor
// EOS 7D (discontinued)
// EOS 7D Mark II
// Full-frame sensor
// EOS 5D (discontinued)
// EOS 5D Mark II (discontinued)
// EOS 5D Mark III (discontinued)
// EOS 5D Mark IV
// EOS 5Ds
// EOS 5Ds R
// EOS 6D (discontinued)
// EOS 6D Mark II
//
// EOS-1D (discontinued)
// EOS-1Ds (discontinued)
// EOS-1D Mark II (discontinued)
// EOS-1Ds Mark II (discontinued)
// EOS-1D Mark II N (discontinued)
// EOS-1D Mark III (discontinued)
// EOS-1Ds Mark III (discontinued)
// EOS-1D Mark IV (discontinued)
// EOS-1D X (discontinued)
// EOS-1D X Mark II (discontinued)
// EOS-1D C (cinema-oriented)
// EOS-1D X Mark III
//
// Canon EOS M
// Canon EOS M2 (not available in North America)
// Canon EOS M3
// Canon EOS M10
// Canon EOS M5
// Canon EOS M6
// Canon EOS M6 Mark II
// Canon EOS M100
// Canon EOS M50
// Canon EOS M50 Mark II
// Canon EOS R7
// Canon EOS R10
//
// Canon EOS R
// Canon EOS RP
// Canon EOS Ra (designed for Astrophotography)
// Canon EOS R5
// Canon EOS R6
// Canon EOS R6 Mark II
// Canon EOS R3
//
// Canon EOS C700 FF
// Canon EOS C700
// Canon EOS C500 Mark II
// Canon EOS C300 Mark III
// Canon EOS C300 Mark II
// Canon EOS C200
// Canon EOS C100 Mark II
// Canon EOS C70
// Canon EOS R5 C - 8K video and 45MP stills; smallest Cinema EOS camera
//
// Canon powerShot S45
// Canon PowerShot S100
// Canon PowerShot S110
// Canon PowerShot S200
// Canon PowerShot S230
// Canon PowerShot S300
// Canon PowerShot S330
// Canon PowerShot S400
// Canon PowerShot S410
// Canon PowerShot S500
// Canon PowerShot SD10
// Canon PowerShot SD20
// Canon PowerShot SD30
// Canon PowerShot SD40
// Canon PowerShot SD100
// Canon PowerShot SD110
// Canon PowerShot SD200
// Canon PowerShot SD300
// Canon PowerShot SD400
// Canon PowerShot SD430
// Canon PowerShot SD450
// Canon PowerShot SD500
// Canon PowerShot SD550
// Canon PowerShot SD600
// Canon PowerShot SD630
// Canon PowerShot SD640 No reference that this camera exist.
// Canon PowerShot SD700 IS
// Canon PowerShot SD750
// Canon PowerShot SD770 IS
// Canon PowerShot SD780 IS
// Canon PowerShot SD790 IS
// Canon PowerShot SD800 IS
// Canon PowerShot SD850 IS
// Canon PowerShot SD870 IS
// Canon PowerShot SD880 IS
// Canon PowerShot SD890 IS
// Canon PowerShot SD900
// Canon PowerShot SD940 IS
// Canon PowerShot SD950 IS
// Canon PowerShot SD960 IS
// Canon PowerShot SD970 IS
// Canon PowerShot SD980 IS
// Canon PowerShot SD990 IS
// Canon PowerShot SD1000
// Canon PowerShot SD1100 IS
// Canon PowerShot SD1200 IS
// Canon PowerShot SD1300 IS
// Canon PowerShot SD1400 IS
// Canon PowerShot SD3500 IS
// Canon PowerShot SD4000 IS
// Canon PowerShot SD4500 IS
// Canon PowerShot 110 HS
// Canon PowerShot 320 HS
// Canon PowerShot 340 HS[4]
// Canon PowerShot 520 HS
//
// Canon PowerShot A5
// Canon PowerShot A5 Zoom
// Canon PowerShot A50
// Canon PowerShot A10
// Canon PowerShot A20
// Canon PowerShot A30
// Canon PowerShot A40
// Canon PowerShot A60
// Canon PowerShot A70
// Canon PowerShot A75
// Canon PowerShot A80
// Canon PowerShot A85
// Canon PowerShot A95
// Canon PowerShot A100
// Canon PowerShot A200
// Canon PowerShot A300
// Canon PowerShot A310
// Canon PowerShot A400
// Canon PowerShot A410
// Canon PowerShot A420
// Canon PowerShot A430
// Canon PowerShot A450
// Canon PowerShot A460
// Canon PowerShot A470
// Canon PowerShot A480
// Canon PowerShot A490 / A495
// Canon PowerShot A510
// Canon PowerShot A520
// Canon PowerShot A530
// Canon PowerShot A540
// Canon PowerShot A550
// Canon PowerShot A560
// Canon PowerShot A570 IS
// Canon PowerShot A580
// Canon PowerShot A590 IS
// Canon PowerShot A610
// Canon PowerShot A620
// Canon PowerShot A630
// Canon PowerShot A640
// Canon PowerShot A650 IS
// Canon PowerShot A700
// Canon PowerShot A710 IS
// Canon PowerShot A720 IS
// Canon PowerShot A800
// Canon PowerShot A810
// Canon PowerShot A1000 IS
// Canon PowerShot A1100 IS
// Canon PowerShot A1200
// Canon PowerShot A1300
// Canon PowerShot A1400
// Canon PowerShot A2000 IS
// Canon PowerShot A2200
// Canon PowerShot A2300 HD
// Canon PowerShot A2500
// Canon PowerShot A2600
// Canon PowerShot A3000 IS
// Canon PowerShot A3100 I
// Canon PowerShot A3150 IS
// Canon PowerShot A3200 IS
// Canon PowerShot A3300 IS
// Canon Powershot A3400 IS
// Canon Powershot A3500 IS
// Canon Powershot A4000 IS
//
// PowerShot D series[edit]
// Canon PowerShot D10
// Canon PowerShot D20
// Canon PowerShot D30
// PowerShot E series[edit]
// Canon PowerShot E1
// PowerShot G series[edit]
// Canon PowerShot G1
// Canon PowerShot G2
// Canon PowerShot G3
// Canon PowerShot G5
// Canon PowerShot G6
// Canon PowerShot G7
// Canon PowerShot G9
// Canon PowerShot G10
// Canon PowerShot G11
// Canon PowerShot G12
// Canon PowerShot G15
// Canon PowerShot G16
// Canon PowerShot G1 X
// Canon PowerShot G1 X MkII
// Canon PowerShot G1 X MkIII
// Canon PowerShot G3 X
// Canon PowerShot G5 X
// Canon PowerShot G5 X MkII
// Canon PowerShot G7 X
// Canon PowerShot G7 X MkII
// Canon PowerShot G7 X MkIII
// Canon PowerShot G9 X
// Canon PowerShot G9 X MkII
//
// PowerShot Pro series[edit]
// Canon Powershot Pro1
// Canon Powershot Pro7d
// Canon Powershot Pro90 IS
// PowerShot S series[edit]
// Canon PowerShot S1 IS
// Canon PowerShot S2 IS
// Canon PowerShot S3 IS
// Canon PowerShot S5 IS
// Canon PowerShot S10
// Canon PowerShot S20
// Canon PowerShot S30
// Canon PowerShot S40
// Canon PowerShot S45
// Canon PowerShot S50
// Canon PowerShot S60
// Canon PowerShot S70
// Canon PowerShot S80
// Canon PowerShot S90
// Canon PowerShot S95
// Canon PowerShot S100
// Canon PowerShot S110
// Canon PowerShot S120
// Canon PowerShot SX1 IS
// Canon PowerShot SX10 IS
// Canon PowerShot SX20 IS
// Canon PowerShot SX30 IS
// Canon PowerShot SX40 HS
// Canon PowerShot SX50 HS
// Canon PowerShot SX60 HS
// Canon PowerShot SX70 HS
// Canon PowerShot SX100 IS
// Canon PowerShot SX110 IS
// Canon PowerShot SX120 IS
// Canon PowerShot SX130 IS
// Canon PowerShot SX150 IS
// Canon PowerShot SX160 IS
// Canon PowerShot SX200 IS
// Canon PowerShot SX210 IS
// Canon PowerShot SX220 HS
// Canon PowerShot SX230 HS (features GPS)
// Canon PowerShot SX240 HS
// Canon PowerShot SX260 HS (features GPS)
// Canon PowerShot SX270 HS
// Canon PowerShot SX280 HS (features GPS)
// Canon PowerShot SX400 IS
// Canon PowerShot SX410 IS
// Canon PowerShot SX420 IS
// (first PowerShot camera with built-in Wi-Fi)
//
// Canon PowerShot SX430 IS
// (not officially sold in North America)
//
// Canon PowerShot SX500 IS
// Canon PowerShot SX510 HS
// Canon PowerShot SX520 HS
// Canon PowerShot SX530 HS
// Canon PowerShot SX540 HS
// Canon PowerShot SX600 HS
// (first SX-Series based PowerShot camera to be more compact)
//
// Canon PowerShot SX610 HS
// Canon PowerShot SX620 HS
// Canon PowerShot SX700 HS
// Canon PowerShot SX710 HS
// Canon PowerShot SX720 HS
// Canon PowerShot SX730 HS
// (first Powershot camera with a flip screen for selfies and vlogs)
//
// Canon PowerShot SX740 HS (features 4K recording)
// PowerShot T series[edit]
// Canon PowerShot TX1
//
