## awesome-go-bot

Link: https://t.me/go_pkg_sender_bot

Interactive bot version of https://awesome-go.com.

Additionally, It has a feature by which you can filter Go repositories by their star counts, 

this way, you can get a list of top Go repositories by simply sending *Top N* to telegram bot where N is a number. 

And to make sure we are up to date with new packages and star counts, another lamda service is running to synchronise 

https://awesome-go.com and star counts.

---

TODOs

- [ ] Decouple backend from frontend.
- [ ] Add LRU cache to reduce function's cold boot time.
- [ ] Add inline mode to,
     - [ ] Group multiple message.
     - [ ] Fetch dev.to Golang articles corresponding to category.

 

