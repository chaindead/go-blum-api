package blumapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type refreshRsp struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (c *Client) authRefresh() error {
	url := "https://user-domain.blum.codes/api/v1/auth/refresh"
	method := "POST"

	payloadString := fmt.Sprintf(`{"refresh":"%s"}`, c.jwtRefresh)
	payload := strings.NewReader(payloadString)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	var rsp refreshRsp
	if err := json.Unmarshal(body, &rsp); err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}

	c.jwtRefresh = rsp.Refresh
	c.jwtAccess = rsp.Access

	return nil
}

type authCreateRsp struct {
	Token struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
		User    struct {
			Id struct {
				Id string `json:"id"`
			} `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"token"`
	JustCreated bool `json:"justCreated"`
}

func (c *Client) authCreate(query string) error {
	url := "https://user-domain.blum.codes/api/v1/auth/provider/PROVIDER_TELEGRAM_MINI_APP"
	method := "POST"

	payloadText := fmt.Sprintf(`{
    "query": "query_id=%s"
}`, query)
	payload := strings.NewReader(payloadText)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	var rsp authCreateRsp
	if err := json.Unmarshal(body, &rsp); err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}

	c.jwtRefresh = rsp.Token.Refresh
	c.jwtAccess = rsp.Token.Access

	return nil
}
