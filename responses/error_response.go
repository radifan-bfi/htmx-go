package responses

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func renderError(code int, err error, message string) error {
	return echo.NewHTTPError(code, &ErrResponse{
		Err:            err,
		HTTPStatusCode: code,
		StatusText:     message,
		ErrorText:      err.Error(),
	})
}

func ErrInvalidRequest(err error) error {
	return renderError(http.StatusBadRequest, err, "Invalid request.")
}
