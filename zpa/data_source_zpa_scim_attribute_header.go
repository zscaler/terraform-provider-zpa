package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/scimattributeheader"
)

func dataSourceScimAttributeHeader() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScimAttributeHeaderRead,
		Schema: map[string]*schema.Schema{
			"canonical_values": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"case_sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"idp_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modifiedby": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"multivalued": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mutability": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"returned": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uniqueness": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceScimAttributeHeaderRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *scimattributeheader.ScimAttributeHeader
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		res, _, err := zClient.scimattributeheader.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	idpName, ok := d.Get("idp_name").(string)
	name, ok2 := d.Get("name").(string)
	if id == "" && ok && ok2 && idpName != "" && name != "" {
		idpResp, _, err := zClient.idpcontroller.GetByName(idpName)
		if err != nil || idpResp == nil {
			log.Printf("[INFO] couldn't find idp by name: %s\n", idpName)
			return err
		}
		res, _, err := zClient.scimattributeheader.GetByName(name, idpResp.ID)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("canonical_values", resp.CanonicalValues)
		_ = d.Set("case_sensitive", resp.CaseSensitive)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("data_type", resp.DataType)
		_ = d.Set("description", resp.Description)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("multivalued", resp.MultiValued)
		_ = d.Set("mutability", resp.Mutability)
		_ = d.Set("name", resp.Name)
		_ = d.Set("required", resp.Required)
		_ = d.Set("returned", resp.Returned)
		_ = d.Set("schema_uri", resp.SchemaURI)
		_ = d.Set("uniqueness", resp.Uniqueness)

	} else {
		return fmt.Errorf("no scim attribute name '%s' & idp name '%s' OR id '%s' was found", name, idpName, id)
	}
	return nil
}
