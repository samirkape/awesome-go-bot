package constant

const Start = `
	Supported commands:
	/listcategories - List all categories
	/top N - e.g /top10
	/description - Get description of the bot
	/start - Get started with the bot
	list packages for category - e.g /10 
`

const Description = `I can provide you with brief information about over 2,000 Go packages, frameworks, and libraries scraped from awesome-go.com. This can be a helpful resource for learning about the Go community's contributions in your free time.

To use the bot:

Send the command /listcategories and reply with the number of any category. For example, to list all the packages for the "Actual Middlewares" category, you would reply with 0.
You can also get information about the top N Go repositories by replying with top N. For example, to get the top 50 repositories, you would reply with top 50. The N value is capped at 200 to prevent the bot from sending too many messages at once. 

For more command options, send /start.
`

const DefaultTopNMessage = "Try /topN where N is a number between 1 and 200"
