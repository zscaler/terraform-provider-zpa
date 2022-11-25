package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/bacertificate"
)

func resourceBaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaCertificateCreate,
		Read:   resourceBaCertificateRead,
		Update: resourceBaCertificateUpdate,
		Delete: resourceBaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.bacertificate.Get(id)
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cert_chain": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cert_blob": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"issued_by": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"issued_to": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"san": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"serial_no": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"valid_from_in_epochsec": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"valid_to_in_epochsec": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceBaCertificateCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandBaCertificate(d)
	log.Printf("[INFO] Creating certificate with request\n%+v\n", req)

	baCertificate, _, err := zClient.bacertificate.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created certificate request. ID: %v\n", baCertificate)

	d.SetId(baCertificate.ID)
	return resourceBaCertificateRead(d, m)

}

func resourceBaCertificateRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.bacertificate.Get(d.Id())
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
	_ = d.Set("cname", resp.CName)
	_ = d.Set("cert_blob", resp.CertBlob)
	_ = d.Set("cert_chain", resp.CertChain)
	_ = d.Set("description", resp.Description)
	_ = d.Set("issued_by", resp.IssuedBy)
	_ = d.Set("issued_to", resp.IssuedTo)
	_ = d.Set("name", resp.Name)
	_ = d.Set("san", resp.San)
	_ = d.Set("serial_no", resp.SerialNo)
	_ = d.Set("status", resp.Status)
	_ = d.Set("valid_from_in_epochsec", resp.ValidFromInEpochSec)
	_ = d.Set("valid_to_in_epochsec", resp.ValidToInEpochSec)
	return nil
}

func resourceBaCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating certificate ID: %v\n", id)
	req := expandBaCertificate(d)

	if _, err := zClient.bacertificate.Update(id, &req); err != nil {
		return err
	}

	return resourceBaCertificateRead(d, m)
}

func resourceBaCertificateDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting certificate ID: %v\n", d.Id())

	if _, err := zClient.bacertificate.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] certificate deleted")
	return nil
}

func expandBaCertificate(d *schema.ResourceData) bacertificate.BaCertificate {
	baCertificate := bacertificate.BaCertificate{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		CertChain:           d.Get("cert_chain").(string),
		CertBlob:            d.Get("cert_blob").(string),
		IssuedBy:            d.Get("issued_by").(string),
		IssuedTo:            d.Get("issued_to").(string),
		SerialNo:            d.Get("serial_no").(string),
		Status:              d.Get("status").(string),
		ValidFromInEpochSec: d.Get("valid_from_in_epochsec").(string),
		ValidToInEpochSec:   d.Get("valid_to_in_epochsec").(string),
	}
	return baCertificate
}
