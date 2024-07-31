package stocks

type Portfolio []Money

func (p Portfolio) Add(money Money) Portfolio {
	return append(p, money)
}

func (p Portfolio) Evaluate(currency string) Money {
	total := 0.0
	for _, m := range p {
		total += convert(m, currency)
	}

	return Money{amount: total, currency: currency}
}

func convert(m Money, currency string) float64 {
	if m.currency == currency {
		return m.amount
	}

	key := m.currency + "->" + currency
	exchangeRates := map[string]float64{
		"EUR->USD": 1.2,
		"USD->KRW": 1100,
	}

	return m.amount * exchangeRates[key]
}
