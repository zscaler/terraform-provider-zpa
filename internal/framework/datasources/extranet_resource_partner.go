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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
)

var (
	_ datasource.DataSource              = &ExtranetResourcePartnerDataSource{}
	_ datasource.DataSourceWithConfigure = &ExtranetResourcePartnerDataSource{}
)

func NewExtranetResourcePartnerDataSource() datasource.DataSource {
	return &ExtranetResourcePartnerDataSource{}
}

type ExtranetResourcePartnerDataSource struct {
	client *client.Client
}

type ExtranetResourcePartnerModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func (d *ExtranetResourcePartnerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extranet_resource_partner"
}

func (d *ExtranetResourcePartnerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves an extranet resource partner summary by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the extranet resource partner.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the extranet resource partner.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *ExtranetResourcePartnerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExtranetResourcePartnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ExtranetResourcePartnerModel
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

	var partner *common.CommonSummary
	var err error

	if name != "" {
		tflog.Debug(ctx, "Retrieving extranet resource partner by name", map[string]any{"name": name})
		partner, _, err = extranet_resource.GetExtranetResourcePartnerByName(ctx, d.client.Service, name)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read extranet resource partner: %v", err))
			return
		}
	} else {
		tflog.Debug(ctx, "Searching extranet resource partners by ID", map[string]any{"id": id})
		partners, _, e := extranet_resource.GetExtranetResourcePartner(ctx, d.client.Service)
		err = e
		if err == nil {
			for i := range partners {
				if partners[i].ID == id {
					partner = &partners[i]
					break
				}
			}
		}
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list extranet resource partners: %v", err))
			return
		}
	}

	if partner == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Extranet resource partner with id %q or name %q not found.", id, name))
		return
	}

	data.ID = types.StringValue(partner.ID)
	data.Name = types.StringValue(partner.Name)
	data.Enabled = types.BoolValue(partner.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
