// Package gobot will be responsible for creation of client
// and router that will handle the incoming requests
package gobot

import (
	"awesome-go-bot/gobot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func New(config *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.GetToken())
	if err != nil {
		return nil, err
	}
	return bot, nil
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

func defaultEditMessageConfig(chatID int64, messageId int, text string) tgbotapi.EditMessageTextConfig {
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageId,
		},
		ParseMode:             "markdown",
		Text:                  text,
		DisableWebPagePreview: true,
	}
}
