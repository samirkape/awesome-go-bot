package chat

type Info interface {
	GetQuery() string
	GetChatId() int64
	IsInline() bool
}

type Chat struct {
	ChatId int64
	Query  string
	Inline bool
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
