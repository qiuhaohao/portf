package main

// Portfolio represents a collection of outstanding positions
type Portfolio interface {
	Symbols() Symbols
	Position(s Symbol) float64
	CashAmount() float64
	Select(ss Symbols) Portfolio
	AggregateEquivalents(market Market, model Model, selector PriceSelector) Portfolio
}

type portfolioImpl struct {
	Cash   float64
	Assets map[Symbol]PortfolioAsset
}

func (p portfolioImpl) Symbols() Symbols {
	ss := NewSymbols()
	for s := range p.Assets {
		ss = ss.Add(s)
	}

	return ss
}

func (p portfolioImpl) Position(s Symbol) float64 {
	if a, ok := p.Assets[s]; ok {
		return a.Position
	}
	return 0
}

func (p portfolioImpl) CashAmount() float64 {
	return p.Cash
}

func (p portfolioImpl) Select(ss Symbols) Portfolio {
	newP := portfolioImpl{
		Cash:   p.Cash,
		Assets: make(map[Symbol]PortfolioAsset),
	}

	for s, a := range p.Assets {
		if ss.Contains(s) {
			newP.Assets[s] = a
		}
	}

	return newP
}

func (p portfolioImpl) AggregateEquivalents(market Market, model Model, selector PriceSelector) Portfolio {
	newP := portfolioImpl{
		Cash:   p.Cash,
		Assets: make(map[Symbol]PortfolioAsset),
	}

	for s, a := range p.Assets {
		newP.Assets[s] = a
	}

	for _, s := range model.Symbols().ToSlice() {
		for _, e := range model.Equivalents(s).ToSlice() {
			if !newP.Symbols().Contains(e) {
				continue
			}
			newP.Assets[s] = PortfolioAsset{
				Position: p.Assets[s].Position +
					(p.Assets[e].Position*selector(market, e))/
						selector(market, s),
			}
			delete(newP.Assets, e)
		}
	}

	return newP
}

type PortfolioAsset struct {
	Position float64
}
