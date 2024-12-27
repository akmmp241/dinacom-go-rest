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
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type ComplaintService interface {
	Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error)
	ExternalWound(ctx context.Context, req *model.ComplaintRequest, user *model.User) (*model.ComplaintResponse, error)
	GetById(ctx context.Context, complaintId string, user *model.User) (*model.ComplaintResponse, error)
	GetAll(ctx context.Context, user *model.User) (*[]model.ComplaintResponse, error)
}

type ComplaintServiceImpl struct {
	Validate      *validator.Validate
	Cnf           *config.Config
	AIClient      *config.AIClient
	AWSClient     *config.AWSClient
	DB            *sql.DB
	ComplaintRepo repository.ComplaintRepository
}

func NewComplaintService(
	validate *validator.Validate,
	cnf *config.Config,
	aiClient *config.AIClient,
	awsClient *config.AWSClient,
	complaintRepo repository.ComplaintRepository,
	db *sql.DB,
) ComplaintService {
	return &ComplaintServiceImpl{
		Validate:      validate,
		Cnf:           cnf,
		AIClient:      aiClient,
		AWSClient:     awsClient,
		ComplaintRepo: complaintRepo,
		DB:            db,
	}
}

func (A ComplaintServiceImpl) Simplifier(ctx context.Context, req *model.SimplifyRequest) (*model.SimplifyResponse, error) {
	err := A.Validate.Struct(req)
	if err != nil {
		return nil, exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
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

func (A ComplaintServiceImpl) ExternalWound(ctx context.Context, req *model.ComplaintRequest, user *model.User) (*model.ComplaintResponse, error) {
	err := A.Validate.Struct(req)
	if err != nil {
		return nil, exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
	}

	generativeModel, err := config.InitModel(A.AIClient.Genai, A.Cnf, config.ExternalWound)
	if err != nil {
		return nil, err
	}

	// upload to gemini and s3 concurrently
	fileURIs, location, err := uploadFilesConcurrently(ctx, req, A)
	if err != nil {
		return nil, err
	}

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

	var geminiComplaintResponse model.GeminiComplaintResponse
	jsonResp := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		jsonResp += fmt.Sprintf("%v\n", part)
	}
	err = json.Unmarshal([]byte(jsonResp), &geminiComplaintResponse)
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
		Title:         geminiComplaintResponse.SuggestedTitle,
		ComplaintsMsg: req.Complaint,
		Response:      jsonResp,
		ImageUrl:      location,
		CreatedAt:     time.Now(),
	}

	_, err = A.ComplaintRepo.Save(ctx, tx, &complaint)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()

	externalWoundResponse := model.ComplaintResponse{
		ComplaintId: generatedId,
		Response:    geminiComplaintResponse,
		ImageUrl:    location,
	}

	return &externalWoundResponse, nil
}

func (A ComplaintServiceImpl) GetById(ctx context.Context, complaintId string, user *model.User) (*model.ComplaintResponse, error) {
	tx, err := A.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	complaint, err := A.ComplaintRepo.FindById(ctx, tx, complaintId)
	if err != nil && errors.Is(err, exceptions.NotFoundError{}) {
		return nil, exceptions.NewHttpNotFoundError("Complaint not found")
	} else if err != nil && !errors.Is(err, exceptions.NotFoundError{}) {
		return nil, err
	}

	if complaint.UserId != user.Id {
		return nil, exceptions.NewForbiddenError("You are not authorized to access this complaint")
	}

	var geminiComplaintResponse model.GeminiComplaintResponse
	err = json.Unmarshal([]byte(complaint.Response), &geminiComplaintResponse)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	complaintResponse := model.ComplaintResponse{
		ComplaintId: complaint.Id,
		Response:    geminiComplaintResponse,
		ImageUrl:    complaint.ImageUrl,
	}

	return &complaintResponse, nil
}

func (A ComplaintServiceImpl) GetAll(ctx context.Context, user *model.User) (*[]model.ComplaintResponse, error) {
	tx, err := A.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	complaints, err := A.ComplaintRepo.FindAll(ctx, tx, user.Id)
	if err != nil {
		return nil, err
	}

	var complaintResponses []model.ComplaintResponse
	for _, complaint := range complaints {
		var geminiComplaintResponse model.GeminiComplaintResponse
		err = json.Unmarshal([]byte(complaint.Response), &geminiComplaintResponse)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		complaintResponse := model.ComplaintResponse{
			ComplaintId: complaint.Id,
			Response:    geminiComplaintResponse,
			ImageUrl:    complaint.ImageUrl,
		}
		complaintResponses = append(complaintResponses, complaintResponse)
	}

	return &complaintResponses, nil
}

func uploadFilesConcurrently(ctx context.Context, req *model.ComplaintRequest, A ComplaintServiceImpl) (fileURIs string, location string, err error) {
	var wg sync.WaitGroup

	fileURIsCh := make(chan string, 1)
	locationCh := make(chan string, 1)
	errorCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		open, err := req.Image.Open()
		if err != nil {
			errorCh <- err
			return
		}
		defer open.Close()

		uri, err := helpers.UploadToGemini(ctx, A.AIClient.Genai, open, "image/png")
		if err != nil {
			errorCh <- err
			return
		}
		fileURIsCh <- uri
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		open, err := req.Image.Open()
		if err != nil {
			errorCh <- err
			return
		}
		defer open.Close()

		loc, err := helpers.UploadS3(ctx, A.AWSClient.Uploader, open, req.Image.Filename, A.Cnf.Env.GetString("AWS_BUCKET_NAME"))
		if err != nil {
			errorCh <- err
			return
		}
		locationCh <- loc
	}()

	wg.Wait()

	close(fileURIsCh)
	close(locationCh)
	close(errorCh)

	if len(errorCh) > 0 {
		err := <-errorCh
		log.Println("Error while uploading:", err)
		return "", "", exceptions.NewInternalServerError()
	}

	fileURI := <-fileURIsCh
	location = <-locationCh

	return fileURI, location, nil
}
