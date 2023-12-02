package main

import (
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
)

var (
	inputFilename string
	logger        *slog.Logger
)

func init() {
	flag.StringVar(&inputFilename, "input", "", "Input filename")
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func firstStar(lines []string) {
	var (
		err        error
		answer     int
		lineResult int
	)
	for lineNumber, eachLine := range lines {
		lineResult, err = resolveNumberByDigits([]rune(eachLine))
		if err != nil {
			logger.Debug(fmt.Sprintf("lineno: %d, line: %#v, result: %d; err: %s", lineNumber, eachLine, lineResult, err.Error()))
		}
		answer += lineResult
	}
	logger.Info(fmt.Sprintf("first star: %d", answer))
}

func secondStar(lines []string) {
	var (
		answer     int
		err        error
		lineResult int
	)

	for lineNumber, eachLine := range lines {
		lineResult, err = resolveNumberByMixedDigitsAndWords(eachLine)
		if err != nil {
			logger.Error(fmt.Sprintf("lineno: %d, line: %#v, result: %d; err: %s", lineNumber, eachLine, lineResult, err.Error()))
		}
		answer += lineResult
	}

	logger.Warn(fmt.Sprintf("second star: %d", answer))
}

func main() {
	var (
		err     error
		inputFs afero.Fs = afero.NewReadOnlyFs(afero.NewOsFs())
		rawData []byte
		lines   []string
	)

	flag.Parse()
	logger.Debug(fmt.Sprintf("input filename: %#v\n", inputFilename))

	rawData, err = afero.ReadFile(inputFs, inputFilename)
	if err != nil {
		panic(err)
	}

	lines = strings.Split(strings.TrimSpace(string(rawData)), "\n")
	// firstStar(lines)
	secondStar(lines)
}
