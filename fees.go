package coinbasepro

type Fees struct {
	MakerFeeRate string `json:"maker_fee_rate"`
	TakerFeeRate string `json:"taker_fee_rate"`
	USDVolume    string `json:"usd_volume"`
}

func (c *Client) GetFees() (Fees, error) {
	var fees Fees

	_, err := c.Request("GET", "/fees", nil, nil, &fees)
	return fees, err
}
