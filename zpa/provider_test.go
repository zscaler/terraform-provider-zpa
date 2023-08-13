package zpa

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProvider          *schema.Provider
	testAccProviders         map[string]*schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
)

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

func testAccPreCheck(t *testing.T) {
	if err := accPreCheck(); err != nil {
		t.Fatal(err)
	}
}

func accPreCheck() error {
	envVars := map[string]string{
		"ZPA_CLIENT_ID":     os.Getenv("ZPA_CLIENT_ID"),
		"ZPA_CLIENT_SECRET": os.Getenv("ZPA_CLIENT_SECRET"),
		"ZPA_CUSTOMER_ID":   os.Getenv("ZPA_CUSTOMER_ID"),
	}

	for key, value := range envVars {
		if value == "" {
			return errors.New(key + " must be set for acceptance tests")
		}
	}

	return nil
}
