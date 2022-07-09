package req

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseUrl string = "https://sample-url.com/"
)

type Requests interface {
	GetFlightPaths(ctx context.Context, user string) ([][]string, error)
}

type requests struct{}

func NewRequests() Requests {
	return &requests{}
}

type flightRequest struct {
	User string `json:"user"`
}

func (r requests) GetFlightPaths(ctx context.Context, user string) ([][]string, error) {
	result := make([][]string, 0)

	data := flightRequest{User: user}
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(buf)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseUrl+"flight-paths", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 400 {
		var data interface{}
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("request failed with status: %d and err: %v", res.StatusCode, data)
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
