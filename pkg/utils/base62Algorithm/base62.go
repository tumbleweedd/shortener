package base62Algorithm

import (
	"errors"
	"math"
	"strings"
)

const (
	alphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length     = uint64(len(alphabet))
	codeLength = 8
)

func Encode(num uint64) string {
	var encoded strings.Builder
	encoded.Grow(codeLength)

	for i := 0; i < codeLength; i++ {
		encoded.WriteByte(alphabet[num%length])
		num = num / length
	}

	return encoded.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, symbol := range encoded {
		alphabeticPosition := strings.Index(alphabet, string(symbol))

		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
