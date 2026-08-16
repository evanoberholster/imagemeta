package canon

import "strconv"

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

// CanonLensTypeSourceFile is the ExifTool source used to generate this table. "Image/ExifTool/Canon.pm"
// CanonLensTypeSourceVersion is the ExifTool version used to generate this table. "13.52"

// CanonLensType
type CanonLensType uint16

const (
	CanonLensUnknown                        CanonLensType = 0
	CanonLensEF50mmF18                      CanonLensType = 1     // Canon EF 50mm f/1.8
	CanonLensEF28mmF28                      CanonLensType = 2     // Canon EF 28mm f/2.8
	CanonLensEF135mmF28Soft                 CanonLensType = 3     // Canon EF 135mm f/2.8 Soft
	CanonLensEF35105mmF3545                 CanonLensType = 4     // Canon EF 35-105mm f/3.5-4.5
	CanonLensEF3570mmF3545                  CanonLensType = 5     // Canon EF 35-70mm f/3.5-4.5
	CanonLensEF2870mmF3545                  CanonLensType = 6     // Canon EF 28-70mm f/3.5-4.5
	CanonLensEF100300mmF56L                 CanonLensType = 7     // Canon EF 100-300mm f/5.6L
	CanonLensEF100300mmF56                  CanonLensType = 8     // Canon EF 100-300mm f/5.6
	CanonLensEF70210mmF4                    CanonLensType = 9     // Canon EF 70-210mm f/4
	CanonLensEF50mmF25Macro                 CanonLensType = 10    // Canon EF 50mm f/2.5 Macro
	CanonLensEF35mmF2                       CanonLensType = 11    // Canon EF 35mm f/2
	CanonLensEF15mmF28Fisheye               CanonLensType = 13    // Canon EF 15mm f/2.8 Fisheye
	CanonLensEF50200mmF3545L                CanonLensType = 14    // Canon EF 50-200mm f/3.5-4.5L
	CanonLensEF50200mmF3545                 CanonLensType = 15    // Canon EF 50-200mm f/3.5-4.5
	CanonLensEF35135mmF3545                 CanonLensType = 16    // Canon EF 35-135mm f/3.5-4.5
	CanonLensEF3570mmF3545A                 CanonLensType = 17    // Canon EF 35-70mm f/3.5-4.5A
	CanonLensEF2870mmF3545ID18              CanonLensType = 18    // Canon EF 28-70mm f/3.5-4.5
	CanonLensEF100200mmF45A                 CanonLensType = 20    // Canon EF 100-200mm f/4.5A
	CanonLensEF80200mmF28L                  CanonLensType = 21    // Canon EF 80-200mm f/2.8L
	CanonLensEF2035mmF28L                   CanonLensType = 22    // Canon EF 20-35mm f/2.8L
	CanonLensEF35105mmF3545ID23             CanonLensType = 23    // Canon EF 35-105mm f/3.5-4.5
	CanonLensEF3580mmF456PowerZoom          CanonLensType = 24    // Canon EF 35-80mm f/4-5.6 Power Zoom
	CanonLensEF3580mmF456PowerZoomID25      CanonLensType = 25    // Canon EF 35-80mm f/4-5.6 Power Zoom
	CanonLensEF100mmF28Macro                CanonLensType = 26    // Canon EF 100mm f/2.8 Macro
	CanonLensEF3580mmF456                   CanonLensType = 27    // Canon EF 35-80mm f/4-5.6
	CanonLensEF80200mmF4556                 CanonLensType = 28    // Canon EF 80-200mm f/4.5-5.6
	CanonLensEF50mmF18II                    CanonLensType = 29    // Canon EF 50mm f/1.8 II
	CanonLensEF35105mmF4556                 CanonLensType = 30    // Canon EF 35-105mm f/4.5-5.6
	CanonLensEF75300mmF456                  CanonLensType = 31    // Canon EF 75-300mm f/4-5.6
	CanonLensEF24mmF28                      CanonLensType = 32    // Canon EF 24mm f/2.8
	CanonLensEF3580mmF456ID35               CanonLensType = 35    // Canon EF 35-80mm f/4-5.6
	CanonLensEF3876mmF4556                  CanonLensType = 36    // Canon EF 38-76mm f/4.5-5.6
	CanonLensEF3580mmF456ID37               CanonLensType = 37    // Canon EF 35-80mm f/4-5.6
	CanonLensEF80200mmF4556II               CanonLensType = 38    // Canon EF 80-200mm f/4.5-5.6 II
	CanonLensEF75300mmF456ID39              CanonLensType = 39    // Canon EF 75-300mm f/4-5.6
	CanonLensEF2880mmF3556                  CanonLensType = 40    // Canon EF 28-80mm f/3.5-5.6
	CanonLensEF2890mmF456                   CanonLensType = 41    // Canon EF 28-90mm f/4-5.6
	CanonLensEF28200mmF3556                 CanonLensType = 42    // Canon EF 28-200mm f/3.5-5.6
	CanonLensEF28105mmF456                  CanonLensType = 43    // Canon EF 28-105mm f/4-5.6
	CanonLensEF90300mmF4556                 CanonLensType = 44    // Canon EF 90-300mm f/4.5-5.6
	CanonLensEFS1855mmF3556II               CanonLensType = 45    // Canon EF-S 18-55mm f/3.5-5.6 [II]
	CanonLensEF2890mmF456ID46               CanonLensType = 46    // Canon EF 28-90mm f/4-5.6
	CanonLensEFS1855mmF3556IS               CanonLensType = 48    // Canon EF-S 18-55mm f/3.5-5.6 IS
	CanonLensEFS55250mmF456IS               CanonLensType = 49    // Canon EF-S 55-250mm f/4-5.6 IS
	CanonLensEFS18200mmF3556IS              CanonLensType = 50    // Canon EF-S 18-200mm f/3.5-5.6 IS
	CanonLensEFS18135mmF3556IS              CanonLensType = 51    // Canon EF-S 18-135mm f/3.5-5.6 IS
	CanonLensEFS1855mmF3556ISII             CanonLensType = 52    // Canon EF-S 18-55mm f/3.5-5.6 IS II
	CanonLensEFS1855mmF3556III              CanonLensType = 53    // Canon EF-S 18-55mm f/3.5-5.6 III
	CanonLensEFS55250mmF456ISII             CanonLensType = 54    // Canon EF-S 55-250mm f/4-5.6 IS II
	CanonLensTSE50mmF28LMacro               CanonLensType = 80    // Canon TS-E 50mm f/2.8L Macro
	CanonLensTSE90mmF28LMacro               CanonLensType = 81    // Canon TS-E 90mm f/2.8L Macro
	CanonLensTSE135mmF4LMacro               CanonLensType = 82    // Canon TS-E 135mm f/4L Macro
	CanonLensTSE17mmF4L                     CanonLensType = 94    // Canon TS-E 17mm f/4L
	CanonLensTSE24mmF35LII                  CanonLensType = 95    // Canon TS-E 24mm f/3.5L II
	CanonLensMPE65mmF2815xMacroPhoto        CanonLensType = 124   // Canon MP-E 65mm f/2.8 1-5x Macro Photo
	CanonLensTSE24mmF35L                    CanonLensType = 125   // Canon TS-E 24mm f/3.5L
	CanonLensTSE45mmF28                     CanonLensType = 126   // Canon TS-E 45mm f/2.8
	CanonLensTSE90mmF28                     CanonLensType = 127   // Canon TS-E 90mm f/2.8
	CanonLensEF300mmF28LUSM                 CanonLensType = 129   // Canon EF 300mm f/2.8L USM
	CanonLensEF50mmF10LUSM                  CanonLensType = 130   // Canon EF 50mm f/1.0L USM
	CanonLensEF2880mmF284LUSM               CanonLensType = 131   // Canon EF 28-80mm f/2.8-4L USM
	CanonLensEF1200mmF56LUSM                CanonLensType = 132   // Canon EF 1200mm f/5.6L USM
	CanonLensEF600mmF4LISUSM                CanonLensType = 134   // Canon EF 600mm f/4L IS USM
	CanonLensEF200mmF18LUSM                 CanonLensType = 135   // Canon EF 200mm f/1.8L USM
	CanonLensEF300mmF28LUSMID136            CanonLensType = 136   // Canon EF 300mm f/2.8L USM
	CanonLensEF85mmF12LUSM                  CanonLensType = 137   // Canon EF 85mm f/1.2L USM
	CanonLensEF2880mmF284L                  CanonLensType = 138   // Canon EF 28-80mm f/2.8-4L
	CanonLensEF400mmF28LUSM                 CanonLensType = 139   // Canon EF 400mm f/2.8L USM
	CanonLensEF500mmF45LUSM                 CanonLensType = 140   // Canon EF 500mm f/4.5L USM
	CanonLensEF500mmF45LUSMID141            CanonLensType = 141   // Canon EF 500mm f/4.5L USM
	CanonLensEF300mmF28LISUSM               CanonLensType = 142   // Canon EF 300mm f/2.8L IS USM
	CanonLensEF500mmF4LISUSM                CanonLensType = 143   // Canon EF 500mm f/4L IS USM
	CanonLensEF35135mmF456USM               CanonLensType = 144   // Canon EF 35-135mm f/4-5.6 USM
	CanonLensEF100300mmF4556USM             CanonLensType = 145   // Canon EF 100-300mm f/4.5-5.6 USM
	CanonLensEF70210mmF3545USM              CanonLensType = 146   // Canon EF 70-210mm f/3.5-4.5 USM
	CanonLensEF35135mmF456USMID147          CanonLensType = 147   // Canon EF 35-135mm f/4-5.6 USM
	CanonLensEF2880mmF3556USM               CanonLensType = 148   // Canon EF 28-80mm f/3.5-5.6 USM
	CanonLensEF100mmF2USM                   CanonLensType = 149   // Canon EF 100mm f/2 USM
	CanonLensEF14mmF28LUSM                  CanonLensType = 150   // Canon EF 14mm f/2.8L USM
	CanonLensEF200mmF28LUSM                 CanonLensType = 151   // Canon EF 200mm f/2.8L USM
	CanonLensEF300mmF4LISUSM                CanonLensType = 152   // Canon EF 300mm f/4L IS USM
	CanonLensEF35350mmF3556LUSM             CanonLensType = 153   // Canon EF 35-350mm f/3.5-5.6L USM
	CanonLensEF20mmF28USM                   CanonLensType = 154   // Canon EF 20mm f/2.8 USM
	CanonLensEF85mmF18USM                   CanonLensType = 155   // Canon EF 85mm f/1.8 USM
	CanonLensEF28105mmF3545USM              CanonLensType = 156   // Canon EF 28-105mm f/3.5-4.5 USM
	CanonLensEF2035mmF3545USM               CanonLensType = 160   // Canon EF 20-35mm f/3.5-4.5 USM
	CanonLensEF2870mmF28LUSM                CanonLensType = 161   // Canon EF 28-70mm f/2.8L USM
	CanonLensEF200mmF28LUSMID162            CanonLensType = 162   // Canon EF 200mm f/2.8L USM
	CanonLensEF300mmF4L                     CanonLensType = 163   // Canon EF 300mm f/4L
	CanonLensEF400mmF56L                    CanonLensType = 164   // Canon EF 400mm f/5.6L
	CanonLensEF70200mmF28LUSM               CanonLensType = 165   // Canon EF 70-200mm f/2.8L USM
	CanonLensEF70200mmF28LUSMPlus14x        CanonLensType = 166   // Canon EF 70-200mm f/2.8L USM + 1.4x
	CanonLensEF70200mmF28LUSMPlus2x         CanonLensType = 167   // Canon EF 70-200mm f/2.8L USM + 2x
	CanonLensEF28mmF18USM                   CanonLensType = 168   // Canon EF 28mm f/1.8 USM
	CanonLensEF1735mmF28LUSM                CanonLensType = 169   // Canon EF 17-35mm f/2.8L USM
	CanonLensEF200mmF28LIIUSM               CanonLensType = 170   // Canon EF 200mm f/2.8L II USM
	CanonLensEF300mmF4LUSM                  CanonLensType = 171   // Canon EF 300mm f/4L USM
	CanonLensEF400mmF56LUSM                 CanonLensType = 172   // Canon EF 400mm f/5.6L USM
	CanonLensEF180mmMacroF35LUSM            CanonLensType = 173   // Canon EF 180mm Macro f/3.5L USM
	CanonLensEF135mmF2LUSM                  CanonLensType = 174   // Canon EF 135mm f/2L USM
	CanonLensEF400mmF28LUSMID175            CanonLensType = 175   // Canon EF 400mm f/2.8L USM
	CanonLensEF2485mmF3545USM               CanonLensType = 176   // Canon EF 24-85mm f/3.5-4.5 USM
	CanonLensEF300mmF4LISUSMID177           CanonLensType = 177   // Canon EF 300mm f/4L IS USM
	CanonLensEF28135mmF3556IS               CanonLensType = 178   // Canon EF 28-135mm f/3.5-5.6 IS
	CanonLensEF24mmF14LUSM                  CanonLensType = 179   // Canon EF 24mm f/1.4L USM
	CanonLensEF35mmF14LUSM                  CanonLensType = 180   // Canon EF 35mm f/1.4L USM
	CanonLensEF100400mmF4556LISUSMPlus14x   CanonLensType = 181   // Canon EF 100-400mm f/4.5-5.6L IS USM + 1.4x
	CanonLensEF100400mmF4556LISUSMPlus2x    CanonLensType = 182   // Canon EF 100-400mm f/4.5-5.6L IS USM + 2x
	CanonLensEF100400mmF4556LISUSM          CanonLensType = 183   // Canon EF 100-400mm f/4.5-5.6L IS USM
	CanonLensEF400mmF28LUSMPlus2x           CanonLensType = 184   // Canon EF 400mm f/2.8L USM + 2x
	CanonLensEF600mmF4LISUSMID185           CanonLensType = 185   // Canon EF 600mm f/4L IS USM
	CanonLensEF70200mmF4LUSM                CanonLensType = 186   // Canon EF 70-200mm f/4L USM
	CanonLensEF70200mmF4LUSMPlus14x         CanonLensType = 187   // Canon EF 70-200mm f/4L USM + 1.4x
	CanonLensEF70200mmF4LUSMPlus2x          CanonLensType = 188   // Canon EF 70-200mm f/4L USM + 2x
	CanonLensEF70200mmF4LUSMPlus28x         CanonLensType = 189   // Canon EF 70-200mm f/4L USM + 2.8x
	CanonLensEF100mmF28MacroUSM             CanonLensType = 190   // Canon EF 100mm f/2.8 Macro USM
	CanonLensEF400mmF4DOIS                  CanonLensType = 191   // Canon EF 400mm f/4 DO IS
	CanonLensEF3580mmF456USM                CanonLensType = 193   // Canon EF 35-80mm f/4-5.6 USM
	CanonLensEF80200mmF4556USM              CanonLensType = 194   // Canon EF 80-200mm f/4.5-5.6 USM
	CanonLensEF35105mmF4556USM              CanonLensType = 195   // Canon EF 35-105mm f/4.5-5.6 USM
	CanonLensEF75300mmF456USM               CanonLensType = 196   // Canon EF 75-300mm f/4-5.6 USM
	CanonLensEF75300mmF456ISUSM             CanonLensType = 197   // Canon EF 75-300mm f/4-5.6 IS USM
	CanonLensEF50mmF14USM                   CanonLensType = 198   // Canon EF 50mm f/1.4 USM
	CanonLensEF2880mmF3556USMID199          CanonLensType = 199   // Canon EF 28-80mm f/3.5-5.6 USM
	CanonLensEF75300mmF456USMID200          CanonLensType = 200   // Canon EF 75-300mm f/4-5.6 USM
	CanonLensEF2880mmF3556USMID201          CanonLensType = 201   // Canon EF 28-80mm f/3.5-5.6 USM
	CanonLensEF2880mmF3556USMIV             CanonLensType = 202   // Canon EF 28-80mm f/3.5-5.6 USM IV
	CanonLensEF2255mmF456USM                CanonLensType = 208   // Canon EF 22-55mm f/4-5.6 USM
	CanonLensEF55200mmF4556                 CanonLensType = 209   // Canon EF 55-200mm f/4.5-5.6
	CanonLensEF2890mmF456USM                CanonLensType = 210   // Canon EF 28-90mm f/4-5.6 USM
	CanonLensEF28200mmF3556USM              CanonLensType = 211   // Canon EF 28-200mm f/3.5-5.6 USM
	CanonLensEF28105mmF456USM               CanonLensType = 212   // Canon EF 28-105mm f/4-5.6 USM
	CanonLensEF90300mmF4556USM              CanonLensType = 213   // Canon EF 90-300mm f/4.5-5.6 USM
	CanonLensEFS1855mmF3556USM              CanonLensType = 214   // Canon EF-S 18-55mm f/3.5-5.6 USM
	CanonLensEF55200mmF4556IIUSM            CanonLensType = 215   // Canon EF 55-200mm f/4.5-5.6 II USM
	CanonLensEF70200mmF28LISUSM             CanonLensType = 224   // Canon EF 70-200mm f/2.8L IS USM
	CanonLensEF70200mmF28LISUSMPlus14x      CanonLensType = 225   // Canon EF 70-200mm f/2.8L IS USM + 1.4x
	CanonLensEF70200mmF28LISUSMPlus2x       CanonLensType = 226   // Canon EF 70-200mm f/2.8L IS USM + 2x
	CanonLensEF70200mmF28LISUSMPlus28x      CanonLensType = 227   // Canon EF 70-200mm f/2.8L IS USM + 2.8x
	CanonLensEF28105mmF3545USMID228         CanonLensType = 228   // Canon EF 28-105mm f/3.5-4.5 USM
	CanonLensEF1635mmF28LUSM                CanonLensType = 229   // Canon EF 16-35mm f/2.8L USM
	CanonLensEF2470mmF28LUSM                CanonLensType = 230   // Canon EF 24-70mm f/2.8L USM
	CanonLensEF1740mmF4LUSM                 CanonLensType = 231   // Canon EF 17-40mm f/4L USM
	CanonLensEF70300mmF4556DOISUSM          CanonLensType = 232   // Canon EF 70-300mm f/4.5-5.6 DO IS USM
	CanonLensEF28300mmF3556LISUSM           CanonLensType = 233   // Canon EF 28-300mm f/3.5-5.6L IS USM
	CanonLensEFS1785mmF456ISUSM             CanonLensType = 234   // Canon EF-S 17-85mm f/4-5.6 IS USM
	CanonLensEFS1022mmF3545USM              CanonLensType = 235   // Canon EF-S 10-22mm f/3.5-4.5 USM
	CanonLensEFS60mmF28MacroUSM             CanonLensType = 236   // Canon EF-S 60mm f/2.8 Macro USM
	CanonLensEF24105mmF4LISUSM              CanonLensType = 237   // Canon EF 24-105mm f/4L IS USM
	CanonLensEF70300mmF456ISUSM             CanonLensType = 238   // Canon EF 70-300mm f/4-5.6 IS USM
	CanonLensEF85mmF12LIIUSM                CanonLensType = 239   // Canon EF 85mm f/1.2L II USM
	CanonLensEFS1755mmF28ISUSM              CanonLensType = 240   // Canon EF-S 17-55mm f/2.8 IS USM
	CanonLensEF50mmF12LUSM                  CanonLensType = 241   // Canon EF 50mm f/1.2L USM
	CanonLensEF70200mmF4LISUSM              CanonLensType = 242   // Canon EF 70-200mm f/4L IS USM
	CanonLensEF70200mmF4LISUSMPlus14x       CanonLensType = 243   // Canon EF 70-200mm f/4L IS USM + 1.4x
	CanonLensEF70200mmF4LISUSMPlus2x        CanonLensType = 244   // Canon EF 70-200mm f/4L IS USM + 2x
	CanonLensEF70200mmF4LISUSMPlus28x       CanonLensType = 245   // Canon EF 70-200mm f/4L IS USM + 2.8x
	CanonLensEF1635mmF28LIIUSM              CanonLensType = 246   // Canon EF 16-35mm f/2.8L II USM
	CanonLensEF14mmF28LIIUSM                CanonLensType = 247   // Canon EF 14mm f/2.8L II USM
	CanonLensEF200mmF2LISUSM                CanonLensType = 248   // Canon EF 200mm f/2L IS USM
	CanonLensEF800mmF56LISUSM               CanonLensType = 249   // Canon EF 800mm f/5.6L IS USM
	CanonLensEF24mmF14LIIUSM                CanonLensType = 250   // Canon EF 24mm f/1.4L II USM
	CanonLensEF70200mmF28LISIIUSM           CanonLensType = 251   // Canon EF 70-200mm f/2.8L IS II USM
	CanonLensEF70200mmF28LISIIUSMPlus14x    CanonLensType = 252   // Canon EF 70-200mm f/2.8L IS II USM + 1.4x
	CanonLensEF70200mmF28LISIIUSMPlus2x     CanonLensType = 253   // Canon EF 70-200mm f/2.8L IS II USM + 2x
	CanonLensEF100mmF28LMacroISUSM          CanonLensType = 254   // Canon EF 100mm f/2.8L Macro IS USM
	CanonLensEFS1585mmF3556ISUSM            CanonLensType = 488   // Canon EF-S 15-85mm f/3.5-5.6 IS USM
	CanonLensEF70300mmF456LISUSM            CanonLensType = 489   // Canon EF 70-300mm f/4-5.6L IS USM
	CanonLensEF815mmF4LFisheyeUSM           CanonLensType = 490   // Canon EF 8-15mm f/4L Fisheye USM
	CanonLensEF300mmF28LISIIUSM             CanonLensType = 491   // Canon EF 300mm f/2.8L IS II USM
	CanonLensEF400mmF28LISIIUSM             CanonLensType = 492   // Canon EF 400mm f/2.8L IS II USM
	CanonLensEF500mmF4LISIIUSM              CanonLensType = 493   // Canon EF 500mm f/4L IS II USM
	CanonLensEF600mmF4LISIIUSM              CanonLensType = 494   // Canon EF 600mm f/4L IS II USM
	CanonLensEF2470mmF28LIIUSM              CanonLensType = 495   // Canon EF 24-70mm f/2.8L II USM
	CanonLensEF200400mmF4LISUSM             CanonLensType = 496   // Canon EF 200-400mm f/4L IS USM
	CanonLensEF200400mmF4LISUSMPlus14x      CanonLensType = 499   // Canon EF 200-400mm f/4L IS USM + 1.4x
	CanonLensEF28mmF28ISUSM                 CanonLensType = 502   // Canon EF 28mm f/2.8 IS USM
	CanonLensEF24mmF28ISUSM                 CanonLensType = 503   // Canon EF 24mm f/2.8 IS USM
	CanonLensEF2470mmF4LISUSM               CanonLensType = 504   // Canon EF 24-70mm f/4L IS USM
	CanonLensEF35mmF2ISUSM                  CanonLensType = 505   // Canon EF 35mm f/2 IS USM
	CanonLensEF400mmF4DOISIIUSM             CanonLensType = 506   // Canon EF 400mm f/4 DO IS II USM
	CanonLensEF1635mmF4LISUSM               CanonLensType = 507   // Canon EF 16-35mm f/4L IS USM
	CanonLensEF1124mmF4LUSM                 CanonLensType = 508   // Canon EF 11-24mm f/4L USM
	CanonLensEF100400mmF4556LISIIUSM        CanonLensType = 747   // Canon EF 100-400mm f/4.5-5.6L IS II USM
	CanonLensEF100400mmF4556LISIIUSMPlus14x CanonLensType = 748   // Canon EF 100-400mm f/4.5-5.6L IS II USM + 1.4x
	CanonLensEF100400mmF4556LISIIUSMPlus2x  CanonLensType = 749   // Canon EF 100-400mm f/4.5-5.6L IS II USM + 2x
	CanonLensEF35mmF14LIIUSM                CanonLensType = 750   // Canon EF 35mm f/1.4L II USM
	CanonLensEF1635mmF28LIIIUSM             CanonLensType = 751   // Canon EF 16-35mm f/2.8L III USM
	CanonLensEF24105mmF4LISIIUSM            CanonLensType = 752   // Canon EF 24-105mm f/4L IS II USM
	CanonLensEF85mmF14LISUSM                CanonLensType = 753   // Canon EF 85mm f/1.4L IS USM
	CanonLensEF70200mmF4LISIIUSM            CanonLensType = 754   // Canon EF 70-200mm f/4L IS II USM
	CanonLensEF400mmF28LISIIIUSM            CanonLensType = 757   // Canon EF 400mm f/2.8L IS III USM
	CanonLensEF600mmF4LISIIIUSM             CanonLensType = 758   // Canon EF 600mm f/4L IS III USM
	CanonLensEFS18135mmF3556ISSTM           CanonLensType = 4142  // Canon EF-S 18-135mm f/3.5-5.6 IS STM
	CanonLensEFM1855mmF3556ISSTM            CanonLensType = 4143  // Canon EF-M 18-55mm f/3.5-5.6 IS STM
	CanonLensEF40mmF28STM                   CanonLensType = 4144  // Canon EF 40mm f/2.8 STM
	CanonLensEFM22mmF2STM                   CanonLensType = 4145  // Canon EF-M 22mm f/2 STM
	CanonLensEFS1855mmF3556ISSTM            CanonLensType = 4146  // Canon EF-S 18-55mm f/3.5-5.6 IS STM
	CanonLensEFM1122mmF456ISSTM             CanonLensType = 4147  // Canon EF-M 11-22mm f/4-5.6 IS STM
	CanonLensEFS55250mmF456ISSTM            CanonLensType = 4148  // Canon EF-S 55-250mm f/4-5.6 IS STM
	CanonLensEFM55200mmF4563ISSTM           CanonLensType = 4149  // Canon EF-M 55-200mm f/4.5-6.3 IS STM
	CanonLensEFS1018mmF4556ISSTM            CanonLensType = 4150  // Canon EF-S 10-18mm f/4.5-5.6 IS STM
	CanonLensEF24105mmF3556ISSTM            CanonLensType = 4152  // Canon EF 24-105mm f/3.5-5.6 IS STM
	CanonLensEFM1545mmF3563ISSTM            CanonLensType = 4153  // Canon EF-M 15-45mm f/3.5-6.3 IS STM
	CanonLensEFS24mmF28STM                  CanonLensType = 4154  // Canon EF-S 24mm f/2.8 STM
	CanonLensEFM28mmF35MacroISSTM           CanonLensType = 4155  // Canon EF-M 28mm f/3.5 Macro IS STM
	CanonLensEF50mmF18STM                   CanonLensType = 4156  // Canon EF 50mm f/1.8 STM
	CanonLensEFM18150mmF3563ISSTM           CanonLensType = 4157  // Canon EF-M 18-150mm f/3.5-6.3 IS STM
	CanonLensEFS1855mmF456ISSTM             CanonLensType = 4158  // Canon EF-S 18-55mm f/4-5.6 IS STM
	CanonLensEFM32mmF14STM                  CanonLensType = 4159  // Canon EF-M 32mm f/1.4 STM
	CanonLensEFS35mmF28MacroISSTM           CanonLensType = 4160  // Canon EF-S 35mm f/2.8 Macro IS STM
	CanonLensEF70300mmF456ISIIUSM           CanonLensType = 36910 // Canon EF 70-300mm f/4-5.6 IS II USM
	CanonLensEFS18135mmF3556ISUSM           CanonLensType = 36912 // Canon EF-S 18-135mm f/3.5-5.6 IS USM
	CanonLensRF50mmF12LUSM                  CanonLensType = 61182 // Canon RF 50mm F1.2L USM
	CanonLensCNE14mmT31LF                   CanonLensType = 61491 // Canon CN-E 14mm T3.1 L F
	CanonLensCNE24mmT15LF                   CanonLensType = 61492 // Canon CN-E 24mm T1.5 L F
	CanonLensCNE85mmT13LF                   CanonLensType = 61494 // Canon CN-E 85mm T1.3 L F
	CanonLensCNE135mmT22LF                  CanonLensType = 61495 // Canon CN-E 135mm T2.2 L F
	CanonLensCNE35mmT15LF                   CanonLensType = 61496 // Canon CN-E 35mm T1.5 L F
)

