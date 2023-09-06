# Awesome-Go Bot

Link: https://t.me/awsmgo_bot

## Description

Awesome-Go Bot is an interactive Telegram bot version of [Awesome-Go](https://awesome-go.com), a curated list of awesome Go libraries, frameworks, and software. This bot provides a convenient way to explore and discover Go packages, with additional features like filtering by GitHub repository star count and retrieving top Go repositories.

### Features

1. **Browse the Awesome-Go List**: You can browse the extensive Awesome-Go list right within Telegram.

2. **Filter by Star Count**: Easily filter the list based on the number of stars a GitHub repository has.

3. **Top Go Repositories**: Get the top Go repositories by sending a message with "Top N," where N is any number between 0 and the maximum number of packages (approximately 2100 as of now).

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

## TODO

Here are some of the planned enhancements and features for the Awesome-Go Bot:

- [x] Search by Tags: You can search for packages using tags.

- [x] Meta Command: Fetch metadata, including the number of packages and the last updated time.

- [ ] LRU Cache: Implement an LRU cache to reduce Google Cloud Function's cold boot time.

- [x] Inline Mode: Group multiple messages and fetch Golang articles from dev.to corresponding to package category/tag.

[Join the Telegram Bot](https://t.me/awsmgo_bot)
