package main

import (
	"context"

	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"

	"wb-l0/internal/config"
	router "wb-l0/internal/http-server"
	"wb-l0/internal/http-server/handlers"
	"wb-l0/internal/repository"
	serviceN "wb-l0/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc, err := stan.Connect("test-cluster", "client-1", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Print(err.Error())
	}
	defer func() {
		err := sc.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}()

	db := config.DbConfig()

	service := serviceN.New(sc, repository.New(db))
	
	if err := service.StartCache(ctx); err != nil {
		log.Println("EFEFEWFEWFEFEFWEFWEFW")
		log.Fatalf("Failed to load cache from the database: %v", err)
	}	else {
		log.Println("amir ")
		
	}

	h := handlers.New(service.Cache)

	go func() {
		service.ListenNats(ctx, "order")
	}()

	go func() {
		<-ctx.Done()
		stop()
	}()

	router.Router(h)
}
