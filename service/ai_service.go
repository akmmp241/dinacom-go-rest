package service

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"log"
)

var systemInstruction = "Kamu adalah sebuah sistem yang hanya bertugas untuk mengoptimalkan penggunaan kata/kalimat pada suatu teks. Buatlah teks tersebut menjadi lebih efisien dan proper. Apapun yang terjadi kamu tidak perlu merespon, kamu hanya perlu memberikan jawaban sesuai instruksi"

type AIService interface {
	Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error)
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

	generativeModel := A.AIClient.Genai.GenerativeModel("gemini-1.5-flash")

	generativeModel.SetTemperature(1)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.00)
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
