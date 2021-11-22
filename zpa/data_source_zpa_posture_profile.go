package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/postureprofile"
)

func postureProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"master_customer_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"posture_udid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"zscaler_cloud": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"zscaler_customer_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"modified_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"modifiedby": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourcePostureProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePostureProfileRead,
		Schema: MergeSchema(postureProfileSchema(),
			map[string]*schema.Schema{
				"list": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: postureProfileSchema(),
					},
				},
			}),
	}
}

func dataSourcePostureProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *postureprofile.PostureProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for posture profile %s\n", id)
		res, _, err := zClient.postureprofile.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for posture profile name %s\n", name)
		res, _, err := zClient.postureprofile.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("domain", resp.Domain)
		_ = d.Set("master_customer_id", resp.MasterCustomerID)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("posture_udid", resp.PostureudID)
		_ = d.Set("zscaler_cloud", resp.ZscalerCloud)
		_ = d.Set("zscaler_customer_id", resp.ZscalerCustomerID)
		_ = d.Set("list", flattenPostureProfileList([]postureprofile.PostureProfile{*resp}))
	} else if id != "" || name != "" {
		return fmt.Errorf("couldn't find any posture profile with name '%s' or id '%s'", name, id)
	} else {
		// get the list
		list, _, err := zClient.postureprofile.GetAll()
		if err != nil {
			return err
		}
		d.SetId("posture-profile-list")
		_ = d.Set("list", flattenPostureProfileList(list))
	}

	return nil
}

func flattenPostureProfileList(list []postureprofile.PostureProfile) []interface{} {
	keys := make([]interface{}, len(list))
	for i, item := range list {
		keys[i] = map[string]interface{}{
			"id":                  item.ID,
			"creation_time":       item.CreationTime,
			"domain":              item.Domain,
			"master_customer_id":  item.MasterCustomerID,
			"modifiedby":          item.ModifiedBy,
			"modified_time":       item.ModifiedTime,
			"name":                item.Name,
			"posture_udid":        item.PostureudID,
			"zscaler_cloud":       item.ZscalerCloud,
			"zscaler_customer_id": item.ZscalerCustomerID,
		}
	}
	return keys
}
