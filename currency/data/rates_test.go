package data

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewRates(t *testing.T) {
	tr, err := NewRates(logrus.New())

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Rates %#v", tr.rates)
}
