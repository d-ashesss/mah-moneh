package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	config *Config
	server *http.Server
}

func NewApp(config *Config, handler http.Handler) *App {
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}
	app := &App{
		config: config,
		server: server,
	}
	return app
}

func (a *App) Run() {
	log.Print("[APP] Starting up")
	signalCtx, signalStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	defer signalStop()

	wg, gCtx := errgroup.WithContext(signalCtx)

	wg.Go(func() error {
		log.Printf("[APP] starting http server on port %s", a.config.Port)
		err := a.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Print("[APP] HTTP server has stopped unexpectedly")
			return err
		}
		return nil
	})
	wg.Go(func() error {
		<-gCtx.Done()
		log.Print("[APP] Shutting down HTTP server")
		serverCtx, serverCancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer serverCancel()
		return a.server.Shutdown(serverCtx)
	})
	wg.Go(func() error {
		<-gCtx.Done()
		signalStop()
		return nil
	})

	if err := wg.Wait(); err != nil {
		log.Printf("[APP] Unexpected exit reason: %s", err)
	}
}
