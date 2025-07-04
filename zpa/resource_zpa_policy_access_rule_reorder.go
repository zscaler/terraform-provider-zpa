package zpa

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

// Define the Terraform resource for reordering policy access rules.
func resourcePolicyAccessRuleReorder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAccessReorderUpdate,
		ReadContext:   resourcePolicyAccessReorderRead,
		UpdateContext: resourcePolicyAccessReorderUpdate,
		DeleteContext: resourcePolicyAccessReorderDelete,
		Schema: map[string]*schema.Schema{
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ACCESS_POLICY",
					"GLOBAL_POLICY",
					"CAPABILITIES_POLICY",
					"BYPASS_POLICY",
					"CLIENT_FORWARDING_POLICY",
					"CREDENTIAL_POLICY",
					"ISOLATION_POLICY",
					"INSPECTION_POLICY",
					"REDIRECTION_POLICY",
					"REAUTH_POLICY",
					"TIMEOUT_POLICY",
					"CLIENTLESS_SESSION_PROTECTION_POLICY",
				}, false),
			},
			"rules": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of rules and their orders",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"order": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(v interface{}, _ cty.Path) diag.Diagnostics {
								order, _ := strconv.Atoi(v.(string))
								if order <= 0 {
									return diag.Diagnostics{
										diag.Diagnostic{
											Severity: diag.Error,
											Summary:  "Rule order 0 is not allowed",
											Detail:   "Orders must start from 1, got:" + v.(string),
										},
									}
								}
								return nil
							},
						},
					},
				},
			},
		},
	}
}

type RulesOrders struct {
	PolicyType string
	Orders     map[string]int
}

// Validate that no two rules have the same order.
func validateRuleOrders(orders *RulesOrders) error {
	// Check for orders <= 0
	for _, order := range orders.Orders {
		if order <= 0 {
			return fmt.Errorf("order must be a positive integer greater than 0")
		}
	}
	// Check for duplicate order values.
	if dupOrder, dupRuleIDs, ok := hasDuplicates(orders.Orders); ok {
		return fmt.Errorf("duplicate order '%d' used by rules with IDs: %v", dupOrder, strings.Join(dupRuleIDs, ", "))
	}

	return nil
}

// Check for duplicate order values.
func hasDuplicates(orders map[string]int) (int, []string, bool) {
	ruleSet := make(map[int][]string)
	for id, order := range orders {
		ruleSet[order] = append(ruleSet[order], id)
	}

	for order, ruleIDs := range ruleSet {
		if len(ruleIDs) > 1 {
			return order, ruleIDs, true
		}
	}
	return 0, nil, false
}

// func getRules(d *schema.ResourceData) (*RulesOrders, error) {
// 	policyType := d.Get("policy_type").(string)
// 	orders := RulesOrders{
// 		PolicyType: policyType,
// 		Orders:     map[string]int{},
// 	}
// 	rulesSet, ok := d.Get("rules").(*schema.Set)
// 	if ok && rulesSet != nil {
// 		for _, r := range rulesSet.List() {
// 			rule, ok := r.(map[string]interface{})
// 			if !ok || rule == nil {
// 				continue
// 			}
// 			id := rule["id"].(string)
// 			order, _ := strconv.Atoi(rule["order"].(string))
// 			orders.Orders[id] = order
// 		}
// 	}
// 	return &orders, nil
// }

// func getRules(d *schema.ResourceData) (*RulesOrders, error) {
// 	policyType := d.Get("policy_type").(string)
// 	orders := RulesOrders{
// 		PolicyType: policyType,
// 		Orders:     map[string]int{},
// 	}

// 	rulesSet, ok := d.Get("rules").(*schema.Set)
// 	if ok && rulesSet != nil {
// 		for _, r := range rulesSet.List() {
// 			rule, ok := r.(map[string]interface{})
// 			if !ok || rule == nil {
// 				continue
// 			}
// 			id := rule["id"].(string)
// 			order, _ := strconv.Atoi(rule["order"].(string))
// 			// Skip if this is the Zscaler Deception rule (rule 1)
// 			if order == 1 {
// 				log.Printf("[INFO] Skipping rule ID %s as it appears to be Zscaler Deception", id)
// 				continue
// 			}
// 			orders.Orders[id] = order
// 		}
// 	}
// 	return &orders, nil
// }

