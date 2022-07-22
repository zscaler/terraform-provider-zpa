package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/segmentgroup"
)

func resourceSegmentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegmentGroupCreate,
		Read:   resourceSegmentGroupRead,
		Update: resourceSegmentGroupUpdate,
		Delete: resourceSegmentGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.segmentgroup.GetByName(id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"applications": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the app group.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether this app group is enabled or not.",
				Optional:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the app group.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"policy_migrated": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tcp_keep_alive_enabled": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  0,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
		},
	}
}

func resourceSegmentGroupCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandSegmentGroup(d)
	log.Printf("[INFO] Creating segment group with request\n%+v\n", req)

	segmentgroup, _, err := zClient.segmentgroup.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created segment group request. ID: %v\n", segmentgroup)

	d.SetId(segmentgroup.ID)
	return resourceSegmentGroupRead(d, m)

}

func resourceSegmentGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.segmentgroup.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing segment group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting segment group:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("config_space", resp.ConfigSpace)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	//_ = d.Set("policy_migrated", resp.PolicyMigrated)
	_ = d.Set("tcp_keep_alive_enabled", resp.TcpKeepAliveEnabled)
	if err := d.Set("applications", flattenSegmentGroupApplicationsSimple(resp)); err != nil {
		return fmt.Errorf("failed to read applications %s", err)
	}
	return nil
}

func flattenSegmentGroupApplicationsSimple(segmentGroup *segmentgroup.SegmentGroup) []interface{} {
	segmentGroupApplications := make([]interface{}, len(segmentGroup.Applications))
	for i, segmentGroupApplication := range segmentGroup.Applications {
		segmentGroupApplications[i] = map[string]interface{}{
			"id": segmentGroupApplication.ID,
		}
	}

	return segmentGroupApplications
}
func resourceSegmentGroupUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating segment group ID: %v\n", id)
	req := expandSegmentGroup(d)

	if _, err := zClient.segmentgroup.Update(id, &req); err != nil {
		return err
	}

	return resourceSegmentGroupRead(d, m)
}

func resourceSegmentGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting segment group ID: %v\n", d.Id())

	if _, err := zClient.segmentgroup.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] segment group deleted")
	return nil
}

func expandSegmentGroup(d *schema.ResourceData) segmentgroup.SegmentGroup {
	segmentGroup := segmentgroup.SegmentGroup{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		PolicyMigrated:      d.Get("policy_migrated").(bool),
		ConfigSpace:         d.Get("config_space").(string),
		TcpKeepAliveEnabled: d.Get("tcp_keep_alive_enabled").(string),
		Applications:        expandSegmentGroupApplications(d.Get("applications").([]interface{})),
	}
	return segmentGroup
}

func expandSegmentGroupApplications(segmentGroupApplication []interface{}) []segmentgroup.Application {
	segmentGroupApplications := make([]segmentgroup.Application, len(segmentGroupApplication))

	for i, segmentGroupApp := range segmentGroupApplication {
		segmentGroupItem := segmentGroupApp.(map[string]interface{})
		segmentGroupApplications[i] = segmentgroup.Application{
			ID: segmentGroupItem["id"].(string),
		}

	}

	return segmentGroupApplications
}
