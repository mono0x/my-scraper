package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/feeds"
)

var sites = []string{
	"aeon.jp",
	"alpark.net",
	"ameblo.jp/emikofruit/",
	"ameblo.jp/kikilala2011/",
	"ameblo.jp/mitsubishimotorsshowroom/",
	"animelab.com",
	"apita-niigatakameda.jp",
	"ario-ageo.jp",
	"ario-sapporo.jp",
	"atgtickets.com",
	"avex.jp",
	"bcl-brand.jp",
	"bhlivetickets.co.uk",
	"blog.itoyokado.co.jp",
	"boatrace.jp",
	"brentcross.co.uk",
	"cafedemiki.jp",
	"centrair.jp",
	"chara-hobby.com",
	"city.ibaraki-koga.lg.jp",
	"city.tama.lg.jp",
	"classic32fes.jp",
	"cutecube-harajuku.jp",
	"dempagumi.dearstage.com",
	"dmdepart.jp",
	"doggypark.jp",
	"drusillas.co.uk",
	"e-a-site.com",
	"e-oubo.com",
	"ekimise.jp",
	"emifull.jp",
	"empmuseum.org",
	"eplus.jp",
	"eterna-takasaki.com",
	"eventbrite.com",
	"eventimapollo.com",
	"exittunesacademy.jp",
	"expo2015.jp",
	"facebook.com/DJ.Hello.Kitty.sanrio.license/",
	"facebook.com/azul.official/",
	"facebook.com/drusillaspark/",
	"facebook.com/evaairwayscorp.jp/",
	"facebook.com/expo2015.jp/",
	"facebook.com/hellokittyrunjapan/",
	"facebook.com/sanriohongkong/",
	"fukutsu-aeonmall.com",
	"gatecity.jp",
	"gilariverarena.com",
	"girls-award.com",
	"gotochikitty.com",
	"goyah.net",
	"grand-tree.jp",
	"hellokittygoaround.com.sg",
	"hellokittylive.com",
	"hiroshimafuchu-aeonmall.com",
	"i-rekihaku.jp",
	"ichinomiya.nha.or.jp",
	"inazumarock.com",
	"japan-monkeypark.jp",
	"jra.go.jp",
	"jra.jp",
	"kasugai.nha.or.jp",
	"kawasaki-roufuku.net",
	"keio-ekiden.com",
	"keio.co.jp",
	"khb-tv.co.jp",
	"kikilala-dreamfunding.com",
	"kiramune.jp",
	"l-tike.com",
	"lalagarden-kawaguchi.com",
	"livenation.co.uk",
	"loft.co.jp",
	"makuharishintoshin-aeonmall.com",
	"mallage.com",
	"mikihouseland.jp",
	"narita-airport.jp",
	"nasuhai.co.jp",
	"nasushiobara-kanko.jp",
	"natori-aeonmall.com",
	"nnr.co.jp",
	"omochaoukoku.co.jp",
	"orleansarena.com",
	"portsmouthguildhall.org.uk",
	"post.japanpost.jp",
	"pref.yamanashi.jp",
	"pride-fish.jp",
	"puroland.jp",
	"pyramidparrhall.com",
	"royalandderngate.co.uk",
	"sacas.net",
	"sanrio.cocolog-nifty.com",
	"sanrio.com",
	"seibuen-yuuenchi.jp",
	"sogo-seibu.jp",
	"southendtheatres.org.uk",
	"stdavidshallcardiff.co.uk",
	"t-expo.jp",
	"takashimaya.co.jp",
	"tama-center.net",
	"tamadairanomori-aeonmall.com",
	"tamajack.com",
	"tamashii.jp",
	"ticket.pia.jp",
	"ticketmaster.co.uk",
	"ticketmaster.com",
	"tickets.yorkbarbican.co.uk",
	"tohoku-bank.co.jp",
	"tokyo-joypolis.com",
	"tokyo-skytree.jp",
	"tonxton.com",
	"toyota.nha.or.jp",
	"trch.co.uk",
	"tsu-kyotei.com",
	"tsukuba-aeonmall.com",
	"twitter.com/CTS_staff/",
	"twitter.com/Hello_Kitty_UK/",
	"twitter.com/KittyNagoya/",
	"twitter.com/atlrs_official/",
	"twitter.com/cafe_de_miki_d/",
	"twitter.com/drusillaspark/",
	"twitter.com/karaage_mayu/",
	"twitter.com/kirimi_sanrio/",
	"twitter.com/loft_shibuya/",
	"twitter.com/mm_wapi/",
	"twitter.com/sanrio_news/",
	"uny.co.jp",
	"yamanashi-kankou.jp",
	"yokohama.lalaport.jp",
}

