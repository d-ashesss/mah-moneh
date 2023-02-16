package main

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/api"
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
	api    *api.Service
}

func NewApp(config *Config, accountService *api.Service) *App {
	server := &http.Server{
		Addr: ":" + config.Port,
	}
	app := &App{
		config: config,
		server: server,
		api:    accountService,
	}
	app.registerHttpHandlers()
	return app
}

func (a *App) Run() {
	log.Println("[app] starting")
	signalCtx, signalStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	defer signalStop()

	wg, gCtx := errgroup.WithContext(signalCtx)

	wg.Go(func() error {
		log.Printf("[app] starting http server on port %s", a.config.Port)
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("[app] http server stopped unexpectedly")
			return err
		}
		return nil
	})
	wg.Go(func() error {
		<-gCtx.Done()
		log.Println("[app] shutting down http server")
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
		log.Println("[app] unexpected exit reason: ", err)
	}
}
