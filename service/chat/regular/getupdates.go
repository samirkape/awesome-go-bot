package regular

import (
	"awesome-go-bot-refactored/service/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type regular struct {
	chat.Info
}

func NewRegularChat(update *tgbotapi.Update) chat.Info {
	return &regular{
		Info: &chat.Chat{
			ChatId: update.Message.Chat.ID,
			Query:  update.Message.Text,
			Inline: false,
		},
	}
}
