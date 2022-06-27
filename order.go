package main

type Order struct {
	Symbol     Symbol
	Type       OrderType
	Side       OrderSide
	Amount     float64
	LimitPrice float64
}

func (o Order) Value() float64 {
	return o.Amount * o.LimitPrice
}

type OrderType int

const (
	OrderTypeLimit = iota + 1
)

type OrderSide int

const (
	OrderSideBuy = iota + 1
	OrderSideSell
)
