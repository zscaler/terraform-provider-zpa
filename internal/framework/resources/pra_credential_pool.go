package resources

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

var (
	_ resource.Resource                = &PRACredentialPoolResource{}
	_ resource.ResourceWithConfigure   = &PRACredentialPoolResource{}
	_ resource.ResourceWithImportState = &PRACredentialPoolResource{}
)

var praCredentialPoolLock sync.Mutex

func NewPRACredentialPoolResource() resource.Resource {
	return &PRACredentialPoolResource{}
}

type PRACredentialPoolResource struct {
	client *client.Client
}

type PRACredentialPoolModel struct {
	ID             types.String             `tfsdk:"id"`
	Name           types.String             `tfsdk:"name"`
	CredentialType types.String             `tfsdk:"credential_type"`
	Credentials    []PRACredentialReference `tfsdk:"credentials"`
	MicrotenantID  types.String             `tfsdk:"microtenant_id"`
}

type PRACredentialReference struct {
	IDs types.Set `tfsdk:"id"`
}

func (r *PRACredentialPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_credential_pool"
}

func (r *PRACredentialPoolResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA PRA credential pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the credential pool.",
			},
			"credential_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("USERNAME_PASSWORD", "SSH_KEY", "PASSWORD"),
				},
				Description: "Credential type. Supported values: USERNAME_PASSWORD, SSH_KEY, PASSWORD.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID for scoping.",
			},
		},
		Blocks: map[string]schema.Block{
			"credentials": schema.ListNestedBlock{
				Description: "List of PRA credential IDs in the pool.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *PRACredentialPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *PRACredentialPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA credential pools.")
		return
	}

	var plan PRACredentialPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPRACredentialPool(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := pracredentialpool.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create PRA credential pool: %v", err))
		return
	}

	state, readDiags := r.readCredentialPool(ctx, service, created.ID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRACredentialPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA credential pools.")
		return
	}

	var state PRACredentialPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	newState, diags := r.readCredentialPool(ctx, service, state.ID.ValueString())
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(d.Detail()), "not found") {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PRACredentialPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA credential pools.")
		return
	}

	var plan PRACredentialPoolModel
	var state PRACredentialPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.CredentialType.IsNull() && !plan.CredentialType.IsNull() &&
		!strings.EqualFold(plan.CredentialType.ValueString(), state.CredentialType.ValueString()) {
		resp.Diagnostics.AddError("Invalid Update", "Changing 'credential_type' is not supported.")
		return
	}

	service := r.serviceForMicrotenant(plan.MicrotenantID)

	payload, diags := expandPRACredentialPool(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := pracredentialpool.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update PRA credential pool: %v", err))
		return
	}

	newState, readDiags := r.readCredentialPool(ctx, service, plan.ID.ValueString())
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *PRACredentialPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before managing PRA credential pools.")
		return
	}

	var state PRACredentialPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicrotenantID)

	if err := detachPRACredentialFromPolicies(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Detach Error", fmt.Sprintf("Failed to detach PRA credential pool from policies: %v", err))
		return
	}

	if _, err := pracredentialpool.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete PRA credential pool: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *PRACredentialPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PRACredentialPoolResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(strings.TrimSpace(microtenantID.ValueString()))
	}
	return service
}

func (r *PRACredentialPoolResource) readCredentialPool(ctx context.Context, service *zscaler.Service, id string) (PRACredentialPoolModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	pool, _, err := pracredentialpool.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return PRACredentialPoolModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("PRA credential pool %s not found", id))}
		}
		return PRACredentialPoolModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read PRA credential pool: %v", err))}
	}

	model, flattenDiags := flattenPRACredentialPool(ctx, pool)
	diags.Append(flattenDiags...)
	return model, diags
}

func expandPRACredentialPool(ctx context.Context, model PRACredentialPoolModel) (pracredentialpool.CredentialPool, diag.Diagnostics) {
	var diags diag.Diagnostics

	credentials, credsDiags := expandPRACredentials(ctx, model.Credentials)
	diags.Append(credsDiags...)

	return pracredentialpool.CredentialPool{
		ID:             model.ID.ValueString(),
		Name:           model.Name.ValueString(),
		CredentialType: model.CredentialType.ValueString(),
		PRACredentials: credentials,
		MicroTenantID:  model.MicrotenantID.ValueString(),
	}, diags
}

func expandPRACredentials(ctx context.Context, models []PRACredentialReference) ([]common.CommonIDName, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(models) == 0 {
		return nil, diags
	}

	var result []common.CommonIDName
	for _, model := range models {
		if model.IDs.IsNull() || model.IDs.IsUnknown() {
			continue
		}
		var ids []string
		diags.Append(model.IDs.ElementsAs(ctx, &ids, false)...)
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			result = append(result, common.CommonIDName{ID: id})
		}
	}

	return result, diags
}

func flattenPRACredentialPool(ctx context.Context, pool *pracredentialpool.CredentialPool) (PRACredentialPoolModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	credentials, credsDiags := flattenPRACredentials(ctx, pool.PRACredentials)
	diags.Append(credsDiags...)

	return PRACredentialPoolModel{
		ID:             types.StringValue(pool.ID),
		Name:           types.StringValue(pool.Name),
		CredentialType: types.StringValue(pool.CredentialType),
		Credentials:    credentials,
		MicrotenantID:  types.StringValue(pool.MicroTenantID),
	}, diags
}

func flattenPRACredentials(ctx context.Context, creds []common.CommonIDName) ([]PRACredentialReference, diag.Diagnostics) {
	if len(creds) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(creds))
	for _, cred := range creds {
		if strings.TrimSpace(cred.ID) != "" {
			ids = append(ids, cred.ID)
		}
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return nil, diags
	}

	return []PRACredentialReference{{IDs: setValue}}, nil
}

func detachPRACredentialFromPolicies(ctx context.Context, service *zscaler.Service, credentialID string) error {
	praCredentialPoolLock.Lock()
	defer praCredentialPoolLock.Unlock()

	policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "CREDENTIAL_POLICY")
	if err != nil {
		return fmt.Errorf("failed to get CREDENTIAL_POLICY policy set: %w", err)
	}

	rules, _, err := policysetcontroller.GetAllByType(ctx, service, "CREDENTIAL_POLICY")
	if err != nil {
		return fmt.Errorf("failed to get CREDENTIAL_POLICY rules: %w", err)
	}

	for _, rule := range rules {
		updated := false

		for i, condition := range rule.Conditions {
			newOperands := make([]policysetcontroller.Operands, 0, len(condition.Operands))
			for _, operand := range condition.Operands {
				if operand.ObjectType == "CREDENTIAL" && operand.LHS == "id" && operand.RHS == credentialID {
					updated = true
					continue
				}
				newOperands = append(newOperands, operand)
			}
			rule.Conditions[i].Operands = newOperands
		}

		if updated {
			if _, err := policysetcontroller.UpdateRule(ctx, service, policySet.ID, rule.ID, &rule); err != nil {
				return fmt.Errorf("failed to update rule %s: %w", rule.ID, err)
			}
		}
	}

	return nil
}
