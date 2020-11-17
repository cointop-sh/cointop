package types

// OrderType

// OrderType in CoinGecko
type OrderType struct {
	MarketCapDesc string
	MarketCapAsc  string
	GeckoDesc     string
	GeckoAsc      string
	VolumeAsc     string
	VolumeDesc    string
}

// OrderTypeObject for certain order
var OrderTypeObject = &OrderType{
	MarketCapDesc: "market_cap_desc",
	MarketCapAsc:  "market_cap_asc",
	GeckoDesc:     "gecko_desc",
	GeckoAsc:      "gecko_asc",
	VolumeAsc:     "volume_asc",
	VolumeDesc:    "volume_desc",
}

// PriceChangePercentage

// PriceChangePercentage in different amount of time
type PriceChangePercentage struct {
	PCP1h   string
	PCP24h  string
	PCP7d   string
	PCP14d  string
	PCP30d  string
	PCP200d string
	PCP1y   string
}

// PriceChangePercentageObject for different amount of time
var PriceChangePercentageObject = &PriceChangePercentage{
	PCP1h:   "1h",
	PCP24h:  "24h",
	PCP7d:   "7d",
	PCP14d:  "14d",
	PCP30d:  "30d",
	PCP200d: "200d",
	PCP1y:   "1y",
}

// SHARED
// coinBaseStruct [private]
type coinBaseStruct struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// AllCurrencies map all currencies (USD, BTC) to float64
type AllCurrencies map[string]float64

// LocalizationItem map all locale (en, zh) into respective string
type LocalizationItem map[string]string

// TYPES

// DescriptionItem map all description (in locale) into respective string
type DescriptionItem map[string]string

// LinksItem map all links
type LinksItem map[string]interface{}

// ChartItem ...
type ChartItem [2]float32

// MarketDataItem map all market data item
type MarketDataItem struct {
	CurrentPrice                           AllCurrencies     `json:"current_price"`
	ROI                                    *ROIItem          `json:"roi"`
	ATH                                    AllCurrencies     `json:"ath"`
	ATHChangePercentage                    AllCurrencies     `json:"ath_change_percentage"`
	ATHDate                                map[string]string `json:"ath_date"`
	MarketCap                              AllCurrencies     `json:"market_cap"`
	MarketCapRank                          uint16            `json:"market_cap_rank"`
	TotalVolume                            AllCurrencies     `json:"total_volume"`
	High24                                 AllCurrencies     `json:"high_24h"`
	Low24                                  AllCurrencies     `json:"low_24h"`
	PriceChange24h                         float64           `json:"price_change_24h"`
	PriceChangePercentage24h               float64           `json:"price_change_percentage_24h"`
	PriceChangePercentage7d                float64           `json:"price_change_percentage_7d"`
	PriceChangePercentage14d               float64           `json:"price_change_percentage_14d"`
	PriceChangePercentage30d               float64           `json:"price_change_percentage_30d"`
	PriceChangePercentage60d               float64           `json:"price_change_percentage_60d"`
	PriceChangePercentage200d              float64           `json:"price_change_percentage_200d"`
	PriceChangePercentage1y                float64           `json:"price_change_percentage_1y"`
	MarketCapChange24h                     float64           `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h           float64           `json:"market_cap_change_percentage_24h"`
	PriceChange24hInCurrency               AllCurrencies     `json:"price_change_24h_in_currency"`
	PriceChangePercentage1hInCurrency      AllCurrencies     `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24hInCurrency     AllCurrencies     `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7dInCurrency      AllCurrencies     `json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14dInCurrency     AllCurrencies     `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30dInCurrency     AllCurrencies     `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage60dInCurrency     AllCurrencies     `json:"price_change_percentage_60d_in_currency"`
	PriceChangePercentage200dInCurrency    AllCurrencies     `json:"price_change_percentage_200d_in_currency"`
	PriceChangePercentage1yInCurrency      AllCurrencies     `json:"price_change_percentage_1y_in_currency"`
	MarketCapChange24hInCurrency           AllCurrencies     `json:"market_cap_change_24h_in_currency"`
	MarketCapChangePercentage24hInCurrency AllCurrencies     `json:"market_cap_change_percentage_24h_in_currency"`
	TotalSupply                            *float64          `json:"total_supply"`
	CirculatingSupply                      float64           `json:"circulating_supply"`
	Sparkline                              *SparklineItem    `json:"sparkline_7d"`
	LastUpdated                            string            `json:"last_updated"`
}

// CommunityDataItem map all community data item
type CommunityDataItem struct {
	FacebookLikes            *uint        `json:"facebook_likes"`
	TwitterFollowers         *uint        `json:"twitter_followers"`
	RedditAveragePosts48h    *float64     `json:"reddit_average_posts_48h"`
	RedditAverageComments48h *float64     `json:"reddit_average_comments_48h"`
	RedditSubscribers        *uint        `json:"reddit_subscribers"`
	RedditAccountsActive48h  *interface{} `json:"reddit_accounts_active_48h"`
	TelegramChannelUserCount *uint        `json:"telegram_channel_user_count"`
}

