package awesome_go_bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/internal/helpers"
	"github.com/samirkape/awesome-go-bot/internal/logger"
	"io"
	"log"
	"net/http"
)

func HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	request, err := parseRequest(r.Body)
	if err != nil {
		log.Fatal("parsing error")
	}
	err = helpers.ExecuteCommand(request)
	if err != nil {
		logger.FieldLogger("failed for chatId: ", request.Message.Chat.ID).Error(err)
	}
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
