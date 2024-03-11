package dto

//type CoinDetailResponse struct {
//	Success bool `json:"success"`
//	Data    Data `json:"data"`
//}
//
//type Data struct {
//	EnglishName        string              `json:"english_name"`
//	KoreanName         string              `json:"korean_name"`
//	Symbol             string              `json:"symbol"`
//	HeaderKeyValues    HeaderKeyValues     `json:"header_key_values"`
//	MarketData         MarketData          `json:"market_data"`
//	MainComponents     []MainComponent     `json:"main_components"`
//	BlockInspectorURLs []BlockInspectorURL `json:"block_inspector_urls"`
//}
//
//type HeaderKeyValues struct {
//	Schema []SchemaItem `json:"__schema__"`
//	// 이 부분은 dynamic keys를 가지므로 map을 사용할 수 있음
//	DynamicKeys map[string]KeyValue `json:"-"`
//}
//
//type SchemaItem struct {
//	Key         string `json:"key"`
//	Name        string `json:"name"`
//	Description string `json:"description"`
//	IsStatic    bool   `json:"is_static"`
//}
//
//type KeyValue struct {
//	Value string `json:"value"`
//	Link  string `json:"link"`
//}
//
//type MarketData struct {
//	CoinMarketCap     MarketInfo  `json:"coin_market_cap"`
//	CoinGecko         MarketInfo  `json:"coin_gecko"`
//	ProjectTeam       ProjectTeam `json:"project_team"`
//	MaxSupply         string      `json:"max_supply"`
//	MaxSupplyProvider string      `json:"max_supply_provider"`
//}
//
//type MarketInfo struct {
//	MarketCap         string `json:"market_cap"`
//	CirculatingSupply string `json:"circulating_supply"`
//	Date              string `json:"date"`
//	Market            string
//}
//
//type ProjectTeam struct {
//	MarketCap         string   `json:"market_cap"`
//	CirculatingSupply string   `json:"circulating_supply"`
//	Date              string   `json:"date"`
//	SupplyPlan        KeyValue `json:"supply_plan"`
//	Contact           string   `json:"contact"`
//}
//
//type MainComponent struct {
//	TypeName string `json:"type_name"`
//	Detail   Detail `json:"detail"`
//}
//
//type Detail struct {
//	Subtitle string `json:"subtitle"`
//	Content  string `json:"content"`
//}
//
//type BlockInspectorURL struct {
//	Value string `json:"value"`
//	Link  string `json:"link"`
//}
