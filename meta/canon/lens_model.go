package canon

//go:generate stringer -type=CanonRFLensType -linecomment -output=lens_model_string.go

// CanonRFLensType is Canon MakerNote FileInfo[0x3d] RF lens model id.
//
// Source: ExifTool Canon FileInfo table (RFLensType).
// https://exiftool.org/TagNames/Canon.html#FileInfo
//
// Values below intentionally match ExifTool's predetermined numeric codes.
type CanonRFLensType uint16

const (
	CanonRFLensTypeNA CanonRFLensType = 0 // n/a

	CanonRFLensTypeRF50mmF12LUSM                  CanonRFLensType = 257 // Canon RF 50mm F1.2L USM
	CanonRFLensTypeRF24105mmF4LISUSM              CanonRFLensType = 258 // Canon RF 24-105mm F4L IS USM
	CanonRFLensTypeRF2870mmF2LUSM                 CanonRFLensType = 259 // Canon RF 28-70mm F2L USM
	CanonRFLensTypeRF35mmF18MacroISSTM            CanonRFLensType = 260 // Canon RF 35mm F1.8 MACRO IS STM
	CanonRFLensTypeRF85mmF12LUSM                  CanonRFLensType = 261 // Canon RF 85mm F1.2L USM
	CanonRFLensTypeRF85mmF12LUSMDS                CanonRFLensType = 262 // Canon RF 85mm F1.2L USM DS
	CanonRFLensTypeRF2470mmF28LISUSM              CanonRFLensType = 263 // Canon RF 24-70mm F2.8L IS USM
	CanonRFLensTypeRF1535mmF28LISUSM              CanonRFLensType = 264 // Canon RF 15-35mm F2.8L IS USM
	CanonRFLensTypeRF24240mmF463ISUSM             CanonRFLensType = 265 // Canon RF 24-240mm F4-6.3 IS USM
	CanonRFLensTypeRF70200mmF28LISUSM             CanonRFLensType = 266 // Canon RF 70-200mm F2.8L IS USM
	CanonRFLensTypeRF85mmF2MacroISSTM             CanonRFLensType = 267 // Canon RF 85mm F2 MACRO IS STM
	CanonRFLensTypeRF600mmF11ISSTM                CanonRFLensType = 268 // Canon RF 600mm F11 IS STM
	CanonRFLensTypeRF600mmF11ISSTMPlusRF14x       CanonRFLensType = 269 // Canon RF 600mm F11 IS STM + RF1.4x
	CanonRFLensTypeRF600mmF11ISSTMPlusRF2x        CanonRFLensType = 270 // Canon RF 600mm F11 IS STM + RF2x
	CanonRFLensTypeRF800mmF11ISSTM                CanonRFLensType = 271 // Canon RF 800mm F11 IS STM
	CanonRFLensTypeRF800mmF11ISSTMPlusRF14x       CanonRFLensType = 272 // Canon RF 800mm F11 IS STM + RF1.4x
	CanonRFLensTypeRF800mmF11ISSTMPlusRF2x        CanonRFLensType = 273 // Canon RF 800mm F11 IS STM + RF2x
	CanonRFLensTypeRF24105mmF471ISSTM             CanonRFLensType = 274 // Canon RF 24-105mm F4-7.1 IS STM
	CanonRFLensTypeRF100500mmF4571LISUSM          CanonRFLensType = 275 // Canon RF 100-500mm F4.5-7.1L IS USM
	CanonRFLensTypeRF100500mmF4571LISUSMPlusRF14x CanonRFLensType = 276 // Canon RF 100-500mm F4.5-7.1L IS USM + RF1.4x
	CanonRFLensTypeRF100500mmF4571LISUSMPlusRF2x  CanonRFLensType = 277 // Canon RF 100-500mm F4.5-7.1L IS USM + RF2x
	CanonRFLensTypeRF70200mmF4LISUSM              CanonRFLensType = 278 // Canon RF 70-200mm F4L IS USM
	CanonRFLensTypeRF100mmF28LMacroISUSM          CanonRFLensType = 279 // Canon RF 100mm F2.8L MACRO IS USM
	CanonRFLensTypeRF50mmF18STM                   CanonRFLensType = 280 // Canon RF 50mm F1.8 STM
	CanonRFLensTypeRF1435mmF4LISUSM               CanonRFLensType = 281 // Canon RF 14-35mm F4L IS USM
	CanonRFLensTypeRFS1845mmF4563ISSTM            CanonRFLensType = 282 // Canon RF-S 18-45mm F4.5-6.3 IS STM
	CanonRFLensTypeRF100400mmF568ISUSM            CanonRFLensType = 283 // Canon RF 100-400mm F5.6-8 IS USM
	CanonRFLensTypeRF100400mmF568ISUSMPlusRF14x   CanonRFLensType = 284 // Canon RF 100-400mm F5.6-8 IS USM + RF1.4x
	CanonRFLensTypeRF100400mmF568ISUSMPlusRF2x    CanonRFLensType = 285 // Canon RF 100-400mm F5.6-8 IS USM + RF2x
	CanonRFLensTypeRFS18150mmF3563ISSTM           CanonRFLensType = 286 // Canon RF-S 18-150mm F3.5-6.3 IS STM
	CanonRFLensTypeRF24mmF18MacroISSTM            CanonRFLensType = 287 // Canon RF 24mm F1.8 MACRO IS STM
	CanonRFLensTypeRF16mmF28STM                   CanonRFLensType = 288 // Canon RF 16mm F2.8 STM
	CanonRFLensTypeRF400mmF28LISUSM               CanonRFLensType = 289 // Canon RF 400mm F2.8L IS USM
	CanonRFLensTypeRF400mmF28LISUSMPlusRF14x      CanonRFLensType = 290 // Canon RF 400mm F2.8L IS USM + RF1.4x
	CanonRFLensTypeRF400mmF28LISUSMPlusRF2x       CanonRFLensType = 291 // Canon RF 400mm F2.8L IS USM + RF2x
	CanonRFLensTypeRF600mmF4LISUSM                CanonRFLensType = 292 // Canon RF 600mm F4L IS USM
	CanonRFLensTypeRF600mmF4LISUSMPlusRF14x       CanonRFLensType = 293 // Canon RF 600mm F4L IS USM + RF1.4x
	CanonRFLensTypeRF600mmF4LISUSMPlusRF2x        CanonRFLensType = 294 // Canon RF 600mm F4L IS USM + RF2x
	CanonRFLensTypeRF800mmF56LISUSM               CanonRFLensType = 295 // Canon RF 800mm F5.6L IS USM
	CanonRFLensTypeRF800mmF56LISUSMPlusRF14x      CanonRFLensType = 296 // Canon RF 800mm F5.6L IS USM + RF1.4x
	CanonRFLensTypeRF800mmF56LISUSMPlusRF2x       CanonRFLensType = 297 // Canon RF 800mm F5.6L IS USM + RF2x
	CanonRFLensTypeRF1200mmF8LISUSM               CanonRFLensType = 298 // Canon RF 1200mm F8L IS USM
	CanonRFLensTypeRF1200mmF8LISUSMPlusRF14x      CanonRFLensType = 299 // Canon RF 1200mm F8L IS USM + RF1.4x
	CanonRFLensTypeRF1200mmF8LISUSMPlusRF2x       CanonRFLensType = 300 // Canon RF 1200mm F8L IS USM + RF2x
	CanonRFLensTypeRF52mmF28LDualFisheye3DVR      CanonRFLensType = 301 // Canon RF 5.2mm F2.8L Dual Fisheye 3D VR
	CanonRFLensTypeRF1530mmF4563ISSTM             CanonRFLensType = 302 // Canon RF 15-30mm F4.5-6.3 IS STM
	CanonRFLensTypeRF135mmF18LISUSM               CanonRFLensType = 303 // Canon RF 135mm F1.8 L IS USM
	CanonRFLensTypeRF2450mmF4563ISSTM             CanonRFLensType = 304 // Canon RF 24-50mm F4.5-6.3 IS STM
	CanonRFLensTypeRFS55210mmF571ISSTM            CanonRFLensType = 305 // Canon RF-S 55-210mm F5-7.1 IS STM
	CanonRFLensTypeRF100300mmF28LISUSM            CanonRFLensType = 306 // Canon RF 100-300mm F2.8L IS USM
	CanonRFLensTypeRF100300mmF28LISUSMPlusRF14x   CanonRFLensType = 307 // Canon RF 100-300mm F2.8L IS USM + RF1.4x
	CanonRFLensTypeRF100300mmF28LISUSMPlusRF2x    CanonRFLensType = 308 // Canon RF 100-300mm F2.8L IS USM + RF2x
	CanonRFLensTypeRF200800mmF639ISUSM            CanonRFLensType = 309 // Canon RF 200-800mm F6.3-9 IS USM
	CanonRFLensTypeRF200800mmF639ISUSMPlusRF14x   CanonRFLensType = 310 // Canon RF 200-800mm F6.3-9 IS USM + RF1.4x
	CanonRFLensTypeRF200800mmF639ISUSMPlusRF2x    CanonRFLensType = 311 // Canon RF 200-800mm F6.3-9 IS USM + RF2x
	CanonRFLensTypeRF1020mmF4LISSTM               CanonRFLensType = 312 // Canon RF 10-20mm F4 L IS STM
	CanonRFLensTypeRF28mmF28STM                   CanonRFLensType = 313 // Canon RF 28mm F2.8 STM
	CanonRFLensTypeRF24105mmF28LISUSMZ            CanonRFLensType = 314 // Canon RF 24-105mm F2.8 L IS USM Z
	CanonRFLensTypeRFS1018mmF4563ISSTM            CanonRFLensType = 315 // Canon RF-S 10-18mm F4.5-6.3 IS STM
	CanonRFLensTypeRF35mmF14LVCM                  CanonRFLensType = 316 // Canon RF 35mm F1.4 L VCM
	CanonRFLensTypeRFS39mmF35STMDualFisheye       CanonRFLensType = 317 // Canon RF-S 3.9mm F3.5 STM DUAL FISHEYE
	CanonRFLensTypeRF2870mmF28ISSTM               CanonRFLensType = 318 // Canon RF 28-70mm F2.8 IS STM
	CanonRFLensTypeRF70200mmF28LISUSMZ            CanonRFLensType = 319 // Canon RF 70-200mm F2.8 L IS USM Z
	CanonRFLensTypeRF70200mmF28LISUSMZPlusRF14x   CanonRFLensType = 320 // Canon RF 70-200mm F2.8 L IS USM Z + RF1.4x
	CanonRFLensTypeRF70200mmF28LISUSMZPlusRF2x    CanonRFLensType = 321 // Canon RF 70-200mm F2.8 L IS USM Z + RF2x
	CanonRFLensTypeRF1628mmF28ISSTM               CanonRFLensType = 323 // Canon RF 16-28mm F2.8 IS STM
	CanonRFLensTypeRFS1430mmF463ISSTMPZ           CanonRFLensType = 324 // Canon RF-S 14-30mm F4-6.3 IS STM PZ
	CanonRFLensTypeRF50mmF14LVCM                  CanonRFLensType = 325 // Canon RF 50mm F1.4 L VCM
	CanonRFLensTypeRF24mmF14LVCM                  CanonRFLensType = 326 // Canon RF 24mm F1.4 L VCM
	CanonRFLensTypeRF20mmF14LVCM                  CanonRFLensType = 327 // Canon RF 20mm F1.4 L VCM
	CanonRFLensTypeRF85mmF14LVCM                  CanonRFLensType = 328 // Canon RF 85mm F1.4 L VCM
	CanonRFLensTypeRF45mmF12STM                   CanonRFLensType = 330 // Canon RF 45mm F1.2 STM
	CanonRFLensTypeRF714mmF2835LFisheyeSTM        CanonRFLensType = 331 // Canon RF 7-14mm F2.8-3.5 L FISHEYE STM
	CanonRFLensTypeRF14mmF14LVCM                  CanonRFLensType = 332 // Canon RF 14mm F1.4 L VCM
)
