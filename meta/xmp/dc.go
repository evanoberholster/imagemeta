package xmp

// Copyright (c) 2017-2018 Alexander Eichhorn
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

import (
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
)

var xmlLang = NewProperty(XMLNS, Lang)

func (dc *DublinCore) parse(p property) (err error) {
	if len(p.Value()) == 0 {
		return ErrPropertyNotSet
	}
	switch p.Name() {
	case Contributor:
		dc.Contributor = append(dc.Contributor, parseString(p.Value()))
	case Coverage:
		dc.Coverage = parseString(p.Value())
	case Format:
		dc.Format = imagetype.FromBytes(p.Value())
	case Creator:
		dc.Creator = append(dc.Creator, parseString(p.Value()))
	case Date:
		// ExifTool-style flattening: keep the first valid dc:date value.
		if dc.Date.IsZero() {
			if t, err := parseDate(p.Value()); err == nil {
				dc.Date = t
			}
		}
	case Identifier:
		dc.Identifier = parseString(p.Value())
	case Language:
		dc.Language = append(dc.Language, parseString(p.Value()))
	case Publisher:
		dc.Publisher = append(dc.Publisher, parseString(p.Value()))
	case Relation:
		dc.Relation = append(dc.Relation, parseString(p.Value()))
	case Source:
		dc.Source = parseString(p.Value())
	case Subject:
		dc.Subject = append(dc.Subject, parseString(p.Value()))
	case Rights:
		parseDublinCoreAltProperty(&dc.Rights, &dc.RightsLang, p)
	case Title:
		parseDublinCoreAltProperty(&dc.Title, &dc.TitleLang, p)
	case Description:
		parseDublinCoreAltProperty(&dc.Description, &dc.DescriptionLang, p)
	case Type:
		dc.Type = append(dc.Type, parseString(p.Value()))
	default:
		return ErrPropertyNotSet
	}
	return nil
}

func parseDublinCoreAltProperty(valueDst *[]string, langDst *[]string, p property) {
	// ExifTool defaults to x-default (or first) language item for RDF Alt.
	if p.pt == tagPType && len(*valueDst) == 0 {
		*valueDst = append(*valueDst, parseString(p.Value()))
	}
	if p.Parent().Equals(xmlLang) {
		*langDst = append(*langDst, parseString(p.Value()))
	}
}

// DublinCore is the "dc:" namespace often seen in xmp meta.
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
	// XMP usage is an ordered list of dates, flattened to the first valid value.
	Date time.Time `xml:"date"`
	// An account of the resource.
	// XMP usage is a list of textual descriptions of the content of the resource, given in various languages.
	Description []string `xml:"description"`
	// XMP usage is a MIME type.
	Format imagetype.ImageType `xml:"format"`
	// An unambiguous reference to the resource within a given context.
	Identifier string `xml:"identifier"`

	// A language of the resource.
	// XMP usage is a list (RDF Bag) of languages used in the resource.
	Language []string `xml:"language"`
	// An entity responsible for making the resource available
	// Examples of a publisher include a person, an organization, or a
	// service. Typically, the name of a publisher should be used to indicate the entity.
	// XMP usage is a list (RDF Bag) of publishers.
	Publisher []string `xml:"publisher"`
	// A related resource.
	// Recommended best practice is to identify the related resource
	// by means of a string conforming to a formal identification system.
	// XMP usage is a list (RDF Bag) of related resources.
	Relation []string `xml:"relation"`
	// Information about rights held in and over the resource.
	// typically, rights information includes a statement about various property
	// rights associated with the resource, including intellectual property rights.
	// XMP usage is a list (RDF Alt) of informal rights statements, given in various languages.
	Rights []string `xml:"rights"`
	// RightsLang stores xml:lang values for dc:rights alternatives.
	RightsLang []string
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
	Title     []string `xml:"title"`
	TitleLang []string
	// The nature or genre of the resource.
	// Recommended best practice is to use a controlled vocabulary such as the DCMI Type
	// Vocabulary [DCMITYPE]. To describe the file format, physical medium, or dimensions of the
	// resource, use the dc:format element.
	// See the dc:format entry for clarification of the XMP usage of that element.
	Type []string `xml:"type"`

	// DescriptionLang stores xml:lang values for dc:description alternatives.
	DescriptionLang []string
}
