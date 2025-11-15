package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cbibannercontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

var (
	_ datasource.DataSource              = &CBIBannersDataSource{}
	_ datasource.DataSourceWithConfigure = &CBIBannersDataSource{}
)

func NewCBIBannersDataSource() datasource.DataSource {
	return &CBIBannersDataSource{}
}

type CBIBannersDataSource struct {
	client *client.Client
}

type CBIBannersModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	PrimaryColor      types.String `tfsdk:"primary_color"`
	TextColor         types.String `tfsdk:"text_color"`
	NotificationTitle types.String `tfsdk:"notification_title"`
	NotificationText  types.String `tfsdk:"notification_text"`
	Logo              types.String `tfsdk:"logo"`
	Banner            types.Bool   `tfsdk:"banner"`
	IsDefault         types.Bool   `tfsdk:"is_default"`
}

func (d *CBIBannersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_banner"
}

func (d *CBIBannersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Cloud Browser Isolation banner by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the CBI banner.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the CBI banner.",
			},
			"primary_color":      schema.StringAttribute{Computed: true},
			"text_color":         schema.StringAttribute{Computed: true},
			"notification_title": schema.StringAttribute{Computed: true},
			"notification_text":  schema.StringAttribute{Computed: true},
			"logo":               schema.StringAttribute{Computed: true},
			"banner":             schema.BoolAttribute{Computed: true},
			"is_default":         schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *CBIBannersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CBIBannersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CBIBannersModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a CBI banner.")
		return
	}

	identifier := id
	if identifier == "" {
		identifier = name
	}

	banner, _, err := cbibannercontroller.GetByNameOrID(ctx, d.client.Service, identifier)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read CBI banner: %v", err))
		return
	}

	if banner == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("CBI banner with identifier %q was not found.", identifier))
		return
	}

	data.ID = types.StringValue(banner.ID)
	data.Name = stringOrNull(banner.Name)
	data.PrimaryColor = stringOrNull(banner.PrimaryColor)
	data.TextColor = stringOrNull(banner.TextColor)
	data.NotificationTitle = stringOrNull(banner.NotificationTitle)
	data.NotificationText = stringOrNull(banner.NotificationText)
	data.Logo = stringOrNull(banner.Logo)
	data.Banner = types.BoolValue(banner.Banner)
	data.IsDefault = types.BoolValue(banner.IsDefault)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
