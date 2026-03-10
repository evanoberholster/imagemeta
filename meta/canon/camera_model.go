package canon

import "strconv"

// CanonCameraModelSourceFile is the ExifTool source used to generate this table. "Image/ExifTool/Canon.pm"
// CanonCameraModelSourceVersion is the ExifTool version used to generate this table. "13.50"

type CanonCameraModel uint32

const (
	CanonModelUnknown             CanonCameraModel = 0
	CanonModelPowerShotA30        CanonCameraModel = 0x1010000
	CanonModelPowerShotS300       CanonCameraModel = 0x1040000
	CanonModelPowerShotA20        CanonCameraModel = 0x1060000
	CanonModelPowerShotA10        CanonCameraModel = 0x1080000
	CanonModelPowerShotS110       CanonCameraModel = 0x1090000
	CanonModelPowerShotG2         CanonCameraModel = 0x1100000
	CanonModelPowerShotS40        CanonCameraModel = 0x1110000
	CanonModelPowerShotS30        CanonCameraModel = 0x1120000
	CanonModelPowerShotA40        CanonCameraModel = 0x1130000
	CanonModelEOSD30              CanonCameraModel = 0x1140000
	CanonModelPowerShotA100       CanonCameraModel = 0x1150000
	CanonModelPowerShotS200       CanonCameraModel = 0x1160000
	CanonModelPowerShotA200       CanonCameraModel = 0x1170000
	CanonModelPowerShotS330       CanonCameraModel = 0x1180000
	CanonModelPowerShotG3         CanonCameraModel = 0x1190000
	CanonModelPowerShotS45        CanonCameraModel = 0x1210000
	CanonModelPowerShotSD100      CanonCameraModel = 0x1230000
	CanonModelPowerShotS230       CanonCameraModel = 0x1240000
	CanonModelPowerShotA70        CanonCameraModel = 0x1250000
	CanonModelPowerShotA60        CanonCameraModel = 0x1260000
	CanonModelPowerShotS400       CanonCameraModel = 0x1270000
	CanonModelPowerShotG5         CanonCameraModel = 0x1290000
	CanonModelPowerShotA300       CanonCameraModel = 0x1300000
	CanonModelPowerShotS50        CanonCameraModel = 0x1310000
	CanonModelPowerShotA80        CanonCameraModel = 0x1340000
	CanonModelPowerShotSD10       CanonCameraModel = 0x1350000
	CanonModelPowerShotS1IS       CanonCameraModel = 0x1360000
	CanonModelPowerShotPro1       CanonCameraModel = 0x1370000
	CanonModelPowerShotS70        CanonCameraModel = 0x1380000
	CanonModelPowerShotS60        CanonCameraModel = 0x1390000
	CanonModelPowerShotG6         CanonCameraModel = 0x1400000
	CanonModelPowerShotS500       CanonCameraModel = 0x1410000
	CanonModelPowerShotA75        CanonCameraModel = 0x1420000
	CanonModelPowerShotSD110      CanonCameraModel = 0x1440000
	CanonModelPowerShotA400       CanonCameraModel = 0x1450000
	CanonModelPowerShotA310       CanonCameraModel = 0x1470000
	CanonModelPowerShotA85        CanonCameraModel = 0x1490000
	CanonModelPowerShotS410       CanonCameraModel = 0x1520000
	CanonModelPowerShotA95        CanonCameraModel = 0x1530000
	CanonModelPowerShotSD300      CanonCameraModel = 0x1540000
	CanonModelPowerShotSD200      CanonCameraModel = 0x1550000
	CanonModelPowerShotA520       CanonCameraModel = 0x1560000
	CanonModelPowerShotA510       CanonCameraModel = 0x1570000
	CanonModelPowerShotSD20       CanonCameraModel = 0x1590000
	CanonModelPowerShotS2IS       CanonCameraModel = 0x1640000
	CanonModelPowerShotSD430      CanonCameraModel = 0x1650000
	CanonModelPowerShotSD500      CanonCameraModel = 0x1660000
	CanonModelEOSD60              CanonCameraModel = 0x1668000
	CanonModelPowerShotSD30       CanonCameraModel = 0x1700000
	CanonModelPowerShotA430       CanonCameraModel = 0x1740000
	CanonModelPowerShotA410       CanonCameraModel = 0x1750000
	CanonModelPowerShotS80        CanonCameraModel = 0x1760000
	CanonModelPowerShotA620       CanonCameraModel = 0x1780000
	CanonModelPowerShotA610       CanonCameraModel = 0x1790000
	CanonModelPowerShotSD630      CanonCameraModel = 0x1800000
	CanonModelPowerShotSD450      CanonCameraModel = 0x1810000
	CanonModelPowerShotTX1        CanonCameraModel = 0x1820000
	CanonModelPowerShotSD400      CanonCameraModel = 0x1870000
	CanonModelPowerShotA420       CanonCameraModel = 0x1880000
	CanonModelPowerShotSD900      CanonCameraModel = 0x1890000
	CanonModelPowerShotSD550      CanonCameraModel = 0x1900000
	CanonModelPowerShotA700       CanonCameraModel = 0x1920000
	CanonModelPowerShotSD700IS    CanonCameraModel = 0x1940000
	CanonModelPowerShotS3IS       CanonCameraModel = 0x1950000
	CanonModelPowerShotA540       CanonCameraModel = 0x1960000
	CanonModelPowerShotSD600      CanonCameraModel = 0x1970000
	CanonModelPowerShotG7         CanonCameraModel = 0x1980000
	CanonModelPowerShotA530       CanonCameraModel = 0x1990000
	CanonModelPowerShotSD800IS    CanonCameraModel = 0x2000000
	CanonModelPowerShotSD40       CanonCameraModel = 0x2010000
	CanonModelPowerShotA710IS     CanonCameraModel = 0x2020000
	CanonModelPowerShotA640       CanonCameraModel = 0x2030000
	CanonModelPowerShotA630       CanonCameraModel = 0x2040000
	CanonModelPowerShotS5IS       CanonCameraModel = 0x2090000
	CanonModelPowerShotA460       CanonCameraModel = 0x2100000
	CanonModelPowerShotSD850IS    CanonCameraModel = 0x2120000
	CanonModelPowerShotA570IS     CanonCameraModel = 0x2130000
	CanonModelPowerShotA560       CanonCameraModel = 0x2140000
	CanonModelPowerShotSD750      CanonCameraModel = 0x2150000
	CanonModelPowerShotSD1000     CanonCameraModel = 0x2160000
	CanonModelPowerShotA550       CanonCameraModel = 0x2180000
	CanonModelPowerShotA450       CanonCameraModel = 0x2190000
	CanonModelPowerShotG9         CanonCameraModel = 0x2230000
	CanonModelPowerShotA650IS     CanonCameraModel = 0x2240000
	CanonModelPowerShotA720IS     CanonCameraModel = 0x2260000
	CanonModelPowerShotSX100IS    CanonCameraModel = 0x2290000
	CanonModelPowerShotSD950IS    CanonCameraModel = 0x2300000
	CanonModelPowerShotSD870IS    CanonCameraModel = 0x2310000
	CanonModelPowerShotSD890IS    CanonCameraModel = 0x2320000
	CanonModelPowerShotSD790IS    CanonCameraModel = 0x2360000
	CanonModelPowerShotSD770IS    CanonCameraModel = 0x2370000
	CanonModelPowerShotA590IS     CanonCameraModel = 0x2380000
	CanonModelPowerShotA580       CanonCameraModel = 0x2390000
	CanonModelPowerShotA470       CanonCameraModel = 0x2420000
	CanonModelPowerShotSD1100IS   CanonCameraModel = 0x2430000
	CanonModelPowerShotSX1IS      CanonCameraModel = 0x2460000
	CanonModelPowerShotSX10IS     CanonCameraModel = 0x2470000
	CanonModelPowerShotA1000IS    CanonCameraModel = 0x2480000
	CanonModelPowerShotG10        CanonCameraModel = 0x2490000
	CanonModelPowerShotA2000IS    CanonCameraModel = 0x2510000
	CanonModelPowerShotSX110IS    CanonCameraModel = 0x2520000
	CanonModelPowerShotSD990IS    CanonCameraModel = 0x2530000
	CanonModelPowerShotSD880IS    CanonCameraModel = 0x2540000
	CanonModelPowerShotE1         CanonCameraModel = 0x2550000
	CanonModelPowerShotD10        CanonCameraModel = 0x2560000
	CanonModelPowerShotSD960IS    CanonCameraModel = 0x2570000
	CanonModelPowerShotA2100IS    CanonCameraModel = 0x2580000
	CanonModelPowerShotA480       CanonCameraModel = 0x2590000
	CanonModelPowerShotSX200IS    CanonCameraModel = 0x2600000
	CanonModelPowerShotSD970IS    CanonCameraModel = 0x2610000
	CanonModelPowerShotSD780IS    CanonCameraModel = 0x2620000
	CanonModelPowerShotA1100IS    CanonCameraModel = 0x2630000
	CanonModelPowerShotSD1200IS   CanonCameraModel = 0x2640000
	CanonModelPowerShotG11        CanonCameraModel = 0x2700000
	CanonModelPowerShotSX120IS    CanonCameraModel = 0x2710000
	CanonModelPowerShotS90        CanonCameraModel = 0x2720000
	CanonModelPowerShotSX20IS     CanonCameraModel = 0x2750000
	CanonModelPowerShotSD980IS    CanonCameraModel = 0x2760000
	CanonModelPowerShotSD940IS    CanonCameraModel = 0x2770000
	CanonModelPowerShotA495       CanonCameraModel = 0x2800000
	CanonModelPowerShotA490       CanonCameraModel = 0x2810000
	CanonModelPowerShotA3100      CanonCameraModel = 0x2820000
	CanonModelPowerShotA3000IS    CanonCameraModel = 0x2830000
	CanonModelPowerShotSD1400IS   CanonCameraModel = 0x2840000
	CanonModelPowerShotSD1300IS   CanonCameraModel = 0x2850000
	CanonModelPowerShotSD3500IS   CanonCameraModel = 0x2860000
	CanonModelPowerShotSX210IS    CanonCameraModel = 0x2870000
	CanonModelPowerShotSD4000IS   CanonCameraModel = 0x2880000
	CanonModelPowerShotSD4500IS   CanonCameraModel = 0x2890000
	CanonModelPowerShotG12        CanonCameraModel = 0x2920000
	CanonModelPowerShotSX30IS     CanonCameraModel = 0x2930000
	CanonModelPowerShotSX130IS    CanonCameraModel = 0x2940000
	CanonModelPowerShotS95        CanonCameraModel = 0x2950000
	CanonModelPowerShotA3300IS    CanonCameraModel = 0x2980000
	CanonModelPowerShotA3200IS    CanonCameraModel = 0x2990000
	CanonModelPowerShotELPH500HS  CanonCameraModel = 0x3000000
	CanonModelPowerShotPro90IS    CanonCameraModel = 0x3010000
	CanonModelPowerShotA800       CanonCameraModel = 0x3010001
	CanonModelPowerShotELPH100HS  CanonCameraModel = 0x3020000
	CanonModelPowerShotSX230HS    CanonCameraModel = 0x3030000
	CanonModelPowerShotELPH300HS  CanonCameraModel = 0x3040000
	CanonModelPowerShotA2200      CanonCameraModel = 0x3050000
	CanonModelPowerShotA1200      CanonCameraModel = 0x3060000
	CanonModelPowerShotSX220HS    CanonCameraModel = 0x3070000
	CanonModelPowerShotG1X        CanonCameraModel = 0x3080000
	CanonModelPowerShotSX150IS    CanonCameraModel = 0x3090000
	CanonModelPowerShotELPH510HS  CanonCameraModel = 0x3100000
	CanonModelPowerShotS100new    CanonCameraModel = 0x3110000
	CanonModelPowerShotSX40HS     CanonCameraModel = 0x3130000
	CanonModelPowerShotELPH310HS  CanonCameraModel = 0x3120000
	CanonModelIXY32S              CanonCameraModel = 0x3140000
	CanonModelPowerShotA1300      CanonCameraModel = 0x3160000
	CanonModelPowerShotA810       CanonCameraModel = 0x3170000
	CanonModelPowerShotELPH320HS  CanonCameraModel = 0x3180000
	CanonModelPowerShotELPH110HS  CanonCameraModel = 0x3190000
	CanonModelPowerShotD20        CanonCameraModel = 0x3200000
	CanonModelPowerShotA4000IS    CanonCameraModel = 0x3210000
	CanonModelPowerShotSX260HS    CanonCameraModel = 0x3220000
	CanonModelPowerShotSX240HS    CanonCameraModel = 0x3230000
	CanonModelPowerShotELPH530HS  CanonCameraModel = 0x3240000
	CanonModelPowerShotELPH520HS  CanonCameraModel = 0x3250000
	CanonModelPowerShotA3400IS    CanonCameraModel = 0x3260000
	CanonModelPowerShotA2400IS    CanonCameraModel = 0x3270000
	CanonModelPowerShotA2300      CanonCameraModel = 0x3280000
	CanonModelPowerShotS100V      CanonCameraModel = 0x3320000
	CanonModelPowerShotG15        CanonCameraModel = 0x3330000
	CanonModelPowerShotSX50HS     CanonCameraModel = 0x3340000
	CanonModelPowerShotSX160IS    CanonCameraModel = 0x3350000
	CanonModelPowerShotS110new    CanonCameraModel = 0x3360000
	CanonModelPowerShotSX500IS    CanonCameraModel = 0x3370000
	CanonModelPowerShotN          CanonCameraModel = 0x3380000
	CanonModelIXUS245HS           CanonCameraModel = 0x3390000
	CanonModelPowerShotSX280HS    CanonCameraModel = 0x3400000
	CanonModelPowerShotSX270HS    CanonCameraModel = 0x3410000
	CanonModelPowerShotA3500IS    CanonCameraModel = 0x3420000
	CanonModelPowerShotA2600      CanonCameraModel = 0x3430000
	CanonModelPowerShotSX275HS    CanonCameraModel = 0x3440000
	CanonModelPowerShotA1400      CanonCameraModel = 0x3450000
	CanonModelPowerShotELPH130IS  CanonCameraModel = 0x3460000
	CanonModelPowerShotELPH115    CanonCameraModel = 0x3470000
	CanonModelPowerShotELPH330HS  CanonCameraModel = 0x3490000
	CanonModelPowerShotA2500      CanonCameraModel = 0x3510000
	CanonModelPowerShotG16        CanonCameraModel = 0x3540000
	CanonModelPowerShotS120       CanonCameraModel = 0x3550000
	CanonModelPowerShotSX170IS    CanonCameraModel = 0x3560000
	CanonModelPowerShotSX510HS    CanonCameraModel = 0x3580000
	CanonModelPowerShotS200new    CanonCameraModel = 0x3590000
	CanonModelIXY620F             CanonCameraModel = 0x3600000
	CanonModelPowerShotN100       CanonCameraModel = 0x3610000
	CanonModelPowerShotG1XMarkII  CanonCameraModel = 0x3640000
	CanonModelPowerShotD30        CanonCameraModel = 0x3650000
	CanonModelPowerShotSX700HS    CanonCameraModel = 0x3660000
	CanonModelPowerShotSX600HS    CanonCameraModel = 0x3670000
	CanonModelPowerShotELPH140IS  CanonCameraModel = 0x3680000
	CanonModelPowerShotELPH135    CanonCameraModel = 0x3690000
	CanonModelPowerShotELPH340HS  CanonCameraModel = 0x3700000
	CanonModelPowerShotELPH150IS  CanonCameraModel = 0x3710000
	CanonModelEOSM3               CanonCameraModel = 0x3740000
	CanonModelPowerShotSX60HS     CanonCameraModel = 0x3750000
	CanonModelPowerShotSX520HS    CanonCameraModel = 0x3760000
	CanonModelPowerShotSX400IS    CanonCameraModel = 0x3770000
	CanonModelPowerShotG7X        CanonCameraModel = 0x3780000
	CanonModelPowerShotN2         CanonCameraModel = 0x3790000
	CanonModelPowerShotSX530HS    CanonCameraModel = 0x3800000
	CanonModelPowerShotSX710HS    CanonCameraModel = 0x3820000
	CanonModelPowerShotSX610HS    CanonCameraModel = 0x3830000
	CanonModelEOSM10              CanonCameraModel = 0x3840000
	CanonModelPowerShotG3X        CanonCameraModel = 0x3850000
	CanonModelPowerShotELPH165HS  CanonCameraModel = 0x3860000
	CanonModelPowerShotELPH160    CanonCameraModel = 0x3870000
	CanonModelPowerShotELPH350HS  CanonCameraModel = 0x3880000
	CanonModelPowerShotELPH170IS  CanonCameraModel = 0x3890000
	CanonModelPowerShotSX410IS    CanonCameraModel = 0x3910000
	CanonModelPowerShotG9X        CanonCameraModel = 0x3930000
	CanonModelEOSM5               CanonCameraModel = 0x3940000
	CanonModelPowerShotG5X        CanonCameraModel = 0x3950000
	CanonModelPowerShotG7XMarkII  CanonCameraModel = 0x3970000
	CanonModelEOSM100             CanonCameraModel = 0x3980000
	CanonModelPowerShotELPH360HS  CanonCameraModel = 0x3990000
	CanonModelPowerShotSX540HS    CanonCameraModel = 0x4010000
	CanonModelPowerShotSX420IS    CanonCameraModel = 0x4020000
	CanonModelPowerShotELPH190IS  CanonCameraModel = 0x4030000
	CanonModelPowerShotG1         CanonCameraModel = 0x4040000
	CanonModelPowerShotELPH180IS  CanonCameraModel = 0x4040001
	CanonModelPowerShotSX720HS    CanonCameraModel = 0x4050000
	CanonModelPowerShotSX620HS    CanonCameraModel = 0x4060000
	CanonModelEOSM6               CanonCameraModel = 0x4070000
	CanonModelPowerShotG9XMarkII  CanonCameraModel = 0x4100000
	CanonModelEOSM50              CanonCameraModel = 0x412
	CanonModelPowerShotELPH185    CanonCameraModel = 0x4150000
	CanonModelPowerShotSX430IS    CanonCameraModel = 0x4160000
	CanonModelPowerShotSX730HS    CanonCameraModel = 0x4170000
	CanonModelPowerShotG1XMarkIII CanonCameraModel = 0x4180000
	CanonModelPowerShotS100       CanonCameraModel = 0x6040000
	CanonModelPowerShotSX740HS    CanonCameraModel = 0x801
	CanonModelPowerShotG5XMarkII  CanonCameraModel = 0x804
	CanonModelPowerShotSX70HS     CanonCameraModel = 0x805
	CanonModelPowerShotG7XMarkIII CanonCameraModel = 0x808
	CanonModelEOSM6MarkII         CanonCameraModel = 0x811
	CanonModelEOSM200             CanonCameraModel = 0x812
	CanonModelEOSC50              CanonCameraModel = 0x40000227
	CanonModelDC19                CanonCameraModel = 0x4007d673
	CanonModelXHA1                CanonCameraModel = 0x4007d674
	CanonModelHV10                CanonCameraModel = 0x4007d675
	CanonModelMD130               CanonCameraModel = 0x4007d676
	CanonModelDC50                CanonCameraModel = 0x4007d777
	CanonModelHV20                CanonCameraModel = 0x4007d778
	CanonModelDC211               CanonCameraModel = 0x4007d779
	CanonModelHG10                CanonCameraModel = 0x4007d77a
	CanonModelHR10                CanonCameraModel = 0x4007d77b
	CanonModelMD255               CanonCameraModel = 0x4007d77d
	CanonModelHF11                CanonCameraModel = 0x4007d81c
	CanonModelHV30                CanonCameraModel = 0x4007d878
	CanonModelXHA1S               CanonCameraModel = 0x4007d87c
	CanonModelDC301               CanonCameraModel = 0x4007d87e
	CanonModelFS100               CanonCameraModel = 0x4007d87f
	CanonModelHF10                CanonCameraModel = 0x4007d880
	CanonModelHG20                CanonCameraModel = 0x4007d882
	CanonModelHF21                CanonCameraModel = 0x4007d925
	CanonModelHFS11               CanonCameraModel = 0x4007d926
	CanonModelHV40                CanonCameraModel = 0x4007d978
	CanonModelDC410               CanonCameraModel = 0x4007d987
	CanonModelFS19                CanonCameraModel = 0x4007d988
	CanonModelHF20                CanonCameraModel = 0x4007d989
	CanonModelHFS10               CanonCameraModel = 0x4007d98a
	CanonModelHFR10               CanonCameraModel = 0x4007da8e
	CanonModelHFM30               CanonCameraModel = 0x4007da8f
	CanonModelHFS20               CanonCameraModel = 0x4007da90
	CanonModelFS31                CanonCameraModel = 0x4007da92
	CanonModelEOSC300             CanonCameraModel = 0x4007dca0
	CanonModelHFG25               CanonCameraModel = 0x4007dda9
	CanonModelXC10                CanonCameraModel = 0x4007dfb4
	CanonModelEOSC200             CanonCameraModel = 0x4007e1c3
	CanonModelEOS1D               CanonCameraModel = 0x80000001
	CanonModelEOS1DS              CanonCameraModel = 0x80000167
	CanonModelEOS10D              CanonCameraModel = 0x80000168
	CanonModelEOS1DMarkIII        CanonCameraModel = 0x80000169
	CanonModelEOSDigitalRebel     CanonCameraModel = 0x80000170
	CanonModelEOS1DMarkII         CanonCameraModel = 0x80000174
	CanonModelEOS20D              CanonCameraModel = 0x80000175
	CanonModelEOSDigitalRebelXSi  CanonCameraModel = 0x80000176
	CanonModelEOS1DsMarkII        CanonCameraModel = 0x80000188
	CanonModelEOSDigitalRebelXT   CanonCameraModel = 0x80000189
	CanonModelEOS40D              CanonCameraModel = 0x80000190
	CanonModelEOS5D               CanonCameraModel = 0x80000213
	CanonModelEOS1DsMarkIII       CanonCameraModel = 0x80000215
	CanonModelEOS5DMarkII         CanonCameraModel = 0x80000218
	CanonModelWFTE1               CanonCameraModel = 0x80000219
	CanonModelEOS1DMarkIIN        CanonCameraModel = 0x80000232
	CanonModelEOS30D              CanonCameraModel = 0x80000234
	CanonModelEOSDigitalRebelXTi  CanonCameraModel = 0x80000236
	CanonModelWFTE2               CanonCameraModel = 0x80000241
	CanonModelWFTE3               CanonCameraModel = 0x80000246
	CanonModelEOS7D               CanonCameraModel = 0x80000250
	CanonModelEOSRebelT1i         CanonCameraModel = 0x80000252
	CanonModelEOSRebelXS          CanonCameraModel = 0x80000254
	CanonModelEOS50D              CanonCameraModel = 0x80000261
	CanonModelEOS1DX              CanonCameraModel = 0x80000269
	CanonModelEOSRebelT2i         CanonCameraModel = 0x80000270
	CanonModelWFTE4               CanonCameraModel = 0x80000271
	CanonModelWFTE5               CanonCameraModel = 0x80000273
	CanonModelEOS1DMarkIV         CanonCameraModel = 0x80000281
	CanonModelEOS5DMarkIII        CanonCameraModel = 0x80000285
	CanonModelEOSRebelT3i         CanonCameraModel = 0x80000286
	CanonModelEOS60D              CanonCameraModel = 0x80000287
	CanonModelEOSRebelT3          CanonCameraModel = 0x80000288
	CanonModelEOS7DMarkII         CanonCameraModel = 0x80000289
	CanonModelWFTE2II             CanonCameraModel = 0x80000297
	CanonModelWFTE4II             CanonCameraModel = 0x80000298
	CanonModelEOSRebelT4i         CanonCameraModel = 0x80000301
	CanonModelEOS6D               CanonCameraModel = 0x80000302
	CanonModelEOS1DC              CanonCameraModel = 0x80000324
	CanonModelEOS70D              CanonCameraModel = 0x80000325
	CanonModelEOSRebelT5i         CanonCameraModel = 0x80000326
	CanonModelEOSRebelT5          CanonCameraModel = 0x80000327
	CanonModelEOS1DXMarkII        CanonCameraModel = 0x80000328
	CanonModelEOSM                CanonCameraModel = 0x80000331
	CanonModelEOS80D              CanonCameraModel = 0x80000350
	CanonModelEOSM2               CanonCameraModel = 0x80000355
	CanonModelEOSRebelSL1         CanonCameraModel = 0x80000346
	CanonModelEOSRebelT6s         CanonCameraModel = 0x80000347
	CanonModelEOS5DMarkIV         CanonCameraModel = 0x80000349
	CanonModelEOS5DS              CanonCameraModel = 0x80000382
	CanonModelEOSRebelT6i         CanonCameraModel = 0x80000393
	CanonModelEOS5DSR             CanonCameraModel = 0x80000401
	CanonModelEOSRebelT6          CanonCameraModel = 0x80000404
	CanonModelEOSRebelT7i         CanonCameraModel = 0x80000405
	CanonModelEOS6DMarkII         CanonCameraModel = 0x80000406
	CanonModelEOS77D              CanonCameraModel = 0x80000408
	CanonModelEOSRebelSL2         CanonCameraModel = 0x80000417
	CanonModelEOSR5               CanonCameraModel = 0x80000421
	CanonModelEOSRebelT100        CanonCameraModel = 0x80000422
	CanonModelEOSR                CanonCameraModel = 0x80000424
	CanonModelEOS1DXMarkIII       CanonCameraModel = 0x80000428
	CanonModelEOSRebelT7          CanonCameraModel = 0x80000432
	CanonModelEOSRP               CanonCameraModel = 0x80000433
	CanonModelEOSRebelT8i         CanonCameraModel = 0x80000435
	CanonModelEOSSL3              CanonCameraModel = 0x80000436
	CanonModelEOS90D              CanonCameraModel = 0x80000437
	CanonModelEOSR3               CanonCameraModel = 0x80000450
	CanonModelEOSR6               CanonCameraModel = 0x80000453
	CanonModelEOSR7               CanonCameraModel = 0x80000464
	CanonModelEOSR10              CanonCameraModel = 0x80000465
	CanonModelPowerShotZOOM       CanonCameraModel = 0x80000467
	CanonModelEOSM50MarkII        CanonCameraModel = 0x80000468
	CanonModelEOSR50              CanonCameraModel = 0x80000480
	CanonModelEOSR6MarkII         CanonCameraModel = 0x80000481
	CanonModelEOSR8               CanonCameraModel = 0x80000487
	CanonModelPowerShotV10        CanonCameraModel = 0x80000491
	CanonModelEOSR1               CanonCameraModel = 0x80000495
	CanonModelEOSR5MarkII         CanonCameraModel = 0x80000496
	CanonModelPowerShotV1         CanonCameraModel = 0x80000497
	CanonModelEOSR100             CanonCameraModel = 0x80000498
	CanonModelEOSR50V             CanonCameraModel = 0x80000516
	CanonModelEOSR6MarkIII        CanonCameraModel = 0x80000518
	CanonModelEOSD2000C           CanonCameraModel = 0x80000520
	CanonModelEOSD6000C           CanonCameraModel = 0x80000560
)

