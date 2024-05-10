package main

import (
	"fmt"

	"github.com/posilva/translatin/sources"
	"github.com/posilva/translatin/translator"
)

const (
	mainsheet = "events"
	filename  = "/Users/pedro.silva/Downloads/base_events.xlsx"
)

func main() {
	trans, err := translator.NewDefaultGeminiTranslator()
	if err != nil {
		panic(err)
	}

	source := sources.NewXLSSourceWithDefaultLanguage(sources.XLSSourceSettings{
		FilePath:  filename,
		SheetName: mainsheet,
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
	translation, err := trans.Translate(first)
	if err != nil {
		panic(err)
	}

	checkQry := translator.Query{
		SourceLanguage: first.TargetLanguage,
		TargetLanguage: first.SourceLanguage,
		Text:           translation.Outcome,
	}

	translationCheck, err := trans.Translate(checkQry)
	if err != nil {
		panic(err)
	}
	fmt.Printf("------ORIGINAL------\n\n")
	fmt.Println(first.Text)
	fmt.Printf("------GENINI-----\n\n")
	fmt.Println(translation.Outcome)
	fmt.Printf("------GEMINI REVERSE CHECK ------\n\n")
	fmt.Println(translationCheck.Outcome)
}
