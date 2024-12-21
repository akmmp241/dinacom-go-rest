package config

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIClient struct {
	Genai *genai.Client
}

func InitAiClient(cnf *Config) *AIClient {
	ctx := context.Background()

	geminiApiKey := cnf.Env.GetString("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		panic(err)
	}

	return &AIClient{
		Genai: client,
	}
}
