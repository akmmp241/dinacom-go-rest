package exceptions

import "net/http"

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

func NewConflictError(msg string) HttpConflictError {
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
