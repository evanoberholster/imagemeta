package xml

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

type namespace interface {
	decodeAttr(attr xml.Attr) (err error)
}

func (xmp *XMP) decodeNS(ns namespace, attr xml.Attr) {
	if err := ns.decodeAttr(attr); err != nil {
		// log error if debug
		if DebugMode {
			fmt.Println(err)
		}
	}
}

// XMPDescription -
func (xmp *XMP) XMPDescription(decoder *xml.Decoder, start *xml.StartElement) {
	for _, i := range start.Attr {
		switch i.Name.Space {
		case "aux":
			xmp.decodeNS(&xmp.Aux, i)
		case "exif":
			xmp.decodeNS(&xmp.Exif, i)
		case "tiff":
			//xmp.decodeNS(&xmp.Tiff, i)
		case "xmp":
			xmp.decodeNS(&xmp.Basic, i)
		case "dc":
			xmp.decodeNS(&xmp.DC, i)
		case "crs":
			xmp.decodeNS(&xmp.CRS, i)
		default:
			if DebugMode {
				fmt.Println("Unsupported: ", i.Name, i.Value)
			}
		}
	}
}

func (basic *Basic) decodeAttr(attr xml.Attr) (err error) {
	switch attr.Name.Local {
	case "CreateDate":
		basic.CreateDate, err = parseDate(attr.Value)
	case "CreatorTool":
		basic.CreatorTool = attr.Value
	case "Label":
		basic.Label = attr.Value
	case "MetadataDate":
		basic.MetadataDate, err = parseDate(attr.Value)
	case "ModifyDate":
		basic.ModifyDate, err = parseDate(attr.Value)
	case "Rating":
		var r float64
		r, err = strconv.ParseFloat(attr.Value, 32)
		basic.Rating = float32(r)
	default:
		err = fmt.Errorf("unknown: %s: %s", attr.Name, attr.Value)
	}
	return err
}

// Basic - the XMP basic namespace contains properties that provide basic descriptive information.
// XMP spec Section 8.4
// xmlns:xmp="http://ns.adobe.com/xap/1.0/"
type Basic struct {
	// The date and time the resource was created. For a digital file, this need not match a
	// file-system  creation time. For a freshly created resource, it should be close to that time,
	// modulo the time taken to write the file. Later file transfer, copying, and so on, can make the
	// file-system time arbitrarily different.
	CreateDate time.Time `xml:"CreateDate"`
	// The name of the first known tool used to create the resource.
	// TODO - the spec says this should be an "AgentName"
	CreatorTool string `xml:"CreatorTool"`
	// A word or short phrase that identifies a resource as a member of a user- defined collection.
	// NOTE: One anticipated usage is to organize resources in a file browser.
	Label string `xml:"Label,attr"`
	// The date and time that any metadata for this resource was last changed.
	// It should be the same as or more recent than xmp:ModifyDate.
	MetadataDate time.Time `xml:"MetadataDate,attr"`
	// The date and time the resource was last modified.
	// NOTE: The value of this property is not necessarily the same as the file’s
	// system modification date because it is typically set before the file is saved.
	ModifyDate time.Time `xml:"ModifyDate,attr"`
	// A user-assigned rating for this file. The value shall be -1 or in the range [0..5],
	// where -1 indicates “rejected” and 0 indicates “unrated”. If xmp:Rating is not present,
	// a value of 0 should be assumed.
	// NOTE: Anticipated usage is for a typical “star rating” UI, with the addition of a notion of rejection.
	Rating float32 `xml:"Rating,attr"`
}

// XMPMM - The XMP Media Management namespace contains properties that provide information
// regarding the identification, composition, and history of a resource.
// XMP spec Section 8.6
type XMPMM struct {
	// DerivedFrom is a reference to the resource from which this was derived.
	// This should be a minimal reference, in which missing components can be
	// assumed to be unchanged.
	// NOTE A rendition might need to specify only the
	// xmpMM:InstanceID and xmpMM:RenditionClass of the original.
	// TODO - this is actually a "ResourceRef" type
	DerivedFrom string
	// DocumentId is the common identifier for all versions and renditions of a resource.
	DocumentId string
	// InstanceId is an identifier for a specific incarnation of a resource,
	// updated each time a file is saved.
	// TODO - this is actually a GUID type
	InstanceId string
	// OriginalDocumentId is the common identifier for the original resource from which the current
	// resource is derived. For example, if you save a resource to a different format,
	// then save that one to another format, each save operation should generate a new
	// xmpMM:DocumentID that uniquely identifies the resource in that format,
	// but should retain the ID of the source file here.
	OriginalDocumentId string
	// The rendition class name for this resource. This property should be absent or
	// set to 'default' for a resource that is not a derived rendition.
	// See definitions of rendition (3.7) and version (3.9).
	// TODO - this is actually a RenditionClass type
	RenditionClass string
	// Can be used to provide additional rendition parameters that are
	// too complex or verbose to encode in xmpMM:RenditionClass.
	RenditionParams string
}
