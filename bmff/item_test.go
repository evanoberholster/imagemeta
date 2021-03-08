package bmff

import (
	"bufio"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockLogger struct{}

func (ml MockLogger) Debug(format string, args ...interface{}) {}

func TestParseItemInfoBox(t *testing.T) {
	ml := MockLogger{}
	DebugLogger(ml)
	array1 := []ItemInfoEntry{
		{ItemID: 513, ProtectionIndex: 0, size: 21, ItemType: ItemTypeMime},
		//{ItemID: 0, ProtectionIndex: 0, size: 0, ItemType: ItemTypeMime},
		//{ItemID: 0, ProtectionIndex: 0, size: 0, ItemType: ItemTypeExif},
	}
	array2 := make([]ItemInfoEntry, 5)
	iinfTests := []struct {
		name string
		data []byte
		val  ItemInfoBox
		err  error
	}{
		{"Test1", []byte{0, 0, 0, 10, 'i', 'i', 'n', 'f'}, ItemInfoBox{}, io.EOF},
		{"Test2", []byte{0, 0, 0, 18, 'i', 'i', 'n', 'f', 0, 0, 0, 3, 0, 5, 0, 0, 18, 'a', 'n'}, ItemInfoBox{array2}, ErrBufLength},
		{"Test3", []byte{0, 0, 0, 18, 'i', 'i', 'n', 'f', 0, 0, 0, 3, 0, 5, 0, 0, 18, 'i', 'n', 'f', 'd', 0, 0, 0, 0, 0}, ItemInfoBox{array2}, io.EOF},                                                              // CloseInnerBox Error
		{"Test4", []byte{0, 0, 0, 18, 'i', 'i', 'n', 'f', 0, 0, 0, 3, 0, 5, 0, 0, 0, 2, 'i', 'n', 'f', 'e', 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ItemInfoBox{array2}, io.EOF},                     // io.EOF
		{"Test5", []byte{0, 0, 0, 18, 'i', 'i', 'n', 'f', 0, 0, 0, 3, 0, 5, 0, 0, 1, 0, 'i', 'n', 'f', 'e', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ItemInfoBox{array2}, ErrInfeVersionNotSupported}, // ErrVersionNotSupported
		{"Test6", []byte{0, 0, 0, 39, 'i', 'i', 'n', 'f', 0, 0, 0, 3, 0, 1, 0, 0, 0, 21, 'i', 'n', 'f', 'e', 2, 0, 0, 0, 2, 1, 0, 0, 'm', 'i', 'm', 'e', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ItemInfoBox{array1}, bufio.ErrNegativeCount},
	}
	// Read ItemInfoEntry: [4]flags, [2]ItemID, [2]ProtectionIndex, [5]ItemType
	//infeHeaderSize := 13
	// Read IinfHeader [4]Flags, [2]ItemCount
	for _, v := range iinfTests {
		outer := newTestBox(v.data)
		inner, err := outer.readInnerBox()
		if err != nil {
			t.Error(err)
		}
		iinf, err := inner.parseItemInfoBox()
		if assert.ErrorIs(t, err, v.err, v.name) {
			assert.Equal(t, iinf, v.val, v.name)
		}
	}

	// General Tests
	iinf := ItemInfoBox{}
	if iinf.Type() != TypeIinf {
		t.Errorf("(iinf) BoxType expected %s got %s", TypeIinf, iinf.Type())
	}
	iinf2 := ItemInfoBox{array2}
	assert.NotEqual(t, iinf.String(), iinf2.String())
}
func TestParseItemLocationBox(t *testing.T) {

	// General Tests
	iloc := ItemLocationBox{}
	assert.Equal(t, TypeIloc, iloc.Type())

	iloc2 := ItemLocationBox{ItemCount: 10} // Add more tests here
	assert.NotEqual(t, iloc.String(), iloc2.String())
}
func TestParseItemPropertiesBox(t *testing.T) {

	// General Tests
	iprp := ItemPropertiesBox{}
	assert.Equal(t, TypeIprp, iprp.Type())

	iprp2 := ItemPropertiesBox{Associations: ItemPropertyAssociation{Entries: make([]ItemPropertyAssociationItem, 2)}}
	assert.NotEqual(t, iprp.String(), iprp2.String())
}
func TestParseItemPropertyContainerBox(t *testing.T) {

	// General Tests
	ipco := ItemPropertyContainerBox{}
	assert.Equal(t, TypeIpco, ipco.Type())
}
func TestParseItemPropertyAssociation(t *testing.T) {

	// General Tests
	ipma := ItemPropertyAssociation{}
	assert.Equal(t, ipma.Type(), TypeIpma)
}
func TestParseImageSpatialExtentsProperty(t *testing.T) {

	// General Tests
	ispe := ImageSpatialExtentsProperty{}
	assert.Equal(t, ispe.Type(), TypeIspe)

	ispe2 := ImageSpatialExtentsProperty{W: 1024, H: 1080}
	assert.NotEqual(t, ispe.String(), ispe2.String())
}
