package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SegmentGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &SegmentGroupsDataSource{}
)

// NewSegmentGroupsDataSource is a helper function to simplify the provider implementation.
func NewSegmentGroupsDataSource() datasource.DataSource {
	return &SegmentGroupsDataSource{}
}

// SegmentGroupsDataSource defines the data source implementation.
type SegmentGroupsDataSource struct {
	client *client.Client
}

// SegmentGroupsDataSourceModel describes the data source data model.
type SegmentGroupsDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	ConfigSpace     types.String `tfsdk:"config_space"`
	CreationTime    types.String `tfsdk:"creation_time"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	MicrotenantID   types.String `tfsdk:"microtenant_id"`
	MicrotenantName types.String `tfsdk:"microtenant_name"`
}

func (d *SegmentGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment_group"
}

func (d *SegmentGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for ZPA Segment Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the segment group.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the segment group.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the segment group.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Status of the segment group.",
				Computed:    true,
			},
			"config_space": schema.StringAttribute{
				Description: "Configuration space.",
				Computed:    true,
			},
			"creation_time": schema.StringAttribute{
				Description: "Creation timestamp.",
				Computed:    true,
			},
			"modified_by": schema.StringAttribute{
				Description: "Last modifier.",
				Computed:    true,
			},
			"modified_time": schema.StringAttribute{
				Description: "Last modification timestamp.",
				Computed:    true,
			},
			"microtenant_id": schema.StringAttribute{
				Description: "Microtenant ID to scope the lookup.",
				Optional:    true,
				Computed:    true,
			},
			"microtenant_name": schema.StringAttribute{
				Description: "Microtenant name.",
				Computed:    true,
			},
		},
	}
}

func (d *SegmentGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *SegmentGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SegmentGroupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		segmentGroup *segmentgroup.SegmentGroup
		err          error
	)

	service := d.client.Service
	if !data.MicrotenantID.IsNull() && !data.MicrotenantID.IsUnknown() {
		service = service.WithMicroTenant(data.MicrotenantID.ValueString())
	}

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		segmentGroupID := data.ID.ValueString()
		tflog.Info(ctx, "Getting data for segment group ID", map[string]interface{}{
			"id": segmentGroupID,
		})
		segmentGroup, _, err = segmentgroup.Get(ctx, service, segmentGroupID)
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		name := data.Name.ValueString()
		tflog.Info(ctx, "Getting data for segment group name", map[string]interface{}{
			"name": name,
		})
		segmentGroup, _, err = segmentgroup.GetByName(ctx, service, name)
	} else {
		resp.Diagnostics.AddError("Missing Required Parameter", "Either 'id' or 'name' must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read segment group, got error: %s", err))
		return
	}

	if diags := flattenSegmentGroupsDataSource(ctx, segmentGroup, &data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Info(ctx, "Read segment group", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper functions

func flattenSegmentGroupsDataSource(ctx context.Context, segmentGroup *segmentgroup.SegmentGroup, data *SegmentGroupsDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(segmentGroup.ID)
	data.Name = types.StringValue(segmentGroup.Name)
	data.Description = types.StringValue(segmentGroup.Description)
	data.Enabled = types.BoolValue(segmentGroup.Enabled)
	data.ConfigSpace = types.StringValue(segmentGroup.ConfigSpace)
	data.CreationTime = types.StringValue(segmentGroup.CreationTime)
	data.ModifiedBy = types.StringValue(segmentGroup.ModifiedBy)
	data.ModifiedTime = types.StringValue(segmentGroup.ModifiedTime)
	data.MicrotenantID = types.StringValue(segmentGroup.MicroTenantID)
	data.MicrotenantName = types.StringValue(segmentGroup.MicroTenantName)

	return diags
}

func flattenSegmentGroupServerGroups(ctx context.Context, attrTypes map[string]attr.Type, groups []segmentgroup.AppServerGroup) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	elements := make([]attr.Value, 0, len(groups))

	for _, group := range groups {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                types.StringValue(group.ID),
			"name":              types.StringValue(group.Name),
			"description":       types.StringValue(group.Description),
			"config_space":      types.StringValue(group.ConfigSpace),
			"creation_time":     types.StringValue(group.CreationTime),
			"modified_by":       types.StringValue(group.ModifiedBy),
			"modified_time":     types.StringValue(group.ModifiedTime),
			"enabled":           types.BoolValue(group.Enabled),
			"dynamic_discovery": types.BoolValue(group.DynamicDiscovery),
		})
		diags.Append(objDiags...)
		elements = append(elements, obj)
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, elements)
	diags.Append(setDiags...)

	return set, diags
}

func stringSliceToSet(value interface{}) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if value == nil {
		return types.SetNull(types.StringType), diags
	}

	toStrings := func(items []interface{}) []string {
		out := make([]string, 0, len(items))
		for _, item := range items {
			switch v := item.(type) {
			case string:
				out = append(out, v)
			case fmt.Stringer:
				out = append(out, v.String())
			case float64:
				out = append(out, fmt.Sprintf("%g", v))
			case int:
				out = append(out, fmt.Sprintf("%d", v))
			case int32:
				out = append(out, fmt.Sprintf("%d", v))
			case int64:
				out = append(out, fmt.Sprintf("%d", v))
			case nil:
				// ignore
			default:
				out = append(out, fmt.Sprint(v))
			}
		}
		return out
	}

	var values []string
	switch v := value.(type) {
	case []string:
		values = v
	case []interface{}:
		values = toStrings(v)
	case string:
		if v != "" {
			values = []string{v}
		}
	default:
		values = []string{fmt.Sprint(v)}
	}

	if len(values) == 0 {
		return types.SetNull(types.StringType), diags
	}

	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}

	set, setDiags := types.SetValue(types.StringType, elements)
	return set, setDiags
}
