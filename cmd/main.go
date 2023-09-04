package main

import (
	"awesome-go-bot/gopackage/mongodb"
	"awesome-go-bot/gopackage/search"
	"awesome-go-bot/internal/service/chat/factory"
	"awesome-go-bot/internal/service/chat/inline"
	"awesome-go-bot/internal/service/chat/regular"
	"awesome-go-bot/internal/service/gobot"
	"awesome-go-bot/internal/service/gobot/config"
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
	URL, found := os.LookupEnv("ATLAS_URI")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}

	token, found := os.LookupEnv("TEST_TOKEN")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}

	goBot, err := gobot.New(&config.Config{Token: token})
	if err != nil {
		return
	}

	dbConfig := mongodb.NewConfig(mongodb.TABLENAME, URL)
	client, err := mongodb.New(dbConfig)
	if err != nil {
		return
	}

	packageService, err := client.GetAllPackages()

	searchService := search.NewService(packageService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates := goBot.GetUpdatesChan(u)
	for update := range updates {
		chat := factory.New(context.Background(), &update)
		if chat.IsInline() {
			err := inline.HandleQuery(goBot, searchService, chat)
			if err != nil {
				log.Println(err)
			}
		} else {
			err := regular.HandleQuery(goBot, packageService, chat)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
