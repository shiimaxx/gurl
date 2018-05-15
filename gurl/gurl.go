package gurl

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/sync/errgroup"
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

	eg := errgroup.Group{}
	tmpFiles := make([]io.Reader, c.Parallel)
	tmpFileNames := make([]string, c.Parallel)

	for p := 0; p < c.Parallel; p++ {
		s := p * chunkSize
		e := s + (chunkSize - 1)
		if p == c.Parallel-1 {
			e += surplus
		}

		i := p
		eg.Go(func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}
			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", s, e))
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			reader, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			filename := fmt.Sprintf("%s.%d.tmp", c.Output, i)
			if err := ioutil.WriteFile(filename, reader, 0644); err != nil {
				return err
			}
			tmpFiles[i], err = os.Open(filename)
			if err != nil {
				return err
			}
			tmpFileNames[i] = filename
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	reader := io.MultiReader(tmpFiles...)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.Output, b, 0644); err != nil {
		return err
	}

	for _, f := range tmpFileNames {
		os.Remove(f)
	}

	return nil
}
