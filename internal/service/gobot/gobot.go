// gobot will be responsible for creation of client
// and router that will handle the incoming requests

package gobot

import (
	"awesome-go-bot-refactored/internal/service/chat"
	"awesome-go-bot-refactored/internal/service/gobot/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func New(config *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.GetToken())
	if err != nil {
		return nil, err
	}
	return bot, nil
}

// Respond will send msg string to user with userid
func Respond(chat chat.Info, bot *tgbotapi.BotAPI, messageText string) error {
	messageConfig := defaultMessageConfig(chat.GetChatId(), messageText)

	_, err := bot.Send(messageConfig)
	if err != nil {
		return fmt.Errorf("message sending failed: %v", err)
	}
	return nil
}

func defaultMessageConfig(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		ParseMode:             "markdown",
		Text:                  text,
		DisableWebPagePreview: true,
	}
}
