package blumapi

import "fmt"

func (c *Client) Bluming(rsp *TokenListRsp) error {
	return c.get("jetton/top/nearest_to_listing?published=include_listed", rsp)
}

func (c *Client) Hot(rsp *TokenListRsp) error {
	return c.get("jetton/top/transactions?published=exclude", rsp)
}

func (c *Client) Published(rsp *TokenListRsp) error {
	return c.get("jetton/top/published_at?published=only", rsp)
}

func (c *Client) Chart(shortName string, step int, rsp *ChartRsp) error {
	url := fmt.Sprintf("jetton/chart/%s?step=%d", shortName, step)

	return c.get(url, rsp)
}

func (c *Client) Live(rsp *LiveDataRsp) error {
	return c.get("jetton/live", rsp)
}
