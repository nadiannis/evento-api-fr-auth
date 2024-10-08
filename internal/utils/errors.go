package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain/response"
	"github.com/rs/zerolog/log"
)

var (
	ErrCustomerNotFound           = errors.New("customer not found")
	ErrTicketTypeNotFound         = errors.New("ticket type not found")
	ErrTicketNotFound             = errors.New("ticket not found")
	ErrEventNotFound              = errors.New("event not found")
	ErrOrderNotFound              = errors.New("order not found")
	ErrCustomerAlreadyExists      = errors.New("customer already exists")
	ErrTicketTypeAlreadyExists    = errors.New("ticket type already exists")
	ErrTicketAlreadyExists        = errors.New("ticket already exists for the event")
	ErrInsufficientTicketQuantity = errors.New("insufficient ticket quantity")
	ErrInsufficientBalance        = errors.New("insufficient balance")
	ErrInvalidID                  = errors.New("invalid id")
	ErrInvalidAction              = errors.New("invalid action")
	ErrInvalidCredentials         = errors.New("invalid authentication credentials")
	ErrUnknownClaimsType          = errors.New("unknown claims type")
)

func errorResponse(c *gin.Context, status int, message any) {
	res := response.ErrorResponse{
		Status:  "error",
		Message: strings.ToLower(http.StatusText(status)),
		Detail:  message,
	}

	WriteJSON(c, status, res)
}

func ServerErrorResponse(c *gin.Context, err error) {
	req := fmt.Sprintf("%s %s %s", c.Request.Proto, c.Request.Method, c.Request.RequestURI)
	log.Error().Str("request", req).Msg(err.Error())

	message := "server encountered a problem"
	errorResponse(c, http.StatusInternalServerError, message)
}

func BadRequestResponse(c *gin.Context, err error) {
	errorResponse(c, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(c *gin.Context, err error) {
	errorResponse(c, http.StatusNotFound, err.Error())
}

func FailedValidationResponse(c *gin.Context, errors map[string]string) {
	errorResponse(c, http.StatusUnprocessableEntity, errors)
}

func InvalidCredentialsResponse(c *gin.Context, err error) {
	errorResponse(c, http.StatusUnauthorized, err.Error())
}

func InvalidAuthenticationTokenResponse(c *gin.Context, err error) {
	req := fmt.Sprintf("%s %s %s", c.Request.Proto, c.Request.Method, c.Request.RequestURI)
	log.Error().Str("request", req).Msg(err.Error())

	c.Header("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	errorResponse(c, http.StatusUnauthorized, message)
}
