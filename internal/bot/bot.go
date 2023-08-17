package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

func init() {
	// Get bot Token from env
	config = &botConfig{
		BotToken: os.Getenv("TOKEN"),
	}

	// Create bot instance
	if botSession == nil {
		botSession = newBot()
	}
}

// newBot will initialize bot instance with
// the token provided in the env
func newBot() *tgbotapi.BotAPI {
	// Check token environment variable read status
	if config.BotToken == "" {
		panic("initBot: empty bot token")
	}

	// Bot instance initializer
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		panic(fmt.Errorf("initBot: error initializing bot: %v", err))
	}
	return bot
}

func defaultMessageConfiguration(chatID int64, text string) tgbotapi.MessageConfig {
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

// SendMessage will send msg string to user with userid
func SendMessage(msg string, userid int) error {
	// Configure msg parameters such as mode, webpreview
	msgCfg := defaultMessageConfiguration(int64(userid), msg)

	// Send message to the respective userid
	_, err := botSession.Send(msgCfg)
	if err != nil {
		return fmt.Errorf("sendmessage: message sending failed: %v", err)
	}

	return nil
}
