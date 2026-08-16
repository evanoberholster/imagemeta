package sony

import "strings"

// SonyModelIDSourceFile is the ExifTool source used to generate this table.
const SonyModelIDSourceFile = "Image/ExifTool/Sony.pm"

// SonyModelIDSourceVersion is the ExifTool version used to generate this table.
const SonyModelIDSourceVersion = "3.86"

//go:generate stringer -type=SonyCameraModel -linecomment -output=models_string.go

// SonyCameraModel is the numeric camera model identifier stored in Sony
// maker-note tag 0xb001 (int16u).  ExifTool assigns these IDs to
// uniquely identify each Sony camera body.
type SonyCameraModel uint16

const (
	SonyCameraUnknown     SonyCameraModel = 0   // Unknown
	SonyCameraDSCR1       SonyCameraModel = 2   // DSC-R1
	SonyCameraDSLRA100    SonyCameraModel = 256 // DSLR-A100
	SonyCameraDSLRA900    SonyCameraModel = 257 // DSLR-A900
	SonyCameraDSLRA700    SonyCameraModel = 258 // DSLR-A700
	SonyCameraDSLRA200    SonyCameraModel = 259 // DSLR-A200
	SonyCameraDSLRA350    SonyCameraModel = 260 // DSLR-A350
	SonyCameraDSLRA300    SonyCameraModel = 261 // DSLR-A300
	SonyCameraDSLRA900C   SonyCameraModel = 262 // DSLR-A900 (APS-C mode)
	SonyCameraDSLRA380    SonyCameraModel = 263 // DSLR-A380/A390
	SonyCameraDSLRA330    SonyCameraModel = 264 // DSLR-A330
	SonyCameraDSLRA230    SonyCameraModel = 265 // DSLR-A230
	SonyCameraDSLRA290    SonyCameraModel = 266 // DSLR-A290
	SonyCameraDSLRA850    SonyCameraModel = 269 // DSLR-A850
	SonyCameraDSLRA850C   SonyCameraModel = 270 // DSLR-A850 (APS-C mode)
	SonyCameraDSLRA550    SonyCameraModel = 273 // DSLR-A550
	SonyCameraDSLRA500    SonyCameraModel = 274 // DSLR-A500
	SonyCameraDSLRA450    SonyCameraModel = 275 // DSLR-A450
	SonyCameraNEX5        SonyCameraModel = 278 // NEX-5
	SonyCameraNEX3        SonyCameraModel = 279 // NEX-3
	SonyCameraSLTA33      SonyCameraModel = 280 // SLT-A33
	SonyCameraSLTA55      SonyCameraModel = 281 // SLT-A55/SLT-A55V
	SonyCameraDSLRA560    SonyCameraModel = 282 // DSLR-A560
	SonyCameraDSLRA580    SonyCameraModel = 283 // DSLR-A580
	SonyCameraNEXC3       SonyCameraModel = 284 // NEX-C3
	SonyCameraSLTA35      SonyCameraModel = 285 // SLT-A35
	SonyCameraSLTA65      SonyCameraModel = 286 // SLT-A65/SLT-A65V
	SonyCameraSLTA77      SonyCameraModel = 287 // SLT-A77/SLT-A77V
	SonyCameraNEX5N       SonyCameraModel = 288 // NEX-5N
	SonyCameraNEX7        SonyCameraModel = 289 // NEX-7
	SonyCameraNEXVG20E    SonyCameraModel = 290 // NEX-VG20E
	SonyCameraSLTA37      SonyCameraModel = 291 // SLT-A37
	SonyCameraSLTA57      SonyCameraModel = 292 // SLT-A57
	SonyCameraNEXF3       SonyCameraModel = 293 // NEX-F3
	SonyCameraSLTA99      SonyCameraModel = 294 // SLT-A99/SLT-A99V
	SonyCameraNEX6        SonyCameraModel = 295 // NEX-6
	SonyCameraNEX5R       SonyCameraModel = 296 // NEX-5R
	SonyCameraDSCRX100    SonyCameraModel = 297 // DSC-RX100
	SonyCameraDSCRX1      SonyCameraModel = 298 // DSC-RX1
	SonyCameraNEXVG900    SonyCameraModel = 299 // NEX-VG900
	SonyCameraNEXVG30E    SonyCameraModel = 300 // NEX-VG30E
	SonyCameraILCE3000    SonyCameraModel = 302 // ILCE-3000/ILCE-3500
	SonyCameraSLTA58      SonyCameraModel = 303 // SLT-A58
	SonyCameraNEX3N       SonyCameraModel = 305 // NEX-3N
	SonyCameraILCE7       SonyCameraModel = 306 // ILCE-7
	SonyCameraNEX5T       SonyCameraModel = 307 // NEX-5T
	SonyCameraDSCRX100M2  SonyCameraModel = 308 // DSC-RX100M2
	SonyCameraDSCRX10     SonyCameraModel = 309 // DSC-RX10
	SonyCameraDSCRX1R     SonyCameraModel = 310 // DSC-RX1R
	SonyCameraILCE7R      SonyCameraModel = 311 // ILCE-7R
	SonyCameraILCE6000    SonyCameraModel = 312 // ILCE-6000
	SonyCameraILCE5000    SonyCameraModel = 313 // ILCE-5000
	SonyCameraDSCRX100M3  SonyCameraModel = 317 // DSC-RX100M3
	SonyCameraILCE7S      SonyCameraModel = 318 // ILCE-7S
	SonyCameraILCA77M2    SonyCameraModel = 319 // ILCA-77M2
	SonyCameraILCE5100    SonyCameraModel = 339 // ILCE-5100
	SonyCameraILCE7M2     SonyCameraModel = 340 // ILCE-7M2
	SonyCameraDSCRX100M4  SonyCameraModel = 341 // DSC-RX100M4
	SonyCameraDSCRX10M2   SonyCameraModel = 342 // DSC-RX10M2
	SonyCameraDSCRX1RM2   SonyCameraModel = 344 // DSC-RX1RM2
	SonyCameraILCEQX1     SonyCameraModel = 346 // ILCE-QX1
	SonyCameraILCE7RM2    SonyCameraModel = 347 // ILCE-7RM2
	SonyCameraILCE7SM2    SonyCameraModel = 350 // ILCE-7SM2
	SonyCameraILCA68      SonyCameraModel = 353 // ILCA-68
	SonyCameraILCA99M2    SonyCameraModel = 354 // ILCA-99M2
	SonyCameraDSCRX10M3   SonyCameraModel = 355 // DSC-RX10M3
	SonyCameraDSCRX100M5  SonyCameraModel = 356 // DSC-RX100M5
	SonyCameraILCE6300    SonyCameraModel = 357 // ILCE-6300
	SonyCameraILCE9       SonyCameraModel = 358 // ILCE-9
	SonyCameraILCE6500    SonyCameraModel = 360 // ILCE-6500
	SonyCameraILCE7RM3    SonyCameraModel = 362 // ILCE-7RM3
	SonyCameraILCE7M3     SonyCameraModel = 363 // ILCE-7M3
	SonyCameraDSCRX0      SonyCameraModel = 364 // DSC-RX0
	SonyCameraDSCRX10M4   SonyCameraModel = 365 // DSC-RX10M4
	SonyCameraDSCRX100M6  SonyCameraModel = 366 // DSC-RX100M6
	SonyCameraDSCHX99     SonyCameraModel = 367 // DSC-HX99
	SonyCameraDSCRX100M5A SonyCameraModel = 369 // DSC-RX100M5A
	SonyCameraILCE6400    SonyCameraModel = 371 // ILCE-6400
	SonyCameraDSCRX0M2    SonyCameraModel = 372 // DSC-RX0M2
	SonyCameraDSCHX95     SonyCameraModel = 373 // DSC-HX95
	SonyCameraDSCRX100M7  SonyCameraModel = 374 // DSC-RX100M7
	SonyCameraILCE7RM4    SonyCameraModel = 375 // ILCE-7RM4
	SonyCameraILCE9M2     SonyCameraModel = 376 // ILCE-9M2
	SonyCameraILCE6600    SonyCameraModel = 378 // ILCE-6600
	SonyCameraILCE6100    SonyCameraModel = 379 // ILCE-6100
	SonyCameraZV1         SonyCameraModel = 380 // ZV-1
	SonyCameraILCE7C      SonyCameraModel = 381 // ILCE-7C
	SonyCameraZVE10       SonyCameraModel = 382 // ZV-E10
	SonyCameraILCE7SM3    SonyCameraModel = 383 // ILCE-7SM3
	SonyCameraILCE1       SonyCameraModel = 384 // ILCE-1
	SonyCameraILMEFX3     SonyCameraModel = 385 // ILME-FX3
	SonyCameraILCE7RM3A   SonyCameraModel = 386 // ILCE-7RM3A
	SonyCameraILCE7RM4A   SonyCameraModel = 387 // ILCE-7RM4A
	SonyCameraILCE7M4     SonyCameraModel = 388 // ILCE-7M4
	SonyCameraZV1F        SonyCameraModel = 389 // ZV-1F
	SonyCameraILCE7RM5    SonyCameraModel = 390 // ILCE-7RM5
	SonyCameraILMEFX30    SonyCameraModel = 391 // ILME-FX30
	SonyCameraILCE9M3     SonyCameraModel = 392 // ILCE-9M3
	SonyCameraZVE1        SonyCameraModel = 393 // ZV-E1
	SonyCameraILCE6700    SonyCameraModel = 394 // ILCE-6700
	SonyCameraZV1M2       SonyCameraModel = 395 // ZV-1M2
	SonyCameraILCE7CR     SonyCameraModel = 396 // ILCE-7CR
	SonyCameraILCE7CM2    SonyCameraModel = 397 // ILCE-7CM2
	SonyCameraILXLR1      SonyCameraModel = 398 // ILX-LR1
	SonyCameraZVE10M2     SonyCameraModel = 399 // ZV-E10M2
	SonyCameraILCE1M2     SonyCameraModel = 400 // ILCE-1M2
	SonyCameraDSCRX1RM3   SonyCameraModel = 401 // DSC-RX1RM3
	SonyCameraILCE6400A   SonyCameraModel = 402 // ILCE-6400A
	SonyCameraILCE6100A   SonyCameraModel = 403 // ILCE-6100A
	SonyCameraDSCRX100M7A SonyCameraModel = 404 // DSC-RX100M7A
	SonyCameraILMEFX2     SonyCameraModel = 406 // ILME-FX2
	SonyCameraILCE7M5     SonyCameraModel = 407 // ILCE-7M5
	SonyCameraZV1A        SonyCameraModel = 408 // ZV-1A
)

