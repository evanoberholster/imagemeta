package meta

import (
	"bytes"
	"strconv"
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
		mm := NewMeteringMode(uint8(i))
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

func BenchmarkShutterSpeed(b *testing.B) {
	for _, bm := range ssList {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				//bm.ss.toBytes()
				bm.ss.UnmarshalText([]byte(bm.str))

			}
		})
	}
}

var ssList = []struct {
	name string
	ss   ShutterSpeed
	str  string
}{
	{"1", ShutterSpeed{1, 1}, "1.0"},
	{"2", ShutterSpeed{1, 50}, "1/50"},
	{"3", ShutterSpeed{1, 250}, "1/250"},
	{"4", ShutterSpeed{6, 10}, "0.6"},
	{"5", ShutterSpeed{1, 4000}, "1/4000"},
	{"6", ShutterSpeed{2, 1}, "2.0"},
	{"7", ShutterSpeed{30, 1}, "30.0"},
	{"8", ShutterSpeed{13, 10}, "1.3"},
	{"8", ShutterSpeed{12, 10}, "1.2"},
	{"9", ShutterSpeed{0, 0}, "0"},
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

		b1 := ShutterSpeed{}
		if err = b1.UnmarshalText([]byte(bm.str)); err != nil {
			t.Error(err)
		}

		m, _ := b1.MarshalText()
		assert.Equal(t, bm.ss, b1, "UnmarshalText #%s, wanted: %s got: %s", bm.name, bm.str, m)

	}
}

var fmList = []struct {
	name string
	fm   FlashMode
	str  string
	m    uint8
	b    bool
}{
	{"test0", NoFlash, "No Flash", 0, false},
	{"test1", FlashFired, "Fired", 1, true},
	{"test5", FlashFired, "Fired", 5, true},
	{"test7", FlashFired, "Fired", 7, true},
	{"test8", FlashOffNotFired, "Off, Did not fire", 8, false},
	{"test16", FlashOffNotFired, "Off, Did not fire", 16, false},
	{"test20", FlashOffNotFired, "Off, Did not fire", 20, false},
	{"test24", FlashAutoNotFired, "Auto, Did not fire", 24, false},
	{"test88", FlashAutoNotFired, "Auto, Did not fire", 88, false},
	{"test25", FlashAutoFired, "Auto, Fired", 25, true},
	{"test29", FlashAutoFired, "Auto, Fired", 29, true},
	{"test31", FlashAutoFired, "Auto, Fired", 31, true},
}

func TestFlashMode(t *testing.T) {
	for _, fm := range fmList {
		if a := NewFlashMode(fm.m); a != fm.fm {
			assert.Equal(t, fm.fm, a, "FlashMode.ParseFlashMode #%s, wanted %s got %s", fm.name, fm.fm, a)
		}

		// TextMarshall
		txt, err := fm.fm.MarshalText()
		if err != nil {
			t.Errorf("Error FlashMode.MarshallText (%s): %s ", fm.name, err.Error())
		}
		b := NewFlashMode(fm.m)
		assert.Equal(t, []byte(strconv.Itoa(int(b))), txt, "FlashMode.MarshalText #%s, wanted %d got %d", fm.name, fm.m, strconv.Itoa(int(b)))

		// UnmarshalText
		if err = b.UnmarshalText([]byte(strconv.Itoa(int(fm.m)))); err != nil {
			t.Errorf("Error FlashMode.UnmarshalText (%s): %s ", fm.name, err.Error())
		}
		assert.Equal(t, fm.fm, b, "FlashMode.UnmarshalText #%s, wanted %d got %d", fm.name, fm.fm, b)

		// String
		assert.Equal(t, fm.str, fm.fm.String(), "FlashMode.String #%s, wanted %s got %s", fm.name, fm.str, fm.fm.String())

		// Bool
		assert.Equal(t, fm.b, fm.fm.Bool(), "FlashMode.Bool #%s, wanted %s got %s", fm.name, fm.fm, fm.fm.Bool())
	}

	// Incorrect FlashMode
	a := FlashMode(250)
	b := NewFlashMode(250)
	assert.Equal(t, a.Bool(), false, "FlashMode.Bool #%s, wanted %s got %s", "Incorrect FlashMode", false, a.Bool())
	assert.NotEqual(t, a, b, "FlashMode.ParseFlashMode #%s, wanted %s got %s", "Incorrect FlashMode", b, a)
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
	ss := NewShutterSpeed(7, 5)
	testSerial(t, &ss)

	eb := NewExposureBias(0, 0)
	testSerial(t, &eb)

	em := NewExposureMode(9)
	testSerial(t, &em)

	ep := NewExposureProgram(8)
	testSerial(t, &ep)

	fm := NewFlashMode(8)
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
