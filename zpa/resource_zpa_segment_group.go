package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func resourceSegmentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegmentGroupCreate,
		Read:   resourceSegmentGroupRead,
		Update: resourceSegmentGroupUpdate,
		Delete: resourceSegmentGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				service := m.(*Client).segmentgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := service.GetByName(id)
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
				Type:        schema.TypeString,
				Description: "Name of the app group.",
				Required:    true,
			},
			"policy_migrated": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "The `policy_migrated` field is now deprecated for the resource `zpa_segment_group`, please remove this attribute to prevent configuration drifts",
			},
			"tcp_keep_alive_enabled": {
				Type:       schema.TypeString,
				Optional:   true,
				Default:    "1",
				Deprecated: "The `tcp_keep_alive_enabled` field is now deprecated for the resource `zpa_segment_group`, please replace all uses of this within the `zpa_application_segment`resources with the attribute `tcp_keep_alive`",
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSegmentGroupCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).segmentgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	req := expandSegmentGroup(d)
	log.Printf("[INFO] Creating segment group with request\n%+v\n", req)

	segmentgroup, _, err := service.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created segment group request. ID: %v\n", segmentgroup)

	d.SetId(segmentgroup.ID)
	return resourceSegmentGroupRead(d, m)
}

func resourceSegmentGroupRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).segmentgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	resp, _, err := service.Get(d.Id())
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
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	_ = d.Set("policy_migrated", resp.PolicyMigrated)
	_ = d.Set("tcp_keep_alive_enabled", resp.TcpKeepAliveEnabled)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
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
	service := m.(*Client).segmentgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))

	id := d.Id()
	log.Printf("[INFO] Updating segment group ID: %v\n", id)
	req := expandSegmentGroup(d)

	if _, _, err := service.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := service.Update(id, &req); err != nil {
		return err
	}

	return resourceSegmentGroupRead(d, m)
}

func detachSegmentGroupFromAllPolicyRules(id string, policySetControllerService *policysetcontroller.Service) {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()
	var rules []policysetcontroller.PolicyRule
	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}
	for _, t := range types {
		policySet, _, err := policySetControllerService.GetByPolicyType(t)
		if err != nil {
			continue
		}
		r, _, err := policySetControllerService.GetAllByType(t)
		if err != nil {
			continue
		}
		for _, rule := range r {
			rule.PolicySetID = policySet.ID
			rules = append(rules, rule)
		}
	}
	for _, rule := range rules {
		changed := false
		for i, condition := range rule.Conditions {
			operands := []policysetcontroller.Operands{}
			for _, op := range condition.Operands {
				if op.ObjectType == "APP_GROUP" && op.LHS == "id" && op.RHS == id {
					changed = true
					continue
				}
				operands = append(operands, op)
			}
			rule.Conditions[i].Operands = operands
		}
		if len(rule.Conditions) == 0 {
			rule.Conditions = []policysetcontroller.Conditions{}
		}
		if changed {
			if _, err := policySetControllerService.UpdateRule(rule.PolicySetID, rule.ID, &rule); err != nil {
				continue
			}
		}
	}
}

func resourceSegmentGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	policySetControllerService := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	service := zClient.segmentgroup.WithMicroTenant(GetString(d.Get("microtenant_id")))
	log.Printf("[INFO] Deleting segment group ID: %v\n", d.Id())

	detachSegmentGroupFromAllPolicyRules(d.Id(), policySetControllerService)

	if _, err := service.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] segment group deleted")
	return nil
}

func expandSegmentGroup(d *schema.ResourceData) segmentgroup.SegmentGroup {
	segmentGroup := segmentgroup.SegmentGroup{
		ID:                  d.Id(),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		PolicyMigrated:      d.Get("policy_migrated").(bool),
		TcpKeepAliveEnabled: d.Get("tcp_keep_alive_enabled").(string),
		MicroTenantID:       d.Get("microtenant_id").(string),
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
