package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/scimgroup"
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
	service := zClient.IDPController

	var resp *scimgroup.ScimGroup
	idpId, okidpId := d.Get("idp_id").(string)
	idpName, okIdpName := d.Get("idp_name").(string)
	if !okIdpName && !okidpId || idpId == "" && idpName == "" {
		log.Printf("[INFO] idp name or id is required\n")
		return fmt.Errorf("idp name or id is required")
	}
	var idpResp *idpcontroller.IdpController
	// getting Idp controller by id or name
	if idpId != "" {
		resp, _, err := idpcontroller.Get(service, idpId)
		if err != nil || resp == nil {
			log.Printf("[INFO] couldn't find idp by id: %s\n", idpId)
			return err
		}
		idpResp = resp
	} else {
		resp, _, err := idpcontroller.GetByName(service, idpName)
		if err != nil || resp == nil {
			log.Printf("[INFO] couldn't find idp by name: %s\n", idpName)
			return err
		}
		idpResp = resp
	}
	// getting scim attribute header by id or name
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		res, _, err := zClient.ScimGroup.Get(idpResp.ID)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		res, _, err := zClient.ScimGroup.GetByName(name, idpResp.ID)
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
