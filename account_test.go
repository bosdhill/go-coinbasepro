package coinbasepro

import (
	"errors"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	// Check for decoding issues
	for _, a := range accounts {
		if StructHasZeroValues(a) {
			t.Error(errors.New("Zero value"))
		}
	}
}

func TestGetAccount(t *testing.T) {
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	for _, a := range accounts {
		account, err := client.GetAccount(a.ID)
		if err != nil {
			t.Error(err)
		}

		// Check for decoding issues
		if StructHasZeroValues(account) {
			t.Error(errors.New("Zero value"))
		}
	}
}

func TestGetCoinbaseAccounts(t *testing.T) {
	client := NewTestClient()
	accounts, err := client.GetCoinbaseAccounts()
	if err != nil {
		t.Error(err)
	}

	if len(accounts) == 0 {
		t.Error(errors.New("No accounts found"))
	}

	sandboxAccounts := map[string]int{
		"Fake":          1,
		"All the Ether": 1,
		"USDC Wallet":   1,
		"BAT Wallet":    1,
		"EUR Wallet":    1,
		"GBP Wallet":    1,
		"Zero Hero":     1,
		"LINK Wallet":   1,
	}

	var totalNames int
	for _, a := range accounts {
		totalNames += sandboxAccounts[a.Name]
	}

	if totalNames != len(sandboxAccounts) {
		t.Error(errors.New("Wrong sandbox account names"))
	}
}

func TestListAccountLedger(t *testing.T) {
	var ledgers []LedgerEntry
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	for _, a := range accounts {
		cursor := client.ListAccountLedger(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledgers); err != nil {
				t.Error(err)
			}

			for _, ledger := range ledgers {
				props := []string{"ID", "CreatedAt", "Amount", "Balance", "Type"}
				if err := EnsureProperties(ledger, props); err != nil {
					t.Error(err)
				}

				if ledger.Type == "match" || ledger.Type == "fee" {
					if err := Ensure(ledger.Details); err != nil {
						t.Error(errors.New("Details is missing"))
					}
				}
			}
		}
	}
}

func TestListHolds(t *testing.T) {
	var holds []Hold
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	for _, a := range accounts {
		cursor := client.ListHolds(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(&holds); err != nil {
				t.Error(err)
			}

			for _, h := range holds {
				// Check for decoding issues
				if StructHasZeroValues(h) {
					t.Error(errors.New("Zero value"))
				}
			}
		}
	}
}
