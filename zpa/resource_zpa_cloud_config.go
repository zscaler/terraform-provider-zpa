package zpa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/custom_config_controller"
)

func resourceZiaCloudConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceZiaCloudConfigRead,
		CreateContext: resourceZiaCloudConfigCreate,
		UpdateContext: resourceZiaCloudConfigUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceZiaCloudConfigRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read zia cloud config import: %s", diags[0].Summary)
				}
				d.SetId("zia_cloud_config")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"zia_cloud_domain": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"zscaler",
					"zscloud",
					"zscalerone",
					"zscalertwo",
					"zscalerthree",
					"zscalerbeta",
					"zscalergov",
					"zscalerten",
					"zspreview",
				}, false),
				Description: "ZIA cloud domain (without .net suffix). Valid values: zscaler, zscloud, zscalerone, zscalertwo, zscalerthree, zscalerbeta, zscalergov, zscalerten, zspreview",
				StateFunc: func(val interface{}) string {
					domain := val.(string)
					// Ensure the domain ends with .net
					if domain != "" && !strings.HasSuffix(domain, ".net") {
						return domain + ".net"
					}
					return domain
				},
			},
			"zia_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zia_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "ZIA password (write-only, not returned by API)",
			},
			"zia_sandbox_api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "ZIA sandbox API token (write-only, not returned by API)",
			},
			"zia_cloud_service_api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "ZIA cloud service API key (write-only, not returned by API)",
			},
		},
	}
}

func resourceZiaCloudConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Normalize the domain by ensuring it has .net suffix
	domain := d.Get("zia_cloud_domain").(string)
	if !strings.HasSuffix(domain, ".net") {
		domain = domain + ".net"
	}

	cloudConfig := custom_config_controller.ZIACloudConfig{
		ZIACloudDomain:        domain,
		ZIAUsername:           d.Get("zia_username").(string),
		ZIAPassword:           d.Get("zia_password").(string),
		ZIASandboxApiToken:    d.Get("zia_sandbox_api_token").(string),
		ZIACloudServiceApiKey: d.Get("zia_cloud_service_api_key").(string),
	}

	_, _, err := custom_config_controller.AddZIACloudConfig(ctx, service, &cloudConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("zia_cloud_config")

	return resourceZiaCloudConfigRead(ctx, d, meta)
}

func resourceZiaCloudConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, _, err := custom_config_controller.GetZIACloudConfig(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("zia_cloud_config")
		_ = d.Set("zia_cloud_domain", resp.ZIACloudDomain)
		_ = d.Set("zia_username", resp.ZIAUsername)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't read ZIA cloud config"))
	}

	return nil
}

func resourceZiaCloudConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Normalize the domain by ensuring it has .net suffix
	domain := d.Get("zia_cloud_domain").(string)
	if !strings.HasSuffix(domain, ".net") {
		domain = domain + ".net"
	}

	cloudConfig := custom_config_controller.ZIACloudConfig{
		ZIACloudDomain:        domain,
		ZIAUsername:           d.Get("zia_username").(string),
		ZIAPassword:           d.Get("zia_password").(string),
		ZIASandboxApiToken:    d.Get("zia_sandbox_api_token").(string),
		ZIACloudServiceApiKey: d.Get("zia_cloud_service_api_key").(string),
	}

	_, _, err := custom_config_controller.AddZIACloudConfig(ctx, service, &cloudConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("zia_cloud_config")

	return resourceZiaCloudConfigRead(ctx, d, meta)
}
