package makernote

// CameraMake identifies a normalized camera manufacturer token parsed from
// IFD0 Make text.
type CameraMake uint8

// CameraMake values.
const (
	CameraMakeUnknown CameraMake = iota
	CameraMakeAgfa
	CameraMakeAmazon
	CameraMakeApple
	CameraMakeARRI
	CameraMakeAsahiOpticalCoLtd
	CameraMakeASUS
	CameraMakeAutelRobotics
	CameraMakeBlackmagicDesign
	CameraMakeCanon
	CameraMakeCasio
	CameraMakeContax
	CameraMakeCosina
	CameraMakeDJI
	CameraMakeDXO
	CameraMakeFujifilm
	CameraMakeGoogle
	CameraMakeGoPro
	CameraMakeHasselblad
	CameraMakeHewlettPackard
	CameraMakeHonor
	CameraMakeHTC
	CameraMakeHuawei
	CameraMakeInsta360
	CameraMakeJKImagingLtd
	CameraMakeKodak
	CameraMakeKonicaMinolta
	CameraMakeKyocera
	CameraMakeLeica
	CameraMakeLGElectronics
	CameraMakeMamiya
	CameraMakeMinoltaCoLtd
	CameraMakeMotorola
	CameraMakeNikon
	CameraMakeNokia
	CameraMakeOlympusCorporation
	CameraMakeOMDigitalSolutions
	CameraMakeOnePlus
	CameraMakeOppo
	CameraMakePanasonic
	CameraMakeParrot
	CameraMakePhaseOne
	CameraMakePentax
	CameraMakePolaroid
	CameraMakeRealme
	CameraMakeRED
	CameraMakeRicoh
	CameraMakeRollei
	CameraMakeSamsung
	CameraMakeSeaLife
	CameraMakeSeikoEpsonCorp
	CameraMakeSigma
	CameraMakeSkydio
	CameraMakeSony
	CameraMakeVivo
	CameraMakeVivitar
	CameraMakeXiaomi
	CameraMakeXiaoyi
	CameraMakeYashica
	CameraMakeYuneec
	CameraMakeZeiss
	CameraMakeZTE
)

const _CameraMake_name = "UnknownAgfaAmazonAppleARRIAsahi OpticalASUSAutel RoboticsBlackmagic DesignCanonCasioContaxCosinaDJIDxOFujifilmGoogleGoProHasselbladHPHonorHTCHuaweiInsta360JK ImagingKodakKonica MinoltaKyoceraLeicaLGMamiyaMinoltaMotorolaNikonNokiaOlympusOM Digital SolutionsOnePlusOPPOPanasonicParrotPhase OnePentaxPolaroidRealmeREDRicohRolleiSamsungSeaLifeEpsonSigmaSkydioSonyvivoVivitarXiaomiXiaoyiYashicaYuneecZeissZTE"

var _CameraMake_index = [...]uint16{0, 7, 11, 17, 22, 26, 39, 43, 57, 74, 79, 84, 90, 96, 99, 102, 110, 116, 121, 131, 133, 138, 141, 147, 155, 165, 170, 184, 191, 196, 198, 204, 211, 219, 224, 229, 236, 256, 263, 267, 276, 282, 291, 297, 305, 311, 314, 319, 325, 332, 339, 344, 349, 355, 359, 363, 370, 376, 382, 389, 395, 400, 403}

// String returns the display name for the camera make value.
func (m CameraMake) String() string {
	i := int(m)
	if i < 0 || i >= len(_CameraMake_index)-1 {
		return CameraMakeUnknown.String()
	}
	return _CameraMake_name[_CameraMake_index[i]:_CameraMake_index[i+1]]
}

// IdentifyCameraMake maps in-place IFD0 Make bytes to a normalized CameraMake.
//
// The input slice is normalized in place (ASCII-lowercased and punctuation/
// whitespace stripped) to avoid allocations.
func IdentifyCameraMake(raw []byte) CameraMake {
	n := 0
	for _, b := range raw {
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		switch b {
		case 0, ' ', '\t', '\n', '\r', '\f', '\v', ',', '.':
			continue
		}
		raw[n] = b
		n++
	}
	return identifyCameraMakeNormalized(raw[:n])
}

// IdentifyCameraMakeString maps IFD0 Make text to the known camera make enum.
func IdentifyCameraMakeString(raw string) CameraMake {
	// Most make values are short; use a stack buffer to avoid allocations.
	var normalized [64]byte
	n := 0
	for i := 0; i < len(raw) && n < len(normalized); i++ {
		b := raw[i]
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		switch b {
		case 0, ' ', '\t', '\n', '\r', '\f', '\v', ',', '.':
			continue
		}
		normalized[n] = b
		n++
	}
	return identifyCameraMakeNormalized(normalized[:n])
}

