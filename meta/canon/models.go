package canon

import "fmt"

//go:generate stringer -type=CanonModelID

// CanonModelID represents a Canon model ID numbers.
// Based on Phil Harvey's exiftool
// Updated: 2024-09-09
// Reference: https://exiftool.org/TagNames/Canon.html
// TODO: Implement support for Third Party Lenses
type CanonModelID uint32

const (
	ModelPowerShotA30      CanonModelID = 0x1010000
	ModelPowerShotS300     CanonModelID = 0x1040000
	ModelPowerShotA20      CanonModelID = 0x1060000
	ModelPowerShotA10      CanonModelID = 0x1080000
	ModelPowerShotS110     CanonModelID = 0x1090000
	ModelPowerShotG2       CanonModelID = 0x1100000
	ModelPowerShotS40      CanonModelID = 0x1110000
	ModelPowerShotS30      CanonModelID = 0x1120000
	ModelPowerShotA40      CanonModelID = 0x1130000
	ModelEOSD30            CanonModelID = 0x1140000
	ModelPowerShotA100     CanonModelID = 0x1150000
	ModelPowerShotS200     CanonModelID = 0x1160000
	ModelPowerShotA200     CanonModelID = 0x1170000
	ModelPowerShotS330     CanonModelID = 0x1180000
	ModelPowerShotG3       CanonModelID = 0x1190000
	ModelPowerShotS45      CanonModelID = 0x1210000
	ModelPowerShotSD100    CanonModelID = 0x1230000
	ModelPowerShotS230     CanonModelID = 0x1240000
	ModelPowerShotA70      CanonModelID = 0x1250000
	ModelPowerShotA60      CanonModelID = 0x1260000
	ModelPowerShotS400     CanonModelID = 0x1270000
	ModelPowerShotG5       CanonModelID = 0x1290000
	ModelPowerShotA300     CanonModelID = 0x1300000
	ModelPowerShotS50      CanonModelID = 0x1310000
	ModelPowerShotA80      CanonModelID = 0x1340000
	ModelPowerShotSD10     CanonModelID = 0x1350000
	ModelPowerShotS1IS     CanonModelID = 0x1360000
	ModelPowerShotPro1     CanonModelID = 0x1370000
	ModelPowerShotS70      CanonModelID = 0x1380000
	ModelPowerShotS60      CanonModelID = 0x1390000
	ModelPowerShotG6       CanonModelID = 0x1400000
	ModelPowerShotS500     CanonModelID = 0x1410000
	ModelPowerShotA75      CanonModelID = 0x1420000
	ModelPowerShotSD110    CanonModelID = 0x1440000
	ModelPowerShotA400     CanonModelID = 0x1450000
	ModelPowerShotA310     CanonModelID = 0x1470000
	ModelPowerShotA85      CanonModelID = 0x1490000
	ModelPowerShotS410     CanonModelID = 0x1520000
	ModelPowerShotA95      CanonModelID = 0x1530000
	ModelPowerShotSD300    CanonModelID = 0x1540000
	ModelPowerShotSD200    CanonModelID = 0x1550000
	ModelPowerShotA520     CanonModelID = 0x1560000
	ModelPowerShotA510     CanonModelID = 0x1570000
	ModelPowerShotSD20     CanonModelID = 0x1590000
	ModelPowerShotS2IS     CanonModelID = 0x1640000
	ModelPowerShotSD430    CanonModelID = 0x1650000
	ModelPowerShotSD500    CanonModelID = 0x1660000
	ModelEOSD60            CanonModelID = 0x1668000
	ModelPowerShotSD30     CanonModelID = 0x1700000
	ModelPowerShotA430     CanonModelID = 0x1740000
	ModelPowerShotA410     CanonModelID = 0x1750000
	ModelPowerShotS80      CanonModelID = 0x1760000
	ModelPowerShotA620     CanonModelID = 0x1780000
	ModelPowerShotA610     CanonModelID = 0x1790000
	ModelPowerShotSD630    CanonModelID = 0x1800000
	ModelPowerShotSD450    CanonModelID = 0x1810000
	ModelPowerShotTX1      CanonModelID = 0x1820000
	ModelPowerShotSD400    CanonModelID = 0x1870000
	ModelPowerShotA420     CanonModelID = 0x1880000
	ModelPowerShotSD900    CanonModelID = 0x1890000
	ModelPowerShotSD550    CanonModelID = 0x1900000
	ModelPowerShotA700     CanonModelID = 0x1920000
	ModelPowerShotSD700    CanonModelID = 0x1940000
	ModelPowerShotS3IS     CanonModelID = 0x1950000
	ModelPowerShotA540     CanonModelID = 0x1960000
	ModelPowerShotSD600    CanonModelID = 0x1970000
	ModelPowerShotG7       CanonModelID = 0x1980000
	ModelPowerShotA530     CanonModelID = 0x1990000
	ModelPowerShotSD800    CanonModelID = 0x2000000
	ModelPowerShotSD40     CanonModelID = 0x2010000
	ModelPowerShotA710     CanonModelID = 0x2020000
	ModelPowerShotA640     CanonModelID = 0x2030000
	ModelPowerShotA630     CanonModelID = 0x2040000
	ModelPowerShotS5IS     CanonModelID = 0x2090000
	ModelPowerShotA460     CanonModelID = 0x2100000
	ModelPowerShotSD850    CanonModelID = 0x2120000
	ModelPowerShotA570     CanonModelID = 0x2130000
	ModelPowerShotA560     CanonModelID = 0x2140000
	ModelPowerShotSD750    CanonModelID = 0x2150000
	ModelPowerShotSD1000   CanonModelID = 0x2160000
	ModelPowerShotA550     CanonModelID = 0x2180000
	ModelPowerShotA450     CanonModelID = 0x2190000
	ModelPowerShotG9       CanonModelID = 0x2230000
	ModelPowerShotA650     CanonModelID = 0x2240000
	ModelPowerShotA720     CanonModelID = 0x2260000
	ModelPowerShotSX100    CanonModelID = 0x2290000
	ModelPowerShotSD950    CanonModelID = 0x2300000
	ModelPowerShotSD870    CanonModelID = 0x2310000
	ModelPowerShotSD890    CanonModelID = 0x2320000
	ModelPowerShotSD790    CanonModelID = 0x2360000
	ModelPowerShotSD770    CanonModelID = 0x2370000
	ModelPowerShotA590     CanonModelID = 0x2380000
	ModelPowerShotA580     CanonModelID = 0x2390000
	ModelPowerShotA470     CanonModelID = 0x2420000
	ModelPowerShotSD1100   CanonModelID = 0x2430000
	ModelPowerShotSX1      CanonModelID = 0x2460000
	ModelPowerShotSX10     CanonModelID = 0x2470000
	ModelPowerShotA1000    CanonModelID = 0x2480000
	ModelPowerShotG10      CanonModelID = 0x2490000
	ModelPowerShotA2000    CanonModelID = 0x2510000
	ModelPowerShotSX110    CanonModelID = 0x2520000
	ModelPowerShotSD990    CanonModelID = 0x2530000
	ModelPowerShotSD880    CanonModelID = 0x2540000
	ModelPowerShotE1       CanonModelID = 0x2550000
	ModelPowerShotD10      CanonModelID = 0x2560000
	ModelPowerShotSD960    CanonModelID = 0x2570000
	ModelPowerShotA2100    CanonModelID = 0x2580000
	ModelPowerShotA480     CanonModelID = 0x2590000
	ModelPowerShotSX200    CanonModelID = 0x2600000
	ModelPowerShotSD970    CanonModelID = 0x2610000
	ModelPowerShotSD780    CanonModelID = 0x2620000
	ModelPowerShotA1100    CanonModelID = 0x2630000
	ModelPowerShotSD1200   CanonModelID = 0x2640000
	ModelPowerShotG11      CanonModelID = 0x2700000
	ModelPowerShotSX120    CanonModelID = 0x2710000
	ModelPowerShotS90      CanonModelID = 0x2720000
	ModelPowerShotSX20     CanonModelID = 0x2750000
	ModelPowerShotSD980    CanonModelID = 0x2760000
	ModelPowerShotSD940    CanonModelID = 0x2770000
	ModelPowerShotA495     CanonModelID = 0x2800000
	ModelPowerShotA490     CanonModelID = 0x2810000
	ModelPowerShotA3100    CanonModelID = 0x2820000
	ModelPowerShotA3000    CanonModelID = 0x2830000
	ModelPowerShotSD1400   CanonModelID = 0x2840000
	ModelPowerShotSD1300   CanonModelID = 0x2850000
	ModelPowerShotSD3500   CanonModelID = 0x2860000
	ModelPowerShotSX210    CanonModelID = 0x2870000
	ModelPowerShotSD4000   CanonModelID = 0x2880000
	ModelPowerShotSD4500   CanonModelID = 0x2890000
	ModelPowerShotG12      CanonModelID = 0x2920000
	ModelPowerShotSX30     CanonModelID = 0x2930000
	ModelPowerShotSX130    CanonModelID = 0x2940000
	ModelPowerShotS95      CanonModelID = 0x2950000
	ModelPowerShotA3300    CanonModelID = 0x2980000
	ModelPowerShotA3200    CanonModelID = 0x2990000
	ModelPowerShotELPH500  CanonModelID = 0x3000000
	ModelPowerShotPro90    CanonModelID = 0x3010000
	ModelPowerShotA800     CanonModelID = 0x3010001
	ModelPowerShotELPH100  CanonModelID = 0x3020000
	ModelPowerShotSX230    CanonModelID = 0x3030000
	ModelPowerShotELPH300  CanonModelID = 0x3040000
	ModelPowerShotA2200    CanonModelID = 0x3050000
	ModelPowerShotA1200    CanonModelID = 0x3060000
	ModelPowerShotSX220    CanonModelID = 0x3070000
	ModelPowerShotG1X      CanonModelID = 0x3080000
	ModelPowerShotSX150    CanonModelID = 0x3090000
	ModelPowerShotELPH510  CanonModelID = 0x3100000
	ModelPowerShotS100New  CanonModelID = 0x3110000
	ModelPowerShotSX40     CanonModelID = 0x3130000
	ModelPowerShotELPH310  CanonModelID = 0x3120000
	ModelIXY32S            CanonModelID = 0x3140000
	ModelPowerShotA1300    CanonModelID = 0x3160000
	ModelPowerShotA810     CanonModelID = 0x3170000
	ModelPowerShotELPH320  CanonModelID = 0x3180000
	ModelPowerShotELPH110  CanonModelID = 0x3190000
	ModelPowerShotD20      CanonModelID = 0x3200000
	ModelPowerShotA4000    CanonModelID = 0x3210000
	ModelPowerShotSX260    CanonModelID = 0x3220000
	ModelPowerShotSX240    CanonModelID = 0x3230000
	ModelPowerShotELPH530  CanonModelID = 0x3240000
	ModelPowerShotELPH520  CanonModelID = 0x3250000
	ModelPowerShotA3400    CanonModelID = 0x3260000
	ModelPowerShotA2400    CanonModelID = 0x3270000
	ModelPowerShotA2300    CanonModelID = 0x3280000
	ModelPowerShotS100V    CanonModelID = 0x3320000
	ModelPowerShotG15      CanonModelID = 0x3330000
	ModelPowerShotSX50     CanonModelID = 0x3340000
	ModelPowerShotSX160    CanonModelID = 0x3350000
	ModelPowerShotS110New  CanonModelID = 0x3360000
	ModelPowerShotSX500    CanonModelID = 0x3370000
	ModelPowerShotN        CanonModelID = 0x3380000
	ModelIXUS245           CanonModelID = 0x3390000
	ModelPowerShotSX280    CanonModelID = 0x3400000
	ModelPowerShotSX270    CanonModelID = 0x3410000
	ModelPowerShotA3500    CanonModelID = 0x3420000
	ModelPowerShotA2600    CanonModelID = 0x3430000
	ModelPowerShotSX275    CanonModelID = 0x3440000
	ModelPowerShotA1400    CanonModelID = 0x3450000
	ModelPowerShotELPH130  CanonModelID = 0x3460000
	ModelPowerShotELPH115  CanonModelID = 0x3470000
	ModelPowerShotELPH330  CanonModelID = 0x3490000
	ModelPowerShotA2500    CanonModelID = 0x3510000
	ModelPowerShotG16      CanonModelID = 0x3540000
	ModelPowerShotS120     CanonModelID = 0x3550000
	ModelPowerShotSX170    CanonModelID = 0x3560000
	ModelPowerShotSX510    CanonModelID = 0x3580000
	ModelPowerShotS200New  CanonModelID = 0x3590000
	ModelIXY620F           CanonModelID = 0x3600000
	ModelPowerShotN100     CanonModelID = 0x3610000
	ModelPowerShotG1XMark2 CanonModelID = 0x3640000
	ModelPowerShotD30      CanonModelID = 0x3650000
	ModelPowerShotSX700    CanonModelID = 0x3660000
	ModelPowerShotSX600    CanonModelID = 0x3670000
	ModelPowerShotELPH140  CanonModelID = 0x3680000
	ModelPowerShotELPH135  CanonModelID = 0x3690000
	ModelPowerShotELPH340  CanonModelID = 0x3700000
	ModelPowerShotELPH150  CanonModelID = 0x3710000
	ModelPowerShotSX60     CanonModelID = 0x3750000
	ModelPowerShotSX520    CanonModelID = 0x3760000
	ModelPowerShotSX400    CanonModelID = 0x3770000
	ModelPowerShotG7X      CanonModelID = 0x3780000
	ModelPowerShotN2       CanonModelID = 0x3790000
	ModelPowerShotSX530    CanonModelID = 0x3800000
	ModelPowerShotSX710    CanonModelID = 0x3820000
	ModelPowerShotSX610    CanonModelID = 0x3830000
	ModelPowerShotG3X      CanonModelID = 0x3850000
	ModelPowerShotELPH165  CanonModelID = 0x3860000
	ModelPowerShotELPH160  CanonModelID = 0x3870000
	ModelPowerShotELPH350  CanonModelID = 0x3880000
	ModelPowerShotELPH170  CanonModelID = 0x3890000
	ModelPowerShotSX410    CanonModelID = 0x3910000
	ModelPowerShotG9X      CanonModelID = 0x3930000
	ModelPowerShotG5X      CanonModelID = 0x3950000
	ModelPowerShotG7XMark2 CanonModelID = 0x3970000
	ModelPowerShotELPH360  CanonModelID = 0x3990000
	ModelPowerShotSX540    CanonModelID = 0x4010000
	ModelPowerShotSX420    CanonModelID = 0x4020000
	ModelPowerShotELPH190  CanonModelID = 0x4030000
	ModelPowerShotG1       CanonModelID = 0x4040000
	ModelPowerShotELPH180  CanonModelID = 0x4040001
	ModelPowerShotSX720    CanonModelID = 0x4050000
	ModelPowerShotSX620    CanonModelID = 0x4060000
	ModelPowerShotG9XMark2 CanonModelID = 0x4100000
	ModelPowerShotELPH185  CanonModelID = 0x4150000
	ModelPowerShotSX430    CanonModelID = 0x4160000
	ModelPowerShotSX730    CanonModelID = 0x4170000
	ModelPowerShotG1XMark3 CanonModelID = 0x4180000
	ModelPowerShotS100     CanonModelID = 0x6040000
	ModelPowerShotSX740    CanonModelID = 0x801
	ModelPowerShotG5XMark2 CanonModelID = 0x804
	ModelPowerShotSX70     CanonModelID = 0x805
	ModelPowerShotG7XMark3 CanonModelID = 0x808

	// Video Cameras
	ModelDC19    CanonModelID = 0x4007d673
	ModelXHA1    CanonModelID = 0x4007d674
	ModelHV10    CanonModelID = 0x4007d675
	ModelMD130   CanonModelID = 0x4007d676
	ModelDC50    CanonModelID = 0x4007d777
	ModelHV20    CanonModelID = 0x4007d778
	ModelDC211   CanonModelID = 0x4007d779
	ModelHG10    CanonModelID = 0x4007d77a
	ModelHR10    CanonModelID = 0x4007d77b
	ModelMD255   CanonModelID = 0x4007d77d
	ModelHF11    CanonModelID = 0x4007d81c
	ModelHV30    CanonModelID = 0x4007d878
	ModelXHA1S   CanonModelID = 0x4007d87c
	ModelDC301   CanonModelID = 0x4007d87e
	ModelFS100   CanonModelID = 0x4007d87f
	ModelHF10    CanonModelID = 0x4007d880
	ModelHG20    CanonModelID = 0x4007d882
	ModelHF21    CanonModelID = 0x4007d925
	ModelHFS11   CanonModelID = 0x4007d926
	ModelHV40    CanonModelID = 0x4007d978
	ModelDC410   CanonModelID = 0x4007d987
	ModelFS19    CanonModelID = 0x4007d988
	ModelHF20    CanonModelID = 0x4007d989
	ModelHFS10   CanonModelID = 0x4007d98a
	ModelHFR10   CanonModelID = 0x4007da8e
	ModelHFM30   CanonModelID = 0x4007da8f
	ModelHFS20   CanonModelID = 0x4007da90
	ModelFS31    CanonModelID = 0x4007da92
	ModelEOSC300 CanonModelID = 0x4007dca0
	ModelHFG25   CanonModelID = 0x4007dda9
	ModelXC10    CanonModelID = 0x4007dfb4
	ModelEOSC200 CanonModelID = 0x4007e1c3

	// Professional EOS Models
	ModelEOS1D         CanonModelID = 0x80000001
	ModelEOS1DS        CanonModelID = 0x80000167
	ModelEOS1DMarkIII  CanonModelID = 0x80000169
	ModelEOS1DMarkII   CanonModelID = 0x80000174
	ModelEOS1DSMarkII  CanonModelID = 0x80000188
	ModelEOS1DSMarkIII CanonModelID = 0x80000215
	ModelEOS1DMarkIIN  CanonModelID = 0x80000232
	ModelEOS1DX        CanonModelID = 0x80000269
	ModelEOS1DMarkIV   CanonModelID = 0x80000281
	ModelEOS1DC        CanonModelID = 0x80000324
	ModelEOS1DXMarkII  CanonModelID = 0x80000328
	ModelEOS1DXMarkIII CanonModelID = 0x80000428

	// Prosumer EOS Models
	ModelEOS5D        CanonModelID = 0x80000213
	ModelEOS5DMarkII  CanonModelID = 0x80000218
	ModelEOS5DMarkIII CanonModelID = 0x80000285
	ModelEOS7D        CanonModelID = 0x80000250
	ModelEOS7DMarkII  CanonModelID = 0x80000289 // IB
	ModelEOS6D        CanonModelID = 0x80000302 // 25

	// Consumer EOS Models
	ModelEOS10D CanonModelID = 0x80000168
	ModelEOS20D CanonModelID = 0x80000175
	ModelEOS30D CanonModelID = 0x80000234
	ModelEOS40D CanonModelID = 0x80000190
	ModelEOS50D CanonModelID = 0x80000261
	ModelEOS60D CanonModelID = 0x80000287
	ModelEOS70D CanonModelID = 0x80000325

	// Rebel/Kiss Series
	// Wireless File Transmitters
	ModelWFTE1   CanonModelID = 0x80000219
	ModelWFTE2   CanonModelID = 0x80000241
	ModelWFTE3   CanonModelID = 0x80000246
	ModelWFTE4   CanonModelID = 0x80000271
	ModelWFTE5   CanonModelID = 0x80000273
	ModelWFTE2II CanonModelID = 0x80000297
	ModelWFTE4II CanonModelID = 0x80000298

	// DSLR Models
	ModelEOS80D      CanonModelID = 0x80000350 // 42
	ModelEOS5DMarkIV CanonModelID = 0x80000349 // 42
	ModelEOS5DS      CanonModelID = 0x80000382
	ModelEOS5DSR     CanonModelID = 0x80000401
	ModelEOS6DMarkII CanonModelID = 0x80000406 // IB/42
	ModelEOS77D      CanonModelID = 0x80000408 // 9000D
	ModelEOS90D      CanonModelID = 0x80000437 // IB

	// Rebel/Kiss Series
	ModelEOSRebel     CanonModelID = 0x80000170 // Digital Rebel / 300D / Kiss Digital
	ModelEOSRebelXT   CanonModelID = 0x80000189 // 350D / Kiss Digital N
	ModelEOSRebelXTi  CanonModelID = 0x80000236 // 400D / Kiss Digital X
	ModelEOSRebelXSi  CanonModelID = 0x80000176 // 450D / Kiss X2
	ModelEOS450D      CanonModelID = 0x80000176 // Same as ModelEOSRebelXSi
	ModelEOSRebelT1i  CanonModelID = 0x80000252 // 500D / Kiss X3
	ModelEOS500D      CanonModelID = 0x80000252 // Same as ModelEOSRebelT1i
	ModelEOSRebelT2i  CanonModelID = 0x80000270 // 550D / Kiss X4
	ModelEOS550D      CanonModelID = 0x80000270 // Same as ModelEOSRebelT2i
	ModelEOSRebelT3i  CanonModelID = 0x80000286 // 600D / Kiss X5
	ModelEOS600D      CanonModelID = 0x80000286 // Same as ModelEOSRebelT3i
	ModelEOSRebelT4i  CanonModelID = 0x80000301 // 650D / Kiss X6i
	ModelEOS650D      CanonModelID = 0x80000301 // Same as ModelEOSRebelT4i
	ModelEOSRebelT5i  CanonModelID = 0x80000326 // 700D / Kiss X7i
	ModelEOS700D      CanonModelID = 0x80000326 // Same as ModelEOSRebelT5i
	ModelEOSRebelXS   CanonModelID = 0x80000254 // 1000D / Kiss F
	ModelEOS1000D     CanonModelID = 0x80000254 // Same as ModelEOSRebelXS
	ModelEOSRebelT3   CanonModelID = 0x80000288 // 1100D / Kiss X50
	ModelEOS1100D     CanonModelID = 0x80000288 // Same as ModelEOSRebelT3
	ModelEOSRebelT5   CanonModelID = 0x80000327 // 1200D / Kiss X70 / Hi
	ModelEOS1200D     CanonModelID = 0x80000327 // Same as ModelEOSRebelT5
	ModelEOSRebelSL1  CanonModelID = 0x80000346 // 100D / Kiss X7
	ModelEOSRebelT6s  CanonModelID = 0x80000347 // 760D / 8000D
	ModelEOS760D      CanonModelID = 0x80000347 // Same as ModelEOSRebelT6s
	ModelEOSRebelT6i  CanonModelID = 0x80000393 // 750D / Kiss X8i
	ModelEOS750D      CanonModelID = 0x80000393 // Same as ModelEOSRebelT6i
	ModelEOSRebelT6   CanonModelID = 0x80000404 // 1300D / Kiss X80
	ModelEOSRebelT7i  CanonModelID = 0x80000405 // 800D / Kiss X9i
	ModelEOSRebelSL2  CanonModelID = 0x80000417 // 200D / Kiss X9
	ModelEOSRebelT100 CanonModelID = 0x80000422 // 4000D / 3000D
	ModelEOSRebelT7   CanonModelID = 0x80000432 // 2000D / 1500D / Kiss X90
	ModelEOSRebelT8i  CanonModelID = 0x80000435 // 850D / X10i
	ModelEOSSL3       CanonModelID = 0x80000436 // 250D / Kiss X10

	// EOS R Series (Mirrorless)
	ModelEOSR5       CanonModelID = 0x80000421
	ModelEOSR        CanonModelID = 0x80000424
	ModelEOSRP       CanonModelID = 0x80000433
	ModelEOSR3       CanonModelID = 0x80000450
	ModelEOSR6       CanonModelID = 0x80000453
	ModelEOSR7       CanonModelID = 0x80000464
	ModelEOSR10      CanonModelID = 0x80000465
	ModelEOSR50      CanonModelID = 0x80000480
	ModelEOSR50V     CanonModelID = 0x80000516
	ModelEOSR6MarkII CanonModelID = 0x80000481
	ModelEOSR8       CanonModelID = 0x80000487
	ModelEOSR1       CanonModelID = 0x80000495
	ModelR5MarkII    CanonModelID = 0x80000496
	ModelEOSR100     CanonModelID = 0x80000498

	// EOS M Series (Mirrorless)
	ModelEOSM         CanonModelID = 0x80000331
	ModelEOSM2        CanonModelID = 0x80000355
	ModelEOSM10       CanonModelID = 0x3840000
	ModelEOSM3        CanonModelID = 0x3740000
	ModelEOSM5        CanonModelID = 0x3940000
	ModelEOSM100      CanonModelID = 0x3980000
	ModelEOSM200      CanonModelID = 0x812
	ModelEOSM50       CanonModelID = 0x412      // Special case without "0000"
	ModelEOSM50MarkII CanonModelID = 0x80000468 // IB
	ModelEOSM6        CanonModelID = 0x4070000
	ModelEOSM6MarkII  CanonModelID = 0x811

	// PowerShot Series
	ModelPowerShotZOOM CanonModelID = 0x80000467
	ModelPowerShotV10  CanonModelID = 0x80000491 // 25
	ModelPowerShotV1   CanonModelID = 0x80000497

	// Cinema Cameras
	ModelEOSD2000C CanonModelID = 0x80000520 // IB
	ModelEOSD6000C CanonModelID = 0x80000560 // PH (guess)
	ModelEOSC500   CanonModelID = 0x4007e1c4
	ModelEOSC700   CanonModelID = 0x4007e1c5
)

