package constant

const SupportedCommands = `
Supported commands:

- /listcategories - List all categories

- get packages for category - e.g /5

- /topN - e.g /top10 or /top

- /description - Get description of the bot

- /commands - Get all supported commands

- @awsmgo_bot <i>query</i>- copy @awsmgo_bot, paste it in Message box and query name or prefix of the package you want to search for
`

const StartMessage = "Hello! I can help you find Go packages, frameworks, and libraries. Send /commands to see what I can do."

const Description = `I can provide you with brief information about over 2500 Go packages, frameworks, and libraries scraped from awesome-go.com. This can be a helpful resource for learning about the Go community's contributions in your free time.

To use the bot:

Send the command /listcategories and reply with the number of any category. For example, to list all the packages for the "Actual Middlewares" category, you would reply with 0.
You can also get information about the top repositories written in Go by replying with /top or you can cap it by providing N suffix e.g /top5 

For more command options, send /start.
`

const CommandPrefix = "/"

const CategoryHelper = `Now respond with the number of the category you want to see. For example, to list all the packages for the <i> Zero Trust </i> category, you would reply with /110.`
