package expand

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func Init() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Println("expand: OPENAI_API_KEY not set, query expansion disabled")
		return
	}
	client = openai.NewClient(key)
	log.Println("expand: query expansion enabled")
}

// Query expands a search query with synonyms and related tech terms using GPT-4o-mini.
// Falls back to original query if OpenAI is unavailable.
func Query(ctx context.Context, query string) string {
	if client == nil {
		return query
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	prompt := `Expand this job search query with common synonyms and related tech terms.
Return ONLY a short comma-separated list of terms, no explanation, no punctuation, lowercase.
Include the original term. Max 6 terms total.

Query: ` + query

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens:   60,
		Temperature: 0,
	})
	if err != nil {
		log.Printf("expand: openai error: %v, using original query", err)
		return query
	}

	expanded := strings.TrimSpace(resp.Choices[0].Message.Content)
	log.Printf("expand: %q → %q", query, expanded)
	return expanded
}
