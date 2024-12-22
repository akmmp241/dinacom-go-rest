package service

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/helpers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"log"
	"time"
)

type AIService interface {
	Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error)
	ExternalWound(ctx context.Context, req *model.ExternalWoundRequest, user *model.User) (*model.ExternalWoundResponse, error)
}

type AIServiceImpl struct {
	Validate      *validator.Validate
	Cnf           *config.Config
	AIClient      *config.AIClient
	DB            *sql.DB
	ComplaintRepo repository.ComplaintRepository
}

func NewAIService(
	validate *validator.Validate,
	cnf *config.Config,
	aiClient *config.AIClient,
	complaintRepo repository.ComplaintRepository,
	db *sql.DB,
) AIService {
	return &AIServiceImpl{
		Validate:      validate,
		Cnf:           cnf,
		AIClient:      aiClient,
		ComplaintRepo: complaintRepo,
		DB:            db,
	}
}

func (A AIServiceImpl) Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error) {
	err := A.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	generativeModel, err := config.InitModel(A.AIClient.Genai, A.Cnf, config.Simplifier)
	if err != nil {
		return nil, err
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

func (A AIServiceImpl) ExternalWound(ctx context.Context, req *model.ExternalWoundRequest, user *model.User) (*model.ExternalWoundResponse, error) {
	err := A.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	generativeModel, err := config.InitModel(A.AIClient.Genai, A.Cnf, config.ExternalWound)
	if err != nil {
		return nil, err
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

	resp, err := session.SendMessage(ctx, genai.Text(req.Complaint))
	if err != nil {
		log.Println("Error while sending message: ", err.Error())
		return nil, exceptions.NewInternalServerError()
	}

	var externalWoundResponse model.ExternalWoundResponse
	jsonResp := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		jsonResp += fmt.Sprintf("%v\n", part)
	}
	err = json.Unmarshal([]byte(jsonResp), &externalWoundResponse)
	if err != nil {
		log.Println("Error while unmarshalling response: ", err.Error())
		return nil, exceptions.NewInternalServerError()
	}

	tx, err := A.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	generatedId := uuid.NewString()
	complaint := model.Complaint{
		Id:            generatedId,
		UserId:        user.Id,
		Title:         "test",
		ComplaintsMsg: req.Complaint,
		Response:      jsonResp,
		CreatedAt:     time.Now(),
	}

	_, err = A.ComplaintRepo.Save(ctx, tx, &complaint)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()

	externalWoundResponse.ComplaintId = generatedId

	return &externalWoundResponse, nil
}
