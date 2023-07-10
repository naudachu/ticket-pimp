package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ticket-pimp/clients/bot/handler"
	"ticket-pimp/internal/controllers"
	"ticket-pimp/internal/controllers/controller"
	"ticket-pimp/internal/storage"

	"github.com/joho/godotenv"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	d "ticket-pimp/internal/domain"
)

func main() {
	log.Print("started")
	env(".dev.env")
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	taskStorage := storage.NewStorage(initDB())
	cloudController := controller.NewCloudController(os.Getenv("CLOUD_BASE_URL"), os.Getenv("CLOUD_USER"), os.Getenv("CLOUD_PASS"))
	gitController := controller.NewGitController(os.Getenv("GIT_BASE_URL"), os.Getenv("GIT_TOKEN"))
	workflow := controllers.NewWorkflowController(os.Getenv("YT_URL"), os.Getenv("YT_TOKEN"), taskStorage, gitController, cloudController)

	if err := runBot(ctx, gitController, cloudController, workflow, taskStorage); err != nil {
		fmt.Println(err)
		defer os.Exit(1)
	}
}

// env
// env function reads provided file and setup envirmental variables;
func env(envFilePath string) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error while loading env file")
	}
}

// runBot ...
// ..function creates new Telegram BOT instance
// ..throw env variables through bot's handlers
// ..setup tg bot router;
// and finally returns tgb.Poller
func runBot(
	ctx context.Context,
	git controller.RepoCreator,
	cloud controller.CloudCreator,
	workflow controllers.IWorkflowController,
	r storage.Storage,
) error {

	client := tg.New(os.Getenv("TG_API"))

	h := handler.NewHandler(git, cloud, workflow, r)

	router := tgb.NewRouter().
		Message(h.NewRepoHandler, tgb.TextHasPrefix("/repo")).
		Message(h.NewFolderHandler, tgb.TextHasPrefix("/folder")).
		Message(h.NewTicketHandler, tgb.TextHasPrefix("/new")).
		Message(h.NewTaskHandler, tgb.TextHasPrefix("/task")).
		Message(h.NewConversion, tgb.TextHasPrefix("/conversion")).
		Message(func(ctx context.Context, mu *tgb.MessageUpdate) error {
			return mu.Answer("...").DoVoid(ctx)
		}, tgb.Command("ping"))

	return tgb.NewPoller(
		router,
		client,
	).Run(ctx)
}

func initDB() *gorm.DB {

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_LINK")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&d.TaskEntity{}, &d.User{})
	return db
}
