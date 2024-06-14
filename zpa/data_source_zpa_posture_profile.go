package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/postureprofile"
)

func dataSourcePostureProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePostureProfileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
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
		},
	}
}

func dataSourcePostureProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.PostureProfile

	var resp *postureprofile.PostureProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for posture profile %s\n", id)
		res, _, err := postureprofile.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for posture profile name %s\n", name)
		res, _, err := postureprofile.GetByName(service, name)
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

	} else {
		return fmt.Errorf("couldn't find any posture profile with name '%s' or id '%s'", name, id)
	}

	return nil
}
