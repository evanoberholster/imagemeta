package meta

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinylib/msgp/msgp"
)

func TestFocalLength(t *testing.T) {
	fl1 := FocalLength(100.25)
	fl2 := NewFocalLength(0, 0)

	buf, err := fl1.MarshalText()
	if err != nil {
		t.Error(err)
	}

	err = fl2.UnmarshalText(buf)
	if err != nil {
		t.Error(err)
	}

	if fl1.String() != "100.25mm" {
		t.Errorf("Incorrect FocalLength.String wanted %s got %s", fl1, fl2)
	}

	if fl1 != fl2 {
		t.Errorf("Incorrect FocalLength.MarshalText wanted %s got %s", fl1, fl2)
	}
	// Insufficient Buffer Length
	err = fl2.UnmarshalText([]byte(""))
	if err != nil {
		t.Errorf("Incorrect Error FocalLength.UnmarshalText wanted %s got %s", "nil", err)
	}

}

func TestMeteringMode(t *testing.T) {
	for i := 0; i < 256; i++ {
		mm := NewMeteringMode(uint16(i))
		mm2 := NewMeteringMode(0)

		buf, err := mm.MarshalText()
		if err != nil {
			t.Error(err)
		}

		err = mm2.UnmarshalText(buf)
		if err != nil {
			t.Error(err)
		}

		if mm2.String() != mm.String() {
			t.Errorf("Incorrect MeteringMode.String wanted %s got %s", mm2, mm)
		}

		if mm != mm2 {
			t.Errorf("Incorrect MeteringMode.MarshallText wanted %s (%d) got %s (%d)", mm, uint8(mm), mm2, uint8(mm2))
		}
		if i > 7 {
			if i == 255 {
				if mm.String() != "Other" {
					t.Errorf("Incorrect MeteringMode.String wanted %s got %s", "Other", mm)
				}
			} else {
				if mm.String() != "Unknown" {
					t.Errorf("Incorrect MeteringMode.String wanted %s got %s", "Unknown", mm)
				}
			}

		}
	}

	if MeteringMode(200).String() != "Unknown" {
		t.Errorf("Incorrect MeteringMode.String wanted %s got %s", "Unknown", NewMeteringMode(200).String())
	}
}

func TestExposureProgram(t *testing.T) {
	ep := ExposureProgram(1)
	str := "Manual"
	if ep.String() != str {
		t.Errorf("Incorrect ExposureProgram.String wanted %s got %s", str, ep.String())
	}
}

func TestExposureMode(t *testing.T) {
	items := []struct {
		str string
		em  ExposureMode
	}{
		{"Auto", 0},
		{"Manual", 1},
		{"Auto bracket", 2},
	}
	for _, v := range items {
		if v.em.String() != v.str {
			t.Errorf("Incorrect ExposureMode.String wanted %s got %s", v.str, v.em.String())
		}
	}

}

func BenchmarkExposureTime(b *testing.B) {
	for _, bm := range ssList {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				//bm.ss.toBytes()
				_, _ = bm.ss.MarshalText()

			}
		})
	}
}

var ssList = []struct {
	name string
	ss   ExposureTime
	str  string
}{
	{"1", ExposureTime(float32(1) / float32(1)), "1.00"},
	{"2", ExposureTime(float32(1) / float32(50)), "1/50"},
	{"3", ExposureTime(float32(1) / float32(250)), "1/250"},
	{"4", ExposureTime(float32(4) / float32(100)), "1/25"},
	{"5", ExposureTime(float32(1) / float32(4000)), "1/4000"},
	{"6", ExposureTime(float32(2) / float32(1)), "2.00"},
	{"7", ExposureTime(float32(30) / float32(1)), "30.00"},
	{"8", ExposureTime(float32(13) / float32(10)), "1.30"},
	{"8", ExposureTime(float32(12) / float32(10)), "1.20"},
	{"9", ExposureTime(0.0), ""},
}

func TestShutterSpeed(t *testing.T) {
	for _, bm := range ssList {
		b, err := bm.ss.MarshalText()
		if err != nil {
			t.Error(err)
		}

		if string(b) != bm.str {
			t.Errorf("Incorrect ShutterSpeed.MarshallText wanted %s got %s ", bm.str, b)
		}

		if bm.ss.String() != bm.str {
			t.Errorf("Incorrect ShutterSpeed.String wanted %s got %s ", bm.str, b)
		}

		//b1 := ExposureTime(0.0)
		//if err = b1.UnmarshalText([]byte(bm.str)); err != nil {
		//	t.Error(err)
		//////}

		//m, _ := b1.MarshalText()
		//assert.Equal(t, bm.ss, b1, "UnmarshalText #%s, wanted: %s got: %s", bm.name, bm.str, m)

	}
}

