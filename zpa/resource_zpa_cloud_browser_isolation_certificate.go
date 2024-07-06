package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func resourceCBICertificates() *schema.Resource {
	return &schema.Resource{
		Create: resourceCBICertificatesCreate,
		Read:   resourceCBICertificatesRead,
		Update: resourceCBICertificatesUpdate,
		Delete: resourceCBICertificatesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.CBICertificateController

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := cbicertificatecontroller.GetByNameOrID(service, id)
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pem": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCBICertificatesCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.CBICertificateController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandCBICertificate(d)
	log.Printf("[INFO] Creating cbi certificate with request\n%+v\n", req)

	cbiCertificate, _, err := cbicertificatecontroller.Create(service, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created cbi certificate request. ID: %v\n", cbiCertificate)

	d.SetId(cbiCertificate.ID)
	return resourceCBICertificatesRead(d, meta)
}

func resourceCBICertificatesRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.CBICertificateController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := cbicertificatecontroller.Get(service, d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing cbi certificate %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting cbi certificate:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("pem", resp.PEM)

	return nil
}

func resourceCBICertificatesUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.CBICertificateController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating cbi certificate ID: %v\n", id)
	req := expandCBICertificate(d)

	if _, _, err := cbicertificatecontroller.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := cbicertificatecontroller.Update(service, id, &req); err != nil {
		return err
	}

	return resourceCBICertificatesRead(d, meta)
}

func resourceCBICertificatesDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.CBICertificateController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting cbi certificate ID: %v\n", d.Id())

	if _, err := cbicertificatecontroller.Delete(service, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] cbi certificate deleted")
	return nil
}

func expandCBICertificate(d *schema.ResourceData) cbicertificatecontroller.CBICertificate {
	cbiCertificate := cbicertificatecontroller.CBICertificate{
		ID:   d.Id(),
		Name: d.Get("name").(string),
		PEM:  d.Get("pem").(string),
	}
	return cbiCertificate
}
