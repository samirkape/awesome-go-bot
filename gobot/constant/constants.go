package constant

const SupportedCommands = `
Supported commands:

- /listcategories - List all categories

- get packages for category - e.g /5

- /topN - e.g /top10

- /description - Get description of the bot

- /commands - Get all supported commands

- @awsmgo_bot <i>query</i>- where query is the name or prefix of the package you want to search for
`

const Description = `I can provide you with brief information about over 2,000 Go packages, frameworks, and libraries scraped from awesome-go.com. This can be a helpful resource for learning about the Go community's contributions in your free time.

To use the bot:

Send the command /listcategories and reply with the number of any category. For example, to list all the packages for the "Actual Middlewares" category, you would reply with 0.
You can also get information about the top N Go repositories by replying with top N. For example, to get the top 50 repositories, you would reply with top 50. The N value is capped at 200 to prevent the bot from sending too many messages at once. 

For more command options, send /start.
`

const DefaultTopNMessage = "Try /topN where N is a number between 1 and 200 for example /top10"

const CommandPrefix = "/"

const CategoryHelper = `Now respond with the number of the category you want to see. For example, to list all the packages for the <i> Actual Middlewares </i> category, you would reply with /1.`
