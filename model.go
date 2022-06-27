package main

import (
	"encoding/json"
	"errors"
	"io"
)

// Model represents an immutable target weighted allocation of assets
type Model interface {
	Symbols() Symbols
	TargetProportion(s Symbol) float64
	Equivalents(s Symbol) Symbols
	Contains(s Symbol) bool
	IsRelevant(s Symbol) bool
	Validate() error
}

type modelImpl struct {
	Assets map[Symbol]modelAsset `json:"assets"`
}

func LoadModelFromJson(r io.Reader) (Model, error) {
	m := new(modelImpl)
	if err := json.NewDecoder(r).Decode(m); err != nil {
		return nil, err
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return m, nil
}

func MustLoadModelFromJson(r io.Reader) Model {
	m, err := LoadModelFromJson(r)
	if err != nil {
		panic(err)
	}
	return m
}

func (m modelImpl) Symbols() Symbols {
	ss := NewSymbols()
	for s := range m.Assets {
		ss = ss.Add(s)
	}

	return ss
}

func (m modelImpl) TargetProportion(s Symbol) float64 {
	if m.Contains(s) {
		return m.Assets[s].Weight / m.totalWeight()
	}
	return 0
}

func (m modelImpl) totalWeight() (totalWeight float64) {
	for _, a := range m.Assets {
		totalWeight += a.Weight
	}
	return
}
func (m modelImpl) Equivalents(s Symbol) Symbols {
	if m.Contains(s) {
		return NewSymbols().Add(m.Assets[s].Equivalents...)
	}
	return NewSymbols()
}

func (m modelImpl) Contains(s Symbol) bool {
	_, ok := m.Assets[s]
	return ok
}

func (m modelImpl) IsRelevant(s Symbol) bool {
	return m.relevantSymbols().Contains(s)
}

func (m modelImpl) relevantSymbols() Symbols {
	rs := m.Symbols()
	for _, a := range m.Assets {
		for _, e := range a.Equivalents {
			rs = rs.Add(e)
		}
	}

	return rs
}

func (m modelImpl) equivalentSymbols() Symbols {
	ss := NewSymbols()
	for _, a := range m.Assets {
		for _, e := range a.Equivalents {
			ss = ss.Add(e)
		}
	}

	return ss
}

var (
	ErrModelIsEmpty          = errors.New("model is empty")
	ErrContainsEquivalent    = errors.New("model contains equivalent")
	ErrDuplicatedEquivalents = errors.New("an equivalent symbol is linked to more than one symbol in the asset")
	ErrNonPositiveWeight     = errors.New("model contains non-positive weight")
)

func (m modelImpl) Validate() error {
	var validationFns = []func() error{
		m.checkNotEmpty,
		m.checkPositiveWeights,
		m.checkNoEquivalent,
		m.checkSymbolEquivalentOneToMany,
	}

	for _, fn := range validationFns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (m modelImpl) checkNoEquivalent() error {
	for _, a := range m.Assets {
		for _, e := range a.Equivalents {
			if m.Contains(e) {
				return ErrContainsEquivalent
			}
		}
	}

	return nil
}

func (m modelImpl) checkSymbolEquivalentOneToMany() error {
	equivalentsSeen := NewSymbols()
	for _, a := range m.Assets {
		for _, e := range a.Equivalents {
			if equivalentsSeen.Contains(e) {
				return ErrDuplicatedEquivalents
			}
			equivalentsSeen = equivalentsSeen.Add(e)
		}
	}
	return nil
}

func (m modelImpl) checkPositiveWeights() error {
	for _, a := range m.Assets {
		if a.Weight < 0 {
			return ErrNonPositiveWeight
		}
	}
	return nil
}

func (m modelImpl) checkNotEmpty() error {
	if len(m.Assets) == 0 {
		return ErrModelIsEmpty
	}
	return nil
}

// modelAsset represents an asset in a Model
type modelAsset struct {
	// Equivalents of an asset are considered to be the same as itself
	// if a Symbol appears as an modelAsset's Equivalents, it is not allowed
	// to appear as a modelAsset by itself in the same Model.
	Equivalents []Symbol `json:"equivalents"`
	Weight      float64  `json:"weight"`
}
