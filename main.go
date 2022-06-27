package main

import (
	"fmt"
	"math"
	"sort"
)

type Calculator struct {
	BuyOnly                bool
	SupportFractionalShare bool
	SlotSize               float64
	buyingPriceSelector    PriceSelector
	sellingPriceSelector   PriceSelector
	valuePriceSelector     PriceSelector
}

// CalculateOrders calculates orders for a model, a portfolio and a value limit
// Default support fractional share, best bid limit
func (c Calculator) CalculateOrders(market Market, model Model, p Portfolio, limit float64) []Order {
	p = p.AggregateEquivalents(market, model, c.valuePriceSelector).
		Select(model.Symbols())

	estimatedTotalOrderValue := math.Min(limit, p.CashAmount())
	preTav := c.calculateTotalAssetValue(market, p)
	postTav := preTav + estimatedTotalOrderValue

	positionDeltaMap := make(map[Symbol]float64)
	for _, s := range model.Symbols().ToSlice() {
		positionDeltaMap[s] = (model.TargetProportion(s)*postTav - c.calculatePositionValue(market, s, p.Position(s))) /
			c.valuePriceSelector(market, s)
	}

	candidateOrders := make([]Order, 0)
	for s, posDelta := range positionDeltaMap {
		var amount float64

		if c.SupportFractionalShare {
			amount = math.Abs(posDelta)
		} else {
			amount = math.Abs(math.Floor(posDelta*c.SlotSize) / c.SlotSize)
		}

		if amount == 0 {
			continue
		}

		var side OrderSide
		if posDelta > 0 {
			side = OrderSideBuy
		} else {
			side = OrderSideSell
		}

		candidateOrders = append(candidateOrders, Order{
			Symbol:     s,
			Type:       OrderTypeLimit,
			LimitPrice: c.buyingPriceSelector(market, s),
			Amount:     amount,
			Side:       side,
		})
	}

	sort.Slice(candidateOrders, func(i, j int) bool {
		return candidateOrders[i].Value() < candidateOrders[j].Value()
	})

	return candidateOrders
}

func (c Calculator) calculateTotalAssetValue(market Market, p Portfolio) (tav float64) {
	for _, s := range p.Symbols().ToSlice() {
		tav += c.calculatePositionValue(market, s, p.Position(s))
	}

	return
}

func (c Calculator) calculatePositionValue(market Market, s Symbol, position float64) (v float64) {
	return c.valuePriceSelector(market, s) * position
}

var market Market = staticMarket{
	"TLT": {
		Last: 112.56,
		Bid:  112.56,
		Ask:  112.56,
	},
	"QQQ": {
		Last: 294.61,
		Bid:  294.61,
		Ask:  294.61,
	},
	"RSP": {
		Last: 138.07,
		Bid:  138.07,
		Ask:  138.07,
	},
	"VOO": {
		Last: 360,
		Bid:  360,
		Ask:  360,
	},
	"CSPX": {
		Last: 398.9,
		Bid:  398.9,
		Ask:  398.9,
	},
	"EFA": {
		Last: 63.77,
		Bid:  63.77,
		Ask:  63.77,
	},
	"MCHI": {
		Last: 55.81,
		Bid:  55.81,
		Ask:  55.81,
	},
	"XLP": {
		Last: 72.85,
		Bid:  72.85,
		Ask:  72.85,
	},
	"GLD": {
		Last: 170.09,
		Bid:  170.09,
		Ask:  170.09,
	},
	"IEF": {
		Last: 101.15,
		Bid:  101.15,
		Ask:  101.15,
	},
	"AGG": {
		Last: 101.05,
		Bid:  101.05,
		Ask:  101.05,
	},
	"KWEB": {
		Last: 33.41,
		Bid:  33.41,
		Ask:  33.41,
	},
	"XLV": {
		Last: 129.2,
		Bid:  129.2,
		Ask:  129.2,
	},
	"IEMG": {
		Last: 49.78,
		Bid:  49.78,
		Ask:  49.78,
	},
	"XLU": {
		Last: 69.01,
		Bid:  69.01,
		Ask:  69.01,
	},
	"XLK": {
		Last: 133.45,
		Bid:  133.45,
		Ask:  133.45,
	},
}
var model = MustLoadModelFromJson(MustOpenFile("./fixture/model-syfe-core.json"))
var portfolio Portfolio = portfolioImpl{
	Cash: 5000,
	Assets: map[Symbol]PortfolioAsset{
		"TLT": {
			Position: 93.7819,
		},
		"QQQ": {
			Position: 33.257,
		},
		"RSP": {
			Position: 62.3899,
		},
		"VOO": {
			Position: 1.1959,
		},
		"CSPX": {
			Position: 15,
		},
		"EFA": {
			Position: 90.2573,
		},
		"MCHI": {
			Position: 115.1172,
		},
		"XLP": {
			Position: 69.1687,
		},
		"GLD": {
			Position: 25.8042,
		},
		"IEF": {
			Position: 43.3724,
		},
		"AGG": {
			Position: 37.9126,
		},
		"KWEB": {
			Position: 112.0926,
		},
		"XLV": {
			Position: 18.6332,
		},
		"IEMG": {
			Position: 39.2866,
		},
		"XLU": {
			Position: 28.2551,
		},
		"XLK": {
			Position: 8.8983,
		},
	},
}

func main() {
	c := Calculator{
		BuyOnly:                true,
		SupportFractionalShare: false,
		SlotSize:               1,
		buyingPriceSelector:    PriceSelectorBid(),
		sellingPriceSelector:   PriceSelectorBid(),
		valuePriceSelector:     PriceSelectorBid(),
	}

	fmt.Printf("%+v", c.CalculateOrders(market, model, portfolio, 5000))
}
