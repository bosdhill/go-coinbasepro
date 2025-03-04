package coinbasepro

import (
	"fmt"
)

type Order struct {
	Type      string `json:"type"`
	Size      string `json:"size,omitempty"`
	Side      string `json:"side"`
	ProductID string `json:"product_id"`
	ClientOID string `json:"client_oid,omitempty"`
	Stp       string `json:"stp,omitempty"`
	Stop      string `json:"stop,omitempty"`
	StopPrice string `json:"stop_price,omitempty"`
	// Limit Order
	Price       string `json:"price,omitempty"`
	TimeInForce string `json:"time_in_force,omitempty"`
	PostOnly    bool   `json:"post_only,omitempty"`
	CancelAfter string `json:"cancel_after,omitempty"`
	// Market Order
	Funds          string `json:"funds,omitempty"`
	SpecifiedFunds string `json:"specified_funds,omitempty"`
	// Response Fields
	ID            string `json:"id"`
	Status        string `json:"status,omitempty"`
	Settled       bool   `json:"settled,omitempty"`
	DoneReason    string `json:"done_reason,omitempty"`
	DoneAt        Time   `json:"done_at,string,omitempty"`
	CreatedAt     Time   `json:"created_at,string,omitempty"`
	FillFees      string `json:"fill_fees,omitempty"`
	FilledSize    string `json:"filled_size,omitempty"`
	ExecutedValue string `json:"executed_value,omitempty"`
}

type CancelAllOrdersParams struct {
	ProductID string
}

type ListOrdersParams struct {
	Status     string
	ProductID  string
	Pagination PaginationParams
}

func (c *Client) CreateOrder(newOrder *Order) (Order, error) {
	var savedOrder Order

	if len(newOrder.Type) == 0 {
		newOrder.Type = "limit"
	}

	_, err := c.Request("POST", "/orders", nil, newOrder, &savedOrder)
	return savedOrder, err
}

func (c *Client) CancelOrder(id string) error {
	url := fmt.Sprintf("/orders/%s", id)
	_, err := c.Request("DELETE", url, nil, nil, nil)
	return err
}

func (c *Client) CancelAllOrders(p ...CancelAllOrdersParams) ([]string, error) {
	var orderIDs []string
	url := "/orders"

	if len(p) > 0 && p[0].ProductID != "" {
		url = fmt.Sprintf("%s?product_id=%s", url, p[0].ProductID)
	}

	_, err := c.Request("DELETE", url, nil, nil, &orderIDs)
	return orderIDs, err
}

func (c *Client) GetOrder(id string) (Order, error) {
	var savedOrder Order

	url := fmt.Sprintf("/orders/%s", id)
	_, err := c.Request("GET", url, nil, nil, &savedOrder)
	return savedOrder, err
}

func (c *Client) ListOrders(p ...ListOrdersParams) *Cursor {
	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
		if p[0].Status != "" {
			paginationParams.AddExtraParam("status", p[0].Status)
		}
		if p[0].ProductID != "" {
			paginationParams.AddExtraParam("product_id", p[0].ProductID)
		}
	}

	return NewCursor(c, "GET", fmt.Sprintf("/orders"),
		&paginationParams)
}
