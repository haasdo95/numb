package utils

import (
	"unicode"
)

func Concrete2Abstract(concrete string) ([]rune, string) {
	abstract := ""
	escape := false
	smeared := make([]rune, 0)
	for _, curr := range concrete {
		var toAppend rune
		if curr == rune('%') || curr == rune('i') {
			escape = true
			toAppend = curr
		} else if unicode.IsDigit(curr) {
			if !escape {
				toAppend = rune('#') // smear off the number!
				smeared = append(smeared, curr)
			} else {
				toAppend = curr
			}
		} else {
			escape = false
			toAppend = curr
		}
		abstract += string(toAppend)
	}
	return smeared, abstract
}
