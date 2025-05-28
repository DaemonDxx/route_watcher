package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	http *http.Client
	req  *http.Request
}

func New() (*Client, error) {
	h := &http.Client{
		Transport: nil,
		Timeout:   20 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://srv.prt24.ru/api/voyage/routes", nil)
	if err != nil {
		return nil, fmt.Errorf("make request failed: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json, text/plain, */*")
	req.Header.Add("Origin", "https://bilet.prt24.ru")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	return &Client{
		http: h,
		req:  req,
	}, nil
}

func (c *Client) IsAvailableRoute(ctx context.Context, routeID int) (bool, error) {
	resCh := make(chan bool, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)

		res, err := c.http.Do(c.req)
		if err != nil {
			errCh <- fmt.Errorf("request failed: %v", err)
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			errCh <- fmt.Errorf("request failed: %v", res.Status)
		}

		routes, err := c.parseBody(res.Body)
		if err != nil {
			errCh <- fmt.Errorf("parse body failed: %v", err)
		}

		_, found := slices.BinarySearch(routes, routeID)

		resCh <- found
	}()

	select {
	case ok := <-resCh:
		return ok, nil
	case err := <-errCh:
		return false, err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (c *Client) parseBody(b io.Reader) ([]int, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %v", err)
	}
	body = body[1:]
	body = body[:len(body)-1]
	arr := strings.Split(string(body), ",")

	res := make([]int, len(arr))
	for _, v := range arr {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("parse body failed: %v", err)
		}
		res = append(res, i)
	}

	return res, nil
}
