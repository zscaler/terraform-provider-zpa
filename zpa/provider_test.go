package zpa

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var testAccProvidersVersionValidation map[string]*schema.Provider
var testAccProviderVersionValidation *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"zscaler": testAccProvider,
	}

	testAccProviderVersionValidation = Provider()
	//testAccProviderVersionValidation.ConfigureFunc = zscalerConfigureWithoutVersionValidation
	testAccProvidersVersionValidation = map[string]*schema.Provider{
		"zscaler": testAccProviderVersionValidation,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ZPA_CLIENT_ID"); v == "" {
		t.Fatal("ZPA_CLIENT_ID must be set for acceptance tests.")
	}
	if v := os.Getenv("ZPA_CLIENT_SECRET"); v == "" {
		t.Fatal("ZPA_CLIENT_SECRET must be set for acceptance tests.")
	}
	if v := os.Getenv("ZPA_CUSTOMER_ID"); v == "" {
		t.Fatal("ZPA_CUSTOMER_ID must be set for acceptance tests.")
	}
}