// canonModelIDMap maps CanonModelID values to their string representations.
var canonModelIDMap = map[CanonModelID]string{
	ModelPowerShotA30:      "PowerShot A30",
	ModelPowerShotS300:     "PowerShot S300 / Digital IXUS 300 / IXY Digital ",
	ModelPowerShotA20:      "PowerShot A20",
	ModelPowerShotA10:      "PowerShot A10",
	ModelPowerShotS110:     "PowerShot S110 / Digital IXUS v / IXY Digital 200",
	ModelPowerShotG2:       "PowerShot G2",
	ModelPowerShotS40:      "PowerShot S40",
	ModelPowerShotS30:      "PowerShot S30",
	ModelPowerShotA40:      "PowerShot A40",
	ModelEOSD30:            "EOS D30",
	ModelPowerShotA100:     "PowerShot A100",
	ModelPowerShotS200:     "PowerShot S200 / Digital IXUS v2 / IXY Digital 200a",
	ModelPowerShotA200:     "PowerShot A200",
	ModelPowerShotS330:     "PowerShot S330 / Digital IXUS 330 / IXY Digital 300a",
	ModelPowerShotG3:       "PowerShot G3",
	ModelPowerShotS45:      "PowerShot S45",
	ModelPowerShotSD100:    "PowerShot SD100 / Digital IXUS II / IXY Digital 30",
	ModelPowerShotS230:     "PowerShot S230 / Digital IXUS v3 / IXY Digital 320",
	ModelPowerShotA70:      "PowerShot A70",
	ModelPowerShotA60:      "PowerShot A60",
	ModelPowerShotS400:     "PowerShot S400 / Digital IXUS 400 / IXY Digital 400",
	ModelPowerShotG5:       "PowerShot G5",
	ModelPowerShotA300:     "PowerShot A300",
	ModelPowerShotS50:      "PowerShot S50",
	ModelPowerShotA80:      "PowerShot A80",
	ModelPowerShotSD10:     "PowerShot SD10 / Digital IXUS i / IXY Digital L",
	ModelPowerShotS1IS:     "PowerShot S1 IS",
	ModelPowerShotPro1:     "PowerShot Pro1",
	ModelPowerShotS70:      "PowerShot S70",
	ModelPowerShotS60:      "PowerShot S60",
	ModelPowerShotG6:       "PowerShot G6",
	ModelPowerShotS500:     "PowerShot S500 / Digital IXUS 500 / IXY Digital 500",
	ModelPowerShotA75:      "PowerShot A75",
	ModelPowerShotSD110:    "PowerShot SD110 / Digital IXUS IIs / IXY Digital 30a",
	ModelPowerShotA400:     "PowerShot A400",
	ModelPowerShotA310:     "PowerShot A310",
	ModelPowerShotA85:      "PowerShot A85",
	ModelPowerShotS410:     "PowerShot S410 / Digital IXUS 430 / IXY Digital 450",
	ModelPowerShotA95:      "PowerShot A95",
	ModelPowerShotSD300:    "PowerShot SD300 / Digital IXUS 40 / IXY Digital 50",
	ModelPowerShotSD200:    "PowerShot SD200 / Digital IXUS 30 / IXY Digital 40",
	ModelPowerShotA520:     "PowerShot A520",
	ModelPowerShotA510:     "PowerShot A510",
	ModelPowerShotSD20:     "PowerShot SD20 / Digital IXUS i5 / IXY Digital L2",
	ModelPowerShotS2IS:     "PowerShot S2 IS",
	ModelPowerShotSD430:    "PowerShot SD430 / Digital IXUS Wireless / IXY Digital Wireless",
	ModelPowerShotSD500:    "PowerShot SD500 / Digital IXUS 700 / IXY Digital 600",
	ModelEOSD60:            "EOS D60",
	ModelPowerShotSD30:     "PowerShot SD30 / Digital IXUS i Zoom / IXY Digital L3",
	ModelPowerShotA430:     "PowerShot A430",
	ModelPowerShotA410:     "PowerShot A410",
	ModelPowerShotS80:      "PowerShot S80",
	ModelPowerShotA620:     "PowerShot A620",
	ModelPowerShotA610:     "PowerShot A610",
	ModelPowerShotSD630:    "PowerShot SD630 / Digital IXUS 65 / IXY Digital 80",
	ModelPowerShotSD450:    "PowerShot SD450 / Digital IXUS 55 / IXY Digital 60",
	ModelPowerShotTX1:      "PowerShot TX1",
	ModelPowerShotSD400:    "PowerShot SD400 / Digital IXUS 50 / IXY Digital 55",
	ModelPowerShotA420:     "PowerShot A420",
	ModelPowerShotSD900:    "PowerShot SD900 / Digital IXUS 900 Ti / IXY Digital 1000",
	ModelPowerShotSD550:    "PowerShot SD550 / Digital IXUS 750 / IXY Digital 700",
	ModelPowerShotA700:     "PowerShot A700",
	ModelPowerShotSD700:    "PowerShot SD700 IS / Digital IXUS 800 IS / IXY Digital 800 IS",
	ModelPowerShotS3IS:     "PowerShot S3 IS",
	ModelPowerShotA540:     "PowerShot A540",
	ModelPowerShotSD600:    "PowerShot SD600 / Digital IXUS 60 / IXY Digital 70",
	ModelPowerShotG7:       "PowerShot G7",
	ModelPowerShotA530:     "PowerShot A530",
	ModelPowerShotSD800:    "PowerShot SD800 IS / Digital IXUS 850 IS / IXY Digital 900 IS",
	ModelPowerShotSD40:     "PowerShot SD40 / Digital IXUS i7 / IXY Digital L4",
	ModelPowerShotA710:     "PowerShot A710 IS",
	ModelPowerShotA640:     "PowerShot A640",
	ModelPowerShotA630:     "PowerShot A630",
	ModelPowerShotS5IS:     "PowerShot S5 IS",
	ModelPowerShotA460:     "PowerShot A460",
	ModelPowerShotSD850:    "PowerShot SD850 IS / Digital IXUS 950 IS / IXY Digital 810 IS",
	ModelPowerShotA570:     "PowerShot A570 IS",
	ModelPowerShotA560:     "PowerShot A560",
	ModelPowerShotSD750:    "PowerShot SD750 / Digital IXUS 75 / IXY Digital 90",
	ModelPowerShotSD1000:   "PowerShot SD1000 / Digital IXUS 70 / IXY Digital 10",
	ModelPowerShotA550:     "PowerShot A550",
	ModelPowerShotA450:     "PowerShot A450",
	ModelPowerShotG9:       "PowerShot G9",
	ModelPowerShotA650:     "PowerShot A650 IS",
	ModelPowerShotA720:     "PowerShot A720 IS",
	ModelPowerShotSX100:    "PowerShot SX100 IS",
	ModelPowerShotSD950:    "PowerShot SD950 IS / Digital IXUS 960 IS / IXY Digital 2000 IS",
	ModelPowerShotSD870:    "PowerShot SD870 IS / Digital IXUS 860 IS / IXY Digital 910 IS",
	ModelPowerShotSD890:    "PowerShot SD890 IS / Digital IXUS 970 IS / IXY Digital 820 IS",
	ModelPowerShotSD790:    "PowerShot SD790 IS / Digital IXUS 90 IS / IXY Digital 95 IS",
	ModelPowerShotSD770:    "PowerShot SD770 IS / Digital IXUS 85 IS / IXY Digital 25 IS",
	ModelPowerShotA590:     "PowerShot A590 IS",
	ModelPowerShotA580:     "PowerShot A580",
	ModelPowerShotA470:     "PowerShot A470",
	ModelPowerShotSD1100:   "PowerShot SD1100 IS / Digital IXUS 80 IS / IXY Digital 20 IS",
	ModelPowerShotSX1:      "PowerShot SX1 IS",
	ModelPowerShotSX10:     "PowerShot SX10 IS",
	ModelPowerShotA1000:    "PowerShot A1000 IS",
	ModelPowerShotG10:      "PowerShot G10",
	ModelPowerShotA2000:    "PowerShot A2000 IS",
	ModelPowerShotSX110:    "PowerShot SX110 IS",
	ModelPowerShotSD990:    "PowerShot SD990 IS / Digital IXUS 980 IS / IXY Digital 3000 IS",
	ModelPowerShotSD880:    "PowerShot SD880 IS / Digital IXUS 870 IS / IXY Digital 920 IS",
	ModelPowerShotE1:       "PowerShot E1",
	ModelPowerShotD10:      "PowerShot D10",
	ModelPowerShotSD960:    "PowerShot SD960 IS / Digital IXUS 110 IS / IXY Digital 510 IS",
	ModelPowerShotA2100:    "PowerShot A2100 IS",
	ModelPowerShotA480:     "PowerShot A480",
	ModelPowerShotSX200:    "PowerShot SX200 IS",
	ModelPowerShotSD970:    "PowerShot SD970 IS / Digital IXUS 990 IS / IXY Digital 830 IS",
	ModelPowerShotSD780:    "PowerShot SD780 IS / Digital IXUS 100 IS / IXY Digital 210 IS",
	ModelPowerShotA1100:    "PowerShot A1100 IS",
	ModelPowerShotSD1200:   "PowerShot SD1200 IS / Digital IXUS 95 IS / IXY Digital 110 IS",
	ModelPowerShotG11:      "PowerShot G11",
	ModelPowerShotSX120:    "PowerShot SX120 IS",
	ModelPowerShotSD980:    "PowerShot SD980 IS / Digital IXUS 200 IS / IXY Digital 930 IS",
	ModelPowerShotSD940:    "PowerShot SD940 IS / Digital IXUS 120 IS / IXY Digital 220 IS",
	ModelPowerShotA495:     "PowerShot A495",
	ModelPowerShotA490:     "PowerShot A490",
	ModelPowerShotA3100:    "PowerShot A3100/A3150 IS",
	ModelPowerShotA3000:    "PowerShot A3000 IS",
	ModelPowerShotSD1400:   "PowerShot SD1400 IS / IXUS 130 / IXY 400F",
	ModelPowerShotSD1300:   "PowerShot SD1300 IS / IXUS 105 / IXY 200F",
	ModelPowerShotSD3500:   "PowerShot SD3500 IS / IXUS 210 / IXY 10S",
	ModelPowerShotSX210:    "PowerShot SX210 IS",
	ModelPowerShotSD4000:   "PowerShot SD4000 IS / IXUS 300 HS / IXY 30S",
	ModelPowerShotSD4500:   "PowerShot SD4500 IS / IXUS 1000 HS / IXY 50S",
	ModelPowerShotG12:      "PowerShot G12",
	ModelPowerShotSX30:     "PowerShot SX30 IS",
	ModelPowerShotSX130:    "PowerShot SX130 IS",
	ModelPowerShotS95:      "PowerShot S95",
	ModelPowerShotA3300:    "PowerShot A3300 IS",
	ModelPowerShotA3200:    "PowerShot A3200 IS",
	ModelPowerShotELPH500:  "PowerShot ELPH 500 HS / IXUS 310 HS / IXY 31S",
	ModelPowerShotPro90:    "PowerShot Pro90 IS",
	ModelPowerShotA800:     "PowerShot A800",
	ModelPowerShotELPH100:  "PowerShot ELPH 100 HS / IXUS 115 HS / IXY 210F",
	ModelPowerShotSX230:    "PowerShot SX230 HS",
	ModelPowerShotELPH300:  "PowerShot ELPH 300 HS / IXUS 220 HS / IXY 410F",
	ModelPowerShotA2200:    "PowerShot A2200",
	ModelPowerShotA1200:    "PowerShot A1200",
	ModelPowerShotSX220:    "PowerShot SX220 HS",
	ModelPowerShotG1X:      "PowerShot G1 X",
	ModelPowerShotSX150:    "PowerShot SX150 IS",
	ModelPowerShotELPH510:  "PowerShot ELPH 510 HS / IXUS 1100 HS / IXY 51S",
	ModelPowerShotS100New:  "PowerShot S100 (new)",
	ModelPowerShotSX40:     "PowerShot SX40 HS",
	ModelPowerShotELPH310:  "PowerShot ELPH 310 HS / IXUS 230 HS / IXY 600F",
	ModelIXY32S:            "IXY 32S", // (PowerShot ELPH 500 HS / IXUS 320 HS ??)
	ModelPowerShotA1300:    "PowerShot A1300",
	ModelPowerShotA810:     "PowerShot A810",
	ModelPowerShotELPH320:  "PowerShot ELPH 320 HS / IXUS 240 HS / IXY 420F",
	ModelPowerShotELPH110:  "PowerShot ELPH 110 HS / IXUS 125 HS / IXY 220F",
	ModelPowerShotD20:      "PowerShot D20",
	ModelPowerShotA4000:    "PowerShot A4000 IS",
	ModelPowerShotSX260:    "PowerShot SX260 HS",
	ModelPowerShotSX240:    "PowerShot SX240 HS",
	ModelPowerShotELPH530:  "PowerShot ELPH 530 HS / IXUS 510 HS / IXY 1",
	ModelPowerShotELPH520:  "PowerShot ELPH 520 HS / IXUS 500 HS / IXY 3",
	ModelPowerShotA3400:    "PowerShot A3400 IS",
	ModelPowerShotA2400:    "PowerShot A2400 IS",
	ModelPowerShotA2300:    "PowerShot A2300",
	ModelPowerShotS100V:    "PowerShot S100V",
	ModelPowerShotG15:      "PowerShot G15",
	ModelPowerShotSX50:     "PowerShot SX50 HS",
	ModelPowerShotSX160:    "PowerShot SX160 IS",
	ModelPowerShotS110New:  "PowerShot S110 (new)",
	ModelPowerShotSX500:    "PowerShot SX500 IS",
	ModelPowerShotN:        "PowerShot N",
	ModelIXUS245:           "IXUS 245 HS / IXY 430F",
	ModelPowerShotSX280:    "PowerShot SX280 HS",
	ModelPowerShotSX270:    "PowerShot SX270 HS",
	ModelPowerShotA3500:    "PowerShot A3500 IS",
	ModelPowerShotA2600:    "PowerShot A2600",
	ModelPowerShotSX275:    "PowerShot SX275 HS",
	ModelPowerShotA1400:    "PowerShot A1400",
	ModelPowerShotELPH130:  "PowerShot ELPH 130 IS / IXUS 140 / IXY 110F",
	ModelPowerShotELPH115:  "PowerShot ELPH 115/120 IS / IXUS 132/135 / IXY 90F/100F",
	ModelPowerShotELPH330:  "PowerShot ELPH 330 HS / IXUS 255 HS / IXY 610F",
	ModelPowerShotA2500:    "PowerShot A2500",
	ModelPowerShotG16:      "PowerShot G16",
	ModelPowerShotS120:     "PowerShot S120",
	ModelPowerShotSX170:    "PowerShot SX170 IS",
	ModelPowerShotSX510:    "PowerShot SX510 HS",
	ModelPowerShotS200New:  "PowerShot S200 (new)",
	ModelIXY620F:           "IXY 620F",
	ModelPowerShotN100:     "PowerShot N100",
	ModelPowerShotG1XMark2: "PowerShot G1 X Mark II",
	ModelPowerShotD30:      "PowerShot D30",
	ModelPowerShotSX700:    "PowerShot SX700 HS",
	ModelPowerShotSX600:    "PowerShot SX600 HS",
	ModelPowerShotELPH140:  "PowerShot ELPH 140 IS / IXUS 150 / IXY 130",
	ModelPowerShotELPH135:  "PowerShot ELPH 135 / IXUS 145 / IXY 120",
	ModelPowerShotELPH340:  "PowerShot ELPH 340 HS / IXUS 265 HS / IXY 630",
	ModelPowerShotELPH150:  "PowerShot ELPH 150 IS / IXUS 155 / IXY 140",
	ModelEOSM3:             "EOS M3",
	ModelPowerShotSX60:     "PowerShot SX60 HS",
	ModelPowerShotSX520:    "PowerShot SX520 HS",
	ModelPowerShotSX400:    "PowerShot SX400 IS",
	ModelPowerShotG7X:      "PowerShot G7 X",
	ModelPowerShotN2:       "PowerShot N2",
	ModelPowerShotSX530:    "PowerShot SX530 HS",
	ModelPowerShotSX710:    "PowerShot SX710 HS",
	ModelPowerShotSX610:    "PowerShot SX610 HS",
	ModelEOSM10:            "EOS M10",
	ModelPowerShotG3X:      "PowerShot G3 X",
	ModelPowerShotELPH165:  "PowerShot ELPH 165 HS / IXUS 165 / IXY 160",
	ModelPowerShotELPH160:  "PowerShot ELPH 160 / IXUS 160",
	ModelPowerShotELPH350:  "PowerShot ELPH 350 HS / IXUS 275 HS / IXY 640",
	ModelPowerShotELPH170:  "PowerShot ELPH 170 IS / IXUS 170",
	ModelPowerShotSX410:    "PowerShot SX410 IS",
	ModelPowerShotG9X:      "PowerShot G9 X",
	ModelEOSM5:             "EOS M5",
	ModelPowerShotG5X:      "PowerShot G5 X",
	ModelPowerShotG7XMark2: "PowerShot G7 X Mark II",
	ModelEOSM100:           "EOS M100",
	ModelPowerShotELPH360:  "PowerShot ELPH 360 HS / IXUS 285 HS / IXY 650",
	ModelPowerShotSX540:    "PowerShot SX540 HS",
	ModelPowerShotSX420:    "PowerShot SX420 IS",
	ModelPowerShotELPH190:  "PowerShot ELPH 190 IS / IXUS 180 / IXY 190",
	ModelPowerShotG1:       "PowerShot G1",
	ModelPowerShotELPH180:  "PowerShot ELPH 180 IS / IXUS 175 / IXY 180",
	ModelPowerShotSX720:    "PowerShot SX720 HS",
	ModelPowerShotSX620:    "PowerShot SX620 HS",
	ModelEOSM6:             "EOS M6",
	ModelPowerShotG9XMark2: "PowerShot G9 X Mark II",
	ModelEOSM50:            "EOS M50 / Kiss M",
	ModelPowerShotELPH185:  "PowerShot ELPH 185 / IXUS 185 / IXY 200",
	ModelPowerShotSX430:    "PowerShot SX430 IS",
	ModelPowerShotSX730:    "PowerShot SX730 HS",
	ModelPowerShotG1XMark3: "PowerShot G1 X Mark III",
	ModelPowerShotS100:     "PowerShot S100 / Digital IXUS / IXY Digital",
	ModelPowerShotSX740:    "PowerShot SX740 HS",
	ModelPowerShotG5XMark2: "PowerShot G5 X Mark II",
	ModelPowerShotSX70:     "PowerShot SX70 HS",
	ModelPowerShotG7XMark3: "PowerShot G7 X Mark III",
	ModelEOSM6MarkII:       "EOS M6 Mark II",
	ModelEOSM200:           "EOS M200",
	ModelDC19:              "DC19/DC21/DC22",
	ModelXHA1:              "XH A1",
	ModelHV10:              "HV10",
	ModelMD130:             "MD130/MD140/MD150/MD160/ZR850",
	ModelDC50:              "DC50", // iVIS
	ModelHV20:              "HV20", // iVIS
	ModelDC211:             "DC211",
	ModelHG10:              "HG10",
	ModelHR10:              "HR10", // iVIS
	ModelMD255:             "MD255/ZR950",
	ModelHF11:              "HF11",
	ModelHV30:              "HV30",
	ModelXHA1S:             "XH A1S",
	ModelDC301:             "DC301/DC310/DC311/DC320/DC330",
	ModelFS100:             "FS100",
	ModelHF10:              "HF10",      // iVIS/VIXIA
	ModelHG20:              "HG20/HG21", // VIXIA
	ModelHF21:              "HF21",      // LEGRIA
	ModelHFS11:             "HF S11",    // LEGRIA
	ModelHV40:              "HV40",      // LEGRIA
	ModelDC410:             "DC410/DC411/DC420",
	ModelFS19:              "FS19/FS20/FS21/FS22/FS200",    // LEGRIA
	ModelHF20:              "HF20/HF200",                   // LEGRIA
	ModelHFS10:             "HF S10/S100",                  // LEGRIA/VIXIA
	ModelHFR10:             "HF R10/R16/R17/R18/R100/R106", // LEGRIA/VIXIA
	ModelHFM30:             "HF M30/M31/M36/M300/M306",     // LEGRIA/VIXIA
	ModelHFS20:             "HF S20/S21/S200",              // LEGRIA/VIXIA
	ModelFS31:              "FS31/FS36/FS37/FS300/FS305/FS306/FS307",
	ModelEOSC300:           "EOS C300",
	ModelHFG25:             "HF G25", // LEGRIA
	ModelXC10:              "XC10",
	ModelEOSC200:           "EOS C200",
	ModelEOS1D:             "EOS-1D",
	ModelEOS1DS:            "EOS-1DS",
	ModelEOS1DMarkIII:      "EOS-1D Mark III",
	ModelEOS1DMarkII:       "EOS-1D Mark II",
	ModelEOS1DSMarkII:      "EOS-1Ds Mark II",
	ModelEOS1DSMarkIII:     "EOS-1Ds Mark III",
	ModelEOS1DMarkIIN:      "EOS-1D Mark II N",
	ModelEOS1DX:            "EOS-1D X",
	ModelEOS1DMarkIV:       "EOS-1D Mark IV",
	ModelEOS1DC:            "EOS-1D C",
	ModelEOS1DXMarkII:      "EOS-1D X Mark II",
	ModelEOS5D:             "EOS 5D",
	ModelEOS5DMarkII:       "EOS 5D Mark II",
	ModelEOS5DMarkIII:      "EOS 5D Mark III",
	ModelEOS7D:             "EOS 7D",
	ModelEOS7DMarkII:       "EOS 7D Mark II",
	ModelEOS6D:             "EOS 6D",
	ModelEOS10D:            "EOS 10D",
	ModelEOS20D:            "EOS 20D",
	ModelEOS30D:            "EOS 30D",
	ModelEOS40D:            "EOS 40D",
	ModelEOS50D:            "EOS 50D",
	ModelEOS60D:            "EOS 60D",
	ModelEOS70D:            "EOS 70D",
	ModelEOSRebel:          "EOS Digital Rebel / 300D / Kiss Digital",
	ModelEOSRebelXT:        "EOS Digital Rebel XT / 350D / Kiss Digital N",
	ModelEOSRebelXTi:       "EOS Digital Rebel XTi / 400D / Kiss Digital X",
	ModelEOSRebelXSi:       "EOS Digital Rebel XSi / 450D / Kiss X2",
	ModelEOSRebelT1i:       "EOS Rebel T1i / 500D / Kiss X3",
	ModelEOSRebelT2i:       "EOS Rebel T2i / 550D / Kiss X4",
	ModelEOSRebelT3i:       "EOS Rebel T3i / 600D / Kiss X5",
	ModelEOSRebelT4i:       "EOS Rebel T4i / 650D / Kiss X6i",
	ModelEOSRebelT5i:       "EOS Rebel T5i / 700D / Kiss X7i",
	ModelEOSRebelXS:        "EOS Rebel XS / 1000D / Kiss F",
	ModelEOSRebelT3:        "EOS Rebel T3 / 1100D / Kiss X50",
	ModelEOSRebelT5:        "EOS Rebel T5 / 1200D / Kiss X70 / Hi",
	ModelWFTE1:             "WFT-E1",
	ModelWFTE2:             "WFT-E2",
	ModelWFTE3:             "WFT-E3",
	ModelWFTE4:             "WFT-E4",
	ModelWFTE5:             "WFT-E5",
	ModelWFTE2II:           "WFT-E2 II",
	ModelWFTE4II:           "WFT-E4 II",
	ModelEOSM:              "EOS M",
	ModelEOS80D:            "EOS 80D",
	ModelEOS5DMarkIV:       "EOS 5D Mark IV",
	ModelEOS5DS:            "EOS 5DS",
	ModelEOS5DSR:           "EOS 5DS R",
	ModelEOS6DMarkII:       "EOS 6D Mark II",
	ModelEOS77D:            "EOS 77D / 9000D",
	ModelEOS90D:            "EOS 90D",
	ModelEOSRebelSL1:       "EOS Rebel SL1 / 100D / Kiss X7",
	ModelEOSRebelT6s:       "EOS Rebel T6s / 760D / 8000D",
	ModelEOSRebelT6i:       "EOS Rebel T6i / 750D / Kiss X8i",
	ModelEOSRebelT6:        "EOS Rebel T6 / 1300D / Kiss X80",
	ModelEOSRebelT7i:       "EOS Rebel T7i / 800D / Kiss X9i",
	ModelEOSRebelSL2:       "EOS Rebel SL2 / 200D / Kiss X9",
	ModelEOSRebelT100:      "EOS Rebel T100 / 4000D / 3000D", // 3000D in China
	ModelEOSRebelT7:        "EOS Rebel T7 / 2000D / 1500D / Kiss X90",
	ModelEOSRebelT8i:       "EOS Rebel T8i / 850D / X10i",
	ModelEOSSL3:            "EOS SL3 / 250D / Kiss X10",
	ModelEOSR5:             "EOS R5",
	ModelEOSR:              "EOS R",
	ModelEOSRP:             "EOS RP",
	ModelEOSR3:             "EOS R3",
	ModelEOSR6:             "EOS R6",
	ModelEOSR7:             "EOS R7",
	ModelEOSR10:            "EOS R10",
	ModelEOSR50:            "EOS R50",
	ModelEOSR6MarkII:       "EOS R6 Mark II",
	ModelEOSR8:             "EOS R8",
	ModelEOSR1:             "EOS R1",
	ModelR5MarkII:          "R5 Mark II",
	ModelEOSR100:           "EOS R100",
	ModelEOSM2:             "EOS M2",
	ModelEOSM50MarkII:      "EOS M50 Mark II / Kiss M2",
	ModelPowerShotZOOM:     "PowerShot ZOOM",
	ModelPowerShotV10:      "PowerShot V10",
	ModelEOSD2000C:         "EOS D2000C",
	ModelEOSD6000C:         "EOS D6000C",
	ModelPowerShotV1:       "PowerShot V1",
	ModelEOSR50V:           "EOS R50 V",
}

