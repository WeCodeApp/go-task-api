package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ErrBadRequest              = "Bad request"
	ErrRequestTimeout          = "Request Timeout"
	ErrInvalidField            = "Invalid field"
	ErrInvalidRequest          = "Bad request"
	ErrUnauthorized            = "Unauthorized"
	ErrAccessDenied            = "Access Denied"
	ErrUnsupportedResponseType = "Unsupported Response Type"
	ErrInvalidScope            = "Invalid Scope"
	ErrServerError             = "Server Error"
	ErrTemporarilyUnavailable  = "Service Unavailable"
	ErrInvalidClient           = "Invalid Client"
	ErrInvalidGrant            = "Invalid Grant"
	ErrUnsupportedGrantType    = "Unsupported Grant Type"
)

var (
	BadRequest            = errors.New("Bad request")
	WrongCredentials      = errors.New("Wrong Credentials")
	NotFound              = errors.New("Not Found")
	Unauthorized          = errors.New("Unauthorized")
	Forbidden             = errors.New("Forbidden")
	PermissionDenied      = errors.New("Permission Denied")
	ExpiredCSRFError      = errors.New("Expired CSRF token")
	WrongCSRFToken        = errors.New("Wrong CSRF token")
	CSRFNotPresented      = errors.New("CSRF not presented")
	NotRequiredFields     = errors.New("No such required fields")
	BadQueryParams        = errors.New("Invalid query params")
	InternalServerError   = errors.New("Internal Server Error")
	RequestTimeoutError   = errors.New("Request Timeout")
	ExistsEmailError      = errors.New("User with given email already exists")
	InvalidJWTToken       = errors.New("Invalid JWT token")
	InvalidJWTClaims      = errors.New("Invalid JWT claims")
	NotAllowedImageHeader = errors.New("Not allowed image header")
	NoCookie              = errors.New("not found cookie header")
	InvalidUUID           = errors.New("invalid uuid")

	// oauth2 errors
	InvalidRequest                 = errors.New("invalid_request")
	UnauthorizedClient             = errors.New("unauthorized_client")
	AccessDenied                   = errors.New("access_denied")
	UnsupportedResponseType        = errors.New("unsupported_response_type")
	InvalidScope                   = errors.New("invalid_scope")
	ServerError                    = errors.New("server_error")
	TemporarilyUnavailable         = errors.New("temporarily_unavailable")
	InvalidClient                  = errors.New("invalid_client")
	InvalidGrant                   = errors.New("invalid_grant")
	UnsupportedGrantType           = errors.New("unsupported_grant_type")
	CodeChallengeRquired           = errors.New("invalid_request")
	UnsupportedCodeChallengeMethod = errors.New("invalid_request")
	InvalidCodeChallengeLen        = errors.New("invalid_request")

	InvalidRedirectURI   = errors.New("invalid redirect uri")
	InvalidAuthorizeCode = errors.New("invalid authorize code")
	InvalidAccessToken   = errors.New("invalid access token")
	InvalidRefreshToken  = errors.New("invalid refresh token")
	ExpiredAccessToken   = errors.New("expired access token")
	ExpiredRefreshToken  = errors.New("expired refresh token")
	MissingCodeVerifier  = errors.New("missing code verifier")
	MissingCodeChallenge = errors.New("missing code challenge")
	InvalidCodeChallenge = errors.New("invalid code challenge")
)

// Rest error interface
type Error interface {
	Code() int
	Error() string
}

// Rest error struct
type HTTPError struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Error  Error() interface method
func (e HTTPError) Error() string {
	return fmt.Sprintf("code: %d - message: %s - data: %v", e.Status, e.Message, e.Data)
}

// Error status
func (e HTTPError) Code() int {
	return e.Status
}

// New Rest Error
func New(status int, err string, causes interface{}) Error {
	return HTTPError{
		Status:  status,
		Message: err,
		Data:    causes,
	}
}

