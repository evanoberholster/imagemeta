// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

const n = 256
const half = (n + 1) / 2

// TorbenMedian calculates the median of a float32 slice using Torben's algorithm, optimized for a fixed slice size of 256.
// Returns the median value or an error if the input slice length is not 256.
func TorbenMedian(m []float32) (float32, error) {
	if len(m) != n {
		return 0, ErrTorbenElementLength
	}

	lo, hi := m[0], m[0]

	for _, v := range m[1:] {
		if v < lo {
			lo = v
		}

		if v > hi {
			hi = v
		}
	}

	for {
		guess := (lo + hi) / 2

		var less, greater, equal int
		maxLTGuess := lo
		minGTGuess := hi

		for _, v := range m {
			if v < guess {
				less++

				if v > maxLTGuess {
					maxLTGuess = v
				}
			} else if v > guess {
				greater++

				if v < minGTGuess {
					minGTGuess = v
				}
			} else {
				equal++
			}

		}

		if less <= half && greater <= half {
			if less >= half {
				return maxLTGuess, nil
			}

			if less+equal >= half {
				return guess, nil
			}

			return minGTGuess, nil
		}

		if less > greater {
			hi = maxLTGuess
		} else {
			lo = minGTGuess
		}
	}
}
