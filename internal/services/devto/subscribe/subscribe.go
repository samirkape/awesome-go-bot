package subscribe

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// Subscriber represents a user subscribed to receive articles
type Subscriber struct {
	gorm.Model
	ChatID   int64
	Interval string // "weekly", "biweekly"
	LastSent time.Time
}

// DevToArticle represents an article from dev.to API
type DevToArticle struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	PublishedAt string   `json:"published_at"`
	ReadingTime int      `json:"reading_time_minutes"`
	Tags        []string `json:"tags"`
}

// SubscriptionService handles article distribution
type SubscriptionService struct {
	db  *gorm.DB
	bot *tgbotapi.BotAPI
}

func NewSubscriptionService(db *gorm.DB, bot *tgbotapi.BotAPI) *SubscriptionService {
	return &SubscriptionService{
		db:  db,
		bot: bot,
	}
}

// Subscribe adds a new subscriber
func (s *SubscriptionService) Subscribe(chatID int64, interval string) error {
	subscriber := Subscriber{
		ChatID:   chatID,
		Interval: interval,
		LastSent: time.Now(),
	}

	return s.db.Create(&subscriber).Error
}

// Unsubscribe removes a subscriber
func (s *SubscriptionService) Unsubscribe(chatID int64) error {
	return s.db.Where("chat_id = ?", chatID).Delete(&Subscriber{}).Error
}

// FetchRandomArticle gets a random article from dev.to
func (s *SubscriptionService) FetchRandomArticle() (*DevToArticle, error) {
	// Fetch articles from last week to ensure freshness
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	url := fmt.Sprintf("https://dev.to/api/articles?per_page=100&published_after=%s", weekAgo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var articles []DevToArticle
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles found")
	}

	// Pick a random article
	randomIndex := rand.Intn(len(articles))
	return &articles[randomIndex], nil
}

// SendArticle sends an article to a specific subscriber
func (s *SubscriptionService) SendArticle(subscriber *Subscriber) error {
	article, err := s.FetchRandomArticle()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("📚 *%s*\n\n%s\n\n🕒 Reading time: %d minutes\n🏷 Tags: %s\n\n🔗 %s",
		article.Title,
		article.Description,
		article.ReadingTime,
		article.Tags,
		article.URL)

	msg := tgbotapi.NewMessage(subscriber.ChatID, message)
	msg.ParseMode = "Markdown"

	_, err = s.bot.Send(msg)
	if err != nil {
		return err
	}

	// Update LastSent time
	subscriber.LastSent = time.Now()
	return s.db.Save(subscriber).Error
}

// ProcessSubscriptions checks and sends articles to subscribers
func (s *SubscriptionService) ProcessSubscriptions() error {
	var subscribers []Subscriber
	if err := s.db.Find(&subscribers).Error; err != nil {
		return err
	}

	for _, subscriber := range subscribers {
		shouldSend := false
		switch subscriber.Interval {
		case "weekly":
			shouldSend = time.Since(subscriber.LastSent) >= 7*24*time.Hour
		case "biweekly":
			shouldSend = time.Since(subscriber.LastSent) >= 14*24*time.Hour
		}

		if shouldSend {
			if err := s.SendArticle(&subscriber); err != nil {
				fmt.Printf("Error sending article to %d: %v\n", subscriber.ChatID, err)
				continue
			}
		}
	}
	return nil
}
