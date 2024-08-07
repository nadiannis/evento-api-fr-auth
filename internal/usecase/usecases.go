package usecase

import (
	"github.com/nadiannis/evento-api-fr-auth/internal/config"
	"github.com/nadiannis/evento-api-fr-auth/internal/repository"
)

type Usecases struct {
	Customers   ICustomerUsecase
	Events      IEventUsecase
	TicketTypes ITicketTypeUsecase
	Tickets     ITicketUsecase
	Orders      IOrderUsecase
}

func NewUsecases(config *config.Config, repositories repository.Repositories) Usecases {
	return Usecases{
		Customers:   NewCustomerUsecase(config, repositories.Customers, repositories.Orders),
		Events:      NewEventUsecase(repositories.Events, repositories.Tickets),
		TicketTypes: NewTicketTypeUsecase(repositories.TicketTypes),
		Tickets:     NewTicketUsecase(repositories.Tickets, repositories.TicketTypes, repositories.Events),
		Orders: NewOrderUsecase(
			repositories.Orders,
			repositories.Customers,
			repositories.Tickets,
			repositories.TicketTypes,
		),
	}
}
