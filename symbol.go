package main

type Symbol string

type Symbols interface {
	Add(ss ...Symbol) Symbols
	Contains(s Symbol) bool
	Union(other Symbols) Symbols
	ToSlice() []Symbol
}

func NewSymbols() Symbols {
	return make(symbolsImpl)
}

type symbolsImpl map[Symbol]struct{}

func (ss symbolsImpl) copy() symbolsImpl {
	newSs := make(symbolsImpl)
	for k := range ss {
		newSs[k] = struct{}{}
	}

	return newSs
}

func (ss symbolsImpl) Add(sList ...Symbol) Symbols {
	newSS := ss.copy()
	for _, s := range sList {
		newSS[s] = struct{}{}
	}
	return newSS
}

func (ss symbolsImpl) Contains(s Symbol) bool {
	_, ok := ss[s]
	return ok
}

func (ss symbolsImpl) ToSlice() []Symbol {
	symbols := make([]Symbol, 0)
	for s := range ss {
		symbols = append(symbols, s)
	}
	return symbols
}

func (ss symbolsImpl) Union(other Symbols) Symbols {
	return ss.Add(other.ToSlice()...)
}
