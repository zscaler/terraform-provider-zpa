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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_controller"
)

var (
	_ datasource.DataSource              = &PrivateCloudControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &PrivateCloudControllerDataSource{}
)

func NewPrivateCloudControllerDataSource() datasource.DataSource {
	return &PrivateCloudControllerDataSource{}
}

type PrivateCloudControllerDataSource struct {
	client *client.Client
}

type PrivateCloudControllerModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	ApplicationStartTime             types.String `tfsdk:"application_start_time"`
	ControlChannelStatus             types.String `tfsdk:"control_channel_status"`
	CreationTime                     types.String `tfsdk:"creation_time"`
	CtrlBrokerName                   types.String `tfsdk:"ctrl_broker_name"`
	CurrentVersion                   types.String `tfsdk:"current_version"`
	Description                      types.String `tfsdk:"description"`
	Enabled                          types.Bool   `tfsdk:"enabled"`
	ExpectedSargeVersion             types.String `tfsdk:"expected_sarge_version"`
	ExpectedUpgradeTime              types.String `tfsdk:"expected_upgrade_time"`
	ExpectedVersion                  types.String `tfsdk:"expected_version"`
	Fingerprint                      types.String `tfsdk:"fingerprint"`
	IPACL                            types.List   `tfsdk:"ip_acl"`
	IssuedCertID                     types.String `tfsdk:"issued_cert_id"`
	LastBrokerConnectTime            types.String `tfsdk:"last_broker_connect_time"`
	LastBrokerConnectTimeDuration    types.String `tfsdk:"last_broker_connect_time_duration"`
	LastBrokerDisconnectTime         types.String `tfsdk:"last_broker_disconnect_time"`
	LastBrokerDisconnectTimeDuration types.String `tfsdk:"last_broker_disconnect_time_duration"`
	LastOSUpgradeTime                types.String `tfsdk:"last_os_upgrade_time"`
	LastSargeUpgradeTime             types.String `tfsdk:"last_sarge_upgrade_time"`
	LastUpgradeTime                  types.String `tfsdk:"last_upgrade_time"`
	Latitude                         types.String `tfsdk:"latitude"`
	ListenIPs                        types.List   `tfsdk:"listen_ips"`
	Location                         types.String `tfsdk:"location"`
	Longitude                        types.String `tfsdk:"longitude"`
	MasterLastSyncTime               types.String `tfsdk:"master_last_sync_time"`
	ModifiedBy                       types.String `tfsdk:"modified_by"`
	ModifiedTime                     types.String `tfsdk:"modified_time"`
	ProvisioningKeyID                types.String `tfsdk:"provisioning_key_id"`
	ProvisioningKeyName              types.String `tfsdk:"provisioning_key_name"`
	OSUpgradeEnabled                 types.Bool   `tfsdk:"os_upgrade_enabled"`
	OSUpgradeStatus                  types.String `tfsdk:"os_upgrade_status"`
	Platform                         types.String `tfsdk:"platform"`
	PlatformDetail                   types.String `tfsdk:"platform_detail"`
	PlatformVersion                  types.String `tfsdk:"platform_version"`
	PreviousVersion                  types.String `tfsdk:"previous_version"`
	PrivateIP                        types.String `tfsdk:"private_ip"`
	PublicIP                         types.String `tfsdk:"public_ip"`
	PublishIPs                       types.List   `tfsdk:"publish_ips"`
	ReadOnly                         types.Bool   `tfsdk:"read_only"`
	RestrictionType                  types.String `tfsdk:"restriction_type"`
	Runtime                          types.String `tfsdk:"runtime"`
	SargeUpgradeAttempt              types.String `tfsdk:"sarge_upgrade_attempt"`
	SargeUpgradeStatus               types.String `tfsdk:"sarge_upgrade_status"`
	SargeVersion                     types.String `tfsdk:"sarge_version"`
	MicroTenantID                    types.String `tfsdk:"microtenant_id"`
	MicroTenantName                  types.String `tfsdk:"microtenant_name"`
	ShardLastSyncTime                types.String `tfsdk:"shard_last_sync_time"`
	EnrollmentCert                   types.Map    `tfsdk:"enrollment_cert"`
	PrivateCloudControllerGroupID    types.String `tfsdk:"private_cloud_controller_group_id"`
	PrivateCloudControllerGroupName  types.String `tfsdk:"private_cloud_controller_group_name"`
	PrivateCloudControllerVersion    types.Map    `tfsdk:"private_cloud_controller_version"`
	SiteSPDNSName                    types.String `tfsdk:"site_sp_dns_name"`
	UpgradeAttempt                   types.String `tfsdk:"upgrade_attempt"`
	UpgradeStatus                    types.String `tfsdk:"upgrade_status"`
	UserDBLastSyncTime               types.String `tfsdk:"userdb_last_sync_time"`
	ZPNSubModuleUpgradeList          types.List   `tfsdk:"zpn_sub_module_upgrade_list"`
	ZscalerManaged                   types.Bool   `tfsdk:"zscaler_managed"`
}

func (d *PrivateCloudControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_cloud_controller"
}

