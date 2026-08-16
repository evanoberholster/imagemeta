package makernote

// CameraMake identifies a normalized camera manufacturer token parsed from
// IFD0 Make text.
type CameraMake uint8

// CameraMake values.
const (
	CameraMakeUnknown CameraMake = iota
	CameraMakeAcer
	CameraMakeAgfa
	CameraMakeAgfaPhoto
	CameraMakeAiptek
	CameraMakeAmazon
	CameraMakeApple
	CameraMakeARRI
	CameraMakeAsahiOpticalCoLtd
	CameraMakeASUS
	CameraMakeAutelRobotics
	CameraMakeBenQ
	CameraMakeBlackBerry
	CameraMakeBlackmagicDesign
	CameraMakeBushnell
	CameraMakeCanon
	CameraMakeCasio
	CameraMakeConcord
	CameraMakeContax
	CameraMakeCosina
	CameraMakeCreative
	CameraMakeDaisy
	CameraMakeDJI
	CameraMakeDigiLife
	CameraMakeDoCoMo
	CameraMakeDXO
	CameraMakeFLIR
	CameraMakeFly
	CameraMakeFujifilm
	CameraMakeGarmin
	CameraMakeGateway
	CameraMakeGeneralImaging
	CameraMakeGenius
	CameraMakeGoogle
	CameraMakeGoPro
	CameraMakeHasselblad
	CameraMakeHewlettPackard
	CameraMakeHitachi
	CameraMakeHonor
	CameraMakeHTC
	CameraMakeHuawei
	CameraMakeInsta360
	CameraMakeJKImagingLtd
	CameraMakeJVC
	CameraMakeJenoptik
	CameraMakeKDDI
	CameraMakeKodak
	CameraMakeKonica
	CameraMakeKonicaMinolta
	CameraMakeKyocera
	CameraMakeLegend
	CameraMakeLeica
	CameraMakeLGElectronics
	CameraMakeLogitech
	CameraMakeLumicron
	CameraMakeMaginon
	CameraMakeMamiya
	CameraMakeMedion
	CameraMakeMercury
	CameraMakeMicrosoft
	CameraMakeMinoltaCoLtd
	CameraMakeMotorola
	CameraMakeMoultrie
	CameraMakeMustek
	CameraMakeNEC
	CameraMakeNikon
	CameraMakeNintendo
	CameraMakeNokia
	CameraMakeNoritsu
	CameraMakeODYS
	CameraMakeOMG
	CameraMakeOlympusCorporation
	CameraMakeOMDigitalSolutions
	CameraMakeOnePlus
	CameraMakeOppo
	CameraMakeOregonScientific
	CameraMakePackardBell
	CameraMakePanasonic
	CameraMakePantech
	CameraMakeParrot
	CameraMakePentacon
	CameraMakePhaseOne
	CameraMakePentax
	CameraMakePolaroid
	CameraMakePraktica
	CameraMakeRealme
	CameraMakeRED
	CameraMakeReconyx
	CameraMakeResearchInMotion
	CameraMakeRicoh
	CameraMakeRollei
	CameraMakeSagem
	CameraMakeSamsung
	CameraMakeSanyo
	CameraMakeSeaLife
	CameraMakeSeikoEpsonCorp
	CameraMakeSharp
	CameraMakeSigma
	CameraMakeSipix
	CameraMakeSkanhex
	CameraMakeSkydio
	CameraMakeSony
	CameraMakeSprint
	CameraMakeSunplus
	CameraMakeToshiba
	CameraMakeTraveler
	CameraMakeTrust
	CameraMakeUMAX
	CameraMakeUniden
	CameraMakeVivo
	CameraMakeVivitar
	CameraMakeVistaQuest
	CameraMakeVodafone
	CameraMakeWWL
	CameraMakeXiaomi
	CameraMakeXiaoyi
	CameraMakeYashica
	CameraMakeYuneec
	CameraMakeZeiss
	CameraMakeZTE
)

