package marketplace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const maxRetries = 3

func (c *Client) doRequest(req *http.Request, result interface{}) error {

	for attempt := 0; attempt < maxRetries; attempt++ {

		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		switch resp.StatusCode {

		case http.StatusOK:
			return json.NewDecoder(resp.Body).Decode(result)

		case http.StatusUnauthorized:
			log.Println("marketplace 401 → refreshing token")

			if err := c.RefreshToken(); err != nil {
				return err
			}

			req.Header.Set("Authorization", "Bearer "+c.token.Get())
			continue

		case http.StatusTooManyRequests:
			wait := time.Duration(attempt+1) * time.Second
			log.Printf("marketplace rate limited → retrying in %s", wait)

			time.Sleep(wait)
			continue

		case http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable:

			wait := time.Duration(attempt+1) * time.Second
			log.Printf("marketplace server error (%d) retrying in %s", resp.StatusCode, wait)

			time.Sleep(wait)
			continue

		default:
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("marketplace error: %s", string(body))
		}
	}

	return errors.New("marketplace request failed after retries")
}
