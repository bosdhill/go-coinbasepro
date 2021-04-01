package coinbasepro

type Currency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
}

func (c *Client) GetCurrencies() ([]Currency, error) {
	var currencies []Currency

	_, err := c.Request("GET", "/currencies", nil, nil, &currencies)
	return currencies, err
}
