// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	terraformVersionCache string
	cacheMu               sync.RWMutex
)

const defaultTerraformVersion = "0.11+compatible"

func (p *ZPAProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ZPAProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	terraformVersion := determineTerraformVersion(req.TerraformVersion)

	config := &client.Config{
		ClientID:         getStringValue(data.ClientID, "ZSCALER_CLIENT_ID"),
		ClientSecret:     getStringValue(data.ClientSecret, "ZSCALER_CLIENT_SECRET"),
		PrivateKey:       getStringValue(data.PrivateKey, "ZSCALER_PRIVATE_KEY"),
		VanityDomain:     getStringValue(data.VanityDomain, "ZSCALER_VANITY_DOMAIN"),
		Cloud:            getStringValue(data.ZscalerCloud, "ZSCALER_CLOUD"),
		CustomerID:       getStringValue(data.CustomerID, "ZPA_CUSTOMER_ID"),
		MicrotenantID:    getStringValue(data.MicrotenantID, "ZPA_MICROTENANT_ID"),
		ZPAClientID:      getStringValue(data.ZPAClientID, "ZPA_CLIENT_ID"),
		ZPAClientSecret:  getStringValue(data.ZPAClientSecret, "ZPA_CLIENT_SECRET"),
		ZPACustomerID:    getStringValue(data.ZPACustomerID, "ZPA_CUSTOMER_ID"),
		ZPACloud:         getStringValue(data.ZPACloud, "ZPA_CLOUD"),
		UseLegacyClient:  getBoolValue(data.UseLegacyClient, "ZSCALER_USE_LEGACY_CLIENT"),
		HTTPProxy:        getStringValue(data.HTTPProxy, "ZSCALER_HTTP_PROXY"),
		RetryCount:       getIntValue(data.MaxRetries, 100),
		Parallelism:      getIntValue(data.Parallelism, 1),
		RequestTimeout:   getIntValue(data.RequestTimeout, 240),
		MinWait:          getIntValue(data.MinWaitSeconds, 2),
		MaxWait:          getIntValue(data.MaxWaitSeconds, 10),
		TerraformVersion: terraformVersion,
		ProviderVersion:  p.version,
	}

	zpaClient, err := client.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create ZPA client",
			fmt.Sprintf("Unable to create ZPA client: %s", err.Error()),
		)
		return
	}

	resp.ResourceData = zpaClient
	resp.DataSourceData = zpaClient
	resp.EphemeralResourceData = zpaClient

	tflog.Info(ctx, "ZPA provider configured successfully")
}

// Helper functions to get values from config or environment
func getStringValue(value types.String, envVar string) string {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueString()
	}
	return os.Getenv(envVar)
}

func getBoolValue(value types.Bool, envVar string) bool {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueBool()
	}
	envValue := os.Getenv(envVar)
	return strings.ToLower(envValue) == "true"
}

func getIntValue(value types.Int64, defaultValue int) int {
	if !value.IsNull() && !value.IsUnknown() {
		return int(value.ValueInt64())
	}
	return defaultValue
}

func determineTerraformVersion(provided string) string {
	candidates := []string{
		strings.TrimSpace(provided),
		strings.TrimSpace(os.Getenv("TF_CLI_VERSION")),
		strings.TrimSpace(os.Getenv("TF_ACC_TERRAFORM_VERSION")),
	}

	for _, candidate := range candidates {
		if candidate != "" {
			setTerraformVersionCache(candidate)
			return candidate
		}
	}

	if detected := detectTerraformVersionFromBinary(); detected != "" {
		setTerraformVersionCache(detected)
		return detected
	}

	cacheMu.RLock()
	cached := terraformVersionCache
	cacheMu.RUnlock()

	if cached != "" {
		return cached
	}

	setTerraformVersionCache(defaultTerraformVersion)
	return defaultTerraformVersion
}

func setTerraformVersionCache(value string) {
	cacheMu.Lock()
	terraformVersionCache = value
	cacheMu.Unlock()
}

func detectTerraformVersionFromBinary() string {
	paths := []string{
		strings.TrimSpace(os.Getenv("TF_ACC_TERRAFORM_PATH")),
		"terraform",
	}

	for _, path := range paths {
		if path == "" {
			continue
		}

		output, err := exec.Command(path, "-version").Output()
		if err != nil {
			continue
		}

		if version := parseTerraformVersion(string(output)); version != "" {
			return version
		}
	}

	return ""
}

func parseTerraformVersion(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Terraform ") {
			version := strings.TrimPrefix(line, "Terraform ")
			version = strings.TrimPrefix(version, "v")
			version = strings.Fields(version)[0]
			return version
		}
	}
	return ""
}
