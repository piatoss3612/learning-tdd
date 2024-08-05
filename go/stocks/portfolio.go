package stocks

import (
	"fmt"
	"strings"
)

type Portfolio []Money

func (p Portfolio) Add(money Money) Portfolio {
	return append(p, money)
}

func (p Portfolio) Evaluate(bank Bank, currency string) (*Money, error) {
	total := 0.0
	failedConversions := make([]string, 0, len(p))

	for _, m := range p {
		if converted, err := bank.Convert(m, currency); err == nil {
			total += converted.amount
		} else {
			failedConversions = append(failedConversions, err.Error())
		}
	}

	if len(failedConversions) == 0 {
		return &Money{amount: total, currency: currency}, nil
	}

	failures := fmt.Sprintf("Missing exchange rate(s):[%s]", strings.Join(failedConversions, ","))

	return nil, fmt.Errorf(failures)
}
