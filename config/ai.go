package config

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"context"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	ExternalWound int8 = 0
	Simplifier    int8 = 1
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

func InitModel(client *genai.Client, cnf *Config, modelType int8) (*genai.GenerativeModel, error) {
	model := cnf.Env.GetString("MODEL")
	generativeModel := client.GenerativeModel(model)

	var systemInstruction = ""
	if modelType == ExternalWound {
		systemInstruction = cnf.Env.GetString("EVIA_SYSTEM_INSTRUCTION")
		externalWoundConfig(generativeModel)
	} else if modelType == Simplifier {
		systemInstruction = cnf.Env.GetString("SIMPLIFIER_SYSTEM_INSTRUCTION")
		simplifierConfig(generativeModel)
	} else {
		return nil, exceptions.NewInternalServerError()
	}

	generativeModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	return generativeModel, nil
}

func simplifierConfig(generativeModel *genai.GenerativeModel) {
	generativeModel.SetTemperature(1)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "text/plain"
}

func externalWoundConfig(generativeModel *genai.GenerativeModel) {
	generativeModel.SetTemperature(1.6)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "application/json"
	generativeModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"overview": {
				Type: genai.TypeString,
			},
			"conclusion": {
				Type: genai.TypeString,
			},
			"details": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"symptoms": {
						Type: genai.TypeString,
					},
					"handling": {
						Type: genai.TypeString,
					},
					"drug": {
						Type: genai.TypeString,
					},
					"reason": {
						Type: genai.TypeString,
					},
					"precautions": {
						Type: genai.TypeString,
					},
				},
			},
		},
	}
}
