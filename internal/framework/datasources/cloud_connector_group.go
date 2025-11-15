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
	cloudconnectorgroup "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector_group"
)

var (
	_ datasource.DataSource              = &CloudConnectorGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudConnectorGroupDataSource{}
)

func NewCloudConnectorGroupDataSource() datasource.DataSource {
	return &CloudConnectorGroupDataSource{}
}

type CloudConnectorGroupDataSource struct {
	client *client.Client
}

type CloudConnectorGroupModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	CreationTime    types.String `tfsdk:"creation_time"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	GeolocationID   types.String `tfsdk:"geolocation_id"`
	ModifiedBy      types.String `tfsdk:"modified_by"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	ZiaCloud        types.String `tfsdk:"zia_cloud"`
	ZiaOrgID        types.String `tfsdk:"zia_org_id"`
	ZnfGroupType    types.String `tfsdk:"znf_group_type"`
	CloudConnectors types.List   `tfsdk:"cloud_connectors"`
}

func (d *CloudConnectorGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_connector_group"
}

func (d *CloudConnectorGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA cloud connector group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the cloud connector group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the cloud connector group.",
			},
			"description":    schema.StringAttribute{Computed: true},
			"creation_time":  schema.StringAttribute{Computed: true},
			"enabled":        schema.BoolAttribute{Computed: true},
			"geolocation_id": schema.StringAttribute{Computed: true},
			"modified_by":    schema.StringAttribute{Computed: true},
			"modified_time":  schema.StringAttribute{Computed: true},
			"zia_cloud":      schema.StringAttribute{Computed: true},
			"zia_org_id":     schema.StringAttribute{Computed: true},
			"znf_group_type": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"cloud_connectors": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"creation_time": schema.StringAttribute{Computed: true},
						"description":   schema.StringAttribute{Computed: true},
						"enabled":       schema.BoolAttribute{Computed: true},
						"fingerprint":   schema.StringAttribute{Computed: true},
						"ipacl": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"issued_cert_id": schema.StringAttribute{Computed: true},
						"modified_by":    schema.StringAttribute{Computed: true},
						"modified_time":  schema.StringAttribute{Computed: true},
						"signing_cert": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"microtenant_id":   schema.StringAttribute{Computed: true},
						"microtenant_name": schema.StringAttribute{Computed: true},
						"read_only":        schema.BoolAttribute{Computed: true},
						"restriction_type": schema.StringAttribute{Computed: true},
						"zscaler_managed":  schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *CloudConnectorGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudConnectorGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data CloudConnectorGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided to read a cloud connector group.")
		return
	}

	var (
		group *cloudconnectorgroup.CloudConnectorGroup
		err   error
	)

	if id != "" {
		tflog.Debug(ctx, "Retrieving cloud connector group by ID", map[string]any{"id": id})
		group, _, err = cloudconnectorgroup.Get(ctx, d.client.Service, id)
	} else {
		tflog.Debug(ctx, "Retrieving cloud connector group by name", map[string]any{"name": name})
		group, _, err = cloudconnectorgroup.GetByName(ctx, d.client.Service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud connector group: %v", err))
		return
	}

	if group == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Cloud connector group with id %q or name %q was not found.", id, name))
		return
	}

	connectors, connDiags := flattenCloudConnectors(ctx, group.CloudConnectors)
	resp.Diagnostics.Append(connDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(group.ID)
	data.Name = stringOrNull(group.Name)
	data.Description = stringOrNull(group.Description)
	data.CreationTime = stringOrNull(group.CreationTime)
	data.Enabled = types.BoolValue(group.Enabled)
	data.GeolocationID = stringOrNull(group.GeolocationID)
	data.ModifiedBy = stringOrNull(group.ModifiedBy)
	data.ModifiedTime = stringOrNull(group.ModifiedTime)
	data.ZiaCloud = stringOrNull(group.ZiaCloud)
	data.ZiaOrgID = stringOrNull(group.ZiaOrgid)
	data.ZnfGroupType = stringOrNull(group.ZnfGroupType)
	data.CloudConnectors = connectors

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenCloudConnectors(ctx context.Context, connectors []cloudconnectorgroup.CloudConnectors) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"creation_time":    types.StringType,
		"description":      types.StringType,
		"enabled":          types.BoolType,
		"fingerprint":      types.StringType,
		"ipacl":            types.ListType{ElemType: types.StringType},
		"issued_cert_id":   types.StringType,
		"modified_by":      types.StringType,
		"modified_time":    types.StringType,
		"signing_cert":     types.MapType{ElemType: types.StringType},
		"microtenant_id":   types.StringType,
		"microtenant_name": types.StringType,
		"read_only":        types.BoolType,
		"restriction_type": types.StringType,
		"zscaler_managed":  types.BoolType,
	}

	if len(connectors) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(connectors))
	var diags diag.Diagnostics
	for _, connector := range connectors {
		ipACL := types.ListNull(types.StringType)
		if len(connector.IPACL) > 0 {
			list, listDiags := types.ListValueFrom(ctx, types.StringType, connector.IPACL)
			diags.Append(listDiags...)
			ipACL = list
		}

		signingCert := types.MapNull(types.StringType)
		if len(connector.SigningCert) > 0 {
			mapValue, mapDiags := mapInterfaceToStringMap(ctx, connector.SigningCert)
			diags.Append(mapDiags...)
			signingCert = mapValue
		}

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":               stringOrNull(connector.ID),
			"name":             stringOrNull(connector.Name),
			"creation_time":    stringOrNull(connector.CreationTime),
			"description":      stringOrNull(connector.Description),
			"enabled":          types.BoolValue(connector.Enabled),
			"fingerprint":      stringOrNull(connector.Fingerprint),
			"ipacl":            ipACL,
			"issued_cert_id":   stringOrNull(connector.IssuedCertID),
			"modified_by":      stringOrNull(connector.ModifiedBy),
			"modified_time":    stringOrNull(connector.ModifiedTime),
			"signing_cert":     signingCert,
			"microtenant_id":   stringOrNull(connector.MicroTenantID),
			"microtenant_name": stringOrNull(connector.MicroTenantName),
			"read_only":        types.BoolValue(connector.ReadOnly),
			"restriction_type": stringOrNull(connector.RestrictionType),
			"zscaler_managed":  types.BoolValue(connector.ZscalerManaged),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
