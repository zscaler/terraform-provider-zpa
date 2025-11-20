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
	cbizpaprofile "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbizpaprofile"
)

var (
	_ datasource.DataSource              = &CBIZPAProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &CBIZPAProfilesDataSource{}
)

func NewCBIZPAProfilesDataSource() datasource.DataSource {
	return &CBIZPAProfilesDataSource{}
}

type CBIZPAProfilesDataSource struct {
	client *client.Client
}

type CBIZPAProfilesModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	CreationTime types.String `tfsdk:"creation_time"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	ModifiedTime types.String `tfsdk:"modified_time"`
	CBITenantID  types.String `tfsdk:"cbi_tenant_id"`
	CBIProfileID types.String `tfsdk:"cbi_profile_id"`
	CBIURL       types.String `tfsdk:"cbi_url"`
}

func (d *CBIZPAProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_zpa_profiles"
}

func (d *CBIZPAProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Cloud Browser Isolation ZPA profile by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the CBI ZPA profile.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the CBI ZPA profile.",
			},
			"description":    schema.StringAttribute{Computed: true},
			"enabled":        schema.BoolAttribute{Computed: true},
			"creation_time":  schema.StringAttribute{Computed: true},
			"modified_by":    schema.StringAttribute{Computed: true},
			"modified_time":  schema.StringAttribute{Computed: true},
			"cbi_tenant_id":  schema.StringAttribute{Computed: true},
			"cbi_profile_id": schema.StringAttribute{Computed: true},
			"cbi_url":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *CBIZPAProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CBIZPAProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CBIZPAProfilesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a CBI ZPA profile.")
		return
	}

	var (
		profile *cbizpaprofile.ZPAProfiles
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving CBI ZPA profile by ID", map[string]any{"id": id})
		profile, _, err = cbizpaprofile.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving CBI ZPA profile by name", map[string]any{"name": name})
		profile, _, err = cbizpaprofile.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read CBI ZPA profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("CBI ZPA profile with id %q or name %q was not found.", id, name))
		return
	}

	data.ID = types.StringValue(profile.ID)
	data.Name = stringOrNull(profile.Name)
	data.Description = stringOrNull(profile.Description)
	data.Enabled = types.BoolValue(profile.Enabled)
	data.CreationTime = stringOrNull(profile.CreationTime)
	data.ModifiedBy = stringOrNull(profile.ModifiedBy)
	data.ModifiedTime = stringOrNull(profile.ModifiedTime)
	data.CBITenantID = stringOrNull(profile.CBITenantID)
	data.CBIProfileID = stringOrNull(profile.CBIProfileID)
	data.CBIURL = stringOrNull(profile.CBIURL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