// String returns the string representation of the CanonModelID value.
//func (id CanonModelID) String() string {
//	if str, ok := canonModelIDMap[id]; ok {
//		return str
//	}
//	return "Unknown"
//}

// IsVideoModelID returns true if the ID is in the video camera range
func IsVideoModelID(id CanonModelID) bool {
	// Check video camera model range (0x4007d673 - 0x4007e1c3)
	return id >= ModelDC19 && id <= ModelEOSC200
}

// IsMirrorlessModelID returns true if the model ID represents a mirrorless camera (EOS R or M series)
func IsMirrorlessModelID(id CanonModelID) bool {
	switch id {
	// R Series
	case ModelEOSR5, ModelEOSR, ModelEOSRP, ModelEOSR3,
		ModelEOSR6, ModelEOSR7, ModelEOSR10, ModelEOSR50,
		ModelEOSR6MarkII, ModelEOSR8, ModelEOSR1,
		ModelR5MarkII, ModelEOSR100:
		return true

	// M Series
	case ModelEOSM, ModelEOSM2, ModelEOSM10, ModelEOSM3,
		ModelEOSM5, ModelEOSM100, ModelEOSM200, ModelEOSM50,
		ModelEOSM50MarkII, ModelEOSM6, ModelEOSM6MarkII:
		return true

	default:
		return false
	}
}

