// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"errors"
	"testing"
)

func TestQuantize_WrongSize(t *testing.T) {
	_, err := Quantize(make([]float32, 100), 0)
	if !errors.Is(err, ErrQuantizeLength) {
		t.Fatalf("expected ErrQuantizeLength, got %v", err)
	}
}

func TestQuantize_AllAboveMedian(t *testing.T) {
	input := make([]float32, 256)
	for i := range input {
		input[i] = 2.0
	}
	hash, err := Quantize(input, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	for i, word := range hash {
		if word != 0xFFFF {
			t.Errorf("hash[%d] = %04x, want ffff", i, word)
		}
	}
}

func TestQuantize_AllBelowMedian(t *testing.T) {
	input := make([]float32, 256)
	hash, err := Quantize(input, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	for i, word := range hash {
		if word != 0 {
			t.Errorf("hash[%d] = %04x, want 0000", i, word)
		}
	}
}

func TestQuantize_EqualToMedianNotSet(t *testing.T) {
	input := make([]float32, 256)
	for i := range input {
		input[i] = 1.0
	}
	hash, err := Quantize(input, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	for i, word := range hash {
		if word != 0 {
			t.Errorf("hash[%d] = %04x, want 0000 (equal-to-median must not be set)", i, word)
		}
	}
}

func TestQuantize_BitPositions(t *testing.T) {
	input := make([]float32, 256)
	input[0] = 2.0
	hash, err := Quantize(input, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	if hash[0] != 0x0001 {
		t.Errorf("hash[0] = %04x, want 0001", hash[0])
	}
	for i := 1; i < 16; i++ {
		if hash[i] != 0 {
			t.Errorf("hash[%d] = %04x, want 0000", i, hash[i])
		}
	}
}
