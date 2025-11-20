package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
)

var (
	_ datasource.DataSource              = &ServiceEdgeControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &ServiceEdgeControllerDataSource{}
)

func NewServiceEdgeControllerDataSource() datasource.DataSource {
	return &ServiceEdgeControllerDataSource{}
}

type ServiceEdgeControllerDataSource struct {
	client *client.Client
}

type ServiceEdgeControllerModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	ApplicationStartTime             types.String `tfsdk:"application_start_time"`
	ServiceEdgeGroupID               types.String `tfsdk:"service_edge_group_id"`
	ServiceEdgeGroupName             types.String `tfsdk:"service_edge_group_name"`
	ControlChannelStatus             types.String `tfsdk:"control_channel_status"`
	CreationTime                     types.String `tfsdk:"creation_time"`
	CtrlBrokerName                   types.String `tfsdk:"ctrl_broker_name"`
	CurrentVersion                   types.String `tfsdk:"current_version"`
	Description                      types.String `tfsdk:"description"`
	Enabled                          types.Bool   `tfsdk:"enabled"`
	ExpectedUpgradeTime              types.String `tfsdk:"expected_upgrade_time"`
	ExpectedVersion                  types.String `tfsdk:"expected_version"`
	Fingerprint                      types.String `tfsdk:"fingerprint"`
	IPACL                            types.String `tfsdk:"ip_acl"`
	IssuedCertID                     types.String `tfsdk:"issued_cert_id"`
	LastBrokerConnectTime            types.String `tfsdk:"last_broker_connect_time"`
	LastBrokerConnectTimeDuration    types.String `tfsdk:"last_broker_connect_time_duration"`
	LastBrokerDisconnectTime         types.String `tfsdk:"last_broker_disconnect_time"`
	LastBrokerDisconnectTimeDuration types.String `tfsdk:"last_broker_disconnect_time_duration"`
	LastUpgradeTime                  types.String `tfsdk:"last_upgrade_time"`
	Latitude                         types.String `tfsdk:"latitude"`
	Location                         types.String `tfsdk:"location"`
	Longitude                        types.String `tfsdk:"longitude"`
	ModifiedBy                       types.String `tfsdk:"modified_by"`
	ModifiedTime                     types.String `tfsdk:"modified_time"`
	ListenIPs                        types.Set    `tfsdk:"listen_ips"`
	PublishIPs                       types.Set    `tfsdk:"publish_ips"`
	PublishIPv6                      types.Bool   `tfsdk:"publish_ipv6"`
	ProvisioningKeyID                types.String `tfsdk:"provisioning_key_id"`
	ProvisioningKeyName              types.String `tfsdk:"provisioning_key_name"`
	Platform                         types.String `tfsdk:"platform"`
	PreviousVersion                  types.String `tfsdk:"previous_version"`
	PrivateIP                        types.String `tfsdk:"private_ip"`
	PublicIP                         types.String `tfsdk:"public_ip"`
	RuntimeOS                        types.String `tfsdk:"runtime_os"`
	SargeVersion                     types.String `tfsdk:"sarge_version"`
	EnrollmentCert                   types.Map    `tfsdk:"enrollment_cert"`
	UpgradeAttempt                   types.String `tfsdk:"upgrade_attempt"`
	UpgradeStatus                    types.String `tfsdk:"upgrade_status"`
	MicroTenantID                    types.String `tfsdk:"microtenant_id"`
	MicroTenantName                  types.String `tfsdk:"microtenant_name"`
	PrivateBrokerVersion             types.List   `tfsdk:"private_broker_version"`
}

func (d *ServiceEdgeControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_edge_controller"
}

