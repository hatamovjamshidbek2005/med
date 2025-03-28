package main

import (
	"context"
	"fmt"
	_ "med/api/openapi"
	"med/internal/configs"
	handlers "med/internal/controllers"
	"os"
	"os/signal"
	"sync"
)

func main() {
	cfg, _ := configs.LoadConfig("./")
	server := handlers.NewServer(cfg)

	go func() {
		server.Run()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdown(server, ctx, cancel)
}

func gracefulShutdown(server handlers.Server, ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		fmt.Println("shutting down")
		server.Stop()
		fmt.Println("shutdown successfully called")
		wg.Done()
	}(&wg)

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
	os.Exit(0)
}
