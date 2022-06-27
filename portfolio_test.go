package main

import (
	"reflect"
	"testing"
)

func Test_portfolioImpl_AggregateEquivalents(t *testing.T) {
	type fields struct {
		Cash   float64
		Assets map[Symbol]PortfolioAsset
	}
	type args struct {
		market   Market
		model    Model
		selector PriceSelector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Portfolio
	}{
		{
			name: "should not asset with no equivalents",
			fields: fields{
				Cash: 100,
				Assets: map[Symbol]PortfolioAsset{
					"TEST": {Position: 1},
				},
			},
			args: args{
				market: staticMarket{},
				model: modelImpl{Assets: map[Symbol]modelAsset{
					"TEST": {
						Equivalents: nil,
						Weight:      100,
					},
				}},
				selector: PriceSelectorBid(),
			},
			want: portfolioImpl{
				Cash: 100,
				Assets: map[Symbol]PortfolioAsset{
					"TEST": {Position: 1},
				},
			},
		},
		{
			name: "should aggregate equivalents",
			fields: fields{
				Cash: 100,
				Assets: map[Symbol]PortfolioAsset{
					"VOO":  {Position: 1},
					"CSPX": {Position: 2},
				},
			},
			args: args{
				market: staticMarket{
					"VOO": {
						Bid: 1,
					},
					"CSPX": {
						Bid: 2,
					},
				},
				model: modelImpl{Assets: map[Symbol]modelAsset{
					"VOO": {
						Equivalents: []Symbol{"CSPX"},
						Weight:      100,
					},
				}},
				selector: PriceSelectorBid(),
			},
			want: portfolioImpl{
				Cash: 100,
				Assets: map[Symbol]PortfolioAsset{
					"VOO": {Position: 5},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portfolioImpl{
				Cash:   tt.fields.Cash,
				Assets: tt.fields.Assets,
			}
			if got := p.AggregateEquivalents(tt.args.market, tt.args.model, tt.args.selector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateEquivalents() = %v, want %v", got, tt.want)
			}
		})
	}
}
