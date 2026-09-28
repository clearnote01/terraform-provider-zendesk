package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nukosuke/go-zendesk/zendesk"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	zd, err := zendesk.NewClient(nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	if err := zd.SetEndpointURL(server.URL); err != nil {
		t.Fatalf("could not set endpoint: %v", err)
	}

	return &Client{Client: *zd}
}

func TestCreateOAuthClient(t *testing.T) {
	enforce := false
	var sent struct {
		Client map[string]interface{} `json:"client"`
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/oauth/clients.json" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &sent); err != nil {
			t.Fatalf("could not decode request body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"client":{"id":42,"name":"Example","identifier":"example","kind":"confidential","secret":"full-secret"}}`))
	})

	created, err := c.CreateOAuthClient(context.Background(), OAuthClient{
		Name:                   "Example",
		Identifier:             "example",
		Kind:                   "confidential",
		EnforceTokenExpiration: &enforce,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if created.ID != 42 || created.Secret != "full-secret" {
		t.Fatalf("unexpected response %+v", created)
	}

	// An explicit false must reach Zendesk, which otherwise defaults it to true.
	if v, ok := sent.Client["enforce_token_expiration"]; !ok || v != false {
		t.Fatalf("enforce_token_expiration was not sent as false: %v", sent.Client)
	}
}

func TestGetOAuthClient(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/oauth/clients/42.json" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"client":{"id":42,"name":"Example","scopes":["tickets:read"]}}`))
	})

	got, err := c.GetOAuthClient(context.Background(), 42)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.ID != 42 || len(got.Scopes) != 1 || got.Scopes[0] != "tickets:read" {
		t.Fatalf("unexpected response %+v", got)
	}
}

func TestUpdateOAuthClient(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/oauth/clients/42.json" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"client":{"id":42,"description":"updated"}}`))
	})

	got, err := c.UpdateOAuthClient(context.Background(), 42, OAuthClient{Description: "updated"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if got.Description != "updated" {
		t.Fatalf("unexpected response %+v", got)
	}
}

func TestDeleteOAuthClient(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/oauth/clients/42.json" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.DeleteOAuthClient(context.Background(), 42); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}
