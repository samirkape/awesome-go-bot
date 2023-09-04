package main

import (
	"awesome-go-bot-refactored/gopackage/mongodb"
	"awesome-go-bot-refactored/gopackage/search"
	"awesome-go-bot-refactored/internal/service/chat/factory"
	"awesome-go-bot-refactored/internal/service/chat/inline"
	"awesome-go-bot-refactored/internal/service/chat/regular"
	"awesome-go-bot-refactored/internal/service/gobot"
	"awesome-go-bot-refactored/internal/service/gobot/config"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func main() {
	mybot, err := gobot.New(&config.Config{Token: "6027173597:AAH-PIKO9IVYS2LBUzTuXlojJYTImQkehdw"})
	if err != nil {
		return
	}

	URL, found := os.LookupEnv("ATLAS_URI")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}
	config := mongodb.NewConfig(mongodb.TABLENAME, URL)
	client, err := mongodb.New(config)
	if err != nil {
		return
	}

	packageService, err := client.GetAllPackages()

	searchService := search.NewService(packageService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates := mybot.GetUpdatesChan(u)
	for update := range updates {
		if update.InlineQuery != nil {
			chat := factory.New(context.Background(), &update)
			if chat.IsInline() {
				err := inline.HandleQuery(mybot, searchService, chat)
				if err != nil {
					log.Println(err)
				}
			} else {
				err := regular.HandleQuery(mybot, packageService, chat)
				if err != nil {
					log.Println(err)
				}
			}

		}
	}

}
