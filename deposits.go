package coinbasepro

type DepositMethod struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
	// Required request field
	PaymentMethodID   string `json:"payment_method_id,omitempty"`   // PaymentMethodID can be determined by calling GetPaymentMethods()
	CoinbaseAccountID string `json:"coinbase_account_id,omitempty"` // CoinbaseAccountID from GetCoinbaseAccounts()
	// Response fields
	ID       string `json:"id,omitempty"`
	PayoutAt Time   `json:"payout_at,string,omitempty"`
}

func (c *Client) CreateDeposit(newDeposit *DepositMethod) (DepositMethod, error) {
	var savedDepositMethod DepositMethod

	_, err := c.Request("POST", "/deposits/payment-method", nil, newDeposit, &savedDepositMethod)
	return savedDepositMethod, err
}

func (c *Client) CreateCoinbaseDeposit(newDeposit *DepositMethod) (DepositMethod, error) {
	var savedDepositMethod DepositMethod

	_, err := c.Request("POST", "/deposits/coinbase-account", nil, newDeposit, &savedDepositMethod)
	return savedDepositMethod, err
}

type Deposit struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	CreatedAt   string      `json:"created_at"`
	CompletedAt string      `json:"completed_at"`
	CanceledAt  interface{} `json:"canceled_at"`
	ProcessedAt string      `json:"processed_at"`
	AccountID   string      `json:"account_id"`
	UserID      string      `json:"user_id"`
	UserNonce   interface{} `json:"user_nonce"`
	Amount      string      `json:"amount"`
	Details     Details
}

type Details struct {
	CryptoAddress         string `json:"crypto_address"`
	DestinationTag        string `json:"destination_tag"`
	CoinbaseAccountID     string `json:"coinbase_account_id"`
	DestinationTagName    string `json:"destination_tag_name"`
	CryptoTransactionID   string `json:"crypto_transaction_id"`
	CoinbaseTransactionID string `json:"coinbase_transaction_id"`
	CryptoTransactionHash string `json:"crypto_transaction_hash"`
}

type DepositQueryParam struct {
	// Set to deposit or internal_deposit (transfer between portfolios)
	Type string `json:"type"`

	// Limit list of deposits to this profile_id. By default, it retrieves deposits using default profile
	ProfileID string `json:"profile_id"`

	// If before is set, then it returns deposits created after the before timestamp, sorted by oldest creation date
	Before string `json:"before"`

	// If after is set, then it returns deposits created before the after timestamp, sorted by newest
	After string `json:"after"`

	// Truncate list to this many deposits, capped at 100. Default is 100.
	Limit string `json:"limit"`
}

func (c *Client) GetDeposits(params *DepositQueryParam) ([]Deposit, error) {
	var transfers []Deposit
	var err error

	// Hack to avoid using type assertions for nil checks in client.go
	if params != nil {
		_, err = c.Request("GET", "/transfers", params, nil, &transfers)
	} else {
		_, err = c.Request("GET", "/transfers", nil, nil, &transfers)
	}
	return transfers, err
}
