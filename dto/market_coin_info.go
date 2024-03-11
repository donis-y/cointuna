package dto

// MarketCoinInfoResponse 시장 정보를 나타내는 구조체
type MarketCoinInfoResponse struct {
	Market      string `json:"market"`
	KoreanName  string `json:"korean_name"`
	EnglishName string `json:"english_name"`
}
