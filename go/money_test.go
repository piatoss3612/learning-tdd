package main

import "testing"

type Dollar struct {
	amount int
}

func (d Dollar) times(multiplier int) Dollar {
	return Dollar{
		amount: d.amount * multiplier,
	}
}

func TestMultiplication(t *testing.T) {
	fiver := Dollar{
		amount: 5,
	}
	tenner := fiver.times(2)
	if tenner.amount != 10 {
		t.Errorf("Expected 10, got %d", tenner.amount)
	}
}
