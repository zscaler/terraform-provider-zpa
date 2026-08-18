package zpa

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// clearAuthEnv blanks every credential-bearing environment variable so the
// tests behave the same on developer machines and CI.
func clearAuthEnv(t *testing.T) {
	t.Helper()
	for _, v := range []string{
		"ZSCALER_CLIENT_ID", "ZSCALER_CLIENT_SECRET", "ZSCALER_PRIVATE_KEY",
		"ZSCALER_VANITY_DOMAIN", "ZSCALER_CLOUD", "ZPA_CUSTOMER_ID",
		"ZPA_CLIENT_ID", "ZPA_CLIENT_SECRET", "ZPA_CLOUD",
		"ZSCALER_USE_LEGACY_CLIENT", "ZSCALER_SKIP_CREDENTIALS_VALIDATION",
	} {
		t.Setenv(v, "")
	}
}

func TestProviderConfigureSkipCredentialsValidation(t *testing.T) {
	clearAuthEnv(t)

	d := schema.TestResourceDataRaw(t, ZPAProvider().Schema, map[string]interface{}{
		"skip_credentials_validation": true,
	})

	meta, diags := providerConfigure(d, "1.0-test")
	if diags.HasError() {
		t.Fatalf("expected no error in skip mode, got: %v", diags)
	}

	foundWarning := false
	for _, diagnostic := range diags {
		if strings.Contains(diagnostic.Summary, "credentials were not validated") {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Errorf("expected a warning diagnostic about skipped credential validation, got: %v", diags)
	}

	client, ok := meta.(*Client)
	if !ok {
		t.Fatalf("expected *Client meta, got %T", meta)
	}
	if !client.skipCredentialsValidation {
		t.Error("expected client to be marked as inert (skipCredentialsValidation)")
	}
	if client.Service != nil {
		t.Error("expected nil Service on the inert client")
	}
}

func TestProviderConfigureSkipCredentialsValidationEnvVar(t *testing.T) {
	clearAuthEnv(t)
	t.Setenv("ZSCALER_SKIP_CREDENTIALS_VALIDATION", "true")

	d := schema.TestResourceDataRaw(t, ZPAProvider().Schema, map[string]interface{}{})

	meta, diags := providerConfigure(d, "1.0-test")
	if diags.HasError() {
		t.Fatalf("expected no error in env-var skip mode, got: %v", diags)
	}
	client, ok := meta.(*Client)
	if !ok || !client.skipCredentialsValidation {
		t.Fatalf("expected inert client via env var, got %#v", meta)
	}
}

func TestProviderConfigureMissingCredentialsStillErrors(t *testing.T) {
	clearAuthEnv(t)

	d := schema.TestResourceDataRaw(t, ZPAProvider().Schema, map[string]interface{}{})

	_, diags := providerConfigure(d, "1.0-test")
	if !diags.HasError() {
		t.Fatal("expected configure to fail without credentials when skip mode is off")
	}
}

func TestInertClientGuardOnResources(t *testing.T) {
	p := ZPAProvider()
	inert := &Client{skipCredentialsValidation: true, policySetIDCache: make(map[string]string)}
	ctx := context.Background()

	r, ok := p.ResourcesMap["zpa_segment_group"]
	if !ok {
		t.Fatal("zpa_segment_group resource not found")
	}
	for name, f := range map[string]func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics{
		"create": r.CreateContext,
		"read":   r.ReadContext,
		"update": r.UpdateContext,
		"delete": r.DeleteContext,
	} {
		if f == nil {
			continue
		}
		diags := f(ctx, r.TestResourceData(), inert)
		if !diags.HasError() {
			t.Errorf("%s: expected guard error with inert client, got none", name)
		} else if !strings.Contains(diags[0].Summary, "skip_credentials_validation") {
			t.Errorf("%s: expected guard error to mention skip_credentials_validation, got: %s", name, diags[0].Summary)
		}
	}

	ds, ok := p.DataSourcesMap["zpa_segment_group"]
	if !ok {
		t.Fatal("zpa_segment_group data source not found")
	}
	diags := ds.ReadContext(ctx, ds.TestResourceData(), inert)
	if !diags.HasError() || !strings.Contains(diags[0].Summary, "skip_credentials_validation") {
		t.Errorf("data source read: expected guard error, got: %v", diags)
	}

	if r.Importer != nil && r.Importer.StateContext != nil {
		if _, err := r.Importer.StateContext(ctx, r.TestResourceData(), inert); err == nil {
			t.Error("import: expected guard error with inert client, got none")
		}
	}
}
