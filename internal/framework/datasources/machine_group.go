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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/machinegroup"
)

var (
	_ datasource.DataSource              = &MachineGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &MachineGroupDataSource{}
)

func NewMachineGroupDataSource() datasource.DataSource {
	return &MachineGroupDataSource{}
}

type MachineGroupDataSource struct {
	client *client.Client
}

type MachineGroupModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	CreationTime    types.String `tfsdk:"creation_time"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	MicroTenantID   types.String `tfsdk:"microtenant_id"`
	MicroTenantName types.String `tfsdk:"microtenant_name"`
	Machines        types.List   `tfsdk:"machines"`
}

func (d *MachineGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_machine_group"
}

func (d *MachineGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	machineAttrTypes := map[string]schema.Attribute{
		"creation_time": schema.StringAttribute{Computed: true},
		"description":   schema.StringAttribute{Computed: true},
		"fingerprint":   schema.StringAttribute{Computed: true},
		"id":            schema.StringAttribute{Computed: true},
		"issued_cert_id": schema.StringAttribute{
			Computed: true,
		},
		"machine_group_id":   schema.StringAttribute{Computed: true},
		"machine_group_name": schema.StringAttribute{Computed: true},
		"machine_token_id":   schema.StringAttribute{Computed: true},
		"modified_by":        schema.StringAttribute{Computed: true},
		"modified_time":      schema.StringAttribute{Computed: true},
		"name":               schema.StringAttribute{Computed: true},
		"signing_cert": schema.MapAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"microtenant_id":   schema.StringAttribute{Computed: true},
		"microtenant_name": schema.StringAttribute{Computed: true},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA machine group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the machine group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the machine group.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"microtenant_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"machines": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: machineAttrTypes,
				},
			},
		},
	}
}

func (d *MachineGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MachineGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data MachineGroupModel
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

	service := d.client.Service
	if !data.MicroTenantID.IsNull() && !data.MicroTenantID.IsUnknown() {
		microID := strings.TrimSpace(data.MicroTenantID.ValueString())
		if microID != "" {
			service = service.WithMicroTenant(microID)
			data.MicroTenantID = types.StringValue(microID)
		}
	}

	var (
		group *machinegroup.MachineGroup
		err   error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving machine group by ID", map[string]any{"id": id})
		group, _, err = machinegroup.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving machine group by name", map[string]any{"name": name})
		group, _, err = machinegroup.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read machine group: %v", err))
		return
	}

	if group == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Machine group with id %q or name %q not found.", id, name))
		return
	}

	state, diags := flattenMachineGroup(ctx, group)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenMachineGroup(ctx context.Context, group *machinegroup.MachineGroup) (MachineGroupModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	machineAttrTypes := map[string]attr.Type{
		"creation_time":      types.StringType,
		"description":        types.StringType,
		"fingerprint":        types.StringType,
		"id":                 types.StringType,
		"issued_cert_id":     types.StringType,
		"machine_group_id":   types.StringType,
		"machine_group_name": types.StringType,
		"machine_token_id":   types.StringType,
		"modified_by":        types.StringType,
		"modified_time":      types.StringType,
		"name":               types.StringType,
		"signing_cert":       types.MapType{ElemType: types.StringType},
		"microtenant_id":     types.StringType,
		"microtenant_name":   types.StringType,
	}

	machineElements := make([]attr.Value, 0, len(group.Machines))
	for _, machine := range group.Machines {
		signingCert, certDiags := mapStringInterfaceToMap(machine.SigningCert)
		diags.Append(certDiags...)

		obj, objDiags := types.ObjectValue(machineAttrTypes, map[string]attr.Value{
			"creation_time":      types.StringValue(machine.CreationTime),
			"description":        types.StringValue(machine.Description),
			"fingerprint":        types.StringValue(machine.Fingerprint),
			"id":                 types.StringValue(machine.ID),
			"issued_cert_id":     types.StringValue(machine.IssuedCertID),
			"machine_group_id":   types.StringValue(machine.MachineGroupID),
			"machine_group_name": types.StringValue(machine.MachineGroupName),
			"machine_token_id":   types.StringValue(machine.MachineTokenID),
			"modified_by":        types.StringValue(machine.ModifiedBy),
			"modified_time":      types.StringValue(machine.ModifiedTime),
			"name":               types.StringValue(machine.Name),
			"signing_cert":       signingCert,
			"microtenant_id":     types.StringValue(machine.MicroTenantID),
			"microtenant_name":   types.StringValue(machine.MicroTenantName),
		})
		diags.Append(objDiags...)
		machineElements = append(machineElements, obj)
	}

	machinesList, listDiags := types.ListValue(types.ObjectType{AttrTypes: machineAttrTypes}, machineElements)
	diags.Append(listDiags...)

	model := MachineGroupModel{
		ID:              types.StringValue(group.ID),
		Name:            types.StringValue(group.Name),
		Description:     types.StringValue(group.Description),
		Enabled:         types.BoolValue(group.Enabled),
		CreationTime:    types.StringValue(group.CreationTime),
		ModifiedBy:      types.StringValue(group.ModifiedBy),
		ModifiedTime:    types.StringValue(group.ModifiedTime),
		MicroTenantID:   types.StringValue(group.MicroTenantID),
		MicroTenantName: types.StringValue(group.MicroTenantName),
		Machines:        machinesList,
	}

	return model, diags
}
