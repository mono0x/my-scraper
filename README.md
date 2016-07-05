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
GO15VENDOREXPERIMENT=1 go build
# Start the application
start_server --port=13000 -- ./my-scraper
```

## Supported Websites

- Character Show Fansite <http://charactershow.jp/>
- Fukkachan Calendar <http://www.fukkachan.com/>
- Gotouchi Chara Calendar <http://gotouchi-chara.jp/calendar_.html>
- Kittychan Info <http://www.kittychan.info/information.html>
- Memoirs of Shibasaki Saki <http://shibasakisaki.web.fc2.com/>
- Mucchan Musao <https://www.facebook.com/mucchan.musao>
- Olympus Camera <https://www.facebook.com/FotoPus>
- Sanrio Events Calendar <http://ameblo.jp/ohtaket/entry-12059393801.html>
- Sanrio Puroland <http://www.puroland.jp/>
- Seibuen Event <http://www.seibuen-yuuenchi.jp/event/index.html?category=e1>
