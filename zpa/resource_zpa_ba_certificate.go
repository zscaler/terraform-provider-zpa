package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
)

func resourceBaCertificate() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBaCertificateCreate,
		Read:          resourceBaCertificateRead,
		UpdateContext: resourceFuncNoOp,
		Delete:        resourceBaCertificateDelete,
		Importer:      nil,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the certificate",
			},
			"cert_blob": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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
	zClient := m.(*Client)
	service := zClient.BACertificate

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandBaCertificate(d)
	log.Printf("[INFO] Creating certificate with request\n%+v\n", req)

	baCertificate, _, err := bacertificate.Create(service, req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created certificate request. ID: %v\n", baCertificate)

	d.SetId(baCertificate.ID)
	return nil
}

func resourceBaCertificateRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.BACertificate

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := bacertificate.Get(service, d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing ba certificate %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting ba certificate:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("certificate", resp.Certificate)
	_ = d.Set("microtenant_id", resp.MicrotenantID)
	return nil
}

func resourceBaCertificateDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.BACertificate

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting certificate ID: %v\n", d.Id())

	if _, err := bacertificate.Delete(service, d.Id()); err != nil {
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
