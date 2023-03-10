package data

import (
	"testing"

	config "github.com/samims/ecommerceGO/currency/configs"
	"github.com/sirupsen/logrus"
)

func TestNewRates(t *testing.T) {
	cfg := config.NewViperConfig()
	tr, err := NewRates(logrus.New(), cfg)

	if err != nil {
		t.Fatal(err)
	}

	tr.log.Infoln("Rates %#v", tr.rates)
}