func getRules(d *schema.ResourceData) (*RulesOrders, error) {
	policyType := d.Get("policy_type").(string)
	orders := RulesOrders{
		PolicyType: policyType,
		Orders:     map[string]int{},
	}

	rulesSet, ok := d.Get("rules").(*schema.Set)
	if ok && rulesSet != nil {
		for _, r := range rulesSet.List() {
			rule, ok := r.(map[string]interface{})
			if !ok || rule == nil {
				continue
			}
			id := rule["id"].(string)
			order, _ := strconv.Atoi(rule["order"].(string))

			// ❌ REMOVE THIS FAULTY SKIP
			// if order == 1 {
			// 	log.Printf("[INFO] Skipping rule ID %s as it appears to be Zscaler Deception", id)
			// 	continue
			// }

			orders.Orders[id] = order
		}
	}
	return &orders, nil
}

func resourcePolicyAccessReorderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policyType := d.Get("policy_type").(string)

	currentRules, _, err := policysetcontroller.GetAllByType(ctx, service, policyType)
	if err != nil {
		log.Printf("[ERROR] failed to get rules: %v\n", err)
		d.SetId("")
		return diag.FromErr(err)
	}

	configuredRules, err := getRules(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] reorder rules on read: %v\n", configuredRules)

	currentOrderMap := make(map[string]int)
	for _, rule := range currentRules {
		if order, err := strconv.Atoi(rule.RuleOrder); err == nil {
			currentOrderMap[rule.ID] = order
		}
	}

	for id := range configuredRules.Orders {
		if currentOrder, exists := currentOrderMap[id]; exists {
			configuredRules.Orders[id] = currentOrder
		}
	}

	rulesMap := []map[string]interface{}{}
	for id, order := range configuredRules.Orders {
		rulesMap = append(rulesMap, map[string]interface{}{
			"id":    id,
			"order": strconv.Itoa(order),
		})
	}

	if err := d.Set("rules", rulesMap); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s-%s", policyType, "reorder"))

	return nil
}

func resourcePolicyAccessReorderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	existingRules, _, err := policysetcontroller.GetAllByType(ctx, service, d.Get("policy_type").(string))
	if err != nil {
		log.Printf("[ERROR] Failed to get existing rules: %v\n", err)
		return diag.FromErr(err)
	}

	deceptionAtOne := false
	deceptionID := ""
	for _, rule := range existingRules {
		if rule.Name == "Zscaler Deception" && rule.RuleOrder == "1" {
			deceptionAtOne = true
			deceptionID = rule.ID
			break
		}
	}

	userDefinedRules, err := getRules(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := validateRuleOrders(userDefinedRules); err != nil {
		log.Printf("[ERROR] Reordering rules failed: %v\n", err)
		return diag.FromErr(err)
	}

	ruleIdToOrder := make(map[string]int)

	if deceptionAtOne {
		_, managedByTerraform := userDefinedRules.Orders[deceptionID]
		if !managedByTerraform {
			ruleIdToOrder[deceptionID] = 1
		}
	}

	for id, order := range userDefinedRules.Orders {
		if id == deceptionID {
			continue
		}

		if deceptionAtOne {
			ruleIdToOrder[id] = order + 1
		} else {
			ruleIdToOrder[id] = order
		}
	}

	if _, err := policysetcontroller.BulkReorder(ctx, service, d.Get("policy_type").(string), ruleIdToOrder); err != nil {
		log.Printf("[ERROR] Bulk reordering rules failed: %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s-%s", d.Get("policy_type").(string), "reorder"))
	return resourcePolicyAccessReorderRead(ctx, d, meta)
}

func resourcePolicyAccessReorderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
