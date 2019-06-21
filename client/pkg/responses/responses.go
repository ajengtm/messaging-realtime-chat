package responses

import (
	"bitbucket.org/kudoindonesia/messaging-realtime-chat/client/pkg/constants"
	"encoding/json"
	"net/http"
)

type ErrorValidation struct {
	Errors interface{} `json:"errors"`
}

// APIResponse defines attributes for api Response
type APIResponse struct {
	HTTPCode   int         `json:"-"`
	Code       int         `json:"code"`
	Message    interface{} `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

// Write writes the data to http response writer
func Write(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(js); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Defaults API Response with standard HTTP Status Code. The default value can
// be changed either using `ModifyMessage` or `ModifyHTTPCode`. You can call it
// directly by :
//
// response.Write(res, response.APIErrorUnknown)
// return
var (
	// A generic error message, given when an unexpected condition was encountered.
	APIErrorUnknown = APIResponse{
		HTTPCode: http.StatusInternalServerError,
		Code:     constants.CodeInternalServerError,
		Message:  constants.MessageInternalServerError,
	}

	// Standard response for successful HTTP requests.
	APIOK = APIResponse{
		HTTPCode: http.StatusOK,
		Code:     constants.CodeGeneralSuccess,
		Message:  constants.MessageGeneralSuccess,
	}

	// The request has been fulfilled, resulting in the creation of a new resource
	APICreated = APIResponse{
		HTTPCode: http.StatusCreated,
		Code:     constants.CodeGeneralSuccess,
		Message:  constants.MessageGeneralSuccess,
	}

	// The request has been accepted for processing, but the processing has not been completed.
	APIAccepted = APIResponse{
		HTTPCode: http.StatusAccepted,
		Code:     constants.CodeGeneralSuccess,
		Message:  constants.MessageGeneralSuccess,
	}

	// The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
	// , size too large, invalid request message framing, or deceptive request routing).
	APIErrorValidation = APIResponse{
		HTTPCode: http.StatusBadRequest,
		Code:     constants.CodeValidationError,
		Message:  constants.MessageGeneralError,
	}

	// The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
	// , size too large, invalid request message framing, or deceptive request routing).
	APIErrorInvalidPassword = APIResponse{
		HTTPCode: http.StatusBadRequest,
		Code:     constants.CodeInvalidAuthentication,
		Message:  constants.MessageInvalidData,
	}

	APIErrorInvalidData = APIResponse{
		HTTPCode: http.StatusBadRequest,
		Code:     constants.CodeInvalidData,
		Message:  constants.MessageInvalidData,
	}

	// The request was valid, but the server is refusing action. The user might not have the necessary permissions for
	// a resource, or may need an account of some sort.
	APIErrorForbidden = APIResponse{
		HTTPCode: http.StatusForbidden,
		Code:     constants.CodeForbidden,
		Message:  constants.MessageForbidden,
	}

	// The request was valid, but the server is refusing action. The user might not have the necessary permissions for
	// a resource, or may need an account of some sort.
	APIErrorUnauthorized = APIResponse{
		HTTPCode: http.StatusUnauthorized,
		Code:     constants.CodeUnauthorized,
		Message:  constants.MessageUnauthorized,
	}
)

// WithMessage modifies api response's message
func (a *APIResponse) WithMessage(message interface{}) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = a.HTTPCode
	new.Code = a.Code
	new.Message = message
	new.Data = a.Data

	return *new
}

// WithData if have data to response
func (a *APIResponse) WithData(data interface{}) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = a.HTTPCode
	new.Code = a.Code
	new.Message = a.Message
	new.Data = data

	return *new
}
