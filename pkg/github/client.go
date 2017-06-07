package github

type Client struct {
}

func (c *Client) PRReviewed(prURL string) (bool, error) {
	return false, nil
}
