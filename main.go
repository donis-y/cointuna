package main

import (
	"cointuna/dto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Schema []SchemaItem `json:"__schema__"`
	// 이 부분은 dynamic keys를 가지므로 map을 사용할 수 있음
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
	marketUrl := "https://api.upbit.com/v1/market/all"
	body := apiGet(marketUrl)

	// JSON 응답을 파싱하여 MarketCoinInfoResponse 슬라이스로 변환
	var marketsAll []dto.MarketCoinInfoResponse
	err := json.Unmarshal(body, &marketsAll)
	if err != nil {
		fmt.Println("JSON 파싱 중 에러 발생:", err)
		return
	}

	var marketsKRW []dto.MarketCoinInfoResponse
	var marketsBTC []dto.MarketCoinInfoResponse
	var marketsUSDT []dto.MarketCoinInfoResponse

	// 파싱된 데이터 출력
	for _, market := range marketsAll {

		if strings.HasPrefix(market.Market, "KRW") {
			marketsKRW = append(marketsKRW, market)
		}

		if strings.HasPrefix(market.Market, "BTC") {
			marketsBTC = append(marketsBTC, market)
		}
		if strings.HasPrefix(market.Market, "USDT") {
			marketsUSDT = append(marketsUSDT, market)
		}

		fmt.Printf("Market: %s, Korean Name: %s, English Name: %s\n", market.Market, market.KoreanName, market.EnglishName)
	}

	fmt.Printf("끝\n\n\n\n\n\n\n")
	var tickerListKRW []string
	var tickerListBTC []string
	var tickerListUSDT []string
	for _, market := range marketsKRW {
		ticker := strings.Split(market.Market, "-")
		tickerListKRW = append(tickerListKRW, ticker[1])
		fmt.Printf("Market: %s, Korean Name: %s, English Name: %s\n", market.Market, market.KoreanName, market.EnglishName)

	}

	for _, market := range marketsBTC {
		ticker := strings.Split(market.Market, "-")
		tickerListBTC = append(tickerListBTC, ticker[1])
		fmt.Printf("Market: %s, Korean Name: %s, English Name: %s\n", market.Market, market.KoreanName, market.EnglishName)

	}

	for _, market := range marketsUSDT {
		ticker := strings.Split(market.Market, "-")
		tickerListUSDT = append(tickerListUSDT, ticker[1])
		fmt.Printf("Market: %s, Korean Name: %s, English Name: %s\n", market.Market, market.KoreanName, market.EnglishName)
	}

	var coinDetailUrlList []string
	for _, ticker := range marketsKRW {
		marketInfo := "https://api-manager.upbit.com/api/v1/coin_info/pub/%s.json"
		fmt.Printf("티커 %s\n", ticker)
		url := fmt.Sprintf(marketInfo, ticker)
		coinDetailUrlList = append(coinDetailUrlList, url)
		fmt.Println(url)
	}
	fmt.Printf("KRW 총개수 %d, ", len(tickerListKRW))
	fmt.Printf("BTC 총개수 %d, ", len(tickerListBTC))
	fmt.Printf("USDT 총개수 %d\n", len(tickerListUSDT))

	index := 0
	for _, market := range marketsKRW {
		index++
		ticker := strings.Split(market.Market, "-")
		urlForm := "https://api-manager.upbit.com/api/v1/coin_info/pub/%s.json"
		coinDetailUrl := fmt.Sprintf(urlForm, ticker[1])
		body = apiGet(coinDetailUrl)
		//fmt.Printf("가자 %s\n", coinDetailUrl)
		// JSON 응답을 파싱하여 CoinDetailResponse 로 변환
		var coinDetailResponse CoinDetailResponse
		//var coinDetailResponse []dto.MarketCoinInfoResponse
		err := json.Unmarshal(body, &coinDetailResponse)
		if err != nil {
			fmt.Println("JSON 파싱 중 에러 발생:", err)
			return
		}

		if coinDetailResponse.Success {
			numberCap := 0
			if strings.Contains(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "조원") {
				multiply := float64(1000000000000)
				c := strings.Replace(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "조원", "", -1)
				cap, err := strconv.ParseFloat(c, 64)
				if err != nil {
					fmt.Println("ParseFloat 중 에러 발생:", err)
				}
				numberCap = int(cap * multiply)
			} else if strings.Contains(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "억원") {
				multiply := float64(100000000)
				c := strings.Replace(coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, "억원", "", -1)
				cap, err := strconv.ParseFloat(c, 64)
				numberCap = int(cap * multiply)
				if err != nil {
					fmt.Println("ParseFloat 중 에러 발생:", err)
				}
			}

			urlForm = "https://api.upbit.com/v1/candles/days?market=%s&count=1"
			coinPriceUrl := fmt.Sprintf(urlForm, market.Market)
			body = apiGet(coinPriceUrl)
			// JSON 응답을 파싱하여 coinPriceInfoResponse 로 변환
			var coinPriceInfoResponse []dto.CoinPriceInfoResponse
			//var coinDetailResponse []dto.MarketCoinInfoResponse
			err := json.Unmarshal(body, &coinPriceInfoResponse)
			if err != nil {
				fmt.Printf("url %s\n", coinPriceUrl)
				fmt.Println("JSON 파싱 중 에러 발생:", err)
				return
			}

			fmt.Printf("%d|%s|%s|%s|%s|%f\n", numberCap, coinDetailResponse.Data.MarketData.CoinMarketCap.MarketCap, ticker[1], market.KoreanName, market.EnglishName, coinPriceInfoResponse[0].TradePrice)

			if index%10 == 0 { // REST API 요청 수 제한으로 sleep. 초당 10회 (종목, 캔들, 체결, 티커, 호가별 각각 적용)
				time.Sleep(time.Second)
			}
		} else {
			fmt.Printf("실패 Market: %s, Korean Name: %s, English Name: %s\n", market.Market, market.KoreanName, market.EnglishName)
		}

	}

}

func apiGet(url string) []byte {
	// Upbit API에서 시장 정보를 가져옴
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("API 호출 중 에러 발생:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("응답 읽기 중 에러 발생:", err)
		return nil
	}
	return body
}
