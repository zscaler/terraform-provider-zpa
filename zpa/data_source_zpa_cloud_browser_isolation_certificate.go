package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func dataSourceCBICertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBICertificatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pem": {
				Type:     schema.TypeString,
				Computed: true,
				// Sensitive: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCBICertificatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cbicertificatecontroller.CBICertificate
	id, idOk := d.Get("id").(string)
	name, nameOk := d.Get("name").(string)

	if idOk && id != "" {
		log.Printf("[INFO] Getting data for CBI certificate with ID: %s\n", id)
		res, _, err := cbicertificatecontroller.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	} else if nameOk && name != "" {
		log.Printf("[INFO] Getting data for CBI certificate with name: %s\n", name)
		res, _, err := cbicertificatecontroller.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	} else {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be specified"))
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("pem", resp.PEM)
		_ = d.Set("is_default", resp.IsDefault)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any CBI certificate with name '%s' or id '%s'", name, id))
	}

	return nil
}
