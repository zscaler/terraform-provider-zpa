package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
)

func resourceBaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaCertificateCreate,
		Read:   resourceBaCertificateRead,
		Update: resourceBaCertificateUpdate,
		Delete: resourceBaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				service := m.(*Client).bacertificate.WithMicroTenant(GetString(d.Get("microtenant_id")))

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := service.GetIssuedByName(id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the certificate",
			},
			"cert_blob": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the certificate",
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate text in PEM format",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the certificate.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the certificate",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Microtenant",
			},
		},
	}
}

func resourceBaCertificateCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).bacertificate.WithMicroTenant(GetString(d.Get("microtenant_id")))

	req := expandBaCertificate(d)
	log.Printf("[INFO] Creating certificate with request\n%+v\n", req)

	baCertificate, _, err := service.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created certificate request. ID: %v\n", baCertificate)

	d.SetId(baCertificate.ID)
	return resourceBaCertificateRead(d, m)

}

func resourceBaCertificateRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).bacertificate.WithMicroTenant(GetString(d.Get("microtenant_id")))

	resp, _, err := service.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing certificate %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting certificate:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("certificate", resp.Certificate)
	_ = d.Set("microtenant_id", resp.MicrotenantID)
	return nil
}

func resourceBaCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	// Update doesn't actually do anything, because an certificates can't be updated.
	// This function is required by the Terraform framework
	return nil
}

func resourceBaCertificateDelete(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).bacertificate.WithMicroTenant(GetString(d.Get("microtenant_id")))

	log.Printf("[INFO] Deleting certificate ID: %v\n", d.Id())

	if _, err := service.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] certificate deleted")
	return nil
}

func expandBaCertificate(d *schema.ResourceData) bacertificate.BaCertificate {
	baCertificate := bacertificate.BaCertificate{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		CertBlob:      d.Get("cert_blob").(string),
		MicrotenantID: d.Get("microtenant_id").(string),
	}
	return baCertificate
}
