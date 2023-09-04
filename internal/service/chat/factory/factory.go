package factory

import (
	"awesome-go-bot/internal/service/chat"
	"awesome-go-bot/internal/service/chat/inline"
	"awesome-go-bot/internal/service/chat/keyboard"
	"awesome-go-bot/internal/service/chat/regular"
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func New(ctx context.Context, request *tgbotapi.Update) chat.Info {
	if request.InlineQuery != nil {
		return inline.NewInlineChat(request)
	} else if request.CallbackQuery != nil {
		return keyboard.NewKeyboardChat(request)
	} else {
		return regular.NewRegularChat(request)
	}
}
