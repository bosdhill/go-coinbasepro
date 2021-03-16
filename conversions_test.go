package coinbasepro

import (
	"errors"
	"testing"
)

func TestCreateConversions(t *testing.T) {
	client := NewTestClient()

	conversion := &Conversion{
		From:   "USD",
		To:     "USDC",
		Amount: "10.00000000",
	}

	savedConversion, err := client.CreateConversion(conversion)

	if err != nil {
		t.Error(err)
	}

	if StructHasZeroValues(savedConversion) {
		t.Error(errors.New("Zero value"))
	}

	props := []string{"From", "To", "Amount"}
	_, err = CompareProperties(conversion, savedConversion, props)
	if err != nil {
		t.Error(err)
	}
}
