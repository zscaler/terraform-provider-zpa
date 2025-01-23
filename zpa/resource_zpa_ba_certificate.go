package zpa

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

func resourceBaCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBaCertificateCreate,
		ReadContext:   resourceBaCertificateRead,
		UpdateContext: resourceFuncNoOp,
		DeleteContext: resourceBaCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func resourceBaCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandBaCertificate(d)
	log.Printf("[INFO] Creating certificate with request\n%+v\n", req)

	baCertificate, _, err := bacertificate.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created certificate request. ID: %v\n", baCertificate)

	d.SetId(baCertificate.ID)
	return nil
}

func resourceBaCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := bacertificate.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing ba certificate %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting ba certificate:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("certificate", resp.Certificate)
	_ = d.Set("microtenant_id", resp.MicrotenantID)
	return nil
}

func resourceBaCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting certificate ID: %v\n", d.Id())

	if _, err := bacertificate.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
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
