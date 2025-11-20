package resources

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

var (
	_ resource.Resource                = &PolicyAccessRuleReorderResource{}
	_ resource.ResourceWithConfigure   = &PolicyAccessRuleReorderResource{}
	_ resource.ResourceWithImportState = &PolicyAccessRuleReorderResource{}
)

var positiveOrderRegex = regexp.MustCompile("^[1-9][0-9]*$")

func NewPolicyAccessRuleReorderResource() resource.Resource {
	return &PolicyAccessRuleReorderResource{}
}

type PolicyAccessRuleReorderResource struct {
	client *client.Client
}

type PolicyAccessRuleReorderModel struct {
	ID            types.String `tfsdk:"id"`
	PolicyType    types.String `tfsdk:"policy_type"`
	Rules         types.Set    `tfsdk:"rules"`
	MicrotenantID types.String `tfsdk:"microtenant_id"`
}

type PolicyAccessRuleReorderRuleModel struct {
	ID    types.String `tfsdk:"id"`
	Order types.String `tfsdk:"order"`
}

type rulesOrders struct {
	PolicyType string
	Orders     map[string]int
}

func (r *PolicyAccessRuleReorderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_access_rule_reorder"
}

func (r *PolicyAccessRuleReorderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the order of rules within a specific policy type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the reorder operation, derived from the policy type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of policy whose rules should be reordered.",
				Validators: []validator.String{
					stringvalidator.OneOf(
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
					),
				},
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant identifier used to scope API calls.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rules": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "Identifier of the policy rule.",
						},
						"order": schema.StringAttribute{
							Required:    true,
							Description: "Desired order of the rule, starting from 1.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(positiveOrderRegex, "order must be a positive integer greater than 0"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *PolicyAccessRuleReorderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *PolicyAccessRuleReorderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan PolicyAccessRuleReorderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.reorderRules(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.refreshState(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PolicyAccessRuleReorderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state PolicyAccessRuleReorderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.refreshState(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PolicyAccessRuleReorderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan PolicyAccessRuleReorderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.reorderRules(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.refreshState(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PolicyAccessRuleReorderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *PolicyAccessRuleReorderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policyType := strings.TrimSpace(req.ID)
	if policyType == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "The import identifier must not be empty.")
		return
	}

	policyType = strings.TrimSuffix(policyType, "-reorder")

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), policyType+"-reorder")...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_type"), policyType)...)
}

func (r *PolicyAccessRuleReorderResource) reorderRules(ctx context.Context, data *PolicyAccessRuleReorderModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.PolicyType.IsNull() || data.PolicyType.IsUnknown() {
		diags.AddError("Missing Policy Type", "The policy_type attribute must be provided.")
		return diags
	}

	orders, rd := r.getRules(ctx, data)
	diags.Append(rd...)
	if diags.HasError() {
		return diags
	}

	if err := validateRuleOrders(orders); err != nil {
		diags.AddError("Invalid Rule Orders", err.Error())
		return diags
	}

	service := r.client.Service
	if !data.MicrotenantID.IsNull() && !data.MicrotenantID.IsUnknown() {
		microtenantID := strings.TrimSpace(data.MicrotenantID.ValueString())
		if microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}
	}

	existingRules, _, err := policysetcontroller.GetAllByType(ctx, service, orders.PolicyType)
	if err != nil {
		diags.AddError("Failed to Retrieve Existing Rules", fmt.Sprintf("Error retrieving rules for policy type %s: %v", orders.PolicyType, err))
		return diags
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

	ruleIDToOrder := make(map[string]int, len(orders.Orders)+1)

	if deceptionAtOne {
		if _, managed := orders.Orders[deceptionID]; !managed {
			ruleIDToOrder[deceptionID] = 1
		}
	}

	for id, order := range orders.Orders {
		if id == deceptionID {
			continue
		}

		if deceptionAtOne {
			ruleIDToOrder[id] = order + 1
		} else {
			ruleIDToOrder[id] = order
		}
	}

	if _, err := policysetcontroller.BulkReorder(ctx, service, orders.PolicyType, ruleIDToOrder); err != nil {
		diags.AddError("Failed to Reorder Rules", fmt.Sprintf("Error reordering rules for policy type %s: %v", orders.PolicyType, err))
		return diags
	}

	data.ID = types.StringValue(orders.PolicyType + "-reorder")
	return diags
}

