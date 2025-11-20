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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

var (
	_ datasource.DataSource              = &AppConnectorControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &AppConnectorControllerDataSource{}
)

func NewAppConnectorControllerDataSource() datasource.DataSource {
	return &AppConnectorControllerDataSource{}
}

type AppConnectorControllerDataSource struct {
	client *client.Client
}

type AppConnectorControllerModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	ApplicationStartTime             types.String `tfsdk:"application_start_time"`
	AppConnectorGroupID              types.String `tfsdk:"app_connector_group_id"`
	AppConnectorGroupName            types.String `tfsdk:"app_connector_group_name"`
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
	ProvisioningKeyID                types.String `tfsdk:"provisioning_key_id"`
	ProvisioningKeyName              types.String `tfsdk:"provisioning_key_name"`
	Platform                         types.String `tfsdk:"platform"`
	PlatformDetail                   types.String `tfsdk:"platform_detail"`
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
	AssistantVersion                 types.List   `tfsdk:"assistant_version"`
	ZPNSubModuleUpgradeList          types.List   `tfsdk:"zpn_sub_module_upgrade_list"`
}

func (d *AppConnectorControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connector_controller"
}

func (d *AppConnectorControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	assistantVersionBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id":                          schema.StringAttribute{Computed: true},
				"application_start_time":      schema.StringAttribute{Computed: true},
				"app_connector_group_id":      schema.StringAttribute{Computed: true},
				"broker_id":                   schema.StringAttribute{Computed: true},
				"creation_time":               schema.StringAttribute{Computed: true},
				"ctrl_channel_status":         schema.StringAttribute{Computed: true},
				"current_version":             schema.StringAttribute{Computed: true},
				"disable_auto_update":         schema.BoolAttribute{Computed: true},
				"expected_version":            schema.StringAttribute{Computed: true},
				"last_broker_connect_time":    schema.StringAttribute{Computed: true},
				"last_broker_disconnect_time": schema.StringAttribute{Computed: true},
				"last_upgraded_time":          schema.StringAttribute{Computed: true},
				"lone_warrior":                schema.BoolAttribute{Computed: true},
				"latitude":                    schema.StringAttribute{Computed: true},
				"longitude":                   schema.StringAttribute{Computed: true},
				"modified_by":                 schema.StringAttribute{Computed: true},
				"modified_time":               schema.StringAttribute{Computed: true},
				"platform":                    schema.StringAttribute{Computed: true},
				"platform_detail":             schema.StringAttribute{Computed: true},
				"previous_version":            schema.StringAttribute{Computed: true},
				"private_ip":                  schema.StringAttribute{Computed: true},
				"public_ip":                   schema.StringAttribute{Computed: true},
				"restart_time_in_sec":         schema.StringAttribute{Computed: true},
				"runtime_os":                  schema.StringAttribute{Computed: true},
				"sarge_version":               schema.StringAttribute{Computed: true},
				"system_start_time":           schema.StringAttribute{Computed: true},
				"mtunnel_id":                  schema.StringAttribute{Computed: true},
				"upgrade_attempt":             schema.StringAttribute{Computed: true},
				"upgrade_status":              schema.StringAttribute{Computed: true},
				"upgrade_now_once":            schema.BoolAttribute{Computed: true},
			},
		},
	}

	zpnSubModuleUpgradeBlock := schema.ListNestedBlock{
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
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves an app connector controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the app connector controller.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the app connector controller.",
			},
			"microtenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"microtenant_name": schema.StringAttribute{
				Computed:    true,
				Description: "Micro-tenant name.",
			},
			"application_start_time":               schema.StringAttribute{Computed: true},
			"app_connector_group_id":               schema.StringAttribute{Computed: true},
			"app_connector_group_name":             schema.StringAttribute{Computed: true},
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
			"provisioning_key_id":                  schema.StringAttribute{Computed: true},
			"provisioning_key_name":                schema.StringAttribute{Computed: true},
			"platform":                             schema.StringAttribute{Computed: true},
			"platform_detail":                      schema.StringAttribute{Computed: true},
			"previous_version":                     schema.StringAttribute{Computed: true},
			"private_ip":                           schema.StringAttribute{Computed: true},
			"public_ip":                            schema.StringAttribute{Computed: true},
			"runtime_os":                           schema.StringAttribute{Computed: true},
			"sarge_version":                        schema.StringAttribute{Computed: true},
			"enrollment_cert": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"upgrade_attempt": schema.StringAttribute{Computed: true},
			"upgrade_status":  schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"assistant_version":           assistantVersionBlock,
			"zpn_sub_module_upgrade_list": zpnSubModuleUpgradeBlock,
		},
	}
}

