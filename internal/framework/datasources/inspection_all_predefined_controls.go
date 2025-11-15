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
	inspectionpredefined "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

var (
	_ datasource.DataSource              = &InspectionAllPredefinedControlsDataSource{}
	_ datasource.DataSourceWithConfigure = &InspectionAllPredefinedControlsDataSource{}
)

func NewInspectionAllPredefinedControlsDataSource() datasource.DataSource {
	return &InspectionAllPredefinedControlsDataSource{}
}

type InspectionAllPredefinedControlsDataSource struct {
	client *client.Client
}

type InspectionAllPredefinedControlsModel struct {
	Version   types.String `tfsdk:"version"`
	GroupName types.String `tfsdk:"group_name"`
	List      types.List   `tfsdk:"list"`
}

func (d *InspectionAllPredefinedControlsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inspection_all_predefined_controls"
}

func (d *InspectionAllPredefinedControlsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves all predefined inspection controls for a given version and optional group name.",
		Attributes: map[string]schema.Attribute{
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "Version of the predefined controls. Defaults to 'OWASP_CRS/3.3.0'.",
			},
			"group_name": schema.StringAttribute{
				Optional:    true,
				Description: "Optional control group name to filter the predefined controls.",
			},
		},
		Blocks: map[string]schema.Block{
			"list": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"action":               schema.StringAttribute{Computed: true},
						"action_value":         schema.StringAttribute{Computed: true},
						"attachment":           schema.StringAttribute{Computed: true},
						"control_group":        schema.StringAttribute{Computed: true},
						"control_number":       schema.StringAttribute{Computed: true},
						"control_type":         schema.StringAttribute{Computed: true},
						"creation_time":        schema.StringAttribute{Computed: true},
						"default_action":       schema.StringAttribute{Computed: true},
						"default_action_value": schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"modifiedby":           schema.StringAttribute{Computed: true},
						"modified_time":        schema.StringAttribute{Computed: true},
						"paranoia_level":       schema.StringAttribute{Computed: true},
						"protocol_type":        schema.StringAttribute{Computed: true},
						"severity":             schema.StringAttribute{Computed: true},
						"version":              schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"associated_inspection_profile_names": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Computed: true},
									"name": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *InspectionAllPredefinedControlsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InspectionAllPredefinedControlsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data InspectionAllPredefinedControlsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version := strings.TrimSpace(data.Version.ValueString())
	if version == "" {
		version = "OWASP_CRS/3.3.0"
	}
	groupName := strings.TrimSpace(data.GroupName.ValueString())

	tflog.Debug(ctx, "Retrieving predefined controls", map[string]any{"version": version, "group": groupName})

	var (
		controls []inspectionpredefined.PredefinedControls
		err      error
	)

	if groupName != "" {
		controls, err = inspectionpredefined.GetAllByGroup(ctx, d.client.Service, version, groupName)
	} else {
		controls, err = inspectionpredefined.GetAll(ctx, d.client.Service, version)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve predefined inspection controls: %v", err))
		return
	}

	list, listDiags := flattenPredefinedControlList(ctx, controls)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Version = types.StringValue(version)
	if groupName != "" {
		data.GroupName = types.StringValue(groupName)
	} else {
		data.GroupName = types.StringNull()
	}
	data.List = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenPredefinedControlList(ctx context.Context, controls []inspectionpredefined.PredefinedControls) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":                                  types.StringType,
		"name":                                types.StringType,
		"action":                              types.StringType,
		"action_value":                        types.StringType,
		"attachment":                          types.StringType,
		"control_group":                       types.StringType,
		"control_number":                      types.StringType,
		"control_type":                        types.StringType,
		"creation_time":                       types.StringType,
		"default_action":                      types.StringType,
		"default_action_value":                types.StringType,
		"description":                         types.StringType,
		"modifiedby":                          types.StringType,
		"modified_time":                       types.StringType,
		"paranoia_level":                      types.StringType,
		"protocol_type":                       types.StringType,
		"severity":                            types.StringType,
		"version":                             types.StringType,
		"associated_inspection_profile_names": types.ListType{ElemType: types.ObjectType{AttrTypes: associatedProfileAttrTypes()}},
	}

	if len(controls) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(controls))
	var diags diag.Diagnostics
	for _, control := range controls {
		// Convert inspectionpredefined.AssociatedInspectionProfileNames to common.AssociatedProfileNames
		commonNames := make([]common.AssociatedProfileNames, 0, len(control.AssociatedInspectionProfileNames))
		for _, name := range control.AssociatedInspectionProfileNames {
			commonNames = append(commonNames, common.AssociatedProfileNames{
				ID:   name.ID,
				Name: name.Name,
			})
		}
		associated, assocDiags := flattenAssociatedProfileNamesAsList(ctx, commonNames)
		diags.Append(assocDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                                  stringOrNull(control.ID),
			"name":                                stringOrNull(control.Name),
			"action":                              stringOrNull(control.Action),
			"action_value":                        stringOrNull(control.ActionValue),
			"attachment":                          stringOrNull(control.Attachment),
			"control_group":                       stringOrNull(control.ControlGroup),
			"control_number":                      stringOrNull(control.ControlNumber),
			"control_type":                        stringOrNull(control.ControlType),
			"creation_time":                       stringOrNull(control.CreationTime),
			"default_action":                      stringOrNull(control.DefaultAction),
			"default_action_value":                stringOrNull(control.DefaultActionValue),
			"description":                         stringOrNull(control.Description),
			"modifiedby":                          stringOrNull(control.ModifiedBy),
			"modified_time":                       stringOrNull(control.ModifiedTime),
			"paranoia_level":                      stringOrNull(control.ParanoiaLevel),
			"protocol_type":                       stringOrNull(control.ProtocolType),
			"severity":                            stringOrNull(control.Severity),
			"version":                             stringOrNull(control.Version),
			"associated_inspection_profile_names": associated,
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