var cameraMakeNames = [...]string{
	"Unknown",
	"Acer",
	"Agfa",
	"AgfaPhoto",
	"Aiptek",
	"Amazon",
	"Apple",
	"ARRI",
	"Asahi Optical",
	"ASUS",
	"Autel Robotics",
	"BenQ",
	"BlackBerry",
	"Blackmagic Design",
	"Bushnell",
	"Canon",
	"Casio",
	"Concord",
	"Contax",
	"Cosina",
	"Creative",
	"Daisy",
	"DJI",
	"DigiLife",
	"DoCoMo",
	"DxO",
	"FLIR",
	"Fly",
	"Fujifilm",
	"Garmin",
	"Gateway",
	"General Imaging",
	"Genius",
	"Google",
	"GoPro",
	"Hasselblad",
	"HP",
	"Hitachi",
	"Honor",
	"HTC",
	"Huawei",
	"Insta360",
	"JK Imaging",
	"JVC",
	"Jenoptik",
	"KDDI",
	"Kodak",
	"Konica",
	"Konica Minolta",
	"Kyocera",
	"Legend",
	"Leica",
	"LG",
	"Logitech",
	"Lumicron",
	"Maginon",
	"Mamiya",
	"Medion",
	"Mercury",
	"Microsoft",
	"Minolta",
	"Motorola",
	"Moultrie",
	"Mustek",
	"NEC",
	"Nikon",
	"Nintendo",
	"Nokia",
	"Noritsu",
	"ODYS",
	"OMG",
	"Olympus",
	"OM Digital Solutions",
	"OnePlus",
	"OPPO",
	"Oregon Scientific",
	"Packard Bell",
	"Panasonic",
	"Pantech",
	"Parrot",
	"Pentacon",
	"Phase One",
	"Pentax",
	"Polaroid",
	"Praktica",
	"Realme",
	"RED",
	"Reconyx",
	"Research In Motion",
	"Ricoh",
	"Rollei",
	"Sagem",
	"Samsung",
	"Sanyo",
	"SeaLife",
	"Epson",
	"Sharp",
	"Sigma",
	"Sipix",
	"Skanhex",
	"Skydio",
	"Sony",
	"Sprint",
	"Sunplus",
	"Toshiba",
	"Traveler",
	"Trust",
	"UMAX",
	"Uniden",
	"vivo",
	"Vivitar",
	"VistaQuest",
	"Vodafone",
	"WWL",
	"Xiaomi",
	"Xiaoyi",
	"Yashica",
	"Yuneec",
	"Zeiss",
	"ZTE",
}

