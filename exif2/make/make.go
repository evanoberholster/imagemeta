package make

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
	_strCameraMake = [51]string{"", "Acer", "Agfa", "Aiptek", "Apple", "Asus", "BenQ", "Canon", "Casio", "DJI", "FujiFilm", "Ge", "Genius", "Google", "GoPro", "Hasselblad", "HP", "Hitachi", "HTC", "Huawei", "Insta360", "Kodak", "Konica", "Kyocera", "Leica", "LG", "Mamyia", "Microsoft", "Minolta", "Motorola", "Nikon", "Nokia", "Olympus", "OnePlus", "Panasonic", "Pentax", "PhaseOne", "Polaroid", "RIM", "Ricoh", "Samsung", "Sanyo", "Sharp", "Sigma", "Sony", "SonyEricsson", "Toshiba", "Vivitar", "Xiamoi", "ZTE", "Hisilicon"}
)

func (cm CameraMake) String() string {
	if int(cm) < len(_strCameraMake) {
		return _strCameraMake[cm]
	}
	return _strCameraMake[0]
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
