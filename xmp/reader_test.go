package xmp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/evanoberholster/imagemeta/xmp/xmpns"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadTagHeader(t *testing.T) {
	var readTagHeaderTests = []struct {
		data   string
		err    error
		tag    Tag
		assert bool
	}{
		{"", io.EOF, Tag{}, false},
		{"         ", ErrBufferFull, Tag{}, false},
		{"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\">", nil, Tag{property: property{parent: xmpns.XMPRootProperty, self: xmpns.NewProperty(xmpns.RdfNamespace, xmpns.RDF), pt: tagPType}, t: startTag}, true},
		{"<rdf:description/>", nil, Tag{property: property{parent: xmpns.XMPRootProperty, self: xmpns.NewProperty(xmpns.RdfNamespace, xmpns.Description), pt: tagPType}, t: soloTag}, true},
		{"\n" + "</rdf:Description>", nil, Tag{property: property{parent: xmpns.XMPRootProperty, self: xmpns.NewProperty(xmpns.RdfNamespace, xmpns.Description), pt: tagPType}, t: stopTag}, true},
		{"<hello >               ", ErrNegativeRead, Tag{}, false},
		{"<? >               ", io.EOF, Tag{}, false},
	}

	var err error
	parentTestTag := Tag{}
	parentTestTag.self = xmpns.XMPRootProperty

	for _, tagTest := range readTagHeaderTests {
		var tag Tag
		br := xmpReader{
			r: bufio.NewReaderSize(bytes.NewReader([]byte(tagTest.data)), 1024),
		}
		tag, err = br.readTagHeader(parentTestTag)
		if err != nil {
			if errors.Cause(err) != tagTest.err {
				t.Error(err)
			}
		}

		if tagTest.assert {
			if !assert.Equal(t, tagTest.tag, tag, tagTest.tag.String()) {
				fmt.Println(tag, tagTest.data)
			}
		}
	}
}

// TestReadAttribute tests reading an attribute
// Todo: needs fixing and making more redundant.
func TestReadAttribute(t *testing.T) {
	//var readAttributeTests = []struct {
	//	data   string
	//	err    error
	//	attr   Attribute
	//	assert bool
	//}{
	//	{"", io.EOF, Attribute{property: property{}}, false},
	//	{"\n" + " xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\">", nil, Attribute{property: property{pt: attrPType, val: []byte("http://www.w3.org/1999/02/22-rdf-syntax-ns#"), self: xmpns.NewProperty(xmpns.XMLnsNS, xmpns.RDF)}}, true},
	//	{"  />", io.EOF, Attribute{property: property{}}, false},
	//	{"  abc:", ErrNegativeRead, Attribute{property: property{}}, false},
	//}
	//
	//var err error
	//for _, attrTest := range readAttributeTests {
	//	var attr Attribute
	//	xr := xmpReader{
	//		r: bufio.NewReaderSize(bytes.NewReader([]byte(attrTest.data)), 120),
	//		a: true,
	//	}
	//	if !xr.hasAttribute() {
	//		t.Error(errors.New("HasAttribute error"))
	//	}
	//	tagTest := Tag{}
	//
	//	if attr, err = xr.readAttribute(&tagTest); err != nil {
	//		if !errors.Is(err, attrTest.err) {
	//			t.Error(err)
	//		}
	//	}
	//
	//	if attrTest.assert {
	//		if !assert.Equal(t, attrTest.attr, attr, attrTest.attr.String()) {
	//			t.Errorf("Attribute error wanted: %s, got: %s", attr, attrTest.attr)
	//			//fmt.Println(attr, attrTest.data)
	//		}
	//	}
	//}
}
