package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/workload_tag_group"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ datasource.DataSource              = &WorkloadTagGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkloadTagGroupDataSource{}
)

func NewWorkloadTagGroupDataSource() datasource.DataSource {
	return &WorkloadTagGroupDataSource{}
}

type WorkloadTagGroupDataSource struct {
	client *client.Client
}

type WorkloadTagGroupDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	MicroTenantID types.String `tfsdk:"microtenant_id"`
}

func (d *WorkloadTagGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workload_tag_group"
}

func (d *WorkloadTagGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA workload tag group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Optional: true},
			"name":           schema.StringAttribute{Optional: true},
			"enabled":        schema.BoolAttribute{Computed: true},
			"microtenant_id": schema.StringAttribute{Optional: true},
		},
	}
}

func (d *WorkloadTagGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkloadTagGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkloadTagGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (data.ID.IsNull() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Identifier", "Either id or name must be specified")
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() && data.MicroTenantID.ValueString() != "" {
		service = service.WithMicroTenant(data.MicroTenantID.ValueString())
	}
	var result *common.CommonSummary

	id := ""
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id = strings.TrimSpace(data.ID.ValueString())
	}

	name := ""
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		name = strings.TrimSpace(data.Name.ValueString())
	}

	tflog.Info(ctx, "Fetching workload tag groups", map[string]any{"id": id, "name": name})
	groups, _, err := workload_tag_group.GetWorkloadTagGroup(ctx, service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workload tag group: %v", err))
		return
	}

	if id != "" {
		for i := range groups {
			if groups[i].ID == id {
				result = &groups[i]
				break
			}
		}
	}

	if result == nil && name != "" {
		for i := range groups {
			if strings.EqualFold(groups[i].Name, name) {
				result = &groups[i]
				break
			}
		}
	}

	if result == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Couldn't find any workload tag group with name '%s' or id '%s'", name, id))
		return
	}

	model := WorkloadTagGroupDataSourceModel{
		ID:            types.StringValue(result.ID),
		Name:          types.StringValue(result.Name),
		Enabled:       types.BoolValue(result.Enabled),
		MicroTenantID: data.MicroTenantID,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
