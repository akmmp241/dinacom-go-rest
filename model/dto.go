package model

import (
	"mime/multipart"
)

type GlobalResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors"`
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=255"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type MeResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SimplifyRequest struct {
	Message string `json:"message" validate:"required"`
}

type SimplifyResponse struct {
	Message       string `json:"message"`
	SimplifiedMsg string `json:"simplified_msg"`
}

type ComplaintRequest struct {
	Complaint string                `json:"complaint" validate:"required"`
	Image     *multipart.FileHeader `json:"image" validate:"required"`
}

type ExternalWoundDetails struct {
	Symptoms    string `json:"symptoms"`
	Handling    string `json:"handling"`
	Drug        string `json:"drug"`
	Reason      string `json:"reason"`
	Precautions string `json:"precautions"`
}

type GeminiComplaintResponse struct {
	Overview       string               `json:"overview"`
	Conclusion     string               `json:"conclusion"`
	SuggestedTitle string               `json:"suggested_title"`
	Details        ExternalWoundDetails `json:"details"`
}

type ComplaintResponse struct {
	ComplaintId string                  `json:"complaint_id"`
	Response    GeminiComplaintResponse `json:"response"`
	ImageUrl    string                  `json:"image_url"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyForgetPasswordOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type VerifyForgetPasswordOtpResponse struct {
	Email              string `json:"email"`
	ResetPasswordToken string `json:"reset_password_token"`
}

type ResetPasswordRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}
