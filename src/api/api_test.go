package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"netpart/api"
	"testing"
	"time"
)

const SERVER = "127.0.0.1:7000"
const BASE_URL = "http://" + SERVER

func wait(ctx context.Context) error {
	for {
		var client http.Client
		req, err := http.NewRequestWithContext(ctx, "GET", BASE_URL+"/ping", nil)
		if err != nil {
			return err
		}
		res, err := client.Do(req)
		if err != nil {
			time.Sleep(250 * time.Millisecond)
			continue
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			return nil
		}

		time.Sleep(250 * time.Millisecond)
	}
}

func TestListInstance(t *testing.T) {
	ctx := context.Background()
	go api.Run(ctx, SERVER)
	err := wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var client http.Client
	req, err := http.NewRequestWithContext(ctx, "GET", BASE_URL+"/instances", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("response not ok. got %v", res.StatusCode)
	}

	_, err = decode[api.ListInstanceResponse](res)
	if err != nil {
		t.Fatal(err)
	}
}

func decode[T any](r *http.Response) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
