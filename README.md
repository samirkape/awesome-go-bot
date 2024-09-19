[![Go Report Card](https://goreportcard.com/badge/github.com/samirkape/awesome-go-bot)](https://goreportcard.com/report/github.com/samirkape/awesome-go-bot)

# Awesome-Go Bot

[Join the Telegram Bot](https://t.me/awsmgo_bot)

## Description

Awesome-Go Bot is an interactive Telegram bot version of [Awesome-Go](https://awesome-go.com), a curated list of awesome Go libraries, frameworks, and software. This bot provides a convenient way to explore and discover Go packages, with additional features like filtering by GitHub repository star count and retrieving top Go repositories.

### Features

1. **Browse the Awesome-Go List**: You can browse the extensive Awesome-Go list right within Telegram.

2. **Filter by Star Count**: Easily filter the list based on the number of stars a GitHub repository has.

3. **Top Go Repositories**: Get the top Go repositories ( sorted by stars ) by sending a message with "/top"

4. **Search by Keyword**: You can search for packages using tags.


### Sync

To ensure that the bot is up to date with the latest packages and their star counts, a [Lambda service](https://github.com/samirkape/awesome-go-sync) runs daily to synchronize data from the [Awesome-Go website](https://awesome-go.com).

## Stack

The Awesome-Go Bot is built using the following technologies:

- **Google Cloud Function**: Used to host the bot's backend logic.

- **MongoDB**: Stores data necessary for the bot's operation.

- **Webhook**: Facilitates communication between the bot and Telegram users.

- **Telegram Bot API**: The bot communicates with users via Telegram's API.

- **Go 1.21**: The programming language used for building the bot's backend.
