package mybot

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

// adding personal user id for debugging perpose
const _MYUSERID = 1346530914

type (
	botResponse struct {
		msgText string
		chatID  int
	}

	// Commands to interact with the bot.
	// they need to be defined in telegram bot settings and handled in code.
	// we get command as a message text in a post request from telegram bot.
	botCommand struct {
		Start          string
		ListCategories string
		ListPackages   string
		GetStats       string
		TopN           string
	}

	// Bot config
	// need to be read from enviroment as a "TOKEN"
	botConfig struct {
		BotToken string
	}
)

func init() {
	// Get bot Token from env
	BotConfig = &botConfig{
		BotToken: os.Getenv("TOKEN"),
	}

	// Create bot instance
	BotInstance = botInit()

	// Set bot commands
	BotCommand = &botCommand{
		Start:          "/start",
		ListCategories: "/listcategories",
		ListPackages:   "/selectentry",
		GetStats:       "/getstats",
		TopN:           "/topn",
	}
}

// Initialize and validate bot
func botInit() *tgbotapi.BotAPI {
	// Check token environment variable read status
	if len(BotConfig.BotToken) == 0 {
		panic("initBot: empty bot token")
	}

	// Bot instance initializer
	bot, err := tgbotapi.NewBotAPI(BotConfig.BotToken)
	if err != nil {
		panic(fmt.Errorf("initBot: error initializing bot: %v", err))
	}

	return bot
}

func executeCommand(response *botResponse, AllData allData) {

	var msgText = response.msgText
	var chatID = response.chatID
	var categories = AllData.CategoryList

	switch msgText {
	case BotCommand.Start:
		SendMessage("Hello, press command button to start", chatID)
	case BotCommand.ListCategories:
		SendMessage("Hold on", chatID)
		SendMessage(listToMsg(categories), chatID)
		SendMessage("Done!", chatID)
		requestCounterIncr(chatID)
	case BotCommand.ListPackages:
		SendMessage("Reply with catergory number", chatID)
	case BotCommand.TopN:
		SendMessage("Reply with top #. e.g top 10", chatID)
	case BotCommand.GetStats:
		SendMessage(fmt.Sprintf("Total requests: %d", RequestCounter), chatID)
	default:
		handleDefaultCommand(msgText, chatID, categories)
	}

}

func checkTopN(msgText string, chatID int) bool {
	// Input validation: Reject response if any alphabet found in the package number
	top := strings.ToLower(msgText)
	if !strings.HasPrefix(top, "top") {
		return false
	}
	pattern := regexp.MustCompile("[0-9]+")
	numbers := pattern.FindAllString(msgText, -1)
	if len(numbers) > 0 {
		num, _ := strconv.Atoi(numbers[0])
		sort.SliceStable(StoreByStars, func(i, j int) bool {
			return StoreByStars[i].Stars > StoreByStars[j].Stars
		})
		pkgs := StoreByStars[:len(StoreByStars)-num]
		if len(pkgs) > MaxAcceptable {
			handleManyPkgs(pkgs, chatID)
		}
	}
	return true
}

func handleDefaultCommand(msgText string, chatID int, colls []string) {
	client := GetDbClient()
	userDBName := GetUserDbName()
	userDBColName := GetUserDbColName()

	// Check for unhandled command and invalid index number
	errString := validateMessage(msgText)
	if errString != "" {
		log.Println(errString)
		return
	}

	if checkTopN(msgText, chatID) {
		return
	}

	// Handle input number(s)
	categoryIdx := strings.Split(msgText, ",")
	if len(categoryIdx) > 0 {
		for _, e := range categoryIdx {

			// Index conversion from String to int
			index, err := strconv.Atoi(e)
			if err != nil {
				log.Println("handleDefaultCommand: Unable to convert msg to integer index")
				SendMessage("Invalid response, numeric input needed", chatID)
				return
			}

			// Input validation: (min >= input number < max)
			if index >= len(colls) || index < 0 {
				ErrMsg := fmt.Sprintf("Invalid response, expected: {0 - %d}, given: %d ", len(colls)-1, index)
				SendMessage(ErrMsg, chatID)
				return
			}

			// Find Packages for respective category index number
			pkgs := LocalPackageByIndex(index, AllData.AllPackages, colls)

			// If too many (>MaxAccepted) packages, merge them.
			if len(pkgs) > MaxAcceptable {
				handleManyPkgs(pkgs, chatID)
			} else {
				for _, pkg := range pkgs {
					SendMessage(pkg.packageToMsg(), chatID)
				}
			}
			SendMessage(fmt.Sprintf("Sent %d packages for %s", len(pkgs), colls[index]), chatID)
		}
	}

	requestCounterIncr(chatID)
	UpdateQueryCount(client, userDBName, userDBColName, bson.M{DBConfig.UserDBKey: RequestCounter})
}

// Check for unhandled command and invalid index number
func validateMessage(msgText string) string {
	// Input validation: Check if it is a unhandled scommand
	if strings.HasPrefix(msgText, "/") {
		return "Invalid response, try numeric input"
	}
	return ""
}

// Merge single Package struct elements into a single message string.
func (input Package) packageToMsg() string {
	msgString := strings.Builder{}
	msgString.WriteString(fmt.Sprintf("[%s](%s)\nStars: %d\n%s\n", input.Name, input.URL, input.Stars, input.Info))
	return msgString.String()
}

// The packagesToList method works on len(reciever)
// and merge them together
func (input Packages) packagesToMsg() string {
	msg := strings.Builder{}
	for _, pkg := range input {
		msg.WriteString(pkg.packageToMsg())
		msg.WriteString("\n")
	}
	return msg.String()
}

// Convert slice of strings into a single string.
func listToMsg(list []string) string {
	msg := strings.Builder{}
	for i, pkg := range list {
		msg.WriteString(fmt.Sprint(i) + ". " + pkg + string("\n")) // 3 = remove ## from start
	}
	return msg.String()
}

// If too many (>"MaxAccepted") packages, merge them
// into "MaxAccepted" packages per message. It calls
// packagesToList() for merging
func handleManyPkgs(p Packages, chatID int) {
	pidx := 0
	mergedCount := int(math.Floor(float64(len(p))/10)) + 1
	for pidx = 0; pidx < mergedCount; pidx++ {
		start := pidx * MergeMessages
		end := pidx*MergeMessages + MergeMessages
		if end > len(p) {
			end = len(p)
		}
		mergedMsg := Packages(p[start:end]).packagesToMsg()
		SendMessage(mergedMsg, chatID)
	}
}

// Keep track of requests apart
// from the one used for debugging and trials
func requestCounterIncr(chatID int) {
	if chatID != _MYUSERID {
		RequestCounter++
	}
}

// Initialize send message data structure which includes information such as
// user id, msg, parsemode, etc.
func newMessageInit(chatID int64, text string) tgbotapi.MessageConfig {
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
	msgcfg := newMessageInit(int64(userid), msg)

	// Send message to the respective userid
	_, err := BotInstance.Send(msgcfg)
	if err != nil {
		return fmt.Errorf("sendmessage: message sending failed: %v", err)
	}

	return nil
}
