package awesome_go_bot

import (
	mongodb2 "awesome-go-bot/domain/gopackage/mongodb"
	"awesome-go-bot/gobot"
	"awesome-go-bot/gobot/config"
	"awesome-go-bot/internal/logger"
	"awesome-go-bot/internal/services/chat/factory"
	"awesome-go-bot/internal/services/packages"
	"awesome-go-bot/internal/services/packages/analytics"
	"awesome-go-bot/internal/services/packages/search"
	"context"
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
)

var queryError = errors.New("unable to handle query")

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request, err := parseRequest(r.Body)
	if err != nil {
		log.Fatal("parsing error")
	}
	err = ExecuteCommand(ctx, request)
	if err != nil {
		logger.FieldLogger("failed for chatId: ", request.Message.Chat.ID).Error(err)
	}
}

func ExecuteCommand(ctx context.Context, incomingRequest *tgbotapi.Update) error {
	// create new bot
	botService, err := gobot.New(config.NewDefaultConfig())
	if err != nil {
		return err
	}
	// create new mongodb client
	dbClient, err := mongodb2.New(mongodb2.NewDefaultConfig())
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
	// create new chat
	chat := factory.NewChatService(incomingRequest, packageService, analyticsService, searchService, botService)
	chat.HandleQuery()
	return nil
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
