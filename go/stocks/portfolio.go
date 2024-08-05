package stocks

import (
	"fmt"
	"strings"
)

type Portfolio []Money

func (p Portfolio) Add(money Money) Portfolio {
	return append(p, money)
}

func (p Portfolio) Evaluate(currency string) (Money, error) {
	total := 0.0
	failedConversions := make([]string, 0, len(p))

	for _, m := range p {
		if converted, ok := convert(m, currency); ok {
			total += converted
		} else {
			failedConversions = append(failedConversions, m.currency+"->"+currency)
		}
	}

	if len(failedConversions) == 0 {
		return Money{amount: total, currency: currency}, nil
	}

	failures := fmt.Sprintf("Missing exchange rate(s):[%s]", strings.Join(failedConversions, ","))

	return Money{}, fmt.Errorf(failures)
}

func convert(m Money, currency string) (float64, bool) {
	if m.currency == currency {
		return m.amount, true
	}

	key := m.currency + "->" + currency
	exchangeRates := map[string]float64{
		"EUR->USD": 1.2,
		"USD->KRW": 1100,
	}

	rate, ok := exchangeRates[key]

	return m.amount * rate, ok
}
