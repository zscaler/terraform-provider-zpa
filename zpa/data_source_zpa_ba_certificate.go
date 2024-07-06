package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
)

func dataSourceBaCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaCertificateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_chain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issued_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issued_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"san": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"serial_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_from_in_epochsec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_to_in_epochsec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.BACertificate

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	var resp *bacertificate.BaCertificate
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for browser certificate %s\n", id)
		res, _, err := bacertificate.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for browser certificate name %s\n", name)
		res, _, err := bacertificate.GetIssuedByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("cname", resp.CName)
		_ = d.Set("cert_chain", resp.CertChain)
		_ = d.Set("certificate", resp.Certificate)
		_ = d.Set("public_key", resp.PublicKey)
		_ = d.Set("issued_by", resp.IssuedBy)
		_ = d.Set("issued_to", resp.IssuedTo)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("san", resp.San)
		_ = d.Set("serial_no", resp.SerialNo)
		_ = d.Set("status", resp.Status)
		_ = d.Set("microtenant_id", resp.MicrotenantID)

		// Convert and set epoch attributes
		creationTime, err := epochToRFC1123(resp.CreationTime, false)
		if err != nil {
			return fmt.Errorf("error formatting creation_time: %s", err)
		}
		if err := d.Set("creation_time", creationTime); err != nil {
			return fmt.Errorf("error setting creation_time: %s", err)
		}

		modifiedTime, err := epochToRFC1123(resp.ModifiedTime, false)
		if err != nil {
			return fmt.Errorf("error formatting modified_time: %s", err)
		}
		if err := d.Set("modified_time", modifiedTime); err != nil {
			return fmt.Errorf("error setting modified_time: %s", err)
		}

		validFromInEpochSec, err := epochToRFC1123(resp.ValidFromInEpochSec, false)
		if err != nil {
			return fmt.Errorf("error formatting valid_from_in_epochsec: %s", err)
		}
		if err := d.Set("valid_from_in_epochsec", validFromInEpochSec); err != nil {
			return fmt.Errorf("error setting valid_from_in_epochsec: %s", err)
		}

		validToInEpochSec, err := epochToRFC1123(resp.ValidToInEpochSec, false)
		if err != nil {
			return fmt.Errorf("error formatting valid_to_in_epochsec: %s", err)
		}
		if err := d.Set("valid_to_in_epochsec", validToInEpochSec); err != nil {
			return fmt.Errorf("error setting valid_to_in_epochsec: %s", err)
		}

	} else {
		return fmt.Errorf("couldn't find any browser certificate with name '%s' or id '%s'", name, id)
	}

	return nil
}
