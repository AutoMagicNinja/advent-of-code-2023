package main

import (
	"fmt"
	"strconv"
	"unicode"
)

func getFirstLastUnicodeDigitIndices(line []rune) (int, int) {
	var (
		first, last int = -1, -1
	)

	for idx, chr := range line {
		logger.Debug(fmt.Sprintf("line: %s; idx=%d; chr=%c (%t); first=%d; last=%d;", string(line), idx, chr, unicode.IsDigit(chr), first, last))
		if unicode.IsDigit(chr) {
			if first < 0 {
				first = idx
			}
			if idx > last {
				last = idx
			}
		}
	}
	return first, last
}

func resolveNumberByDigits(line []rune) (result int, err error) {
	var (
		first, last int = getFirstLastUnicodeDigitIndices(line)
	)
	return strconv.Atoi(string([]rune{line[first], line[last]}))
}

