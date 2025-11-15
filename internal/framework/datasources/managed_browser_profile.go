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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/managed_browser"
)

var (
	_ datasource.DataSource              = &ManagedBrowserProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &ManagedBrowserProfileDataSource{}
)

func NewManagedBrowserProfileDataSource() datasource.DataSource {
	return &ManagedBrowserProfileDataSource{}
}

type ManagedBrowserProfileDataSource struct {
	client *client.Client
}

type ManagedBrowserProfileModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	BrowserType          types.String `tfsdk:"browser_type"`
	CustomerID           types.String `tfsdk:"customer_id"`
	MicroTenantID        types.String `tfsdk:"microtenant_id"`
	MicroTenantName      types.String `tfsdk:"microtenant_name"`
	ChromePostureProfile types.List   `tfsdk:"chrome_posture_profile"`
	CreationTime         types.String `tfsdk:"creation_time"`
	ModifiedBy           types.String `tfsdk:"modified_by"`
	ModifiedTime         types.String `tfsdk:"modified_time"`
}

func (d *ManagedBrowserProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_browser_profile"
}

func (d *ManagedBrowserProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA managed browser profile by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the managed browser profile.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the managed browser profile.",
			},
			"description":      schema.StringAttribute{Computed: true},
			"browser_type":     schema.StringAttribute{Computed: true},
			"customer_id":      schema.StringAttribute{Computed: true},
			"microtenant_id":   schema.StringAttribute{Computed: true},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"creation_time":    schema.StringAttribute{Computed: true},
			"modified_by":      schema.StringAttribute{Computed: true},
			"modified_time":    schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"chrome_posture_profile": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                 schema.StringAttribute{Computed: true},
						"browser_type":       schema.StringAttribute{Computed: true},
						"crowd_strike_agent": schema.BoolAttribute{Computed: true},
						"creation_time":      schema.StringAttribute{Computed: true},
						"modified_by":        schema.StringAttribute{Computed: true},
						"modified_time":      schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ManagedBrowserProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ManagedBrowserProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ManagedBrowserProfileModel
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

	var profile *managed_browser.ManagedBrowserProfile
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving managed browser profile by ID", map[string]any{"id": id})
		profiles, _, ferr := managed_browser.GetAll(ctx, service)
		if ferr != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list managed browser profiles: %v", ferr))
			return
		}
		for _, candidate := range profiles {
			if candidate.ID == id {
				profile = &candidate
				break
			}
		}
	} else {
		tflog.Debug(ctx, "Retrieving managed browser profile by name", map[string]any{"name": name})
		profile, _, err = managed_browser.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read managed browser profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Managed browser profile with id %q or name %q was not found.", id, name))
		return
	}

	chromeProfile, chromeDiags := flattenChromePostureProfile(ctx, profile.ChromePostureProfile)
	resp.Diagnostics.Append(chromeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(profile.ID)
	data.Name = types.StringValue(profile.Name)
	data.Description = types.StringValue(profile.Description)
	data.BrowserType = types.StringValue(profile.BrowserType)
	data.CustomerID = types.StringValue(profile.CustomerID)
	data.MicroTenantID = types.StringValue(profile.MicrotenantID)
	data.MicroTenantName = types.StringValue(profile.MicrotenantName)
	data.ChromePostureProfile = chromeProfile
	data.CreationTime = types.StringValue(profile.CreationTime)
	data.ModifiedBy = types.StringValue(profile.ModifiedBy)
	data.ModifiedTime = types.StringValue(profile.ModifiedTime)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenChromePostureProfile(ctx context.Context, profile managed_browser.ChromePostureProfile) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":                 types.StringType,
		"browser_type":       types.StringType,
		"crowd_strike_agent": types.BoolType,
		"creation_time":      types.StringType,
		"modified_by":        types.StringType,
		"modified_time":      types.StringType,
	}

	if profile.ID == "" {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":                 types.StringValue(profile.ID),
		"browser_type":       types.StringValue(profile.BrowserType),
		"crowd_strike_agent": types.BoolValue(profile.CrowdStrikeAgent),
		"creation_time":      types.StringValue(profile.CreationTime),
		"modified_by":        types.StringValue(profile.ModifiedBy),
		"modified_time":      types.StringValue(profile.ModifiedTime),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}
