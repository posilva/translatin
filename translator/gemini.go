package translator

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	defaultTimeout = 5 * time.Second
)

// GeminiTranslator translates text using the Gemini API.
type GeminiTranslator struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewDefaultGeminiTranslatorFromEnv creates a new GeminiTranslator from environment variables.
func NewDefaultGeminiTranslator() (*GeminiTranslator, error) {
	return NewDefaultGeminiTranslatorFromEnv()
}

// NewGeminiTranslator creates a new GeminiTranslator.
func NewGeminiTranslator(client *genai.Client) *GeminiTranslator {
	return &GeminiTranslator{
		client: client,
		model:  defaultModel(client),
	}
}

// NewDefaultGeminiTranslatorFromEnv creates a new GeminiTranslator from environment variables.
func NewDefaultGeminiTranslatorFromEnv() (*GeminiTranslator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return NewGeminiTranslator(client), nil
}

// Translate translates text using the Gemini API.
func (t *GeminiTranslator) Translate(query Query) (QueryResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	// default query
	qry := "translate from " + query.SourceLanguage + " to " + query.TargetLanguage + ": " + query.Text
	resp, err := t.model.GenerateContent(ctx, genai.Text(qry))
	if err != nil {
		return QueryResult{}, err
	}

	var buffer strings.Builder
	for _, v := range resp.Candidates {
		for _, p := range v.Content.Parts {
			buffer.WriteString(fmt.Sprintf("%s", p))
		}
	}
	return QueryResult{Query: query, Outcome: buffer.String()}, nil
}

func defaultModel(client *genai.Client) *genai.GenerativeModel {
	return client.GenerativeModel("gemini-1.0-pro")
}
