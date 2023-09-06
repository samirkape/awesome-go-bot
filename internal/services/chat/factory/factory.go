package factory

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/inline"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/keyboard"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/regular"
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
)

func NewChatService(request *tgbotapi.Update, packageService packages.Service, analyticsService analytics.Service, searchService search.Service, botService *tgbotapi.BotAPI) (chat.Info, error) {
	if request.InlineQuery != nil {
		return inline.NewInlineChat(request, searchService, botService)
	} else if request.CallbackQuery != nil {
		return keyboard.NewDefaultKeyboardChat(request, analyticsService, botService)
	} else {
		return regular.NewRegularChat(request, analyticsService, botService)
	}
}
