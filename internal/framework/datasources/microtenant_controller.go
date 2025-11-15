package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

var (
	_ datasource.DataSource              = &MicrotenantControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &MicrotenantControllerDataSource{}
)

func NewMicrotenantControllerDataSource() datasource.DataSource {
	return &MicrotenantControllerDataSource{}
}

type MicrotenantControllerDataSource struct {
	client *client.Client
}

type MicrotenantControllerModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	Enabled                    types.Bool   `tfsdk:"enabled"`
	CriteriaAttribute          types.String `tfsdk:"criteria_attribute"`
	CriteriaAttributeValues    types.List   `tfsdk:"criteria_attribute_values"`
	Operator                   types.String `tfsdk:"operator"`
	Priority                   types.String `tfsdk:"priority"`
	PrivilegedApprovalsEnabled types.Bool   `tfsdk:"privileged_approvals_enabled"`
	CreationTime               types.String `tfsdk:"creation_time"`
	ModifiedBy                 types.String `tfsdk:"modified_by"`
	ModifiedTime               types.String `tfsdk:"modified_time"`
	Roles                      types.Set    `tfsdk:"roles"`
	Users                      types.Set    `tfsdk:"user"`
}

func (d *MicrotenantControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_microtenant_controller"
}

func (d *MicrotenantControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA micro-tenant by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the micro-tenant.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the micro-tenant.",
			},
			"description":        schema.StringAttribute{Computed: true},
			"enabled":            schema.BoolAttribute{Computed: true},
			"criteria_attribute": schema.StringAttribute{Computed: true},
			"criteria_attribute_values": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"operator":                     schema.StringAttribute{Computed: true},
			"priority":                     schema.StringAttribute{Computed: true},
			"privileged_approvals_enabled": schema.BoolAttribute{Computed: true},
			"creation_time":                schema.StringAttribute{Computed: true},
			"modified_by":                  schema.StringAttribute{Computed: true},
			"modified_time":                schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"roles": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"custom_role": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"user": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":               schema.StringAttribute{Computed: true},
						"name":             schema.StringAttribute{Computed: true},
						"description":      schema.StringAttribute{Computed: true},
						"enabled":          schema.BoolAttribute{Computed: true},
						"comments":         schema.StringAttribute{Computed: true},
						"customer_id":      schema.StringAttribute{Computed: true},
						"display_name":     schema.StringAttribute{Computed: true},
						"delivery_tag":     schema.StringAttribute{Computed: true},
						"email":            schema.StringAttribute{Computed: true},
						"eula":             schema.StringAttribute{Computed: true},
						"force_pwd_change": schema.BoolAttribute{Computed: true},
						"group_ids": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"iam_user_id":             schema.StringAttribute{Computed: true},
						"is_enabled":              schema.BoolAttribute{Computed: true},
						"is_locked":               schema.BoolAttribute{Computed: true},
						"language_code":           schema.StringAttribute{Computed: true},
						"local_login_disabled":    schema.BoolAttribute{Computed: true},
						"password":                schema.StringAttribute{Computed: true},
						"one_identity_user":       schema.BoolAttribute{Computed: true},
						"operation_type":          schema.StringAttribute{Computed: true},
						"phone_number":            schema.BoolAttribute{Computed: true},
						"pin_session":             schema.StringAttribute{Computed: true},
						"role_id":                 schema.BoolAttribute{Computed: true},
						"sync_version":            schema.StringAttribute{Computed: true},
						"microtenant_id":          schema.StringAttribute{Computed: true},
						"microtenant_name":        schema.StringAttribute{Computed: true},
						"timezone":                schema.StringAttribute{Computed: true},
						"tmp_password":            schema.StringAttribute{Computed: true},
						"token_id":                schema.StringAttribute{Computed: true},
						"two_factor_auth_enabled": schema.BoolAttribute{Computed: true},
						"two_factor_auth_type":    schema.StringAttribute{Computed: true},
						"username":                schema.StringAttribute{Computed: true},
						"creation_time":           schema.StringAttribute{Computed: true},
						"modifiedby":              schema.StringAttribute{Computed: true},
						"modified_time":           schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *MicrotenantControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *MicrotenantControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data MicrotenantControllerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var tenant *microtenants.MicroTenant
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving micro-tenant by ID", map[string]any{"id": id})
		tenant, _, err = microtenants.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving micro-tenant by name", map[string]any{"name": name})
		tenant, _, err = microtenants.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read micro-tenant: %v", err))
		return
	}

	if tenant == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Micro-tenant with id %q or name %q was not found.", id, name))
		return
	}

	state, diags := flattenMicrotenant(ctx, tenant)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMicrotenant(ctx context.Context, tenant *microtenants.MicroTenant) (MicrotenantControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	criteriaValues, listDiags := types.ListValueFrom(ctx, types.StringType, tenant.CriteriaAttributeValues)
	diags.Append(listDiags...)

	roles, roleDiags := flattenMicrotenantRoles(ctx, tenant.Roles)
	diags.Append(roleDiags...)

	users, userDiags := flattenMicrotenantUser(ctx, tenant.UserResource)
	diags.Append(userDiags...)

	state := MicrotenantControllerModel{
		ID:                         types.StringValue(tenant.ID),
		Name:                       types.StringValue(tenant.Name),
		Description:                types.StringValue(tenant.Description),
		Enabled:                    types.BoolValue(tenant.Enabled),
		CriteriaAttribute:          types.StringValue(tenant.CriteriaAttribute),
		CriteriaAttributeValues:    criteriaValues,
		Operator:                   types.StringValue(tenant.Operator),
		Priority:                   types.StringValue(tenant.Priority),
		PrivilegedApprovalsEnabled: types.BoolValue(tenant.PrivilegedApprovalsEnabled),
		CreationTime:               types.StringValue(tenant.CreationTime),
		ModifiedBy:                 types.StringValue(tenant.ModifiedBy),
		ModifiedTime:               types.StringValue(tenant.ModifiedTime),
		Roles:                      roles,
		Users:                      users,
	}

	return state, diags
}

