package marketplace

import (
	"fmt"
	"net/http"
)

func (c *Client) ListOrders() ([]Order, error) {

	if err := c.ensureToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/order/list", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token.Get())

	var result OrderListResponse

	if err := c.doRequest(req, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
