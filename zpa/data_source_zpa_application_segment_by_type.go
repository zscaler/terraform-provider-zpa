package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbytype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplicationSegmentByType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationSegmentByTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Browser Access application",
			},
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the application",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Browser Access application",
			},
			"application_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of application, BROWSER_ACCESS, INSPECT or SECURE_REMOTE_ACCESS",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Browser Access application is enabled or not",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain of the Browser Access application",
			},
			"application_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port for the Browser Access application",
			},
			"application_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol for the Browser Access application",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Browser Access certificate",
			},
			"certificate_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Browser Access certificate",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant",
			},
			"microtenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Microtenant",
			},
		},
	}
}

func dataSourceApplicationSegmentByTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationType := d.Get("application_type").(string)
	if applicationType != "BROWSER_ACCESS" && applicationType != "SECURE_REMOTE_ACCESS" && applicationType != "INSPECT" {
		return diag.FromErr(fmt.Errorf("invalid application_type '%s'. Valid types are 'BROWSER_ACCESS', 'SECURE_REMOTE_ACCESS', 'INSPECT'", applicationType))
	}

	name := d.Get("name").(string)
	log.Printf("[INFO] Getting data for application segment with type %s", applicationType)
	if name != "" {
		log.Printf("[INFO] Getting data for application segment with name %s and type %s", name, applicationType)
	}

	// Call the SDK function
	resp, _, err := applicationsegmentbytype.GetByApplicationType(ctx, service, name, applicationType, true)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resp) == 0 {
		return diag.FromErr(fmt.Errorf("no application segment found for name '%s' and type '%s'", name, applicationType))
	}

	// Assuming we are only interested in the first result for simplicity
	appSegment := resp[0]

	d.SetId(appSegment.ID)
	_ = d.Set("app_id", appSegment.AppID)
	_ = d.Set("name", appSegment.Name)
	_ = d.Set("enabled", appSegment.Enabled)
	_ = d.Set("domain", appSegment.Domain)
	_ = d.Set("application_port", appSegment.ApplicationPort)
	_ = d.Set("application_protocol", appSegment.ApplicationProtocol)
	_ = d.Set("certificate_id", appSegment.CertificateID)
	_ = d.Set("certificate_name", appSegment.CertificateName)
	_ = d.Set("microtenant_id", appSegment.MicroTenantID)
	_ = d.Set("microtenant_name", appSegment.MicroTenantName)

	return nil
}
