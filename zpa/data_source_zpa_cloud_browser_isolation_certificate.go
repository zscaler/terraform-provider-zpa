package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func dataSourceCBICertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCBICertificatesRead,
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

func dataSourceCBICertificatesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.CBICertificateController

	var resp *cbicertificatecontroller.CBICertificate
	id, idOk := d.Get("id").(string)
	name, nameOk := d.Get("name").(string)

	if idOk && id != "" {
		log.Printf("[INFO] Getting data for CBI certificate with ID: %s\n", id)
		res, _, err := cbicertificatecontroller.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	} else if nameOk && name != "" {
		log.Printf("[INFO] Getting data for CBI certificate with name: %s\n", name)
		res, _, err := cbicertificatecontroller.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("pem", resp.PEM)
		_ = d.Set("is_default", resp.IsDefault)
	} else {
		return fmt.Errorf("couldn't find any CBI certificate with name '%s' or id '%s'", name, id)
	}

	return nil
}
