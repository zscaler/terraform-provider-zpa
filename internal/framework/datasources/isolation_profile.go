package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/isolationprofile"
)

var (
	_ datasource.DataSource              = &IsolationProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &IsolationProfileDataSource{}
)

func NewIsolationProfileDataSource() datasource.DataSource {
	return &IsolationProfileDataSource{}
}

type IsolationProfileDataSource struct {
	client *client.Client
}

type IsolationProfileModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	CreationTime       types.String `tfsdk:"creation_time"`
	ModifiedTime       types.String `tfsdk:"modified_time"`
	ModifiedBy         types.String `tfsdk:"modified_by"`
	IsolationProfileID types.String `tfsdk:"isolation_profile_id"`
	IsolationTenantID  types.String `tfsdk:"isolation_tenant_id"`
	IsolationURL       types.String `tfsdk:"isolation_url"`
}

func (d *IsolationProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_isolation_profile"
}

func (d *IsolationProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a cloud browser isolation profile by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the isolation profile.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the isolation profile.",
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"isolation_profile_id": schema.StringAttribute{
				Computed: true,
			},
			"isolation_tenant_id": schema.StringAttribute{
				Computed: true,
			},
			"isolation_url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *IsolationProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IsolationProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data IsolationProfileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var (
		profile *isolationprofile.IsolationProfile
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving isolation profile by id", map[string]any{"id": id})
		// API does not have direct get by ID, so fetch all and locate
		profile, err = getIsolationProfileByID(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving isolation profile by name", map[string]any{"name": name})
		profile, _, err = isolationprofile.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read isolation profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Isolation profile with id %q or name %q not found.", id, name))
		return
	}

	data.ID = types.StringValue(profile.ID)
	data.Name = types.StringValue(profile.Name)
	data.Description = types.StringValue(profile.Description)
	data.Enabled = types.BoolValue(profile.Enabled)
	data.CreationTime = types.StringValue(profile.CreationTime)
	data.ModifiedBy = types.StringValue(profile.ModifiedBy)
	data.ModifiedTime = types.StringValue(profile.ModifiedTime)
	data.IsolationProfileID = types.StringValue(profile.IsolationProfileID)
	data.IsolationTenantID = types.StringValue(profile.IsolationTenantID)
	data.IsolationURL = types.StringValue(profile.IsolationURL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getIsolationProfileByID(ctx context.Context, service *zscaler.Service, id string) (*isolationprofile.IsolationProfile, error) {
	// The isolation profile SDK only supports listing and lookup by name, so retrieve all and search.
	profiles, _, err := isolationprofile.GetAll(ctx, service)
	if err != nil {
		return nil, err
	}

	for _, profile := range profiles {
		if profile.ID == id {
			return &profile, nil
		}
	}

	return nil, nil
}
