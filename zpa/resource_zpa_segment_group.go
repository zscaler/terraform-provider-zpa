package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func resourceSegmentGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSegmentGroupCreate,
		ReadContext:   resourceSegmentGroupRead,
		UpdateContext: resourceSegmentGroupUpdate,
		DeleteContext: resourceSegmentGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := segmentgroup.GetByName(ctx, service, id)
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSegmentGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandSegmentGroup(d)
	log.Printf("[INFO] Creating segment group with request\n%+v\n", req)

	segmentgroup, _, err := segmentgroup.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created segment group request. ID: %v\n", segmentgroup)

	d.SetId(segmentgroup.ID)
	return resourceSegmentGroupRead(ctx, d, meta)
}

func resourceSegmentGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if service.LegacyClient != nil && service.LegacyClient.ZpaClient != nil {
		// Handle v2-specific logic here
	} else if service.Client != nil {
		// Handle v3-specific logic here
	} else {
		return diag.FromErr(fmt.Errorf("no valid client available for resource creation"))
	}

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := segmentgroup.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing segment group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting segment group:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	if err := d.Set("applications", flattenSegmentGroupApplicationsSimple(resp)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read applications %s", err))
	}
	return nil
}

func resourceSegmentGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if service.LegacyClient != nil && service.LegacyClient.ZpaClient != nil {
		// Handle v2-specific logic here
	} else if service.Client != nil {
		// Handle v3-specific logic here
	} else {
		return diag.FromErr(fmt.Errorf("no valid client available for resource creation"))
	}

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating segment group ID: %v\n", id)
	req := expandSegmentGroup(d)

	if _, _, err := segmentgroup.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := segmentgroup.UpdateV2(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceSegmentGroupRead(ctx, d, meta)
}

func resourceSegmentGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if service.LegacyClient != nil && service.LegacyClient.ZpaClient != nil {
		// Handle v2-specific logic here
	} else if service.Client != nil {
		// Handle v3-specific logic here
	} else {
		return diag.FromErr(fmt.Errorf("no valid client available for resource creation"))
	}

	// Use MicroTenant if available
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting segment group ID: %v\n", d.Id())

	if err := detachSegmentGroupFromAllPolicyRules(ctx, d.Id(), service); err != nil {
		return diag.FromErr(fmt.Errorf("error detaching SegmentGroup with ID %s from PolicySetControllers: %s", d.Id(), err))
	}

	if _, err := segmentgroup.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting SegmentGroup with ID %s: %s", d.Id(), err))
	}

	log.Printf("[INFO] Segment group with ID %s deleted", d.Id())
	d.SetId("")
	return nil
}

func expandSegmentGroup(d *schema.ResourceData) segmentgroup.SegmentGroup {
	segmentGroup := segmentgroup.SegmentGroup{
		ID:            d.Id(),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		MicroTenantID: d.Get("microtenant_id").(string),
		Applications:  expandSegmentGroupApplications(d.Get("applications").([]interface{})),
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

func detachSegmentGroupFromAllPolicyRules(ctx context.Context, id string, service *zscaler.Service) error {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()

	// Process V1 policies
	if err := detachSegmentGroupFromV1Policies(ctx, id, service); err != nil {
		return fmt.Errorf("failed to detach from v1 policies: %w", err)
	}

	// Process V2 policies
	if err := detachSegmentGroupFromV2Policies(ctx, id, service); err != nil {
		return fmt.Errorf("failed to detach from v2 policies: %w", err)
	}

	return nil
}

// detachSegmentGroupFromV1Policies handles detaching segment groups from v1 policy rules
func detachSegmentGroupFromV1Policies(ctx context.Context, id string, service *zscaler.Service) error {
	var rules []policysetcontroller.PolicyRule
	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}

	for _, t := range types {
		// Fetch the policy set by type
		policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, t)
		if err != nil {
			log.Printf("[WARN] Failed to fetch policy set of type %s: %v", t, err)
			continue
		}

		// Fetch all rules associated with the policy set
		r, _, err := policysetcontroller.GetAllByType(ctx, service, t)
		if err != nil {
			log.Printf("[WARN] Failed to fetch policy rules of type %s: %v", t, err)
			continue
		}

		// Update the policy rules with the fetched policy set ID
		for _, rule := range r {
			rule.PolicySetID = policySet.ID
			rules = append(rules, rule)
		}
	}

	log.Printf("[INFO] Detaching Segment Groups from V1 Policy Rules, count: %d", len(rules))

	for _, rr := range rules {
		rule := rr
		changed := false
		newConditions := []policysetcontroller.Conditions{}

		// Verify this rule actually exists in v1 API before processing
		// This prevents trying to update v2 rules via v1 API
		_, _, err := policysetcontroller.GetPolicyRule(ctx, service, rule.PolicySetID, rule.ID)
		if err != nil {
			if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
				log.Printf("[DEBUG] Rule %s not found in v1 API, skipping (likely a v2 rule)", rule.ID)
				continue
			}
			log.Printf("[WARN] Failed to verify v1 policy rule %s: %v", rule.ID, err)
			continue
		}

		// Process each condition
		for _, condition := range rr.Conditions {
			operands := []policysetcontroller.Operands{}
			for _, op := range condition.Operands {
				if op.ObjectType == "APP_GROUP" && op.LHS == "id" && op.RHS == id {
					changed = true
					continue
				}
				operands = append(operands, op)
			}

			// Only keep the condition if there are operands remaining
			if len(operands) > 0 {
				condition.Operands = operands
				newConditions = append(newConditions, condition)
			} else {
				// If all operands were removed, mark as changed but don't include the condition
				changed = true
			}
		}

		// Ensure conditions array is not nil
		if len(newConditions) == 0 {
			newConditions = []policysetcontroller.Conditions{}
		}

		// If the rule was changed, update it
		if changed {
			rule.Conditions = newConditions
			if _, err := policysetcontroller.UpdateRule(ctx, service, rule.PolicySetID, rule.ID, &rule); err != nil {
				log.Printf("[WARN] Failed to update v1 policy rule %s: %v", rule.ID, err)
				return fmt.Errorf("failed to update v1 policy rule %s: %w", rule.ID, err)
			}
		}
	}

	return nil
}

