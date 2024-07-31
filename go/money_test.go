package main

import (
	s "tdd/stocks"
	"testing"
)

func TestMultiplication(t *testing.T) {
	fiver := s.NewMoney(5, "USD")
	tenner := fiver.Times(2)
	expectedTenner := s.NewMoney(10, "USD")
	assertEqual(t, expectedTenner, tenner)
}

func TestDivision(t *testing.T) {
	originalMoney := s.NewMoney(4002, "KRW")
	actualMoneyAfterDivision := originalMoney.Divide(4)
	expectedMoneyAfterDivision := s.NewMoney(1000.5, "KRW")
	assertEqual(t, expectedMoneyAfterDivision, actualMoneyAfterDivision)
}

func TestAddition(t *testing.T) {
	var portfolio s.Portfolio
	var portfolioInDollars s.Money

	fiveDollars := s.NewMoney(5, "USD")
	tenDollars := s.NewMoney(10, "USD")
	fifteenDollars := s.NewMoney(15, "USD")

	portfolio = portfolio.Add(fiveDollars)
	portfolio = portfolio.Add(tenDollars)
	portfolioInDollars = portfolio.Evaluate("USD")

	assertEqual(t, fifteenDollars, portfolioInDollars)
}

func TestAdditionDollarsAndEuros(t *testing.T) {
	var portfolio s.Portfolio

	fiveDollars := s.NewMoney(5, "USD")
	tenEuros := s.NewMoney(10, "EUR")

	portfolio = portfolio.Add(fiveDollars)
	portfolio = portfolio.Add(tenEuros)

	expectedValue := s.NewMoney(17, "USD")
	actualValue := portfolio.Evaluate("USD")
	assertEqual(t, expectedValue, actualValue)
}

func assertEqual(t *testing.T, expected, actual s.Money) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
