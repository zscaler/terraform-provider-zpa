package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
)

var (
	_ datasource.DataSource              = &CustomerVersionProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &CustomerVersionProfileDataSource{}
)

func NewCustomerVersionProfileDataSource() datasource.DataSource {
	return &CustomerVersionProfileDataSource{}
}

type CustomerVersionProfileDataSource struct {
	client *client.Client
}

type CustomerVersionProfileModel struct {
	ID                             types.String `tfsdk:"id"`
	CustomerID                     types.String `tfsdk:"customer_id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	CreationTime                   types.String `tfsdk:"creation_time"`
	ModifiedBy                     types.String `tfsdk:"modified_by"`
	ModifiedTime                   types.String `tfsdk:"modified_time"`
	NumberOfAssistants             types.String `tfsdk:"number_of_assistants"`
	NumberOfCustomers              types.String `tfsdk:"number_of_customers"`
	NumberOfPrivateBrokers         types.String `tfsdk:"number_of_private_brokers"`
	NumberOfSiteControllers        types.String `tfsdk:"number_of_site_controllers"`
	NumberOfUpdatedAssistants      types.String `tfsdk:"number_of_updated_assistants"`
	NumberOfUpdatedPrivateBrokers  types.String `tfsdk:"number_of_updated_private_brokers"`
	NumberOfUpdatedSiteControllers types.String `tfsdk:"number_of_updated_site_controllers"`
	UpgradePriority                types.String `tfsdk:"upgrade_priority"`
	VisibilityScope                types.String `tfsdk:"visibility_scope"`
	CustomScopeCustomerIDs         types.List   `tfsdk:"custom_scope_customer_ids"`
	CustomScopeRequestCustomerIDs  types.List   `tfsdk:"custom_scope_request_customer_ids"`
	Versions                       types.List   `tfsdk:"versions"`
}

func (d *CustomerVersionProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer_version_profile"
}

func (d *CustomerVersionProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a customer version profile by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"customer_id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"number_of_assistants": schema.StringAttribute{
				Computed: true,
			},
			"number_of_customers": schema.StringAttribute{
				Computed: true,
			},
			"number_of_private_brokers": schema.StringAttribute{
				Computed: true,
			},
			"number_of_site_controllers": schema.StringAttribute{
				Computed: true,
			},
			"number_of_updated_assistants": schema.StringAttribute{
				Computed: true,
			},
			"number_of_updated_private_brokers": schema.StringAttribute{
				Computed: true,
			},
			"number_of_updated_site_controllers": schema.StringAttribute{
				Computed: true,
			},
			"upgrade_priority": schema.StringAttribute{
				Computed: true,
			},
			"visibility_scope": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"custom_scope_customer_ids": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"customer_id":           schema.StringAttribute{Computed: true},
						"name":                  schema.StringAttribute{Computed: true},
						"exclude_constellation": schema.BoolAttribute{Computed: true},
						"is_partner":            schema.BoolAttribute{Computed: true},
					},
				},
			},
			"custom_scope_request_customer_ids": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"add_customer_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"delete_customer_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"versions": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"creation_time":                schema.StringAttribute{Computed: true},
						"customer_id":                  schema.StringAttribute{Computed: true},
						"id":                           schema.StringAttribute{Computed: true},
						"modified_by":                  schema.StringAttribute{Computed: true},
						"modified_time":                schema.StringAttribute{Computed: true},
						"platform":                     schema.StringAttribute{Computed: true},
						"restart_after_uptime_in_days": schema.StringAttribute{Computed: true},
						"role":                         schema.StringAttribute{Computed: true},
						"version":                      schema.StringAttribute{Computed: true},
						"version_profile_gid":          schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *CustomerVersionProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomerVersionProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CustomerVersionProfileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(data.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddError("Missing Name", "The 'name' attribute is required.")
		return
	}

	tflog.Debug(ctx, "Retrieving customer version profile", map[string]any{"name": name})

	profile, _, err := customerversionprofile.GetByName(ctx, d.client.Service, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read customer version profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Customer version profile with name %q not found.", name))
		return
	}

	state, diags := flattenCustomerVersionProfile(ctx, profile)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenCustomerVersionProfile(ctx context.Context, profile *customerversionprofile.CustomerVersionProfile) (CustomerVersionProfileModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	customScopeCustomers, diags1 := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"customer_id":           types.StringType,
		"name":                  types.StringType,
		"exclude_constellation": types.BoolType,
		"is_partner":            types.BoolType,
	}}, convertScopeCustomerIDs(profile.CustomScopeCustomerIDs))
	diags.Append(diags1...)

	customScopeRequests, diags2 := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"add_customer_ids":    types.ListType{ElemType: types.StringType},
		"delete_customer_ids": types.ListType{ElemType: types.StringType},
	}}, convertScopeRequestCustomers(profile.CustomScopeRequestCustomerIDs))
	diags.Append(diags2...)

	versions, diags3 := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"creation_time":                types.StringType,
		"customer_id":                  types.StringType,
		"id":                           types.StringType,
		"modified_by":                  types.StringType,
		"modified_time":                types.StringType,
		"platform":                     types.StringType,
		"restart_after_uptime_in_days": types.StringType,
		"role":                         types.StringType,
		"version":                      types.StringType,
		"version_profile_gid":          types.StringType,
	}}, convertVersions(profile.Versions))
	diags.Append(diags3...)

	model := CustomerVersionProfileModel{
		ID:                             types.StringValue(profile.ID),
		CustomerID:                     types.StringValue(profile.CustomerID),
		Name:                           types.StringValue(profile.Name),
		Description:                    types.StringValue(profile.Description),
		CreationTime:                   types.StringValue(profile.CreationTime),
		ModifiedBy:                     types.StringValue(profile.ModifiedBy),
		ModifiedTime:                   types.StringValue(profile.ModifiedTime),
		NumberOfAssistants:             types.StringValue(profile.NumberOfAssistants),
		NumberOfCustomers:              types.StringValue(profile.NumberOfCustomers),
		NumberOfPrivateBrokers:         types.StringValue(profile.NumberOfPrivateBrokers),
		NumberOfSiteControllers:        types.StringValue(profile.NumberOfSiteControllers),
		NumberOfUpdatedAssistants:      types.StringValue(profile.NumberOfUpdatedAssistants),
		NumberOfUpdatedPrivateBrokers:  types.StringValue(profile.NumberOfUpdatedPrivateBrokers),
		NumberOfUpdatedSiteControllers: types.StringValue(profile.NumberOfUpdatedSiteControllers),
		UpgradePriority:                types.StringValue(profile.UpgradePriority),
		VisibilityScope:                types.StringValue(profile.VisibilityScope),
		CustomScopeCustomerIDs:         customScopeCustomers,
		CustomScopeRequestCustomerIDs:  customScopeRequests,
		Versions:                       versions,
	}

	return model, diags
}

func convertScopeCustomerIDs(scope []customerversionprofile.CustomScopeCustomerIDs) []map[string]attr.Value {
	result := make([]map[string]attr.Value, 0, len(scope))
	for _, item := range scope {
		result = append(result, map[string]attr.Value{
			"customer_id":           types.StringValue(item.CustomerID),
			"name":                  types.StringValue(item.Name),
			"exclude_constellation": types.BoolValue(item.ExcludeConstellation),
			"is_partner":            types.BoolValue(item.IsPartner),
		})
	}
	return result
}

func convertScopeRequestCustomers(req customerversionprofile.CustomScopeRequestCustomerIDs) []map[string]attr.Value {
	if req.AddCustomerIDs == "" && req.DeletecustomerIDs == "" {
		return nil
	}

	return []map[string]attr.Value{
		{
			"add_customer_ids":    listFromCSV(req.AddCustomerIDs),
			"delete_customer_ids": listFromCSV(req.DeletecustomerIDs),
		},
	}
}

func convertVersions(versions []customerversionprofile.Versions) []map[string]attr.Value {
	result := make([]map[string]attr.Value, 0, len(versions))
	for _, v := range versions {
		result = append(result, map[string]attr.Value{
			"creation_time":                types.StringValue(v.CreationTime),
			"customer_id":                  types.StringValue(v.CustomerID),
			"id":                           types.StringValue(v.ID),
			"modified_by":                  types.StringValue(v.ModifiedBy),
			"modified_time":                types.StringValue(v.ModifiedTime),
			"platform":                     types.StringValue(v.Platform),
			"restart_after_uptime_in_days": types.StringValue(v.RestartAfterUptimeInDays),
			"role":                         types.StringValue(v.Role),
			"version":                      types.StringValue(v.Version),
			"version_profile_gid":          types.StringValue(v.VersionProfileGID),
		})
	}
	return result
}

func listFromCSV(csv string) types.List {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return types.ListNull(types.StringType)
	}
	parts := strings.Split(csv, ",")
	values := make([]attr.Value, 0, len(parts))
	for _, part := range parts {
		values = append(values, types.StringValue(strings.TrimSpace(part)))
	}
	list, _ := types.ListValue(types.StringType, values)
	return list
}
