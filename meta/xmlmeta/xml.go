// Package xmlmeta provides types and functions for identifying and decoding
// image xml information as well as xmp files.
package xmlmeta

import (
	"encoding/xml"
	"io"
)

var (
	rdfSeq         = xml.Name{Space: "rdf", Local: "Seq"}
	rdfLi          = xml.Name{Space: "rdf", Local: "li"}
	rdfDescription = xml.Name{Space: "rdf", Local: "Description"}
)

func Walk(r io.Reader) (xmp XMPPacket, err error) {
	xmp = XMPPacket{}

	decoder := xml.NewDecoder(r)
	var t xml.Token
	for {
		if t, err = decoder.RawToken(); err != nil {
			if err == io.EOF {
				return xmp, nil
			}
			return
		}
		switch x := t.(type) {
		case xml.StartElement:
			switch x.Name.Space {
			case "rdf":
				xmp.XMPDescription(decoder, &x)
			case "dc":
				xmp.DC.decode(decoder, &x)
				//default:
				//fmt.Println(x.Name)
				//xmp.XMPDescription(decoder, &x)
				//fmt.Printf("Unexpected SE {%s}%s\n", x.Name.Space, x.Name.Local)
			}

			//case xml.EndElement:
			//	switch x.Name {
			//	default:
			//		//fmt.Printf("Unexpected EE {%s}%s\n", x.Name.Space, x.Name.Local)
			//	}
		}
	}
}

// XMPPacket is an Instance of the XMP Data Model.
// Also known as an XMP document.
type XMPPacket2 struct {
	RDF XMPRDFMeta `xml:"RDF"`
}

// XMP is itself a subset of the Resource Description Framework.
// RDF is an OG (Original Gangter) web (W3) standard for data interchange.
// https://www.w3.org/TR/rdf11-concepts/
// All XMPPackets are a single instance of an RDF document, with limitations
// defined by the XMP spec
type XMPRDFMeta struct {
	Description XMPDescription `xml:"Description"`
}

type XMPDescription struct {
	DublinCore
	XMPBasic
	//XMPRights
	//XMPMM
	//XMPIdq
}