func (d *ServiceEdgeControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	privateBrokerBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{Computed: true},
				"application_start_time": schema.StringAttribute{
					Computed: true,
				},
				"broker_id": schema.StringAttribute{Computed: true},
				"creation_time": schema.StringAttribute{
					Computed: true,
				},
				"current_version": schema.StringAttribute{
					Computed: true,
				},
				"disable_auto_update": schema.BoolAttribute{
					Computed: true,
				},
				"last_connect_time":     schema.StringAttribute{Computed: true},
				"last_disconnect_time":  schema.StringAttribute{Computed: true},
				"last_upgraded_time":    schema.StringAttribute{Computed: true},
				"lone_warrior":          schema.BoolAttribute{Computed: true},
				"modified_by":           schema.StringAttribute{Computed: true},
				"modified_time":         schema.StringAttribute{Computed: true},
				"platform":              schema.StringAttribute{Computed: true},
				"platform_detail":       schema.StringAttribute{Computed: true},
				"previous_version":      schema.StringAttribute{Computed: true},
				"service_edge_group_id": schema.StringAttribute{Computed: true},
				"private_ip":            schema.StringAttribute{Computed: true},
				"public_ip":             schema.StringAttribute{Computed: true},
				"restart_instructions":  schema.StringAttribute{Computed: true},
				"restart_time_in_sec":   schema.StringAttribute{Computed: true},
				"runtime_os":            schema.StringAttribute{Computed: true},
				"sarge_version":         schema.StringAttribute{Computed: true},
				"system_start_time":     schema.StringAttribute{Computed: true},
				"tunnel_id":             schema.StringAttribute{Computed: true},
				"upgrade_attempt":       schema.StringAttribute{Computed: true},
				"upgrade_status":        schema.StringAttribute{Computed: true},
				"upgrade_now_once":      schema.BoolAttribute{Computed: true},
			},
			Blocks: map[string]schema.Block{
				"zpn_sub_module_upgrade": schema.ListNestedBlock{
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"id":               schema.StringAttribute{Computed: true},
							"creation_time":    schema.StringAttribute{Computed: true},
							"current_version":  schema.StringAttribute{Computed: true},
							"entity_gid":       schema.StringAttribute{Computed: true},
							"entity_type":      schema.StringAttribute{Computed: true},
							"expected_version": schema.StringAttribute{Computed: true},
							"modified_by":      schema.StringAttribute{Computed: true},
							"modified_time":    schema.StringAttribute{Computed: true},
							"previous_version": schema.StringAttribute{Computed: true},
							"role":             schema.StringAttribute{Computed: true},
							"upgrade_status":   schema.StringAttribute{Computed: true},
							"upgrade_time":     schema.StringAttribute{Computed: true},
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves a service edge controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the service edge controller.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the service edge controller.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"application_start_time":               schema.StringAttribute{Computed: true},
			"service_edge_group_id":                schema.StringAttribute{Computed: true},
			"service_edge_group_name":              schema.StringAttribute{Computed: true},
			"control_channel_status":               schema.StringAttribute{Computed: true},
			"creation_time":                        schema.StringAttribute{Computed: true},
			"ctrl_broker_name":                     schema.StringAttribute{Computed: true},
			"current_version":                      schema.StringAttribute{Computed: true},
			"description":                          schema.StringAttribute{Computed: true},
			"enabled":                              schema.BoolAttribute{Computed: true},
			"expected_upgrade_time":                schema.StringAttribute{Computed: true},
			"expected_version":                     schema.StringAttribute{Computed: true},
			"fingerprint":                          schema.StringAttribute{Computed: true},
			"ip_acl":                               schema.StringAttribute{Computed: true},
			"issued_cert_id":                       schema.StringAttribute{Computed: true},
			"last_broker_connect_time":             schema.StringAttribute{Computed: true},
			"last_broker_connect_time_duration":    schema.StringAttribute{Computed: true},
			"last_broker_disconnect_time":          schema.StringAttribute{Computed: true},
			"last_broker_disconnect_time_duration": schema.StringAttribute{Computed: true},
			"last_upgrade_time":                    schema.StringAttribute{Computed: true},
			"latitude":                             schema.StringAttribute{Computed: true},
			"location":                             schema.StringAttribute{Computed: true},
			"longitude":                            schema.StringAttribute{Computed: true},
			"modified_by":                          schema.StringAttribute{Computed: true},
			"modified_time":                        schema.StringAttribute{Computed: true},
			"listen_ips": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"publish_ips": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"publish_ipv6":          schema.BoolAttribute{Computed: true},
			"provisioning_key_id":   schema.StringAttribute{Computed: true},
			"provisioning_key_name": schema.StringAttribute{Computed: true},
			"platform":              schema.StringAttribute{Computed: true},
			"previous_version":      schema.StringAttribute{Computed: true},
			"private_ip":            schema.StringAttribute{Computed: true},
			"public_ip":             schema.StringAttribute{Computed: true},
			"runtime_os":            schema.StringAttribute{Computed: true},
			"sarge_version":         schema.StringAttribute{Computed: true},
			"enrollment_cert": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"upgrade_attempt":  schema.StringAttribute{Computed: true},
			"upgrade_status":   schema.StringAttribute{Computed: true},
			"microtenant_name": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"private_broker_version": privateBrokerBlock,
		},
	}
}

