# my-scraper

[![Build Status](https://travis-ci.org/mono0x/my-scraper.svg)](https://travis-ci.org/mono0x/my-scraper)
[![Coverage Status](https://coveralls.io/repos/github/mono0x/my-scraper/badge.svg?branch=master)](https://coveralls.io/github/mono0x/my-scraper?branch=master)

`my-scraper` is an atom feed generator for my favorite websites.

```sh
# Install retool
make setup
# Test & build the app
make
# Start the app
retool do start_server --port=8080 -- ./my-scraper
```

## Supported Websites

- Facebook
- Fukoku Life <https://act.fukoku-life.co.jp/event/index.php>
- Google Calendar
- Harmonyland <http://www.harmonyland.jp/welcome.html>
- Instagram
- Kittychan Info <http://www.kittychan.info/information.html>
- PR TIMES (Sanrio) <http://prtimes.jp/main/action.php?run=html&page=searchkey&search_word=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA&search_pattern=1>
- Sanrio News Release <http://www.sanrio.co.jp/corporate/release/>
- Sanrio Puroland <http://www.puroland.jp/>
- Seibuen Event <http://www.seibu-leisure.co.jp/event/index.html?category=e1>
- Twitter
- ValuePress! (Sanrio) <https://www.value-press.com/search?q=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA>
- Yuyakekoyake News <http://yuyakekoyake.jp/news/index.php>
