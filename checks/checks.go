package checks

import "github.com/posilva/translatin/translator"

type CheckResult struct {
	Outcome string
}

type Checker interface {
	Check(query translator.QueryResult) (CheckResult, error)
}
