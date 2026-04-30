package canon

import "testing"

func TestParseBatteryType(t *testing.T) {
	withHeader := func(payload []byte) []byte {
		raw := make([]byte, BatteryTypePayloadSize)
		copy(raw[:BatteryTypeHeaderLen], []byte{0xde, 0xad, 0xbe, 0xef})
		copy(raw[BatteryTypeHeaderLen:], payload)
		return raw
	}

	tests := []struct {
		name    string
		payload []byte
		want    string
	}{
		{
			name:    "too short",
			payload: make([]byte, 10),
			want:    "",
		},
		{
			name:    "starts with NUL",
			payload: withHeader([]byte{0}),
			want:    "",
		},
		{
			name:    "LP-E6N",
			payload: withHeader([]byte("LP-E6N\x00TRAILING")),
			want:    "LP-E6N",
		},
		{
			name:    "LP-E6",
			payload: withHeader([]byte("LP-E6\x00")),
			want:    "LP-E6",
		},
		{
			name:    "LP-E6NH",
			payload: withHeader([]byte("LP-E6NH\x00")),
			want:    "LP-E6NH",
		},
		{
			name:    "LP-E6P",
			payload: withHeader([]byte("LP-E6P\x00")),
			want:    "LP-E6P",
		},
		{
			name:    "LP-E12",
			payload: withHeader([]byte("LP-E12\x00")),
			want:    "LP-E12",
		},
		{
			name:    "LP-E17",
			payload: withHeader([]byte("LP-E17\x00")),
			want:    "LP-E17",
		},
		{
			name:    "LP-E19",
			payload: withHeader([]byte("LP-E19\x00")),
			want:    "LP-E19",
		},
		{
			name:    "NB-13L",
			payload: withHeader([]byte("NB-13L\x00")),
			want:    "NB-13L",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ParseBatteryType(tc.payload[BatteryTypeHeaderLen:])
			if got != tc.want {
				t.Fatalf("got = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNormalizeFirmwareVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"strips prefix", "Firmware Version 1.0.1", "1.0.1"},
		{"passthrough", "1.0.1", "1.0.1"},
		{"trim space", "  Firmware Version 2.0.0  ", "2.0.0"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := NormalizeFirmwareVersion(tc.input)
			if got != tc.want {
				t.Fatalf("got = %q, want %q", got, tc.want)
			}
		})
	}
}