// CanonLensType represents Canon lens type identifiers
// Note: Values are incorrect for EOS 7D images with lenses of type 256 or greater
// Based on Phil Harvey's exiftool
// Updated Dec-5-2024
type CanonLensType uint16

const (
	CanonLens50mm18           CanonLensType = 1   // Canon EF 50mm f/1.8
	CanonLens28mm28           CanonLensType = 2   // Canon EF 28mm f/2.8
	CanonLens135mmSoft        CanonLensType = 3   // Canon EF 135mm f/2.8 Soft
	CanonLens35105mm          CanonLensType = 4   // Canon EF 35-105mm f/3.5-4.5
	CanonLens3570mm           CanonLensType = 5   // Canon EF 35-70mm f/3.5-4.5
	CanonLens2870mm           CanonLensType = 6   // Canon EF 28-70mm f/3.5-4.5
	CanonLens100300mmL        CanonLensType = 7   // Canon EF 100-300mm f/5.6L
	CanonLens100300mm         CanonLensType = 8   // Canon EF 100-300mm f/5.6
	CanonLens70210mm          CanonLensType = 9   // Canon EF 70-210mm f/4
	CanonLens50mmMacro        CanonLensType = 10  // Canon EF 50mm f/2.5 Macro
	CanonLens35mm             CanonLensType = 11  // Canon EF 35mm f/2
	CanonLens15mmFisheye      CanonLensType = 13  // Canon EF 15mm f/2.8 Fisheye
	CanonLens50200mmL         CanonLensType = 14  // Canon EF 50-200mm f/3.5-4.5L
	CanonLens50200mm          CanonLensType = 15  // Canon EF 50-200mm f/3.5-4.5
	CanonLens35135mm          CanonLensType = 16  // Canon EF 35-135mm f/3.5-4.5
	CanonLens3570mmA          CanonLensType = 17  // Canon EF 35-70mm f/3.5-4.5A
	CanonLens2870mmII         CanonLensType = 18  // Canon EF 28-70mm f/3.5-4.5
	CanonLens100200mm         CanonLensType = 20  // Canon EF 100-200mm f/4.5A
	CanonLens80200mmL         CanonLensType = 21  // Canon EF 80-200mm f/2.8L
	CanonLens2035mmL          CanonLensType = 22  // Canon EF 20-35mm f/2.8L
	CanonLens35105mm2         CanonLensType = 23  // Canon EF 35-105mm f/3.5-4.5
	CanonLens3580mmPZ         CanonLensType = 24  // Canon EF 35-80mm f/4-5.6 Power Zoom
	CanonLens3580mmPZ2        CanonLensType = 25  // Canon EF 35-80mm f/4-5.6 Power Zoom
	CanonLens100mmMacro       CanonLensType = 26  // Canon EF 100mm f/2.8 Macro
	CanonLens3580mm           CanonLensType = 27  // Canon EF 35-80mm f/4-5.6
	CanonLens80200mm          CanonLensType = 28  // Canon EF 80-200mm f/4.5-5.6
	CanonLens50mm18II         CanonLensType = 29  // Canon EF 50mm f/1.8 II
	CanonLens35105mm3         CanonLensType = 30  // Canon EF 35-105mm f/4.5-5.6
	CanonLens75300mm          CanonLensType = 31  // Canon EF 75-300mm f/4-5.6
	CanonLens24mm28           CanonLensType = 32  // Canon EF 24mm f/2.8
	CanonLens3580mm2          CanonLensType = 35  // Canon EF 35-80mm f/4-5.6
	CanonLens3876mm           CanonLensType = 36  // Canon EF 38-76mm f/4.5-5.6
	CanonLens3580mm3          CanonLensType = 37  // Canon EF 35-80mm f/4-5.6
	CanonLens80200mm2         CanonLensType = 38  // Canon EF 80-200mm f/4.5-5.6 II
	CanonLens75300mm2         CanonLensType = 39  // Canon EF 75-300mm f/4-5.6
	CanonLens2880mm           CanonLensType = 40  // Canon EF 28-80mm f/3.5-5.6
	CanonLens2890mm           CanonLensType = 41  // Canon EF 28-90mm f/4-5.6
	CanonLens28200mm          CanonLensType = 42  // Canon EF 28-200mm f/3.5-5.6
	CanonLens28105mm          CanonLensType = 43  // Canon EF 28-105mm f/4-5.6
	CanonLens90300mm          CanonLensType = 44  // Canon EF 90-300mm f/4.5-5.6
	CanonLensEFS1855mm        CanonLensType = 45  // Canon EF-S 18-55mm f/3.5-5.6 [II]
	CanonLens2890mm2          CanonLensType = 46  // Canon EF 28-90mm f/4-5.6
	CanonLensEFS1855mmIS      CanonLensType = 48  // Canon EF-S 18-55mm f/3.5-5.6 IS
	CanonLensEFS55250mmIS     CanonLensType = 49  // Canon EF-S 55-250mm f/4-5.6 IS
	CanonLensEFS18200mmIS     CanonLensType = 50  // Canon EF-S 18-200mm f/3.5-5.6 IS
	CanonLensEFS18135mmIS     CanonLensType = 51  // Canon EF-S 18-135mm f/3.5-5.6 IS
	CanonLensEFS1855mmIS2     CanonLensType = 52  // Canon EF-S 18-55mm f/3.5-5.6 IS II
	CanonLensEFS1855mm3       CanonLensType = 53  // Canon EF-S 18-55mm f/3.5-5.6 III
	CanonLensEFS55250mmIS2    CanonLensType = 54  // Canon EF-S 55-250mm f/4-5.6 IS II
	CanonLensTS50mm           CanonLensType = 80  // Canon TS-E 50mm f/2.8L Macro
	CanonLensTS90mm           CanonLensType = 81  // Canon TS-E 90mm f/2.8L Macro
	CanonLensTS135mm          CanonLensType = 82  // Canon TS-E 135mm f/4L Macro
	CanonLensTS17mm           CanonLensType = 94  // Canon TS-E 17mm f/4L
	CanonLensTS24mm2          CanonLensType = 95  // Canon TS-E 24mm f/3.5L II
	CanonLensMP65mm           CanonLensType = 124 // Canon MP-E 65mm f/2.8 1-5x Macro Photo
	CanonLensTS24mm           CanonLensType = 125 // Canon TS-E 24mm f/3.5L
	CanonLensTS45mm           CanonLensType = 126 // Canon TS-E 45mm f/2.8
	CanonLensTS90mm2          CanonLensType = 127 // Canon TS-E 90mm f/2.8
	CanonLens300mmL           CanonLensType = 129 // Canon EF 300mm f/2.8L USM
	CanonLens50mmF10L         CanonLensType = 130 // Canon EF 50mm f/1.0L USM
	CanonLens2880mmL          CanonLensType = 131 // Canon EF 28-80mm f/2.8-4L USM
	CanonLens1200mmL          CanonLensType = 132 // Canon EF 1200mm f/5.6L USM
	CanonLens600mmL           CanonLensType = 134 // Canon EF 600mm f/4L IS USM
	CanonLens200mmL           CanonLensType = 135 // Canon EF 200mm f/1.8L USM
	CanonLens300mmL2          CanonLensType = 136 // Canon EF 300mm f/2.8L USM
	CanonLens85mmL            CanonLensType = 137 // Canon EF 85mm f/1.2L USM
	CanonLens2880mmL2         CanonLensType = 138 // Canon EF 28-80mm f/2.8-4L
	CanonLens400mmL           CanonLensType = 139 // Canon EF 400mm f/2.8L USM
	CanonLens500mmL           CanonLensType = 140 // Canon EF 500mm f/4.5L USM
	CanonLens500mmL2          CanonLensType = 141 // Canon EF 500mm f/4.5L USM
	CanonLens300mmLIS         CanonLensType = 142 // Canon EF 300mm f/2.8L IS USM
	CanonLens500mmLIS         CanonLensType = 143 // Canon EF 500mm f/4L IS USM
	CanonLens35135mmUSM       CanonLensType = 144 // Canon EF 35-135mm f/4-5.6 USM
	CanonLens100300mmUSM      CanonLensType = 145 // Canon EF 100-300mm f/4.5-5.6 USM
	CanonLens70210mmUSM       CanonLensType = 146 // Canon EF 70-210mm f/3.5-4.5 USM
	CanonLens35135mmUSM2      CanonLensType = 147 // Canon EF 35-135mm f/4-5.6 USM
	CanonLens2880mmUSM        CanonLensType = 148 // Canon EF 28-80mm f/3.5-5.6 USM
	CanonLens100mmUSM         CanonLensType = 149 // Canon EF 100mm f/2 USM
	CanonLens14mmL            CanonLensType = 150 // Canon EF 14mm f/2.8L USM
	CanonLens200mmL2          CanonLensType = 151 // Canon EF 200mm f/2.8L USM
	CanonLens300mmLIS2        CanonLensType = 152 // Canon EF 300mm f/4L IS USM
	CanonLens35350mmL         CanonLensType = 153 // Canon EF 35-350mm f/3.5-5.6L USM
	CanonLens20mmUSM          CanonLensType = 154 // Canon EF 20mm f/2.8 USM
	CanonLens85mmUSM          CanonLensType = 155 // Canon EF 85mm f/1.8 USM
	CanonLens28105mmUSM       CanonLensType = 156 // Canon EF 28-105mm f/3.5-4.5 USM
	CanonLens2035mmUSM        CanonLensType = 160 // Canon EF 20-35mm f/3.5-4.5 USM
	CanonLens2870mmL2         CanonLensType = 161 // Canon EF 28-70mm f/2.8L USM
	CanonLens200mmL3          CanonLensType = 162 // Canon EF 200mm f/2.8L USM
	CanonLens300mmLf4         CanonLensType = 163 // Canon EF 300mm f/4L
	CanonLens400mmL2          CanonLensType = 164 // Canon EF 400mm f/5.6L
	CanonLens70200mmL         CanonLensType = 165 // Canon EF 70-200mm f/2.8L USM
	CanonLens70200mmL14x      CanonLensType = 166 // Canon EF 70-200mm f/2.8L USM + 1.4x
	CanonLens70200mmL2x       CanonLensType = 167 // Canon EF 70-200mm f/2.8L USM + 2x
	CanonLens28mmUSM          CanonLensType = 168 // Canon EF 28mm f/1.8 USM
	CanonLens1735mmL          CanonLensType = 169 // Canon EF 17-35mm f/2.8L USM
	CanonLens200mmf28         CanonLensType = 170 // Canon EF 200mm f/2.8L II USM
	CanonLens300mmf4          CanonLensType = 171 // Canon EF 300mm f/4L USM
	CanonLens400mmf56         CanonLensType = 172 // Canon EF 400mm f/5.6L USM
	CanonLens180mmL           CanonLensType = 173 // Canon EF 180mm Macro f/3.5L USM
	CanonLens135mmL           CanonLensType = 174 // Canon EF 135mm f/2L USM
	CanonLens400mmL4          CanonLensType = 175 // Canon EF 400mm f/2.8L USM
	CanonLens2485mmUSM        CanonLensType = 176 // Canon EF 24-85mm f/3.5-4.5 USM
	CanonLens300mmLIS3        CanonLensType = 177 // Canon EF 300mm f/4L IS USM
	CanonLens28135mmIS        CanonLensType = 178 // Canon EF 28-135mm f/3.5-5.6 IS
	CanonLens24mmL            CanonLensType = 179 // Canon EF 24mm f/1.4L USM
	CanonLens35mmL            CanonLensType = 180 // Canon EF 35mm f/1.4L USM
	CanonLens100400mmLIS      CanonLensType = 181 // Canon EF 100-400mm f/4.5-5.6L IS USM + 1.4x
	CanonLens100400mmLIS2     CanonLensType = 182 // Canon EF 100-400mm f/4.5-5.6L IS USM + 2x
	CanonLens100400mmLIS3     CanonLensType = 183 // Canon EF 100-400mm f/4.5-5.6L IS USM
	CanonLens400mmL2x         CanonLensType = 184 // Canon EF 400mm f/2.8L USM + 2x
	CanonLens600mmLIS         CanonLensType = 185 // Canon EF 600mm f/4L IS USM
	CanonLens70200mmL2        CanonLensType = 186 // Canon EF 70-200mm f/4L USM
	CanonLens70200mmLf414x    CanonLensType = 187 // Canon EF 70-200mm f/4L USM + 1.4x
	CanonLens70200mmL2x2      CanonLensType = 188 // Canon EF 70-200mm f/4L USM + 2x
	CanonLens70200mmL28x      CanonLensType = 189 // Canon EF 70-200mm f/4L USM + 2.8x
	CanonLens100mmMacroUSM    CanonLensType = 190 // Canon EF 100mm f/2.8 Macro USM
	CanonLens400mmDOIS        CanonLensType = 191 // Canon EF 400mm f/4 DO IS
	CanonLens3580mmUSM        CanonLensType = 193 // Canon EF 35-80mm f/4-5.6 USM
	CanonLens80200mmUSM       CanonLensType = 194 // Canon EF 80-200mm f/4.5-5.6 USM
	CanonLens35105mmUSM       CanonLensType = 195 // Canon EF 35-105mm f/4.5-5.6 USM
	CanonLens75300mmUSM       CanonLensType = 196 // Canon EF 75-300mm f/4-5.6 USM
	CanonLens75300mmISUSM     CanonLensType = 197 // Canon EF 75-300mm f/4-5.6 IS USM
	CanonLens50mmUSM          CanonLensType = 198 // Canon EF 50mm f/1.4 USM
	CanonLens2880mmUSM2       CanonLensType = 199 // Canon EF 28-80mm f/3.5-5.6 USM
	CanonLens75300mmUSM2      CanonLensType = 200 // Canon EF 75-300mm f/4-5.6 USM
	CanonLens2880mmUSM3       CanonLensType = 201 // Canon EF 28-80mm f/3.5-5.6 USM
	CanonLens2880mmUSM4       CanonLensType = 202 // Canon EF 28-80mm f/3.5-5.6 USM IV
	CanonLens2255mmUSM        CanonLensType = 208 // Canon EF 22-55mm f/4-5.6 USM
	CanonLens55200mm          CanonLensType = 209 // Canon EF 55-200mm f/4.5-5.6
	CanonLens2890mmUSM        CanonLensType = 210 // Canon EF 28-90mm f/4-5.6 USM
	CanonLens28200mmUSM       CanonLensType = 211 // Canon EF 28-200mm f/3.5-5.6 USM
	CanonLens28105mmUSM2      CanonLensType = 212 // Canon EF 28-105mm f/4-5.6 USM
	CanonLens90300mmUSM       CanonLensType = 213 // Canon EF 90-300mm f/4.5-5.6 USM
	CanonLensEFS1855mmUSM     CanonLensType = 214 // Canon EF-S 18-55mm f/3.5-5.6 USM
	CanonLens55200mm2         CanonLensType = 215 // Canon EF 55-200mm f/4.5-5.6 II USM
	CanonLens70200mmLIS       CanonLensType = 224 // Canon EF 70-200mm f/2.8L IS USM
	CanonLens70200mmLIS14x    CanonLensType = 225 // Canon EF 70-200mm f/2.8L IS USM + 1.4x
	CanonLens70200mmLIS2x     CanonLensType = 226 // Canon EF 70-200mm f/2.8L IS USM + 2x
	CanonLens70200mmLIS28x    CanonLensType = 227 // Canon EF 70-200mm f/2.8L IS USM + 2.8x
	CanonLens28105mmUSM3      CanonLensType = 228 // Canon EF 28-105mm f/3.5-4.5 USM
	CanonLens1635mmL          CanonLensType = 229 // Canon EF 16-35mm f/2.8L USM
	CanonLens2470mmL          CanonLensType = 230 // Canon EF 24-70mm f/2.8L USM
	CanonLens1740mmL          CanonLensType = 231 // Canon EF 17-40mm f/4L USM
	CanonLens70300mmDOIS      CanonLensType = 232 // Canon EF 70-300mm f/4.5-5.6 DO IS USM
	CanonLens28300mmLIS       CanonLensType = 233 // Canon EF 28-300mm f/3.5-5.6L IS USM
	CanonLensEFS1785mmIS      CanonLensType = 234 // Canon EF-S 17-85mm f/4-5.6 IS USM
	CanonLensEFS1022mm        CanonLensType = 235 // Canon EF-S 10-22mm f/3.5-4.5 USM
	CanonLensEFS60mmMacro     CanonLensType = 236 // Canon EF-S 60mm f/2.8 Macro USM
	CanonLens24105mmLIS       CanonLensType = 237 // Canon EF 24-105mm f/4L IS USM
	CanonLens70300mmIS        CanonLensType = 238 // Canon EF 70-300mm f/4-5.6 IS USM
	CanonLens85mmL2           CanonLensType = 239 // Canon EF 85mm f/1.2L II USM
	CanonLensEFS1755mmIS      CanonLensType = 240 // Canon EF-S 17-55mm f/2.8 IS USM
	CanonLens50mmL            CanonLensType = 241 // Canon EF 50mm f/1.2L USM
	CanonLens70200mmLIS2      CanonLensType = 242 // Canon EF 70-200mm f/4L IS USM
	CanonLens70200mmLIS214x   CanonLensType = 243 // Canon EF 70-200mm f/4L IS USM + 1.4x
	CanonLens70200mmLIS22x    CanonLensType = 244 // Canon EF 70-200mm f/4L IS USM + 2x
	CanonLens70200mmLIS228x   CanonLensType = 245 // Canon EF 70-200mm f/4L IS USM + 2.8x
	CanonLens1635mmL2         CanonLensType = 246 // Canon EF 16-35mm f/2.8L II USM
	CanonLens14mmL2           CanonLensType = 247 // Canon EF 14mm f/2.8L II USM
	CanonLens200mmL2IS        CanonLensType = 248 // Canon EF 200mm f/2L IS USM
	CanonLens800mmLIS         CanonLensType = 249 // Canon EF 800mm f/5.6L IS USM
	CanonLens24mmL2           CanonLensType = 250 // Canon EF 24mm f/1.4L II USM
	CanonLens70200mmLIS2II    CanonLensType = 251 // Canon EF 70-200mm f/2.8L IS II/III USM
	CanonLens70200mmLIS2II14x CanonLensType = 252 // Canon EF 70-200mm f/2.8L IS II/III USM + 1.4x
	CanonLens70200mmLIS2II2x  CanonLensType = 253 // Canon EF 70-200mm f/2.8L IS II/III USM + 2x
	CanonLens100mmLMacroIS    CanonLensType = 254 // Canon EF 100mm f/2.8L Macro IS USM
	CanonLensEFS1585mmIS      CanonLensType = 488 // Canon EF-S 15-85mm f/3.5-5.6 IS USM
	CanonLens70300mmLIS       CanonLensType = 489 // Canon EF 70-300mm f/4-5.6L IS USM
	CanonLens815mmLFisheye    CanonLensType = 490 // Canon EF 8-15mm f/4L Fisheye USM
	CanonLens300mmLf28IS2     CanonLensType = 491 // Canon EF 300mm f/2.8L IS II USM
	CanonLens400mmLIS2        CanonLensType = 492 // Canon EF 400mm f/2.8L IS II USM
	CanonLens500mmLIS2        CanonLensType = 493 // Canon EF 500mm f/4L IS II USM
	CanonLens24105mmL         CanonLensType = 493 // Canon EF 24-105mm f/4L IS USM
	CanonLens600mmLIS2        CanonLensType = 494 // Canon EF 600mm f/4L IS II USM
	CanonLens2470mmL2         CanonLensType = 495 // Canon EF 24-70mm f/2.8L II USM
	CanonLens200400mmL        CanonLensType = 496 // Canon EF 200-400mm f/4L IS USM
	CanonLens200400mmL14x     CanonLensType = 499 // Canon EF 200-400mm f/4L IS USM + 1.4x
	CanonLens28mmIS           CanonLensType = 502 // Canon EF 28mm f/2.8 IS USM
	CanonLens24mmIS           CanonLensType = 503 // Canon EF 24mm f/2.8 IS USM
	CanonLens2470mmL4IS       CanonLensType = 504 // Canon EF 24-70mm f/4L IS USM
	CanonLens35mmIS           CanonLensType = 505 // Canon EF 35mm f/2 IS USM
	CanonLens400mmDOIS2       CanonLensType = 506 // Canon EF 400mm f/4 DO IS II USM
	CanonLens1635mmL4IS       CanonLensType = 507 // Canon EF 16-35mm f/4L IS USM
	CanonLens1124mmL          CanonLensType = 508 // Canon EF 11-24mm f/4L USM
	CanonLens100400mmL2IS     CanonLensType = 747 // Canon EF 100-400mm f/4.5-5.6L IS II USM
	CanonLens100400mmL2IS14x  CanonLensType = 748 // Canon EF 100-400mm f/4.5-5.6L IS II USM + 1.4x
	CanonLens35mmL2           CanonLensType = 750 // Canon EF 35mm f/1.4L II USM
	CanonLens1635mmL3         CanonLensType = 751 // Canon EF 16-35mm f/2.8L III USM
	CanonLens24105mmL2IS      CanonLensType = 752 // Canon EF 24-105mm f/4L IS II USM
	CanonLens85mmL4IS         CanonLensType = 753 // Canon EF 85mm f/1.4L IS USM
	CanonLens70200mmL4IS2     CanonLensType = 754 // Canon EF 70-200mm f/4L IS II USM
	CanonLens400mmL3IS        CanonLensType = 757 // Canon EF 400mm f/2.8L IS III USM
	CanonLens600mmL3IS        CanonLensType = 758 // Canon EF 600mm f/4L IS III USM
	// STM Lenses (0x1000 range)
	CanonLensEFS18135mmSTM   CanonLensType = 4142 // Canon EF-S 18-135mm f/3.5-5.6 IS STM
	CanonLensEFM1855mmSTM    CanonLensType = 4143 // Canon EF-M 18-55mm f/3.5-5.6 IS STM
	CanonLens40mmSTM         CanonLensType = 4144 // Canon EF 40mm f/2.8 STM
	CanonLensEFM22mmSTM      CanonLensType = 4145 // Canon EF-M 22mm f/2 STM
	CanonLensEFS1855mmSTM    CanonLensType = 4146 // Canon EF-S 18-55mm f/3.5-5.6 IS STM
	CanonLensEFM1122mmSTM    CanonLensType = 4147 // Canon EF-M 11-22mm f/4-5.6 IS STM
	CanonLensEFS55250mmSTM   CanonLensType = 4148 // Canon EF-S 55-250mm f/4-5.6 IS STM
	CanonLensEFM55200mmSTM   CanonLensType = 4149 // Canon EF-M 55-200mm f/4.5-6.3 IS STM
	CanonLensEFS1018mmSTM    CanonLensType = 4150 // Canon EF-S 10-18mm f/4.5-5.6 IS STM
	CanonLens24105mmSTM      CanonLensType = 4152 // Canon EF 24-105mm f/3.5-5.6 IS STM
	CanonLensEFM1545mmSTM    CanonLensType = 4153 // Canon EF-M 15-45mm f/3.5-6.3 IS STM
	CanonLensEFS24mmSTM      CanonLensType = 4154 // Canon EF-S 24mm f/2.8 STM
	CanonLensEFM28mmMacroSTM CanonLensType = 4155 // Canon EF-M 28mm f/3.5 Macro IS STM
	CanonLens50mmSTM         CanonLensType = 4156 // Canon EF 50mm f/1.8 STM
	CanonLensEFM18150mmSTM   CanonLensType = 4157 // Canon EF-M 18-150mm f/3.5-6.3 IS STM
	CanonLensEFS1855mm4STM   CanonLensType = 4158 // Canon EF-S 18-55mm f/4-5.6 IS STM
	CanonLensEFM32mmSTM      CanonLensType = 4159 // Canon EF-M 32mm f/1.4 STM
	CanonLensEFS35mmMacroSTM CanonLensType = 4160 // Canon EF-S 35mm f/2.8 Macro IS STM
	// Nano USM Lenses (0x9000 range)
	CanonLens70300mmNanoUSM CanonLensType = 36910 // Canon EF 70-300mm f/4-5.6 IS II USM
	CanonLens18135mmNanoUSM CanonLensType = 36912 // Canon EF-S 18-135mm f/3.5-5.6 IS USM
	// CN-E Cinema Lenses (0xf000 range)
	CanonLensCNE14mm  CanonLensType = 61491 // Canon CN-E 14mm T3.1 L F
	CanonLensCNE24mm  CanonLensType = 61492 // Canon CN-E 24mm T1.5 L F
	CanonLensCNE85mm  CanonLensType = 61494 // Canon CN-E 85mm T1.3 L F
	CanonLensCNE135mm CanonLensType = 61495 // Canon CN-E 135mm T2.2 L F
	CanonLensCNE35mm  CanonLensType = 61496 // Canon CN-E 35mm T1.5 L F
	CanonRLens        CanonLensType = 61182 // Canon R Lens Type
	CanonLensUnkown   CanonLensType = 65535 // Unknown Lens Type
)