var allCanonLensTypes = []CanonLensType{
	CanonLensEF50mmF18,
	CanonLensEF28mmF28,
	CanonLensEF135mmF28Soft,
	CanonLensEF35105mmF3545,
	CanonLensEF3570mmF3545,
	CanonLensEF2870mmF3545,
	CanonLensEF100300mmF56L,
	CanonLensEF100300mmF56,
	CanonLensEF70210mmF4,
	CanonLensEF50mmF25Macro,
	CanonLensEF35mmF2,
	CanonLensEF15mmF28Fisheye,
	CanonLensEF50200mmF3545L,
	CanonLensEF50200mmF3545,
	CanonLensEF35135mmF3545,
	CanonLensEF3570mmF3545A,
	CanonLensEF2870mmF3545ID18,
	CanonLensEF100200mmF45A,
	CanonLensEF80200mmF28L,
	CanonLensEF2035mmF28L,
	CanonLensEF35105mmF3545ID23,
	CanonLensEF3580mmF456PowerZoom,
	CanonLensEF3580mmF456PowerZoomID25,
	CanonLensEF100mmF28Macro,
	CanonLensEF3580mmF456,
	CanonLensEF80200mmF4556,
	CanonLensEF50mmF18II,
	CanonLensEF35105mmF4556,
	CanonLensEF75300mmF456,
	CanonLensEF24mmF28,
	CanonLensEF3580mmF456ID35,
	CanonLensEF3876mmF4556,
	CanonLensEF3580mmF456ID37,
	CanonLensEF80200mmF4556II,
	CanonLensEF75300mmF456ID39,
	CanonLensEF2880mmF3556,
	CanonLensEF2890mmF456,
	CanonLensEF28200mmF3556,
	CanonLensEF28105mmF456,
	CanonLensEF90300mmF4556,
	CanonLensEFS1855mmF3556II,
	CanonLensEF2890mmF456ID46,
	CanonLensEFS1855mmF3556IS,
	CanonLensEFS55250mmF456IS,
	CanonLensEFS18200mmF3556IS,
	CanonLensEFS18135mmF3556IS,
	CanonLensEFS1855mmF3556ISII,
	CanonLensEFS1855mmF3556III,
	CanonLensEFS55250mmF456ISII,
	CanonLensTSE50mmF28LMacro,
	CanonLensTSE90mmF28LMacro,
	CanonLensTSE135mmF4LMacro,
	CanonLensTSE17mmF4L,
	CanonLensTSE24mmF35LII,
	CanonLensMPE65mmF2815xMacroPhoto,
	CanonLensTSE24mmF35L,
	CanonLensTSE45mmF28,
	CanonLensTSE90mmF28,
	CanonLensEF300mmF28LUSM,
	CanonLensEF50mmF10LUSM,
	CanonLensEF2880mmF284LUSM,
	CanonLensEF1200mmF56LUSM,
	CanonLensEF600mmF4LISUSM,
	CanonLensEF200mmF18LUSM,
	CanonLensEF300mmF28LUSMID136,
	CanonLensEF85mmF12LUSM,
	CanonLensEF2880mmF284L,
	CanonLensEF400mmF28LUSM,
	CanonLensEF500mmF45LUSM,
	CanonLensEF500mmF45LUSMID141,
	CanonLensEF300mmF28LISUSM,
	CanonLensEF500mmF4LISUSM,
	CanonLensEF35135mmF456USM,
	CanonLensEF100300mmF4556USM,
	CanonLensEF70210mmF3545USM,
	CanonLensEF35135mmF456USMID147,
	CanonLensEF2880mmF3556USM,
	CanonLensEF100mmF2USM,
	CanonLensEF14mmF28LUSM,
	CanonLensEF200mmF28LUSM,
	CanonLensEF300mmF4LISUSM,
	CanonLensEF35350mmF3556LUSM,
	CanonLensEF20mmF28USM,
	CanonLensEF85mmF18USM,
	CanonLensEF28105mmF3545USM,
	CanonLensEF2035mmF3545USM,
	CanonLensEF2870mmF28LUSM,
	CanonLensEF200mmF28LUSMID162,
	CanonLensEF300mmF4L,
	CanonLensEF400mmF56L,
	CanonLensEF70200mmF28LUSM,
	CanonLensEF70200mmF28LUSMPlus14x,
	CanonLensEF70200mmF28LUSMPlus2x,
	CanonLensEF28mmF18USM,
	CanonLensEF1735mmF28LUSM,
	CanonLensEF200mmF28LIIUSM,
	CanonLensEF300mmF4LUSM,
	CanonLensEF400mmF56LUSM,
	CanonLensEF180mmMacroF35LUSM,
	CanonLensEF135mmF2LUSM,
	CanonLensEF400mmF28LUSMID175,
	CanonLensEF2485mmF3545USM,
	CanonLensEF300mmF4LISUSMID177,
	CanonLensEF28135mmF3556IS,
	CanonLensEF24mmF14LUSM,
	CanonLensEF35mmF14LUSM,
	CanonLensEF100400mmF4556LISUSMPlus14x,
	CanonLensEF100400mmF4556LISUSMPlus2x,
	CanonLensEF100400mmF4556LISUSM,
	CanonLensEF400mmF28LUSMPlus2x,
	CanonLensEF600mmF4LISUSMID185,
	CanonLensEF70200mmF4LUSM,
	CanonLensEF70200mmF4LUSMPlus14x,
	CanonLensEF70200mmF4LUSMPlus2x,
	CanonLensEF70200mmF4LUSMPlus28x,
	CanonLensEF100mmF28MacroUSM,
	CanonLensEF400mmF4DOIS,
	CanonLensEF3580mmF456USM,
	CanonLensEF80200mmF4556USM,
	CanonLensEF35105mmF4556USM,
	CanonLensEF75300mmF456USM,
	CanonLensEF75300mmF456ISUSM,
	CanonLensEF50mmF14USM,
	CanonLensEF2880mmF3556USMID199,
	CanonLensEF75300mmF456USMID200,
	CanonLensEF2880mmF3556USMID201,
	CanonLensEF2880mmF3556USMIV,
	CanonLensEF2255mmF456USM,
	CanonLensEF55200mmF4556,
	CanonLensEF2890mmF456USM,
	CanonLensEF28200mmF3556USM,
	CanonLensEF28105mmF456USM,
	CanonLensEF90300mmF4556USM,
	CanonLensEFS1855mmF3556USM,
	CanonLensEF55200mmF4556IIUSM,
	CanonLensEF70200mmF28LISUSM,
	CanonLensEF70200mmF28LISUSMPlus14x,
	CanonLensEF70200mmF28LISUSMPlus2x,
	CanonLensEF70200mmF28LISUSMPlus28x,
	CanonLensEF28105mmF3545USMID228,
	CanonLensEF1635mmF28LUSM,
	CanonLensEF2470mmF28LUSM,
	CanonLensEF1740mmF4LUSM,
	CanonLensEF70300mmF4556DOISUSM,
	CanonLensEF28300mmF3556LISUSM,
	CanonLensEFS1785mmF456ISUSM,
	CanonLensEFS1022mmF3545USM,
	CanonLensEFS60mmF28MacroUSM,
	CanonLensEF24105mmF4LISUSM,
	CanonLensEF70300mmF456ISUSM,
	CanonLensEF85mmF12LIIUSM,
	CanonLensEFS1755mmF28ISUSM,
	CanonLensEF50mmF12LUSM,
	CanonLensEF70200mmF4LISUSM,
	CanonLensEF70200mmF4LISUSMPlus14x,
	CanonLensEF70200mmF4LISUSMPlus2x,
	CanonLensEF70200mmF4LISUSMPlus28x,
	CanonLensEF1635mmF28LIIUSM,
	CanonLensEF14mmF28LIIUSM,
	CanonLensEF200mmF2LISUSM,
	CanonLensEF800mmF56LISUSM,
	CanonLensEF24mmF14LIIUSM,
	CanonLensEF70200mmF28LISIIUSM,
	CanonLensEF70200mmF28LISIIUSMPlus14x,
	CanonLensEF70200mmF28LISIIUSMPlus2x,
	CanonLensEF100mmF28LMacroISUSM,
	CanonLensEFS1585mmF3556ISUSM,
	CanonLensEF70300mmF456LISUSM,
	CanonLensEF815mmF4LFisheyeUSM,
	CanonLensEF300mmF28LISIIUSM,
	CanonLensEF400mmF28LISIIUSM,
	CanonLensEF500mmF4LISIIUSM,
	CanonLensEF600mmF4LISIIUSM,
	CanonLensEF2470mmF28LIIUSM,
	CanonLensEF200400mmF4LISUSM,
	CanonLensEF200400mmF4LISUSMPlus14x,
	CanonLensEF28mmF28ISUSM,
	CanonLensEF24mmF28ISUSM,
	CanonLensEF2470mmF4LISUSM,
	CanonLensEF35mmF2ISUSM,
	CanonLensEF400mmF4DOISIIUSM,
	CanonLensEF1635mmF4LISUSM,
	CanonLensEF1124mmF4LUSM,
	CanonLensEF100400mmF4556LISIIUSM,
	CanonLensEF100400mmF4556LISIIUSMPlus14x,
	CanonLensEF100400mmF4556LISIIUSMPlus2x,
	CanonLensEF35mmF14LIIUSM,
	CanonLensEF1635mmF28LIIIUSM,
	CanonLensEF24105mmF4LISIIUSM,
	CanonLensEF85mmF14LISUSM,
	CanonLensEF70200mmF4LISIIUSM,
	CanonLensEF400mmF28LISIIIUSM,
	CanonLensEF600mmF4LISIIIUSM,
	CanonLensEFS18135mmF3556ISSTM,
	CanonLensEFM1855mmF3556ISSTM,
	CanonLensEF40mmF28STM,
	CanonLensEFM22mmF2STM,
	CanonLensEFS1855mmF3556ISSTM,
	CanonLensEFM1122mmF456ISSTM,
	CanonLensEFS55250mmF456ISSTM,
	CanonLensEFM55200mmF4563ISSTM,
	CanonLensEFS1018mmF4556ISSTM,
	CanonLensEF24105mmF3556ISSTM,
	CanonLensEFM1545mmF3563ISSTM,
	CanonLensEFS24mmF28STM,
	CanonLensEFM28mmF35MacroISSTM,
	CanonLensEF50mmF18STM,
	CanonLensEFM18150mmF3563ISSTM,
	CanonLensEFS1855mmF456ISSTM,
	CanonLensEFM32mmF14STM,
	CanonLensEFS35mmF28MacroISSTM,
	CanonLensEF70300mmF456ISIIUSM,
	CanonLensEFS18135mmF3556ISUSM,
	CanonLensRF50mmF12LUSM,
	CanonLensCNE14mmT31LF,
	CanonLensCNE24mmT15LF,
	CanonLensCNE85mmT13LF,
	CanonLensCNE135mmT22LF,
	CanonLensCNE35mmT15LF,
}