// CanonCameraModelCount is the number of Canon model entries from ExifTool %canonModelID.
const CanonCameraModelCount = 357

const _CanonCameraModel_name = "UnknownPowerShot A30PowerShot S300PowerShot A20PowerShot A10PowerShot S110PowerShot G2PowerShot S40PowerShot S30PowerShot A40EOS D30PowerShot A100PowerShot S200PowerShot A200PowerShot S330PowerShot G3PowerShot S45PowerShot SD100PowerShot S230PowerShot A70PowerShot A60PowerShot S400PowerShot G5PowerShot A300PowerShot S50PowerShot A80PowerShot SD10PowerShot S1 ISPowerShot Pro1PowerShot S70PowerShot S60PowerShot G6PowerShot S500PowerShot A75PowerShot SD110PowerShot A400PowerShot A310PowerShot A85PowerShot S410PowerShot A95PowerShot SD300PowerShot SD200PowerShot A520PowerShot A510PowerShot SD20PowerShot S2 ISPowerShot SD430PowerShot SD500EOS D60PowerShot SD30PowerShot A430PowerShot A410PowerShot S80PowerShot A620PowerShot A610PowerShot SD630PowerShot SD450PowerShot TX1PowerShot SD400PowerShot A420PowerShot SD900PowerShot SD550PowerShot A700PowerShot SD700 ISPowerShot S3 ISPowerShot A540PowerShot SD600PowerShot G7PowerShot A530PowerShot SD800 ISPowerShot SD40PowerShot A710 ISPowerShot A640PowerShot A630PowerShot S5 ISPowerShot A460PowerShot SD850 ISPowerShot A570 ISPowerShot A560PowerShot SD750PowerShot SD1000PowerShot A550PowerShot A450PowerShot G9PowerShot A650 ISPowerShot A720 ISPowerShot SX100 ISPowerShot SD950 ISPowerShot SD870 ISPowerShot SD890 ISPowerShot SD790 ISPowerShot SD770 ISPowerShot A590 ISPowerShot A580PowerShot A470PowerShot SD1100 ISPowerShot SX1 ISPowerShot SX10 ISPowerShot A1000 ISPowerShot G10PowerShot A2000 ISPowerShot SX110 ISPowerShot SD990 ISPowerShot SD880 ISPowerShot E1PowerShot D10PowerShot SD960 ISPowerShot A2100 ISPowerShot A480PowerShot SX200 ISPowerShot SD970 ISPowerShot SD780 ISPowerShot A1100 ISPowerShot SD1200 ISPowerShot G11PowerShot SX120 ISPowerShot S90PowerShot SX20 ISPowerShot SD980 ISPowerShot SD940 ISPowerShot A495PowerShot A490PowerShot A3100PowerShot A3000 ISPowerShot SD1400 ISPowerShot SD1300 ISPowerShot SD3500 ISPowerShot SX210 ISPowerShot SD4000 ISPowerShot SD4500 ISPowerShot G12PowerShot SX30 ISPowerShot SX130 ISPowerShot S95PowerShot A3300 ISPowerShot A3200 ISPowerShot ELPH 500 HSPowerShot Pro90 ISPowerShot A800PowerShot ELPH 100 HSPowerShot SX230 HSPowerShot ELPH 300 HSPowerShot A2200PowerShot A1200PowerShot SX220 HSPowerShot G1 XPowerShot SX150 ISPowerShot ELPH 510 HSPowerShot S100 (new)PowerShot SX40 HSPowerShot ELPH 310 HSIXY 32SPowerShot A1300PowerShot A810PowerShot ELPH 320 HSPowerShot ELPH 110 HSPowerShot D20PowerShot A4000 ISPowerShot SX260 HSPowerShot SX240 HSPowerShot ELPH 530 HSPowerShot ELPH 520 HSPowerShot A3400 ISPowerShot A2400 ISPowerShot A2300PowerShot S100VPowerShot G15PowerShot SX50 HSPowerShot SX160 ISPowerShot S110 (new)PowerShot SX500 ISPowerShot NIXUS 245 HSPowerShot SX280 HSPowerShot SX270 HSPowerShot A3500 ISPowerShot A2600PowerShot SX275 HSPowerShot A1400PowerShot ELPH 130 ISPowerShot ELPH 115PowerShot ELPH 330 HSPowerShot A2500PowerShot G16PowerShot S120PowerShot SX170 ISPowerShot SX510 HSPowerShot S200 (new)IXY 620FPowerShot N100PowerShot G1 X Mark IIPowerShot D30PowerShot SX700 HSPowerShot SX600 HSPowerShot ELPH 140 ISPowerShot ELPH 135PowerShot ELPH 340 HSPowerShot ELPH 150 ISEOS M3PowerShot SX60 HSPowerShot SX520 HSPowerShot SX400 ISPowerShot G7 XPowerShot N2PowerShot SX530 HSPowerShot SX710 HSPowerShot SX610 HSEOS M10PowerShot G3 XPowerShot ELPH 165 HSPowerShot ELPH 160PowerShot ELPH 350 HSPowerShot ELPH 170 ISPowerShot SX410 ISPowerShot G9 XEOS M5PowerShot G5 XPowerShot G7 X Mark IIEOS M100PowerShot ELPH 360 HSPowerShot SX540 HSPowerShot SX420 ISPowerShot ELPH 190 ISPowerShot G1PowerShot ELPH 180 ISPowerShot SX720 HSPowerShot SX620 HSEOS M6PowerShot G9 X Mark IIEOS M50PowerShot ELPH 185PowerShot SX430 ISPowerShot SX730 HSPowerShot G1 X Mark IIIPowerShot S100PowerShot SX740 HSPowerShot G5 X Mark IIPowerShot SX70 HSPowerShot G7 X Mark IIIEOS M6 Mark IIEOS M200EOS C50DC19XH A1HV10MD130DC50HV20DC211HG10HR10MD255HF11HV30XH A1SDC301FS100HF10HG20HF21HF S11HV40DC410FS19HF20HF S10HF R10HF M30HF S20FS31EOS C300HF G25XC10EOS C200EOS-1DEOS-1DSEOS 10DEOS-1D Mark IIIEOS Digital RebelEOS-1D Mark IIEOS 20DEOS Digital Rebel XSiEOS-1Ds Mark IIEOS Digital Rebel XTEOS 40DEOS 5DEOS-1Ds Mark IIIEOS 5D Mark IIWFT-E1EOS-1D Mark II NEOS 30DEOS Digital Rebel XTiWFT-E2WFT-E3EOS 7DEOS Rebel T1iEOS Rebel XSEOS 50DEOS-1D XEOS Rebel T2iWFT-E4WFT-E5EOS-1D Mark IVEOS 5D Mark IIIEOS Rebel T3iEOS 60DEOS Rebel T3EOS 7D Mark IIWFT-E2 IIWFT-E4 IIEOS Rebel T4iEOS 6DEOS-1D CEOS 70DEOS Rebel T5iEOS Rebel T5EOS-1D X Mark IIEOS MEOS 80DEOS M2EOS Rebel SL1EOS Rebel T6sEOS 5D Mark IVEOS 5DSEOS Rebel T6iEOS 5DS REOS Rebel T6EOS Rebel T7iEOS 6D Mark IIEOS 77DEOS Rebel SL2EOS R5EOS Rebel T100EOS REOS-1D X Mark IIIEOS Rebel T7EOS RPEOS Rebel T8iEOS SL3EOS 90DEOS R3EOS R6EOS R7EOS R10PowerShot ZOOMEOS M50 Mark IIEOS R50EOS R6 Mark IIEOS R8PowerShot V10EOS R1EOS R5 Mark IIPowerShot V1EOS R100EOS R50 VEOS R6 Mark IIIEOS D2000CEOS D6000C"

