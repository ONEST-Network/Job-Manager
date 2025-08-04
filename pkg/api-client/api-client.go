package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Interface interface {
	ApiCall(request interface{}, url string, response interface{}, method string) error
}

type APIClient struct{}

func NewAPIClient() Interface {
	return &APIClient{}
}

func (a *APIClient) ApiCall(request interface{}, url string, response interface{}, method string) error {
	var (
		data []byte = nil
		err  error
	)

	if request != nil {
		data, err = json.Marshal(request)
		if err != nil {
			return err
		}
	}

	statusCode, body, err := restCall(method, url, data)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("unexpected call response: %s, statuscode %v", string(body), statusCode)
	}

	if response != nil {
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
	}

	return nil
}

func restCall(method string, urlStr string, payload []byte) (int, []byte, error) {
	transport := &http.Transport{
		DisableKeepAlives: true,
		Proxy:             http.ProxyFromEnvironment,
	}

	client := http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	defer client.CloseIdleConnections()

	req, err := http.NewRequestWithContext(context.Background(), method, urlStr, bytes.NewBuffer(payload))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request, %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send request, %v", err)
	}

	if resp == nil || resp.Body == nil {
		return 0, nil, errors.New("received nil response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("failed to read response body, %v", err)
	}

	defer resp.Body.Close()

	return resp.StatusCode, body, nil
}
