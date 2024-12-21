package service

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/helpers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"log"
)

type AIService interface {
	Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error)
	ExternalWound(ctx context.Context, req *model.ExternalWoundRequest) (*model.ExternalWoundResponse, error)
}

type AIServiceImpl struct {
	Validate *validator.Validate
	Cnf      *config.Config
	AIClient *config.AIClient
}

func NewAIService(
	validate *validator.Validate,
	cnf *config.Config,
	aiClient *config.AIClient,
) AIService {
	return &AIServiceImpl{
		Validate: validate,
		Cnf:      cnf,
		AIClient: aiClient,
	}
}

func (A AIServiceImpl) Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error) {
	err := A.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	systemInstruction := A.Cnf.Env.GetString("SIMPLIFIER_SYSTEM_INSTRUCTION")

	generativeModel := A.AIClient.Genai.GenerativeModel("gemini-1.5-flash")

	generativeModel.SetTemperature(1)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "text/plain"
	generativeModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	session := generativeModel.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(req.Message))
	if err != nil {
		log.Println("Error while sending message: ", err.Error())
		return nil, exceptions.NewInternalServerError()
	}

	simplifiedMsg := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		simplifiedMsg += fmt.Sprintf("%v\n", part)
	}

	return &model.SimplifyResponse{
		Message:       req.Message,
		SimplifiedMsg: simplifiedMsg,
	}, nil
}

func (A AIServiceImpl) ExternalWound(ctx context.Context, req *model.ExternalWoundRequest) (*model.ExternalWoundResponse, error) {
	err := A.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	systemInstruction := A.Cnf.Env.GetString("EVIA_SYSTEM_INSTRUCTION")

	generativeModel := A.AIClient.Genai.GenerativeModel("gemini-1.5-flash")

	generativeModel.SetTemperature(1.6)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "application/json"
	generativeModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}
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

	fileURIs, err := helpers.UploadToGemini(ctx, A.AIClient.Genai, req.Image, "image/png")

	session := generativeModel.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.FileData{URI: fileURIs},
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("test"))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	var externalWoundResponse model.ExternalWoundResponse
	for _, part := range resp.Candidates[0].Content.Parts {
		jsonPart := fmt.Sprintf("%v\n", part)
		err = json.Unmarshal([]byte(jsonPart), &externalWoundResponse)
		if err != nil {
			log.Println("Error while unmarshalling response: ", err.Error())
			return nil, exceptions.NewInternalServerError()
		}
	}

	return &externalWoundResponse, nil
}
