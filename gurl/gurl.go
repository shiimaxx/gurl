package gurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// Client gurl client
type Client struct {
	Parallel int
	Output   string
}

// Content actual content
type Content struct {
	Name   string
	Length int
}

// NewClient constractor for Client
func NewClient(parallel int, output string) *Client {
	return &Client{
		Parallel: parallel,
		Output:   output,
	}
}

// Get content of url
func (c *Client) Get(url string) error {
	resp, err := http.Head(url)
	if err != nil {
		return err
	}

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	chunkSize := contentLength / c.Parallel
	surplus := contentLength % c.Parallel

	var wg sync.WaitGroup
	for i := 0; i < c.Parallel; i++ {
		s := i * chunkSize
		e := s + (chunkSize - 1)
		if i == c.Parallel-1 {
			e += surplus
		}

		wg.Add(1)
		go func(startRange, endRange, i int) {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startRange, endRange))
			resp, _ := client.Do(req)
			defer resp.Body.Close()

			reader, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile(fmt.Sprintf("%s.%d.tmp", c.Output, i), reader, 0755)
			wg.Done()
		}(s, e, i)
	}
	wg.Wait()

	return nil
}
