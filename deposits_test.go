package coinbasepro

import (
	"errors"
	"testing"
	"time"
)

func TestCreateDepositFromPaymentMethod(t *testing.T) {
	client := NewTestClient()
	paymentMethods, err := client.GetPaymentMethods()
	if err != nil {
		t.Error(err)
	}

	if len(paymentMethods) == 0 {
		t.Error(errors.New("no payment methods found"))
	}

	for _, p := range paymentMethods {
		if p.Type == "ach_bank_account" {
			depositMethod := &DepositMethod{
				Currency:        "USD",
				Amount:          "10.00000000",
				PaymentMethodID: p.ID,
			}

			d, err := client.CreateDeposit(depositMethod)
			if err != nil {
				t.Error(err)
			}

			err = EnsureProperties(d, []string{"ID", "Amount", "Currency", "PayoutAt"})
			if err != nil {
				t.Error(err)
			}

			_, err = CompareProperties(d, depositMethod, []string{"Amount", "Currency"})
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestCreateDepositFromCoinbase(t *testing.T) {
	client := NewTestClient()
	accounts, err := client.GetCoinbaseAccounts()
	if err != nil {
		t.Error(err)
	}

	if len(accounts) == 0 {
		t.Error(errors.New("No accounts found"))
	}

	for _, a := range accounts {
		if a.Name == "Fake" {
			depositMethod := &DepositMethod{
				Currency:          "BTC",
				Amount:            "10.00000000",
				CoinbaseAccountID: a.ID,
			}

			d, err := client.CreateCoinbaseDeposit(depositMethod)
			if err != nil {
				t.Error(err)
			}

			err = EnsureProperties(d, []string{"ID", "Amount", "Currency"})
			if err != nil {
				t.Error(err)
			}

			_, err = CompareProperties(d, depositMethod, []string{"Amount", "Currency"})
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestGetDeposits(t *testing.T) {
	client := NewTestClient()

	tests := []struct {
		query   *DepositQueryParam
		atLeast int
	}{
		{
			query:   nil,
			atLeast: 1,
		},
		{
			query: &DepositQueryParam{
				Type:  "internal_deposit",
				Limit: "1",
				After: time.Now().UTC().Format(time.RFC3339),
			},
			atLeast: 0,
		},
	}

	for _, tt := range tests {
		deposits, err := client.GetDeposits(tt.query)
		if err != nil {
			t.Error(err)
		}

		if len(deposits) < tt.atLeast {
			t.Fatalf("Expected at least %v deposits. Got %v", tt.atLeast, len(deposits))
		}
	}
}
