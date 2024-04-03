package zpa

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Define the Terraform resource for reordering policy access rules.
func resourcePolicyAccessRuleReorder() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyAccessReorderUpdate,
		Read:   resourcePolicyAccessReorderRead,
		Update: resourcePolicyAccessReorderUpdate,
		Delete: resourcePolicyAccessReorderDelete,
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
			orders.Orders[id] = order
		}
	}
	return &orders, nil
}

func resourcePolicyAccessReorderRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	rulesOrders, err := getRules(d)
	if err != nil {
		return err
	}
	rules, _, err := zClient.policysetcontroller.GetAllByType(rulesOrders.PolicyType)
	if err != nil {
		log.Printf("[ERROR] failed to get rules: %v\n", err)
		return err
	}
	log.Printf("[INFO] reorder rules on read: %v\n", rulesOrders)
	for _, r := range rules {
		for id := range rulesOrders.Orders {
			if r.ID == id {
				rulesOrders.Orders[id], _ = strconv.Atoi(r.RuleOrder)
			}
		}
	}
	rulesMap := []map[string]interface{}{}
	for id, order := range rulesOrders.Orders {
		rulesMap = append(rulesMap, map[string]interface{}{
			"id":    id,
			"order": strconv.Itoa(order),
		})
	}
	_ = d.Set("rules", rulesMap)
	return nil
}

func resourcePolicyAccessReorderUpdate(d *schema.ResourceData, m interface{}) error {
	// Convert the interface to a client instance.
	zClient := m.(*Client)
	// Fetch and sort the rule orders from the provided data.
	rules, err := getRules(d)
	if err != nil {
		return err
	}
	// Validate the fetched rule orders.
	if err := validateRuleOrders(rules); err != nil {
		log.Printf("[ERROR] reordering rules failed: %v\n", err)
		return err
	}
	d.SetId(rules.PolicyType)
	_, err = zClient.policysetcontroller.BulkReorder(rules.PolicyType, rules.Orders)
	if err != nil {
		log.Printf("[ERROR] reordering rules failed: %v\n", err)
	}
	// Read the updated rule set.
	return resourcePolicyAccessReorderRead(d, m)
}

func resourcePolicyAccessReorderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
