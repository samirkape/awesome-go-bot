package awesome_go_bot

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/domain/gopackage/mongodb"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/config"
	"github.com/samirkape/awesome-go-bot/internal/logger"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/factory"
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
	"io"
	"log"
	"net/http"
)

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
	dbClient, err := mongodb.New(mongodb.NewDefaultConfig())
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
	chat, err := factory.NewChatService(incomingRequest, packageService, analyticsService, searchService, botService)
	if chat == nil {
		return err
	}
	return chat.HandleQuery()
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
