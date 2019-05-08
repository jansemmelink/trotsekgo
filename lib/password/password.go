package password

import (
	"unicode"

	"github.com/jansemmelink/log"
)

//IsStrong checks password strenght
func IsStrong(pw string) error {
	if len(pw) < 8 {
		return log.Wrapf(nil, "password has less than 8 characters")
	}

	nrOfDigits := 0
	nrOfSymbols := 0
	nrOfLower := 0
	nrOfUpper := 0
	for _, c := range pw {
		if unicode.IsLetter(c) {
			if unicode.IsLower(c) {
				nrOfLower++
			}
			if unicode.IsUpper(c) {
				nrOfUpper++
			}
		}
		if unicode.IsDigit(c) {
			nrOfDigits++
		}
		if unicode.IsGraphic(c) {
			nrOfSymbols++
		}
	}

	if nrOfLower <= 0 || nrOfUpper <= 0 {
		return log.Wrapf(nil, "password does not have lowercase and uppercase letters")
	}
	if nrOfDigits <= 0 {
		return log.Wrapf(nil, "password does not contain any numeric digits")
	}
	if nrOfSymbols <= 0 {
		return log.Wrapf(nil, "password does not contain any symbols")
	}
	//string enough
	return nil
}
