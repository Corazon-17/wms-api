package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) Authorize() (string, error) {

	timestamp := time.Now().Unix()

	path := "/oauth/authorize"

	base := fmt.Sprintf("%s%s%d%s",
		c.partnerID,
		path,
		timestamp,
		c.shopID,
	)

	sign := Sign(c.partnerKey, base)

	url := fmt.Sprintf(
		"%s%s?shop_id=%s&partner_id=%s&timestamp=%d&sign=%s&redirect=%s",
		c.baseURL,
		path,
		c.shopID,
		c.partnerID,
		timestamp,
		sign,
		c.redirect,
	)

	resp, err := c.http.Get(url)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Code string `json:"code"`
		} `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Data.Code, nil
}

func (c *Client) ExchangeToken(code string) error {

	url := fmt.Sprintf("%s/oauth/token", c.baseURL)

	body := map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}

	payload, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	var result TokenResponse

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	c.token.Set(
		result.Data.AccessToken,
		result.Data.RefreshToken,
		result.Data.ExpiresIn,
	)

	return nil
}

func (c *Client) RefreshToken() error {

	url := fmt.Sprintf("%s/oauth/token", c.baseURL)

	body := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": c.token.RefreshToken,
	}

	payload, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	var result TokenResponse

	if err := c.doRequest(req, &result); err != nil {
		return err
	}

	c.token.Set(
		result.Data.AccessToken,
		result.Data.RefreshToken,
		result.Data.ExpiresIn,
	)

	return nil
}

func (c *Client) ensureToken() error {

	if c.token.AccessToken == "" {
		code, err := c.Authorize()
		if err != nil {
			return err
		}

		return c.ExchangeToken(code)
	}

	if c.token.IsExpired() {
		return c.RefreshToken()
	}

	return nil
}
