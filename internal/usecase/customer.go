package usecase

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain/request"
	"github.com/nadiannis/evento-api-fr-auth/internal/domain/response"
	"github.com/nadiannis/evento-api-fr-auth/internal/repository"
	"github.com/nadiannis/evento-api-fr-auth/internal/utils"
)

type CustomerUsecase struct {
	customerRepository repository.ICustomerRepository
	orderRepository    repository.IOrderRepository
}

func NewCustomerUsecase(customerRepository repository.ICustomerRepository, orderRepository repository.IOrderRepository) ICustomerUsecase {
	return &CustomerUsecase{
		customerRepository: customerRepository,
		orderRepository:    orderRepository,
	}
}

func (u *CustomerUsecase) Login(input *request.CustomerRequest) (*string, error) {
	customer, err := u.customerRepository.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrCustomerNotFound):
			return nil, utils.ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	match, err := customer.Password.Matches(input.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, utils.ErrInvalidCredentials
	}

	claims := utils.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(customer.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api.evento.com",
			Audience:  []string{"api.evento.com"},
		},
	}

	token, err := utils.GenerateJWTToken(claims)

	return token, err
}

func (u *CustomerUsecase) GetAll() ([]*response.CustomerResponse, error) {
	customers, err := u.customerRepository.GetAll()
	if err != nil {
		return nil, err
	}

	customerResponses := make([]*response.CustomerResponse, 0)

	for _, customer := range customers {
		orders, err := u.orderRepository.GetByCustomerID(customer.ID)
		if err != nil {
			return nil, err
		}

		customerResponse := &response.CustomerResponse{
			ID:       customer.ID,
			Username: customer.Username,
			Balance:  customer.Balance,
			Orders:   orders,
		}
		customerResponses = append(customerResponses, customerResponse)
	}

	return customerResponses, nil
}

func (u *CustomerUsecase) Add(input *request.CustomerRequest) (*domain.Customer, error) {
	customer := &domain.Customer{
		Username: input.Username,
		Balance:  0,
	}

	err := customer.Password.Set(input.Password)
	if err != nil {
		return nil, err
	}

	err = u.customerRepository.Add(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (u *CustomerUsecase) GetByID(customerID int64) (*response.CustomerResponse, error) {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	orders, err := u.orderRepository.GetByCustomerID(customer.ID)
	if err != nil {
		return nil, err
	}

	customerResponse := &response.CustomerResponse{
		ID:       customer.ID,
		Username: customer.Username,
		Balance:  customer.Balance,
		Orders:   orders,
	}

	return customerResponse, nil
}

func (u *CustomerUsecase) UpdateBalance(customerID int64, input *request.CustomerBalanceRequest) (*domain.Customer, error) {
	_, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	var customer *domain.Customer

	switch input.Action {
	case request.ActionAdd:
		customer, err = u.customerRepository.AddBalance(customerID, input.Balance)
	case request.ActionDeduct:
		customer, err = u.customerRepository.DeductBalance(customerID, input.Balance)
	default:
		return nil, utils.ErrInvalidAction
	}

	return customer, err
}
