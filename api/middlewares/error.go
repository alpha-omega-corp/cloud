package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bunrouter"
	"io"
	"net/http"
	"net/http/httptest"
)

type HTTPError struct {
	statusCode int

	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(err error) HTTPError {
	switch {
	case err == io.EOF:
		return HTTPError{
			statusCode: http.StatusBadRequest,

			Code:    "eof",
			Message: "EOF reading HTTP request body",
		}
	case errors.Is(err, sql.ErrNoRows):
		return HTTPError{
			statusCode: http.StatusNotFound,

			Code:    "not_found",
			Message: "Page Not Found",
		}
	}

	return HTTPError{
		statusCode: http.StatusInternalServerError,

		Code:    "internal",
		Message: "Internal server error",
	}
}

func NewErrorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := next(w, req)
		rec := httptest.NewRecorder()
		fmt.Print(err)

		switch err := err.(type) {
		case nil:
		case HTTPError:
			w.WriteHeader(err.statusCode)
			_ = bunrouter.JSON(w, err)
		default:
			httpErr := NewHTTPError(err)
			w.WriteHeader(httpErr.statusCode)
			_ = bunrouter.JSON(w, httpErr)
		}

		fmt.Printf("%s %s returned %d\n", req.Method, req.Route(), rec.Code)

		return nil
	}
}
