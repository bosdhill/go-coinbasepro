package coinbasepro

import (
	"testing"
)

func TestGetPaymentMethods(t *testing.T) {
	client := NewTestClient()
	paymentMethods, err := client.GetPaymentMethods()
	if err != nil {
		t.Error(err)
	}

	for _, p := range paymentMethods {
		err = EnsureProperties(p, []string{"ID", "Type", "Name", "Currency"})
		if err != nil {
			t.Error(err)
		}
	}
}
