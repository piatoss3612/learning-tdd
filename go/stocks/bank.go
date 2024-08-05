package stocks

import "errors"

type Bank struct {
	exchangeRates map[string]float64
}

func (b Bank) AddExchangeRate(currencyFrom string, currencyTo string, rate float64) {
	key := currencyFrom + "->" + currencyTo
	b.exchangeRates[key] = rate
}

func (b Bank) Convert(money Money, currencyTo string) (*Money, error) {
	if money.currency == currencyTo {
		return &money, nil
	}

	key := money.currency + "->" + currencyTo
	rate, ok := b.exchangeRates[key]
	if !ok {
		return nil, errors.New(key)
	}

	return &Money{amount: money.amount * rate, currency: currencyTo}, nil
}

func NewBank() Bank {
	return Bank{exchangeRates: make(map[string]float64)}
}
