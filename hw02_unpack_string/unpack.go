package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var sb strings.Builder

	runeInput := []rune(input)
	escaped := false

	for i, rune := range runeInput {
		isDigitNext := (len(runeInput) > i+1) && unicode.IsDigit(runeInput[i+1])
		isDigit := unicode.IsDigit(rune)
		isSlash := rune == 92 // '\' ~ 92

		if isDigit && !escaped {
			if i == 0 || isDigitNext {
				return "", ErrInvalidString
			}
			continue
		}

		if isSlash && !escaped {
			escaped = true
			continue
		}

		if escaped && !isDigit && !isSlash {
			return "", ErrInvalidString
		}

		if isDigit || isSlash || !escaped {
			if isDigitNext {
				repeatTimes, _ := strconv.Atoi(string(runeInput[i+1]))
				repeatedRune := strings.Repeat(string(rune), repeatTimes)

				sb.WriteString(repeatedRune)
			} else {
				sb.WriteRune(rune)
			}
		}

		escaped = false
	}

	if escaped {
		return "", ErrInvalidString
	}

	return sb.String(), nil
}