// Map of lens types to their string descriptions
var canonLensTypeMap = map[CanonLensType]string{
	1:   "Canon EF 50mm f/1.8",
	2:   "Canon EF 28mm f/2.8",
	3:   "Canon EF 135mm f/2.8 Soft",
	4:   "Canon EF 35-105mm f/3.5-4.5",
	5:   "Canon EF 35-70mm f/3.5-4.5",
	6:   "Canon EF 28-70mm f/3.5-4.5",
	7:   "Canon EF 100-300mm f/5.6L",
	8:   "Canon EF 100-300mm f/5.6",
	9:   "Canon EF 70-210mm f/4",
	10:  "Canon EF 50mm f/2.5 Macro",
	11:  "Canon EF 35mm f/2",
	13:  "Canon EF 15mm f/2.8 Fisheye",
	14:  "Canon EF 50-200mm f/3.5-4.5L",
	15:  "Canon EF 50-200mm f/3.5-4.5",
	16:  "Canon EF 35-135mm f/3.5-4.5",
	17:  "Canon EF 35-70mm f/3.5-4.5A",
	18:  "Canon EF 28-70mm f/3.5-4.5",
	20:  "Canon EF 100-200mm f/4.5A",
	21:  "Canon EF 80-200mm f/2.8L",
	22:  "Canon EF 20-35mm f/2.8L",
	23:  "Canon EF 35-105mm f/3.5-4.5",
	24:  "Canon EF 35-80mm f/4-5.6 Power Zoom",
	25:  "Canon EF 35-80mm f/4-5.6 Power Zoom",
	26:  "Canon EF 100mm f/2.8 Macro",
	27:  "Canon EF 35-80mm f/4-5.6",
	28:  "Canon EF 80-200mm f/4.5-5.6",
	29:  "Canon EF 50mm f/1.8 II",
	30:  "Canon EF 35-105mm f/4.5-5.6",
	31:  "Canon EF 75-300mm f/4-5.6",
	32:  "Canon EF 24mm f/2.8",
	35:  "Canon EF 35-80mm f/4-5.6",
	36:  "Canon EF 38-76mm f/4.5-5.6",
	37:  "Canon EF 35-80mm f/4-5.6",
	38:  "Canon EF 80-200mm f/4.5-5.6 II",
	39:  "Canon EF 75-300mm f/4-5.6",
	40:  "Canon EF 28-80mm f/3.5-5.6",
	41:  "Canon EF 28-90mm f/4-5.6",
	42:  "Canon EF 28-200mm f/3.5-5.6",
	43:  "Canon EF 28-105mm f/4-5.6",
	44:  "Canon EF 90-300mm f/4.5-5.6",
	45:  "Canon EF-S 18-55mm f/3.5-5.6 [II]",
	46:  "Canon EF 28-90mm f/4-5.6",
	48:  "Canon EF-S 18-55mm f/3.5-5.6 IS",
	49:  "Canon EF-S 55-250mm f/4-5.6 IS",
	50:  "Canon EF-S 18-200mm f/3.5-5.6 IS",
	51:  "Canon EF-S 18-135mm f/3.5-5.6 IS",
	52:  "Canon EF-S 18-55mm f/3.5-5.6 IS II",
	53:  "Canon EF-S 18-55mm f/3.5-5.6 III",
	54:  "Canon EF-S 55-250mm f/4-5.6 IS II",
	80:  "Canon TS-E 50mm f/2.8L Macro",
	81:  "Canon TS-E 90mm f/2.8L Macro",
	82:  "Canon TS-E 135mm f/4L Macro",
	94:  "Canon TS-E 17mm f/4L",
	95:  "Canon TS-E 24mm f/3.5L II",
	124: "Canon MP-E 65mm f/2.8 1-5x Macro Photo",
	125: "Canon TS-E 24mm f/3.5L",
	126: "Canon TS-E 45mm f/2.8",
	127: "Canon TS-E 90mm f/2.8",
	129: "Canon EF 300mm f/2.8L USM",
	130: "Canon EF 50mm f/1.0L USM",
	131: "Canon EF 28-80mm f/2.8-4L USM",
	132: "Canon EF 1200mm f/5.6L USM",
	134: "Canon EF 600mm f/4L IS USM",
	135: "Canon EF 200mm f/1.8L USM",
	136: "Canon EF 300mm f/2.8L USM",
	137: "Canon EF 85mm f/1.2L USM",
	138: "Canon EF 28-80mm f/2.8-4L",
	139: "Canon EF 400mm f/2.8L USM",
	140: "Canon EF 500mm f/4.5L USM",
	141: "Canon EF 500mm f/4.5L USM",
	142: "Canon EF 300mm f/2.8L IS USM",
	143: "Canon EF 500mm f/4L IS USM",
	144: "Canon EF 35-135mm f/4-5.6 USM",
	145: "Canon EF 100-300mm f/4.5-5.6 USM",
	146: "Canon EF 70-210mm f/3.5-4.5 USM",
	147: "Canon EF 35-135mm f/4-5.6 USM",
	148: "Canon EF 28-80mm f/3.5-5.6 USM",
	149: "Canon EF 100mm f/2 USM",
	150: "Canon EF 14mm f/2.8L USM",
	151: "Canon EF 200mm f/2.8L USM",
	152: "Canon EF 300mm f/4L IS USM",
	153: "Canon EF 35-350mm f/3.5-5.6L USM",
	154: "Canon EF 20mm f/2.8 USM",
	155: "Canon EF 85mm f/1.8 USM",
	156: "Canon EF 28-105mm f/3.5-4.5 USM",
	160: "Canon EF 20-35mm f/3.5-4.5 USM",
	161: "Canon EF 28-70mm f/2.8L USM",
	162: "Canon EF 200mm f/2.8L USM",
	163: "Canon EF 300mm f/4L",
	164: "Canon EF 400mm f/5.6L",
	165: "Canon EF 70-200mm f/2.8L USM",
	166: "Canon EF 70-200mm f/2.8L USM + 1.4x",
	167: "Canon EF 70-200mm f/2.8L USM + 2x",
	168: "Canon EF 28mm f/1.8 USM",
	169: "Canon EF 17-35mm f/2.8L USM",
	170: "Canon EF 200mm f/2.8L II USM",
	171: "Canon EF 300mm f/4L USM",
	172: "Canon EF 400mm f/5.6L USM",
	173: "Canon EF 180mm Macro f/3.5L USM",
	174: "Canon EF 135mm f/2L USM",
	175: "Canon EF 400mm f/2.8L USM",
	176: "Canon EF 24-85mm f/3.5-4.5 USM",
	177: "Canon EF 300mm f/4L IS USM",
	178: "Canon EF 28-135mm f/3.5-5.6 IS",
	179: "Canon EF 24mm f/1.4L USM",
	180: "Canon EF 35mm f/1.4L USM",
	181: "Canon EF 100-400mm f/4.5-5.6L IS USM + 1.4x",
	182: "Canon EF 100-400mm f/4.5-5.6L IS USM + 2x",
	183: "Canon EF 100-400mm f/4.5-5.6L IS USM",
	184: "Canon EF 400mm f/2.8L USM + 2x",
	185: "Canon EF 600mm f/4L IS USM",
	186: "Canon EF 70-200mm f/4L USM",
	187: "Canon EF 70-200mm f/4L USM + 1.4x",
	188: "Canon EF 70-200mm f/4L USM + 2x",
	189: "Canon EF 70-200mm f/4L USM + 2.8x",
	190: "Canon EF 100mm f/2.8 Macro USM",
	191: "Canon EF 400mm f/4 DO IS",
	193: "Canon EF 35-80mm f/4-5.6 USM",
	194: "Canon EF 80-200mm f/4.5-5.6 USM",
	195: "Canon EF 35-105mm f/4.5-5.6 USM",
	196: "Canon EF 75-300mm f/4-5.6 USM",
	197: "Canon EF 75-300mm f/4-5.6 IS USM",
	198: "Canon EF 50mm f/1.4 USM",
	199: "Canon EF 28-80mm f/3.5-5.6 USM",
	200: "Canon EF 75-300mm f/4-5.6 USM",
	201: "Canon EF 28-80mm f/3.5-5.6 USM",
	202: "Canon EF 28-80mm f/3.5-5.6 USM IV",
	208: "Canon EF 22-55mm f/4-5.6 USM",
	209: "Canon EF 55-200mm f/4.5-5.6",
	210: "Canon EF 28-90mm f/4-5.6 USM",
	211: "Canon EF 28-200mm f/3.5-5.6 USM",
	212: "Canon EF 28-105mm f/4-5.6 USM",
	213: "Canon EF 90-300mm f/4.5-5.6 USM",
	214: "Canon EF-S 18-55mm f/3.5-5.6 USM",
	215: "Canon EF 55-200mm f/4.5-5.6 II USM",
	224: "Canon EF 70-200mm f/2.8L IS USM",
	225: "Canon EF 70-200mm f/2.8L IS USM + 1.4x",
	226: "Canon EF 70-200mm f/2.8L IS USM + 2x",
	227: "Canon EF 70-200mm f/2.8L IS USM + 2.8x",
	228: "Canon EF 28-105mm f/3.5-4.5 USM",
	229: "Canon EF 16-35mm f/2.8L USM",
	230: "Canon EF 24-70mm f/2.8L USM",
	231: "Canon EF 17-40mm f/4L USM",
	232: "Canon EF 70-300mm f/4.5-5.6 DO IS USM",
	233: "Canon EF 28-300mm f/3.5-5.6L IS USM",
	234: "Canon EF-S 17-85mm f/4-5.6 IS USM",
	235: "Canon EF-S 10-22mm f/3.5-4.5 USM",
	236: "Canon EF-S 60mm f/2.8 Macro USM",
	237: "Canon EF 24-105mm f/4L IS USM",
	238: "Canon EF 70-300mm f/4-5.6 IS USM",
	239: "Canon EF 85mm f/1.2L II USM",
	240: "Canon EF-S 17-55mm f/2.8 IS USM",
	241: "Canon EF 50mm f/1.2L USM",
	242: "Canon EF 70-200mm f/4L IS USM",
	243: "Canon EF 70-200mm f/4L IS USM + 1.4x",
	244: "Canon EF 70-200mm f/4L IS USM + 2x",
	245: "Canon EF 70-200mm f/4L IS USM + 2.8x",
	246: "Canon EF 16-35mm f/2.8L II USM",
	247: "Canon EF 14mm f/2.8L II USM",
	248: "Canon EF 200mm f/2L IS USM",
	249: "Canon EF 800mm f/5.6L IS USM",
	250: "Canon EF 24mm f/1.4L II USM",
	251: "Canon EF 70-200mm f/2.8L IS II/III USM",
	252: "Canon EF 70-200mm f/2.8L IS II/III USM + 1.4x",
	253: "Canon EF 70-200mm f/2.8L IS II/III USM + 2x",
	254: "Canon EF 100mm f/2.8L Macro IS USM",
	489: "Canon EF 70-300mm f/4-5.6L IS USM",
	490: "Canon EF 8-15mm f/4L Fisheye USM",
	491: "Canon EF 300mm f/2.8L IS II USM",
	492: "Canon EF 400mm f/2.8L IS II USM",
	493: "Canon EF 500mm f/4L IS II USM",
	494: "Canon EF 600mm f/4L IS II USM",
	495: "Canon EF 24-70mm f/2.8L II USM",
	496: "Canon EF 200-400mm f/4L IS USM",
	499: "Canon EF 200-400mm f/4L IS USM + 1.4x",
	502: "Canon EF 28mm f/2.8 IS USM",
	503: "Canon EF 24mm f/2.8 IS USM",
	504: "Canon EF 24-70mm f/4L IS USM",
	505: "Canon EF 35mm f/2 IS USM",
	506: "Canon EF 400mm f/4 DO IS II USM",
	507: "Canon EF 16-35mm f/4L IS USM",
	508: "Canon EF 11-24mm f/4L USM",
	747: "Canon EF 100-400mm f/4.5-5.6L IS II USM",
	748: "Canon EF 100-400mm f/4.5-5.6L IS II USM + 1.4x",
	750: "Canon EF 35mm f/1.4L II USM",
	751: "Canon EF 16-35mm f/2.8L III USM",
	752: "Canon EF 24-105mm f/4L IS II USM",
	753: "Canon EF 85mm f/1.4L IS USM",
	754: "Canon EF 70-200mm f/4L IS II USM",
	757: "Canon EF 400mm f/2.8L IS III USM",
	758: "Canon EF 600mm f/4L IS III USM",

	// STM Lenses
	4142: "Canon EF-S 18-135mm f/3.5-5.6 IS STM",
	4143: "Canon EF-M 18-55mm f/3.5-5.6 IS STM",
	4144: "Canon EF 40mm f/2.8 STM",
	4145: "Canon EF-M 22mm f/2 STM",
	4146: "Canon EF-S 18-55mm f/3.5-5.6 IS STM",
	4147: "Canon EF-M 11-22mm f/4-5.6 IS STM",
	4148: "Canon EF-S 55-250mm f/4-5.6 IS STM",
	4149: "Canon EF-M 55-200mm f/4.5-6.3 IS STM",
	4150: "Canon EF-S 10-18mm f/4.5-5.6 IS STM",
	4152: "Canon EF 24-105mm f/3.5-5.6 IS STM",
	4153: "Canon EF-M 15-45mm f/3.5-6.3 IS STM",
	4154: "Canon EF-S 24mm f/2.8 STM",
	4155: "Canon EF-M 28mm f/3.5 Macro IS STM",
	4156: "Canon EF 50mm f/1.8 STM",
	4157: "Canon EF-M 18-150mm f/3.5-6.3 IS STM",
	4158: "Canon EF-S 18-55mm f/4-5.6 IS STM",
	4159: "Canon EF-M 32mm f/1.4 STM",
	4160: "Canon EF-S 35mm f/2.8 Macro IS STM",

	// Nano USM Lenses
	36910: "Canon EF 70-300mm f/4-5.6 IS II USM",
	36912: "Canon EF-S 18-135mm f/3.5-5.6 IS USM",

	// CN-E Cinema Lenses
	61491: "Canon CN-E 14mm T3.1 L F",
	61492: "Canon CN-E 24mm T1.5 L F",
	61494: "Canon CN-E 85mm T1.3 L F",
	61495: "Canon CN-E 135mm T2.2 L F",
	61496: "Canon CN-E 35mm T1.5 L F",
	61182: "Canon R Lens Type",
	65535: "Unknown Lens Type",
}

