package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/posilva/translatin/checks"
	"github.com/posilva/translatin/sources"
	"github.com/posilva/translatin/translator"
)

const (
	mainsheet = "events"
	filename  = "/Users/pedro.silva/Downloads/base_events.xlsx"
)

var (
	filepathFlag  = flag.String("filepath", filename, "path to the XLS file")
	mainSheetFlag = flag.String("sheet", mainsheet, "name of the source sheet")
)

func main() {
	flag.Parse()

	trans, err := translator.NewDefaultGeminiTranslator()
	if err != nil {
		panic(err)
	}

	checker := checks.NewGeminiReverseChecker(trans)

	source := sources.NewXLSSourceWithDefaultLanguage(sources.XLSSourceSettings{
		FilePath:  *filepathFlag,
		SheetName: *mainSheetFlag,
	})

	if err := source.Open(); err != nil {
		panic(err)
	}
	defer source.Close()

	entries, err := source.BuildQueries()
	if err != nil {
		panic(err)
	}

	first := entries[0]
	for _, e := range entries {
		if strings.Compare(e.TargetLanguage, "ko-KR") == 0 {
			first = e
			break
		}
	}
	translation, err := trans.Translate(first)
	if err != nil {
		panic(err)
	}

	checkResult, err := checker.Check(translation)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n------ORIGINAL------\n\n")
	fmt.Println(first.Text)
	fmt.Printf("\n------GENINI-----\n\n")
	fmt.Println(translation.Outcome)
	fmt.Printf("\n------GEMINI REVERSE CHECK ------\n\n")
	fmt.Println(checkResult.Outcome)
}
