package apple

import (
	"encoding/json"
	"testing"
)

func TestImageCaptureTypeString(t *testing.T) {
	tests := []struct {
		name     string
		input    ImageCaptureType
		expected string
	}{
		{"ProRAW", ImageCaptureProRAW, "ProRAW"},
		{"Portrait", ImageCapturePortrait, "Portrait"},
		{"Photo", ImageCapturePhoto, "Photo"},
		{"Manual Focus", ImageCaptureManualFocus, "Manual Focus"},
		{"Scene", ImageCaptureScene, "Scene"},
		{"Unknown", ImageCaptureType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("ImageCaptureType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCameraTypeString(t *testing.T) {
	tests := []struct {
		name     string
		input    CameraType
		expected string
	}{
		{"Back Wide Angle", CameraBackWideAngle, "Back Wide Angle"},
		{"Back Normal", CameraBackNormal, "Back Normal"},
		{"Front", CameraFront, "Front"},
		{"Unknown", CameraType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("CameraType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAFStateString(t *testing.T) {
	tests := []struct {
		name     string
		input    AFState
		expected string
	}{
		{"Stable", AFState(true), "Yes"},
		{"Unstable", AFState(false), "No"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("AFState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppleMakerNoteJSON(t *testing.T) {
	note := AppleMakerNote{
		MakerNoteVersion:  1,
		AEStable:          true,
		AETarget:          100,
		ImageCaptureType:  ImageCaptureProRAW,
		CameraType:        CameraBackWideAngle,
		BurstUUID:         "test-uuid",
		ContentIdentifier: "test-content-id",
	}

	data, err := json.Marshal(note)
	if err != nil {
		t.Fatalf("Failed to marshal AppleMakerNote: %v", err)
	}

	var decoded AppleMakerNote
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal AppleMakerNote: %v", err)
	}

	if decoded.MakerNoteVersion != note.MakerNoteVersion {
		t.Errorf("MakerNoteVersion = %v, want %v", decoded.MakerNoteVersion, note.MakerNoteVersion)
	}
	if decoded.ImageCaptureType != note.ImageCaptureType {
		t.Errorf("ImageCaptureType = %v, want %v", decoded.ImageCaptureType, note.ImageCaptureType)
	}
	if decoded.CameraType != note.CameraType {
		t.Errorf("CameraType = %v, want %v", decoded.CameraType, note.CameraType)
	}
}

func TestFocusDistanceRange(t *testing.T) {
	tests := []struct {
		name       string
		focusRange FocusDistanceRange
		expected   string
	}{
		{
			name:       "Normal order",
			focusRange: FocusDistanceRange{{120, 100}, {250, 100}}, // 1.2m, 2.5m
			expected:   "1.20 - 2.50 m",
		},
		{
			name:       "Reversed order",
			focusRange: FocusDistanceRange{{250, 100}, {120, 100}}, // 2.5m, 1.2m
			expected:   "1.20 - 2.50 m",
		},
		{
			name:       "Same values",
			focusRange: FocusDistanceRange{{100, 100}, {100, 100}}, // 1.0m, 1.0m
			expected:   "1.00 - 1.00 m",
		},
		{
			name:       "Zero values",
			focusRange: FocusDistanceRange{{0, 1}, {0, 1}},
			expected:   "0.00 - 0.00 m",
		},
		{
			name:       "Large numbers",
			focusRange: FocusDistanceRange{{10000, 1000}, {20000, 1000}}, // 10m, 20m
			expected:   "10.00 - 20.00 m",
		},
		{
			name:       "Fractional values",
			focusRange: FocusDistanceRange{{1, 3}, {2, 3}}, // ~0.33m, ~0.67m
			expected:   "0.33 - 0.67 m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.focusRange.String()
			if got != tt.expected {
				t.Errorf("FocusDistanceRange.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccelerationVectorString(t *testing.T) {
	tests := []struct {
		name     string
		vector   AccelerationVector
		expected string
	}{
		{
			name:     "Zero vector",
			vector:   AccelerationVector{{0, 1}, {0, 1}, {0, 1}},
			expected: "X:0.00, Y:0.00, Z:0.00",
		},
		{
			name:     "Unit vector",
			vector:   AccelerationVector{{1, 1}, {1, 1}, {1, 1}},
			expected: "X:1.00, Y:1.00, Z:1.00",
		},
		{
			name:     "Mixed values",
			vector:   AccelerationVector{{-100, 100}, {150, 100}, {-200, 100}},
			expected: "X:-1.00, Y:1.50, Z:-2.00",
		},
		{
			name:     "From raw values",
			vector:   AccelerationVectorfromRaw([]int32{100, 100, -100, 100, 200, 100}),
			expected: "X:1.00, Y:-1.00, Z:2.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.vector.String()
			if got != tt.expected {
				t.Errorf("AccelerationVector.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHDRImageTypeString(t *testing.T) {
	tests := []struct {
		name     string
		hdrType  HDRImageType
		expected string
	}{
		{
			name:     "HDR Processed",
			hdrType:  HDRProcessed,
			expected: "HDR Image",
		},
		{
			name:     "HDR Original",
			hdrType:  HDROriginal,
			expected: "Original Image",
		},
		{
			name:     "HDR Unknown",
			hdrType:  HDRUnknown,
			expected: "Unknown HDRImage Type 0",
		},
		{
			name:     "Invalid value",
			hdrType:  HDRImageType(99),
			expected: "Unknown HDRImage Type 99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hdrType.String()
			if got != tt.expected {
				t.Errorf("HDRImageType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAFStateFromRAW(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected AFState
	}{
		{
			name:     "Stable (1)",
			input:    1,
			expected: true,
		},
		{
			name:     "Unstable (0)",
			input:    0,
			expected: false,
		},
		{
			name:     "Negative value",
			input:    -1,
			expected: false,
		},
		{
			name:     "Large positive value",
			input:    100,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AFStateFromRAW(tt.input)
			if got != tt.expected {
				t.Errorf("AFStateFromRAW(%d) = %v, want %v", tt.input, got, tt.expected)
			}

			// Also verify String() output
			expectedStr := "No"
			if tt.expected {
				expectedStr = "Yes"
			}
			if got.String() != expectedStr {
				t.Errorf("AFStateFromRAW(%d).String() = %v, want %v", tt.input, got.String(), expectedStr)
			}
		})
	}
}