func (d *PrivateCloudControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA Private Cloud Controller by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the private cloud controller.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the private cloud controller.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"application_start_time": schema.StringAttribute{Computed: true},
			"control_channel_status": schema.StringAttribute{Computed: true},
			"creation_time":          schema.StringAttribute{Computed: true},
			"ctrl_broker_name":       schema.StringAttribute{Computed: true},
			"current_version":        schema.StringAttribute{Computed: true},
			"description":            schema.StringAttribute{Computed: true},
			"enabled":                schema.BoolAttribute{Computed: true},
			"expected_sarge_version": schema.StringAttribute{Computed: true},
			"expected_upgrade_time":  schema.StringAttribute{Computed: true},
			"expected_version":       schema.StringAttribute{Computed: true},
			"fingerprint":            schema.StringAttribute{Computed: true},
			"ip_acl": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"issued_cert_id": schema.StringAttribute{Computed: true},
			"last_broker_connect_time": schema.StringAttribute{
				Computed: true,
			},
			"last_broker_connect_time_duration": schema.StringAttribute{
				Computed: true,
			},
			"last_broker_disconnect_time": schema.StringAttribute{
				Computed: true,
			},
			"last_broker_disconnect_time_duration": schema.StringAttribute{
				Computed: true,
			},
			"last_os_upgrade_time":    schema.StringAttribute{Computed: true},
			"last_sarge_upgrade_time": schema.StringAttribute{Computed: true},
			"last_upgrade_time":       schema.StringAttribute{Computed: true},
			"latitude":                schema.StringAttribute{Computed: true},
			"listen_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"location":              schema.StringAttribute{Computed: true},
			"longitude":             schema.StringAttribute{Computed: true},
			"master_last_sync_time": schema.StringAttribute{Computed: true},
			"modified_by":           schema.StringAttribute{Computed: true},
			"modified_time":         schema.StringAttribute{Computed: true},
			"provisioning_key_id":   schema.StringAttribute{Computed: true},
			"provisioning_key_name": schema.StringAttribute{Computed: true},
			"os_upgrade_enabled":    schema.BoolAttribute{Computed: true},
			"os_upgrade_status":     schema.StringAttribute{Computed: true},
			"platform":              schema.StringAttribute{Computed: true},
			"platform_detail":       schema.StringAttribute{Computed: true},
			"platform_version":      schema.StringAttribute{Computed: true},
			"previous_version":      schema.StringAttribute{Computed: true},
			"private_ip":            schema.StringAttribute{Computed: true},
			"public_ip":             schema.StringAttribute{Computed: true},
			"publish_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"read_only":             schema.BoolAttribute{Computed: true},
			"restriction_type":      schema.StringAttribute{Computed: true},
			"runtime":               schema.StringAttribute{Computed: true},
			"sarge_upgrade_attempt": schema.StringAttribute{Computed: true},
			"sarge_upgrade_status":  schema.StringAttribute{Computed: true},
			"sarge_version":         schema.StringAttribute{Computed: true},
			"microtenant_name":      schema.StringAttribute{Computed: true},
			"shard_last_sync_time":  schema.StringAttribute{Computed: true},
			"enrollment_cert": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"private_cloud_controller_group_id":   schema.StringAttribute{Computed: true},
			"private_cloud_controller_group_name": schema.StringAttribute{Computed: true},
			"private_cloud_controller_version": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"site_sp_dns_name":      schema.StringAttribute{Computed: true},
			"upgrade_attempt":       schema.StringAttribute{Computed: true},
			"upgrade_status":        schema.StringAttribute{Computed: true},
			"userdb_last_sync_time": schema.StringAttribute{Computed: true},
			"zpn_sub_module_upgrade_list": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"zscaler_managed": schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *PrivateCloudControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PrivateCloudControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PrivateCloudControllerModel
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

	var controller *private_cloud_controller.PrivateCloudController
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving private cloud controller by ID", map[string]any{"id": id})
		controller, _, err = private_cloud_controller.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving private cloud controller by name", map[string]any{"name": name})
		controller, _, err = private_cloud_controller.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read private cloud controller: %v", err))
		return
	}

	if controller == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Private cloud controller with id %q or name %q was not found.", id, name))
		return
	}

	state, diags := flattenPrivateCloudController(ctx, controller)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		state.MicroTenantID = data.MicroTenantID
	} else if controller.MicrotenantId != "" {
		state.MicroTenantID = types.StringValue(controller.MicrotenantId)
	} else {
		state.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenPrivateCloudController(ctx context.Context, controller *private_cloud_controller.PrivateCloudController) (PrivateCloudControllerModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	ipACL, d := types.ListValueFrom(ctx, types.StringType, controller.IpAcl)
	diags.Append(d...)

	listenIPs, d := types.ListValueFrom(ctx, types.StringType, controller.ListenIps)
	diags.Append(d...)

	publishIPs, d := types.ListValueFrom(ctx, types.StringType, controller.PublishIps)
	diags.Append(d...)

	zpnValues := make([]string, 0, len(controller.ZpnSubModuleUpgradeList))
	for _, item := range controller.ZpnSubModuleUpgradeList {
		zpnValues = append(zpnValues, fmt.Sprintf("%v", item))
	}
	zpnList, d := types.ListValueFrom(ctx, types.StringType, zpnValues)
	diags.Append(d...)

	enrollment, mapDiags := mapStringInterfaceToMap(controller.EnrollmentCert)
	diags.Append(mapDiags...)

	versionMap, versionDiags := mapStringInterfaceToMap(controller.PrivateCloudControllerVersion)
	diags.Append(versionDiags...)

	state := PrivateCloudControllerModel{
		ID:                               types.StringValue(controller.Id),
		Name:                             types.StringValue(controller.Name),
		ApplicationStartTime:             types.StringValue(controller.ApplicationStartTime),
		ControlChannelStatus:             types.StringValue(controller.ControlChannelStatus),
		CreationTime:                     types.StringValue(controller.CreationTime),
		CtrlBrokerName:                   types.StringValue(controller.CtrlBrokerName),
		CurrentVersion:                   types.StringValue(controller.CurrentVersion),
		Description:                      types.StringValue(controller.Description),
		Enabled:                          types.BoolValue(controller.Enabled),
		ExpectedSargeVersion:             types.StringValue(controller.ExpectedSargeVersion),
		ExpectedUpgradeTime:              types.StringValue(controller.ExpectedUpgradeTime),
		ExpectedVersion:                  types.StringValue(controller.ExpectedVersion),
		Fingerprint:                      types.StringValue(controller.Fingerprint),
		IPACL:                            ipACL,
		IssuedCertID:                     types.StringValue(controller.IssuedCertId),
		LastBrokerConnectTime:            types.StringValue(controller.LastBrokerConnectTime),
		LastBrokerConnectTimeDuration:    types.StringValue(controller.LastBrokerConnectTimeDuration),
		LastBrokerDisconnectTime:         types.StringValue(controller.LastBrokerDisconnectTime),
		LastBrokerDisconnectTimeDuration: types.StringValue(controller.LastBrokerDisconnectTimeDuration),
		LastOSUpgradeTime:                types.StringValue(controller.LastOsUpgradeTime),
		LastSargeUpgradeTime:             types.StringValue(controller.LastSargeUpgradeTime),
		LastUpgradeTime:                  types.StringValue(controller.LastUpgradeTime),
		Latitude:                         types.StringValue(controller.Latitude),
		ListenIPs:                        listenIPs,
		Location:                         types.StringValue(controller.Location),
		Longitude:                        types.StringValue(controller.Longitude),
		MasterLastSyncTime:               types.StringValue(controller.MasterLastSyncTime),
		ModifiedBy:                       types.StringValue(controller.ModifiedBy),
		ModifiedTime:                     types.StringValue(controller.ModifiedTime),
		ProvisioningKeyID:                types.StringValue(controller.ProvisioningKeyId),
		ProvisioningKeyName:              types.StringValue(controller.ProvisioningKeyName),
		OSUpgradeEnabled:                 types.BoolValue(controller.OsUpgradeEnabled),
		OSUpgradeStatus:                  types.StringValue(controller.OsUpgradeStatus),
		Platform:                         types.StringValue(controller.Platform),
		PlatformDetail:                   types.StringValue(controller.PlatformDetail),
		PlatformVersion:                  types.StringValue(controller.PlatformVersion),
		PreviousVersion:                  types.StringValue(controller.PreviousVersion),
		PrivateIP:                        types.StringValue(controller.PrivateIp),
		PublicIP:                         types.StringValue(controller.PublicIp),
		PublishIPs:                       publishIPs,
		ReadOnly:                         types.BoolValue(controller.ReadOnly),
		RestrictionType:                  types.StringValue(controller.RestrictionType),
		Runtime:                          types.StringValue(controller.Runtime),
		SargeUpgradeAttempt:              types.StringValue(controller.SargeUpgradeAttempt),
		SargeUpgradeStatus:               types.StringValue(controller.SargeUpgradeStatus),
		SargeVersion:                     types.StringValue(controller.SargeVersion),
		MicroTenantName:                  types.StringValue(controller.MicrotenantName),
		ShardLastSyncTime:                types.StringValue(controller.ShardLastSyncTime),
		EnrollmentCert:                   enrollment,
		PrivateCloudControllerGroupID:    types.StringValue(controller.PrivateCloudControllerGroupId),
		PrivateCloudControllerGroupName:  types.StringValue(controller.PrivateCloudControllerGroupName),
		PrivateCloudControllerVersion:    versionMap,
		SiteSPDNSName:                    types.StringValue(controller.SiteSpDnsName),
		UpgradeAttempt:                   types.StringValue(controller.UpgradeAttempt),
		UpgradeStatus:                    types.StringValue(controller.UpgradeStatus),
		UserDBLastSyncTime:               types.StringValue(controller.UserdbLastSyncTime),
		ZPNSubModuleUpgradeList:          zpnList,
		ZscalerManaged:                   types.BoolValue(controller.ZscalerManaged),
	}

	return state, diags
}
