package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zip-archive_Golang/api"
	"zip-archive_Golang/model"
)

func main() {
	cfg := model.LoadConfig()
	router := api.NewRouter(cfg)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	go func() {
		log.Println("Server running on port", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {

			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		log.Fatal("server shutdown error:", err)
	}
}
