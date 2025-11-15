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
	cbiregions "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiregions"
)

var (
	_ datasource.DataSource              = &CBIRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &CBIRegionsDataSource{}
)

func NewCBIRegionsDataSource() datasource.DataSource {
	return &CBIRegionsDataSource{}
}

type CBIRegionsDataSource struct {
	client *client.Client
}

type CBIRegionsModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *CBIRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_browser_isolation_region"
}

func (d *CBIRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Cloud Browser Isolation region by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the CBI region.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the CBI region.",
			},
		},
	}
}

func (d *CBIRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CBIRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CBIRegionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a CBI region.")
		return
	}

	var (
		region *cbiregions.CBIRegions
		err    error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving CBI region by ID", map[string]any{"id": id})
		regions, _, fetchErr := cbiregions.GetAll(ctx, d.client.Service)
		if fetchErr != nil {
			err = fetchErr
		} else {
			for _, candidate := range regions {
				candidate := candidate
				if candidate.ID == id {
					region = &candidate
					break
				}
			}
			if region == nil {
				err = fmt.Errorf("no CBI region with id %s was found", id)
			}
		}
	} else {
		tflog.Debug(ctx, "Retrieving CBI region by name", map[string]any{"name": name})
		region, _, err = cbiregions.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read CBI region: %v", err))
		return
	}

	if region == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("CBI region with id %q or name %q was not found.", id, name))
		return
	}

	data.ID = types.StringValue(region.ID)
	data.Name = stringOrNull(region.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
