package isobmff

func readIprp(b *box) (err error) {
	if logLevelInfo() {
		logInfoBox(b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeIpma:
			_, err = readIpma(&inner)
		case typeIpco:
			err = readIpco(&inner)
		default:
			if logLevelInfo() {
				logInfoBox(&inner).Send()
			}
		}
		if err != nil && logLevelError() {
			logError().Object("box", inner).Err(err).Send()
		}
		if err = inner.close(); err != nil && logLevelError() {
			logError().Object("box", inner).Err(err).Send()
		}
	}
	return b.close()
}

func readIpco(b *box) (err error) {
	if logLevelInfo() {
		logInfoBox(b).Send()
	}
	//var inner box
	//var ok bool
	//for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
	//	if logLevelInfo() {
	//		logInfoBox(inner)
	//	}
	//	inner.close()
	//}

	return b.close()
}

// ItemPropertiesBox is an ISOBMFF "iprp" box
type ItemPropertiesBox struct {
	PropertyContainer ItemPropertyContainerBox
	//Associations      []ItemPropertyAssociation // at least 1
	Associations ItemPropertyAssociation
}

type ItemPropertyContainerBox struct{}

// ItemPropertyAssociation is an ISOBMFF "ipma" box
type ItemPropertyAssociation struct {
	//Flags      Flags
	//EntryCount uint32
	Entries []IpmaItem
}

// ItemPropertyAssociationItem is not a box
type IpmaItem struct {
	ItemID       uint32
	Associations [6]uint16
	//AssociationsCount uint32 // as declared
	//Associations      []ItemProperty // as parsed
}

func readIpma(b *box) (ipma ItemPropertyAssociation, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		return
	}
	b.readFlagsFromBuf(buf[:4])
	count := int(bmffEndian.Uint32(buf[4:8]))
	if logLevelInfo() {
		logInfoBox(b).Uint32("entries", uint32(count)).Send()
	}
	// Entries
	// /ipma.Entries = make([]IpmaItem, count)
	// /if buf, err = b.Peek(b.remain); err != nil {
	// /	return
	// /}
	return ipma, b.close()
}