// modelIDMapping maps camera model string prefixes to SonyCameraModel values.
var modelIDMapping = []struct {
	ModelPrefix string
	ID          SonyCameraModel
}{
	{"DSLR-A100", SonyCameraDSLRA100},
	{"DSLR-A900", SonyCameraDSLRA900},
	{"DSLR-A700", SonyCameraDSLRA700},
	{"DSLR-A200", SonyCameraDSLRA200},
	{"DSLR-A350", SonyCameraDSLRA350},
	{"DSLR-A300", SonyCameraDSLRA300},
	{"DSLR-A390", SonyCameraDSLRA380},
	{"DSLR-A380", SonyCameraDSLRA380},
	{"DSLR-A330", SonyCameraDSLRA330},
	{"DSLR-A230", SonyCameraDSLRA230},
	{"DSLR-A290", SonyCameraDSLRA290},
	{"DSLR-A850", SonyCameraDSLRA850},
	{"DSLR-A550", SonyCameraDSLRA550},
	{"DSLR-A500", SonyCameraDSLRA500},
	{"DSLR-A450", SonyCameraDSLRA450},
	{"DSLR-A560", SonyCameraDSLRA560},
	{"DSLR-A580", SonyCameraDSLRA580},
	{"SLT-A58", SonyCameraSLTA58},
	{"SLT-A37", SonyCameraSLTA37},
	{"SLT-A57", SonyCameraSLTA57},
	{"SLT-A35", SonyCameraSLTA35},
	{"SLT-A33", SonyCameraSLTA33},
	{"SLT-A55", SonyCameraSLTA55},
	{"SLT-A65", SonyCameraSLTA65},
	{"SLT-A77", SonyCameraSLTA77},
	{"SLT-A99", SonyCameraSLTA99},
	{"NEX-VG900", SonyCameraNEXVG900},
	{"NEX-VG30", SonyCameraNEXVG30E},
	{"NEX-VG20", SonyCameraNEXVG20E},
	{"NEX-C3", SonyCameraNEXC3},
	{"NEX-F3", SonyCameraNEXF3},
	{"NEX-3N", SonyCameraNEX3N},
	{"NEX-5N", SonyCameraNEX5N},
	{"NEX-5R", SonyCameraNEX5R},
	{"NEX-5T", SonyCameraNEX5T},
	{"NEX-5", SonyCameraNEX5},
	{"NEX-3", SonyCameraNEX3},
	{"NEX-6", SonyCameraNEX6},
	{"NEX-7", SonyCameraNEX7},
	{"ILCA-99M2", SonyCameraILCA99M2},
	{"ILCA-77M2", SonyCameraILCA77M2},
	{"ILCA-68", SonyCameraILCA68},
	{"ILCE-7RM5", SonyCameraILCE7RM5},
	{"ILCE-7RM4A", SonyCameraILCE7RM4A},
	{"ILCE-7RM4", SonyCameraILCE7RM4},
	{"ILCE-7RM3A", SonyCameraILCE7RM3A},
	{"ILCE-7RM3", SonyCameraILCE7RM3},
	{"ILCE-7RM2", SonyCameraILCE7RM2},
	{"ILCE-7SM3", SonyCameraILCE7SM3},
	{"ILCE-7SM2", SonyCameraILCE7SM2},
	{"ILCE-7S", SonyCameraILCE7S},
	{"ILCE-7M5", SonyCameraILCE7M5},
	{"ILCE-7M4", SonyCameraILCE7M4},
	{"ILCE-7M3", SonyCameraILCE7M3},
	{"ILCE-7M2", SonyCameraILCE7M2},
	{"ILCE-7CR", SonyCameraILCE7CR},
	{"ILCE-7CM2", SonyCameraILCE7CM2},
	{"ILCE-7C", SonyCameraILCE7C},
	{"ILCE-7R", SonyCameraILCE7R},
	{"ILCE-7", SonyCameraILCE7},
	{"ILCE-9M3", SonyCameraILCE9M3},
	{"ILCE-9M2", SonyCameraILCE9M2},
	{"ILCE-9", SonyCameraILCE9},
	{"ILCE-1M2", SonyCameraILCE1M2},
	{"ILCE-1", SonyCameraILCE1},
	{"ILCE-6700", SonyCameraILCE6700},
	{"ILCE-6600", SonyCameraILCE6600},
	{"ILCE-6500", SonyCameraILCE6500},
	{"ILCE-6400", SonyCameraILCE6400},
	{"ILCE-6300", SonyCameraILCE6300},
	{"ILCE-6100", SonyCameraILCE6100},
	{"ILCE-6000", SonyCameraILCE6000},
	{"ILCE-5100", SonyCameraILCE5100},
	{"ILCE-5000", SonyCameraILCE5000},
	{"ILCE-3500", SonyCameraILCE3000},
	{"ILCE-3000", SonyCameraILCE3000},
	{"ILCE-QX1", SonyCameraILCEQX1},
	{"ILME-FX30", SonyCameraILMEFX30},
	{"ILME-FX3", SonyCameraILMEFX3},
	{"ILME-FX2", SonyCameraILMEFX2},
	{"ILX-LR1", SonyCameraILXLR1},
	{"ZV-E10M2", SonyCameraZVE10M2},
	{"ZV-E10", SonyCameraZVE10},
	{"ZV-E1", SonyCameraZVE1},
	{"ZV-1M2", SonyCameraZV1M2},
	{"ZV-1F", SonyCameraZV1F},
	{"ZV-1", SonyCameraZV1},
	{"DSC-RX1RM3", SonyCameraDSCRX1RM3},
	{"DSC-RX1RM2", SonyCameraDSCRX1RM2},
	{"DSC-RX1R", SonyCameraDSCRX1R},
	{"DSC-RX10M4", SonyCameraDSCRX10M4},
	{"DSC-RX10M3", SonyCameraDSCRX10M3},
	{"DSC-RX10M2", SonyCameraDSCRX10M2},
	{"DSC-RX10", SonyCameraDSCRX10},
	{"DSC-RX100M7A", SonyCameraDSCRX100M7A},
	{"DSC-RX100M7", SonyCameraDSCRX100M7},
	{"DSC-RX100M6", SonyCameraDSCRX100M6},
	{"DSC-RX100M5A", SonyCameraDSCRX100M5A},
	{"DSC-RX100M5", SonyCameraDSCRX100M5},
	{"DSC-RX100M4", SonyCameraDSCRX100M4},
	{"DSC-RX100M3", SonyCameraDSCRX100M3},
	{"DSC-RX100M2", SonyCameraDSCRX100M2},
	{"DSC-RX100", SonyCameraDSCRX100},
	{"DSC-RX0M2", SonyCameraDSCRX0M2},
	{"DSC-RX0", SonyCameraDSCRX0},
	{"DSC-RX1", SonyCameraDSCRX1},
	{"DSC-R1", SonyCameraDSCR1},
	{"DSC-HX99", SonyCameraDSCHX99},
	{"DSC-HX95", SonyCameraDSCHX95},
}

// ModelIDFromModel maps a camera model string to a SonyCameraModel ID value.
// It uses the EXIF Model string to identify the camera body.  Returns 0
// (SonyCameraUnknown) when the model is not recognized.
func ModelIDFromModel(model string) uint16 {
	model = strings.ToUpper(strings.TrimSpace(model))
	for _, m := range modelIDMapping {
		if strings.HasPrefix(model, strings.ToUpper(m.ModelPrefix)) {
			return uint16(m.ID)
		}
	}
	return 0
}
