package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/idpcontroller"
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
	idpId, okidpId := d.Get("idp_id").(string)
	idpName, okIdpName := d.Get("idp_name").(string)
	if !okIdpName && !okidpId || idpId == "" && idpName == "" {
		log.Printf("[INFO] idp name or id is required\n")
		return fmt.Errorf("idp name or id is required")
	}
	var idpResp *idpcontroller.IdpController
	// getting Idp controller by id or name
	if idpId != "" {
		resp, _, err := zClient.idpcontroller.Get(idpId)
		if err != nil || resp == nil {
			log.Printf("[INFO] couldn't find idp by id: %s\n", idpId)
			return err
		}
		idpResp = resp
	} else {
		resp, _, err := zClient.idpcontroller.GetByName(idpName)
		if err != nil || resp == nil {
			log.Printf("[INFO] couldn't find idp by name: %s\n", idpName)
			return err
		}
		idpResp = resp
	}
	// getting scim attribute header by id or name
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		res, _, err := zClient.samlattribute.Get(idpResp.ID)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
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