// String returns a human-readable representation of the lens type
func (lt CanonLensType) String() string {
	if str, ok := canonLensTypeMap[lt]; ok {
		return str
	}
	return fmt.Sprintf("Unknown lens type: %d", lt)
}

// Valid checks if the lens type is known
func (lt CanonLensType) Valid() error {
	if _, ok := canonLensTypeMap[lt]; !ok {
		return fmt.Errorf("invalid lens type: %d", lt)
	}
	return nil
}

// CanonRFLensType represents Canon RF lens identifiers
type CanonRFLensType uint16

const (
	RFLens50mmF12L          CanonRFLensType = 0  // Canon RF 50mm F1.2L USM
	RFLens24105mmF4L        CanonRFLensType = 1  // Canon RF 24-105mm F4L IS USM
	RFLens2870mmF2L         CanonRFLensType = 2  // Canon RF 28-70mm F2L USM
	RFLens35mmF18Macro      CanonRFLensType = 3  // Canon RF 35mm F1.8 MACRO IS STM
	RFLens85mmF12L          CanonRFLensType = 4  // Canon RF 85mm F1.2L USM
	RFLens85mmF12LDS        CanonRFLensType = 5  // Canon RF 85mm F1.2L USM DS
	RFLens2470mmF28L        CanonRFLensType = 6  // Canon RF 24-70mm F2.8L IS USM
	RFLens1535mmF28L        CanonRFLensType = 7  // Canon RF 15-35mm F2.8L IS USM
	RFLens24240mmF463       CanonRFLensType = 8  // Canon RF 24-240mm F4-6.3 IS USM
	RFLens70200mmF28L       CanonRFLensType = 9  // Canon RF 70-200mm F2.8L IS USM
	RFLens85mmF2Macro       CanonRFLensType = 10 // Canon RF 85mm F2 MACRO IS STM
	RFLens600mmF11          CanonRFLensType = 11 // Canon RF 600mm F11 IS STM
	RFLens600mmF11x14       CanonRFLensType = 12 // Canon RF 600mm F11 IS STM + RF1.4x
	RFLens600mmF11x2        CanonRFLensType = 13 // Canon RF 600mm F11 IS STM + RF2x
	RFLens800mmF11          CanonRFLensType = 14 // Canon RF 800mm F11 IS STM
	RFLens800mmF11x14       CanonRFLensType = 15 // Canon RF 800mm F11 IS STM + RF1.4x
	RFLens800mmF11x2        CanonRFLensType = 16 // Canon RF 800mm F11 IS STM + RF2x
	RFLens24105mmF471       CanonRFLensType = 17 // Canon RF 24-105mm F4-7.1 IS STM
	RFLens100500mmF4571L    CanonRFLensType = 18 // Canon RF 100-500mm F4.5-7.1L IS USM
	RFLens100500mmF4571Lx14 CanonRFLensType = 19 // Canon RF 100-500mm F4.5-7.1L IS USM + RF1.4x
	RFLens100500mmF4571Lx2  CanonRFLensType = 20 // Canon RF 100-500mm F4.5-7.1L IS USM + RF2x
	RFLens70200mmF4L        CanonRFLensType = 21 // Canon RF 70-200mm F4L IS USM
	RFLens100mmF28LMacro    CanonRFLensType = 22 // Canon RF 100mm F2.8L MACRO IS USM
	RFLens50mmF18STM        CanonRFLensType = 23 // Canon RF 50mm F1.8 STM
	RFLens1435mmF4L         CanonRFLensType = 24 // Canon RF 14-35mm F4L IS USM
	RFLensS1845mmF4563      CanonRFLensType = 25 // Canon RF-S 18-45mm F4.5-6.3 IS STM
	RFLens100400mmF568      CanonRFLensType = 26 // Canon RF 100-400mm F5.6-8 IS USM
	RFLens100400mmF568x14   CanonRFLensType = 27 // Canon RF 100-400mm F5.6-8 IS USM + RF1.4x
	RFLens100400mmF568x2    CanonRFLensType = 28 // Canon RF 100-400mm F5.6-8 IS USM + RF2x
	RFLensS18150mmF3563     CanonRFLensType = 29 // Canon RF-S 18-150mm F3.5-6.3 IS STM
	RFLens24mmF18Macro      CanonRFLensType = 30 // Canon RF 24mm F1.8 MACRO IS STM
	RFLens16mmF28STM        CanonRFLensType = 31 // Canon RF 16mm F2.8 STM
	RFLens400mmF28L         CanonRFLensType = 32 // Canon RF 400mm F2.8L IS USM
	RFLens400mmF28Lx14      CanonRFLensType = 33 // Canon RF 400mm F2.8L IS USM + RF1.4x
	RFLens400mmF28Lx2       CanonRFLensType = 34 // Canon RF 400mm F2.8L IS USM + RF2x
	RFLens600mmF4L          CanonRFLensType = 35 // Canon RF 600mm F4L IS USM
	RFLens600mmF4Lx14       CanonRFLensType = 36 // Canon RF 600mm F4L IS USM + RF1.4x
	RFLens600mmF4Lx2        CanonRFLensType = 37 // Canon RF 600mm F4L IS USM + RF2x
	RFLens800mmF56L         CanonRFLensType = 38 // Canon RF 800mm F5.6L IS USM
	RFLens800mmF56Lx14      CanonRFLensType = 39 // Canon RF 800mm F5.6L IS USM + RF1.4x
	RFLens800mmF56Lx2       CanonRFLensType = 40 // Canon RF 800mm F5.6L IS USM + RF2x
	RFLens1200mmF8L         CanonRFLensType = 41 // Canon RF 1200mm F8L IS USM
	RFLens1200mmF8Lx14      CanonRFLensType = 42 // Canon RF 1200mm F8L IS USM + RF1.4x
	RFLens1200mmF8Lx2       CanonRFLensType = 43 // Canon RF 1200mm F8L IS USM + RF2x
	RFLens52mmF28LFisheye   CanonRFLensType = 44 // Canon RF 5.2mm F2.8L Dual Fisheye 3D VR
	RFLens1530mmF4563       CanonRFLensType = 45 // Canon RF 15-30mm F4.5-6.3 IS STM
	RFLens135mmF18L         CanonRFLensType = 46 // Canon RF 135mm F1.8L IS USM
	RFLens2450mmF4563       CanonRFLensType = 47 // Canon RF 24-50mm F4.5-6.3 IS STM
	RFLensS55210mmF571      CanonRFLensType = 48 // Canon RF-S 55-210mm F5-7.1 IS STM
	RFLens100300mmF28L      CanonRFLensType = 49 // Canon RF 100-300mm F2.8L IS USM
	RFLens100300mmF28Lx14   CanonRFLensType = 50 // Canon RF 100-300mm F2.8L IS USM + RF1.4x
	RFLens100300mmF28Lx2    CanonRFLensType = 51 // Canon RF 100-300mm F2.8L IS USM + RF2x
	RFLens1020mmF4L         CanonRFLensType = 52 // Canon RF 10-20mm F4L IS STM
	RFLens28mmF28STM        CanonRFLensType = 53 // Canon RF 28mm F2.8 STM
	RFLens24105mmF28LZ      CanonRFLensType = 54 // Canon RF 24-105mm F2.8L IS USM Z
	RFLensS1018mmF4563      CanonRFLensType = 55 // Canon RF-S 10-18mm F4.5-6.3 IS STM
	RFLens35mmF14LVCM       CanonRFLensType = 56 // Canon RF 35mm F1.4L VCM
	RFLens70200mmF28LZ      CanonRFLensType = 57 // Canon RF 70-200mm F2.8L IS USM Z
	RFLens50mmF14LVCM       CanonRFLensType = 58 // Canon RF 50mm F1.4L VCM
	RFLens24mmF14LVCM       CanonRFLensType = 59 // Canon RF 24mm F1.4L VCM
)

