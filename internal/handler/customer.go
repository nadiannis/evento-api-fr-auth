package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain/request"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain/response"
	"github.com/nadiannis/evento-api-fr-auth/internal/usecase"
	"github.com/nadiannis/evento-api-fr-auth/internal/utils"
)

type CustomerHandler struct {
	usecase usecase.ICustomerUsecase
}

func NewCustomerHandler(usecase usecase.ICustomerUsecase) ICustomerHandler {
	return &CustomerHandler{
		usecase: usecase,
	}
}

func (h *CustomerHandler) Login(c *gin.Context) {
	var input request.CustomerRequest

	err := utils.ReadJSON(c, &input)
	if err != nil {
		utils.BadRequestResponse(c, err)
		return
	}

	v := utils.NewValidator()

	v.Check(input.Username != "", "username", "username is required")
	v.Check(input.Password != "", "password", "password is required")

	if !v.Valid() {
		utils.FailedValidationResponse(c, v.Errors)
		return
	}

	token, err := h.usecase.Login(&input)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidCredentials):
			utils.InvalidCredentialsResponse(c, err)
		default:
			utils.ServerErrorResponse(c, err)
		}
		return
	}

	res := response.SuccessResponse{
		Status:  response.Success,
		Message: "customer authenticated successfully",
		Data:    token,
	}

	utils.WriteJSON(c, http.StatusOK, res)
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.usecase.GetAll()
	if err != nil {
		utils.ServerErrorResponse(c, err)
		return
	}

	time.Sleep(2 * time.Second) // Simulate real processing time

	cust, _ := c.Get("customer")
	fmt.Printf("CUSTOMER: %+v\n", cust)

	res := response.SuccessResponse{
		Status:  response.Success,
		Message: "customers retrieved successfully",
		Data:    customers,
	}

	utils.WriteJSON(c, http.StatusOK, res)
}

func (h *CustomerHandler) Add(c *gin.Context) {
	var input request.CustomerRequest

	err := utils.ReadJSON(c, &input)
	if err != nil {
		utils.BadRequestResponse(c, err)
		return
	}

	v := utils.NewValidator()

	v.Check(input.Username != "", "username", "username is required")
	v.Check(utils.Matches(input.Username, utils.UsernameRX), "username", "username is invalid")
	v.Check(input.Password != "", "password", "password is required")
	v.Check(len(input.Password) >= 8, "password", "password must be at least 8 characters long")
	v.Check(len(input.Password) <= 72, "password", "password must not be more than 72 characters long")

	if !v.Valid() {
		utils.FailedValidationResponse(c, v.Errors)
		return
	}

	customer, err := h.usecase.Add(&input)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrCustomerAlreadyExists):
			utils.BadRequestResponse(c, err)
		default:
			utils.ServerErrorResponse(c, err)
		}
		return
	}

	res := response.SuccessResponse{
		Status:  response.Success,
		Message: "customer added successfully",
		Data:    customer,
	}

	utils.WriteJSON(c, http.StatusCreated, res)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		utils.BadRequestResponse(c, utils.ErrInvalidID)
		return
	}

	customer, err := h.usecase.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrCustomerNotFound):
			utils.NotFoundResponse(c, err)
		default:
			utils.ServerErrorResponse(c, err)
		}
		return
	}

	res := response.SuccessResponse{
		Status:  response.Success,
		Message: "customer retrieved successfully",
		Data:    customer,
	}

	utils.WriteJSON(c, http.StatusOK, res)
}

func (h *CustomerHandler) UpdateBalance(c *gin.Context) {
	id, err := utils.ReadIDParam(c)
	if err != nil {
		utils.BadRequestResponse(c, utils.ErrInvalidID)
		return
	}

	var input request.CustomerBalanceRequest

	err = utils.ReadJSON(c, &input)
	if err != nil {
		utils.BadRequestResponse(c, err)
		return
	}

	v := utils.NewValidator()

	v.Check(input.Action != "", "action", "action is required")
	v.Check(utils.PermittedValue(input.Action, request.ActionAdd, request.ActionDeduct), "action", "action should be 'add' or 'deduct'")
	v.Check(input.Balance != 0, "balance", "balance is required")
	v.Check(input.Balance > 0, "balance", "balance should not be a negative number")

	if !v.Valid() {
		utils.FailedValidationResponse(c, v.Errors)
		return
	}

	customer, err := h.usecase.UpdateBalance(id, &input)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrCustomerNotFound):
			utils.NotFoundResponse(c, err)
		case errors.Is(err, utils.ErrInsufficientBalance):
			utils.BadRequestResponse(c, err)
		default:
			utils.ServerErrorResponse(c, err)
		}
		return
	}

	res := response.SuccessResponse{
		Status:  response.Success,
		Message: "customer balance updated successfully",
		Data:    customer,
	}

	utils.WriteJSON(c, http.StatusOK, res)
}
