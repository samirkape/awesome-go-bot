## awesome-go-bot

Link: https://t.me/go_pkg_sender_bot

### Description 

Interactive bot version of https://awesome-go.com.

Additionally, It has a feature by which you can filter Go repositories by their star 

counts. This way, you can get a list of top Go repositories by simply sending *Top N* to telegram

bot where N is any number in betwen 0 to maximum number of packages (around 2100 as of now). 

### Updates

To make sure we are up to date with the new packages and star counts, [another lamda service](https://github.com/samirkape/awesome-go-sync) is running 

at the interval of 15 days to synchronise https://awesome-go.com and star counts.


### TODO

- [ ] Decouple backend from frontend.
- [ ] Add LRU cache to reduce Google cloud function's cold boot time.
- [ ] Add inline mode to,
     - [ ] Group multiple message.
     - [ ] Fetch dev.to Golang articles corresponding to package category/tag.

 

