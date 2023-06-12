package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ticket-pimp/handler"

	"github.com/joho/godotenv"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func main() {
	log.Print("started")
	env(".env")
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		defer os.Exit(1)
	}
}

func env(envFilePath string) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func run(ctx context.Context) error {
	client := tg.New(os.Getenv("TG_API"))

	router := tgb.NewRouter().
		Message(handler.NewTicketHandler, tgb.TextHasPrefix("/new")).
		Message(handler.PingHandler, tgb.Command("ping")).
		Message(handler.NewRepoHandler, tgb.TextHasPrefix("/repo"))

	return tgb.NewPoller(
		router,
		client,
	).Run(ctx)
}