func identifyCameraMakeNormalized(normalized []byte) CameraMake {
	switch string(normalized) {
	case "agfa":
		return CameraMakeAgfa
	case "amazon":
		return CameraMakeAmazon
	case "apple":
		return CameraMakeApple
	case "arri":
		return CameraMakeARRI
	case "asahioptical":
		return CameraMakeAsahiOpticalCoLtd
	case "asahiopticalcoltd":
		return CameraMakeAsahiOpticalCoLtd
	case "asus":
		return CameraMakeASUS
	case "autelrobotics":
		return CameraMakeAutelRobotics
	case "blackmagicdesign":
		return CameraMakeBlackmagicDesign
	case "canon":
		return CameraMakeCanon
	case "casio":
		return CameraMakeCasio
	case "casiocomputercoltd":
		return CameraMakeCasio
	case "contax":
		return CameraMakeContax
	case "cosina":
		return CameraMakeCosina
	case "dji":
		return CameraMakeDJI
	case "dxo":
		return CameraMakeDXO
	case "eastmankodakcompany":
		return CameraMakeKodak
	case "epson":
		return CameraMakeSeikoEpsonCorp
	case "fujifilm":
		return CameraMakeFujifilm
	case "google":
		return CameraMakeGoogle
	case "gopro":
		return CameraMakeGoPro
	case "hasselblad":
		return CameraMakeHasselblad
	case "hewlett-packard":
		return CameraMakeHewlettPackard
	case "honor":
		return CameraMakeHonor
	case "hp":
		return CameraMakeHewlettPackard
	case "htc":
		return CameraMakeHTC
	case "huawei":
		return CameraMakeHuawei
	case "insta360":
		return CameraMakeInsta360
	case "jkimaging":
		return CameraMakeJKImagingLtd
	case "jkimagingltd":
		return CameraMakeJKImagingLtd
	case "kodak":
		return CameraMakeKodak
	case "konicaminolta":
		return CameraMakeKonicaMinolta
	case "konicaminoltacamerainc":
		return CameraMakeKonicaMinolta
	case "kyocera":
		return CameraMakeKyocera
	case "leica":
		return CameraMakeLeica
	case "leicacameraag":
		return CameraMakeLeica
	case "lg":
		return CameraMakeLGElectronics
	case "lge":
		return CameraMakeLGElectronics
	case "lgelectronics":
		return CameraMakeLGElectronics
	case "mamiya":
		return CameraMakeMamiya
	case "minolta":
		return CameraMakeMinoltaCoLtd
	case "minoltacoltd":
		return CameraMakeMinoltaCoLtd
	case "motorola":
		return CameraMakeMotorola
	case "nikon":
		return CameraMakeNikon
	case "nikoncorporation":
		return CameraMakeNikon
	case "nokia":
		return CameraMakeNokia
	case "olympus":
		return CameraMakeOlympusCorporation
	case "olympuscorporation":
		return CameraMakeOlympusCorporation
	case "olympusimagingcorp":
		return CameraMakeOlympusCorporation
	case "olympusopticalcoltd":
		return CameraMakeOlympusCorporation
	case "omdigitalsolutions":
		return CameraMakeOMDigitalSolutions
	case "oneplus":
		return CameraMakeOnePlus
	case "oppo":
		return CameraMakeOppo
	case "panasonic":
		return CameraMakePanasonic
	case "parrot":
		return CameraMakeParrot
	case "pentax":
		return CameraMakePentax
	case "pentaxcorporation":
		return CameraMakePentax
	case "pentaxricohimaging":
		return CameraMakePentax
	case "phaseone":
		return CameraMakePhaseOne
	case "polaroid":
		return CameraMakePolaroid
	case "realme":
		return CameraMakeRealme
	case "red":
		return CameraMakeRED
	case "ricoh":
		return CameraMakeRicoh
	case "ricohimagingcompanyltd":
		return CameraMakeRicoh
	case "rollei":
		return CameraMakeRollei
	case "samsung":
		return CameraMakeSamsung
	case "samsungtechwin":
		return CameraMakeSamsung
	case "sealife":
		return CameraMakeSeaLife
	case "seikoepsoncorp":
		return CameraMakeSeikoEpsonCorp
	case "sigma":
		return CameraMakeSigma
	case "skydio":
		return CameraMakeSkydio
	case "sony":
		return CameraMakeSony
	case "vivo":
		return CameraMakeVivo
	case "vivitar":
		return CameraMakeVivitar
	case "xiaomi":
		return CameraMakeXiaomi
	case "xiaoyi":
		return CameraMakeXiaoyi
	case "yashica":
		return CameraMakeYashica
	case "yuneec":
		return CameraMakeYuneec
	case "zeiss":
		return CameraMakeZeiss
	case "zte":
		return CameraMakeZTE
	}
	return CameraMakeUnknown
}