// String returns the display name for the camera make value.
func (m CameraMake) String() string {
	i := int(m)
	if i < 0 || i >= len(cameraMakeNames) {
		return cameraMakeNames[CameraMakeUnknown]
	}
	return cameraMakeNames[i]
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
	case "acer", "acercorporation", "acerinc":
		return CameraMakeAcer
	case "agfa", "agfagevaert":
		return CameraMakeAgfa
	case "agfaphoto":
		return CameraMakeAgfaPhoto
	case "aiptek", "aiptekinternationalinc":
		return CameraMakeAiptek
	case "amazon":
		return CameraMakeAmazon
	case "apple":
		return CameraMakeApple
	case "arri":
		return CameraMakeARRI
	case "asahioptical", "asahiopticalcoltd":
		return CameraMakeAsahiOpticalCoLtd
	case "asus":
		return CameraMakeASUS
	case "autelrobotics":
		return CameraMakeAutelRobotics
	case "benq", "benqcorporation", "benq-siemens", "benq_e72":
		return CameraMakeBenQ
	case "blackberry":
		return CameraMakeBlackBerry
	case "blackmagicdesign":
		return CameraMakeBlackmagicDesign
	case "bushnell":
		return CameraMakeBushnell
	case "canon", "canoninc":
		return CameraMakeCanon
	case "casio", "casiocomputercoltd":
		return CameraMakeCasio
	case "concordcameracorp", "concordcameragmbh", "concordcorporation":
		return CameraMakeConcord
	case "contax":
		return CameraMakeContax
	case "cosina":
		return CameraMakeCosina
	case "creative", "creativelabs", "creativetechnologyltd":
		return CameraMakeCreative
	case "daisymultimedia", "daisymultimedialtd":
		return CameraMakeDaisy
	case "dji":
		return CameraMakeDJI
	case "digilife":
		return CameraMakeDigiLife
	case "docomo", "nttdocomo":
		return CameraMakeDoCoMo
	case "dxo":
		return CameraMakeDXO
	case "eastmankodak", "eastmankodakcompany", "kodak":
		return CameraMakeKodak
	case "epson", "seikoepsoncorp":
		return CameraMakeSeikoEpsonCorp
	case "flir", "flirsystems", "flirsystemsab", "teledyneflir":
		return CameraMakeFLIR
	case "fly":
		return CameraMakeFly
	case "fujifilm", "fujiphotofilmcoltd":
		return CameraMakeFujifilm
	case "garmin":
		return CameraMakeGarmin
	case "gateway", "(c)2003gatewayinc":
		return CameraMakeGateway
	case "ge", "gedscimagingcorp", "generalimagingco":
		return CameraMakeGeneralImaging
	case "genius", "kyesystemscorp":
		return CameraMakeGenius
	case "google":
		return CameraMakeGoogle
	case "gopro":
		return CameraMakeGoPro
	case "hasselblad":
		return CameraMakeHasselblad
	case "hewlettpackard", "hewlett-packard", "hewlett-packardco", "hewlett-packardcompany", "hp", "hpphotosmart120":
		return CameraMakeHewlettPackard
	case "hitachi", "hitachilivingsystemsltd":
		return CameraMakeHitachi
	case "honor":
		return CameraMakeHonor
	case "ht", "htc", "htc-8900", "htc-p4600":
		return CameraMakeHTC
	case "huawei":
		return CameraMakeHuawei
	case "arashivision", "insta360":
		return CameraMakeInsta360
	case "jkimaging", "jkimagingltd":
		return CameraMakeJKImagingLtd
	case "jvc", "victor":
		return CameraMakeJVC
	case "jenimage", "jenimageeuropegmbh", "jenoptified", "jenoptik", "jenoptikcameragmbh", "jenoptikopticalcoltd":
		return CameraMakeJenoptik
	case "kddi-aa", "kddi-ca", "kddi-hi", "kddi-kc", "kddi-ma", "kddi-sa", "kddi-se", "kddi-sh", "kddi-sn", "kddi-st", "kddi-ts":
		return CameraMakeKDDI
	case "konica", "konicacoltd", "konicacorporation":
		return CameraMakeKonica
	case "konicaminolta", "konicaminoltacamerainc", "konicaminoltaphotoimaginginc":
		return CameraMakeKonicaMinolta
	case "kyocera":
		return CameraMakeKyocera
	case "legend", "legenddsc", "legendgrouplimited":
		return CameraMakeLegend
	case "leica", "leicacameraag":
		return CameraMakeLeica
	case "lg", "lgcyon", "lgelec", "lge", "lgelectronics", "lgelectronicsinc", "lgmobile", "lg_electronics":
		return CameraMakeLGElectronics
	case "logitechinc":
		return CameraMakeLogitech
	case "lumicron", "lumicrontechnologyinc":
		return CameraMakeLumicron
	case "maginon", "maginonopticalcoltd":
		return CameraMakeMaginon
	case "mamiya", "mamiya-opcoltd":
		return CameraMakeMamiya
	case "medion", "medion5mpdigitcam", "medionag", "medionopticalcoltd":
		return CameraMakeMedion
	case "mercuryperipheralsinc":
		return CameraMakeMercury
	case "microsoft", "microsoftmobile":
		return CameraMakeMicrosoft
	case "minolta", "minoltacoltd":
		return CameraMakeMinoltaCoLtd
	case "motorola", "motorol", "motorolakoreainc", "motorolamobility":
		return CameraMakeMotorola
	case "moultrie":
		return CameraMakeMoultrie
	case "mustek":
		return CameraMakeMustek
	case "nec":
		return CameraMakeNEC
	case "fs-nikon", "nikon", "nikoncorporation":
		return CameraMakeNikon
	case "nintendo":
		return CameraMakeNintendo
	case "hmdglobal", "nokia":
		return CameraMakeNokia
	case "noritsu", "noritsukoki":
		return CameraMakeNoritsu
	case "odys", "odyscorp":
		return CameraMakeODYS
	case "omglife":
		return CameraMakeOMG
	case "olympus", "olympuscorp", "olympuscorporation", "olympusimagingcorp", "olympusopticalcoltd", "olympus_imaging_corp":
		return CameraMakeOlympusCorporation
	case "omdigitalsolutions":
		return CameraMakeOMDigitalSolutions
	case "oneplus":
		return CameraMakeOnePlus
	case "oppo":
		return CameraMakeOppo
	case "oregonscientific":
		return CameraMakeOregonScientific
	case "packardbell":
		return CameraMakePackardBell
	case "panasonic":
		return CameraMakePanasonic
	case "pantech", "pantechwirelessinc":
		return CameraMakePantech
	case "parrot":
		return CameraMakeParrot
	case "pentacon", "pentacongermany":
		return CameraMakePentacon
	case "phaseone":
		return CameraMakePhaseOne
	case "pentax", "pentaxcorporation", "pentaxricohimaging":
		return CameraMakePentax
	case "madebypolaroid", "polaroid", "polaroidpdc1050", "polaroidpdc6350":
		return CameraMakePolaroid
	case "praktica":
		return CameraMakePraktica
	case "realme":
		return CameraMakeRealme
	case "red":
		return CameraMakeRED
	case "reconyx":
		return CameraMakeReconyx
	case "researchinmotion":
		return CameraMakeResearchInMotion
	case "ricoh", "ricohimaging", "ricohimagingcompanyltd":
		return CameraMakeRicoh
	case "rollei", "rolleifototechnicgmbh":
		return CameraMakeRollei
	case "sagem":
		return CameraMakeSagem
	case "samsung", "samsunganycall", "samsungcorporation", "samsungelec", "samsungelectronics", "samsungelectronicscoltd", "samsungopticalcoltd", "samsungtechwin", "samsungtechwinco", "samsungtechwincoltd":
		return CameraMakeSamsung
	case "sanyo", "sanyoelectriccoltd":
		return CameraMakeSanyo
	case "sealife":
		return CameraMakeSeaLife
	case "sharp":
		return CameraMakeSharp
	case "sigma":
		return CameraMakeSigma
	case "sipixcorporation":
		return CameraMakeSipix
	case "skanhex", "skanhexopticalcoltd", "skanhextech", "skanhextechnologyinc", "skanhextechwincoltd":
		return CameraMakeSkanhex
	case "skydio":
		return CameraMakeSkydio
	case "semc", "sony", "sonycomputerentertainmentinc", "sonyericsson", "sonycorporation":
		return CameraMakeSony
	case "sprint":
		return CameraMakeSprint
	case "sunplus", "sunplustechnologycoltd":
		return CameraMakeSunplus
	case "toshiba", "toshibacorporation":
		return CameraMakeToshiba
	case "traveler", "traveleropticalcoltd":
		return CameraMakeTraveler
	case "trust", "trustcomputerproducts":
		return CameraMakeTrust
	case "umax":
		return CameraMakeUMAX
	case "unidencorporation":
		return CameraMakeUniden
	case "vivo":
		return CameraMakeVivo
	case "vivitar", "vivicam", "vivitarcorporation":
		return CameraMakeVivitar
	case "viewquest", "vistaquest":
		return CameraMakeVistaQuest
	case "t-mobile(r)", "t-mobileshadow", "vodafone", "vodafonegroup", "vodafone710", "vodafonev720", "vodafonev810":
		return CameraMakeVodafone
	case "wwl", "wwlltd", "wwlcorporation":
		return CameraMakeWWL
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
