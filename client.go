package blumapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const basePath = "https://mempad-domain.blum.codes/api/v1/"

type Client struct {
	BaseURL     string
	rateLimiter *rate.Limiter
	HTTPClient  *http.Client

	jwtAccess, jwtRefresh string
}

func (c *Client) authHeader() string {
	return fmt.Sprintf("Bearer %s", c.jwtAccess)
}

func NewClient(query string) (*Client, error) {
	c := &Client{
		BaseURL: basePath,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		rateLimiter: rate.NewLimiter(rate.Limit(10), 1),
	}

	err := c.authCreate(query)
	if err != nil {
		return nil, fmt.Errorf("refresh JWT: %w", err)
	}

	go c.authUpdater()

	return c, nil
}

func (c *Client) get(endpoint string, result any) error {
	if err := c.rateLimiter.Wait(context.Background()); err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("authorization", c.authHeader())
	req.Header.Add("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("cache-control", "no-cache")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) authUpdater() {
	ticker := time.NewTicker(50 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C

		err := c.authRefresh()
		if err != nil {
			log.Printf("refresh JWT: %v", err)
		}
	}
}