var ebList = []struct {
	name string
	eb   ExposureBias
	str  string
	n    int16
	d    int16
}{
	{"test0", ExposureBias(259), "+1/3", 1, 3},
	{"test1", ExposureBias(515), "+2/3", 2, 3},
	{"test2", ExposureBias(-1021), "-4/3", -4, 3},
	{"test3", ExposureBias(0), "0/0", 0, 0},
	{"test4", ExposureBias(1283), "+5/3", 5, 3},
}

func TestExposureBias(t *testing.T) {
	for _, eb := range ebList {
		testEB := ExposureBias(0)
		err := testEB.UnmarshalText([]byte(eb.str))
		if testEB != eb.eb || err != nil {
			t.Errorf("ExposureBias.MarshalText #%s, wanted %d got %d", eb.name, eb.eb, testEB)
		}

		// TextMarshall
		_, err = eb.eb.MarshalText()
		if err != nil {
			t.Errorf("Error ExposureBias.MarshallText (%s): %s ", eb.name, err.Error())
		}
		assert.Equal(t, eb.eb, testEB, "ExposureBias.MarshalText #%s, wanted %d got %d", eb.name, eb.eb, testEB)

		if eb.str != eb.eb.String() {
			t.Errorf("ExposureBias.String #%s, wanted %s got %s", eb.name, eb.str, eb.eb.String())
		}
		testEB = NewExposureBias(eb.n, eb.d)
		if testEB != eb.eb {
			t.Errorf("NewExposureBias #%s, wanted %s got %s", eb.name, testEB.String(), eb.eb.String())
		}
	}
	str := "6/3"
	eb2 := NewExposureBias(6, 3)
	if eb2.String() != "+"+str {
		t.Errorf("NewExposureBias #%s, wanted %s got %s", "Unsigned test", eb2.String(), "+"+str)
	}
}

func BenchmarkExposureBias(b *testing.B) {
	for _, eb := range ebList {
		b.Run(eb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				eb.eb.MarshalText()
			}
		})
	}
	//for _, eb := range ebList {
	//	b.Run(eb.name, func(b *testing.B) {
	//		b.ReportAllocs()
	//		b.ResetTimer()
	//		for i := 0; i < b.N; i++ {
	//			testEB := ExposureBias(0)
	//			testEB.UnmarshalText([]byte(eb.str))
	//		}
	//	})
	//}
}

////
// Automated Tests for MessagePack
////

func TestApertureMsgPack(t *testing.T) {
	v := NewAperture(10, 5)

	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Log("WARNING: TestAperture Msgsize() is inaccurate")
	}

	vn := NewAperture(0, 0)
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, &v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}

	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

type MsgPackInterface interface {
	//Encode(w io.Writer, e msgp.Encodable) error
	//Decode(r io.Reader, d msgp.Decodable) error
	DecodeMsg(dc *msgp.Reader) (err error)
	EncodeMsg(en *msgp.Writer) (err error)
	MarshalMsg(b []byte) (o []byte, err error)
	UnmarshalMsg(bts []byte) (o []byte, err error)
	//Skip(b []byte) ([]byte, error)
	Msgsize() (s int)
}

func TestMsgPack(t *testing.T) {
	// Aperture
	a := NewAperture(7, 5)
	testSerial(t, &a)

	// ShutterSpeed
	ss := ExposureTime(float32(7) / float32(5))
	testSerial(t, &ss)

	eb := NewExposureBias(0, 0)
	testSerial(t, &eb)

	em := NewExposureMode(9)
	testSerial(t, &em)

	ep := NewExposureProgram(8)
	testSerial(t, &ep)

	fm := NewFlash(8)
	testSerial(t, &fm)

	fl := NewFocalLength(8, 5)
	testSerial(t, &fl)

	mm := NewMeteringMode(10)
	testSerial(t, &mm)
}

func testSerial(t *testing.T, v MsgPackInterface) {
	var buf bytes.Buffer
	msgp.Encode(&buf, v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Log("WARNING: TestAperture Msgsize() is inaccurate")
	}

	err := msgp.Decode(&buf, v)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	msgp.Encode(&buf, v)
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}

	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

