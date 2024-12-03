package blumapi

import (
	"math/big"
	"sort"
	"strconv"
)

type TokenListRsp []Token

type Token struct {
	ID               int    `json:"id"`
	Address          string `json:"address"`
	Ticker           string `json:"ticker"`
	Name             string `json:"name"`
	Shortname        string `json:"shortname"`
	Description      string `json:"description"`
	BannerFileKey    string `json:"bannerFileKey"`
	IconFileKey      string `json:"iconFileKey"`
	IsNSFW           bool   `json:"isNSFW"`
	Status           string `json:"status"`
	Socials          []any  `json:"socials"`
	ReleaseTimestamp int64  `json:"releaseTimestamp"`
	Stats            Stats  `json:"stats"`
}

type Stats struct {
	MarketCap         string `json:"marketCap"`
	Volume            string `json:"volume"`
	TransactionsCount int    `json:"transactionsCount"`
	HoldersCount      int    `json:"holdersCount"`
	TonCollected      int64  `json:"tonCollected"`
}

type ChartRsp struct {
	Points [][]any `json:"points"`
}

type Points []Point
type Point struct {
	Ts    int
	Price *big.Float
}

type LiveDataRsp struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Address      string `json:"address"`
	User         User   `json:"user"`
	Ticker       string `json:"ticker"`
	Shortname    string `json:"shortname"`
	IconFileKey  string `json:"iconFileKey"`
	Type         string `json:"type"`
	Amount       string `json:"amount"`
	JettonAmount string `json:"jettonAmount"`
	Timestamp    int    `json:"timestamp"`
}

func (t *Transaction) AmountInt() int64 {
	res, _ := strconv.ParseFloat(t.Amount, 64)

	return int64(res)
}

func (t *Transaction) JettonAmountInt() int64 {
	res, _ := strconv.ParseFloat(t.JettonAmount, 64)

	return int64(res)
}

type User struct {
	Address string `json:"address"`
}

func (c *ChartRsp) ToPoints() Points {
	points := make(Points, len(c.Points))

	for i, p := range c.Points {
		priceString := p[1].(string)
		priceFloat := new(big.Float)
		_, _ = priceFloat.SetString(priceString)

		points[i] = Point{
			Ts:    int(p[0].(float64)),
			Price: priceFloat,
		}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Ts < points[j].Ts
	})

	return points
}
