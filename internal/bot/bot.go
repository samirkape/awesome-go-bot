package mybot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

// personal user id for debugging
const _MYUSERID = 1346530914

const StartCmd = "/start"
const ListCategoriesCmd = "/listcategories"
const ListPackagesCmd = "/selectentry"
const TopNCmd = "/topn"
const DescriptionCmd = "/description"

const MAXACCEPTABLE = 1
const MERGEN = 10
const MAXTOPN = 200

var botSession *tgbotapi.BotAPI
var config *botConfig

type BotResponse struct {
	Command string
	ChatID  int
}

type botConfig struct {
	BotToken string
}

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
