package coinbasepro

import (
	"fmt"
)

type Account struct {
	ID        string `json:"id"`
	Balance   string `json:"balance"`
	Hold      string `json:"hold"`
	Available string `json:"available"`
	Currency  string `json:"currency"`
}
type WireDepositInfo struct {
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	BankName      string `json:"bank_name"`
	BankAddress   string `json:"bank_address"`
	BankCountry   struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"bank_country"`
	AccountName    string `json:"account_name"`
	AccountAddress string `json:"account_address"`
	Reference      string `json:"reference"`
}

type SepaDepsoitInfo struct {
	Iban            string `json:"iban"`
	Swift           string `json:"swift"`
	BankName        string `json:"bank_name"`
	BankAddress     string `json:"bank_address"`
	BankCountryName string `json:"bank_country_name"`
	AccountName     string `json:"account_name"`
	AccountAddress  string `json:"account_address"`
	Reference       string `json:"reference"`
}

type CoinbaseAccount struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Balance  string          `json:"balance"`
	Currency string          `json:"currency"`
	Type     string          `json:"type"`
	Primary  bool            `json:"primary"`
	Active   bool            `json:"active"`
	Wire     WireDepositInfo `json:"wire_deposit_information"`
	Sepa     SepaDepsoitInfo `json:"sepa_deposit_information"`
}

// Ledger

type LedgerEntry struct {
	ID        string        `json:"id,number"`
	CreatedAt Time          `json:"created_at,string"`
	Amount    string        `json:"amount"`
	Balance   string        `json:"balance"`
	Type      string        `json:"type"`
	Details   LedgerDetails `json:"details"`
}

type LedgerDetails struct {
	OrderID   string `json:"order_id"`
	TradeID   string `json:"trade_id"`
	ProductID string `json:"product_id"`
}

type GetAccountLedgerParams struct {
	Pagination PaginationParams
}

// Holds

type Hold struct {
	AccountID string `json:"account_id"`
	CreatedAt Time   `json:"created_at,string"`
	UpdatedAt Time   `json:"updated_at,string"`
	Amount    string `json:"amount"`
	Type      string `json:"type"`
	Ref       string `json:"ref"`
}

type ListHoldsParams struct {
	Pagination PaginationParams
}

// Client Funcs
func (c *Client) GetAccounts() ([]Account, error) {
	var accounts []Account
	_, err := c.Request("GET", "/accounts", nil, nil, &accounts)

	return accounts, err
}

func (c *Client) GetAccount(id string) (Account, error) {
	account := Account{}

	url := fmt.Sprintf("/accounts/%s", id)
	_, err := c.Request("GET", url, nil, nil, &account)
	return account, err
}

func (c *Client) GetCoinbaseAccounts() ([]CoinbaseAccount, error) {
	var accounts []CoinbaseAccount

	_, err := c.Request("GET", "/coinbase-accounts/", nil, nil, &accounts)
	return accounts, err
}

func (c *Client) ListAccountLedger(id string,
	p ...GetAccountLedgerParams) *Cursor {
	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
	}

	return NewCursor(c, "GET", fmt.Sprintf("/accounts/%s/ledger", id),
		&paginationParams)
}

func (c *Client) ListHolds(id string, p ...ListHoldsParams) *Cursor {
	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
	}

	return NewCursor(c, "GET", fmt.Sprintf("/accounts/%s/holds", id),
		&paginationParams)
}
