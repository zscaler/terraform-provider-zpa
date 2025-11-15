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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
)

var (
	_ datasource.DataSource              = &LocationGroupControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationGroupControllerDataSource{}
)

func NewLocationGroupControllerDataSource() datasource.DataSource {
	return &LocationGroupControllerDataSource{}
}

type LocationGroupControllerDataSource struct {
	client *client.Client
}

type LocationGroupControllerModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	LocationName      types.String `tfsdk:"location_name"`
	ZIAErName         types.String `tfsdk:"zia_er_name"`
	LocationGroupID   types.String `tfsdk:"location_group_id"`
	LocationGroupName types.String `tfsdk:"location_group_name"`
	ZiaLocations      types.List   `tfsdk:"zia_locations"`
}

func (d *LocationGroupControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location_group_controller"
}

func (d *LocationGroupControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a location group associated with an extranet resource partner, including ZIA locations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the location group.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the location group.",
			},
			"location_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the ZIA location to search for.",
			},
			"zia_er_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the extranet resource partner.",
			},
			"location_group_id": schema.StringAttribute{
				Computed: true,
			},
			"location_group_name": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"zia_locations": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true},
						"name":    schema.StringAttribute{Computed: true},
						"enabled": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *LocationGroupControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationGroupControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data LocationGroupControllerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationName := strings.TrimSpace(data.LocationName.ValueString())
	extranetName := strings.TrimSpace(data.ZIAErName.ValueString())

	if locationName == "" || extranetName == "" {
		resp.Diagnostics.AddError("Missing Required Attributes", "'location_name' and 'zia_er_name' must be provided.")
		return
	}

	tflog.Debug(ctx, "Resolving extranet resource partner", map[string]any{"name": extranetName})
	extranet, _, err := extranet_resource.GetExtranetResourcePartnerByName(ctx, d.client.Service, extranetName)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read extranet resource partner %q: %v", extranetName, err))
		return
	}
	if extranet == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Extranet resource partner %q not found.", extranetName))
		return
	}

	tflog.Debug(ctx, "Retrieving location groups for extranet resource", map[string]any{"er_id": extranet.ID})
	groups, _, err := location_controller.GetLocationGroupExtranetResource(ctx, d.client.Service, extranet.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read location groups for extranet resource %q: %v", extranetName, err))
		return
	}

	var targetGroup *common.LocationGroupDTO
	var targetLocation *common.CommonSummary
	for i := range groups {
		for j := range groups[i].ZiaLocations {
			if strings.EqualFold(groups[i].ZiaLocations[j].Name, locationName) {
				targetGroup = &groups[i]
				targetLocation = &groups[i].ZiaLocations[j]
				break
			}
		}
		if targetGroup != nil {
			break
		}
	}

	if targetGroup == nil || targetLocation == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Location %q not found in groups for extranet resource %q.", locationName, extranetName))
		return
	}

	locationsList, diags := flattenZiaLocationsList(targetGroup.ZiaLocations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(targetGroup.ID)
	data.Name = types.StringValue(targetGroup.Name)
	data.LocationGroupID = types.StringValue(targetGroup.ID)
	data.LocationGroupName = types.StringValue(targetGroup.Name)
	data.ZiaLocations = locationsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenZiaLocationsList(locations []common.CommonSummary) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":      types.StringType,
		"name":    types.StringType,
		"enabled": types.BoolType,
	}

	values := make([]attr.Value, 0, len(locations))
	for _, loc := range locations {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":      types.StringValue(loc.ID),
			"name":    types.StringValue(loc.Name),
			"enabled": types.BoolValue(loc.Enabled),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
