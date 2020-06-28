package xmlmeta

import "time"

// TODO: Create Unmarshal for XMP Info

//xmlns:xmp="http://ns.adobe.com/xap/1.0/"
//xmlns:dc="http://purl.org/dc/elements/1.1/"
//xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
//xmlns:exifEX="http://cipa.jp/exif/1.0/"
//xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
//xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
//xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
//xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
//xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
//xmlns:darktable="http://darktable.sf.net/"

func (xpckt XMPPacket) Rating() float32 {
	return xpckt.XMP.Rating
}

func (xpckt XMPPacket) ModifyDate() time.Time {
	return xpckt.XMP.ModifyDate
}

func (xpckt XMPPacket) CreateDate() time.Time {
	return xpckt.XMP.CreateDate
}

func (xpckt XMPPacket) CreatorTool() string {
	return xpckt.XMP.CreatorTool
}

func (xpckt XMPPacket) Label() string {
	return xpckt.XMP.Label
}

func (xpckt XMPPacket) Description() []string {
	return xpckt.DC.Description
}

func (xpckt XMPPacket) Creator() []string {
	return xpckt.DC.Creator
}
