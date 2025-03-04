package coinbasepro

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Secret     string
	Key        string
	Passphrase string
	HTTPClient *http.Client
	RetryCount int
}

type ClientConfig struct {
	BaseURL    string
	Key        string
	Passphrase string
	Secret     string
}

func NewClient() *Client {
	baseURL := os.Getenv("COINBASE_PRO_BASEURL")
	if baseURL == "" {
		baseURL = "https://api.pro.coinbase.com"
	}

	client := Client{
		BaseURL:    baseURL,
		Key:        os.Getenv("COINBASE_PRO_KEY"),
		Passphrase: os.Getenv("COINBASE_PRO_PASSPHRASE"),
		Secret:     os.Getenv("COINBASE_PRO_SECRET"),
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		RetryCount: 0,
	}

	if os.Getenv("COINBASE_PRO_SANDBOX") == "1" {
		client.UpdateConfig(&ClientConfig{
			BaseURL: "https://api-public.sandbox.pro.coinbase.com",
		})
	}

	return &client
}

func (c *Client) UpdateConfig(config *ClientConfig) {
	baseURL := config.BaseURL
	key := config.Key
	passphrase := config.Passphrase
	secret := config.Secret

	if baseURL != "" {
		c.BaseURL = baseURL
	}
	if key != "" {
		c.Key = key
	}
	if passphrase != "" {
		c.Passphrase = passphrase
	}
	if secret != "" {
		c.Secret = secret
	}
}

func (c *Client) Request(method, path string,
	params, body, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.RetryCount+1; i++ {
		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2*1000) * time.Millisecond
		time.Sleep(retryDuration)

		res, err = c.request(method, path, params, body, result)
		if res != nil && res.StatusCode == 429 {
			continue
		} else {
			break
		}
	}

	return res, err
}

func encodeParams(params interface{}) (string, error) {
	p := url.Values{}
	v := reflect.Indirect(reflect.ValueOf(params))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		field := v.FieldByName(name)
		if !field.IsZero() {
			value, ok := v.FieldByName(name).Interface().(string)
			if !ok {
				return "", errors.New("could not encode params")
			}
			p.Add(strings.ToLower(name), value)
		}
	}
	return p.Encode(), nil
}

func (c *Client) request(method, path string,
	params, body, result interface{}) (res *http.Response, err error) {
	if params != nil {
		urlPath, err := url.Parse(path)
		if err != nil {
			return nil, err
		}

		query, err := encodeParams(params)
		if err != nil {
			return nil, err
		}

		urlPath.RawQuery = query
		path = urlPath.String()
	}

	var data []byte
	b := bytes.NewReader(make([]byte, 0))
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return res, err
		}

		b = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, fullURL, b)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// XXX: Sandbox time is off right now
	if os.Getenv("TEST_COINBASE_OFFSET") != "" {
		inc, err := strconv.Atoi(os.Getenv("TEST_COINBASE_OFFSET"))
		if err != nil {
			return res, err
		}

		timestamp = strconv.FormatInt(time.Now().Unix()+int64(inc), 10)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Go Coinbase Pro Client 1.0")

	h, err := c.Headers(method, path, timestamp, string(data))
	if err != nil {
		return res, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		defer res.Body.Close()
		coinbaseError := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&coinbaseError); err != nil {
			return res, err
		}

		return res, error(coinbaseError)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(result); err != nil {
			return res, err
		}
	}

	return res, nil
}

// Headers generates a map that can be used as headers to authenticate a request
func (c *Client) Headers(method, url, timestamp, data string) (map[string]string, error) {
	h := make(map[string]string)
	h["CB-ACCESS-KEY"] = c.Key
	h["CB-ACCESS-PASSPHRASE"] = c.Passphrase
	h["CB-ACCESS-TIMESTAMP"] = timestamp

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := generateSig(message, c.Secret)
	if err != nil {
		return nil, err
	}
	h["CB-ACCESS-SIGN"] = sig
	return h, nil
}
