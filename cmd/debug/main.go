package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	awesome_go_bot "github.com/samirkape/awesome-go-bot"
	"github.com/samirkape/awesome-go-bot/domain/gopackage/mongodb"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/config"
	"github.com/samirkape/awesome-go-bot/internal/logger"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/factory"
	"github.com/samirkape/awesome-go-bot/internal/services/internalerrors"
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func main() {
	_ = defaultTest()
}

func defaultTest() error {
	// create new bot
	botService, err := gobot.New(config.NewConfig(os.Getenv("TRIAL_BOT_TOKEN")))
	if err != nil {
		return err
	}
	// create new mongodb client
	dbClient, err := mongodb.New(mongodb.WithDefaultConfig())
	if err != nil {
		return err
	}
	packageService := packages.NewService(dbClient)
	// get all analyticsService from the database
	analyticsService := analytics.NewService(packageService)
	if err != nil {
		return err
	}
	// create new search service
	searchService := search.NewService(packageService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600
	updates := botService.GetUpdatesChan(u)

	for update := range updates {
		// create new chat
		newChat, err := factory.NewChat(&update)
		if err != nil {
			internalerrors.RespondToError(err, botService, newChat)
		}
		chatService, err := factory.NewChatService(newChat, analyticsService, searchService, botService)
		log.Println(err)
		logger.FieldLogger("query: ", chatService.GetQuery()).Info("query received")
		chatService.HandleQuery()
	}
	return nil
}

func searchTest() {
	for i := 0; i < 5; i++ {
		start := time.Now()
		dbClient, _ := mongodb.New(mongodb.WithDefaultConfig())
		packageService := packages.NewService(dbClient)
		searchService := search.NewService(packageService)
		res := searchService.Search("kube")
		fmt.Println("time taken: ", time.Since(start))
		if len(res) > 0 {
			fmt.Println(res[0].Name)
		}
	}
}

func webhookTest() {
	// Define the URL to which you want to send the request.
	url := "https://example.com/api/endpoint"

	// Define the payload as a string.
	var payload = `{
		"update_id":10000,
		"message":{
		  "date":1441645532,
		  "chat":{
			 "last_name":"Test Lastname",
			 "id":1111111,
			 "type": "private",
			 "first_name":"Test Firstname",
			 "username":"Testusername"
		  },
		  "message_id":1365,
		  "from":{
			 "last_name":"Test Lastname",
			 "id":1111111,
			 "first_name":"Test Firstname",
			 "username":"Testusername"
		  },
		  "text":"/start"
		}
		}`

	// Create a new POST request with the payload.
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	awesome_go_bot.HandleTelegramWebHook(nil, req)
}
