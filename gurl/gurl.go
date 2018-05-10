package gurl

// Client clinet
type Client struct {
	Output string
}

// NewClient constractor for Client
func NewClient(output string) *Client {
	c := new(Client)
	c.Output = output
	return c
}

// Get access url
func (c *Client) Get(url string) error {
	return nil
}