var _CanonCameraModel_index = [...]uint16{0, 7, 20, 34, 47, 60, 74, 86, 99, 112, 125, 132, 146, 160, 174, 188, 200, 213, 228, 242, 255, 268, 282, 294, 308, 321, 334, 348, 363, 377, 390, 403, 415, 429, 442, 457, 471, 485, 498, 512, 525, 540, 555, 569, 583, 597, 612, 627, 642, 649, 663, 677, 691, 704, 718, 732, 747, 762, 775, 790, 804, 819, 834, 848, 866, 881, 895, 910, 922, 936, 954, 968, 985, 999, 1013, 1028, 1042, 1060, 1077, 1091, 1106, 1122, 1136, 1150, 1162, 1179, 1196, 1214, 1232, 1250, 1268, 1286, 1304, 1321, 1335, 1349, 1368, 1384, 1401, 1419, 1432, 1450, 1468, 1486, 1504, 1516, 1529, 1547, 1565, 1579, 1597, 1615, 1633, 1651, 1670, 1683, 1701, 1714, 1731, 1749, 1767, 1781, 1795, 1810, 1828, 1847, 1866, 1885, 1903, 1922, 1941, 1954, 1971, 1989, 2002, 2020, 2038, 2059, 2077, 2091, 2112, 2130, 2151, 2166, 2181, 2199, 2213, 2231, 2252, 2272, 2289, 2310, 2317, 2332, 2346, 2367, 2388, 2401, 2419, 2437, 2455, 2476, 2497, 2515, 2533, 2548, 2563, 2576, 2593, 2611, 2631, 2649, 2660, 2671, 2689, 2707, 2725, 2740, 2758, 2773, 2794, 2812, 2833, 2848, 2861, 2875, 2893, 2911, 2931, 2939, 2953, 2975, 2988, 3006, 3024, 3045, 3063, 3084, 3105, 3111, 3128, 3146, 3164, 3178, 3190, 3208, 3226, 3244, 3251, 3265, 3286, 3304, 3325, 3346, 3364, 3378, 3384, 3398, 3420, 3428, 3449, 3467, 3485, 3506, 3518, 3539, 3557, 3575, 3581, 3603, 3610, 3628, 3646, 3664, 3687, 3701, 3719, 3741, 3758, 3781, 3795, 3803, 3810, 3814, 3819, 3823, 3828, 3832, 3836, 3841, 3845, 3849, 3854, 3858, 3862, 3868, 3873, 3878, 3882, 3886, 3890, 3896, 3900, 3905, 3909, 3913, 3919, 3925, 3931, 3937, 3941, 3949, 3955, 3959, 3967, 3973, 3980, 3987, 4002, 4019, 4033, 4040, 4061, 4076, 4096, 4103, 4109, 4125, 4139, 4145, 4161, 4168, 4189, 4195, 4201, 4207, 4220, 4232, 4239, 4247, 4260, 4266, 4272, 4286, 4301, 4314, 4321, 4333, 4347, 4356, 4365, 4378, 4384, 4392, 4399, 4412, 4424, 4440, 4445, 4452, 4458, 4471, 4484, 4498, 4505, 4518, 4527, 4539, 4552, 4566, 4573, 4586, 4592, 4606, 4611, 4628, 4640, 4646, 4659, 4666, 4673, 4679, 4685, 4691, 4698, 4712, 4727, 4734, 4748, 4754, 4767, 4773, 4787, 4799, 4807, 4816, 4831, 4841, 4851}

