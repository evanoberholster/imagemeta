package xmp

const (
	regionAreaXIndex = iota
	regionAreaYIndex
	regionAreaWIndex
	regionAreaHIndex
)

// RegionType is the MWG region semantic type (for example "Face").
type RegionType string

// Dimensions stores the dimensions and unit used by mwg-rs:AppliedToDimensions.
type Dimensions struct {
	H    uint32
	W    uint32
	Unit string
}

// RegionExtensions stores optional apple-fi extension values associated with a region.
type RegionExtensions struct {
	AngleInfoRoll   float32
	AngleInfoYaw    float32
	ConfidenceLevel float32
	FaceID          string
	TimeStamp       uint64
}

// Region is one MWG region-list entry.
type Region struct {
	Name       string     `xml:"name"` // Person name if present.
	Role       string     `xml:"role"` // Optional role metadata.
	Type       RegionType `xml:"type"` // Region type, usually "Face".
	Area       RegionArea `xml:"area"` // [x, y, w, h] normalized/pixel values.
	Extensions RegionExtensions
}

type RegionArea struct {
	X float32
	Y float32
	W float32
	H float32
}

// RegionInfo stores MWG Regions metadata.
type RegionInfo struct {
	AppliedToDimensions Dimensions `xml:"appliedToDimensions"`
	RegionList          []Region   `xml:"regionList>Bag>li"`
}

func (ri *RegionInfo) parse(p property) error {
	switch p.Namespace() {
	case MwgRSNS:
		return ri.parseMWG(p)
	case StDimNS:
		return ri.parseStDim(p)
	case StAreaNS:
		return ri.parseStArea(p)
	case AppleFiNS:
		return ri.parseAppleFaceInfo(p)
	default:
		return ErrPropertyNotSet
	}
}

func (ri *RegionInfo) parseMWG(p property) error {
	switch p.Name() {
	case Type, RegionTypeTag:
		region := ri.regionForProperty(p)
		region.Type = RegionType(parseString(p.Value()))
	case NameTag:
		region := ri.regionForProperty(p)
		region.Name = parseString(p.Value())
	case RoleTag:
		region := ri.regionForProperty(p)
		region.Role = parseString(p.Value())
	case RegionAppliedToDimensionsH:
		ri.AppliedToDimensions.H = parseUint32(p.Value())
	case RegionAppliedToDimensionsW:
		ri.AppliedToDimensions.W = parseUint32(p.Value())
	case RegionAppliedToDimensionsUnit:
		ri.AppliedToDimensions.Unit = parseString(p.Value())
	case RegionAreaX:
		region := ri.regionForProperty(p)
		region.setAreaComponent(regionAreaXIndex, parseFloat64(p.Value()))
	case RegionAreaY:
		region := ri.regionForProperty(p)
		region.setAreaComponent(regionAreaYIndex, parseFloat64(p.Value()))
	case RegionAreaW:
		region := ri.regionForProperty(p)
		region.setAreaComponent(regionAreaWIndex, parseFloat64(p.Value()))
	case RegionAreaH:
		region := ri.regionForProperty(p)
		region.setAreaComponent(regionAreaHIndex, parseFloat64(p.Value()))
	case RegionAreaUnit:
		region := ri.regionForProperty(p)
		region.setAreaUnit(parseString(p.Value()))
	case RegionExtensionsAngleInfoRoll:
		region := ri.regionForProperty(p)
		region.Extensions.AngleInfoRoll = float32(parseFloat64(p.Value()))
	case RegionExtensionsAngleInfoYaw:
		region := ri.regionForProperty(p)
		region.Extensions.AngleInfoYaw = float32(parseFloat64(p.Value()))
	case RegionExtensionsConfidenceLevel:
		region := ri.regionForProperty(p)
		region.Extensions.ConfidenceLevel = float32(parseFloat64(p.Value()))
	case RegionExtensionsFaceID:
		region := ri.regionForProperty(p)
		region.Extensions.FaceID = parseString(p.Value())
	case RegionExtensionsTimeStamp:
		region := ri.regionForProperty(p)
		region.Extensions.TimeStamp = parseUint(p.Value())
	default:
		return ErrPropertyNotSet
	}

	return nil
}

func (ri *RegionInfo) parseStDim(p property) error {
	switch p.Name() {
	case H:
		ri.AppliedToDimensions.H = parseUint32(p.Value())
	case W:
		ri.AppliedToDimensions.W = parseUint32(p.Value())
	case Unit:
		ri.AppliedToDimensions.Unit = parseString(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return nil
}

func (ri *RegionInfo) parseStArea(p property) error {
	region := ri.regionForProperty(p)
	switch p.Name() {
	case X:
		region.setAreaComponent(regionAreaXIndex, parseFloat64(p.Value()))
	case Y:
		region.setAreaComponent(regionAreaYIndex, parseFloat64(p.Value()))
	case W:
		region.setAreaComponent(regionAreaWIndex, parseFloat64(p.Value()))
	case H:
		region.setAreaComponent(regionAreaHIndex, parseFloat64(p.Value()))
	case Unit:
		region.setAreaUnit(parseString(p.Value()))
	default:
		return ErrPropertyNotSet
	}
	return nil
}

func (ri *RegionInfo) parseAppleFaceInfo(p property) error {
	region := ri.regionForProperty(p)
	switch p.Name() {
	case RegionExtensionsAngleInfoRoll:
		region.Extensions.AngleInfoRoll = float32(parseFloat64(p.Value()))
	case RegionExtensionsAngleInfoYaw:
		region.Extensions.AngleInfoYaw = float32(parseFloat64(p.Value()))
	case RegionExtensionsConfidenceLevel:
		region.Extensions.ConfidenceLevel = float32(parseFloat64(p.Value()))
	case RegionExtensionsFaceID:
		region.Extensions.FaceID = parseString(p.Value())
	case RegionExtensionsTimeStamp:
		region.Extensions.TimeStamp = parseUint(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return nil
}

func (ri *RegionInfo) ensureRegion(index int) *Region {
	if index < 0 {
		index = 0
	}
	for len(ri.RegionList) <= index {
		ri.RegionList = append(ri.RegionList, Region{})
	}
	return &ri.RegionList[index]
}

func (ri *RegionInfo) regionForProperty(p property) *Region {
	if index := p.RegionIndex(); index >= 0 {
		return ri.ensureRegion(index)
	}

	if len(ri.RegionList) == 0 {
		ri.RegionList = append(ri.RegionList, Region{})
	}
	return &ri.RegionList[len(ri.RegionList)-1]
}

func (r *Region) setAreaComponent(index int, raw float64) {
	if index < 0 || index > regionAreaHIndex {
		return
	}

	value := float32(raw)
	switch index {
	case regionAreaXIndex:
		r.Area.X = value
	case regionAreaYIndex:
		r.Area.Y = value
	case regionAreaWIndex:
		r.Area.W = value
	case regionAreaHIndex:
		r.Area.H = value
	}
}

func (r *Region) setAreaUnit(unit string) {
	_ = unit // Unit is currently ignored. Area is stored as parsed float values.
}
