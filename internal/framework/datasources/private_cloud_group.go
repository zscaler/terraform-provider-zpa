package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

var (
	_ datasource.DataSource              = &PrivateCloudGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &PrivateCloudGroupDataSource{}
)

func NewPrivateCloudGroupDataSource() datasource.DataSource {
	return &PrivateCloudGroupDataSource{}
}

type PrivateCloudGroupDataSource struct {
	client *client.Client
}

type PrivateCloudGroupModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	CityCountry            types.String `tfsdk:"city_country"`
	CountryCode            types.String `tfsdk:"country_code"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	GeoLocationID          types.String `tfsdk:"geo_location_id"`
	IsPublic               types.String `tfsdk:"is_public"`
	Latitude               types.String `tfsdk:"latitude"`
	Location               types.String `tfsdk:"location"`
	Longitude              types.String `tfsdk:"longitude"`
	OverrideVersionProfile types.Bool   `tfsdk:"override_version_profile"`
	ReadOnly               types.Bool   `tfsdk:"read_only"`
	RestrictionType        types.String `tfsdk:"restriction_type"`
	MicroTenantID          types.String `tfsdk:"microtenant_id"`
	MicroTenantName        types.String `tfsdk:"microtenant_name"`
	SiteID                 types.String `tfsdk:"site_id"`
	SiteName               types.String `tfsdk:"site_name"`
	UpgradeDay             types.String `tfsdk:"upgrade_day"`
	UpgradeTimeInSecs      types.String `tfsdk:"upgrade_time_in_secs"`
	VersionProfileID       types.String `tfsdk:"version_profile_id"`
	VersionProfileName     types.String `tfsdk:"version_profile_name"`
	ZscalerManaged         types.Bool   `tfsdk:"zscaler_managed"`
	CreationTime           types.String `tfsdk:"creation_time"`
	ModifiedBy             types.String `tfsdk:"modified_by"`
	ModifiedTime           types.String `tfsdk:"modified_time"`
}

func (d *PrivateCloudGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_cloud_group"
}

func (d *PrivateCloudGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA Private Cloud Group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the private cloud group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the private cloud group.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"city_country": schema.StringAttribute{Computed: true},
			"country_code": schema.StringAttribute{Computed: true},
			"description":  schema.StringAttribute{Computed: true},
			"enabled":      schema.BoolAttribute{Computed: true},
			"geo_location_id": schema.StringAttribute{
				Computed: true,
			},
			"is_public": schema.StringAttribute{Computed: true},
			"latitude":  schema.StringAttribute{Computed: true},
			"location":  schema.StringAttribute{Computed: true},
			"longitude": schema.StringAttribute{Computed: true},
			"override_version_profile": schema.BoolAttribute{
				Computed: true,
			},
			"read_only":        schema.BoolAttribute{Computed: true},
			"restriction_type": schema.StringAttribute{Computed: true},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"site_id":          schema.StringAttribute{Computed: true},
			"site_name":        schema.StringAttribute{Computed: true},
			"upgrade_day":      schema.StringAttribute{Computed: true},
			"upgrade_time_in_secs": schema.StringAttribute{
				Computed: true,
			},
			"version_profile_id": schema.StringAttribute{Computed: true},
			"version_profile_name": schema.StringAttribute{
				Computed: true,
			},
			"zscaler_managed": schema.BoolAttribute{Computed: true},
			"creation_time":   schema.StringAttribute{Computed: true},
			"modified_by":     schema.StringAttribute{Computed: true},
			"modified_time":   schema.StringAttribute{Computed: true},
		},
	}
}

func (d *PrivateCloudGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PrivateCloudGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PrivateCloudGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		if microID := strings.TrimSpace(data.MicroTenantID.ValueString()); microID != "" {
			service = service.WithMicroTenant(microID)
			data.MicroTenantID = types.StringValue(microID)
		}
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var group *private_cloud_group.PrivateCloudGroup
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving private cloud group by ID", map[string]any{"id": id})
		group, _, err = private_cloud_group.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving private cloud group by name", map[string]any{"name": name})
		group, _, err = private_cloud_group.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read private cloud group: %v", err))
		return
	}

	if group == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Private cloud group with id %q or name %q was not found.", id, name))
		return
	}

	state := PrivateCloudGroupModel{
		ID:                     types.StringValue(group.ID),
		Name:                   types.StringValue(group.Name),
		CityCountry:            types.StringValue(group.CityCountry),
		CountryCode:            types.StringValue(group.CountryCode),
		Description:            types.StringValue(group.Description),
		Enabled:                types.BoolValue(group.Enabled),
		GeoLocationID:          types.StringValue(group.GeoLocationID),
		IsPublic:               types.StringValue(group.IsPublic),
		Latitude:               types.StringValue(group.Latitude),
		Location:               types.StringValue(group.Location),
		Longitude:              types.StringValue(group.Longitude),
		OverrideVersionProfile: types.BoolValue(group.OverrideVersionProfile),
		ReadOnly:               types.BoolValue(group.ReadOnly),
		RestrictionType:        types.StringValue(group.RestrictionType),
		MicroTenantName:        types.StringValue(group.MicrotenantName),
		SiteID:                 types.StringValue(group.SiteID),
		SiteName:               types.StringValue(group.SiteName),
		UpgradeDay:             types.StringValue(group.UpgradeDay),
		UpgradeTimeInSecs:      types.StringValue(group.UpgradeTimeInSecs),
		VersionProfileID:       types.StringValue(group.VersionProfileID),
		VersionProfileName:     types.StringValue(group.VersionProfileName),
		ZscalerManaged:         types.BoolValue(group.ZscalerManaged),
		CreationTime:           types.StringValue(group.CreationTime),
		ModifiedBy:             types.StringValue(group.ModifiedBy),
		ModifiedTime:           types.StringValue(group.ModifiedTime),
	}

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		state.MicroTenantID = data.MicroTenantID
	} else if group.MicrotenantID != "" {
		state.MicroTenantID = types.StringValue(group.MicrotenantID)
	} else {
		state.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
