package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/scimattributeheader"
)

func dataSourceScimUserAttributeHeader() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScimUserAttributeHeaderRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
		},
	}
}

func dataSourceScimUserAttributeHeaderRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *scimattributeheader.ScimAttributeHeader
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		res, _, err := zClient.scimuserattributeheader.GetAll(id)
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
		res, _, err := zClient.scimuserattributeheader.GetByName(name, idpResp.ID)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("list", resp.List)

	} else {
		return fmt.Errorf("no scim attribute name '%s' & idp name '%s' OR id '%s' was found", name, idpName, id)
	}
	return nil
}
