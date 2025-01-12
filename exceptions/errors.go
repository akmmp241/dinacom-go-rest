package exceptions

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

type GlobalError interface {
	Error() string
	GetCode() int
}

type HttpBadRequestError struct {
	Msg  string
	Code int
}

func NewBadRequestError(msg string) HttpBadRequestError {
	return HttpBadRequestError{Msg: msg, Code: http.StatusBadRequest}
}

func (v HttpBadRequestError) Error() string {
	return v.Msg
}

func (v HttpBadRequestError) GetCode() int {
	return v.Code
}

type NotFoundError struct {
	Msg string
}

func NewNotFoundError() NotFoundError {
	return NotFoundError{}
}

func (u NotFoundError) Error() string {
	return u.Msg
}

type HttpConflictError struct {
	Msg  string
	Code int
}

func (h HttpConflictError) Error() string {
	return h.Msg
}

func (h HttpConflictError) GetCode() int {
	return h.Code
}

func NewHttpConflictError(msg string) HttpConflictError {
	return HttpConflictError{Msg: msg, Code: http.StatusConflict}
}

type HttpInternalServerError struct {
	Msg  string
	Code int
}

func (i HttpInternalServerError) Error() string {
	return i.Msg
}

func (i HttpInternalServerError) GetCode() int {
	return i.Code
}

func NewInternalServerError() HttpInternalServerError {
	return HttpInternalServerError{Msg: "Internal Server Error", Code: http.StatusInternalServerError}
}

type HttpUnauthorized struct {
	Msg  string
	Code int
}

func (u HttpUnauthorized) Error() string {
	return u.Msg
}

func (u HttpUnauthorized) GetCode() int {
	return u.Code
}

func NewUnauthorizedError(msg string) HttpUnauthorized {
	return HttpUnauthorized{Msg: msg, Code: http.StatusUnauthorized}
}

type HttpForbiddenError struct {
	Msg  string
	Code int
}

func (f HttpForbiddenError) Error() string {
	return f.Msg
}

func (f HttpForbiddenError) GetCode() int {
	return f.Code
}

func NewForbiddenError(msg string) HttpForbiddenError {
	return HttpForbiddenError{Msg: msg, Code: http.StatusForbidden}
}

type HttpNotFoundError struct {
	Msg  string
	Code int
}

func (f HttpNotFoundError) Error() string {
	return f.Msg
}

func (f HttpNotFoundError) GetCode() int {
	return f.Code
}

func NewHttpNotFoundError(msg string) HttpNotFoundError {
	return HttpNotFoundError{Msg: msg, Code: http.StatusNotFound}
}

type FailedValidationError struct {
	Msg    string
	Code   int
	Errors map[string]interface{}
}

func (f FailedValidationError) Error() string {
	return f.Msg
}

func (f FailedValidationError) GetCode() int {
	return f.Code
}

func NewFailedValidationError(obj interface{}, err validator.ValidationErrors) FailedValidationError {

	objRef := reflect.TypeOf(obj)

	errMsgs := make(map[string]interface{})

	for i := 0; i < objRef.NumField(); i++ {
		structField := objRef.Field(i)
		errMsgs[structField.Tag.Get("json")] = nil
	}

	for _, err := range err {
		structField, _ := objRef.FieldByName(err.Field())
		field := structField.Tag.Get("json")
		errMsgs[field] = handleValidationErrorMessage(err.Tag(), err.Param(), err.Field())
	}

	return FailedValidationError{Msg: "Failed Validation", Code: http.StatusUnprocessableEntity, Errors: errMsgs}
}
