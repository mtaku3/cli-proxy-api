package configaccess

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdkaccess "github.com/router-for-me/CLIProxyAPI/v8/sdk/access"
	sdkconfig "github.com/router-for-me/CLIProxyAPI/v8/sdk/config"
)

func authenticateWithKeys(t *testing.T, keys []string, header string) (*sdkaccess.Result, *sdkaccess.AuthError) {
	t.Helper()
	sdkaccess.UnregisterProvider(sdkaccess.AccessProviderTypeConfigAPIKey)
	t.Cleanup(func() { sdkaccess.UnregisterProvider(sdkaccess.AccessProviderTypeConfigAPIKey) })

	Register(&sdkconfig.SDKConfig{APIKeys: keys})
	manager := sdkaccess.NewManager()
	manager.SetProviders(sdkaccess.RegisteredProviders())

	req := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	if header != "" {
		req.Header.Set("Authorization", "Bearer "+header)
	}
	return manager.Authenticate(context.Background(), req)
}

func TestRequireAPIKeyUnsetAllowsEmptyKeys(t *testing.T) {
	t.Setenv("REQUIRE_API_KEY", "")

	if _, err := authenticateWithKeys(t, nil, ""); err != nil {
		t.Fatalf("expected upstream behavior (no auth) without REQUIRE_API_KEY, got %v", err)
	}
}

func TestRequireAPIKeyRejectsWhenNoKeysConfigured(t *testing.T) {
	t.Setenv("REQUIRE_API_KEY", "true")

	for _, keys := range [][]string{nil, {}, {"", "  "}} {
		for _, header := range []string{"", "anything"} {
			_, err := authenticateWithKeys(t, keys, header)
			if err == nil {
				t.Fatalf("keys=%q header=%q: expected rejection, got success", keys, header)
			}
			if err.HTTPStatusCode() != http.StatusUnauthorized {
				t.Fatalf("keys=%q header=%q: expected 401, got %d", keys, header, err.HTTPStatusCode())
			}
		}
	}
}

func TestRequireAPIKeyKeepsConfiguredKeysWorking(t *testing.T) {
	t.Setenv("REQUIRE_API_KEY", "true")

	if _, err := authenticateWithKeys(t, []string{"sk-valid"}, "sk-valid"); err != nil {
		t.Fatalf("expected valid key to authenticate, got %v", err)
	}
	if _, err := authenticateWithKeys(t, []string{"sk-valid"}, "sk-wrong"); err == nil {
		t.Fatal("expected wrong key to be rejected")
	}
}

func TestRequireAPIKeyParsing(t *testing.T) {
	for value, want := range map[string]bool{"true": true, "1": true, " TRUE ": true, "false": false, "0": false, "yes": false, "": false} {
		t.Setenv("REQUIRE_API_KEY", value)
		if got := requireAPIKey(); got != want {
			t.Fatalf("REQUIRE_API_KEY=%q: got %v, want %v", value, got, want)
		}
	}
}
