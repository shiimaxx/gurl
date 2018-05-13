package gurl

// Doer http request doer
type Doer interface {
	doRequest(string) error
}

// Client gurl client
type Client struct {
	GurlDoer Doer
	Output   string
}

// Content actual content
type Content struct {
	Name   string
	Length int
}

// NewDoer constractor for Doer
func NewDoer() Doer {
	d := new(Doer)
	return *d
}

// NewClient constractor for Client
func NewClient(doer Doer, output string) *Client {
	return &Client{
		GurlDoer: doer,
		Output:   output,
	}
}

// Get content of url
func (c *Client) Get(url string) error {
	return nil
}
