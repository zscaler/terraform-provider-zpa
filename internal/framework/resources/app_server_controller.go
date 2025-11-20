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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &AppServerControllerResource{}
	_ resource.ResourceWithConfigure   = &AppServerControllerResource{}
	_ resource.ResourceWithImportState = &AppServerControllerResource{}
)

func NewAppServerControllerResource() resource.Resource {
	return &AppServerControllerResource{}
}

type AppServerControllerResource struct {
	client *client.Client
}

type AppServerControllerResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Address           types.String `tfsdk:"address"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	AppServerGroupIDs types.Set    `tfsdk:"app_server_group_ids"`
	ConfigSpace       types.String `tfsdk:"config_space"`
	MicroTenantID     types.String `tfsdk:"microtenant_id"`
	MicroTenantName   types.String `tfsdk:"microtenant_name"`
	CreationTime      types.String `tfsdk:"creation_time"`
	ModifiedBy        types.String `tfsdk:"modifiedby"`
	ModifiedTime      types.String `tfsdk:"modified_time"`
}

func (r *AppServerControllerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_server"
}

func (r *AppServerControllerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Application Server (App Connector Server)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the application server.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the application server.",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "Domain or IP address of the application server.",
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"app_server_group_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"config_space": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("DEFAULT", "SIEM"),
				},
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

func (r *AppServerControllerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppServerControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppServerControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	payload, diags := expandApplicationServer(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, _, err := appservercontroller.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create application server: %v", err))
		return
	}

	plan.ID = types.StringValue(created.ID)

	state, diags := r.readApplicationServer(ctx, service, created.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppServerControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppServerControllerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "Application server ID is required to read the resource")
		return
	}

	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	newState, diags := r.readApplicationServer(ctx, service, state.ID.ValueString())
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

func (r *AppServerControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppServerControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !plan.MicroTenantID.IsNull() && plan.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(plan.MicroTenantID.ValueString())
	}

	// Check if resource still exists before updating
	if _, _, err := appservercontroller.Get(ctx, service, plan.ID.ValueString()); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	payload, diags := expandApplicationServer(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := appservercontroller.Update(ctx, service, plan.ID.ValueString(), payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update application server: %v", err))
		return
	}

	state, stateDiags := r.readApplicationServer(ctx, service, plan.ID.ValueString())
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

func (r *AppServerControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppServerControllerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.client.Service
	if !state.MicroTenantID.IsNull() && state.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(state.MicroTenantID.ValueString())
	}

	if _, err := appservercontroller.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete application server: %v", err))
		return
	}
}

func (r *AppServerControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *AppServerControllerResource) readApplicationServer(ctx context.Context, service *zscaler.Service, id string) (AppServerControllerResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state AppServerControllerResourceModel

	server, _, err := appservercontroller.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			return AppServerControllerResourceModel{}, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read application server: %v", err))
		return state, diags
	}

	model, flattenDiags := flattenAppServer(ctx, server)
	diags.Append(flattenDiags...)

	state = AppServerControllerResourceModel(model)
	return state, diags
}

func expandApplicationServer(ctx context.Context, plan AppServerControllerResourceModel) (appservercontroller.ApplicationServer, diag.Diagnostics) {
	var diags diag.Diagnostics

	groupIDs := make([]string, 0)
	if !plan.AppServerGroupIDs.IsNull() && !plan.AppServerGroupIDs.IsUnknown() {
		var values []string
		convertDiags := plan.AppServerGroupIDs.ElementsAs(ctx, &values, false)
		diags.Append(convertDiags...)
		if !diags.HasError() {
			groupIDs = values
		}
	}

	payload := appservercontroller.ApplicationServer{
		ID:                plan.ID.ValueString(),
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		Address:           plan.Address.ValueString(),
		Enabled:           helpers.BoolValue(plan.Enabled, true),
		AppServerGroupIds: groupIDs,
		ConfigSpace:       plan.ConfigSpace.ValueString(),
		MicroTenantID:     plan.MicroTenantID.ValueString(),
	}

	return payload, diags
}

func flattenAppServer(ctx context.Context, server *appservercontroller.ApplicationServer) (AppServerControllerResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	groupIDs, setDiags := types.SetValueFrom(ctx, types.StringType, server.AppServerGroupIds)
	diags.Append(setDiags...)

	model := AppServerControllerResourceModel{
		ID:                types.StringValue(server.ID),
		Name:              types.StringValue(server.Name),
		Description:       types.StringValue(server.Description),
		Address:           types.StringValue(server.Address),
		Enabled:           types.BoolValue(server.Enabled),
		AppServerGroupIDs: groupIDs,
		ConfigSpace:       types.StringValue(server.ConfigSpace),
		MicroTenantID:     types.StringValue(server.MicroTenantID),
		MicroTenantName:   types.StringValue(server.MicroTenantName),
		CreationTime:      types.StringValue(server.CreationTime),
		ModifiedBy:        types.StringValue(server.ModifiedBy),
		ModifiedTime:      types.StringValue(server.ModifiedTime),
	}

	return model, diags
}