// Concatenated string of all RF lens names
const strCanonRFLensTypeString = "Canon RF 50mm F1.2L USMCanon RF 24-105mm F4L IS USMCanon RF 28-70mm F2L USM" +
	"Canon RF 35mm F1.8 MACRO IS STMCanon RF 85mm F1.2L USMCanon RF 85mm F1.2L USM DS" +
	"Canon RF 24-70mm F2.8L IS USMCanon RF 15-35mm F2.8L IS USMCanon RF 24-240mm F4-6.3 IS USM" +
	"Canon RF 70-200mm F2.8L IS USMCanon RF 85mm F2 MACRO IS STMCanon RF 600mm F11 IS STM" +
	"Canon RF 600mm F11 IS STM + RF1.4xCanon RF 600mm F11 IS STM + RF2xCanon RF 800mm F11 IS STM" +
	"Canon RF 800mm F11 IS STM + RF1.4xCanon RF 800mm F11 IS STM + RF2xCanon RF 24-105mm F4-7.1 IS STM" +
	"Canon RF 100-500mm F4.5-7.1L IS USMCanon RF 100-500mm F4.5-7.1L IS USM + RF1.4x" +
	"Canon RF 100-500mm F4.5-7.1L IS USM + RF2xCanon RF 70-200mm F4L IS USM" +
	"Canon RF 100mm F2.8L MACRO IS USMCanon RF 50mm F1.8 STMCanon RF 14-35mm F4L IS USM" +
	"Canon RF-S 18-45mm F4.5-6.3 IS STMCanon RF 100-400mm F5.6-8 IS USM" +
	"Canon RF 100-400mm F5.6-8 IS USM + RF1.4xCanon RF 100-400mm F5.6-8 IS USM + RF2x" +
	"Canon RF-S 18-150mm F3.5-6.3 IS STMCanon RF 24mm F1.8 MACRO IS STMCanon RF 16mm F2.8 STM" +
	"Canon RF 400mm F2.8L IS USMCanon RF 400mm F2.8L IS USM + RF1.4xCanon RF 400mm F2.8L IS USM + RF2x" +
	"Canon RF 600mm F4L IS USMCanon RF 600mm F4L IS USM + RF1.4xCanon RF 600mm F4L IS USM + RF2x" +
	"Canon RF 800mm F5.6L IS USMCanon RF 800mm F5.6L IS USM + RF1.4xCanon RF 800mm F5.6L IS USM + RF2x" +
	"Canon RF 1200mm F8L IS USMCanon RF 1200mm F8L IS USM + RF1.4xCanon RF 1200mm F8L IS USM + RF2x" +
	"Canon RF 5.2mm F2.8L Dual Fisheye 3D VRCanon RF 15-30mm F4.5-6.3 IS STMCanon RF 135mm F1.8L IS USM" +
	"Canon RF 24-50mm F4.5-6.3 IS STMCanon RF-S 55-210mm F5-7.1 IS STMCanon RF 100-300mm F2.8L IS USM" +
	"Canon RF 100-300mm F2.8L IS USM + RF1.4xCanon RF 100-300mm F2.8L IS USM + RF2x" +
	"Canon RF 10-20mm F4L IS STMCanon RF 28mm F2.8 STMCanon RF 24-105mm F2.8L IS USM Z" +
	"Canon RF-S 10-18mm F4.5-6.3 IS STMCanon RF 35mm F1.4L VCMCanon RF 70-200mm F2.8L IS USM Z" +
	"Canon RF 50mm F1.4L VCMCanon RF 24mm F1.4L VCM"

// Indices into the concatenated string for efficient slicing
var strCanonRFLensTypeDist = []int{
	0, 23, 51, 75, 105, 129, 156, 185, 214, 245, 274, 302, 329, 361, 391, 417,
	449, 479, 509, 543, 584, 625, 654, 687, 710, 737, 771, 802, 841, 880, 916,
	947, 971, 1000, 1033, 1066, 1099, 1127, 1160, 1193, 1226, 1259, 1292, 1325,
	1363, 1395, 1423, 1454, 1485, 1516, 1557, 1598, 1629, 1656, 1689, 1725,
	1756, 1787, 1820, 1847,
}

// String returns a human-readable representation of the RF lens type
func (rt CanonRFLensType) String() string {
	if int(rt) < len(strCanonRFLensTypeDist)-1 {
		return strCanonRFLensTypeString[strCanonRFLensTypeDist[rt]:strCanonRFLensTypeDist[rt+1]]
	}
	return fmt.Sprintf("Unknown RF lens type: %d", rt)
}
