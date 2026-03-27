package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

func Fetch(url string) (string, error) {
	res, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("network Error %w", err)
	}

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("error reading body :%w", err)
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				fmt.Printf("Error closing response body: %v\n", err)
			}
		}()
		return string(body), nil
	}

	return "", fmt.Errorf("http error: bad request %d", res.StatusCode)

}