// New Rest Error With Message
func NewWithMessage(status int, err string, causes interface{}) Error {
	return HTTPError{
		Status:  status,
		Message: err,
		Data:    causes,
	}
}

// New Rest Error From Bytes
func NewFromBytes(bytes []byte) (Error, error) {
	var apiErr HTTPError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

// New Bad Request Error
func NewBadRequestError(causes interface{}) Error {
	return HTTPError{
		Status:  http.StatusBadRequest,
		Message: BadRequest.Error(),
		Data:    causes,
	}
}

// New Not Found Error
func NewNotFoundError(causes interface{}) Error {
	return HTTPError{
		Status:  http.StatusNotFound,
		Message: NotFound.Error(),
		Data:    causes,
	}
}

// New Unauthorized Error
func NewUnauthorizedError(causes interface{}) Error {
	return HTTPError{
		Status:  http.StatusUnauthorized,
		Message: Unauthorized.Error(),
		Data:    causes,
	}
}

// New Forbidden Error
func NewForbiddenError(causes interface{}) Error {
	return HTTPError{
		Status:  http.StatusForbidden,
		Message: Forbidden.Error(),
		Data:    causes,
	}
}

// New Internal Server Error
func NewInternalServerError(causes interface{}) Error {
	result := HTTPError{
		Status:  http.StatusInternalServerError,
		Message: InternalServerError.Error(),
		Data:    causes,
	}
	return result
}

// Parser of error string messages returns HTTPError
func ParseErrors(err error) Error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return New(http.StatusRequestTimeout, ErrRequestTimeout, nil)
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		return parseValidatorError(err)
	case strings.Contains(strings.ToLower(err.Error()), "unmarshal"):
		return New(http.StatusBadRequest, ErrBadRequest, err)
	case strings.Contains(strings.ToLower(err.Error()), "invalid_request"):
		return New(http.StatusBadRequest, ErrInvalidRequest, err)
	case strings.Contains(strings.ToLower(err.Error()), "unauthorized_client"):
		return New(http.StatusUnauthorized, ErrUnauthorized, err)
	case strings.Contains(strings.ToLower(err.Error()), "access_denied"):
		return New(http.StatusForbidden, ErrAccessDenied, err)
	case strings.Contains(strings.ToLower(err.Error()), "unsupported_response_type"):
		return New(http.StatusUnauthorized, ErrUnsupportedResponseType, err)
	case strings.Contains(strings.ToLower(err.Error()), "invalid_scope"):
		return New(http.StatusBadRequest, ErrInvalidScope, err)
	case strings.Contains(strings.ToLower(err.Error()), "server_error"):
		return New(http.StatusInternalServerError, ErrServerError, err)
	case strings.Contains(strings.ToLower(err.Error()), "temporarily_unavailable"):
		return New(http.StatusServiceUnavailable, ErrTemporarilyUnavailable, err)
	case strings.Contains(strings.ToLower(err.Error()), "invalid_client"):
		return New(http.StatusUnauthorized, ErrInvalidClient, err)
	case strings.Contains(strings.ToLower(err.Error()), "invalid_grant"):
		return New(http.StatusUnauthorized, ErrInvalidGrant, err)
	case strings.Contains(strings.ToLower(err.Error()), "unsupported_grant_type"):
		return New(http.StatusUnauthorized, ErrUnsupportedGrantType, err)

	default:
		if restErr, ok := err.(Error); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func parseValidatorError(err error) Error {
	return New(http.StatusBadRequest, ErrInvalidField, err)
}

// Error response
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Code(), ParseErrors(err)
}

// Error response object and status code
func HTTPErrorResponse(ctx *gin.Context, err error) {
	restErr := ParseErrors(err)
	ctx.JSON(restErr.Code(), restErr)
}

func TokenErrorResponse(c *gin.Context, err error) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	c.Writer.Header().Set("Cache-Control", "no-store")
	c.Writer.Header().Set("Pragma", "no-cache")
	restErr := ParseErrors(err)
	c.JSON(restErr.Code(), restErr)
}