var _CanonCameraModel_value = [...]CanonCameraModel{
	CanonModelUnknown,
	CanonModelPowerShotA30,
	CanonModelPowerShotS300,
	CanonModelPowerShotA20,
	CanonModelPowerShotA10,
	CanonModelPowerShotS110,
	CanonModelPowerShotG2,
	CanonModelPowerShotS40,
	CanonModelPowerShotS30,
	CanonModelPowerShotA40,
	CanonModelEOSD30,
	CanonModelPowerShotA100,
	CanonModelPowerShotS200,
	CanonModelPowerShotA200,
	CanonModelPowerShotS330,
	CanonModelPowerShotG3,
	CanonModelPowerShotS45,
	CanonModelPowerShotSD100,
	CanonModelPowerShotS230,
	CanonModelPowerShotA70,
	CanonModelPowerShotA60,
	CanonModelPowerShotS400,
	CanonModelPowerShotG5,
	CanonModelPowerShotA300,
	CanonModelPowerShotS50,
	CanonModelPowerShotA80,
	CanonModelPowerShotSD10,
	CanonModelPowerShotS1IS,
	CanonModelPowerShotPro1,
	CanonModelPowerShotS70,
	CanonModelPowerShotS60,
	CanonModelPowerShotG6,
	CanonModelPowerShotS500,
	CanonModelPowerShotA75,
	CanonModelPowerShotSD110,
	CanonModelPowerShotA400,
	CanonModelPowerShotA310,
	CanonModelPowerShotA85,
	CanonModelPowerShotS410,
	CanonModelPowerShotA95,
	CanonModelPowerShotSD300,
	CanonModelPowerShotSD200,
	CanonModelPowerShotA520,
	CanonModelPowerShotA510,
	CanonModelPowerShotSD20,
	CanonModelPowerShotS2IS,
	CanonModelPowerShotSD430,
	CanonModelPowerShotSD500,
	CanonModelEOSD60,
	CanonModelPowerShotSD30,
	CanonModelPowerShotA430,
	CanonModelPowerShotA410,
	CanonModelPowerShotS80,
	CanonModelPowerShotA620,
	CanonModelPowerShotA610,
	CanonModelPowerShotSD630,
	CanonModelPowerShotSD450,
	CanonModelPowerShotTX1,
	CanonModelPowerShotSD400,
	CanonModelPowerShotA420,
	CanonModelPowerShotSD900,
	CanonModelPowerShotSD550,
	CanonModelPowerShotA700,
	CanonModelPowerShotSD700IS,
	CanonModelPowerShotS3IS,
	CanonModelPowerShotA540,
	CanonModelPowerShotSD600,
	CanonModelPowerShotG7,
	CanonModelPowerShotA530,
	CanonModelPowerShotSD800IS,
	CanonModelPowerShotSD40,
	CanonModelPowerShotA710IS,
	CanonModelPowerShotA640,
	CanonModelPowerShotA630,
	CanonModelPowerShotS5IS,
	CanonModelPowerShotA460,
	CanonModelPowerShotSD850IS,
	CanonModelPowerShotA570IS,
	CanonModelPowerShotA560,
	CanonModelPowerShotSD750,
	CanonModelPowerShotSD1000,
	CanonModelPowerShotA550,
	CanonModelPowerShotA450,
	CanonModelPowerShotG9,
	CanonModelPowerShotA650IS,
	CanonModelPowerShotA720IS,
	CanonModelPowerShotSX100IS,
	CanonModelPowerShotSD950IS,
	CanonModelPowerShotSD870IS,
	CanonModelPowerShotSD890IS,
	CanonModelPowerShotSD790IS,
	CanonModelPowerShotSD770IS,
	CanonModelPowerShotA590IS,
	CanonModelPowerShotA580,
	CanonModelPowerShotA470,
	CanonModelPowerShotSD1100IS,
	CanonModelPowerShotSX1IS,
	CanonModelPowerShotSX10IS,
	CanonModelPowerShotA1000IS,
	CanonModelPowerShotG10,
	CanonModelPowerShotA2000IS,
	CanonModelPowerShotSX110IS,
	CanonModelPowerShotSD990IS,
	CanonModelPowerShotSD880IS,
	CanonModelPowerShotE1,
	CanonModelPowerShotD10,
	CanonModelPowerShotSD960IS,
	CanonModelPowerShotA2100IS,
	CanonModelPowerShotA480,
	CanonModelPowerShotSX200IS,
	CanonModelPowerShotSD970IS,
	CanonModelPowerShotSD780IS,
	CanonModelPowerShotA1100IS,
	CanonModelPowerShotSD1200IS,
	CanonModelPowerShotG11,
	CanonModelPowerShotSX120IS,
	CanonModelPowerShotS90,
	CanonModelPowerShotSX20IS,
	CanonModelPowerShotSD980IS,
	CanonModelPowerShotSD940IS,
	CanonModelPowerShotA495,
	CanonModelPowerShotA490,
	CanonModelPowerShotA3100,
	CanonModelPowerShotA3000IS,
	CanonModelPowerShotSD1400IS,
	CanonModelPowerShotSD1300IS,
	CanonModelPowerShotSD3500IS,
	CanonModelPowerShotSX210IS,
	CanonModelPowerShotSD4000IS,
	CanonModelPowerShotSD4500IS,
	CanonModelPowerShotG12,
	CanonModelPowerShotSX30IS,
	CanonModelPowerShotSX130IS,
	CanonModelPowerShotS95,
	CanonModelPowerShotA3300IS,
	CanonModelPowerShotA3200IS,
	CanonModelPowerShotELPH500HS,
	CanonModelPowerShotPro90IS,
	CanonModelPowerShotA800,
	CanonModelPowerShotELPH100HS,
	CanonModelPowerShotSX230HS,
	CanonModelPowerShotELPH300HS,
	CanonModelPowerShotA2200,
	CanonModelPowerShotA1200,
	CanonModelPowerShotSX220HS,
	CanonModelPowerShotG1X,
	CanonModelPowerShotSX150IS,
	CanonModelPowerShotELPH510HS,
	CanonModelPowerShotS100new,
	CanonModelPowerShotSX40HS,
	CanonModelPowerShotELPH310HS,
	CanonModelIXY32S,
	CanonModelPowerShotA1300,
	CanonModelPowerShotA810,
	CanonModelPowerShotELPH320HS,
	CanonModelPowerShotELPH110HS,
	CanonModelPowerShotD20,
	CanonModelPowerShotA4000IS,
	CanonModelPowerShotSX260HS,
	CanonModelPowerShotSX240HS,
	CanonModelPowerShotELPH530HS,
	CanonModelPowerShotELPH520HS,
	CanonModelPowerShotA3400IS,
	CanonModelPowerShotA2400IS,
	CanonModelPowerShotA2300,
	CanonModelPowerShotS100V,
	CanonModelPowerShotG15,
	CanonModelPowerShotSX50HS,
	CanonModelPowerShotSX160IS,
	CanonModelPowerShotS110new,
	CanonModelPowerShotSX500IS,
	CanonModelPowerShotN,
	CanonModelIXUS245HS,
	CanonModelPowerShotSX280HS,
	CanonModelPowerShotSX270HS,
	CanonModelPowerShotA3500IS,
	CanonModelPowerShotA2600,
	CanonModelPowerShotSX275HS,
	CanonModelPowerShotA1400,
	CanonModelPowerShotELPH130IS,
	CanonModelPowerShotELPH115,
	CanonModelPowerShotELPH330HS,
	CanonModelPowerShotA2500,
	CanonModelPowerShotG16,
	CanonModelPowerShotS120,
	CanonModelPowerShotSX170IS,
	CanonModelPowerShotSX510HS,
	CanonModelPowerShotS200new,
	CanonModelIXY620F,
	CanonModelPowerShotN100,
	CanonModelPowerShotG1XMarkII,
	CanonModelPowerShotD30,
	CanonModelPowerShotSX700HS,
	CanonModelPowerShotSX600HS,
	CanonModelPowerShotELPH140IS,
	CanonModelPowerShotELPH135,
	CanonModelPowerShotELPH340HS,
	CanonModelPowerShotELPH150IS,
	CanonModelEOSM3,
	CanonModelPowerShotSX60HS,
	CanonModelPowerShotSX520HS,
	CanonModelPowerShotSX400IS,
	CanonModelPowerShotG7X,
	CanonModelPowerShotN2,
	CanonModelPowerShotSX530HS,
	CanonModelPowerShotSX710HS,
	CanonModelPowerShotSX610HS,
	CanonModelEOSM10,
	CanonModelPowerShotG3X,
	CanonModelPowerShotELPH165HS,
	CanonModelPowerShotELPH160,
	CanonModelPowerShotELPH350HS,
	CanonModelPowerShotELPH170IS,
	CanonModelPowerShotSX410IS,
	CanonModelPowerShotG9X,
	CanonModelEOSM5,
	CanonModelPowerShotG5X,
	CanonModelPowerShotG7XMarkII,
	CanonModelEOSM100,
	CanonModelPowerShotELPH360HS,
	CanonModelPowerShotSX540HS,
	CanonModelPowerShotSX420IS,
	CanonModelPowerShotELPH190IS,
	CanonModelPowerShotG1,
	CanonModelPowerShotELPH180IS,
	CanonModelPowerShotSX720HS,
	CanonModelPowerShotSX620HS,
	CanonModelEOSM6,
	CanonModelPowerShotG9XMarkII,
	CanonModelEOSM50,
	CanonModelPowerShotELPH185,
	CanonModelPowerShotSX430IS,
	CanonModelPowerShotSX730HS,
	CanonModelPowerShotG1XMarkIII,
	CanonModelPowerShotS100,
	CanonModelPowerShotSX740HS,
	CanonModelPowerShotG5XMarkII,
	CanonModelPowerShotSX70HS,
	CanonModelPowerShotG7XMarkIII,
	CanonModelEOSM6MarkII,
	CanonModelEOSM200,
	CanonModelEOSC50,
	CanonModelDC19,
	CanonModelXHA1,
	CanonModelHV10,
	CanonModelMD130,
	CanonModelDC50,
	CanonModelHV20,
	CanonModelDC211,
	CanonModelHG10,
	CanonModelHR10,
	CanonModelMD255,
	CanonModelHF11,
	CanonModelHV30,
	CanonModelXHA1S,
	CanonModelDC301,
	CanonModelFS100,
	CanonModelHF10,
	CanonModelHG20,
	CanonModelHF21,
	CanonModelHFS11,
	CanonModelHV40,
	CanonModelDC410,
	CanonModelFS19,
	CanonModelHF20,
	CanonModelHFS10,
	CanonModelHFR10,
	CanonModelHFM30,
	CanonModelHFS20,
	CanonModelFS31,
	CanonModelEOSC300,
	CanonModelHFG25,
	CanonModelXC10,
	CanonModelEOSC200,
	CanonModelEOS1D,
	CanonModelEOS1DS,
	CanonModelEOS10D,
	CanonModelEOS1DMarkIII,
	CanonModelEOSDigitalRebel,
	CanonModelEOS1DMarkII,
	CanonModelEOS20D,
	CanonModelEOSDigitalRebelXSi,
	CanonModelEOS1DsMarkII,
	CanonModelEOSDigitalRebelXT,
	CanonModelEOS40D,
	CanonModelEOS5D,
	CanonModelEOS1DsMarkIII,
	CanonModelEOS5DMarkII,
	CanonModelWFTE1,
	CanonModelEOS1DMarkIIN,
	CanonModelEOS30D,
	CanonModelEOSDigitalRebelXTi,
	CanonModelWFTE2,
	CanonModelWFTE3,
	CanonModelEOS7D,
	CanonModelEOSRebelT1i,
	CanonModelEOSRebelXS,
	CanonModelEOS50D,
	CanonModelEOS1DX,
	CanonModelEOSRebelT2i,
	CanonModelWFTE4,
	CanonModelWFTE5,
	CanonModelEOS1DMarkIV,
	CanonModelEOS5DMarkIII,
	CanonModelEOSRebelT3i,
	CanonModelEOS60D,
	CanonModelEOSRebelT3,
	CanonModelEOS7DMarkII,
	CanonModelWFTE2II,
	CanonModelWFTE4II,
	CanonModelEOSRebelT4i,
	CanonModelEOS6D,
	CanonModelEOS1DC,
	CanonModelEOS70D,
	CanonModelEOSRebelT5i,
	CanonModelEOSRebelT5,
	CanonModelEOS1DXMarkII,
	CanonModelEOSM,
	CanonModelEOS80D,
	CanonModelEOSM2,
	CanonModelEOSRebelSL1,
	CanonModelEOSRebelT6s,
	CanonModelEOS5DMarkIV,
	CanonModelEOS5DS,
	CanonModelEOSRebelT6i,
	CanonModelEOS5DSR,
	CanonModelEOSRebelT6,
	CanonModelEOSRebelT7i,
	CanonModelEOS6DMarkII,
	CanonModelEOS77D,
	CanonModelEOSRebelSL2,
	CanonModelEOSR5,
	CanonModelEOSRebelT100,
	CanonModelEOSR,
	CanonModelEOS1DXMarkIII,
	CanonModelEOSRebelT7,
	CanonModelEOSRP,
	CanonModelEOSRebelT8i,
	CanonModelEOSSL3,
	CanonModelEOS90D,
	CanonModelEOSR3,
	CanonModelEOSR6,
	CanonModelEOSR7,
	CanonModelEOSR10,
	CanonModelPowerShotZOOM,
	CanonModelEOSM50MarkII,
	CanonModelEOSR50,
	CanonModelEOSR6MarkII,
	CanonModelEOSR8,
	CanonModelPowerShotV10,
	CanonModelEOSR1,
	CanonModelEOSR5MarkII,
	CanonModelPowerShotV1,
	CanonModelEOSR100,
	CanonModelEOSR50V,
	CanonModelEOSR6MarkIII,
	CanonModelEOSD2000C,
	CanonModelEOSD6000C,
}

func (m CanonCameraModel) String() string {
	for i, v := range _CanonCameraModel_value {
		if m == v {
			return _CanonCameraModel_name[_CanonCameraModel_index[i]:_CanonCameraModel_index[i+1]]
		}
	}
	return "CanonCameraModel(0x" + strconv.FormatUint(uint64(m), 16) + ")"
}
