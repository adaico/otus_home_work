package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	idleState = iota
	escapeState
	runeState
)

func Unpack(input string) (string, error) {
	state := idleState

	var sb strings.Builder
	var currentRune rune

	for _, rune := range input {
		isSlash := rune == '\\'
		isDigit := unicode.IsDigit(rune)
		isLetter := !isSlash && !isDigit

		switch state {
		case idleState:
			switch {
			case isSlash:
				state = escapeState
			case isDigit:
				return "", ErrInvalidString
			default:
				currentRune = rune
				state = runeState
			}
		case escapeState:
			if isLetter {
				return "", ErrInvalidString
			}

			currentRune = rune
			state = runeState
		case runeState:
			switch {
			case isDigit:
				write(&sb, currentRune, int(rune-'0'))
				state = idleState
			case isSlash:
				write(&sb, currentRune, 1)
				state = escapeState
			default:
				write(&sb, currentRune, 1)
				currentRune = rune
			}
		}
	}

	switch state {
	case escapeState:
		return "", ErrInvalidString
	case runeState:
		write(&sb, currentRune, 1)
	}

	return sb.String(), nil
}

func write(sb *strings.Builder, symbol rune, times int) {
	if times == 1 {
		sb.WriteRune(symbol)
	} else {
		repeatedRune := strings.Repeat(string(symbol), times)

		sb.WriteString(repeatedRune)
	}
}