var keywords = []string{
	//	"ハローキティ",
	//	"マイメロディ",
	//	"シナモン",
	//	"ポムポムプリン",
	//	"クロミ",
	//	"リトルツインスターズ",
	//	"リルリルフェアリル",
	//	"ジュエルペット",
	//	"ABシ～ナモン",
	//	"DJハローキティ",
	//	"Doraemon meets Hello Kitty",
	//	"HAPPY GIFT CHRISTMAS",
	//	"HELLO KITTY RUN",
	//	"サンリオハロウィンパレード",
	//	"SANRIO EXPO",
	//	"Sanrio Days",
	//	"お花を咲かせよう",
	//	"ぐでーてぃんぐ",
	//	"ハッピーフレンドクリスマス",
	//	"キティズパラダイス",
	//	"チャレンジタイム",
	//	"ハッピーキャラバン",
	"ハローキティが遊びにくるよ",
	"ハローキティのラブリーステージ",
	"ハローキティのワンダーランド 科学の国へようこそ",
	"ハローキティアート展",
	"ハートフルタイム",
	"ミラクルダンスタイム",
	"ラブリーステージ",
}

const (
	BingSearchApiEndpoint = "https://api.cognitive.microsoft.com/bing/v5.0/search"
)

type SanrioSearchSource struct {
}

type BingSearchWebPage struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	DisplayUrl string `json:"displayUrl"`
	Snippet    string `json:"snippet"`
	DeepLinks  []*struct {
	} `json:"deepLinks"`
}

type BingSearchResponse struct {
	Type     string `json:"_type"`
	WebPages struct {
		WebSearchUrl          string               `json:"webSearchUrl"`
		TotalEstimatedMatches uint                 `json:"totalEstimatedMatches"`
		Value                 []*BingSearchWebPage `json:"value"`
	} `json:"webPages"`
	DateLastCrawled string `json:"dateLastCrawled"`
}

type BingErrorResponse struct {
	Type   string `json:"_type"`
	Errors []*struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func NewSanrioSearchSource() *SanrioSearchSource {
	return &SanrioSearchSource{}
}

func (s *SanrioSearchSource) Scrape() (*feeds.Feed, error) {
	webPages, err := s.Fetch()
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromWebPages(webPages)
}

func (s *SanrioSearchSource) Fetch() ([]*BingSearchWebPage, error) {
	webPages := make([]*BingSearchWebPage, 0, len(sites)/1*50)

	siteConditions := make([]string, 0, 1)
	for _, site := range sites {
		siteConditions = append(siteConditions, fmt.Sprintf(`site:"%s"`, site))
		if len(siteConditions) >= 1 {
			response, err := processChunk(siteConditions)
			if err != nil {
				return nil, err
			}
			webPages = append(webPages, response.WebPages.Value...)
			siteConditions = nil
		}
	}
	if len(siteConditions) > 0 {
		response, err := processChunk(siteConditions)
		if err != nil {
			return nil, err
		}
		webPages = append(webPages, response.WebPages.Value...)
	}
	return webPages, nil
}

func processChunk(siteConditions []string) (*BingSearchResponse, error) {
	keywordConditions := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		keywordConditions = append(keywordConditions, fmt.Sprintf(`(%s)`, keyword))
	}

	values := url.Values{}
	values.Add("q", fmt.Sprintf("(%s) (%s)", strings.Join(keywordConditions, " OR "), strings.Join(siteConditions, " OR ")))
	values.Add("count", "50")
	values.Add("cc", "JP")

	fmt.Println(values)
	return doSearch(values)
}

func doSearch(values url.Values) (*BingSearchResponse, error) {
	req, err := http.NewRequest("GET", BingSearchApiEndpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_SUBSCRIPTION_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))

	var searchResponse BingSearchResponse
	if err := json.Unmarshal(data, &searchResponse); err != nil {
		return nil, err
	}
	return &searchResponse, nil
}

func (s *SanrioSearchSource) ScrapeFromWebPages(webPages []*BingSearchWebPage) (*feeds.Feed, error) {
	items := make([]*feeds.Item, 0, len(webPages))
	for _, page := range webPages {
		link := canonicalizeURL(page.Url)

		items = append(items, &feeds.Item{
			Id:          link,
			Link:        &feeds.Link{Href: link},
			Title:       page.Name,
			Description: page.Snippet,
		})
	}

	feed := &feeds.Feed{
		Title: "Sanrio Search",
		Link:  &feeds.Link{Href: "https://scraper.mono0x.net/sanrio-search"},
		Items: items,
	}
	return feed, nil
}

func canonicalizeURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return u
	}
	link := parsed.Query().Get("r")
	if link == "" {
		return u
	}
	return link
}
