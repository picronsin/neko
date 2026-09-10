package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/m1k1o/neko/server/pkg/protocol"
)

func HttpJsonRequest(w http.ResponseWriter, r *http.Request, res any) error {
	err := json.NewDecoder(r.Body).Decode(res)

	if err == nil {
		return nil
	}

	if err == io.EOF {
		return HttpBadRequest("no data provided").WithInternalErr(err)
	}

	return HttpBadRequest("unable to parse provided data").WithInternalErr(err)
}

func HttpJsonResponse(w http.ResponseWriter, code int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Err(err).Str("module", "http").Msg("sending http json response failed")
	}
}

func HttpSuccess(w http.ResponseWriter, res ...any) error {
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		HttpJsonResponse(w, http.StatusOK, res[0])
	}

	return nil
}

// HTTPError is an error with a message and an HTTP status code.
type HTTPError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`

	InternalErr error  `json:"-"`
	InternalMsg string `json:"-"`
}

// WithErrorCode attaches the stable protocol/application code used by typed
// clients. HTTP status remains available for generic HTTP tooling.
func (e *HTTPError) WithErrorCode(code string) *HTTPError {
	e.ErrorCode = code
	return e
}

func (e *HTTPError) Error() string {
	if e.InternalMsg != "" {
		return e.InternalMsg
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *HTTPError) Cause() error {
	if e.InternalErr != nil {
		return e.InternalErr
	}
	return e
}

// WithInternalErr adds internal error information to the error
func (e *HTTPError) WithInternalErr(err error) *HTTPError {
	e.InternalErr = err
	return e
}

// WithInternalMsg adds internal message information to the error
func (e *HTTPError) WithInternalMsg(msg string) *HTTPError {
	e.InternalMsg = msg
	return e
}

// WithInternalMsg adds internal formated message information to the error
func (e *HTTPError) WithInternalMsgf(fmtStr string, args ...any) *HTTPError {
	e.InternalMsg = fmt.Sprintf(fmtStr, args...)
	return e
}

// Sends error with custom formated message
func (e *HTTPError) Msgf(fmtSt string, args ...any) *HTTPError {
	e.Message = fmt.Sprintf(fmtSt, args...)
	return e
}

// Sends error with custom message
func (e *HTTPError) Msg(str string) *HTTPError {
	e.Message = str
	return e
}

func HttpError(code int, res ...string) *HTTPError {
	err := &HTTPError{
		Code:      code,
		Message:   http.StatusText(code),
		ErrorCode: defaultHTTPErrorCode(code),
	}

	if len(res) == 1 {
		err.Message = res[0]
	}

	return err
}

func defaultHTTPErrorCode(status int) string {
	switch status {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return string(protocol.InvalidPayload)
	case http.StatusUnauthorized, http.StatusForbidden:
		return string(protocol.PermissionDenied)
	case http.StatusNotFound:
		return string(protocol.NotFound)
	case http.StatusConflict:
		return string(protocol.ControlConflict)
	default:
		return string(protocol.InternalError)
	}
}

func HttpBadRequest(res ...string) *HTTPError {
	return HttpError(http.StatusBadRequest, res...)
}

func HttpUnauthorized(res ...string) *HTTPError {
	return HttpError(http.StatusUnauthorized, res...)
}

func HttpForbidden(res ...string) *HTTPError {
	return HttpError(http.StatusForbidden, res...)
}

func HttpNotFound(res ...string) *HTTPError {
	return HttpError(http.StatusNotFound, res...)
}

func HttpUnprocessableEntity(res ...string) *HTTPError {
	return HttpError(http.StatusUnprocessableEntity, res...)
}

func HttpInternalServerError(res ...string) *HTTPError {
	return HttpError(http.StatusInternalServerError, res...)
}
