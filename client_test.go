package coinbasepro

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

func TestClientErrorsOnNotFound(t *testing.T) {
	client := NewTestClient()
	_, err := client.Request("GET", "/fake", nil, nil, nil)
	if err == nil {
		t.Error(errors.New("Should have thrown 404 error"))
	}
}

func TestEncodeParams(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	params := DepositQueryParam{
		After: now,
		Limit: "1",
		Type:  "internal_deposit",
	}

	p, err := encodeParams(params)
	if err != nil {
		t.Error(err)
	}

	v := url.Values{}
	v.Add("after", params.After)
	v.Add("limit", params.Limit)
	v.Add("type", params.Type)

	if p != v.Encode() {
		t.Errorf("could not encode params")
	}
}
