package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"authy/auth"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	srvPort := os.Getenv("SERVER_PORT")
	if srvPort == "" {
		srvPort = "8080"
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", srvPort),
		Handler: auth.Router(),
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(s *http.Server) {
		log.Printf("authy http now available at %s\n", s.Addr)
		if serr := s.ListenAndServe(); serr != http.ErrServerClosed {
			log.Fatal(serr)
		}
	}(s)

	<-sigChan

	err := s.Shutdown(context.Background())
	if err != nil {
		log.Fatal("something wrong when stopping server : ", err)
		return
	}

	log.Printf("server stopped")
}
