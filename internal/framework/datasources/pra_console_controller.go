package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

var (
	_ datasource.DataSource              = &PRAConsoleControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &PRAConsoleControllerDataSource{}
)

func NewPRAConsoleControllerDataSource() datasource.DataSource {
	return &PRAConsoleControllerDataSource{}
}

type PRAConsoleControllerDataSource struct {
	client *client.Client
}

type PRAConsoleControllerModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	IconText        types.String `tfsdk:"icon_text"`
	CreationTime    types.String `tfsdk:"creation_time"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	MicroTenantName types.String `tfsdk:"microtenant_name"`
	PRAApplication  types.List   `tfsdk:"pra_application"`
	PRAPortals      types.List   `tfsdk:"pra_portals"`
}

func (d *PRAConsoleControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_console_controller"
}

func (d *PRAConsoleControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA PRA console controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the PRA console controller.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the PRA console controller.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"description": schema.StringAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"icon_text":   schema.StringAttribute{Computed: true},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{Computed: true},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"microtenant_name": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"pra_application": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"pra_portals": schema.ListNestedBlock{
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

func (d *PRAConsoleControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PRAConsoleControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PRAConsoleControllerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		if microID := strings.TrimSpace(data.MicroTenantID.ValueString()); microID != "" {
			service = service.WithMicroTenant(microID)
			data.MicroTenantID = types.StringValue(microID)
		}
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var console *praconsole.PRAConsole
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving PRA console controller by ID", map[string]any{"id": id})
		console, _, err = praconsole.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving PRA console controller by name", map[string]any{"name": name})
		console, _, err = praconsole.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PRA console controller: %v", err))
		return
	}

	if console == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PRA console controller with id %q or name %q was not found.", id, name))
		return
	}

	appList, appDiags := flattenPRAConsoleApplication(ctx, console.PRAApplication)
	resp.Diagnostics.Append(appDiags...)
	portalsList, portalDiags := flattenPRAConsolePortals(ctx, console.PRAPortals)
	resp.Diagnostics.Append(portalDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(console.ID)
	data.Name = types.StringValue(console.Name)
	data.Description = types.StringValue(console.Description)
	data.Enabled = types.BoolValue(console.Enabled)
	data.IconText = types.StringValue(console.IconText)
	data.CreationTime = types.StringValue(console.CreationTime)
	data.ModifiedBy = types.StringValue(console.ModifiedBy)
	data.ModifiedTime = types.StringValue(console.ModifiedTime)
	data.PRAApplication = appList
	data.PRAPortals = portalsList

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		// retain user provided value
	} else if console.MicroTenantID != "" {
		data.MicroTenantID = types.StringValue(console.MicroTenantID)
	} else {
		data.MicroTenantID = types.StringNull()
	}

	if console.MicroTenantName != "" {
		data.MicroTenantName = types.StringValue(console.MicroTenantName)
	} else {
		data.MicroTenantName = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenPRAConsoleApplication(ctx context.Context, app praconsole.PRAApplication) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if app.ID == "" && app.Name == "" {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":   types.StringValue(app.ID),
		"name": types.StringValue(app.Name),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func flattenPRAConsolePortals(ctx context.Context, portals []praconsole.PRAPortals) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if len(portals) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(portals))
	for _, portal := range portals {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(portal.ID),
			"name": types.StringValue(portal.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
