// Package gobot will be responsible for creation of client
// and router that will handle the incoming requests
package gobot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/gobot/config"
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
		ParseMode:             tgbotapi.ModeMarkdown,
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
		ParseMode:             tgbotapi.ModeMarkdown,
		Text:                  text,
		DisableWebPagePreview: true,
	}
}
