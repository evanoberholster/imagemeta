package xmp

import (
	"bufio"
	"errors"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/evanoberholster/imagemeta/meta/jpeg"
	"github.com/orisano/gosax"
)

const (
	parserBufferSize = 64 * 1024
	scanBufferSize   = 4 * 1024
	minPeekSize      = 64
)

var rdfRootProperty = NewProperty(RdfNS, RDF)
var mwgRegionListProperty = NewProperty(MwgRSNS, RegionListTag)
var xmpRootToken = []byte("<x:xmpmeta")
var rdfRootToken = []byte("<rdf:RDF")

var saxReaderPool = sync.Pool{
	New: func() interface{} {
		return gosax.NewReaderSize(nil, parserBufferSize)
	},
}

var scanReaderPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReaderSize(pooledReaderReset{}, scanBufferSize)
	},
}

func getPooledScanReader() *bufio.Reader {
	br, ok := scanReaderPool.Get().(*bufio.Reader)
	if !ok || br == nil {
		return bufio.NewReaderSize(pooledReaderReset{}, scanBufferSize)
	}
	return br
}

func getPooledSAXReader() *gosax.Reader {
	r, ok := saxReaderPool.Get().(*gosax.Reader)
	if !ok || r == nil {
		return gosax.NewReaderSize(nil, parserBufferSize)
	}
	return r
}

type pooledReaderReset struct{}

func (pooledReaderReset) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

// Parse reads XMP metadata from sidecar files and embedded image payloads.
//
// Supported embedded sources include JPEG APP1 packets and ISOBMFF files such as
// CR3, HEIC/HEIF, AVIF and JXL containers.
func Parse(r io.Reader) (XMP, error) {
	return ParseWithOptions(r, ParseOptions{})
}

// ParseWithOptions reads XMP metadata from sidecar files and embedded image payloads
// with additional parser options.
func ParseWithOptions(r io.Reader, opts ParseOptions) (XMP, error) {
	br, pooled := asBufferedReader(r)
	if pooled {
		defer releaseBufferedReader(br)
	}

	fileType, err := imagetype.ScanBuf(br)
	if err != nil {
		return ParseXmpWithOptions(br, opts)
	}

	switch fileType {
	case imagetype.ImageJPEG:
		return parseJPEG(br, opts)
	case imagetype.ImageCR3, imagetype.ImageHEIC, imagetype.ImageHEIF, imagetype.ImageAVIF:
		return parseISOBMFF(br, opts)
	case imagetype.ImageJXL:
		// JPEG XL may be codestream-based and not always wrapped in BMFF boxes.
		// Direct XMP stream parsing is robust for both sidecar-like and embedded packets.
		return ParseXmpWithOptions(br, opts)
	default:
		return ParseXmpWithOptions(br, opts)
	}
}

func parseJPEG(r io.Reader, opts ParseOptions) (XMP, error) {
	var out XMP
	found := false

	err := jpeg.ScanJPEG(
		r,
		nil,
		func(packet io.Reader) error {
			x, err := ParseXmpWithOptions(packet, opts)
			if err != nil {
				return err
			}
			out = x
			found = true
			return nil
		},
	)
	if err != nil {
		return XMP{}, err
	}
	if !found {
		return XMP{}, ErrNoXMP
	}
	return out, nil
}

func parseISOBMFF(r io.Reader, opts ParseOptions) (XMP, error) {
	var out XMP
	found := false

	bmffReader := isobmff.NewReader(
		r,
		nil,
		func(packet io.Reader, _ isobmff.XPacketHeader) error {
			x, err := ParseXmpWithOptions(packet, opts)
			if err != nil {
				return err
			}
			out = x
			found = true
			return nil
		},
		nil,
	)
	defer bmffReader.Close()

	if err := bmffReader.ReadFTYP(); err != nil {
		return XMP{}, err
	}

	for {
		err := bmffReader.ReadMetadata()
		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		return XMP{}, err
	}

	if !found {
		return XMP{}, ErrNoXMP
	}
	return out, nil
}

