package main

import (
	"database/sql"
	"flag"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nadiannis/evento-api-fr-auth/internal/config"
	"github.com/nadiannis/evento-api-fr-auth/internal/handler"
	"github.com/nadiannis/evento-api-fr-auth/internal/repository"
	"github.com/nadiannis/evento-api-fr-auth/internal/usecase"
	"github.com/rs/zerolog/log"
)

type application struct {
	config   *config.Config
	usecases usecase.Usecases
	handlers handler.Handlers
}

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "", "PostgreSQL data source name")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", "", "JWT secret")

	flag.Parse()

	db, err := openDB(&cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()
	log.Info().Msg("connected to database successfully")

	repos := repository.NewRepositories(db)
	usecases := usecase.NewUsecases(&cfg, repos)
	handlers := handler.NewHandlers(usecases)

	app := &application{
		config:   &cfg,
		usecases: usecases,
		handlers: handlers,
	}

	log.Info().Msg("add ticket types")
	prepopulateTicketTypes(usecases.TicketTypes)

	log.Info().Msg("add events and tickets")
	prepopulateEventsAndTickets(usecases.Events, usecases.Tickets)

	err = app.serve()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
