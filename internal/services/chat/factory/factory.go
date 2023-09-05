package factory

import (
	"awesome-go-bot/internal/services/chat"
	"awesome-go-bot/internal/services/chat/inline"
	"awesome-go-bot/internal/services/chat/keyboard"
	"awesome-go-bot/internal/services/chat/regular"
	"awesome-go-bot/internal/services/packages"
	"awesome-go-bot/internal/services/packages/analytics"
	"awesome-go-bot/internal/services/packages/search"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewChatService(request *tgbotapi.Update, packageService packages.Service, analyticsService analytics.Service, searchService search.Service, botService *tgbotapi.BotAPI) chat.Info {
	if request.InlineQuery != nil {
		return inline.NewInlineChat(request, searchService, botService)
	} else if request.CallbackQuery != nil {
		return keyboard.NewDefaultKeyboardChat(request, analyticsService, botService)
	} else {
		return regular.NewRegularChat(request, analyticsService, botService)
	}
}
