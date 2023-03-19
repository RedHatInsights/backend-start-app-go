package main

import (
	"consoledot-go-template/internal/config"
	"consoledot-go-template/internal/db"
	"consoledot-go-template/internal/logging"
	"consoledot-go-template/internal/routes"
	// DAO import for pgx implementation
	_ "consoledot-go-template/internal/dao/pgx"

	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	mainCtx := context.Background()
	config.Initialize("config/api.env")

	logger, closeFn := logging.InitializeLogger()
	defer closeFn()
	log.Logger = logger

	// initialize the rest
	err := db.Initialize(mainCtx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
		panic(err)
	}
	defer db.Close()

	//log.Info().Msgf("Starting an instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	log.Info().Msgf("Starting an instance on port %d", config.Application.Port)
	router := routes.RootRouter()
	apiServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Application.Port),
		Handler: router,
	}

	waitForSignal := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		if err := apiServer.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Main service shutdown error")
		}
		close(waitForSignal)
	}()

	if err := apiServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Main service listen error")
		}
	}

	<-waitForSignal
}
