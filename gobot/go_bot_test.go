package gobot

import (
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/gobot/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock implementation of the config.Config interface.
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetToken() string {
	args := m.Called()
	return args.String(0)
}

func TestNew(t *testing.T) {
	// Arrange
	mockConfig := new(MockConfig)
	mockConfig.On("GetToken").Return("fake-token")

	// Act
	bot, err := New(mockConfig)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, bot)
	assert.Equal(t, "fake-token", bot.Token)
}

func TestDefaultMessageConfig(t *testing.T) {
	// Arrange
	chatID := int64(12345)
	text := "Hello, World!"

	// Act
	config := defaultMessageConfig(chatID, text)

	// Assert
	assert.Equal(t, chatID, config.ChatID)
	assert.Equal(t, 0, config.ReplyToMessageID)
	assert.Equal(t, tgbotapi.ModeMarkdown, config.ParseMode)
	assert.Equal(t, text, config.Text)
	assert.True(t, config.DisableWebPagePreview)
}

func TestDefaultEditMessageConfig(t *testing.T) {
	// Arrange
	chatID := int64(12345)
	messageID := 67890
	text := "Updated text"

	// Act
	config := defaultEditMessageConfig(chatID, messageID, text)

	// Assert
	assert.Equal(t, chatID, config.ChatID)
	assert.Equal(t, messageID, config.MessageID)
	assert.Equal(t, tgbotapi.ModeMarkdown, config.ParseMode)
	assert.Equal(t, text, config.Text)
	assert.True(t, config.DisableWebPagePreview)
}
