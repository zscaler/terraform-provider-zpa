package zpa

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var deceptionAccessPolicyRuleExist *bool = nil
var m sync.Mutex

func validateAccessPolicyRuleOrder(order string, zClient *Client) error {
	m.Lock()
	defer m.Unlock()
	if deceptionAccessPolicyRuleExist == nil {
		policy, _, err := zClient.policysetcontroller.GetByNameAndType("ACCESS_POLICY", "Zscaler Deception")
		if err != nil || policy == nil {
			f := false
			deceptionAccessPolicyRuleExist = &f
		}
	}
	if !*deceptionAccessPolicyRuleExist {
		return nil
	}

	if order == "" {
		return nil
	}

	o, err := strconv.Atoi(order)
	if err != nil {
		return nil
	}

	if o == 1 {
		return fmt.Errorf("policy Zscaler Deception exists, order must start from 2")
	}
	return nil
}

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

func validateRuleOrders(orders *RulesOrders) error {
	sort.Slice(orders.Orders, func(i, j int) bool {
		return orders.Orders[i].Order < orders.Orders[j].Order
	})
	for i := 0; i < len(orders.Orders)-1; i++ {
		if orders.Orders[i].Order == orders.Orders[i+1].Order {
			return fmt.Errorf("duplicate order '%d' used by two rule: '%s' & '%s'", orders.Orders[i].Order, orders.Orders[i].ID, orders.Orders[i+1].ID)
		}
	}
	return nil
}

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
	sort.Slice(orders.Orders, func(i, j int) bool {
		return orders.Orders[i].Order < orders.Orders[j].Order
	})
	return &orders, nil
}

func resourcePolicyAccessReorderCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	rules, err := getRules(d, zClient)
	if err != nil {
		return err
	}
	log.Printf("[INFO] reorder rules on create: %v\n", rules)

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
	zClient := m.(*Client)
	rules, err := getRules(d, zClient)
	if err != nil {
		return err
	}
	if err := validateRuleOrders(rules); err != nil {
		log.Printf("[ERROR] reordering rules failed: %v\n", err)
		return err
	}

	remoteRules, _, err := zClient.policysetcontroller.GetAllByType(rules.PolicyType)
	if err != nil {
		log.Printf("[ERROR] failed to get rules: %v\n", err)
		return err
	}
	log.Printf("[INFO] reorder rules on update: %v\n", rules)
	orders := map[int]RuleOrder{}
	ordersList := []RuleOrder{}
	for _, r := range rules.Orders {
		orderchanged := false
		originalOrder := r.Order
		found := false
		for _, r2 := range remoteRules {
			if r.ID == r2.ID {
				found = true
				if strconv.Itoa(r.Order) != r2.RuleOrder {
					orderchanged = true
					originalOrder, _ = strconv.Atoi(r2.RuleOrder)
				}
			}
		}
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
	sort.SliceStable(ordersList, func(i, j int) bool {
		return ordersList[i].Order < ordersList[j].Order
	})
	for _, r := range ordersList {
		orderchanged := false
		originalOrder := r.Order
		found := false
		for _, r2 := range remoteRules {
			if r.ID == r2.ID {
				found = true
				if strconv.Itoa(r.Order) != r2.RuleOrder {
					orderchanged = true
					originalOrder, _ = strconv.Atoi(r2.RuleOrder)
				}
			}
		}
		if !found || !orderchanged {
			continue
		}

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
		// avoid NO adjacent rules issue
		if replacedByRule, ok := orders[r.OriginalOrder]; ok && replacedByRule.OriginalOrder == r.Order {
			continue
		}
		// reconcile the remote rules copy
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
	return resourcePolicyAccessReorderRead(d, m)
}

func resourcePolicyAccessReorderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