func flattenMicrotenantRoles(ctx context.Context, roles []microtenants.Roles) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"custom_role": types.BoolType,
	}

	if len(roles) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(roles))
	for _, role := range roles {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":          types.StringValue(role.ID),
			"name":        types.StringValue(role.Name),
			"custom_role": types.BoolValue(role.CustomRole),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(setDiags...)
	return set, diags
}

func flattenMicrotenantUser(ctx context.Context, user *microtenants.UserResource) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	groupIDs, groupDiags := types.SetValueFrom(ctx, types.StringType, user.GroupIDs)
	diags.Append(groupDiags...)

	attrTypes := map[string]attr.Type{
		"id":                      types.StringType,
		"name":                    types.StringType,
		"description":             types.StringType,
		"enabled":                 types.BoolType,
		"comments":                types.StringType,
		"customer_id":             types.StringType,
		"display_name":            types.StringType,
		"delivery_tag":            types.StringType,
		"email":                   types.StringType,
		"eula":                    types.StringType,
		"force_pwd_change":        types.BoolType,
		"group_ids":               types.SetType{ElemType: types.StringType},
		"iam_user_id":             types.StringType,
		"is_enabled":              types.BoolType,
		"is_locked":               types.BoolType,
		"language_code":           types.StringType,
		"local_login_disabled":    types.BoolType,
		"password":                types.StringType,
		"one_identity_user":       types.BoolType,
		"operation_type":          types.StringType,
		"phone_number":            types.BoolType,
		"pin_session":             types.StringType,
		"role_id":                 types.BoolType,
		"sync_version":            types.StringType,
		"microtenant_id":          types.StringType,
		"microtenant_name":        types.StringType,
		"timezone":                types.StringType,
		"tmp_password":            types.StringType,
		"token_id":                types.StringType,
		"two_factor_auth_enabled": types.BoolType,
		"two_factor_auth_type":    types.StringType,
		"username":                types.StringType,
		"creation_time":           types.StringType,
		"modifiedby":              types.StringType,
		"modified_time":           types.StringType,
	}

	if user == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":                      helpers.StringValueOrNull(user.ID),
		"name":                    helpers.StringValueOrNull(user.Name),
		"description":             helpers.StringValueOrNull(user.Description),
		"enabled":                 types.BoolValue(user.Enabled),
		"comments":                helpers.StringValueOrNull(user.Comments),
		"customer_id":             helpers.StringValueOrNull(user.CustomerID),
		"display_name":            helpers.StringValueOrNull(user.DisplayName),
		"delivery_tag":            helpers.StringValueOrNull(user.DeliveryTag),
		"email":                   helpers.StringValueOrNull(user.Email),
		"eula":                    helpers.StringValueOrNull(user.Eula),
		"force_pwd_change":        types.BoolValue(user.ForcePwdChange),
		"group_ids":               groupIDs,
		"iam_user_id":             helpers.StringValueOrNull(user.IAMUserID),
		"is_enabled":              types.BoolValue(user.IsEnabled),
		"is_locked":               types.BoolValue(user.IsLocked),
		"language_code":           helpers.StringValueOrNull(user.LanguageCode),
		"local_login_disabled":    types.BoolValue(user.LocalLoginDisabled),
		"password":                helpers.StringValueOrNull(user.Password),
		"one_identity_user":       types.BoolValue(user.OneIdentityUser),
		"operation_type":          helpers.StringValueOrNull(user.OperationType),
		"phone_number":            parseBoolFromString(user.PhoneNumber),
		"pin_session":             boolToString(user.PinSession),
		"role_id":                 parseBoolFromString(user.RoleID),
		"sync_version":            helpers.StringValueOrNull(user.SyncVersion),
		"microtenant_id":          helpers.StringValueOrNull(user.MicrotenantID),
		"microtenant_name":        helpers.StringValueOrNull(user.MicrotenantName),
		"timezone":                helpers.StringValueOrNull(user.Timezone),
		"tmp_password":            helpers.StringValueOrNull(user.TmpPassword),
		"token_id":                helpers.StringValueOrNull(user.TokenID),
		"two_factor_auth_enabled": types.BoolValue(user.TwoFactorAuthEnabled),
		"two_factor_auth_type":    helpers.StringValueOrNull(user.TwoFactorAuthType),
		"username":                helpers.StringValueOrNull(user.Username),
		"creation_time":           helpers.StringValueOrNull(user.CreationTime),
		"modifiedby":              helpers.StringValueOrNull(user.ModifiedBy),
		"modified_time":           helpers.StringValueOrNull(user.ModifiedTime),
	})
	diags.Append(objDiags...)

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(setDiags...)
	return set, diags
}

func parseBoolFromString(value string) types.Bool {
	if strings.TrimSpace(value) == "" {
		return types.BoolNull()
	}
	// Try to parse as bool
	if value == "true" || value == "1" {
		return types.BoolValue(true)
	}
	return types.BoolValue(false)
}

func boolToString(value bool) types.String {
	if value {
		return types.StringValue("true")
	}
	return types.StringValue("false")
}
