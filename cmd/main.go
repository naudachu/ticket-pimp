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
	env(".dev.env")
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

	h := handler.NewHandler(
		os.Getenv("GIT_BASE_URL"),
		os.Getenv("GIT_TOKEN"),
		os.Getenv("CLOUD_BASE_URL"),
		os.Getenv("CLOUD_USER"),
		os.Getenv("CLOUD_PASS"),
		os.Getenv("YT_URL"),
		os.Getenv("YT_TOKEN"),
	)

	router := tgb.NewRouter().
		Message(h.NewTicketHandler, tgb.TextHasPrefix("/new")).
		Message(h.PingHandler, tgb.Command("ping")).
		Message(h.NewRepoHandler, tgb.TextHasPrefix("/repo")).
		Message(h.NewFolderHandler, tgb.TextHasPrefix("/folder"))

	return tgb.NewPoller(
		router,
		client,
	).Run(ctx)
}
