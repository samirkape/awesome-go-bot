package main

import (
	bot "github.com/samirkape/awesome-go-bot"
)

func main() {
	response := bot.BotResponse{"top 4", 1346530914}
	// Head package list from the databse
	allPackages := bot.GetAllData()

	// Handle command given in the msgText
	// e.g /listpackages, /getStats
	bot.ExecuteCommand(&response, allPackages)
}
