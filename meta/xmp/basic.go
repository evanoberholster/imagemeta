package xmp

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

var (
	xmpMMDerivedFrom = NewProperty(XmpMMNS, DerivedFrom)
	xapMMDerivedFrom = NewProperty(XapMMNS, DerivedFrom)
)

func (basic *Basic) parse(p property) (err error) {
	switch p.Property().Name() {
	case CreateDate:
		basic.CreateDate, err = parseDate(p.Value())
	case CreatorTool:
		basic.CreatorTool = parseString(p.Value())
	case Label:
		basic.Label = parseString(p.Value())
	case MetadataDate:
		basic.MetadataDate, err = parseDate(p.Value())
	case ModifyDate:
		basic.ModifyDate, err = parseDate(p.Value())
	case Rating:
		basic.Rating = parseInt8(p.Value())
	case XMPToolkit:
		basic.Toolkit = parseString(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return
}

func (mm *XMPMM) parse(p property) (err error) {
	switch p.Property().Name() {
	case DocumentID:
		if p.Namespace() == StRefNS && (p.Parent().Equals(xmpMMDerivedFrom) || p.Parent().Equals(xapMMDerivedFrom)) {
			mm.DerivedFromDocumentID = parseUUID(p.Value())
			return nil
		}
		mm.DocumentID = parseUUID(p.Value())
	case OriginalDocumentID:
		if p.Namespace() == StRefNS && (p.Parent().Equals(xmpMMDerivedFrom) || p.Parent().Equals(xapMMDerivedFrom)) {
			mm.DerivedFromOriginalDocumentID = parseUUID(p.Value())
			return nil
		}
		mm.OriginalDocumentID = parseUUID(p.Value())
	case DerivedFromDocumentID:
		mm.DerivedFromDocumentID = parseUUID(p.Value())
	case DerivedFromOriginalDocumentID:
		mm.DerivedFromOriginalDocumentID = parseUUID(p.Value())
	case PreservedFileName:
		mm.PreservedFileName = parseString(p.Value())
	case InstanceID:
		mm.InstanceID = parseUUID(p.Value())
	case HistoryTag:
		return mm.parseHistory(p)
	default:
		return ErrPropertyNotSet
	}
	return
}

func (mm *XMPMM) parseHistory(p property) (err error) {
	if p.Parent().Namespace() != StEvtNS {
		return ErrPropertyNotSet
	}

	h := mm.ensureHistory()
	switch p.Parent().Name() {
	case Action:
		h.Action = parseString(p.Value())
		mm.HistoryAction = h.Action
	case Changed:
		h.Changed = parseString(p.Value())
		mm.HistoryChanged = h.Changed
	case InstanceID:
		h.InstanceID = parseUUID(p.Value())
		mm.HistoryInstanceID = h.InstanceID
	case SoftwareAgent:
		h.Software = parseString(p.Value())
		mm.HistorySoftwareAgent = h.Software
	case When:
		h.Date, err = parseDate(p.Value())
		mm.HistoryWhen = h.Date
	case Parameters:
		h.Parameters = parseString(p.Value())
		mm.HistoryParameters = h.Parameters
	default:
		return ErrPropertyNotSet
	}
	return err
}

func (mm *XMPMM) ensureHistory() *History {
	if len(mm.History) == 0 {
		mm.History = append(mm.History, History{})
	}
	return &mm.History[0]
}

// Basic - the XMP basic namespace contains properties that provide basic descriptive information.
// XMP spec Section 8.4
// xmlns:xmp="http://ns.adobe.com/xap/1.0/"
type Basic struct {
	// Toolkit is x:xmptk from the XMP root element.
	Toolkit string
	// The date and time the resource was created. For a digital file, this need not match a
	// file-system  creation time. For a freshly created resource, it should be close to that time,
	// modulo the time taken to write the file. Later file transfer, copying, and so on, can make the
	// file-system time arbitrarily different.
	CreateDate time.Time `xml:"CreateDate"`
	// The name of the first known tool used to create the resource.
	CreatorTool string `xml:"CreatorTool"`
	// A word or short phrase that identifies a resource as a member of a user-defined collection.
	Label string `xml:"Label,attr"`
	// The date and time that any metadata for this resource was last changed.
	// It should be the same as or more recent than xmp:ModifyDate.
	MetadataDate time.Time `xml:"MetadataDate,attr"`
	// The date and time the resource was last modified.
	ModifyDate time.Time `xml:"ModifyDate,attr"`
	// A user-assigned rating for this file. The value shall be -1 or in the range [0..5],
	// where -1 indicates “rejected” and 0 indicates “unrated”. If xmp:Rating is not present,
	// a value of 0 should be assumed.
	Rating int8 `xml:"Rating,attr"`
}

// XMPMM - The XMP Media Management namespace contains properties that provide information
// regarding the identification, composition, and history of a resource.
// XMP spec Section 8.6
// Incomplete
type XMPMM struct {
	// DocumentId is the common identifier for all versions and renditions of a resource.
	DocumentID meta.UUID
	// InstanceId is an identifier for a specific incarnation of a resource,
	// updated each time a file is saved.
	InstanceID meta.UUID
	// OriginalDocumentId is the common identifier for the original resource from which the current
	// resource is derived. For example, if you save a resource to a different format,
	// then save that one to another format, each save operation should generate a new
	// xmpMM:DocumentID that uniquely identifies the resource in that format,
	// but should retain the ID of the source file here.
	OriginalDocumentID meta.UUID

	History []History

	// DerivedFromDocumentID identifies the direct source resource document ID.
	DerivedFromDocumentID meta.UUID
	// DerivedFromOriginalDocumentID identifies the original source document ID.
	DerivedFromOriginalDocumentID meta.UUID

	// ExifTool-compatible flattened history values (first event).
	HistoryAction        string
	HistoryChanged       string
	HistoryInstanceID    meta.UUID
	HistoryWhen          time.Time
	HistorySoftwareAgent string
	HistoryParameters    string

	PreservedFileName string
}

// History is an XMPMM History sequence
type History struct {
	Changed    string
	Action     string
	InstanceID meta.UUID
	Date       time.Time
	Software   string // softwareAgent
	Parameters string
}
