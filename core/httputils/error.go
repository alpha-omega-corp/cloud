package httputils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	errEOF      = BadRequest("eof", "EOF reading HTTP request body")
	ErrNotFound = NotFound("not found")
	ErrInternal = NewError(http.StatusInternalServerError, "internal", "internal server error")
)

func NotFound(msg string, args ...interface{}) Error {
	return NewError(http.StatusNotFound, "not_found", msg, args...)
}

func Forbidden(msg string, args ...interface{}) Error {
	return NewError(http.StatusForbidden, "forbidden", msg, args...)
}

func BadRequest(code, msg string, args ...interface{}) Error {
	return NewError(http.StatusBadRequest, code, msg, args...)
}

//------------------------------------------------------------------------------

type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(status int, code, msg string, args ...interface{}) Error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return Error{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}

func (e Error) Error() string {
	return e.Message
}

//------------------------------------------------------------------------------

func From(err error, debug bool) Error {
	switch err {
	case io.EOF:
		return errEOF
	case sql.ErrNoRows:
		return ErrNotFound
	}

	switch err := err.(type) {
	case Error:
		return err
	case *json.SyntaxError:
		return BadRequest("json_syntax", err.Error())
	}

	if debug {
		return NewError(http.StatusInternalServerError, "internal", err.Error())
	}
	return ErrInternal
}
