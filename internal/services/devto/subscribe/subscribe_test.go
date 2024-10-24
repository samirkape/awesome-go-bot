package subscribe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockBotAPI implements necessary methods from tgbotapi.BotAPI
type MockBotAPI struct {
	mock.Mock
}

// Send implements the method from BotAPI interface
func (m *MockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(c)
	return args.Get(0).(tgbotapi.Message), args.Error(1)
}

// Request implements the method from BotAPI interface
func (m *MockBotAPI) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	args := m.Called(c)
	return args.Get(0).(*tgbotapi.APIResponse), args.Error(1)
}

// GetUpdatesChan implements the method from BotAPI interface
func (m *MockBotAPI) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	args := m.Called(config)
	return args.Get(0).(tgbotapi.UpdatesChannel)
}

// GetFile implements the method from BotAPI interface
func (m *MockBotAPI) GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error) {
	args := m.Called(config)
	return args.Get(0).(tgbotapi.File), args.Error(1)
}

// GetFileDirectURL implements the method from BotAPI interface
func (m *MockBotAPI) GetFileDirectURL(fileID string) (string, error) {
	args := m.Called(fileID)
	return args.String(0), args.Error(1)
}

// setupTestDB creates a test database using SQLite in-memory
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&Subscriber{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// setupMockDevToAPI creates a test server that returns mock Dev.to articles
func setupMockDevToAPI(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articles := []DevToArticle{
			{
				Title:       "Test Article 1",
				URL:         "https://dev.to/test1",
				Description: "Test Description 1",
				PublishedAt: time.Now().Format(time.RFC3339),
				ReadingTime: 5,
				Tags:        []string{"go", "testing"},
			},
			{
				Title:       "Test Article 2",
				URL:         "https://dev.to/test2",
				Description: "Test Description 2",
				PublishedAt: time.Now().Format(time.RFC3339),
				ReadingTime: 8,
				Tags:        []string{"api", "web"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}))
}

func TestSubscribe(t *testing.T) {
	db := setupTestDB(t)
	mockBot := &MockBotAPI{}
	service := NewSubscriptionService(db, mockBot)

	tests := []struct {
		name     string
		chatID   int64
		interval string
		wantErr  bool
	}{
		{
			name:     "Valid weekly subscription",
			chatID:   123456,
			interval: "weekly",
			wantErr:  false,
		},
		{
			name:     "Valid biweekly subscription",
			chatID:   789012,
			interval: "biweekly",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Subscribe(tt.chatID, tt.interval)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify subscriber was created
				var sub Subscriber
				result := db.Where("chat_id = ?", tt.chatID).First(&sub)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.chatID, sub.ChatID)
				assert.Equal(t, tt.interval, sub.Interval)
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	db := setupTestDB(t)
	mockBot := &MockBotAPI{}
	service := NewSubscriptionService(db, mockBot)

	// Create a test subscriber first
	testChatID := int64(123456)
	_ = service.Subscribe(testChatID, "weekly")

	tests := []struct {
		name    string
		chatID  int64
		wantErr bool
	}{
		{
			name:    "Unsubscribe existing user",
			chatID:  testChatID,
			wantErr: false,
		},
		{
			name:    "Unsubscribe non-existing user",
			chatID:  999999,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Unsubscribe(tt.chatID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify subscriber was deleted
				var sub Subscriber
				result := db.Where("chat_id = ?", tt.chatID).First(&sub)
				assert.Error(t, result.Error)
				assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
			}
		})
	}
}

func TestFetchRandomArticle(t *testing.T) {
	mockServer := setupMockDevToAPI(t)
	defer mockServer.Close()

	db := setupTestDB(t)
	mockBot := &MockBotAPI{}
	service := NewSubscriptionService(db, mockBot)

	t.Run("Successfully fetch random article", func(t *testing.T) {
		article, err := service.FetchRandomArticle()
		assert.NoError(t, err)
		assert.NotNil(t, article)
		assert.NotEmpty(t, article.Title)
		assert.NotEmpty(t, article.URL)
		assert.NotEmpty(t, article.Description)
		assert.NotEmpty(t, article.Tags)
	})
}

func TestSendArticle(t *testing.T) {
	db := setupTestDB(t)
	mockBot := &MockBotAPI{}
	service := NewSubscriptionService(db, mockBot)
	mockServer := setupMockDevToAPI(t)
	defer mockServer.Close()

	subscriber := &Subscriber{
		ChatID:   123456,
		Interval: "weekly",
		LastSent: time.Now().Add(-8 * 24 * time.Hour),
	}

	// Setup mock expectations
	mockBot.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).Return(
		tgbotapi.Message{}, nil,
	)

	t.Run("Successfully send article", func(t *testing.T) {
		err := service.SendArticle(subscriber)
		assert.NoError(t, err)
		mockBot.AssertExpectations(t)

		// Verify LastSent was updated
		assert.True(t, subscriber.LastSent.After(time.Now().Add(-1*time.Minute)))
	})
}

func TestProcessSubscriptions(t *testing.T) {
	db := setupTestDB(t)
	mockBot := &MockBotAPI{}
	service := NewSubscriptionService(db, mockBot)
	mockServer := setupMockDevToAPI(t)
	defer mockServer.Close()

	// Create test subscribers
	subscribers := []Subscriber{
		{
			ChatID:   123,
			Interval: "weekly",
			LastSent: time.Now().Add(-8 * 24 * time.Hour), // Should be processed
		},
		{
			ChatID:   456,
			Interval: "weekly",
			LastSent: time.Now().Add(-3 * 24 * time.Hour), // Should not be processed
		},
		{
			ChatID:   789,
			Interval: "biweekly",
			LastSent: time.Now().Add(-15 * 24 * time.Hour), // Should be processed
		},
	}

	for _, sub := range subscribers {
		db.Create(&sub)
	}

	// Setup mock expectations for messages that should be sent
	mockBot.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).Return(
		tgbotapi.Message{}, nil,
	).Times(2) // Expecting 2 messages (for subscribers 123 and 789)

	t.Run("Process subscriptions successfully", func(t *testing.T) {
		err := service.ProcessSubscriptions()
		assert.NoError(t, err)
		mockBot.AssertExpectations(t)

		// Verify that appropriate subscribers were processed
		var sub Subscriber
		db.First(&sub, "chat_id = ?", 123)
		assert.True(t, sub.LastSent.After(time.Now().Add(-1*time.Minute)))

		db.First(&sub, "chat_id = ?", 456)
		assert.True(t, sub.LastSent.Before(time.Now().Add(-2*24*time.Hour)))

		db.First(&sub, "chat_id = ?", 789)
		assert.True(t, sub.LastSent.After(time.Now().Add(-1*time.Minute)))
	})
}
