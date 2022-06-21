package zpa

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"zpa": testAccProvider,
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"zpa": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = Provider()
}

func testAccPreCheck(t *testing.T) {
	err := accPreCheck()
	if err != nil {
		t.Fatalf("%v", err)
	}
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

func accPreCheck() error {
	if v := os.Getenv("ZPA_CLIENT_ID"); v == "" {
		return errors.New("ZPA_CLIENT_ID must be set for acceptance tests")
	}
	client_id := os.Getenv("ZPA_CLIENT_ID")
	client_secret := os.Getenv("ZPA_CLIENT_SECRET")
	customer_id := os.Getenv("ZPA_CUSTOMER_ID")
	if client_id == "" && (client_id == "" || client_secret == "" || customer_id == "") {
		return errors.New("either ZPA_CLIENT_ID or ZPA_CLIENT_SECRET, and ZPA_CUSTOMER_ID must be set for acceptance tests")
	}
	return nil
}
