package factory

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/inline"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/keyboard"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/regular"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
)

func NewChatService(chat chat.Info, analyticsService analytics.Service, searchService search.Service, botService *tgbotapi.BotAPI) (chat.Info, error) {
	if chat.IsInline() {
		return inline.NewInlineChat(chat, searchService, botService)
	} else if chat.IsCallBack() {
		return keyboard.NewDefaultKeyboardChat(chat, analyticsService, botService)
	} else {
		return regular.NewRegularChat(chat, analyticsService, botService)
	}
}

func NewChat(request *tgbotapi.Update) (chat.Info, error) {
	if request.InlineQuery != nil {
		return inline.NewChat(request)
	} else if request.CallbackQuery != nil {
		return keyboard.NewChat(request)
	} else {
		return regular.NewValidatedChat(request)
	}
}