func asBufferedReader(r io.Reader) (*bufio.Reader, bool) {
	if br, ok := r.(*bufio.Reader); ok {
		if br.Size() >= minPeekSize {
			return br, false
		}
	}

	br := getPooledScanReader()
	br.Reset(r)
	return br, true
}

func releaseBufferedReader(br *bufio.Reader) {
	br.Reset(pooledReaderReset{})
	scanReaderPool.Put(br)
}

func parseXMPStream(r io.Reader, opts ParseOptions) (x XMP, err error) {
	br, pooled := asBufferedReader(r)
	if pooled {
		defer releaseBufferedReader(br)
	}

	found, err := seekXMPStart(br)
	if err != nil {
		return XMP{}, err
	}
	if !found {
		return XMP{}, ErrNoXMP
	}

	sax := getPooledSAXReader()
	sax.Reset(br)
	sax.EmitSelfClosingTag = true
	defer saxReaderPool.Put(sax)

	parser := xmpStreamParser{r: sax, debug: opts.Debug, regionIndex: -1}
	return parser.parse()
}

func seekXMPStart(br *bufio.Reader) (bool, error) {
	for {
		b1, err := br.Peek(1)
		if err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
		if b1[0] == '<' {
			if hasReaderPrefix(br, xmpRootToken) || hasReaderPrefix(br, rdfRootToken) {
				return true, nil
			}
		}

		if _, err := br.ReadByte(); err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
	}
}

func hasReaderPrefix(br *bufio.Reader, token []byte) bool {
	buf, err := br.Peek(len(token))
	if err != nil && len(buf) < len(token) {
		return false
	}
	for i := 0; i < len(token); i++ {
		if buf[i] != token[i] {
			return false
		}
	}
	return true
}

type xmpStreamParser struct {
	r       *gosax.Reader
	xmp     XMP
	started bool
	debug   bool
	root    Property

	stack [64]Property
	depth int

	regionIndex   int16
	regionLiDepth int
}

func (p *xmpStreamParser) parse() (XMP, error) {
	for {
		ev, err := p.r.Event()
		if err != nil {
			return XMP{}, err
		}

		switch ev.Type() {
		case gosax.EventEOF:
			if !p.started {
				return XMP{}, ErrNoXMP
			}
			return p.xmp, nil
		case gosax.EventStart:
			p.parseStart(ev.Bytes)
		case gosax.EventEnd:
			if p.parseEnd(ev.Bytes) {
				return p.xmp, nil
			}
		case gosax.EventText:
			p.parseText(ev.Bytes)
		case gosax.EventCData:
			p.parseCData(ev.Bytes)
		}
	}
}

func (p *xmpStreamParser) parseStart(tag []byte) {
	tagName, attrs := gosax.Name(tag)
	tagProp := identifyProperty(tagName)

	if !p.started {
		if tagProp != XMPRootProperty && tagProp != rdfRootProperty {
			return
		}
		p.started = true
		p.root = tagProp
	}

	if p.isMWGRegionListItem(tagProp) {
		p.regionIndex++
		p.regionLiDepth = p.depth + 1
	}

	listTarget, hasListTarget := p.listTargetFor(tagProp)
	for len(attrs) > 0 {
		attr, next, err := gosax.NextAttribute(attrs)
		if err != nil {
			break
		}
		attrs = next
		if len(attr.Key) == 0 {
			break
		}
		if len(attr.Value) < 2 {
			continue
		}

		attrProp := identifyProperty(attr.Key)
		if attrProp == 0 {
			continue
		}

		parent := tagProp
		self := attrProp
		if hasListTarget {
			parent = attrProp
			self = listTarget
		}
		attrVal := unescapeXML(attr.Value[1 : len(attr.Value)-1])
		p.emitProperty(attrPType, parent, self, attrVal)
	}

	p.push(tagProp)
}

func (p *xmpStreamParser) parseEnd(tag []byte) bool {
	if !p.started {
		return false
	}
	tagName, _ := gosax.Name(tag)
	tagProp := identifyProperty(tagName)
	if tagProp == RDFLi && p.regionLiDepth == p.depth {
		p.regionLiDepth = 0
	}
	p.pop(tagProp)
	return tagProp == p.root
}

