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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

var (
	_ resource.Resource                = &PRACredentialControllerResource{}
	_ resource.ResourceWithConfigure   = &PRACredentialControllerResource{}
	_ resource.ResourceWithImportState = &PRACredentialControllerResource{}
)

var praCredentialPolicyLock sync.Mutex

func NewPRACredentialControllerResource() resource.Resource {
	return &PRACredentialControllerResource{}
}

type PRACredentialControllerResource struct {
	client *client.Client
}

type PRACredentialControllerResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	CredentialType  types.String `tfsdk:"credential_type"`
	Passphrase      types.String `tfsdk:"passphrase"`
	Password        types.String `tfsdk:"password"`
	PrivateKey      types.String `tfsdk:"private_key"`
	UserDomain      types.String `tfsdk:"user_domain"`
	Username        types.String `tfsdk:"username"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	MicroTenantName types.String `tfsdk:"microtenant_name"`
}

func (r *PRACredentialControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_credential_controller"
}

func (r *PRACredentialControllerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA PRA credential controller.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the credential controller.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the credential controller.",
			},
			"credential_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("USERNAME_PASSWORD", "SSH_KEY", "PASSWORD"),
				},
				Description: "Credential type. Supported values: USERNAME_PASSWORD, SSH_KEY, PASSWORD.",
			},
			"passphrase": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase protecting the SSH private key.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password associated with the credential.",
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "SSH private key associated with the credential.",
			},
			"user_domain": schema.StringAttribute{
				Optional:    true,
				Description: "Domain name associated with the username.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username associated with the credential.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID to scope the resource.",
			},
			"microtenant_name": schema.StringAttribute{
				Computed:    true,
				Description: "Micro-tenant name.",
			},
		},
	}
}

func (r *PRACredentialControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PRACredentialControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before creating this resource.",
		)
		return
	}

	var plan PRACredentialControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload := expandPRACredentialController(plan)

	created, _, err := pracredential.Create(ctx, service, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create PRA credential controller: %v", err))
		return
	}

	tflog.Info(ctx, "Created PRA credential controller", map[string]any{"id": created.ID})

	state, diags := r.readPRACredential(ctx, service, created.ID, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sensitive fields and user_domain from plan as they are not returned by API
	state.Passphrase = plan.Passphrase
	state.Password = plan.Password
	state.PrivateKey = plan.PrivateKey
	state.UserDomain = plan.UserDomain

	// Preserve username from plan if it was set, otherwise use what API returned
	// This handles the case where API might return empty string but plan had null
	if !plan.Username.IsNull() && !plan.Username.IsUnknown() {
		state.Username = plan.Username
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PRACredentialControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this resource.",
		)
		return
	}

	var state PRACredentialControllerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "PRA credential controller ID is required to read the resource.")
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	currentState, diags := r.readPRACredential(ctx, service, state.ID.ValueString(), state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if currentState.ID.IsNull() || currentState.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &currentState)...)
}

func (r *PRACredentialControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before updating this resource.",
		)
		return
	}

	var plan PRACredentialControllerResourceModel
	var state PRACredentialControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.CredentialType.IsNull() && !plan.CredentialType.IsNull() && !strings.EqualFold(plan.CredentialType.ValueString(), state.CredentialType.ValueString()) {
		resp.Diagnostics.AddError("Invalid Update", "Changing 'credential_type' is not supported.")
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	existing, _, err := pracredential.Get(ctx, service, state.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PRA credential controller: %v", err))
		return
	}

	payload := expandPRACredentialController(plan)
	if _, err := pracredential.Update(ctx, service, existing.ID, &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update PRA credential controller: %v", err))
		return
	}

	updatedState, diags := r.readPRACredential(ctx, service, existing.ID, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sensitive fields and user_domain from plan if they were provided, otherwise from state
	if !plan.Passphrase.IsNull() && !plan.Passphrase.IsUnknown() {
		updatedState.Passphrase = plan.Passphrase
	} else {
		updatedState.Passphrase = state.Passphrase
	}
	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		updatedState.Password = plan.Password
	} else {
		updatedState.Password = state.Password
	}
	if !plan.PrivateKey.IsNull() && !plan.PrivateKey.IsUnknown() {
		updatedState.PrivateKey = plan.PrivateKey
	} else {
		updatedState.PrivateKey = state.PrivateKey
	}
	if !plan.UserDomain.IsNull() && !plan.UserDomain.IsUnknown() {
		updatedState.UserDomain = plan.UserDomain
	} else {
		updatedState.UserDomain = state.UserDomain
	}

	// Preserve username from plan if it was set, otherwise use what API returned
	if !plan.Username.IsNull() && !plan.Username.IsUnknown() {
		updatedState.Username = plan.Username
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *PRACredentialControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before deleting this resource.",
		)
		return
	}

	var state PRACredentialControllerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	if err := detachCredentialFromPolicies(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Detach Error", fmt.Sprintf("Unable to detach PRA credential controller from policies: %v", err))
		return
	}

	if _, err := pracredential.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete PRA credential controller: %v", err))
		return
	}
}

func (r *PRACredentialControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PRACredentialControllerResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(strings.TrimSpace(microtenantID.ValueString()))
	}
	return service
}

func (r *PRACredentialControllerResource) readPRACredential(ctx context.Context, service *zscaler.Service, id string, existingState PRACredentialControllerResourceModel) (PRACredentialControllerResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state PRACredentialControllerResourceModel

	credential, _, err := pracredential.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			state.ID = types.StringNull()
			return state, diags
		}

		diags.AddError("Client Error", fmt.Sprintf("Unable to read PRA credential controller: %v", err))
		return state, diags
	}

	state = flattenPRACredentialControllerResource(credential, existingState)
	return state, diags
}

func expandPRACredentialController(plan PRACredentialControllerResourceModel) pracredential.Credential {
	return pracredential.Credential{
		ID:             plan.ID.ValueString(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		CredentialType: plan.CredentialType.ValueString(),
		Passphrase:     plan.Passphrase.ValueString(),
		Password:       plan.Password.ValueString(),
		PrivateKey:     plan.PrivateKey.ValueString(),
		UserDomain:     plan.UserDomain.ValueString(),
		UserName:       plan.Username.ValueString(),
		MicroTenantID:  plan.MicroTenantID.ValueString(),
	}
}

func flattenPRACredentialControllerResource(cred *pracredential.Credential, existingState PRACredentialControllerResourceModel) PRACredentialControllerResourceModel {
	// Preserve sensitive fields from existing state as they are not returned by the API
	// This matches SDKv2 behavior where password, passphrase, and private_key are never set in Read
	passphrase := existingState.Passphrase
	password := existingState.Password
	privateKey := existingState.PrivateKey

	// Preserve user_domain from existing state - SDKv2 does not set this field in Read
	userDomain := existingState.UserDomain

	// Username is returned by API in SDKv2, so update it if returned and non-empty
	// If API returns empty string but existingState has a value (null or otherwise), preserve existingState
	username := existingState.Username
	if cred.UserName != "" {
		username = types.StringValue(cred.UserName)
	} else if existingState.Username.IsNull() || existingState.Username.IsUnknown() {
		// If API returns empty and existingState is null/unknown, set to null (not empty string)
		username = types.StringNull()
	}

	return PRACredentialControllerResourceModel{
		ID:              types.StringValue(cred.ID),
		Name:            types.StringValue(cred.Name),
		Description:     types.StringValue(cred.Description),
		CredentialType:  types.StringValue(cred.CredentialType),
		Passphrase:      passphrase,
		Password:        password,
		PrivateKey:      privateKey,
		UserDomain:      userDomain,
		Username:        username,
		MicroTenantID:   types.StringValue(cred.MicroTenantID),
		MicroTenantName: types.StringValue(cred.MicroTenantName),
	}
}

func detachCredentialFromPolicies(ctx context.Context, service *zscaler.Service, credentialID string) error {
	praCredentialPolicyLock.Lock()
	defer praCredentialPolicyLock.Unlock()

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
