package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ProvisioningKeyResource{}
	_ resource.ResourceWithConfigure   = &ProvisioningKeyResource{}
	_ resource.ResourceWithImportState = &ProvisioningKeyResource{}
)

func NewProvisioningKeyResource() resource.Resource {
	return &ProvisioningKeyResource{}
}

type ProvisioningKeyResource struct {
	client *client.Client
}

type ProvisioningKeyResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	AssociationType       types.String `tfsdk:"association_type"`
	AppConnectorGroupID   types.String `tfsdk:"app_connector_group_id"`
	AppConnectorGroupName types.String `tfsdk:"app_connector_group_name"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	MaxUsage              types.String `tfsdk:"max_usage"`
	EnrollmentCertID      types.String `tfsdk:"enrollment_cert_id"`
	// UIConfig              types.String `tfsdk:"ui_config"`
	UsageCount      types.String `tfsdk:"usage_count"`
	ZComponentID    types.String `tfsdk:"zcomponent_id"`
	ZComponentName  types.String `tfsdk:"zcomponent_name"`
	ProvisioningKey types.String `tfsdk:"provisioning_key"`
	IPAcl           types.Set    `tfsdk:"ip_acl"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	MicroTenantName types.String `tfsdk:"microtenant_name"`
	CreationTime    types.String `tfsdk:"creation_time"`
	ModifiedBy      types.String `tfsdk:"modifiedby"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
}

func (r *ProvisioningKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provisioning_key"
}

func (r *ProvisioningKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA provisioning key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the provisioning key.",
			},
			"association_type": schema.StringAttribute{
				Required:    true,
				Description: "Provisioning key association type.",
				Validators: []validator.String{
					stringvalidator.OneOf(provisioningkey.ProvisioningKeyAssociationTypes...),
				},
			},
			"app_connector_group_id": schema.StringAttribute{
				Optional: true,
			},
			"app_connector_group_name": schema.StringAttribute{Computed: true},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"max_usage": schema.StringAttribute{
				Required: true,
			},
			"enrollment_cert_id": schema.StringAttribute{
				Required: true,
			},
			// "ui_config":   schema.StringAttribute{Optional: true},
			"usage_count": schema.StringAttribute{Computed: true},
			"zcomponent_id": schema.StringAttribute{
				Required: true,
			},
			"zcomponent_name": schema.StringAttribute{Computed: true},
			"provisioning_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "For backward compatibility. The value is cleared after reads; use the `zpa_provisioning_key` ephemeral resource to retrieve the live key value.",
			},
			"ip_acl": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"microtenant_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"creation_time":    schema.StringAttribute{Computed: true},
			"modifiedby":       schema.StringAttribute{Computed: true},
			"modified_time":    schema.StringAttribute{Computed: true},
		},
	}
}

func (r *ProvisioningKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProvisioningKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProvisioningKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	associationType := plan.AssociationType.ValueString()
	payload, diags := expandProvisioningKey(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := provisioningkey.Create(ctx, service, associationType, &payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create provisioning key: %v", err))
		return
	}

	plan.ID = types.StringValue(created.ID)

	state, stateDiags := r.readProvisioningKey(ctx, service, associationType, created.ID)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProvisioningKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProvisioningKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "Provisioning key ID is required to read the resource")
		return
	}

	associationType := state.AssociationType.ValueString()
	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	newState, diags := r.readProvisioningKey(ctx, service, associationType, state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState.ID.IsNull() || newState.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ProvisioningKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProvisioningKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	associationType := plan.AssociationType.ValueString()

	// Check if resource still exists before updating
	if _, _, err := provisioningkey.Get(ctx, service, associationType, plan.ID.ValueString()); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	payload, diags := expandProvisioningKey(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := provisioningkey.Update(ctx, service, associationType, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update provisioning key: %v", err))
		return
	}

	state, stateDiags := r.readProvisioningKey(ctx, service, associationType, plan.ID.ValueString())
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProvisioningKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProvisioningKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	associationType := state.AssociationType.ValueString()
	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	if _, err := provisioningkey.Delete(ctx, service, associationType, state.ID.ValueString()); err != nil {
		// If resource is already deleted (not found), that's fine
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete provisioning key: %v", err))
		return
	}
}

func (r *ProvisioningKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	service := r.client.Service

	key, associationType, _, err := provisioningkey.GetByIDAllAssociations(ctx, service, id)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to lookup provisioning key association type for %s: %v", id, err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("association_type"), types.StringValue(associationType))...)

	if key != nil {
		if key.MicroTenantID != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("microtenant_id"), types.StringValue(key.MicroTenantID))...)
		}
	}
}

func (r *ProvisioningKeyResource) readProvisioningKey(ctx context.Context, service *zscaler.Service, associationType, id string) (ProvisioningKeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state ProvisioningKeyResourceModel

	key, _, err := provisioningkey.Get(ctx, service, associationType, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			return ProvisioningKeyResourceModel{}, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read provisioning key: %v", err))
		return state, diags
	}

	flattened, flattenDiags := flattenProvisioningKey(ctx, key, associationType)
	diags.Append(flattenDiags...)
	state = ProvisioningKeyResourceModel(flattened)
	return state, diags
}

func expandProvisioningKey(ctx context.Context, plan ProvisioningKeyResourceModel) (provisioningkey.ProvisioningKey, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipACL []string
	if !plan.IPAcl.IsNull() && !plan.IPAcl.IsUnknown() {
		elementsDiags := plan.IPAcl.ElementsAs(ctx, &ipACL, false)
		diags.Append(elementsDiags...)
	}

	payload := provisioningkey.ProvisioningKey{
		ID:                    plan.ID.ValueString(),
		Name:                  plan.Name.ValueString(),
		AppConnectorGroupID:   plan.AppConnectorGroupID.ValueString(),
		AppConnectorGroupName: plan.AppConnectorGroupName.ValueString(),
		Enabled:               helpers.BoolValue(plan.Enabled, true),
		MaxUsage:              plan.MaxUsage.ValueString(),
		EnrollmentCertID:      plan.EnrollmentCertID.ValueString(),
		// UIConfig:              plan.UIConfig.ValueString(),
		UsageCount:      plan.UsageCount.ValueString(),
		ZcomponentID:    plan.ZComponentID.ValueString(),
		ZcomponentName:  plan.ZComponentName.ValueString(),
		ProvisioningKey: plan.ProvisioningKey.ValueString(),
		MicroTenantID:   plan.MicroTenantID.ValueString(),
		IPACL:           ipACL,
	}

	return payload, diags
}

func flattenProvisioningKey(ctx context.Context, key *provisioningkey.ProvisioningKey, associationType string) (ProvisioningKeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	ipAcl, setDiags := types.SetValueFrom(ctx, types.StringType, key.IPACL)
	diags.Append(setDiags...)

	model := ProvisioningKeyResourceModel{
		ID:                    helpers.StringValueOrNull(key.ID),
		Name:                  helpers.StringValueOrNull(key.Name),
		AssociationType:       types.StringValue(associationType),
		AppConnectorGroupID:   helpers.StringValueOrNull(key.AppConnectorGroupID),
		AppConnectorGroupName: helpers.StringValueOrNull(key.AppConnectorGroupName),
		Enabled:               types.BoolValue(key.Enabled),
		MaxUsage:              helpers.StringValueOrNull(key.MaxUsage),
		EnrollmentCertID:      helpers.StringValueOrNull(key.EnrollmentCertID),
		// UIConfig:              types.StringValue(key.UIConfig),
		UsageCount:      helpers.StringValueOrNull(key.UsageCount),
		ZComponentID:    helpers.StringValueOrNull(key.ZcomponentID),
		ZComponentName:  helpers.StringValueOrNull(key.ZcomponentName),
		ProvisioningKey: helpers.StringValueOrNull(key.ProvisioningKey),
		IPAcl:           ipAcl,
		MicroTenantID:   helpers.StringValueOrNull(key.MicroTenantID),
		MicroTenantName: helpers.StringValueOrNull(key.MicroTenantName),
		CreationTime:    helpers.StringValueOrNull(key.CreationTime),
		ModifiedBy:      helpers.StringValueOrNull(key.ModifiedBy),
		ModifiedTime:    helpers.StringValueOrNull(key.ModifiedTime),
	}

	return model, diags
}
