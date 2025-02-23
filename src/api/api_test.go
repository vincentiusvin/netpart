package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"netpart/api"
	"netpart/control"
	"testing"
	"time"
)

const SERVER = "127.0.0.1:7000"
const BASE_URL = "http://" + SERVER

func TestMain(m *testing.M) {
	ctx := context.Background()
	go api.Run(ctx, SERVER)
	err := wait(ctx)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestListInstance(t *testing.T) {
	ctx := context.Background()

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

func TestAddInstance(t *testing.T) {
	ctx := context.Background()
	go api.Run(ctx, SERVER)
	err := wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var inst control.Instance

	t.Run("add instance", func(t *testing.T) {
		inst, err = addRequest(ctx, "test")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("delete instance", func(t *testing.T) {
		_, err = deleteRequest(ctx, inst.Name)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func deleteRequest(ctx context.Context, name string) (api.KillInstanceResponse, error) {
	var client http.Client
	var resp api.KillInstanceResponse

	req, err := http.NewRequestWithContext(ctx, "DELETE", BASE_URL+"/instances/"+name, nil)
	if err != nil {
		return resp, err
	}

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if res.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("response not ok. got %v", res.StatusCode)
	}

	val, err := decode[api.KillInstanceResponse](res)
	if err != nil {
		return resp, err
	}
	return val, nil
}

func addRequest(ctx context.Context, name string) (api.AddInstanceSuccessResponse, error) {
	var client http.Client
	var resp api.AddInstanceSuccessResponse
	body, err := encode(api.AddInstanceBody{
		Name: name,
	})
	if err != nil {
		return resp, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", BASE_URL+"/instances", body)
	if err != nil {
		return resp, err
	}

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if res.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("response not ok. got %v", res.StatusCode)
	}

	val, err := decode[api.AddInstanceSuccessResponse](res)
	if err != nil {
		return resp, err
	}
	return val, nil
}

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

func decode[T any](r *http.Response) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func encode[T any](v T) (io.Reader, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
