# my-scraper

[![Build Status](https://travis-ci.org/mono0x/my-scraper.svg)](https://travis-ci.org/mono0x/my-scraper)

`my-scraper` is an atom feed generator for my favorite websites.

Either [Server::Starter](https://metacpan.org/pod/Server::Starter) or [go-server-starter](https://github.com/lestrrat/go-server-starter) is required to run this application.

```sh
# Install glide
brew install glide
# Install go-server-starter
go get github.com/lestrrat/go-server-starter/cmd/start_server
# Install dependencies
glide install
# Build the application
go build
# Start the application
start_server --port=13000 -- ./my-scraper
```

## Supported Websites

- Character Show Fansite <http://charactershow.jp/>
- Facebook
- Fukoku Life <https://act.fukoku-life.co.jp/event/index.php>
- Google Calendar
- Harmonyland <http://www.harmonyland.jp/welcome.html>
- Instagram
- Kittychan Info <http://www.kittychan.info/information.html>
- Memoirs of Shibasaki Saki <http://shibasakisaki.web.fc2.com/>
- PR TIMES (Sanrio) <http://prtimes.jp/main/action.php?run=html&page=searchkey&search_word=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA&search_pattern=1>
- Sanrio News Release <http://www.sanrio.co.jp/corporate/release/>
- Sanrio Puroland <http://www.puroland.jp/>
- Seibuen Event <http://www.seibuen-yuuenchi.jp/event/index.html?category=e1>
- Twitter
- ValuePress! (Sanrio) <https://www.value-press.com/search?q=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA>
