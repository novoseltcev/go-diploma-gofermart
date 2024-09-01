package utils

const (
	asciiZero = 48
	asciiTen  = 57
)

func ValidateLunhNumber(number string) bool {
	// Copy-paste: https://github.com/ShiraazMoollatjie/goluhn/blob/master/goluhn.go#L18

	var sum int64
	parity := len(number) % 2

	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return false
		}

		d = d - asciiZero
		// Double the value of every second digit.
		if i % 2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}

		// Take the sum of all the digits.
		sum += int64(d)
	}

	return sum % 10 == 0
}
