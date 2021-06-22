// types is a part of mparser, responsible for maintaining types, variables and constants
package mybot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// Config holds database related configurations such as
	// mongodb URI to establish connection with DB.
	DBConfig *dbconfig

	// DBclient requires one time initialization
	DBClient *mongo.Client

	// Bot config includes fetching token from env.
	BotConfig *botConfig

	// BotInstance requires one time initialization.
	BotInstance *tgbotapi.BotAPI

	// Command are communication interface of bot and the app
	BotCommand *botCommand

	// RequestCounter serve as a counter to count the user queries
	RequestCounter int
)

type Packages []Package

// If any category contains packages  more than `MaxAcceptable`
// Merge them into a group of `MergeMessages` and send as a single message
const (
	MaxAcceptable = 10
	MergeMessages = 10
)

// Below structs are used for parsing the incoming POST request from telegram bot.
// root level structure
type ReceiveMessage struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message struct
// holds information about complete message that includes chat id, msg text etc.
type Message struct {
	MessageID int        `json:"message_id"`
	From      From       `json:"from"`
	Chat      Chat       `json:"chat"`
	Date      int        `json:"date"`
	Text      string     `json:"text"`
	Entities  []Entities `json:"entities"`
}

// From struct
// holds information about the sender.
type From struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat struct
// it holds the meta data of chat that includes id, msg type (text, image, etc.)
type Chat struct {
	ID                          int    `json:"id"`
	FirstName                   string `json:"first_name"`
	UserName                    string `json:"username"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

// Entities struct
// Unused, written for future use
type Entities struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

// below data structures are related to parsing of markdown file
// root structure
type Category struct {
	Details Meta
}

// root structure of a category meta data
type Meta struct {
	Title    string
	Line     LineMeta
	SubTitle string
	Count    int
}

// this structure holds a information for multiple single lines.
// i.e it stores multiple raw lines related with package that belong to certain category.
type LineMeta struct {
	Packages []Package
	FullLink []string
}

// this is final structure of parser which will also be use for inserting package into database.
type Package struct {
	Name string             `bson:"name" json:"name"`
	URL  string             `bson:"url" json:"url"`
	Info string             `bson:"info" json:"info"`
	ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

// DB related structs
type UserRequestCounter struct {
	Count int `bson:"count" json:"count"`
}
