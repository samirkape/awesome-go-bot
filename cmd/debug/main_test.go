package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	awesome_go_bot "github.com/samirkape/awesome-go-bot"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/config"
	"github.com/samirkape/awesome-go-bot/internal/errors"
	"github.com/samirkape/awesome-go-bot/internal/logger"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/factory"
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	mongodb2 "github.com/samirkape/awesome-go-bot/internal/services/packages/mongodb"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func TestMemory(t *testing.T) {
	// create new bot
	botService, err := gobot.New(config.NewConfig(os.Getenv("TRIAL_BOT_TOKEN")))
	if err != nil {
		t.Fatalf("unable to create bot")
	}
	// create new mongodb client
	dbClient, err := mongodb2.New(mongodb2.WithDefaultConfig())
	if err != nil {
		t.Fatalf("unable to create db client")
	}
	packageService := packages.NewService(dbClient)
	// get all analyticsService from the database
	analyticsService := analytics.NewService(packageService)
	if err != nil {
		t.Fatalf("unable to create analytics service")
	}
	// create new search service
	searchService := search.NewService(packageService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600
	updates := botService.GetUpdatesChan(u)

	for update := range updates {
		// create new chat
		newChat, err := factory.New(&update)
		if err != nil {
			errors.RespondToError(err, botService, newChat)
		}
		chatService, err := factory.NewService(newChat, analyticsService, searchService, botService)
		log.Println(err)
		logger.FieldLogger("query: ", chatService.GetQuery()).Info("query received")
		err = chatService.HandleQuery()
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestSearch(t *testing.T) {
	for i := 0; i < 5; i++ {
		start := time.Now()
		dbClient, _ := mongodb2.New(mongodb2.WithDefaultConfig())
		packageService := packages.NewService(dbClient)
		searchService := search.NewService(packageService)
		res := searchService.Search("kube")
		fmt.Println("time taken: ", time.Since(start))
		if len(res) > 0 {
			fmt.Println(res[0].Name)
		}
	}
}

func TestWebhook(t *testing.T) {
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
