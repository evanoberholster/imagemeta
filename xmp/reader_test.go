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

var readTagHeaderTests = []struct {
	data   string
	err    error
	tag    Tag
	assert bool
}{
	{"", io.EOF, Tag{}, false},
	{"         ", io.EOF, Tag{}, false},
	{"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\">", nil, Tag{property: property{parent: xmpns.XMPRootProperty, self: xmpns.NewProperty(xmpns.RdfNS, xmpns.RDF), pt: tagPType}, t: startTag}, true},
	{"<rdf:description/>", nil, Tag{property: property{parent: xmpns.XMPRootProperty, self: xmpns.NewProperty(xmpns.RdfNS, xmpns.Description), pt: tagPType}, t: soloTag}, true},
	{"\n" + "</rdf:Description>", nil, Tag{property: property{parent: xmpns.XMPRootProperty, self: xmpns.NewProperty(xmpns.RdfNS, xmpns.Description), pt: tagPType}, t: stopTag}, true},
	{"<hello >               ", ErrNegativeRead, Tag{}, false},
}

func TestReadTagHeader(t *testing.T) {
	var err error
	parentTestTag := Tag{}
	parentTestTag.self = xmpns.XMPRootProperty

	for _, tagTest := range readTagHeaderTests {
		var tag Tag
		br := bufReader{
			r: bufio.NewReaderSize(bytes.NewReader([]byte(tagTest.data)), 120),
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

var readAttributeTests = []struct{}{}

func TestReadAttribute(t *testing.T) {

}

var readAttrValue = []struct{}{}

func TestReadAttrValue(t *testing.T) {

}