func (d *AppConnectorControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppConnectorControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before reading data sources.")
		return
	}

	var data AppConnectorControllerModel
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

	var connector *appconnectorcontroller.AppConnector
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving app connector controller by ID", map[string]any{"id": id})
		connector, _, err = appconnectorcontroller.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving app connector controller by name", map[string]any{"name": name})
		connector, _, err = appconnectorcontroller.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read app connector controller: %v", err))
		return
	}

	if connector == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("App connector controller with id %q or name %q not found.", id, name))
		return
	}

	state, diags := d.flattenAppConnectorController(ctx, connector)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *AppConnectorControllerDataSource) flattenAppConnectorController(ctx context.Context, connector *appconnectorcontroller.AppConnector) (AppConnectorControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	enrollmentCert, certDiags := types.MapValueFrom(ctx, types.StringType, connector.EnrollmentCert)
	diags.Append(certDiags...)

	assistantVersion, avDiags := d.flattenAssistantVersion(ctx, &connector.AssistantVersion)
	diags.Append(avDiags...)

	zpnSubModuleUpgradeList, zpnDiags := d.flattenZPNSubModuleUpgradeList(ctx, connector.ZPNSubModuleUpgrade)
	diags.Append(zpnDiags...)

	model := AppConnectorControllerModel{
		ID:                               types.StringValue(connector.ID),
		Name:                             types.StringValue(connector.Name),
		ApplicationStartTime:             types.StringValue(connector.ApplicationStartTime),
		AppConnectorGroupID:              types.StringValue(connector.AppConnectorGroupID),
		AppConnectorGroupName:            types.StringValue(connector.AppConnectorGroupName),
		ControlChannelStatus:             types.StringValue(connector.ControlChannelStatus),
		CreationTime:                     types.StringValue(connector.CreationTime),
		CtrlBrokerName:                   types.StringValue(connector.CtrlBrokerName),
		CurrentVersion:                   types.StringValue(connector.CurrentVersion),
		Description:                      types.StringValue(connector.Description),
		Enabled:                          types.BoolValue(connector.Enabled),
		ExpectedUpgradeTime:              types.StringValue(connector.ExpectedUpgradeTime),
		ExpectedVersion:                  types.StringValue(connector.ExpectedVersion),
		Fingerprint:                      types.StringValue(connector.Fingerprint),
		IPACL:                            types.StringValue(connector.IPACL),
		IssuedCertID:                     types.StringValue(connector.IssuedCertID),
		LastBrokerConnectTime:            types.StringValue(connector.LastBrokerConnectTime),
		LastBrokerConnectTimeDuration:    types.StringValue(connector.LastBrokerConnectTimeDuration),
		LastBrokerDisconnectTime:         types.StringValue(connector.LastBrokerDisconnectTime),
		LastBrokerDisconnectTimeDuration: types.StringValue(connector.LastBrokerDisconnectTimeDuration),
		LastUpgradeTime:                  types.StringValue(connector.LastUpgradeTime),
		Latitude:                         types.StringValue(connector.Latitude),
		Location:                         types.StringValue(connector.Location),
		Longitude:                        types.StringValue(connector.Longitude),
		ModifiedBy:                       types.StringValue(connector.ModifiedBy),
		ModifiedTime:                     types.StringValue(connector.ModifiedTime),
		ProvisioningKeyID:                types.StringValue(connector.ProvisioningKeyID),
		ProvisioningKeyName:              types.StringValue(connector.ProvisioningKeyName),
		Platform:                         types.StringValue(connector.Platform),
		PlatformDetail:                   types.StringValue(connector.PlatformDetail),
		PreviousVersion:                  types.StringValue(connector.PreviousVersion),
		PrivateIP:                        types.StringValue(connector.PrivateIP),
		PublicIP:                         types.StringValue(connector.PublicIP),
		RuntimeOS:                        types.StringValue(connector.RuntimeOS),
		SargeVersion:                     types.StringValue(connector.SargeVersion),
		EnrollmentCert:                   enrollmentCert,
		UpgradeAttempt:                   types.StringValue(connector.UpgradeAttempt),
		UpgradeStatus:                    types.StringValue(connector.UpgradeStatus),
		MicroTenantID:                    types.StringValue(connector.MicroTenantID),
		MicroTenantName:                  types.StringValue(connector.MicroTenantName),
		AssistantVersion:                 assistantVersion,
		ZPNSubModuleUpgradeList:          zpnSubModuleUpgradeList,
	}

	return model, diags
}

