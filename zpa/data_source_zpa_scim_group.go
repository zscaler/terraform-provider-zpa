package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/scimgroup"
)

func dataSourceScimGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScimGroupRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"idp_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceScimGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	var resp *scimgroup.ScimGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		res, _, err := zClient.scimgroup.Get(id)
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
		res, _, err := zClient.scimgroup.GetByName(name, idpResp.ID)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(strconv.FormatInt(int64(resp.ID), 10))
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("idp_group_id", resp.IdpGroupID)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("idp_name", resp.IdpName)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
	} else {
		return fmt.Errorf("no scim name '%s' & idp name '%s' OR id '%s' was found", name, idpName, id)
	}
	return nil
}
