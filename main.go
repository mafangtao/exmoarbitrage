package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"github.com/tusupov/exmoarbitrage/api"
	"github.com/tusupov/exmoarbitrage/config"
	"github.com/tusupov/exmoarbitrage/route"
	"github.com/tusupov/exmoarbitrage/service"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Server config
	cfg := config.Init()
	abs, err := filepath.Abs(cfg.TemplateDirectory)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Template directory: %s", abs)

	// API config
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	exmoApi := api.NewExmo(api.ExmoBaseUrl, client)

	// Service
	serviceApi := service.NewArbitrage(exmoApi)

	// Init route and view templates
	router, err := route.Init(cfg, serviceApi)
	if err != nil {
		log.Fatal(err)
	}

	// New server
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerPort),
		Handler: router,
	}

	// Start server
	go func() {
		log.Printf("Listening [%s] ...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Waiting stop signal
	<-stop

	// Safe shutdown server
	shutdown(srv, time.Second*59)

}

// Safe shutdown server
func shutdown(srv *http.Server, timeout time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("Shutdown with timeout: %s\n", timeout)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}

}
