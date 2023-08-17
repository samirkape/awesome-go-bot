// types is a part of mparser, responsible for maintaining types, variables and constants
package mybot

import "github.com/samirkape/awesome-go-bot/internal/packages"

var (
	AllPackages packages.AllData
)

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

// DB related structs
type UserRequestCounter struct {
	Count int `bson:"count" json:"count"`
}

var Description string = `I can send you short details about 2000+ Go packages, frameworks, and libraries that I scraped from awesome-go.com
You can use this bot in your free time to get familiar with the Go community contribution.

How to use?

1. Click /listcategories and reply with any category number. 
e.g. Send 0 to list all the packages for *Actual Middlewares* 

2. You can also get information about the top number of Go repositories by replying with top N. e.g top 50
N is capped to 200 to stop bot from sending lots of messages at once 
`
