package main

import (
	"cointuna/dto"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type CoinDetailResponse struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	EnglishName        string              `json:"english_name"`
	KoreanName         string              `json:"korean_name"`
	Symbol             string              `json:"symbol"`
	HeaderKeyValues    HeaderKeyValues     `json:"header_key_values"`
	MarketData         MarketData          `json:"market_data"`
	MainComponents     []MainComponent     `json:"main_components"`
	BlockInspectorURLs []BlockInspectorURL `json:"block_inspector_urls"`
}

type HeaderKeyValues struct {
	Schema      []SchemaItem          `json:"__schema__"`
	DynamicKeys map[string]KeyValue `json:"-"`
}

type SchemaItem struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsStatic    bool   `json:"is_static"`
}

type KeyValue struct {
	Value string `json:"value"`
	Link  string `json:"link"`
}

type MarketData struct {
	CoinMarketCap     MarketInfo  `json:"coin_market_cap"`
	CoinGecko         MarketInfo  `json:"coin_gecko"`
	ProjectTeam       ProjectTeam `json:"project_team"`
	MaxSupply         string      `json:"max_supply"`
	MaxSupplyProvider string      `json:"max_supply_provider"`
}

type MarketInfo struct {
	MarketCap         string `json:"market_cap"`
	CirculatingSupply string `json:"circulating_supply"`
	Date              string `json:"date"`
	Market            string
}

type ProjectTeam struct {
	MarketCap         string   `json:"market_cap"`
	CirculatingSupply string   `json:"circulating_supply"`
	Date              string   `json:"date"`
	SupplyPlan        KeyValue `json:"supply_plan"`
	Contact           string   `json:"contact"`
}

type MainComponent struct {
	TypeName string `json:"type_name"`
	Detail   Detail `json:"detail"`
}

type Detail struct {
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

type BlockInspectorURL struct {
	Value string `json:"value"`
	Link  string `json:"link"`
}

func main() {
	// .env 파일 로드
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 환경 변수에서 데이터베이스 연결 정보 읽기
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// DSN (Data Source Name) 생성
	// 사용자 이름이나 비밀번호에 특수문자가 포함될 경우를 대비해 URL 인코딩을 수행합니다.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		url.QueryEscape(dbUser),
		url.QueryEscape(dbPassword),
		dbHost,
		dbPort,
		dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to MariaDB!")

	marketUrl := "https://api.upbit.com/v1/market/all"
	body := apiGet(marketUrl)

	var marketsAll []dto.MarketCoinInfoResponse
	err = json.Unmarshal(body, &marketsAll)
	if err != nil {
		log.Println("JSON 파싱 중 에러 발생:", err)
		return
	}

	var marketsKRW []dto.MarketCoinInfoResponse
	for _, market := range marketsAll {
		if strings.HasPrefix(market.Market, "KRW") {
			marketsKRW = append(marketsKRW, market)
		}
	}

	index := 0
	for _, market := range marketsKRW {
		index++
		ticker := strings.Split(market.Market, "-")
		urlForm := "https://api-manager.upbit.com/api/v1/coin_info/pub/%s.json"
		coinDetailUrl := fmt.Sprintf(urlForm, ticker[1])
		body = apiGet(coinDetailUrl)

		var coinDetailResponse CoinDetailResponse
		err := json.Unmarshal(body, &coinDetailResponse)
		if err != nil {
			log.Println("JSON 파싱 중 에러 발생:", err)
			return
		}

		if coinDetailResponse.Success {
			numberCap := 0
			if strings.Contains(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "조원") {
				multiply := float64(1000000000000)
				c := strings.Replace(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "조원", "", -1)
				cap, err := strconv.ParseFloat(c, 64)
				if err == nil {
					numberCap = int(cap * multiply)
				}
			} else if strings.Contains(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "억원") {
				multiply := float64(100000000)
				c := strings.Replace(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "억원", "", -1)
				cap, err := strconv.ParseFloat(c, 64)
				if err == nil {
					numberCap = int(cap * multiply)
				}
			}

			urlForm = "https://api.upbit.com/v1/candles/days?market=%s&count=1"
			coinPriceUrl := fmt.Sprintf(urlForm, market.Market)
			body = apiGet(coinPriceUrl)

			var coinPriceInfoResponse []dto.CoinPriceInfoResponse
			err := json.Unmarshal(body, &coinPriceInfoResponse)
			if err != nil {
				log.Printf("url %s", coinPriceUrl)
				log.Println("JSON 파싱 중 에러 발생:", err)
				return
			}

			// 데이터베이스에 저장
			stmt, err := db.Prepare("INSERT INTO coin_info(market_cap_num, market_cap_str, ticker, korean_name, english_name, trade_price) VALUES(?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Println(err)
				continue
			}
			defer stmt.Close()

			_, err = stmt.Exec(numberCap, coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, ticker[1], market.KoreanName, market.EnglishName, coinPriceInfoResponse[0].TradePrice)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Saved: (%d/%d) %s (%s)", index, len(marketsKRW), market.KoreanName, market.Market)

			if index%10 == 0 {
				time.Sleep(time.Second)
			}
		} else {
			log.Printf("실패 Market: %s, Korean Name: %s, English Name: %s", market.Market, market.KoreanName, market.EnglishName)
		}
	}
}

func apiGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("API 호출 중 에러 발생:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("응답 읽기 중 에러 발생:", err)
		return nil
	}
	return body
}