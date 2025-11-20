package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/postureprofile"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &PostureProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &PostureProfileDataSource{}
)

func NewPostureProfileDataSource() datasource.DataSource {
	return &PostureProfileDataSource{}
}

type PostureProfileDataSource struct {
	client *client.Client
}

type PostureProfileDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	MasterCustomerID  types.String `tfsdk:"master_customer_id"`
	Domain            types.String `tfsdk:"domain"`
	PostureUDID       types.String `tfsdk:"posture_udid"`
	ZscalerCloud      types.String `tfsdk:"zscaler_cloud"`
	ZscalerCustomerID types.String `tfsdk:"zscaler_customer_id"`
	CreationTime      types.String `tfsdk:"creation_time"`
	ModifiedTime      types.String `tfsdk:"modified_time"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
}

func (d *PostureProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_posture_profile"
}

func (d *PostureProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA posture profile by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"master_customer_id":  schema.StringAttribute{Computed: true},
			"domain":              schema.StringAttribute{Computed: true},
			"posture_udid":        schema.StringAttribute{Computed: true},
			"zscaler_cloud":       schema.StringAttribute{Computed: true},
			"zscaler_customer_id": schema.StringAttribute{Computed: true},
			"creation_time":       schema.StringAttribute{Computed: true},
			"modified_time":       schema.StringAttribute{Computed: true},
			"modified_by":         schema.StringAttribute{Computed: true},
		},
	}
}

func (d *PostureProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *PostureProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PostureProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.ID.IsNull() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be specified")
		return
	}

	service := d.client.Service

	var (
		profile *postureprofile.PostureProfile
		err     error
	)

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		id := data.ID.ValueString()
		tflog.Info(ctx, "Fetching posture profile", map[string]any{"id": id})
		profile, _, err = postureprofile.Get(ctx, service, id)
	} else {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Fetching posture profile", map[string]any{"name": name})
		profile, _, err = postureprofile.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read posture profile: %v", err))
		return
	}

	model := PostureProfileDataSourceModel{
		ID:                types.StringValue(profile.ID),
		Name:              types.StringValue(profile.Name),
		MasterCustomerID:  types.StringValue(profile.MasterCustomerID),
		Domain:            types.StringValue(profile.Domain),
		PostureUDID:       types.StringValue(profile.PostureudID),
		ZscalerCloud:      types.StringValue(profile.ZscalerCloud),
		ZscalerCustomerID: types.StringValue(profile.ZscalerCustomerID),
		CreationTime:      types.StringValue(profile.CreationTime),
		ModifiedTime:      types.StringValue(profile.ModifiedTime),
		ModifiedBy:        types.StringValue(profile.ModifiedBy),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
