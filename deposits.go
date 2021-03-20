package coinbasepro

import (
	"encoding/json"
)

type Deposit struct {
	Currency        string `json:"currency"`
	Amount          string `json:"amount"`
	PaymentMethodID string `json:"payment_method_id"` // PaymentMethodID can be determined by calling GetPaymentMethods()
	// Response fields
	ID       string `json:"id,omitempty"`
	PayoutAt Time   `json:"payout_at,string,omitempty"`
}

type Total struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
type Remaining struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
type Limit struct {
	PeriodInDays int       `json:"period_in_days"`
	Total        Total     `json:"total"`
	Remaining    Remaining `json:"remaining"`
}
type Limits struct {
	Buy        []Limit
	InstantBuy []Limit
	Sell       []Limit
	Desposit   []Limit
}
type PaymentMethod struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Currency      string `json:"currency"`
	PrimaryBuy    bool   `json:"primary_buy"`
	PrimarySell   bool   `json:"primary_sell"`
	AllowBuy      bool   `json:"allow_buy"`
	AllowSell     bool   `json:"allow_sell"`
	AllowDeposit  bool   `json:"allow_deposit"`
	AllowWithdraw bool   `json:"allow_withdraw"`
	Limits        Limits `json:"limits"`
}

func (p PaymentMethod) String() string {
	s, _ := json.MarshalIndent(p, "", "\t")
	return string(s)
}

func (c *Client) CreateDeposit(newDeposit *Deposit) (Deposit, error) {
	var savedDeposit Deposit

	_, err := c.Request("POST", "/deposits/payment-method", newDeposit, &savedDeposit)
	return savedDeposit, err
}

func (c *Client) GetPaymentMethods() ([]PaymentMethod, error) {
	var paymentMethods []PaymentMethod

	_, err := c.Request("GET", "/payment-methods", nil, &paymentMethods)
	return paymentMethods, err
}
