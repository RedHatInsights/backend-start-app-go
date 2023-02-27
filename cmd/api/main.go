package main

import (
	"consoledot-go-template/internal/routes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	router := routes.RootRouter()

	apiServer := http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
		Handler: router,
	}

	waitForSignal := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		if err := apiServer.Shutdown(context.Background()); err != nil {
			//log.Fatal().Err(err).Msg("Main service shutdown error")
		}
		close(waitForSignal)
	}()

	if err := apiServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			//log.Fatal().Err(err).Msg("Main service listen error")
		}
	}

	<-waitForSignal
}
