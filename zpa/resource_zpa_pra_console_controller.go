package zpa

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

func resourcePRAConsoleController() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePRAConsoleControllerCreate,
		ReadContext:   resourcePRAConsoleControllerRead,
		UpdateContext: resourcePRAConsoleControllerUpdate,
		DeleteContext: resourcePRAConsoleControllerDelete,
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
					resp, _, err := praconsole.GetByName(ctx, service, id)
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the privileged console",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the privileged console",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:    true,
				Description: "The description of the privileged console",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not the privileged console is enabled",
			},
			"icon_text": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:    true,
				Description: "The privileged console icon. The icon image is converted to base64 encoded text format",
			},
			"pra_portals": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							// Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The unique identifier of the privileged portal",
						},
					},
				},
			},
			"pra_application": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant.",
			},
		},
	}
}

func resourcePRAConsoleControllerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandPRAConsole(d)
	log.Printf("[INFO] Creating pra console with request\n%+v\n", req)

	praConsole, _, err := praconsole.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created pra console request. ID: %v\n", praConsole)

	d.SetId(praConsole.ID)
	return resourcePRAConsoleControllerRead(ctx, d, meta)
}

func resourcePRAConsoleControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := praconsole.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing pra console %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting pra console controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("icon_text", resp.IconText)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("pra_portals", flattenPRAPortalIDSimple(resp.PRAPortals))

	if v := flattenPRAApplicationIDSimple(resp.PRAApplication); v != nil {
		log.Printf("[DEBUG] Setting pra_application in state: %+v\n", v)
		_ = d.Set("pra_application", v)
	}

	return nil
}

func resourcePRAConsoleControllerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating pra console ID: %v\n", id)
	req := expandPRAConsole(d)

	if _, _, err := praconsole.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := praconsole.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePRAConsoleControllerRead(ctx, d, meta)
}

func resourcePRAConsoleControllerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	svc := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		svc = svc.WithMicroTenant(microTenantID)
	}

	// Detach the PRA console from any credential policies.
	if err := detachPRAConsoleFromPolicy(ctx, d.Id(), svc); err != nil {
		return diag.FromErr(fmt.Errorf("failed to detach PRA console from policy: %w", err))
	}

	log.Printf("[INFO] Deleting PRA console ID: %v", d.Id())
	if _, err := praconsole.Delete(ctx, svc, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting PRA console with ID %s: %s", d.Id(), err))
	}
	d.SetId("")
	log.Printf("[INFO] PRA console deleted")
	return nil
}

func expandPRAConsole(d *schema.ResourceData) praconsole.PRAConsole {
	result := praconsole.PRAConsole{
		ID:            d.Id(),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		IconText:      d.Get("icon_text").(string),
		MicroTenantID: d.Get("microtenant_id").(string),
		PRAPortals:    expandPRAPortal(d),
	}
	application := expandPRAApplication(d)
	if application != nil {
		result.PRAApplication = *application // TODO: Need to fix pointer to PRAApplication Struct
	}
	return result
}

func expandPRAApplication(d *schema.ResourceData) *praconsole.PRAApplication {
	if v, ok := d.GetOk("pra_application"); ok {
		applicationList := v.([]interface{})
		if len(applicationList) > 0 {
			firstApplication, ok := applicationList[0].(map[string]interface{})
			if !ok || firstApplication == nil {
				return nil
			}
			id, ok := firstApplication["id"].(string)
			if !ok || id == "" {
				return nil
			}
			return &praconsole.PRAApplication{
				ID: id,
			}
		}
	}
	return nil
}

