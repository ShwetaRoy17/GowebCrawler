package fetcher

import (
	"io"
	"net/http"
	"time"
	"errors"
)


var client = &http.Client{
	Timeout: time.Second * 10,
}


func Get(url string)(string, error) {
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode == 200 {
		defer res.Body.Close()
		body,err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	if res.StatusCode == 400 {
		return "",errors.New("http error: bad request")
	
	}
	return "", nil
}