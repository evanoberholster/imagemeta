package ifds

//go:generate stringer -type=CameraMake

// CameraMake is Camera Make found in Exif
type CameraMake uint16

const (
	CameraMakeUnknown CameraMake = iota
	Acer
	Agfa
	Aiptek
	Apple
	Asus
	BenQ
	Canon
	Casio
	DJI
	FujiFilm
	Ge
	Genius
	Google
	GoPro
	Hasselblad
	HP
	Hitachi
	HTC
	Huawei
	Insta360
	Kodak
	Konica
	Kyocera
	Leica
	LG
	Mamyia
	Microsoft
	Minolta
	Motorola
	Nikon
	Nokia
	Olympus
	OnePlus
	Panasonic
	Pentax
	PhaseOne
	Polaroid
	RIM
	Ricoh
	Samsung
	Sanyo
	Sharp
	Sigma
	Sony
	SonyEricsson
	Toshiba
	Vivitar
	Xiamoi
	ZTE

	//Additions
	Hisilicon
)

var (
	_strCameraMake       = "AcerAgfaAiptekAppleAsusBenQCanonCasioDJIFujiFilmGeGeniusGoogleGoProHasselbladHPHitachiHTCHuaweiInsta360KodakKonicaKyoceraLeicaLGMamyiaMicrosoftMinoltaMotorolaNikonNokiaOlympusOnePlusPanasonicPentaxPhaseOnePolaroidRIMRicohSamsungSanyoSharpSigmaSonySonyEricssonToshibaVivitarXiamoiZTEHisilicon"
	_strCameraMakeOffset = []uint16{0, 0, 4, 8, 14, 19, 23, 27, 32, 37, 40, 48, 50, 56, 62, 67, 77, 79, 86, 89, 95, 103, 108, 114, 121, 126, 128, 134, 143, 150, 158, 163, 168, 175, 182, 191, 197, 205, 213, 216, 221, 228, 233, 238, 243, 247, 259, 266, 273, 279, 282, 291}
)

func (cm CameraMake) String() string {
	if cm <= Hisilicon {
		return _strCameraMake[_strCameraMakeOffset[cm]:_strCameraMakeOffset[cm+1]]
	}
	return ""
}

// CameraMakeFromString returns a camera make from the given string
func CameraMakeFromString(str string) (CameraMake, bool) {
	if cm, ok := mapStringCameraMake[str]; ok {
		return cm, true
	}
	return CameraMakeUnknown, false
}

var mapStringCameraMake = map[string]CameraMake{
	"":                  CameraMakeUnknown,
	"Acer":              Acer,
	"Agfa":              Agfa,
	"Aiptek":            Aiptek,
	"Apple":             Apple,
	"Asus":              Asus,
	"BenQ":              BenQ,
	"Canon":             Canon,
	"Casio":             Casio,
	"DJI":               DJI,
	"FujiFilm":          FujiFilm,
	"Ge":                Ge,
	"Genius":            Genius,
	"Google":            Google,
	"GoPro":             GoPro,
	"Hasselblad":        Hasselblad,
	"HP":                HP,
	"Hitachi":           Hitachi,
	"HTC":               HTC,
	"HUAWEI":            Huawei,
	"Insta360":          Insta360,
	"Kodak":             Kodak,
	"Konica":            Konica,
	"Kyocera":           Kyocera,
	"Leica":             Leica,
	"LG":                LG,
	"Mamyia":            Mamyia,
	"Microsoft":         Microsoft,
	"Minolta":           Minolta,
	"Motorola":          Motorola,
	"Nikon":             Nikon,
	"NIKON CORPORATION": Nikon,
	"Nokia":             Nokia,
	"Olympus":           Olympus,
	"OnePlus":           OnePlus,
	"Panasonic":         Panasonic,
	"Pentax":            Pentax,
	"PhaseOne":          PhaseOne,
	"Polaroid":          Polaroid,
	"RIM":               RIM,
	"Ricoh":             Ricoh,
	"Samsung":           Samsung,
	"Sanyo":             Sanyo,
	"Sharp":             Sharp,
	"Sigma":             Sigma,
	"Sony":              Sony,
	"SONY":              Sony,
	"SonyEricsson":      SonyEricsson,
	"Toshiba":           Toshiba,
	"Vivitar":           Vivitar,
	"Xiamoi":            Xiamoi,
	"ZTE":               ZTE,
	"Hisilicon":         Hisilicon,
}
