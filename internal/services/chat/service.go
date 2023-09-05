package chat

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Info interface {
	GetQuery() string
	GetChatId() int64
	GetQueryId() string
	IsInline() bool
	IsCallBack() bool
	GetCallBackQuery() *tgbotapi.CallbackQuery
	SetMessageId(int)
	GetMessageId() int
	HandleQuery() error
}

type Chat struct {
	ChatId        int64
	Query         string
	QueryId       string
	Inline        bool
	CallBack      bool
	MessageId     int
	CallBackQuery *tgbotapi.CallbackQuery
}

func (c *Chat) GetQuery() string {
	return c.Query
}

func (c *Chat) GetChatId() int64 {
	return c.ChatId
}

func (c *Chat) IsInline() bool {
	return c.Inline
}

func (c *Chat) GetQueryId() string {
	return c.QueryId
}

func (c *Chat) IsCallBack() bool {
	return c.CallBack
}

func (c *Chat) GetCallBackQuery() *tgbotapi.CallbackQuery {
	return c.CallBackQuery
}

func (c *Chat) SetMessageId(id int) {
	c.MessageId = id
}

func (c *Chat) GetMessageId() int {
	return c.MessageId
}

func (c *Chat) HandleQuery() error {
	return c.HandleQuery()
}
