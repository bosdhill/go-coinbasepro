package coinbasepro

type Transfer struct {
	Type              string `json:"type"`
	Amount            string `json:"amount"`
	CoinbaseAccountID string `json:"coinbase_account_id,string"`
}

func (c *Client) CreateTransfer(newTransfer *Transfer) (Transfer, error) {
	var savedTransfer Transfer

	_, err := c.Request("POST", "/deposits/coinbase-account", nil, newTransfer, &savedTransfer)
	return savedTransfer, err
}
