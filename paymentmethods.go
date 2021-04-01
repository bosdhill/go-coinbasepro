package coinbasepro

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

type Limits struct {
	Buy        []Limit
	InstantBuy []Limit
	Sell       []Limit
	Desposit   []Limit
}

type Limit struct {
	PeriodInDays int       `json:"period_in_days"`
	Total        Total     `json:"total"`
	Remaining    Remaining `json:"remaining"`
}

type Total struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Remaining struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func (c *Client) GetPaymentMethods() ([]PaymentMethod, error) {
	var paymentMethods []PaymentMethod

	_, err := c.Request("GET", "/payment-methods", nil, nil, &paymentMethods)
	return paymentMethods, err
}
