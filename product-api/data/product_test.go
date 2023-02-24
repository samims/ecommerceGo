package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "latte",
		Price: 1.0,
		SKU:   "abc-addf-xdf",
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
