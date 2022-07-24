package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/samlattribute"
)

func dataSourceSamlAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSamlAttributeRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"saml_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_attribute": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSamlAttributeRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *samlattribute.SamlAttribute
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for saml attribute %s\n", id)
		res, _, err := zClient.samlattribute.Get(id)
		if err != nil {
			return err
		}
		resp = res

	}
	name, ok := d.Get("name").(string)
	if ok && id == "" && name != "" {
		log.Printf("[INFO] Getting data for saml attribute %s\n", name)
		res, _, err := zClient.samlattribute.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("idp_name", resp.IdpName)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("saml_name", resp.SamlName)
		_ = d.Set("user_attribute", resp.UserAttribute)
	} else {
		return fmt.Errorf("couldn't find any saml attribute with name '%s' or id '%s'", name, id)
	}
	return nil
}
