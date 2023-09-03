package awesome_go_bot

import (
	"awesome-go-bot-refactored/internal/logger"
	"awesome-go-bot-refactored/service/chat"
	"awesome-go-bot-refactored/service/chat/inline"
	"awesome-go-bot-refactored/service/chat/poll"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	chat := getChat(ctx, request)
	err = ExecuteCommand(ctx, chat)
	if err != nil {
		logger.FieldLogger("failed for chatId: ", request.Message.Chat.ID).Error(err)
	}
}

func ExecuteCommand(ctx context.Context, chat chat.Info) error {
	if chat.IsInline() {
		fmt.Println("inline query")
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

func getChat(ctx context.Context, request *tgbotapi.Update) chat.Info {
	if request.InlineQuery != nil {
		return inline.NewInlineChat(request)
	} else {
		return poll.NewPollingChat(request)
	}
}