func (p *xmpStreamParser) parseText(text []byte) {
	if !p.started || p.depth == 0 {
		return
	}

	value := trimSpace(text)
	if len(value) == 0 {
		return
	}
	value = unescapeXML(value)

	self := p.stack[p.depth-1]
	var parent Property
	if p.depth >= 2 {
		parent = p.stack[p.depth-2]
	}

	if listTarget, ok := p.listTargetFor(self); ok {
		self = listTarget
	}

	p.emitProperty(tagPType, parent, self, value)
}

func (p *xmpStreamParser) emitProperty(pt pType, parent Property, self Property, value []byte) {
	prop := property{
		pt:          pt,
		parent:      parent,
		self:        self,
		val:         value,
		regionIndex: p.currentRegionIndex(),
	}
	if err := p.xmp.parser(prop, p.debug); err != nil {
		return
	}
}

func (p *xmpStreamParser) isMWGRegionListItem(prop Property) bool {
	return prop == RDFLi &&
		p.depth >= 2 &&
		p.stack[p.depth-1] == RDFSeq &&
		p.stack[p.depth-2] == mwgRegionListProperty
}

func (p *xmpStreamParser) currentRegionIndex() int16 {
	if p.regionLiDepth == 0 || p.depth < p.regionLiDepth {
		return -1
	}
	return p.regionIndex
}

func (p *xmpStreamParser) parseCData(cdata []byte) {
	// <![CDATA[value]]>
	if len(cdata) >= 12 {
		p.parseText(cdata[9 : len(cdata)-3])
		return
	}
	p.parseText(cdata)
}

func (p *xmpStreamParser) push(prop Property) {
	if p.depth >= len(p.stack) {
		return
	}
	p.stack[p.depth] = prop
	p.depth++
}

func (p *xmpStreamParser) pop(prop Property) {
	for p.depth > 0 {
		top := p.stack[p.depth-1]
		p.depth--
		if top == prop {
			return
		}
	}
}

func (p *xmpStreamParser) listTargetFor(tagProp Property) (Property, bool) {
	if tagProp != RDFLi {
		return 0, false
	}

	// Attribute parsing: current tag is not yet pushed on stack.
	if p.depth >= 2 && isListContainer(p.stack[p.depth-1]) {
		return p.stack[p.depth-2], true
	}

	// Text parsing: current tag is already on stack.
	if p.depth >= 3 && p.stack[p.depth-1] == RDFLi && isListContainer(p.stack[p.depth-2]) {
		return p.stack[p.depth-3], true
	}

	return 0, false
}

func isListContainer(prop Property) bool {
	switch prop {
	case RDFSeq, RDFBag, RDFAlt:
		return true
	default:
		return false
	}
}

func trimSpace(buf []byte) []byte {
	i := 0
	for i < len(buf) {
		switch buf[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			goto right
		}
	}
	return buf[:0]

right:
	j := len(buf) - 1
	for j >= i {
		switch buf[j] {
		case ' ', '\n', '\r', '\t':
			j--
		default:
			return buf[i : j+1]
		}
	}
	return buf[:0]
}

func unescapeXML(buf []byte) []byte {
	if len(buf) == 0 {
		return buf
	}
	if !containsXMLEscape(buf) {
		return buf
	}
	unescaped, err := gosax.Unescape(buf)
	if err != nil {
		return decodeXMLEntities(buf)
	}
	return unescaped
}

func containsXMLEscape(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		switch buf[i] {
		case '&', '\r':
			return true
		}
	}
	return false
}

func identifyProperty(name []byte) Property {
	colon := -1
	for i := 0; i < len(name); i++ {
		if name[i] == ':' {
			colon = i
			break
		}
	}
	if colon <= 0 || colon+1 >= len(name) {
		return 0
	}

	prefix := name[:colon]
	local := name[colon+1:]

	ns := IdentifyNamespace(prefix)
	n := IdentifyName(local)

	if ns == UnknownNS || n == UnknownPropertyName {
		return 0
	}

	return NewProperty(ns, n)
}