func (r *PolicyAccessRuleReorderResource) refreshState(ctx context.Context, data *PolicyAccessRuleReorderModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.PolicyType.IsNull() || data.PolicyType.IsUnknown() {
		diags.AddError("Missing Policy Type", "The policy_type attribute must be provided.")
		return diags
	}

	orders, rd := r.getRules(ctx, data)
	diags.Append(rd...)
	if diags.HasError() {
		return diags
	}

	service := r.client.Service
	if !data.MicrotenantID.IsNull() && !data.MicrotenantID.IsUnknown() {
		microtenantID := strings.TrimSpace(data.MicrotenantID.ValueString())
		if microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}
	}

	currentRules, _, err := policysetcontroller.GetAllByType(ctx, service, orders.PolicyType)
	if err != nil {
		diags.AddError("Failed to Retrieve Rules", fmt.Sprintf("Error retrieving rules for policy type %s: %v", orders.PolicyType, err))
		return diags
	}

	currentOrderMap := make(map[string]int, len(currentRules))
	for _, rule := range currentRules {
		if order, convErr := strconv.Atoi(rule.RuleOrder); convErr == nil {
			currentOrderMap[rule.ID] = order
		}
	}

	for id := range orders.Orders {
		if order, exists := currentOrderMap[id]; exists {
			orders.Orders[id] = order
		}
	}

	data.ID = types.StringValue(orders.PolicyType + "-reorder")

	ruleObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"order": types.StringType,
		},
	}

	ruleModels := make([]PolicyAccessRuleReorderRuleModel, 0, len(orders.Orders))
	for id, order := range orders.Orders {
		ruleModels = append(ruleModels, PolicyAccessRuleReorderRuleModel{
			ID:    types.StringValue(id),
			Order: types.StringValue(strconv.Itoa(order)),
		})
	}

	sort.Slice(ruleModels, func(i, j int) bool {
		return ruleModels[i].ID.ValueString() < ruleModels[j].ID.ValueString()
	})

	setValue, setDiags := types.SetValueFrom(ctx, ruleObjectType, ruleModels)
	diags.Append(setDiags...)
	if diags.HasError() {
		return diags
	}

	data.Rules = setValue
	return diags
}

func (r *PolicyAccessRuleReorderResource) getRules(ctx context.Context, data *PolicyAccessRuleReorderModel) (*rulesOrders, diag.Diagnostics) {
	var diags diag.Diagnostics

	if data.Rules.IsNull() || data.Rules.IsUnknown() {
		diags.AddError("Missing Rules", "The rules set must be provided.")
		return nil, diags
	}

	ruleModels := make([]PolicyAccessRuleReorderRuleModel, 0)
	diags.Append(data.Rules.ElementsAs(ctx, &ruleModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := &rulesOrders{
		PolicyType: data.PolicyType.ValueString(),
		Orders:     make(map[string]int, len(ruleModels)),
	}

	for _, rule := range ruleModels {
		if rule.ID.IsNull() || rule.ID.IsUnknown() || strings.TrimSpace(rule.ID.ValueString()) == "" {
			diags.AddError("Invalid Rule", "Each rule must include a non-empty id.")
			continue
		}

		if rule.Order.IsNull() || rule.Order.IsUnknown() || strings.TrimSpace(rule.Order.ValueString()) == "" {
			diags.AddError("Invalid Rule Order", fmt.Sprintf("Rule %s must include an order value.", rule.ID.ValueString()))
			continue
		}

		order, err := strconv.Atoi(strings.TrimSpace(rule.Order.ValueString()))
		if err != nil {
			diags.AddError("Invalid Rule Order", fmt.Sprintf("Rule %s has an invalid order value: %v", rule.ID.ValueString(), err))
			continue
		}

		result.Orders[rule.ID.ValueString()] = order
	}

	return result, diags
}

func validateRuleOrders(orders *rulesOrders) error {
	for _, order := range orders.Orders {
		if order <= 0 {
			return fmt.Errorf("order must be a positive integer greater than 0")
		}
	}

	if dupOrder, dupRuleIDs, ok := hasDuplicates(orders.Orders); ok {
		return fmt.Errorf("duplicate order '%d' used by rules with IDs: %s", dupOrder, strings.Join(dupRuleIDs, ", "))
	}

	return nil
}

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