// DeveloperDataItem map all developer data item
type DeveloperDataItem struct {
	Forks              *uint `json:"forks"`
	Stars              *uint `json:"stars"`
	Subscribers        *uint `json:"subscribers"`
	TotalIssues        *uint `json:"total_issues"`
	ClosedIssues       *uint `json:"closed_issues"`
	PRMerged           *uint `json:"pull_requests_merged"`
	PRContributors     *uint `json:"pull_request_contributors"`
	CommitsCount4Weeks *uint `json:"commit_count_4_weeks"`
}

// PublicInterestItem map all public interest item
type PublicInterestItem struct {
	AlexaRank   uint `json:"alexa_rank"`
	BingMatches uint `json:"bing_matches"`
}

// ImageItem struct for all sizes of image
type ImageItem struct {
	Thumb string `json:"thumb"`
	Small string `json:"small"`
	Large string `json:"large"`
}

// ROIItem ROI Item
type ROIItem struct {
	Times      float64 `json:"times"`
	Currency   string  `json:"currency"`
	Percentage float64 `json:"percentage"`
}

// SparklineItem for sparkline
type SparklineItem struct {
	Price []float64 `json:"price"`
}

// TickerItem for ticker
type TickerItem struct {
	Base   string `json:"base"`
	Target string `json:"target"`
	Market struct {
		Name             string `json:"name"`
		Identifier       string `json:"identifier"`
		TradingIncentive bool   `json:"has_trading_incentive"`
	} `json:"market"`
	Last            float64            `json:"last"`
	ConvertedLast   map[string]float64 `json:"converted_last"`
	Volume          float64            `json:"volume"`
	ConvertedVolume map[string]float64 `json:"converted_volume"`
	Timestamp       string             `json:"timestamp"`
	IsAnomaly       bool               `json:"is_anomaly"`
	IsStale         bool               `json:"is_stale"`
	CoinID          string             `json:"coin_id"`
}

// StatusUpdateItem for BEAM
type StatusUpdateItem struct {
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	User        string `json:"user"`
	UserTitle   string `json:"user_title"`
	Pin         bool   `json:"pin"`
	Project     struct {
		coinBaseStruct
		Type  string    `json:"type"`
		Image ImageItem `json:"image"`
	} `json:"project"`
}

// CoinsListItem item in CoinList
type CoinsListItem struct {
	coinBaseStruct
}

// CoinsMarketItem item in CoinMarket
type CoinsMarketItem struct {
	coinBaseStruct
	Image                               string         `json:"image"`
	CurrentPrice                        float64        `json:"current_price"`
	MarketCap                           float64        `json:"market_cap"`
	MarketCapRank                       int16          `json:"market_cap_rank"`
	TotalVolume                         float64        `json:"total_volume"`
	High24                              float64        `json:"high_24h"`
	Low24                               float64        `json:"low_24h"`
	PriceChange24h                      float64        `json:"price_change_24h"`
	PriceChangePercentage24h            float64        `json:"price_change_percentage_24h"`
	MarketCapChange24h                  float64        `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h        float64        `json:"market_cap_change_percentage_24h"`
	CirculatingSupply                   float64        `json:"circulating_supply"`
	TotalSupply                         float64        `json:"total_supply"`
	ATH                                 float64        `json:"ath"`
	ATHChangePercentage                 float64        `json:"ath_change_percentage"`
	ATHDate                             string         `json:"ath_date"`
	ROI                                 *ROIItem       `json:"roi"`
	LastUpdated                         string         `json:"last_updated"`
	SparklineIn7d                       *SparklineItem `json:"sparkline_in_7d"`
	PriceChangePercentage1hInCurrency   *float64       `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24hInCurrency  *float64       `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7dInCurrency   *float64       `json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14dInCurrency  *float64       `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30dInCurrency  *float64       `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage200dInCurrency *float64       `json:"price_change_percentage_200d_in_currency"`
	PriceChangePercentage1yInCurrency   *float64       `json:"price_change_percentage_1y_in_currency"`
}

// EventCountryItem item in EventsCountries
type EventCountryItem struct {
	Country string `json:"country"`
	Code    string `json:"code"`
}

// ExchangeRatesItem item in ExchangeRate
type ExchangeRatesItem map[string]ExchangeRatesItemStruct

// ExchangeRatesItemStruct struct in ExchangeRateItem
type ExchangeRatesItemStruct struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// Global for data of /global
type Global struct {
	ActiveCryptocurrencies          uint16        `json:"active_cryptocurrencies"`
	UpcomingICOs                    uint16        `json:"upcoming_icos"`
	EndedICOs                       uint16        `json:"ended_icos"`
	Markets                         uint16        `json:"markets"`
	MarketCapChangePercentage24hUSD float32       `json:"market_cap_change_percentage_24h_usd"`
	TotalMarketCap                  AllCurrencies `json:"total_market_cap"`
	TotalVolume                     AllCurrencies `json:"total_volume"`
	MarketCapPercentage             AllCurrencies `json:"market_cap_percentage"`
	UpdatedAt                       int64         `json:"updated_at"`
}
