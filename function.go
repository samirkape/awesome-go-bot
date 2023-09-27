package awesome_go_bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/config"
	"github.com/samirkape/awesome-go-bot/internal/errors"
	"github.com/samirkape/awesome-go-bot/internal/logger"
	chatfactory "github.com/samirkape/awesome-go-bot/internal/services/chat/factory"
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	mongodb "github.com/samirkape/awesome-go-bot/internal/services/packages/mongodb"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
	"io"
	"log"
	"net/http"
)

func HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	request, err := parseRequest(r.Body)
	if err != nil {
		log.Fatal("parsing error")
	}
	err = ExecuteCommand(request)
	if err != nil {
		logger.FieldLogger("failed for chatId: ", request.Message.Chat.ID).Error(err)
	}
}

func ExecuteCommand(incomingRequest *tgbotapi.Update) error {
	// create new bot
	botService, err := gobot.New(config.WithDefaultConfig())
	if err != nil {
		return err
	}
	// create new chat
	chatInfo, err := chatfactory.New(incomingRequest)
	if err != nil {
		errors.RespondToError(err, botService, chatInfo)
		return err
	}
	// create new mongodb client
	client, err := mongodb.New(mongodb.WithDefaultConfig())
	if err != nil {
		return err
	}
	// create package service from mongodb client
	packageService := packages.NewService(client)
	// create analytics interface from package service
	analyticsService := analytics.NewService(packageService)
	if err != nil {
		return err
	}
	// create new search service based on package service
	searchService := search.NewService(packageService)
	// create new chat service
	chatService, err := chatfactory.NewService(chatInfo, analyticsService, searchService, botService)
	if chatService == nil {
		return err
	}
	logger.FieldLogger("query: ", chatService.GetQuery()).Info("query received")
	return chatService.HandleQuery()
}

func parseRequest(body io.ReadCloser) (*tgbotapi.Update, error) {
	var update *tgbotapi.Update
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&update)
	if err != nil {
		return nil, err
	}
	return update, nil
}
