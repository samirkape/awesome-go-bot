## awesome-go-bot

Link: https://t.me/awsmgo_bot

### Description 

Interactive bot version of https://awesome-go.com.  

Additionally, you can 

* Filter the list by Github repo star count.  
* Get top Go repositories by simply sending *Top N* where N is any number in betwen 0 to maximum number of packages (around 2100 as of now).

### Sync

To make sure we are up to date with the new packages and star counts, [another lambda service](https://github.com/samirkape/awesome-go-sync) \
is running once every day to synchronise https://awesome-go.com and star counts.

### Stack
* Google Cloud Function
* MongoDB
* Webhook
* Telegram Bot API
* Go1.13

### TODO

- [x] Search by tags
- [x] Add /meta command to fetch metadata such as,
     * Number of packages. 
     * Last updated time.
- [x] Decouple backend from frontend.
- [ ] Add LRU cache to reduce Google cloud function's cold boot time.
- [x] Add inline mode to,
     * Group multiple message.
     * Fetch dev.to Golang articles corresponding to package category/tag.

### Known Bugs

1. **Category:** Style guide\
**Kind:** Filter by star count\
**Issue:** Github api get request is failing due to wrong URL path. This is happening\
because every package in the style guide category is listed as a path to .md file instead of\
Github repository.

2. Some sub categories in the awesome-go.com are\
being recorded as a separate categories due to inconsistent\
fomatting.
