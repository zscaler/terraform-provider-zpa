package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/scimattributeheader"
)

func dataSourceScimAttributeHeader() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScimAttributeHeaderRead,
		Schema: map[string]*schema.Schema{
			"canonical_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"case_sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"values": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceScimAttributeHeaderRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *scimattributeheader.ScimAttributeHeader
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
		res, _, err := zClient.scimattributeheader.Get(idpResp.ID, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		res, _, err := zClient.scimattributeheader.GetByName(name, idpResp.ID)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		values, _ := zClient.scimattributeheader.GetValues(resp.IdpID, resp.ID)
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
		_ = d.Set("values", values)

	} else {
		return fmt.Errorf("no scim attribute name '%s' & idp name '%s' OR id '%s' was found", name, idpName, id)
	}
	return nil
}
