package commands

import (
	"errors"
	"fmt"
	"github.com/samirkape/awesome-go-bot/internal/bot"
	"github.com/samirkape/awesome-go-bot/internal/packages"
	"github.com/samirkape/awesome-go-bot/internal/repository"
	"github.com/samirkape/awesome-go-bot/internal/util"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Execute is the main function that handles all the user requests.
// When user posts some message to bot, it will be parsed and received in
// the botResponse struct that has userID to respond back and the message
// which you can find in the switch case defined in BotCMD to proccess the user request.
func Execute(userQuery *bot.Request, AllData packages.AllData) {
	var inputCommand = userQuery.Command
	var chatID = userQuery.ChatID
	var categories = AllData.CategoryList

	switch inputCommand {
	case bot.StartCmd:
		err := bot.SendMessage("hello 🖐, press command button to start", chatID)
		if err != nil {
			return
		}
	case bot.ListCategoriesCmd:
		err := bot.SendMessage(util.ListToMsg(categories), chatID)
		if err != nil {
			return
		}
	case bot.ListPackagesCmd:
		err := bot.SendMessage("reply with category number", chatID)
		if err != nil {
			return
		}
	case bot.TopNCmd:
		err := bot.SendMessage("reply with top #. e.g top 10", chatID)
		if err != nil {
			return
		}
	case bot.DescriptionCmd:
		err := bot.SendMessage(Description, chatID)
		if err != nil {
			return
		}
	default:
		handleDefaultCommand(inputCommand, chatID, categories)
	}
}

// We have set some default commands in the Bot config, such as, /listcategories.
// But there are some commands that needs some variable msg such as N or top N.
// ( where N is a number ). Those commands are handled here.
func handleDefaultCommand(inputCommand string, chatID int, colls []string) {
	// Check for unhandled command and invalid index number
	err := util.ValidateMessage(inputCommand)
	if err != nil {
		log.Println(err)
		return
	}

	if topN(inputCommand, chatID) {
		return
	}

	// Handle input number(s)
	categoryIdx := strings.Split(inputCommand, ",")
	if len(categoryIdx) > 0 {
		for _, e := range categoryIdx {

			// Index conversion from String to int
			index, err := strconv.Atoi(e)
			if err != nil {
				log.Printf("handleDefaultCommand: unable to convert msg to integer index: %v", err.Error())
				err := bot.SendMessage("invalid request, numeric input needed", chatID)
				if err != nil {
					return
				}
			}

			// Input validation: (min >= input number < max)
			if index >= len(colls) || index < 0 {
				ErrMsg := fmt.Sprintf("invalid request, expected range: {0 - %d}, given: %d ", len(colls)-1, index)
				err := bot.SendMessage(ErrMsg, chatID)
				if err != nil {
					return
				}
			}

			// Find Packages for respective category index number
			pkgs := repository.LocalPackageByIndex(index, AllPackages.AllPackages, colls)

			// If too many (>MaxAccepted) packages, merge them.
			if len(pkgs) > bot.MAXACCEPTABLE {
				limitPackageCount(pkgs, chatID, false)
			} else {
				for _, pkg := range pkgs {
					err := bot.SendMessage(pkg.PackageToMsg(false), chatID)
					if err != nil {
						return
					}
				}
			}
			err = bot.SendMessage(fmt.Sprintf("sent %d packages for %s", len(pkgs), colls[index]), chatID)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// topN responds to user request for top N packages.
// At the start of an instance, we load all the data
// from the mongodb repo and store it in AllPackages struct
// which has packages sorted by their repository stars. When user
// sends e.g Top 5, we send him the first 5 elements of SortedByStars struct.
func topN(command string, chatID int) bool {
	// input validation
	num, err := validateAndParseTopCommand(command)
	if err != nil {
		_ = bot.SendMessage(err.Error(), chatID)
		return false
	}

	sort.SliceStable(packages.SortedByStars, func(i, j int) bool {
		return packages.SortedByStars[i].Stars > packages.SortedByStars[j].Stars
	})

	pkgs := packages.SortedByStars[:num]
	if len(pkgs) > bot.MAXACCEPTABLE {
		limitPackageCount(pkgs, chatID, true)

	} else {
		for _, pkg := range pkgs {
			err := bot.SendMessage(pkg.PackageToMsg(true), chatID)
			if err != nil {
				return false
			}
		}
	}
	return true
}

func validateAndParseTopCommand(command string) (int, error) {
	var num int
	top := strings.ToLower(command)
	if !strings.HasPrefix(top, "top") {
		return 0, errors.New("invalid request, expected top N, given: " + command)
	}
	pattern := regexp.MustCompile("[0-9]+")
	numbers := pattern.FindAllString(command, -1)
	if len(numbers) > 0 {
		num, _ = strconv.Atoi(numbers[0])
		// Input validation: (min >= input number < max)
		if num >= bot.MAXTOPN || num < 0 {
			return 0, fmt.Errorf("invalid request, N is capped to < 200, given: %d ", num)
		}
	}
	return num, nil
}

// limitPackageCount limits the number of packages to be sent to user.
func limitPackageCount(packages packages.Packages, chatID int, forTop bool) {
	pIdx := 0
	mergedCount := int(math.Floor(float64(len(packages))/10)) + 1
	for pIdx = 0; pIdx < mergedCount; pIdx++ {
		start := pIdx * bot.MERGEN
		end := pIdx*bot.MERGEN + bot.MERGEN
		if end > len(packages) {
			end = len(packages)
		}
		mergedMsg := packages[start:end].PackagesToMsg(forTop)
		err := bot.SendMessage(mergedMsg, chatID)
		if err != nil {
			return
		}
	}
}
