// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"errors"
	"math/rand"
	"sort"
	"testing"
)

func TestTorbenMedian_WrongSize(t *testing.T) {
	_, err := TorbenMedian(make([]float32, 100))

	if !errors.Is(err, ErrTorbenElementLength) {
		t.Fatalf("expected error for TorbenMedian with wrong size, got %v", err)
	}
}

func TestTorbenMedian_AllSameValue(t *testing.T) {
	m := make([]float32, 256)
	for i := range m {
		m[i] = 42.0
	}
	got, err := TorbenMedian(m)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42.0 {
		t.Errorf("got %v, want 42.0", got)
	}
}

func TestTorbenMedian_SortedSequence(t *testing.T) {
	m := make([]float32, 256)
	for i := range m {
		m[i] = float32(i)
	}

	got, err := TorbenMedian(m)
	if err != nil {
		t.Fatal(err)
	}

	if got != 127.0 {
		t.Errorf("got %v, want 127.0", got)
	}
}

func TestTorbenMedian_MatchesSortMedian(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	m := make([]float32, 256)
	for i := range m {
		m[i] = rng.Float32() * 1000
	}

	got, err := TorbenMedian(m)
	if err != nil {
		t.Fatal(err)
	}

	sorted := make([]float32, 256)
	copy(sorted, m)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	want := sorted[127]
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
