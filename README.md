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
- Fukkachan Calendar <http://www.fukkachan.com/>
- Fukoku Life <https://act.fukoku-life.co.jp/event/index.php>
- Gotouchi Chara Calendar <http://gotouchi-chara.jp/calendar_.html>
- Kittychan Info <http://www.kittychan.info/information.html>
- Life Corporation <https://www.facebook.com/lifecorp428>
- Memoirs of Shibasaki Saki <http://shibasakisaki.web.fc2.com/>
- Mucchan Musao <https://www.facebook.com/mucchan.musao>
- Olympus Camera <https://www.facebook.com/FotoPus>
- PR TIMES (Sanrio) <http://prtimes.jp/main/action.php?run=html&page=searchkey&search_word=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA&search_pattern=1>
- Sanrio Events Calendar <http://ameblo.jp/ohtaket/entry-12059393801.html>
- Sanrio News Release <http://www.sanrio.co.jp/corporate/release/>
- Sanrio Puroland <http://www.puroland.jp/>
- Seibuen Event <http://www.seibuen-yuuenchi.jp/event/index.html?category=e1>
- ValuePress! (Sanrio) <https://www.value-press.com/search?q=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA>
- Yufu Terashima Calendar <http://sp.yufuterashima.com/>
