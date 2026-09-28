package zendesk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nukosuke/go-zendesk/zendesk"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func newOAuthTestClient(t *testing.T, response string) *newClient.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	zd, err := zendesk.NewClient(nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	if err := zd.SetEndpointURL(server.URL); err != nil {
		t.Fatalf("could not set endpoint: %v", err)
	}

	return &newClient.Client{Client: *zd}
}

func TestUnmarshalOAuthClient(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "42",
		mapGetterSetter: mapGetterSetter{
			"name":                     "Example",
			"identifier":               "example",
			"kind":                     "confidential",
			"scopes":                   []interface{}{"tickets:read", "users:write"},
			"enforce_token_expiration": false,
		},
	}

	c, err := unmarshalOAuthClient(m)
	if err != nil {
		t.Fatalf("could not unmarshal: %v", err)
	}

	if c.ID != 42 || c.Name != "Example" || c.Identifier != "example" || c.Kind != "confidential" {
		t.Fatalf("unexpected oauth client %+v", c)
	}
	if len(c.Scopes) != 2 || c.Scopes[1] != "users:write" {
		t.Fatalf("unexpected scopes %v", c.Scopes)
	}
	if c.EnforceTokenExpiration == nil || *c.EnforceTokenExpiration {
		t.Fatalf("enforce_token_expiration should be an explicit false")
	}
}

func TestCreateOAuthClientStoresSecret(t *testing.T) {
	zd := newOAuthTestClient(t, `{"client":{"id":42,"name":"Example","secret":"full-secret"}}`)
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{
			"name":                     "Example",
			"identifier":               "example",
			"kind":                     "confidential",
			"enforce_token_expiration": true,
		},
	}

	if diags := createOAuthClient(context.Background(), m, zd); diags.HasError() {
		t.Fatalf("create failed: %v", diags)
	}

	if m.Id() != "42" {
		t.Fatalf("unexpected id %q", m.Id())
	}
	if v := m.Get("secret"); v != "full-secret" {
		t.Fatalf("secret from the create response was not stored, got %v", v)
	}
}

func TestReadOAuthClientKeepsSecret(t *testing.T) {
	// Reads only ever return a truncated stub, which must not replace the secret captured on create.
	zd := newOAuthTestClient(t, `{"client":{"id":42,"name":"Example","secret":"stub"}}`)
	m := &identifiableMapGetterSetter{
		id: "42",
		mapGetterSetter: mapGetterSetter{
			"secret": "full-secret",
		},
	}

	if diags := readOAuthClient(context.Background(), m, zd); diags.HasError() {
		t.Fatalf("read failed: %v", diags)
	}

	if v := m.Get("secret"); v != "full-secret" {
		t.Fatalf("read overwrote the secret with %v", v)
	}
	if v := m.Get("name"); v != "Example" {
		t.Fatalf("read did not refresh name, got %v", v)
	}
}
