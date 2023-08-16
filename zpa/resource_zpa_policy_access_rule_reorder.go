package zpa

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Global variables for state management
var deceptionAccessPolicyRuleExist *bool // Pointer to check if deception rule exists.
var m sync.Mutex                         // Mutex to ensure thread safety.

// Validate the access policy rule's order.
func validateAccessPolicyRuleOrder(order string, zClient *Client) error {
	m.Lock()
	defer m.Unlock()

	// Check if we've already verified the existence of the Deception rule.
	if deceptionAccessPolicyRuleExist == nil {
		policy, _, err := zClient.policysetcontroller.GetByNameAndType("ACCESS_POLICY", "Zscaler Deception")
		if err != nil || policy == nil {
			f := false
			deceptionAccessPolicyRuleExist = &f
		} else {
			t := true
			deceptionAccessPolicyRuleExist = &t
		}
	}
	// If Deception rule doesn't exist or the order is empty, no further checks needed.
	if deceptionAccessPolicyRuleExist != nil && !*deceptionAccessPolicyRuleExist || order == "" {
		return nil
	}

	if order == "" {
		return nil
	}
	// Convert string order to integer.
	o, err := strconv.Atoi(order)
	if err != nil {
		return nil
	}
	// If the Deception rule exists, order should start from 2.
	if o == 1 {
		return fmt.Errorf("policy Zscaler Deception exists, order must start from 2")
	}
	return nil
}

// Define the Terraform resource for reordering policy access rules.
func resourcePolicyAccessRuleReorder() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyAccessReorderCreate,
		Read:   resourcePolicyAccessReorderRead,
		Update: resourcePolicyAccessReorderUpdate,
		Delete: resourcePolicyAccessReorderDelete,
		Schema: map[string]*schema.Schema{
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of rules and their orders",
				MaxItems:    1000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"order": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Data structures for rule ordering.
type RuleOrder struct {
	ID            string
	Order         int
	OriginalOrder int
}

type RulesOrders struct {
	PolicySetID string
	PolicyType  string
	Orders      []RuleOrder
}

// Validate that no two rules have the same order.
func validateRuleOrders(orders *RulesOrders) error {
	// Sort rules by order.
	sort.Slice(orders.Orders, func(i, j int) bool {
		return orders.Orders[i].Order < orders.Orders[j].Order
	})
	// Check for duplicate order values.
	for i := 0; i < len(orders.Orders)-1; i++ {
		if orders.Orders[i].Order == orders.Orders[i+1].Order {
			return fmt.Errorf("duplicate order '%d' used by two rules: '%s' & '%s'", orders.Orders[i].Order, orders.Orders[i].ID, orders.Orders[i+1].ID)
		}
	}
	return nil
}

// Fetch and sort the rule orders from the provided data.
func getRules(d *schema.ResourceData, zClient *Client) (*RulesOrders, error) {
	policyType := d.Get("policy_type").(string)
	globalPolicySet, err := GetGlobalPolicySetByPolicyType(zClient.policysetcontroller, policyType)
	if err != nil {
		log.Printf("[ERROR] reordering rules failed getting global policy set '%s': %v\n", policyType, err)
		return nil, err
	}
	orders := RulesOrders{
		PolicySetID: globalPolicySet.ID,
		PolicyType:  policyType,
		Orders:      []RuleOrder{},
	}
	// Extract rules from the data.
	rulesSet, ok := d.Get("rules").(*schema.Set)
	if ok && rulesSet != nil {
		for _, r := range rulesSet.List() {
			rule, ok := r.(map[string]interface{})
			if !ok || rule == nil {
				continue
			}
			id := rule["id"].(string)
			order, _ := strconv.Atoi(rule["order"].(string))
			orders.Orders = append(orders.Orders, RuleOrder{
				ID:    id,
				Order: order,
			})
		}
	}
	// Sort the rules by their order.
	sort.Slice(orders.Orders, func(i, j int) bool {
		return orders.Orders[i].Order < orders.Orders[j].Order
	})
	return &orders, nil
}

// Create operation for the Terraform resource.
func resourcePolicyAccessReorderCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	rules, err := getRules(d, zClient)
	if err != nil {
		return err
	}
	log.Printf("[INFO] reorder rules on create: %v\n", rules)

	// Validate the orders.
	if err := validateRuleOrders(rules); err != nil {
		log.Printf("[ERROR] reordering rules failed: %v\n", err)
		return err
	}

	for _, r := range rules.Orders {
		if rules.PolicyType == "ACCESS_POLICY" {
			if err := validateAccessPolicyRuleOrder(strconv.Itoa(r.Order), zClient); err != nil {
				log.Printf("[ERROR] reordering rule ID '%s' failed, order validation error: %v\n", r.ID, err)
				continue
			}
		}
		_, err := zClient.policysetcontroller.Reorder(rules.PolicySetID, r.ID, r.Order)
		if err != nil {
			log.Printf("[ERROR] reordering rule ID '%s' failed: %v\n", r.ID, err)
		}
	}

	d.SetId(rules.PolicySetID)
	return resourcePolicyAccessReorderRead(d, m)
}

func resourcePolicyAccessReorderRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	rulesOrders, err := getRules(d, zClient)
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
		for i, r2 := range rulesOrders.Orders {
			if r.ID == r2.ID {
				rulesOrders.Orders[i].Order, _ = strconv.Atoi(r.RuleOrder)
			}
		}
	}
	rulesMap := []map[string]interface{}{}
	for _, r := range rulesOrders.Orders {
		rulesMap = append(rulesMap, map[string]interface{}{
			"id":    r.ID,
			"order": strconv.Itoa(r.Order),
		})
	}
	_ = d.Set("rules", rulesMap)
	return nil
}

func resourcePolicyAccessReorderUpdate(d *schema.ResourceData, m interface{}) error {
	// Convert the interface to a client instance.
	zClient := m.(*Client)
	// Fetch and sort the rule orders from the provided data.
	rules, err := getRules(d, zClient)
	if err != nil {
		return err
	}
	// Validate the fetched rule orders.
	if err := validateRuleOrders(rules); err != nil {
		log.Printf("[ERROR] reordering rules failed: %v\n", err)
		return err
	}
	// Fetch the existing remote rules based on the policy type.
	remoteRules, _, err := zClient.policysetcontroller.GetAllByType(rules.PolicyType)
	if err != nil {
		log.Printf("[ERROR] failed to get rules: %v\n", err)
		return err
	}
	log.Printf("[INFO] reorder rules on update: %v\n", rules)
	// Maps and slices for storing rule orders.
	orders := map[int]RuleOrder{}
	ordersList := []RuleOrder{}

	// Iterate over the fetched rule orders to determine changes.
	for _, r := range rules.Orders {
		orderchanged := false
		originalOrder := r.Order
		found := false

		// Check if there's a change in order for each rule against the remote set.
		for _, r2 := range remoteRules {
			if r.ID == r2.ID {
				found = true
				if strconv.Itoa(r.Order) != r2.RuleOrder {
					orderchanged = true
					originalOrder, _ = strconv.Atoi(r2.RuleOrder)
				}
			}
		}

		// If no match was found or order did not change, skip to the next iteration.
		if !found || !orderchanged {
			continue
		}
		o := RuleOrder{
			ID:            r.ID,
			Order:         r.Order,
			OriginalOrder: originalOrder,
		}
		orders[r.Order] = o
		ordersList = append(ordersList, o)
	}

	// Sort rules based on the order field.
	sort.SliceStable(ordersList, func(i, j int) bool {
		return ordersList[i].Order < ordersList[j].Order
	})

	// Re-check and re-order the rule set.
	for _, r := range ordersList {
		orderchanged := false
		originalOrder := r.Order
		found := false

		// Check if there's a change in order for each rule against the remote set.
		for _, r2 := range remoteRules {
			if r.ID == r2.ID {
				found = true
				if strconv.Itoa(r.Order) != r2.RuleOrder {
					orderchanged = true
					originalOrder, _ = strconv.Atoi(r2.RuleOrder)
				}
			}
		}
		// If no match was found or order did not change, skip to the next iteration.
		if !found || !orderchanged {
			continue
		}

		// Check for special rules related to 'ACCESS_POLICY'.
		if rules.PolicyType == "ACCESS_POLICY" {
			if err := validateAccessPolicyRuleOrder(strconv.Itoa(r.Order), zClient); err != nil {
				log.Printf("[ERROR] reordering rule ID '%s' failed, order validation error: %v\n", r.ID, err)
				continue
			}
		}

		// Request the service to reorder the rules.
		_, err := zClient.policysetcontroller.Reorder(rules.PolicySetID, r.ID, r.Order)
		if err != nil {
			log.Printf("[ERROR] reordering rule ID '%s' failed: %v\n", r.ID, err)
		}
		// avoid NO adjacent rules issue
		// Handle potential ordering issue related to adjacency.
		if replacedByRule, ok := orders[r.OriginalOrder]; ok && replacedByRule.OriginalOrder == r.Order && r.Order != replacedByRule.Order+1 && r.Order != replacedByRule.Order-1 {
			continue
		}
		// reconcile the remote rules copy
		// Re-adjust the order of rules in the remote copy for consistency.
		for i := range remoteRules {
			if r.ID == remoteRules[i].ID {
				continue
			}
			o, _ := strconv.Atoi(remoteRules[i].RuleOrder)
			if originalOrder > r.Order && o >= r.Order && o < originalOrder {
				remoteRules[i].RuleOrder = strconv.Itoa(o + 1)
			} else if originalOrder < r.Order && o <= r.Order && o > originalOrder {
				remoteRules[i].RuleOrder = strconv.Itoa(o - 1)
			}
		}
	}
	// Read the updated rule set.
	return resourcePolicyAccessReorderRead(d, m)
}

func resourcePolicyAccessReorderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
