package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	mybot "github.com/samirkape/awesome-go-bot/internal/commands"
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

type Request struct {
	Command string
	ChatID  int
}

type RequestData struct {
	UpdateID int           `json:"update_id"`
	Message  mybot.Message `json:"message"`
}

type botConfig struct {
	BotToken string
}
