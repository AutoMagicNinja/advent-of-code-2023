package main

import (
	"fmt"
	"maps"
	"strconv"
	"unicode"

	"github.com/dlclark/regexp2"
)

func getFirstLastUnicodeDigitIndices(line []rune) (int, int) {
	var (
		first, last int = -1, -1
	)

	for idx, chr := range line {
		logger.Debug(fmt.Sprintf("type=rune; line: %s; idx=%d; chr=%c (%t); first=%d; last=%d;", string(line), idx, chr, unicode.IsDigit(chr), first, last))
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

// I could probably do this using stringer, and some unstring method, but I really can't be bothered right now
func convertToDigit(s string) int {
	switch s {
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	}
	return -1
}

// returns map[starting_index]captured_string
func getAllCaptures(pattern *regexp2.Regexp, str []rune) (matches map[int]string, err error) {
	var match *regexp2.Match
	match, err = pattern.FindRunesMatch(str)
	matches = make(map[int]string)
	for err == nil && match != nil {
		for gCt, group := range match.Groups() {
			for cCt, capture := range group.Captures {
				logger.Debug(fmt.Sprintf("type=match; str=%s gCt=%d; cCt=%d; idx=%d; len=%d; s=%s", string(str), gCt, cCt, capture.Index, capture.Length, capture.String()))
				if capture.Length > 0 {
					matches[capture.Index] = capture.String()
				}
			}
		}
		match, err = pattern.FindNextMatch(match)
	}
	return
}

func getSpelledOutNumberLocations(line []rune) (results map[int]int) {
	var (
		spelledOutNumberRE = regexp2.MustCompile(`(?=(one|two|three|four|five|six|seven|eight|nine))`, regexp2.RE2)
		allMatches, err    = getAllCaptures(spelledOutNumberRE, line)
		conversionResult   int
	)
	if err != nil {
		panic(err)
	}
	logger.Debug(fmt.Sprintf("type=word; line=%s; matchCount=%d;", string(line), len(allMatches)))
	results = make(map[int]int)
	for matchIdx, eachMatch := range allMatches {
		conversionResult = convertToDigit(eachMatch)
		logger.Debug(fmt.Sprintf("type=word; line=%s; start=%d; match=%s; result=%d;",
			string(line), matchIdx, eachMatch, conversionResult))
		if conversionResult > 0 {
			results[matchIdx] = conversionResult
		}
	}

	return
}

func resolveNumberByDigits(line []rune) (result int, err error) {
	var (
		first, last int = getFirstLastUnicodeDigitIndices(line)
	)
	return strconv.Atoi(string([]rune{line[first], line[last]}))
}

func getUnicodeDigitLocations(line string) (results map[int]int) {
	results = make(map[int]int)
	var (
		runeLine    []rune = []rune(line)
		first, last        = getFirstLastUnicodeDigitIndices(runeLine)
		err         error
	)
	if first >= 0 {
		results[first], err = strconv.Atoi(string(runeLine[first]))
		if err != nil {
			panic(err)
		}
		// TODO: DRY
		results[last], err = strconv.Atoi(string(runeLine[last]))
		if err != nil {
			panic(err)
		}
	}
	return
}

func resolveNumberByMixedDigitsAndWords(line string) (result int, err error) {
	var (
		locations   map[int]int = make(map[int]int)
		first, last int         = -1, -1
	)
	unicodeLocations := getUnicodeDigitLocations(line)
	logger.Debug(fmt.Sprintf("line=%s; Runes=%#v", line, unicodeLocations))
	maps.Copy[map[int]int](locations, unicodeLocations)

	wordLocations := getSpelledOutNumberLocations([]rune(line))
	logger.Debug(fmt.Sprintf("line=%s; Words=%#v", line, wordLocations))
	maps.Copy[map[int]int](locations, wordLocations)

	logger.Debug(fmt.Sprintf("line=%s; combo=%#v", line, locations))

	for idx := range locations {
		if first < 0 || first > idx {
			first = idx
		}
		if last < idx {
			last = idx
		}
	}
	if first < 0 || last < 0 {
		err = fmt.Errorf("could not find two distinct numbers (somehow)")
	}
	result = locations[first]*10 + locations[last]
	logger.Info(fmt.Sprintf("line=%s; first=%d (%d); last=%d (%d); result=%d;", line, first, locations[first], last, locations[last], result))
	return
}
