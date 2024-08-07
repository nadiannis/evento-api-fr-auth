package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain/response"
	"github.com/nadiannis/evento-api-fr-auth/internal/utils"
	"github.com/rs/zerolog/log"
)

func (app *application) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := fmt.Sprintf("%dns", time.Since(start).Nanoseconds())
		request := fmt.Sprintf("%s %s %s", c.Request.Proto, c.Request.Method, c.Request.RequestURI)
		message := utils.GetLogMessage(c)

		status := response.Success
		logEvent := log.Info()
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			status = response.Error
			logEvent = log.Error()
		}

		logEvent.
			Str("request", request).
			Str("status", string(status)).
			Int("status_code", statusCode).
			Str("status_description", strings.ToLower(http.StatusText(statusCode))).
			Interface("message", message).
			Str("process_time", duration).
			Send()
	}
}

func (app *application) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Authorization")

		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			utils.InvalidAuthenticationTokenResponse(c, errors.New("no authorization header value"))
			c.Abort()
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.InvalidAuthenticationTokenResponse(c, errors.New("invalid authorization header value"))
			c.Abort()
			return
		}

		token := headerParts[1]

		claims, err := utils.ValidateJWTToken(app.config.JWT.Secret, token)
		if err != nil {
			utils.InvalidAuthenticationTokenResponse(c, err)
			c.Abort()
			return
		}

		subject, err := claims.GetSubject()
		if err != nil {
			utils.InvalidAuthenticationTokenResponse(c, err)
			c.Abort()
			return
		}

		customerID, err := strconv.ParseInt(subject, 10, 64)
		if err != nil {
			utils.ServerErrorResponse(c, err)
			c.Abort()
			return
		}

		customer, err := app.usecases.Customers.GetByID(customerID)
		if err != nil {
			switch {
			case errors.Is(err, utils.ErrCustomerNotFound):
				utils.InvalidAuthenticationTokenResponse(c, err)
			default:
				utils.ServerErrorResponse(c, err)
			}
			c.Abort()
			return
		}

		c.Set("customer", customer)
		c.Next()
	}
}