func (d *AppConnectorControllerDataSource) flattenAssistantVersion(ctx context.Context, version *appconnectorcontroller.AssistantVersion) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if version == nil || version.ID == "" {
		return types.ListNull(types.ObjectType{AttrTypes: assistantVersionAttrTypes()}), diags
	}

	attrTypes := assistantVersionAttrTypes()
	objValue, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":                          types.StringValue(version.ID),
		"application_start_time":      types.StringValue(version.ApplicationStartTime),
		"app_connector_group_id":      types.StringValue(version.AppConnectorGroupID),
		"broker_id":                   types.StringValue(version.BrokerId),
		"creation_time":               types.StringValue(version.CreationTime),
		"ctrl_channel_status":         types.StringValue(version.CtrlChannelStatus),
		"current_version":             types.StringValue(version.CurrentVersion),
		"disable_auto_update":         types.BoolValue(version.DisableAutoUpdate),
		"expected_version":            types.StringValue(version.ExpectedVersion),
		"last_broker_connect_time":    types.StringValue(version.LastBrokerConnectTime),
		"last_broker_disconnect_time": types.StringValue(version.LastBrokerDisconnectTime),
		"last_upgraded_time":          types.StringValue(version.LastUpgradedTime),
		"lone_warrior":                types.BoolValue(version.LoneWarrior),
		"latitude":                    types.StringValue(version.Latitude),
		"longitude":                   types.StringValue(version.Longitude),
		"modified_by":                 types.StringValue(version.ModifiedBy),
		"modified_time":               types.StringValue(version.ModifiedTime),
		"platform":                    types.StringValue(version.Platform),
		"platform_detail":             types.StringValue(version.PlatformDetail),
		"previous_version":            types.StringValue(version.PreviousVersion),
		"private_ip":                  types.StringValue(version.PrivateIP),
		"public_ip":                   types.StringValue(version.PublicIP),
		"restart_time_in_sec":         types.StringValue(version.RestartTimeInSec),
		"runtime_os":                  types.StringValue(version.RuntimeOS),
		"sarge_version":               types.StringValue(version.SargeVersion),
		"system_start_time":           types.StringValue(version.SystemStartTime),
		"mtunnel_id":                  types.StringValue(version.MtunnelID),
		"upgrade_attempt":             types.StringValue(version.UpgradeAttempt),
		"upgrade_status":              types.StringValue(version.UpgradeStatus),
		"upgrade_now_once":            types.BoolValue(version.UpgradeNowOnce),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return list, diags
}

func (d *AppConnectorControllerDataSource) flattenZPNSubModuleUpgradeList(ctx context.Context, upgrades []common.ZPNSubModuleUpgrade) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := zpnSubModuleUpgradeAttrTypes()

	if len(upgrades) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(upgrades))
	for _, upgrade := range upgrades {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":               types.StringValue(upgrade.ID),
			"creation_time":    types.StringValue(upgrade.CreationTime),
			"current_version":  types.StringValue(upgrade.CurrentVersion),
			"entity_gid":       types.StringValue(upgrade.EntityGid),
			"entity_type":      types.StringValue(upgrade.EntityType),
			"expected_version": types.StringValue(upgrade.ExpectedVersion),
			"modified_by":      types.StringValue(upgrade.ModifiedBy),
			"modified_time":    types.StringValue(upgrade.ModifiedTime),
			"previous_version": types.StringValue(upgrade.PreviousVersion),
			"role":             types.StringValue(upgrade.Role),
			"upgrade_status":   types.StringValue(upgrade.UpgradeStatus),
			"upgrade_time":     types.StringValue(upgrade.UpgradeTime),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)

	return list, diags
}

func assistantVersionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                          types.StringType,
		"application_start_time":      types.StringType,
		"app_connector_group_id":      types.StringType,
		"broker_id":                   types.StringType,
		"creation_time":               types.StringType,
		"ctrl_channel_status":         types.StringType,
		"current_version":             types.StringType,
		"disable_auto_update":         types.BoolType,
		"expected_version":            types.StringType,
		"last_broker_connect_time":    types.StringType,
		"last_broker_disconnect_time": types.StringType,
		"last_upgraded_time":          types.StringType,
		"lone_warrior":                types.BoolType,
		"latitude":                    types.StringType,
		"longitude":                   types.StringType,
		"modified_by":                 types.StringType,
		"modified_time":               types.StringType,
		"platform":                    types.StringType,
		"platform_detail":             types.StringType,
		"previous_version":            types.StringType,
		"private_ip":                  types.StringType,
		"public_ip":                   types.StringType,
		"restart_time_in_sec":         types.StringType,
		"runtime_os":                  types.StringType,
		"sarge_version":               types.StringType,
		"system_start_time":           types.StringType,
		"mtunnel_id":                  types.StringType,
		"upgrade_attempt":             types.StringType,
		"upgrade_status":              types.StringType,
		"upgrade_now_once":            types.BoolType,
	}
}

func zpnSubModuleUpgradeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"creation_time":    types.StringType,
		"current_version":  types.StringType,
		"entity_gid":       types.StringType,
		"entity_type":      types.StringType,
		"expected_version": types.StringType,
		"modified_by":      types.StringType,
		"modified_time":    types.StringType,
		"previous_version": types.StringType,
		"role":             types.StringType,
		"upgrade_status":   types.StringType,
		"upgrade_time":     types.StringType,
	}
}
