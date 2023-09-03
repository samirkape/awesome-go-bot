package inline

import (
	"awesome-go-bot-refactored/service/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type inlineChat struct {
	chat.Info
}

func NewInlineChat(update *tgbotapi.Update) chat.Info {
	return &inlineChat{
		Info: &chat.Chat{
			QueryId: update.InlineQuery.ID,
			Query:   update.InlineQuery.Query,
			Inline:  true,
		},
	}
}
