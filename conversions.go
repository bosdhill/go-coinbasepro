package coinbasepro

type Conversion struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	// Response fields
	ID            string `json:"id"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
}

func (c *Client) CreateConversion(newConversion *Conversion) (Conversion, error) {
	var savedConversion Conversion

	_, err := c.Request("POST", "/conversions", nil, newConversion, &savedConversion)
	return savedConversion, err
}
