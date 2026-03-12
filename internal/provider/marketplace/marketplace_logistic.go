package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ShipOrder(orderSN string, channelID string) (*ShipResponse, error) {

	if err := c.ensureToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/logistic/ship", c.baseURL)

	body := map[string]string{
		"order_sn":   orderSN,
		"channel_id": channelID,
	}

	payload, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token.Get())
	req.Header.Set("Content-Type", "application/json")

	var result ShipResponse

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
