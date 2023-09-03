package poll

import (
	"awesome-go-bot-refactored/service/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type poll struct {
	chat.Info
}

func NewPollingChat(update *tgbotapi.Update) chat.Info {
	return &poll{
		Info: &chat.Chat{
			ChatId: update.Message.Chat.ID,
			Query:  update.Message.Text,
			Inline: false,
		},
	}
}
