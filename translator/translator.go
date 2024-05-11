package translator

type Query struct {
	ID             string
	Key            string
	SourceLanguage string
	TargetLanguage string
	Text           string
}

type QueryResult struct {
	Query   Query
	Outcome string
}

type Translator interface {
	Translate(query Query) (QueryResult, error)
}
