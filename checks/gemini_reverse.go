package checks

import "github.com/posilva/translatin/translator"

type GeminiReverseChecker struct {
	translator *translator.GeminiTranslator
}

func NewGeminiReverseChecker(t *translator.GeminiTranslator) *GeminiReverseChecker {
	return &GeminiReverseChecker{
		translator: t,
	}
}

func (c *GeminiReverseChecker) Check(queryResult translator.QueryResult) (CheckResult, error) {
	checkQry := translator.Query{
		SourceLanguage: queryResult.Query.TargetLanguage,
		TargetLanguage: queryResult.Query.SourceLanguage,
		Text:           queryResult.Outcome,
	}

	translation, err := c.translator.Translate(checkQry)
	if err != nil {
		return CheckResult{}, err
	}

	return CheckResult{
		Outcome: translation.Outcome,
	}, nil
}
