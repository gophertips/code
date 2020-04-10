package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// API endpoints.
const (
	apiListAllBreeds = "https://dog.ceo/api/breeds/list/all"
	apiListByBreed   = "https://dog.ceo/api/breed/%s/list"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type envelope struct {
	Message json.RawMessage `json:"message"`
	Status  string          `json:"status"`
}

// Client implements the logic.
type Client struct {
	HTTPClient httpClient
}

// Breeds is used to fetch all the breeds.
func (c Client) Breeds(ctx context.Context) (breeds []string, err error) {
	payload, err := c.doRequest(ctx, apiListAllBreeds)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the request: %w", err)
	}

	keys, err := processKeys(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the breeds: %w", err)
	}

	return keys, nil
}

// SubBreeds returns all sub-breads of a given breed.
func (c Client) SubBreeds(ctx context.Context, parentBreed string) (subBreeds []string, err error) {
	payload, err := c.doRequest(ctx, fmt.Sprintf(apiListByBreed, parentBreed))
	if err != nil {
		return nil, fmt.Errorf("failed to execute the request: %w", err)
	}

	resp, err := processValues(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the sub-breeds: %w", err)
	}

	return resp, nil
}

func (c Client) doRequest(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the HTTP request: %w", err)
	}

	req = req.WithContext(ctx)
	rawResp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the breeds: %w", err)
	}
	defer rawResp.Body.Close()

	if rawResp.StatusCode < 200 || rawResp.StatusCode >= 300 {
		return nil, fmt.Errorf("invalid status code '%d'", rawResp.StatusCode)
	}

	payload, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response payload: %w", err)
	}

	var e envelope
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if e.Status != "success" {
		return nil, errors.New("invalid response")
	}

	return e.Message, nil
}

func processKeys(payload []byte) ([]string, error) {
	content := make(map[string]interface{})
	if err := json.Unmarshal(payload, &content); err != nil {
		return nil, fmt.Errorf("failed to parse the response: %w", err)
	}

	response := make([]string, 0, len(content))
	for key := range content {
		response = append(response, key)
	}

	return response, nil
}

func processValues(payload []byte) ([]string, error) {
	var fragments []interface{}
	if err := json.Unmarshal(payload, &fragments); err != nil {
		return nil, errors.New("could not parse the response")
	}

	resp := make([]string, 0, len(fragments))
	for _, rawFragment := range fragments {
		fragment, ok := rawFragment.(string)
		if !ok {
			return nil, errors.New("could not parse the response")
		}
		resp = append(resp, fragment)
	}

	return resp, nil
}