func expandPRAPortal(d *schema.ResourceData) []praconsole.PRAPortals {
	praPortalInterface, ok := d.GetOk("pra_portals")
	if ok {
		praPortal := praPortalInterface.(*schema.Set)
		log.Printf("[INFO] pra portal data: %+v\n", praPortal)
		var praPortals []praconsole.PRAPortals
		for _, praPortal := range praPortal.List() {
			praPortal, ok := praPortal.(map[string]interface{})
			if ok {
				for _, id := range praPortal["id"].(*schema.Set).List() {
					praPortals = append(praPortals, praconsole.PRAPortals{
						ID: id.(string),
					})
				}
			}
		}
		return praPortals
	}

	return []praconsole.PRAPortals{}
}

func flattenPRAPortalIDSimple(praPortals []praconsole.PRAPortals) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(praPortals))
	for i, portal := range praPortals {
		ids[i] = portal.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenPRAApplicationIDSimple(praApplication praconsole.PRAApplication) []interface{} {
	if praApplication.ID == "" {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": praApplication.ID,
		},
	}
}

func detachPRAConsoleFromPolicy(ctx context.Context, id string, svc *zscaler.Service) error {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()

	// Get all rules for relevant policy types (expand if needed)
	types := []string{"CREDENTIAL_POLICY"} // Add others if needed
	var allRules []policysetcontrollerv2.PolicyRuleResource

	for _, t := range types {
		policySet, _, err := policysetcontrollerv2.GetByPolicyType(ctx, svc, t)
		if err != nil {
			return fmt.Errorf("failed to get policy set for type %s: %w", t, err)
		}

		rules, _, err := policysetcontrollerv2.GetAllByType(ctx, svc, t)
		if err != nil {
			return fmt.Errorf("failed to get rules for policy type %s: %w", t, err)
		}

		for _, rule := range rules {
			rule.PolicySetID = policySet.ID
			allRules = append(allRules, rule)
		}
	}

	// Track rules we modified
	modifiedRules := make(map[string]bool)
	// Phase 1: Update all rules to remove the PRA console reference
	for _, rule := range allRules {
		modifiedRules[rule.ID] = false
		newConditions := []policysetcontrollerv2.PolicyRuleResourceConditions{}
		for _, cond := range rule.Conditions {
			newOperands := []policysetcontrollerv2.PolicyRuleResourceOperands{}
			for _, op := range cond.Operands {
				if strings.EqualFold(op.ObjectType, "CONSOLE") && strings.EqualFold(op.LHS, "id") {
					log.Printf("[DEBUG] Detaching PRA console from rule %s, op: %+v , id: %s", rule.ID, op, id)
					// Filter out our PRA console ID
					if len(op.Values) > 0 {
						filteredValues := []string{}
						for _, v := range op.Values {
							if v == id {
								modifiedRules[rule.ID] = true
								continue
							}
							filteredValues = append(filteredValues, v)
						}
						op.Values = filteredValues
						if len(filteredValues) > 0 {
							newOperands = append(newOperands, op)
						}
					} else if op.RHS == id {
						modifiedRules[rule.ID] = true
						continue
					} else {
						newOperands = append(newOperands, op)
					}
				} else {
					newOperands = append(newOperands, op)
				}
			}

			if len(newOperands) > 0 {
				cond.Operands = newOperands
				newConditions = append(newConditions, cond)
			}
		}

		if modifiedRules[rule.ID] {
			data, _ := json.MarshalIndent(rule.Conditions, "", "  ")
			log.Printf("[DEBUG] rule Conditions before update: %s", string(data))
			rule.Conditions = newConditions
			data, _ = json.MarshalIndent(rule.Conditions, "", "  ")
			log.Printf("[DEBUG] rule Conditionsafter update: %s", string(data))

			// Retry updates in case of transient errors
			maxRetries := 3
			for i := 0; i < maxRetries; i++ {
				convertedRule := ConvertV1ResponseToV2Request(rule)
				_, err := policysetcontrollerv2.UpdateRule(ctx, svc, rule.PolicySetID, rule.ID, &convertedRule)
				if err == nil {
					break
				}
				if i == maxRetries-1 {
					return fmt.Errorf("failed to update rule %s after %d attempts: %w", rule.ID, maxRetries, err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}
