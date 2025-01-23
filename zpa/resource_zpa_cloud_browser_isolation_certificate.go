package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func resourceCBICertificates() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCBICertificatesCreate,
		ReadContext:   resourceCBICertificatesRead,
		UpdateContext: resourceCBICertificatesUpdate,
		DeleteContext: resourceCBICertificatesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := cbicertificatecontroller.GetByNameOrID(ctx, service, id)
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

func resourceCBICertificatesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandCBICertificate(d)
	log.Printf("[INFO] Creating cbi certificate with request\n%+v\n", req)

	cbiCertificate, _, err := cbicertificatecontroller.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created cbi certificate request. ID: %v\n", cbiCertificate)

	d.SetId(cbiCertificate.ID)
	return resourceCBICertificatesRead(ctx, d, meta)
}

func resourceCBICertificatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := cbicertificatecontroller.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing cbi certificate %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting cbi certificate:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("pem", resp.PEM)

	return nil
}

func resourceCBICertificatesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating cbi certificate ID: %v\n", id)
	req := expandCBICertificate(d)

	if _, _, err := cbicertificatecontroller.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := cbicertificatecontroller.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceCBICertificatesRead(ctx, d, meta)
}

func resourceCBICertificatesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting cbi certificate ID: %v\n", d.Id())

	if _, err := cbicertificatecontroller.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
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