// detachSegmentGroupFromV2Policies handles detaching segment groups from v2 policy rules
func detachSegmentGroupFromV2Policies(ctx context.Context, id string, service *zscaler.Service) error {
	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}
	var allRules []policysetcontrollerv2.PolicyRuleResource

	for _, t := range types {
		// Fetch the policy set by type
		policySet, _, err := policysetcontrollerv2.GetByPolicyType(ctx, service, t)
		if err != nil {
			log.Printf("[WARN] Failed to fetch v2 policy set of type %s: %v", t, err)
			continue
		}

		// Fetch all rules associated with the policy set
		rules, _, err := policysetcontrollerv2.GetAllByType(ctx, service, t)
		if err != nil {
			log.Printf("[WARN] Failed to fetch v2 policy rules of type %s: %v", t, err)
			continue
		}

		// Update the policy rules with the fetched policy set ID
		for _, rule := range rules {
			rule.PolicySetID = policySet.ID
			allRules = append(allRules, rule)
		}
	}

	log.Printf("[INFO] Detaching Segment Groups from V2 Policy Rules, count: %d", len(allRules))

	for _, rule := range allRules {
		changed := false
		newConditions := []policysetcontrollerv2.PolicyRuleResourceConditions{}

		// Process each condition
		for _, cond := range rule.Conditions {
			newOperands := []policysetcontrollerv2.PolicyRuleResourceOperands{}

			// Process each operand
			for _, op := range cond.Operands {
				if strings.EqualFold(op.ObjectType, "APP_GROUP") {
					// Check if Values array contains the segment group ID
					if len(op.Values) > 0 {
						filteredValues := []string{}
						for _, v := range op.Values {
							if v == id {
								changed = true
								continue
							}
							filteredValues = append(filteredValues, v)
						}
						// Only keep the operand if there are remaining values
						if len(filteredValues) > 0 {
							op.Values = filteredValues
							newOperands = append(newOperands, op)
						} else {
							changed = true
						}
					} else if op.RHS == id {
						// Fallback: check RHS for backward compatibility
						changed = true
						continue
					} else {
						newOperands = append(newOperands, op)
					}
				} else {
					newOperands = append(newOperands, op)
				}
			}

			// Only keep the condition if there are operands remaining
			if len(newOperands) > 0 {
				cond.Operands = newOperands
				newConditions = append(newConditions, cond)
			}
		}

		// If the rule was changed, update it
		if changed {
			rule.Conditions = newConditions
			convertedRule := ConvertV1ResponseToV2Request(rule)
			if _, err := policysetcontrollerv2.UpdateRule(ctx, service, rule.PolicySetID, rule.ID, &convertedRule); err != nil {
				log.Printf("[WARN] Failed to update v2 policy rule %s: %v", rule.ID, err)
				return fmt.Errorf("failed to update v2 policy rule %s: %w", rule.ID, err)
			}
		}
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