// flashTestList contains test data for Flash
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html#Flash (23/09/2019)
var flashTestList = []struct {
	flash       Flash
	flashFired  bool
	noFn        bool
	flashReturn FlashMode
	redEye      bool
	flashMode   FlashMode
	str         string
}{
	{0, false, false, FlashModeNone, false, FlashModeNone, "No Flash"},
	{1, true, false, FlashModeNone, false, FlashModeNone, "Fired"},
	{5, true, false, FlashNoReturn, false, FlashModeNone, "Fired, Return not detected"},
	{7, true, false, FlashReturn, false, FlashModeNone, "Fired, Return detected"},
	{8, false, false, FlashModeNone, false, FlashModeOn, "On, Did not fire"},
	{9, true, false, FlashModeNone, false, FlashModeOn, "On, Fired"},
	{13, true, false, FlashNoReturn, false, FlashModeOn, "On, Return not detected"},
	{15, true, false, FlashReturn, false, FlashModeOn, "On, Return detected"},
	{16, false, false, FlashModeNone, false, FlashModeOff, "Off, Did not fire"},
	{20, false, false, FlashNoReturn, false, FlashModeOff, "Off, Did not fire, Return not detected"},
	{24, false, false, FlashModeNone, false, FlashModeAuto, "Auto, Did not fire"},
	{25, true, false, FlashModeNone, false, FlashModeAuto, "Auto, Fired"},
	{29, true, false, FlashNoReturn, false, FlashModeAuto, "Auto, Fired, Return not detected"},
	{31, true, false, FlashReturn, false, FlashModeAuto, "Auto, Fired, Return detected"},
	{32, false, true, FlashModeNone, false, FlashModeNone, "No flash function"},
	{48, false, true, FlashModeNone, false, FlashModeOff, "Off, No flash function"},
	{65, true, false, FlashModeNone, true, FlashModeNone, "Fired, Red-eye reduction"},
	{69, true, false, FlashNoReturn, true, FlashModeNone, "Fired, Red-eye reduction, Return not detected"},
	{71, true, false, FlashReturn, true, FlashModeNone, "Fired, Red-eye reduction, Return detected"},
	{73, true, false, FlashModeNone, true, FlashModeOn, "On, Red-eye reduction"},
	{77, true, false, FlashNoReturn, true, FlashModeOn, "On, Red-eye reduction, Return not detected"},
	{79, true, false, FlashReturn, true, FlashModeOn, "On, Red-eye reduction, Return detected"},
	{80, false, false, FlashModeNone, true, FlashModeOff, "Off, Red-eye reduction"},
	{88, false, false, FlashModeNone, true, FlashModeAuto, "Auto, Did not fire, Red-eye reduction"},
	{89, true, false, FlashModeNone, true, FlashModeAuto, "Auto, Fired, Red-eye reduction"},
	{93, true, false, FlashNoReturn, true, FlashModeAuto, "Auto, Fired, Red-eye reduction, Return not detected"},
	{95, true, false, FlashReturn, true, FlashModeAuto, "Auto, Fired, Red-eye reduction, Return detected"},
}

func TestFlash(t *testing.T) {
	for _, f := range flashTestList {
		// Test Fired
		if f.flash.Fired() != f.flashFired {
			t.Errorf("Incorrect Flash Fired on %d wanted %v got %v", f.flash, f.flashFired, f.flash.Fired())
		}
		// Test NoFunction
		if f.flash.FlashFunction() != f.noFn {
			t.Errorf("Incorrect Flash Function on %d wanted %v got %v", f.flash, f.noFn, f.flash.FlashFunction())
		}
		// Test Return and NoReturn
		if f.flash.ReturnStatus() != f.flashReturn {
			t.Errorf("Incorrect Flash Return Status on %d wanted %v got %v", f.flash, f.flashReturn, f.flash.ReturnStatus())
		}
		// Test Redeye
		if f.flash.Redeye() != f.redEye {
			t.Errorf("Incorrect Flash Red-eye on %d wanted %v got %v", f.flash, f.redEye, f.flash.Redeye())
		}
		// Test FlashMode
		if f.flash.Mode() != f.flashMode {
			t.Errorf("Incorrect Flash Mode on %d wanted %v got %v", f.flash, f.flashMode, f.flash.Mode())
		}
		// Test Stringer
		if f.flash.String() != f.str {
			t.Errorf("Incorrect Flash String on %d wanted %v got %v", f.flash, f.str, f.flash.String())
		}
	}

	// Test TextMarshall
	// Test TextUnMarshall

}
