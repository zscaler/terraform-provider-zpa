package zpa

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// providerData builds the provider block configuration Configure receives.
func providerData(t *testing.T, raw map[string]any) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, ZPAProvider().Schema, raw)
}

// Provider block attributes win over environment variables, which fill in what
// the block leaves out; retries and timeouts have defaults.
func TestNewConfig(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("ZSCALER_CLIENT_ID", "env-client")
	t.Setenv("ZSCALER_CLIENT_SECRET", "env-secret")
	t.Setenv("ZSCALER_VANITY_DOMAIN", "acme")
	t.Setenv("ZPA_CUSTOMER_ID", "123")
	t.Setenv("ZSCALER_CLOUD", "beta")
	t.Setenv("ZPA_MICROTENANT_ID", "mt-env")
	t.Setenv("ZSCALER_HTTP_PROXY", "http://proxy.acme.com:3128")
	t.Setenv("ZSCALER_USE_LEGACY_CLIENT", "TRUE")
	t.Setenv("ZPA_CLIENT_ID", "legacy-id")
	t.Setenv("ZPA_CLIENT_SECRET", "legacy-secret")
	t.Setenv("ZPA_CLOUD", "BETA")

	c := NewConfig(providerData(t, map[string]any{"client_id": "hcl-client", "max_retries": 10}))
	checks := map[string][2]any{
		"client_id from HCL":      {c.clientID, "hcl-client"},
		"client_secret from env":  {c.clientSecret, "env-secret"},
		"vanity_domain from env":  {c.vanityDomain, "acme"},
		"customer_id from env":    {c.customerID, "123"},
		"zscaler_cloud from env":  {c.cloud, "beta"},
		"microtenant_id from env": {c.microtenantID, "mt-env"},
		"http_proxy from env":     {c.httpProxy, "http://proxy.acme.com:3128"},
		"legacy client from env":  {c.useLegacyClient, true},
		"legacy client id":        {c.zpaClientID, "legacy-id"},
		"legacy client secret":    {c.zpaClientSecret, "legacy-secret"},
		"legacy customer id":      {c.zpaCustomerID, "123"},
		"legacy cloud":            {c.BaseURL, "BETA"},
		"max_retries from HCL":    {c.retryCount, 10},
		"request_timeout default": {c.requestTimeout, 240},
	}
	for name, check := range checks {
		if check[0] != check[1] {
			t.Errorf("%s: got %v, want %v", name, check[0], check[1])
		}
	}
}

// Configuration errors are reported before any attempt to log in.
func TestProviderClientConfigurationErrors(t *testing.T) {
	valid := map[string]any{"client_id": "id", "client_secret": "secret", "vanity_domain": "acme", "customer_id": "123"}
	with := func(extra map[string]any) map[string]any {
		out := map[string]any{}
		for k, v := range valid {
			out[k] = v
		}
		for k, v := range extra {
			out[k] = v
		}
		return out
	}
	cases := map[string]struct {
		config  map[string]any
		wantErr string
	}{
		"missing credentials":   {map[string]any{"client_id": "id"}, "missing required parameters"},
		"proxy port not a port": {with(map[string]any{"http_proxy": "http://proxy.acme.com:99999"}), "invalid port number"},
		"proxy url invalid":     {with(map[string]any{"http_proxy": "http://proxy.acme.com:abc"}), "invalid port"},
		"legacy proxy invalid":  {with(map[string]any{"use_legacy_client": true, "http_proxy": "http://proxy.acme.com:0"}), "invalid port number"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			clearAuthEnv(t)
			_, diags := providerConfigure(providerData(t, tc.config), "1.0-test")
			if !diags.HasError() || !strings.Contains(diags[0].Summary, tc.wantErr) {
				t.Fatalf("got %v, want an error containing %q", diags, tc.wantErr)
			}
		})
	}
}

// The client built (and authenticated) by loadClients is reused, not built again.
func TestConfigClientReusesLoadedClient(t *testing.T) {
	api := common.NewZPATest(t)
	c := &Config{zscalerSDKClientV3: api.Service.Client}
	client, err := c.Client()
	if err != nil {
		t.Fatalf("Client: %v", err)
	}
	if client.Service.Client != api.Service.Client {
		t.Error("Client() built a new SDK client instead of reusing the loaded one")
	}
}

func TestGenerateUserAgent(t *testing.T) {
	ua := generateUserAgent("1.9.0", "123")
	for _, part := range []string{"Terraform/1.9.0", "Provider/", "Customer/123"} {
		if !strings.Contains(ua, part) {
			t.Errorf("user agent %q is missing %q", ua, part)
		}
	}
}