func AllCanonLensTypes() []CanonLensType {
	out := make([]CanonLensType, len(allCanonLensTypes))
	copy(out, allCanonLensTypes)
	return out
}

var canonLensTypeLabels = map[CanonLensType]string{
	CanonLensEF50mmF18:                      "Canon EF 50mm f/1.8",
	CanonLensEF28mmF28:                      "Canon EF 28mm f/2.8",
	CanonLensEF135mmF28Soft:                 "Canon EF 135mm f/2.8 Soft",
	CanonLensEF35105mmF3545:                 "Canon EF 35-105mm f/3.5-4.5",
	CanonLensEF3570mmF3545:                  "Canon EF 35-70mm f/3.5-4.5",
	CanonLensEF2870mmF3545:                  "Canon EF 28-70mm f/3.5-4.5",
	CanonLensEF100300mmF56L:                 "Canon EF 100-300mm f/5.6L",
	CanonLensEF100300mmF56:                  "Canon EF 100-300mm f/5.6",
	CanonLensEF70210mmF4:                    "Canon EF 70-210mm f/4",
	CanonLensEF50mmF25Macro:                 "Canon EF 50mm f/2.5 Macro",
	CanonLensEF35mmF2:                       "Canon EF 35mm f/2",
	CanonLensEF15mmF28Fisheye:               "Canon EF 15mm f/2.8 Fisheye",
	CanonLensEF50200mmF3545L:                "Canon EF 50-200mm f/3.5-4.5L",
	CanonLensEF50200mmF3545:                 "Canon EF 50-200mm f/3.5-4.5",
	CanonLensEF35135mmF3545:                 "Canon EF 35-135mm f/3.5-4.5",
	CanonLensEF3570mmF3545A:                 "Canon EF 35-70mm f/3.5-4.5A",
	CanonLensEF2870mmF3545ID18:              "Canon EF 28-70mm f/3.5-4.5",
	CanonLensEF100200mmF45A:                 "Canon EF 100-200mm f/4.5A",
	CanonLensEF80200mmF28L:                  "Canon EF 80-200mm f/2.8L",
	CanonLensEF2035mmF28L:                   "Canon EF 20-35mm f/2.8L",
	CanonLensEF35105mmF3545ID23:             "Canon EF 35-105mm f/3.5-4.5",
	CanonLensEF3580mmF456PowerZoom:          "Canon EF 35-80mm f/4-5.6 Power Zoom",
	CanonLensEF3580mmF456PowerZoomID25:      "Canon EF 35-80mm f/4-5.6 Power Zoom",
	CanonLensEF100mmF28Macro:                "Canon EF 100mm f/2.8 Macro",
	CanonLensEF3580mmF456:                   "Canon EF 35-80mm f/4-5.6",
	CanonLensEF80200mmF4556:                 "Canon EF 80-200mm f/4.5-5.6",
	CanonLensEF50mmF18II:                    "Canon EF 50mm f/1.8 II",
	CanonLensEF35105mmF4556:                 "Canon EF 35-105mm f/4.5-5.6",
	CanonLensEF75300mmF456:                  "Canon EF 75-300mm f/4-5.6",
	CanonLensEF24mmF28:                      "Canon EF 24mm f/2.8",
	CanonLensEF3580mmF456ID35:               "Canon EF 35-80mm f/4-5.6",
	CanonLensEF3876mmF4556:                  "Canon EF 38-76mm f/4.5-5.6",
	CanonLensEF3580mmF456ID37:               "Canon EF 35-80mm f/4-5.6",
	CanonLensEF80200mmF4556II:               "Canon EF 80-200mm f/4.5-5.6 II",
	CanonLensEF75300mmF456ID39:              "Canon EF 75-300mm f/4-5.6",
	CanonLensEF2880mmF3556:                  "Canon EF 28-80mm f/3.5-5.6",
	CanonLensEF2890mmF456:                   "Canon EF 28-90mm f/4-5.6",
	CanonLensEF28200mmF3556:                 "Canon EF 28-200mm f/3.5-5.6",
	CanonLensEF28105mmF456:                  "Canon EF 28-105mm f/4-5.6",
	CanonLensEF90300mmF4556:                 "Canon EF 90-300mm f/4.5-5.6",
	CanonLensEFS1855mmF3556II:               "Canon EF-S 18-55mm f/3.5-5.6 [II]",
	CanonLensEF2890mmF456ID46:               "Canon EF 28-90mm f/4-5.6",
	CanonLensEFS1855mmF3556IS:               "Canon EF-S 18-55mm f/3.5-5.6 IS",
	CanonLensEFS55250mmF456IS:               "Canon EF-S 55-250mm f/4-5.6 IS",
	CanonLensEFS18200mmF3556IS:              "Canon EF-S 18-200mm f/3.5-5.6 IS",
	CanonLensEFS18135mmF3556IS:              "Canon EF-S 18-135mm f/3.5-5.6 IS",
	CanonLensEFS1855mmF3556ISII:             "Canon EF-S 18-55mm f/3.5-5.6 IS II",
	CanonLensEFS1855mmF3556III:              "Canon EF-S 18-55mm f/3.5-5.6 III",
	CanonLensEFS55250mmF456ISII:             "Canon EF-S 55-250mm f/4-5.6 IS II",
	CanonLensTSE50mmF28LMacro:               "Canon TS-E 50mm f/2.8L Macro",
	CanonLensTSE90mmF28LMacro:               "Canon TS-E 90mm f/2.8L Macro",
	CanonLensTSE135mmF4LMacro:               "Canon TS-E 135mm f/4L Macro",
	CanonLensTSE17mmF4L:                     "Canon TS-E 17mm f/4L",
	CanonLensTSE24mmF35LII:                  "Canon TS-E 24mm f/3.5L II",
	CanonLensMPE65mmF2815xMacroPhoto:        "Canon MP-E 65mm f/2.8 1-5x Macro Photo",
	CanonLensTSE24mmF35L:                    "Canon TS-E 24mm f/3.5L",
	CanonLensTSE45mmF28:                     "Canon TS-E 45mm f/2.8",
	CanonLensTSE90mmF28:                     "Canon TS-E 90mm f/2.8",
	CanonLensEF300mmF28LUSM:                 "Canon EF 300mm f/2.8L USM",
	CanonLensEF50mmF10LUSM:                  "Canon EF 50mm f/1.0L USM",
	CanonLensEF2880mmF284LUSM:               "Canon EF 28-80mm f/2.8-4L USM",
	CanonLensEF1200mmF56LUSM:                "Canon EF 1200mm f/5.6L USM",
	CanonLensEF600mmF4LISUSM:                "Canon EF 600mm f/4L IS USM",
	CanonLensEF200mmF18LUSM:                 "Canon EF 200mm f/1.8L USM",
	CanonLensEF300mmF28LUSMID136:            "Canon EF 300mm f/2.8L USM",
	CanonLensEF85mmF12LUSM:                  "Canon EF 85mm f/1.2L USM",
	CanonLensEF2880mmF284L:                  "Canon EF 28-80mm f/2.8-4L",
	CanonLensEF400mmF28LUSM:                 "Canon EF 400mm f/2.8L USM",
	CanonLensEF500mmF45LUSM:                 "Canon EF 500mm f/4.5L USM",
	CanonLensEF500mmF45LUSMID141:            "Canon EF 500mm f/4.5L USM",
	CanonLensEF300mmF28LISUSM:               "Canon EF 300mm f/2.8L IS USM",
	CanonLensEF500mmF4LISUSM:                "Canon EF 500mm f/4L IS USM",
	CanonLensEF35135mmF456USM:               "Canon EF 35-135mm f/4-5.6 USM",
	CanonLensEF100300mmF4556USM:             "Canon EF 100-300mm f/4.5-5.6 USM",
	CanonLensEF70210mmF3545USM:              "Canon EF 70-210mm f/3.5-4.5 USM",
	CanonLensEF35135mmF456USMID147:          "Canon EF 35-135mm f/4-5.6 USM",
	CanonLensEF2880mmF3556USM:               "Canon EF 28-80mm f/3.5-5.6 USM",
	CanonLensEF100mmF2USM:                   "Canon EF 100mm f/2 USM",
	CanonLensEF14mmF28LUSM:                  "Canon EF 14mm f/2.8L USM",
	CanonLensEF200mmF28LUSM:                 "Canon EF 200mm f/2.8L USM",
	CanonLensEF300mmF4LISUSM:                "Canon EF 300mm f/4L IS USM",
	CanonLensEF35350mmF3556LUSM:             "Canon EF 35-350mm f/3.5-5.6L USM",
	CanonLensEF20mmF28USM:                   "Canon EF 20mm f/2.8 USM",
	CanonLensEF85mmF18USM:                   "Canon EF 85mm f/1.8 USM",
	CanonLensEF28105mmF3545USM:              "Canon EF 28-105mm f/3.5-4.5 USM",
	CanonLensEF2035mmF3545USM:               "Canon EF 20-35mm f/3.5-4.5 USM",
	CanonLensEF2870mmF28LUSM:                "Canon EF 28-70mm f/2.8L USM",
	CanonLensEF200mmF28LUSMID162:            "Canon EF 200mm f/2.8L USM",
	CanonLensEF300mmF4L:                     "Canon EF 300mm f/4L",
	CanonLensEF400mmF56L:                    "Canon EF 400mm f/5.6L",
	CanonLensEF70200mmF28LUSM:               "Canon EF 70-200mm f/2.8L USM",
	CanonLensEF70200mmF28LUSMPlus14x:        "Canon EF 70-200mm f/2.8L USM + 1.4x",
	CanonLensEF70200mmF28LUSMPlus2x:         "Canon EF 70-200mm f/2.8L USM + 2x",
	CanonLensEF28mmF18USM:                   "Canon EF 28mm f/1.8 USM",
	CanonLensEF1735mmF28LUSM:                "Canon EF 17-35mm f/2.8L USM",
	CanonLensEF200mmF28LIIUSM:               "Canon EF 200mm f/2.8L II USM",
	CanonLensEF300mmF4LUSM:                  "Canon EF 300mm f/4L USM",
	CanonLensEF400mmF56LUSM:                 "Canon EF 400mm f/5.6L USM",
	CanonLensEF180mmMacroF35LUSM:            "Canon EF 180mm Macro f/3.5L USM",
	CanonLensEF135mmF2LUSM:                  "Canon EF 135mm f/2L USM",
	CanonLensEF400mmF28LUSMID175:            "Canon EF 400mm f/2.8L USM",
	CanonLensEF2485mmF3545USM:               "Canon EF 24-85mm f/3.5-4.5 USM",
	CanonLensEF300mmF4LISUSMID177:           "Canon EF 300mm f/4L IS USM",
	CanonLensEF28135mmF3556IS:               "Canon EF 28-135mm f/3.5-5.6 IS",
	CanonLensEF24mmF14LUSM:                  "Canon EF 24mm f/1.4L USM",
	CanonLensEF35mmF14LUSM:                  "Canon EF 35mm f/1.4L USM",
	CanonLensEF100400mmF4556LISUSMPlus14x:   "Canon EF 100-400mm f/4.5-5.6L IS USM + 1.4x",
	CanonLensEF100400mmF4556LISUSMPlus2x:    "Canon EF 100-400mm f/4.5-5.6L IS USM + 2x",
	CanonLensEF100400mmF4556LISUSM:          "Canon EF 100-400mm f/4.5-5.6L IS USM",
	CanonLensEF400mmF28LUSMPlus2x:           "Canon EF 400mm f/2.8L USM + 2x",
	CanonLensEF600mmF4LISUSMID185:           "Canon EF 600mm f/4L IS USM",
	CanonLensEF70200mmF4LUSM:                "Canon EF 70-200mm f/4L USM",
	CanonLensEF70200mmF4LUSMPlus14x:         "Canon EF 70-200mm f/4L USM + 1.4x",
	CanonLensEF70200mmF4LUSMPlus2x:          "Canon EF 70-200mm f/4L USM + 2x",
	CanonLensEF70200mmF4LUSMPlus28x:         "Canon EF 70-200mm f/4L USM + 2.8x",
	CanonLensEF100mmF28MacroUSM:             "Canon EF 100mm f/2.8 Macro USM",
	CanonLensEF400mmF4DOIS:                  "Canon EF 400mm f/4 DO IS",
	CanonLensEF3580mmF456USM:                "Canon EF 35-80mm f/4-5.6 USM",
	CanonLensEF80200mmF4556USM:              "Canon EF 80-200mm f/4.5-5.6 USM",
	CanonLensEF35105mmF4556USM:              "Canon EF 35-105mm f/4.5-5.6 USM",
	CanonLensEF75300mmF456USM:               "Canon EF 75-300mm f/4-5.6 USM",
	CanonLensEF75300mmF456ISUSM:             "Canon EF 75-300mm f/4-5.6 IS USM",
	CanonLensEF50mmF14USM:                   "Canon EF 50mm f/1.4 USM",
	CanonLensEF2880mmF3556USMID199:          "Canon EF 28-80mm f/3.5-5.6 USM",
	CanonLensEF75300mmF456USMID200:          "Canon EF 75-300mm f/4-5.6 USM",
	CanonLensEF2880mmF3556USMID201:          "Canon EF 28-80mm f/3.5-5.6 USM",
	CanonLensEF2880mmF3556USMIV:             "Canon EF 28-80mm f/3.5-5.6 USM IV",
	CanonLensEF2255mmF456USM:                "Canon EF 22-55mm f/4-5.6 USM",
	CanonLensEF55200mmF4556:                 "Canon EF 55-200mm f/4.5-5.6",
	CanonLensEF2890mmF456USM:                "Canon EF 28-90mm f/4-5.6 USM",
	CanonLensEF28200mmF3556USM:              "Canon EF 28-200mm f/3.5-5.6 USM",
	CanonLensEF28105mmF456USM:               "Canon EF 28-105mm f/4-5.6 USM",
	CanonLensEF90300mmF4556USM:              "Canon EF 90-300mm f/4.5-5.6 USM",
	CanonLensEFS1855mmF3556USM:              "Canon EF-S 18-55mm f/3.5-5.6 USM",
	CanonLensEF55200mmF4556IIUSM:            "Canon EF 55-200mm f/4.5-5.6 II USM",
	CanonLensEF70200mmF28LISUSM:             "Canon EF 70-200mm f/2.8L IS USM",
	CanonLensEF70200mmF28LISUSMPlus14x:      "Canon EF 70-200mm f/2.8L IS USM + 1.4x",
	CanonLensEF70200mmF28LISUSMPlus2x:       "Canon EF 70-200mm f/2.8L IS USM + 2x",
	CanonLensEF70200mmF28LISUSMPlus28x:      "Canon EF 70-200mm f/2.8L IS USM + 2.8x",
	CanonLensEF28105mmF3545USMID228:         "Canon EF 28-105mm f/3.5-4.5 USM",
	CanonLensEF1635mmF28LUSM:                "Canon EF 16-35mm f/2.8L USM",
	CanonLensEF2470mmF28LUSM:                "Canon EF 24-70mm f/2.8L USM",
	CanonLensEF1740mmF4LUSM:                 "Canon EF 17-40mm f/4L USM",
	CanonLensEF70300mmF4556DOISUSM:          "Canon EF 70-300mm f/4.5-5.6 DO IS USM",
	CanonLensEF28300mmF3556LISUSM:           "Canon EF 28-300mm f/3.5-5.6L IS USM",
	CanonLensEFS1785mmF456ISUSM:             "Canon EF-S 17-85mm f/4-5.6 IS USM",
	CanonLensEFS1022mmF3545USM:              "Canon EF-S 10-22mm f/3.5-4.5 USM",
	CanonLensEFS60mmF28MacroUSM:             "Canon EF-S 60mm f/2.8 Macro USM",
	CanonLensEF24105mmF4LISUSM:              "Canon EF 24-105mm f/4L IS USM",
	CanonLensEF70300mmF456ISUSM:             "Canon EF 70-300mm f/4-5.6 IS USM",
	CanonLensEF85mmF12LIIUSM:                "Canon EF 85mm f/1.2L II USM",
	CanonLensEFS1755mmF28ISUSM:              "Canon EF-S 17-55mm f/2.8 IS USM",
	CanonLensEF50mmF12LUSM:                  "Canon EF 50mm f/1.2L USM",
	CanonLensEF70200mmF4LISUSM:              "Canon EF 70-200mm f/4L IS USM",
	CanonLensEF70200mmF4LISUSMPlus14x:       "Canon EF 70-200mm f/4L IS USM + 1.4x",
	CanonLensEF70200mmF4LISUSMPlus2x:        "Canon EF 70-200mm f/4L IS USM + 2x",
	CanonLensEF70200mmF4LISUSMPlus28x:       "Canon EF 70-200mm f/4L IS USM + 2.8x",
	CanonLensEF1635mmF28LIIUSM:              "Canon EF 16-35mm f/2.8L II USM",
	CanonLensEF14mmF28LIIUSM:                "Canon EF 14mm f/2.8L II USM",
	CanonLensEF200mmF2LISUSM:                "Canon EF 200mm f/2L IS USM",
	CanonLensEF800mmF56LISUSM:               "Canon EF 800mm f/5.6L IS USM",
	CanonLensEF24mmF14LIIUSM:                "Canon EF 24mm f/1.4L II USM",
	CanonLensEF70200mmF28LISIIUSM:           "Canon EF 70-200mm f/2.8L IS II USM",
	CanonLensEF70200mmF28LISIIUSMPlus14x:    "Canon EF 70-200mm f/2.8L IS II USM + 1.4x",
	CanonLensEF70200mmF28LISIIUSMPlus2x:     "Canon EF 70-200mm f/2.8L IS II USM + 2x",
	CanonLensEF100mmF28LMacroISUSM:          "Canon EF 100mm f/2.8L Macro IS USM",
	CanonLensEFS1585mmF3556ISUSM:            "Canon EF-S 15-85mm f/3.5-5.6 IS USM",
	CanonLensEF70300mmF456LISUSM:            "Canon EF 70-300mm f/4-5.6L IS USM",
	CanonLensEF815mmF4LFisheyeUSM:           "Canon EF 8-15mm f/4L Fisheye USM",
	CanonLensEF300mmF28LISIIUSM:             "Canon EF 300mm f/2.8L IS II USM",
	CanonLensEF400mmF28LISIIUSM:             "Canon EF 400mm f/2.8L IS II USM",
	CanonLensEF500mmF4LISIIUSM:              "Canon EF 500mm f/4L IS II USM",
	CanonLensEF600mmF4LISIIUSM:              "Canon EF 600mm f/4L IS II USM",
	CanonLensEF2470mmF28LIIUSM:              "Canon EF 24-70mm f/2.8L II USM",
	CanonLensEF200400mmF4LISUSM:             "Canon EF 200-400mm f/4L IS USM",
	CanonLensEF200400mmF4LISUSMPlus14x:      "Canon EF 200-400mm f/4L IS USM + 1.4x",
	CanonLensEF28mmF28ISUSM:                 "Canon EF 28mm f/2.8 IS USM",
	CanonLensEF24mmF28ISUSM:                 "Canon EF 24mm f/2.8 IS USM",
	CanonLensEF2470mmF4LISUSM:               "Canon EF 24-70mm f/4L IS USM",
	CanonLensEF35mmF2ISUSM:                  "Canon EF 35mm f/2 IS USM",
	CanonLensEF400mmF4DOISIIUSM:             "Canon EF 400mm f/4 DO IS II USM",
	CanonLensEF1635mmF4LISUSM:               "Canon EF 16-35mm f/4L IS USM",
	CanonLensEF1124mmF4LUSM:                 "Canon EF 11-24mm f/4L USM",
	CanonLensEF100400mmF4556LISIIUSM:        "Canon EF 100-400mm f/4.5-5.6L IS II USM",
	CanonLensEF100400mmF4556LISIIUSMPlus14x: "Canon EF 100-400mm f/4.5-5.6L IS II USM + 1.4x",
	CanonLensEF100400mmF4556LISIIUSMPlus2x:  "Canon EF 100-400mm f/4.5-5.6L IS II USM + 2x",
	CanonLensEF35mmF14LIIUSM:                "Canon EF 35mm f/1.4L II USM",
	CanonLensEF1635mmF28LIIIUSM:             "Canon EF 16-35mm f/2.8L III USM",
	CanonLensEF24105mmF4LISIIUSM:            "Canon EF 24-105mm f/4L IS II USM",
	CanonLensEF85mmF14LISUSM:                "Canon EF 85mm f/1.4L IS USM",
	CanonLensEF70200mmF4LISIIUSM:            "Canon EF 70-200mm f/4L IS II USM",
	CanonLensEF400mmF28LISIIIUSM:            "Canon EF 400mm f/2.8L IS III USM",
	CanonLensEF600mmF4LISIIIUSM:             "Canon EF 600mm f/4L IS III USM",
	CanonLensEFS18135mmF3556ISSTM:           "Canon EF-S 18-135mm f/3.5-5.6 IS STM",
	CanonLensEFM1855mmF3556ISSTM:            "Canon EF-M 18-55mm f/3.5-5.6 IS STM",
	CanonLensEF40mmF28STM:                   "Canon EF 40mm f/2.8 STM",
	CanonLensEFM22mmF2STM:                   "Canon EF-M 22mm f/2 STM",
	CanonLensEFS1855mmF3556ISSTM:            "Canon EF-S 18-55mm f/3.5-5.6 IS STM",
	CanonLensEFM1122mmF456ISSTM:             "Canon EF-M 11-22mm f/4-5.6 IS STM",
	CanonLensEFS55250mmF456ISSTM:            "Canon EF-S 55-250mm f/4-5.6 IS STM",
	CanonLensEFM55200mmF4563ISSTM:           "Canon EF-M 55-200mm f/4.5-6.3 IS STM",
	CanonLensEFS1018mmF4556ISSTM:            "Canon EF-S 10-18mm f/4.5-5.6 IS STM",
	CanonLensEF24105mmF3556ISSTM:            "Canon EF 24-105mm f/3.5-5.6 IS STM",
	CanonLensEFM1545mmF3563ISSTM:            "Canon EF-M 15-45mm f/3.5-6.3 IS STM",
	CanonLensEFS24mmF28STM:                  "Canon EF-S 24mm f/2.8 STM",
	CanonLensEFM28mmF35MacroISSTM:           "Canon EF-M 28mm f/3.5 Macro IS STM",
	CanonLensEF50mmF18STM:                   "Canon EF 50mm f/1.8 STM",
	CanonLensEFM18150mmF3563ISSTM:           "Canon EF-M 18-150mm f/3.5-6.3 IS STM",
	CanonLensEFS1855mmF456ISSTM:             "Canon EF-S 18-55mm f/4-5.6 IS STM",
	CanonLensEFM32mmF14STM:                  "Canon EF-M 32mm f/1.4 STM",
	CanonLensEFS35mmF28MacroISSTM:           "Canon EF-S 35mm f/2.8 Macro IS STM",
	CanonLensEF70300mmF456ISIIUSM:           "Canon EF 70-300mm f/4-5.6 IS II USM",
	CanonLensEFS18135mmF3556ISUSM:           "Canon EF-S 18-135mm f/3.5-5.6 IS USM",
	CanonLensRF50mmF12LUSM:                  "Canon RF 50mm F1.2L USM",
	CanonLensCNE14mmT31LF:                   "Canon CN-E 14mm T3.1 L F",
	CanonLensCNE24mmT15LF:                   "Canon CN-E 24mm T1.5 L F",
	CanonLensCNE85mmT13LF:                   "Canon CN-E 85mm T1.3 L F",
	CanonLensCNE135mmT22LF:                  "Canon CN-E 135mm T2.2 L F",
	CanonLensCNE35mmT15LF:                   "Canon CN-E 35mm T1.5 L F",
}

func (l CanonLensType) String() string {
	if l == CanonLensUnknown {
		return "Unknown"
	}
	if s, ok := canonLensTypeLabels[l]; ok {
		return s
	}
	return "CanonLensType(" + strconv.FormatUint(uint64(l), 10) + ")"
}
