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
	inspectionpredefined "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

var (
	_ datasource.DataSource              = &InspectionPredefinedControlsDataSource{}
	_ datasource.DataSourceWithConfigure = &InspectionPredefinedControlsDataSource{}
)

func NewInspectionPredefinedControlsDataSource() datasource.DataSource {
	return &InspectionPredefinedControlsDataSource{}
}

type InspectionPredefinedControlsDataSource struct {
	client *client.Client
}

type InspectionPredefinedControlsModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Version                          types.String `tfsdk:"version"`
	Action                           types.String `tfsdk:"action"`
	ActionValue                      types.String `tfsdk:"action_value"`
	Attachment                       types.String `tfsdk:"attachment"`
	ControlGroup                     types.String `tfsdk:"control_group"`
	ControlNumber                    types.String `tfsdk:"control_number"`
	ControlType                      types.String `tfsdk:"control_type"`
	CreationTime                     types.String `tfsdk:"creation_time"`
	DefaultAction                    types.String `tfsdk:"default_action"`
	DefaultActionValue               types.String `tfsdk:"default_action_value"`
	Description                      types.String `tfsdk:"description"`
	ModifiedBy                       types.String `tfsdk:"modifiedby"`
	ModifiedTime                     types.String `tfsdk:"modified_time"`
	ParanoiaLevel                    types.String `tfsdk:"paranoia_level"`
	ProtocolType                     types.String `tfsdk:"protocol_type"`
	Severity                         types.String `tfsdk:"severity"`
	AssociatedInspectionProfileNames types.Set    `tfsdk:"associated_inspection_profile_names"`
}

func (d *InspectionPredefinedControlsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inspection_predefined_controls"
}

func (d *InspectionPredefinedControlsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a predefined inspection control by ID or name/version.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the predefined control.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the predefined control.",
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Version of the predefined control. Required when 'name' is specified.",
			},
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
	}
}

func (d *InspectionPredefinedControlsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InspectionPredefinedControlsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data InspectionPredefinedControlsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	version := strings.TrimSpace(data.Version.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a predefined control.")
		return
	}

	if id == "" && name != "" && version == "" {
		resp.Diagnostics.AddError("Missing Version", "When 'name' is specified, 'version' must also be provided.")
		return
	}

	var (
		control *inspectionpredefined.PredefinedControls
		err     error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving predefined control by ID", map[string]any{"id": id})
		control, _, err = inspectionpredefined.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving predefined control by name", map[string]any{"name": name, "version": version})
		control, _, err = inspectionpredefined.GetByName(ctx, d.client.Service, name, version)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read predefined inspection control: %v", err))
		return
	}

	if control == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Predefined inspection control with id %q or name %q was not found.", id, name))
		return
	}

	associated, assocDiags := flattenAssociatedProfileNames(ctx, control.AssociatedInspectionProfileNames)
	resp.Diagnostics.Append(assocDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(control.ID)
	data.Name = stringOrNull(control.Name)
	if version != "" {
		data.Version = types.StringValue(version)
	} else {
		data.Version = stringOrNull(control.Version)
	}
	data.Action = stringOrNull(control.Action)
	data.ActionValue = stringOrNull(control.ActionValue)
	data.Attachment = stringOrNull(control.Attachment)
	data.ControlGroup = stringOrNull(control.ControlGroup)
	data.ControlNumber = stringOrNull(control.ControlNumber)
	data.ControlType = stringOrNull(control.ControlType)
	data.CreationTime = stringOrNull(control.CreationTime)
	data.DefaultAction = stringOrNull(control.DefaultAction)
	data.DefaultActionValue = stringOrNull(control.DefaultActionValue)
	data.Description = stringOrNull(control.Description)
	data.ModifiedBy = stringOrNull(control.ModifiedBy)
	data.ModifiedTime = stringOrNull(control.ModifiedTime)
	data.ParanoiaLevel = stringOrNull(control.ParanoiaLevel)
	data.ProtocolType = stringOrNull(control.ProtocolType)
	data.Severity = stringOrNull(control.Severity)
	data.AssociatedInspectionProfileNames = associated

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
