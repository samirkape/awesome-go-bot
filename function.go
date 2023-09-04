package awesome_go_bot

import (
	"awesome-go-bot/gopackage/mongodb"
	"awesome-go-bot/gopackage/search"
	"awesome-go-bot/internal/logger"
	"awesome-go-bot/internal/service/chat"
	"awesome-go-bot/internal/service/chat/factory"
	"awesome-go-bot/internal/service/chat/inline"
	"awesome-go-bot/internal/service/chat/keyboard"
	"awesome-go-bot/internal/service/chat/regular"
	"awesome-go-bot/internal/service/gobot"
	"awesome-go-bot/internal/service/gobot/config"
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
	chatInfo := factory.New(ctx, request)
	err = ExecuteCommand(ctx, chatInfo)
	if err != nil {
		logger.FieldLogger("failed for chatId: ", request.Message.Chat.ID).Error(err)
	}
}

func ExecuteCommand(ctx context.Context, chat chat.Info) error {
	// create new bot
	botService, err := gobot.New(config.NewDefaultConfig())
	if err != nil {
		return err
	}
	// create new mongodb client
	client, err := mongodb.New(mongodb.NewDefaultConfig())
	if err != nil {
		return err
	}
	// get all packageService from the database
	packageService, err := client.GetAllPackages()
	if err != nil {
		return err
	}
	// create new search service
	searchService := search.NewService(packageService)

	// handle query
	if chat.IsInline() {
		err := inline.HandleQuery(botService, searchService, chat)
		if err != nil {
			return err
		}
	} else if chat.IsCallBack() {
		err := keyboard.ProcessUsingInlineKeyboard(botService, packageService, chat)
		if err != nil {
			return err
		}
	} else {
		err := regular.HandleQuery(botService, packageService, chat)
		if err != nil {
			return err
		}
	}
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
