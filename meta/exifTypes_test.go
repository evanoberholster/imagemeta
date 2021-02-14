package meta

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFocalLength(t *testing.T) {
	fl1 := FocalLength(100.25)
	fl2 := FocalLength(0)

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
		t.Errorf("Incorrect FocalLength.MarshallText wanted %s got %s", fl1, fl2)
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

	}
}

func TextExposureMode(t *testing.T) {

}

func BenchmarkShutterSpeed100(b *testing.B) {
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
		if a := ParseFlashMode(fm.m); a != fm.fm {
			assert.Equal(t, fm.fm, a, "FlashMode.ParseFlashMode #%s, wanted %s got %s", fm.name, fm.fm, a)
		}

		// TextMarshall
		txt, err := fm.fm.MarshalText()
		if err != nil {
			t.Errorf("Error FlashMode.MarshallText (%s): %s ", fm.name, err.Error())
		}
		b := ParseFlashMode(fm.m)
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
	b := ParseFlashMode(250)
	assert.Equal(t, a.Bool(), false, "FlashMode.Bool #%s, wanted %s got %s", "Incorrect FlashMode", false, a.Bool())
	assert.NotEqual(t, a, b, "FlashMode.ParseFlashMode #%s, wanted %s got %s", "Incorrect FlashMode", b, a)
}

var ebList = []struct {
	name string
	eb   ExposureBias
	str  string
}{
	{"test0", ExposureBias{1, 3}, "1/3"},
	{"test1", ExposureBias{2, 3}, "2/3"},
	{"test2", ExposureBias{-4, 3}, "-4/3"},
	{"test3", ExposureBias{1, 0}, "0/0"},
}

func TestExposureBias(t *testing.T) {
	for _, eb := range ebList {

		// TextMarshall
		txt, err := eb.eb.MarshalText()
		if err != nil {
			t.Errorf("Error ExposureBias.MarshallText (%s): %s ", eb.name, err.Error())
		}
		assert.Equal(t, []byte(eb.str), txt, "ExposureBias.MarshalText #%s, wanted %d got %d", eb.name, eb.eb, []byte(eb.str))

		// UnmarshalText
		//if err = b.UnmarshalText([]byte(strconv.Itoa(int(fm.m)))); err != nil {
		//	t.Errorf("Error FlashMode.UnmarshalText (%s): %s ", fm.name, err.Error())
		//}
		//assert.Equal(t, fm.fm, b, "FlashMode.UnmarshalText #%s, wanted %d got %d", fm.name, fm.fm, b)

		// String
		assert.Equal(t, eb.str, eb.eb.String(), "ExposureBias.String #%s, wanted %s got %s", eb.name, eb.str, eb.eb.String())
	}
}

func BenchmarkExposureBias100(b *testing.B) {
	for _, eb := range ebList {
		b.Run(eb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				eb.eb.MarshalText()
			}
		})
	}
}
