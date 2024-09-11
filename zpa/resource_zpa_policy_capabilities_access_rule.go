package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontrollerv2"
)

func resourcePolicyCapabilitiesAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCapabilitiesAccessRuleCreate,
		Read:   resourcePolicyCapabilitiesAccessRuleRead,
		Update: resourcePolicyCapabilitiesAccessRuleUpdate,
		Delete: resourcePolicyCapabilitiesAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"CAPABILITIES_POLICY"}),
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the name of the policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the description of the access policy.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "  This is for providing the rule action.",
				ValidateFunc: validation.StringInSlice([]string{
					"CHECK_CAPABILITIES",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "This is for proviidng the set of conditions for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AND",
								"OR",
							}, false),
						},
						"operands": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "This signifies the various policy criteria.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "This denotes a list of values for the given object type. The value depend upon the key. If rhs is defined this list will be ignored",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"object_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "  This is for specifying the policy critiera.",
										ValidateFunc: validation.StringInSlice([]string{
											"APP",
											"APP_GROUP",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
										}, false),
									},
									"entry_values": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rhs": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"lhs": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"privileged_capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clipboard_copy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA Clipboard Copy function",
						},
						"clipboard_paste": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA Clipboard Paste function",
						},
						"file_upload": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA File Transfer capabilities that enables the File Upload function",
						},
						"file_download": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA File Transfer capabilities that enables the File Download function",
						},
						"inspect_file_download": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Inspects the file via ZIA sandbox (if you have set up the ZIA cloud and the Integrations settings) and downloads the file following the inspection",
						},
						"inspect_file_upload": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Inspects the file via ZIA sandbox (if you have set up the ZIA cloud and the Integrations settings) and uploads the file following the inspection",
						},
						"monitor_session": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA Monitoring Capabilities to enable the PRA Session Monitoring function",
						},
						"record_session": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA Session Recording capabilities to enable PRA Session Recording",
						},
						"share_session": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the PRA Session Control and Monitoring capabilities to enable PRA Session Monitoring",
						},
					},
				},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePolicyCapabilitiesAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	// Automatically determining policy_set_id for "CAPABILITIES_POLICY"
	policySetID, err := fetchPolicySetIDByType(zClient, "CAPABILITIES_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandPrivilegedCapabilitiesRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy capabilities rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	resp, _, err := policysetcontrollerv2.CreateRule(service, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePolicyCapabilitiesAccessRuleRead(d, meta)
}

func resourcePolicyCapabilitiesAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(zClient, "CAPABILITIES_POLICY", microTenantID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontrollerv2.GetPolicyRule(service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	v2PolicyRule := ConvertV1ResponseToV2Request(*resp)

	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", v2PolicyRule.Name)
	_ = d.Set("description", v2PolicyRule.Description)
	_ = d.Set("action", v2PolicyRule.Action)
	_ = d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	_ = d.Set("policy_set_id", policySetID) // Here, you're setting it based on fetched ID
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))
	if len(resp.PrivilegedCapabilities.Capabilities) > 0 {
		_ = d.Set("privileged_capabilities", flattenPrivilegedCapabilities(resp.PrivilegedCapabilities))
	}
	return nil
}

func resourcePolicyCapabilitiesAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "CAPABILITIES_POLICY"
	policySetID, err := fetchPolicySetIDByType(zClient, "CAPABILITIES_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating policy capabilities rule ID: %v\n", ruleID)
	req, err := expandPrivilegedCapabilitiesRule(d, policySetID)
	if err != nil {
		return err
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	// Checking the current state of the rule to handle cases where it might have been deleted outside Terraform
	_, respErr, err := policysetcontrollerv2.GetPolicyRule(service, policySetID, ruleID)
	if err != nil {
		if respErr != nil && (respErr.StatusCode == http.StatusNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}
	_, err = policysetcontrollerv2.UpdateRule(service, policySetID, ruleID, req)
	if err != nil {
		return err
	}

	return resourcePolicyCapabilitiesAccessRuleRead(d, meta)
}

func resourcePolicyCapabilitiesAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	// Assume "CAPABILITIES_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(zClient, "CAPABILITIES_POLICY", microTenantID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	service := zClient.PolicySetControllerV2
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	if _, err := policysetcontrollerv2.Delete(service, policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func flattenPrivilegedCapabilities(capabilities policysetcontrollerv2.PrivilegedCapabilities) []interface{} {
	capMap := make(map[string]bool)
	for _, cap := range capabilities.Capabilities {
		switch cap {
		case "CLIPBOARD_COPY":
			capMap["clipboard_copy"] = true
		case "CLIPBOARD_PASTE":
			capMap["clipboard_paste"] = true
		case "FILE_DOWNLOAD":
			capMap["file_download"] = true
		case "FILE_UPLOAD":
			capMap["file_upload"] = true
		case "INSPECT_FILE_DOWNLOAD":
			capMap["inspect_file_download"] = true
		case "INSPECT_FILE_UPLOAD":
			capMap["inspect_file_upload"] = true
		case "MONITOR_SESSION":
			capMap["monitor_session"] = true
		case "RECORD_SESSION":
			capMap["record_session"] = true
		case "SHARE_SESSION":
			capMap["share_session"] = true
		}
	}

	return []interface{}{map[string]interface{}{
		"clipboard_copy":        capMap["clipboard_copy"],
		"clipboard_paste":       capMap["clipboard_paste"],
		"file_download":         capMap["file_download"],
		"file_upload":           capMap["file_upload"],
		"inspect_file_download": capMap["inspect_file_download"],
		"inspect_file_upload":   capMap["inspect_file_upload"],
		"monitor_session":       capMap["monitor_session"],
		"record_session":        capMap["record_session"],
		"share_session":         capMap["share_session"],
	}}
}

func expandPrivilegedCapabilitiesRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice for capabilities
	capabilities := []string{}

	// Check if the privileged_capabilities block exists
	if v, ok := d.GetOk("privileged_capabilities"); ok {
		privCapsList := v.([]interface{})
		if len(privCapsList) > 0 {
			// Extract the map from the first item of the list (MaxItems: 1)
			privCapsMap := privCapsList[0].(map[string]interface{})

			// Convert Boolean values to the API expected string values
			if privCapsMap["clipboard_copy"].(bool) {
				capabilities = append(capabilities, "CLIPBOARD_COPY")
			}
			if privCapsMap["clipboard_paste"].(bool) {
				capabilities = append(capabilities, "CLIPBOARD_PASTE")
			}
			if privCapsMap["file_download"].(bool) {
				capabilities = append(capabilities, "FILE_DOWNLOAD")
			}
			if privCapsMap["file_upload"].(bool) {
				capabilities = append(capabilities, "FILE_UPLOAD")
			}
			if privCapsMap["inspect_file_download"].(bool) {
				capabilities = append(capabilities, "INSPECT_FILE_DOWNLOAD")
			}
			if privCapsMap["inspect_file_upload"].(bool) {
				capabilities = append(capabilities, "INSPECT_FILE_UPLOAD")
			}
			if privCapsMap["monitor_session"].(bool) {
				capabilities = append(capabilities, "MONITOR_SESSION")
			}
			if privCapsMap["record_session"].(bool) {
				capabilities = append(capabilities, "RECORD_SESSION")
			}
			if privCapsMap["share_session"].(bool) {
				capabilities = append(capabilities, "SHARE_SESSION")
			}
		}
	}

	// Construct the PrivilegedCapabilities struct
	privilegedCapabilities := policysetcontrollerv2.PrivilegedCapabilities{
		Capabilities: capabilities,
	}

	// Construct the PolicyRule struct
	policyRule := &policysetcontrollerv2.PolicyRule{
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Action:                 d.Get("action").(string),
		MicroTenantID:          d.Get("microtenant_id").(string),
		PolicySetID:            policySetID,
		Conditions:             conditions,
		PrivilegedCapabilities: privilegedCapabilities,
	}

	return policyRule, nil
}
