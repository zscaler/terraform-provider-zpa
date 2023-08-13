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
			"policy_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
	ID    string
	Order int
}

type RulesOrders struct {
	PolicySetID string
	PolicyType  string
	Orders      []RuleOrder
}

func getRules(d *schema.ResourceData) RulesOrders {
	policySetID := d.Get("policy_set_id").(string)
	policyType := d.Get("policy_type").(string)
	orders := RulesOrders{
		PolicySetID: policySetID,
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
	return orders
}

func resourcePolicyAccessReorderCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	rules := getRules(d)
	log.Printf("[INFO] reorder rules on create: %v\n", rules)
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
	rulesOrders := getRules(d)
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
	rules := getRules(d)
	log.Printf("[INFO] reorder rules on update: %v\n", rules)
	for _, r := range rules.Orders {
		if err := validateAccessPolicyRuleOrder(strconv.Itoa(r.Order), zClient); err != nil {
			log.Printf("[ERROR] reordering rule ID '%s' failed, order validation error: %v\n", r.ID, err)
			continue
		}
		_, err := zClient.policysetcontroller.Reorder(rules.PolicySetID, r.ID, r.Order)
		if err != nil {
			log.Printf("[ERROR] reordering rule ID '%s' failed: %v\n", r.ID, err)
		}
	}
	return resourcePolicyAccessReorderRead(d, m)
}

func resourcePolicyAccessReorderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
