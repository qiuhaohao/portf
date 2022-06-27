package main

type Market interface {
	Last(s Symbol) float64
	Bid(s Symbol) float64
	Ask(s Symbol) float64
}

type PriceSelector func(market Market, s Symbol) float64

func PriceSelectorLast() PriceSelector {
	return func(market Market, s Symbol) float64 {
		return market.Last(s)
	}
}

func PriceSelectorBid() PriceSelector {
	return func(market Market, s Symbol) float64 {
		return market.Bid(s)
	}
}

func PriceSelectorAsk() PriceSelector {
	return func(market Market, s Symbol) float64 {
		return market.Ask(s)
	}
}

type staticMarket map[Symbol]marketPrice

type marketPrice struct {
	Last float64
	Bid  float64
	Ask  float64
}

func (sm staticMarket) Last(s Symbol) float64 {
	if p, ok := sm[s]; ok {
		return p.Last
	}
	return 0
}

func (sm staticMarket) Bid(s Symbol) float64 {
	if p, ok := sm[s]; ok {
		return p.Bid
	}
	return 0
}

func (sm staticMarket) Ask(s Symbol) float64 {
	if p, ok := sm[s]; ok {
		return p.Ask
	}
	return 0
}