func (d *ServiceEdgeControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceEdgeControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data ServiceEdgeControllerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
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

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	var edge *serviceedgecontroller.ServiceEdgeController
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving service edge controller by ID", map[string]any{"id": id})
		edge, _, err = serviceedgecontroller.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving service edge controller by name", map[string]any{"name": name})
		edge, _, err = serviceedgecontroller.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service edge controller: %v", err))
		return
	}

	if edge == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Service edge controller with id %q or name %q not found.", id, name))
		return
	}

	state, diags := flattenServiceEdgeController(ctx, edge)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenServiceEdgeController(ctx context.Context, edge *serviceedgecontroller.ServiceEdgeController) (ServiceEdgeControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	listenIPs, listenDiags := types.SetValueFrom(ctx, types.StringType, edge.ListenIPs)
	diags.Append(listenDiags...)

	publishIPs, publishDiags := types.SetValueFrom(ctx, types.StringType, edge.PublishIPs)
	diags.Append(publishDiags...)

	enrollmentCert, certDiags := types.MapValueFrom(ctx, types.StringType, edge.EnrollmentCert)
	diags.Append(certDiags...)

	privateBrokerVersion, pbDiags := helpers.FlattenPrivateBrokerVersionToList(ctx, edge.PrivateBrokerVersion)
	diags.Append(pbDiags...)

	model := ServiceEdgeControllerModel{
		ID:                               types.StringValue(edge.ID),
		Name:                             types.StringValue(edge.Name),
		ApplicationStartTime:             types.StringValue(edge.ApplicationStartTime),
		ServiceEdgeGroupID:               types.StringValue(edge.ServiceEdgeGroupID),
		ServiceEdgeGroupName:             types.StringValue(edge.ServiceEdgeGroupName),
		ControlChannelStatus:             types.StringValue(edge.ControlChannelStatus),
		CreationTime:                     types.StringValue(edge.CreationTime),
		CtrlBrokerName:                   types.StringValue(edge.CtrlBrokerName),
		CurrentVersion:                   types.StringValue(edge.CurrentVersion),
		Description:                      types.StringValue(edge.Description),
		Enabled:                          types.BoolValue(edge.Enabled),
		ExpectedUpgradeTime:              types.StringValue(edge.ExpectedUpgradeTime),
		ExpectedVersion:                  types.StringValue(edge.ExpectedVersion),
		Fingerprint:                      types.StringValue(edge.Fingerprint),
		IPACL:                            types.StringValue(edge.IPACL),
		IssuedCertID:                     types.StringValue(edge.IssuedCertID),
		LastBrokerConnectTime:            types.StringValue(edge.LastBrokerConnectTime),
		LastBrokerConnectTimeDuration:    types.StringValue(edge.LastBrokerConnectTimeDuration),
		LastBrokerDisconnectTime:         types.StringValue(edge.LastBrokerDisconnectTime),
		LastBrokerDisconnectTimeDuration: types.StringValue(edge.LastBrokerDisconnectTimeDuration),
		LastUpgradeTime:                  types.StringValue(edge.LastUpgradeTime),
		Latitude:                         types.StringValue(edge.Latitude),
		Location:                         types.StringValue(edge.Location),
		Longitude:                        types.StringValue(edge.Longitude),
		ModifiedBy:                       types.StringValue(edge.ModifiedBy),
		ModifiedTime:                     types.StringValue(edge.ModifiedTime),
		ListenIPs:                        listenIPs,
		PublishIPs:                       publishIPs,
		PublishIPv6:                      types.BoolValue(edge.PublishIPv6),
		ProvisioningKeyID:                types.StringValue(edge.ProvisioningKeyID),
		ProvisioningKeyName:              types.StringValue(edge.ProvisioningKeyName),
		Platform:                         types.StringValue(edge.Platform),
		PreviousVersion:                  types.StringValue(edge.PreviousVersion),
		PrivateIP:                        types.StringValue(edge.PrivateIP),
		PublicIP:                         types.StringValue(edge.PublicIP),
		RuntimeOS:                        types.StringValue(edge.RuntimeOS),
		SargeVersion:                     types.StringValue(edge.SargeVersion),
		EnrollmentCert:                   enrollmentCert,
		UpgradeAttempt:                   types.StringValue(edge.UpgradeAttempt),
		UpgradeStatus:                    types.StringValue(edge.UpgradeStatus),
		MicroTenantID:                    types.StringValue(edge.MicroTenantID),
		MicroTenantName:                  types.StringValue(edge.MicroTenantName),
		PrivateBrokerVersion:             privateBrokerVersion,
	}

	return model, diags
}
