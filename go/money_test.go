package main

import "testing"

type Money struct {
	amount   float64
	currency string
}

func (m Money) times(multiplier float64) Money {
	return Money{
		amount:   m.amount * multiplier,
		currency: m.currency,
	}
}

func (m Money) Divide(divisor float64) Money {
	return Money{
		amount:   m.amount / divisor,
		currency: m.currency,
	}
}

func TestMultiplication(t *testing.T) {
	fiver := Money{
		amount:   5,
		currency: "USD",
	}
	tenner := fiver.times(2)
	expectedTenner := Money{
		amount:   10,
		currency: "USD",
	}
	assertEqual(t, expectedTenner, tenner)
}

func TestMultiplicationInEuros(t *testing.T) {
	tenEuros := Money{
		amount:   10,
		currency: "EUR",
	}
	twentyEuros := tenEuros.times(2)
	expectedTwentyEuros := Money{
		amount:   20,
		currency: "EUR",
	}
	assertEqual(t, expectedTwentyEuros, twentyEuros)
}

func TestDivision(t *testing.T) {
	originalMoney := Money{amount: 4002, currency: "KRW"}
	actualMoneyAfterDivision := originalMoney.Divide(4)
	expectedMoneyAfterDivision := Money{amount: 1000.5, currency: "KRW"}
	assertEqual(t, expectedMoneyAfterDivision, actualMoneyAfterDivision)
}

func assertEqual(t *testing.T, expected, actual Money) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
