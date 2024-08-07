package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadiannis/evento-api-fr-auth/internal/utils"
)

func (app *application) routes() *gin.Engine {
	r := gin.New()

	r.Use(app.RequestLogger())
	r.Use(gin.Recovery())

	r.GET("/api", func(c *gin.Context) {
		message := fmt.Sprintf("API is running on port %d", app.config.Port)
		c.String(http.StatusOK, message)
		utils.SetLogMessage(c, message)
	})

	r.POST("/api/customers/authentication", app.handlers.Customers.Login)
	r.POST("/api/customers", app.handlers.Customers.Add)
	r.GET("/api/customers", app.Authenticate(), app.handlers.Customers.GetAll)
	r.GET("/api/customers/:id", app.Authenticate(), app.handlers.Customers.GetByID)
	r.PATCH("/api/customers/:id/balances", app.Authenticate(), app.handlers.Customers.UpdateBalance)

	r.GET("/api/events", app.handlers.Events.GetAll)
	r.GET("/api/events/:id", app.handlers.Events.GetByID)

	r.GET("/api/tickets", app.handlers.Tickets.GetAll)
	r.GET("/api/tickets/:id", app.handlers.Tickets.GetByID)
	r.PATCH("/api/tickets/:id/quantities", app.Authenticate(), app.handlers.Tickets.UpdateQuantity) // Intended solely for concurrency testing purpose

	r.GET("/api/orders", app.Authenticate(), app.handlers.Orders.GetAll)
	r.POST("/api/orders", app.Authenticate(), app.handlers.Orders.Add)
	r.DELETE("/api/orders", app.Authenticate(), app.handlers.Orders.DeleteAll) // Intended solely for concurrency testing purpose

	return r
}
