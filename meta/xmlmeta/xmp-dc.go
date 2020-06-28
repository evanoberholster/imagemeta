package xmlmeta

import (
	"encoding/xml"
	"time"
)

func (dc *DublinCore) decodeAttr(attr xml.Attr) (err error) {
	switch attr.Name.Local {
	case "format":
		dc.Format = attr.Value
	default:
		//fmt.Printf("unknown DC: %s: %s", attr.Name, attr.Value)
		//err = fmt.Errorf("unknown DC: %s: %s", attr.Name, attr.Value)
	}
	return err
}

func (dc *DublinCore) decode(decoder *xml.Decoder, start *xml.StartElement) (err error) {
	switch start.Name.Local {
	case "creator":
		dc.Creator = decodeRDF(decoder, start)
	case "rights":
		dc.Rights = decodeRDF(decoder, start)
	case "subject":
		dc.Subject = decodeRDF(decoder, start)
	case "contributor":
		dc.Contributor = decodeRDF(decoder, start)
	case "description":
		dc.Description = decodeRDF(decoder, start)
	case "title":
		dc.Title = decodeRDF(decoder, start)
	}
	return err
}

// DublinCore is the "dc:" namespace often seen in xmp meta. It's pretty sweet.
// https://en.wikipedia.org/wiki/Dublin_Core
// http://dublincore.org
// For the XMP flavour, see XMP section 8.3
//
// xmlns:dc="http://purl.org/dc/elements/1.1/"
type DublinCore struct {
	// An entity responsible for making contributions to the resource
	// Examples of a contributor include a person, an organization, or a service.
	// Typically, the name of a contributor should be used to indicate the entity.
	// XMP usage is a list of contributors. These contributors should not include those listed in dc:creator.
	Contributor []string `xml:"contributor"`
	// The spatial or temporal topic of the resource, the spatial applicability of the resource,
	// or the jurisdiction under which the resource is relevant.
	// XMP usage is the extent or scope of the resource.
	Coverage string `xml:"coverage"`
	// An entity primarily responsible for making the resource.
	// Examples of a creator include a person, an organization, or a
	// service. Typically, the name of a creator should be used to indicate the entity.
	// XMP usage is a list of creators. Entities should be listed in order of decreasing precedence,
	// if such order is significant.
	Creator []string `xml:"creator"`
	// A point or period of time associated with an event in the life cycle of the resource.
	Date time.Time `xml:"date"`
	// An account of the resource.
	// XMP usage is a list of textual descriptions of the content of the resource, given in various languages.
	Description []string `xml:"description"`
	// XMP usage is a MIME type.
	Format string `xml:"format"`
	// An unambiguous reference to the resource within a given context.
	Identifier string `xml:"identifier"`

	// A language of the resource.
	// XMP usage is a list of languages used in the content of the resource.
	// TODO - RDFSeq is a guess
	Language []string `xml:"language"`
	// An entity responsible for making the resource available
	// Examples of a publisher include a person, an organization, or a
	// service. Typically, the name of a publisher should be used to indicate the entity.
	//  XMP usage is a list of publishers.
	// TODO - RDFSeq is a guess
	//Publisher RDFSeq `xml:"publisher"`
	// A related resource.
	// Recommended best practice is to identify the related resource
	// by means of a string conforming to a formal identification system.
	// XMP usage is a list of related resources.
	// TODO - RDFSeq is a guess
	//Relation RDFSeq `xml:"relation"`
	// Information about rights held in and over the resource.
	// typically, rights information includes a statement about various property
	// rights associated with the resource, including intellectual property rights.
	// XMP usage is a list of informal rights statements, given in various languages.
	// TODO - RDFAlt is a guess
	Rights []string `xml:"rights"`
	// A related resource from which the described resource is derived.
	// The described resource may be derived from the related resource in whole or in part.
	// Recommended best practice is to identify the related resource by means of a string
	// conforming to a formal identification system.
	Source string `xml:"source"`
	// The topic of the resource.
	// Typically, the subject will be represented using keywords, key phrases, or
	// classification codes. Recommended best practice is to use a controlled vocabulary.
	// To describe the spatial or temporal topic of the resource, use the dc:coverage element.
	// XMP usage is a list of descriptive phrases or keywords that specify the content of the resource.
	Subject []string `xml:"subject"`
	// A name given to the resource.
	// Typically, a title will be a name by which the resource is formally known.
	// XMP usage is a title or name, given in various languages.
	Title []string `xml:"title"`
	// The nature or genre of the resource.
	// Recommended best practice is to use a controlled vocabulary such as the DCMI Type
	// Vocabulary [DCMITYPE]. To describe the file format, physical medium, or dimensions of the
	// resource, use the dc:format element.
	// See the dc:format entry for clarification of the XMP usage of that element.
	//Type RDFBag `xml:"type"`
}