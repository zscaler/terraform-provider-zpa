package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ApplicationSegmentInspectionResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSegmentInspectionResource{}
	_ resource.ResourceWithImportState = &ApplicationSegmentInspectionResource{}
)

func NewApplicationSegmentInspectionResource() resource.Resource {
	return &ApplicationSegmentInspectionResource{}
}

type ApplicationSegmentInspectionResource struct {
	client *client.Client
}

type ApplicationSegmentInspectionModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	SegmentGroupID            types.String `tfsdk:"segment_group_id"`
	SegmentGroupName          types.String `tfsdk:"segment_group_name"`
	BypassType                types.String `tfsdk:"bypass_type"`
	BypassOnReauth            types.Bool   `tfsdk:"bypass_on_reauth"`
	ConfigSpace               types.String `tfsdk:"config_space"`
	DomainNames               types.List   `tfsdk:"domain_names"`
	DoubleEncrypt             types.Bool   `tfsdk:"double_encrypt"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	PassiveHealthEnabled      types.Bool   `tfsdk:"passive_health_enabled"`
	HealthCheckType           types.String `tfsdk:"health_check_type"`
	HealthReporting           types.String `tfsdk:"health_reporting"`
	ICMPAccessType            types.String `tfsdk:"icmp_access_type"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	TCPKeepAlive              types.String `tfsdk:"tcp_keep_alive"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	AdpEnabled                types.Bool   `tfsdk:"adp_enabled"`
	AutoAppProtectEnabled     types.Bool   `tfsdk:"auto_app_protect_enabled"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.List   `tfsdk:"tcp_port_range"`
	UDPPortRange              types.List   `tfsdk:"udp_port_range"`
	TCPProtocols              types.List   `tfsdk:"tcp_protocols"`
	UDPProtocols              types.List   `tfsdk:"udp_protocols"`
	ServerGroups              types.List   `tfsdk:"server_groups"`
	CommonAppsDto             types.List   `tfsdk:"common_apps_dto"`
	InspectionApps            types.List   `tfsdk:"inspection_apps"`
}

type inspectionCommonAppsModel struct {
	AppsConfig types.List `tfsdk:"apps_config"`
}

type inspectionAppsConfigModel struct {
	AppID               types.String `tfsdk:"app_id"`
	InspectAppID        types.String `tfsdk:"inspect_app_id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	AppTypes            types.Set    `tfsdk:"app_types"`
	ApplicationPort     types.String `tfsdk:"application_port"`
	ApplicationProtocol types.String `tfsdk:"application_protocol"`
	CertificateID       types.String `tfsdk:"certificate_id"`
	Domain              types.String `tfsdk:"domain"`
	TrustUntrustedCert  types.Bool   `tfsdk:"trust_untrusted_cert"`
}

type inspectionAppModel struct {
	ID                  types.String `tfsdk:"id"`
	AppID               types.String `tfsdk:"app_id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	ApplicationPort     types.String `tfsdk:"application_port"`
	ApplicationProtocol types.String `tfsdk:"application_protocol"`
	CertificateID       types.String `tfsdk:"certificate_id"`
	CertificateName     types.String `tfsdk:"certificate_name"`
	Domain              types.String `tfsdk:"domain"`
	Protocols           types.List   `tfsdk:"protocols"`
	TrustUntrustedCert  types.Bool   `tfsdk:"trust_untrusted_cert"`
	MicroTenantID       types.String `tfsdk:"microtenant_id"`
	MicroTenantName     types.String `tfsdk:"microtenant_name"`
}

func (r *ApplicationSegmentInspectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_inspection"
}

func (r *ApplicationSegmentInspectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Inspection application segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the inspection application segment.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"segment_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Segment group identifier associated with the application.",
			},
			"segment_group_name": schema.StringAttribute{
				Computed: true,
			},
			"bypass_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NEVER"),
				Validators: []validator.String{
					stringvalidator.OneOf("ALWAYS", "NEVER", "ON_NET"),
				},
			},
			"bypass_on_reauth": schema.BoolAttribute{Optional: true, Computed: true},
			"config_space": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("DEFAULT", "SIEM"),
				},
			},
			"domain_names": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of domains and IPs for the application segment.",
			},
			"double_encrypt":         schema.BoolAttribute{Optional: true, Computed: true},
			"enabled":                schema.BoolAttribute{Optional: true, Computed: true},
			"passive_health_enabled": schema.BoolAttribute{Optional: true, Computed: true},
			"health_check_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"health_reporting": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"icmp_access_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
			},
			"ip_anchored":      schema.BoolAttribute{Optional: true, Computed: true},
			"is_cname_enabled": schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_keep_alive":   schema.StringAttribute{Optional: true, Computed: true},
			"use_in_dr_mode":   schema.BoolAttribute{Optional: true, Computed: true},
			"is_incomplete_dr_config": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"select_connector_close_to_app": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"adp_enabled":              schema.BoolAttribute{Optional: true, Computed: true},
			"auto_app_protect_enabled": schema.BoolAttribute{Optional: true, Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant identifier used to scope API requests.",
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"tcp_port_ranges": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"udp_port_ranges": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"tcp_protocols": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"udp_protocols": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"tcp_port_range": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Optional: true, Computed: true},
						"to":   schema.StringAttribute{Optional: true, Computed: true},
					},
				},
			},
			"udp_port_range": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{Optional: true, Computed: true},
						"to":   schema.StringAttribute{Optional: true, Computed: true},
					},
				},
			},
			"inspection_apps": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"app_id":               schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"enabled":              schema.BoolAttribute{Computed: true},
						"application_port":     schema.StringAttribute{Computed: true},
						"application_protocol": schema.StringAttribute{Computed: true},
						"certificate_id":       schema.StringAttribute{Computed: true},
						"certificate_name":     schema.StringAttribute{Computed: true},
						"domain":               schema.StringAttribute{Computed: true},
						"protocols": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"trust_untrusted_cert": schema.BoolAttribute{Computed: true},
						"microtenant_id":       schema.StringAttribute{Computed: true},
						"microtenant_name":     schema.StringAttribute{Computed: true},
					},
				},
			},
			// server_groups: TypeList in SDKv2, id is TypeSet (Required)
			"server_groups": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{ElementType: types.StringType, Required: true},
					},
				},
			},
			// common_apps_dto: TypeSet in SDKv2, apps_config is TypeSet
			"common_apps_dto": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"apps_config": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"app_id":         schema.StringAttribute{Optional: true, Computed: true},
									"inspect_app_id": schema.StringAttribute{Optional: true, Computed: true},
									"name":           schema.StringAttribute{Optional: true, Computed: true},
									"description":    schema.StringAttribute{Optional: true, Computed: true},
									"app_types": schema.SetAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
									},
									"application_port": schema.StringAttribute{Optional: true, Computed: true},
									"application_protocol": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("HTTP", "HTTPS"),
										},
									},
									"certificate_id":       schema.StringAttribute{Optional: true, Computed: true},
									"domain":               schema.StringAttribute{Optional: true, Computed: true},
									"trust_untrusted_cert": schema.BoolAttribute{Optional: true, Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ApplicationSegmentInspectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = cl
}

func (r *ApplicationSegmentInspectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan ApplicationSegmentInspectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if validation := validateInspectionCustomizeDiff(ctx, &plan); validation.HasError() {
		resp.Diagnostics.Append(validation...)
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload, diags := r.expandInspectionSegment(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if validation := helpers.ValidateAppPorts(helpers.BoolValue(plan.SelectConnectorCloseToApp, false), payload.UDPAppPortRange, payload.UDPPortRanges); validation.HasError() {
		resp.Diagnostics.Append(validation...)
		return
	}

	created, _, err := applicationsegmentinspection.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inspection application segment: %v", err))
		return
	}

	state, readDiags := r.readInspectionSegment(ctx, service, created.ID, plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentInspectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var state ApplicationSegmentInspectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)
	newState, diags := r.readInspectionSegment(ctx, service, state.ID.ValueString(), state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		if len(diags) > 0 {
			for _, diagnostic := range diags {
				if diagnostic.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(diagnostic.Detail()), "not found") {
					resp.State.RemoveResource(ctx)
					return
				}
			}
		}
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ApplicationSegmentInspectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var plan ApplicationSegmentInspectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if validation := validateInspectionCustomizeDiff(ctx, &plan); validation.HasError() {
		resp.Diagnostics.Append(validation...)
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload, diags := r.expandInspectionSegment(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if validation := helpers.ValidateAppPorts(helpers.BoolValue(plan.SelectConnectorCloseToApp, false), payload.UDPAppPortRange, payload.UDPPortRanges); validation.HasError() {
		resp.Diagnostics.Append(validation...)
		return
	}

	if _, err := applicationsegmentinspection.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inspection application segment: %v", err))
		return
	}

	state, readDiags := r.readInspectionSegment(ctx, service, plan.ID.ValueString(), plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentInspectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider was not configured before use.")
		return
	}

	var state ApplicationSegmentInspectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.SegmentGroupID.IsNull() && state.SegmentGroupID.ValueString() != "" {
		if err := helpers.DetachSegmentGroup(ctx, r.client, state.ID.ValueString(), state.SegmentGroupID.ValueString()); err != nil {
			resp.Diagnostics.AddError("Detach Segment Group Error", fmt.Sprintf("Failed to detach segment group: %v", err))
			return
		}
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)
	if _, err := applicationsegmentinspection.Delete(ctx, service, state.ID.ValueString()); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete inspection application segment: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ApplicationSegmentInspectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *ApplicationSegmentInspectionResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func (r *ApplicationSegmentInspectionResource) expandInspectionSegment(ctx context.Context, plan ApplicationSegmentInspectionModel) (applicationsegmentinspection.AppSegmentInspection, diag.Diagnostics) {
	var diags diag.Diagnostics

	domainNames, domainDiags := helpers.ListValueToStringSlice(ctx, plan.DomainNames)
	diags.Append(domainDiags...)

	tcpRanges, tcpRangeDiags := helpers.ListValueToStringSlice(ctx, plan.TCPPortRanges)
	diags.Append(tcpRangeDiags...)

	udpRanges, udpRangeDiags := helpers.ListValueToStringSlice(ctx, plan.UDPPortRanges)
	diags.Append(udpRangeDiags...)

	tcpPorts, tcpPortDiags := helpers.ExpandNetworkPorts(ctx, plan.TCPPortRange)
	diags.Append(tcpPortDiags...)

	udpPorts, udpPortDiags := helpers.ExpandNetworkPorts(ctx, plan.UDPPortRange)
	diags.Append(udpPortDiags...)

	tcpProtocols, tcpProtoDiags := helpers.ListValueToStringSlice(ctx, plan.TCPProtocols)
	diags.Append(tcpProtoDiags...)

	udpProtocols, udpProtoDiags := helpers.ListValueToStringSlice(ctx, plan.UDPProtocols)
	diags.Append(udpProtoDiags...)

	serverGroups, serverGroupDiags := helpers.ExpandServerGroups(ctx, plan.ServerGroups)
	diags.Append(serverGroupDiags...)

	commonApps, commonAppsDiags := expandInspectionCommonAppsDto(ctx, plan.CommonAppsDto)
	diags.Append(commonAppsDiags...)

	payload := applicationsegmentinspection.AppSegmentInspection{
		ID:                        plan.ID.ValueString(),
		Name:                      plan.Name.ValueString(),
		Description:               plan.Description.ValueString(),
		SegmentGroupID:            plan.SegmentGroupID.ValueString(),
		BypassType:                plan.BypassType.ValueString(),
		BypassOnReauth:            helpers.BoolValue(plan.BypassOnReauth, false),
		ConfigSpace:               plan.ConfigSpace.ValueString(),
		DomainNames:               domainNames,
		DoubleEncrypt:             helpers.BoolValue(plan.DoubleEncrypt, false),
		Enabled:                   helpers.BoolValue(plan.Enabled, false),
		PassiveHealthEnabled:      helpers.BoolValue(plan.PassiveHealthEnabled, false),
		HealthCheckType:           plan.HealthCheckType.ValueString(),
		HealthReporting:           plan.HealthReporting.ValueString(),
		ICMPAccessType:            plan.ICMPAccessType.ValueString(),
		IPAnchored:                helpers.BoolValue(plan.IPAnchored, false),
		IsCnameEnabled:            helpers.BoolValue(plan.IsCnameEnabled, false),
		TCPKeepAlive:              plan.TCPKeepAlive.ValueString(),
		SelectConnectorCloseToApp: helpers.BoolValue(plan.SelectConnectorCloseToApp, false),
		UseInDrMode:               helpers.BoolValue(plan.UseInDrMode, false),
		IsIncompleteDRConfig:      helpers.BoolValue(plan.IsIncompleteDRConfig, false),
		AdpEnabled:                helpers.BoolValue(plan.AdpEnabled, false),
		AutoAppProtectEnabled:     helpers.BoolValue(plan.AutoAppProtectEnabled, false),
		MicroTenantID:             plan.MicroTenantID.ValueString(),
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPAppPortRange:           tcpPorts,
		UDPAppPortRange:           udpPorts,
		TCPProtocols:              tcpProtocols,
		UDPProtocols:              udpProtocols,
		AppServerGroups:           serverGroups,
		CommonAppsDto:             commonApps,
	}

	return payload, diags
}

func (r *ApplicationSegmentInspectionResource) readInspectionSegment(ctx context.Context, service *zscaler.Service, id string, existingState ApplicationSegmentInspectionModel) (ApplicationSegmentInspectionModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state ApplicationSegmentInspectionModel

	segment, _, err := applicationsegmentinspection.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			diags.AddError("Not Found", fmt.Sprintf("Inspection application segment %s was not found", id))
			return state, diags
		}
		diags.AddError("Client Error", fmt.Sprintf("Unable to read inspection application segment: %v", err))
		return state, diags
	}

	domainNames, domainDiags := helpers.StringSliceToList(ctx, segment.DomainNames)
	diags.Append(domainDiags...)

	tcpRanges, tcpRangeDiags := helpers.StringSliceToList(ctx, segment.TCPPortRanges)
	diags.Append(tcpRangeDiags...)

	udpRanges, udpRangeDiags := helpers.StringSliceToList(ctx, segment.UDPPortRanges)
	diags.Append(udpRangeDiags...)

	tcpPorts, tcpPortDiags := helpers.FlattenNetworkPorts(ctx, segment.TCPAppPortRange)
	diags.Append(tcpPortDiags...)

	udpPorts, udpPortDiags := helpers.FlattenNetworkPorts(ctx, segment.UDPAppPortRange)
	diags.Append(udpPortDiags...)

	tcpProtocols, tcpProtoDiags := helpers.StringSliceToList(ctx, segment.TCPProtocols)
	diags.Append(tcpProtoDiags...)

	udpProtocols, udpProtoDiags := helpers.StringSliceToList(ctx, segment.UDPProtocols)
	diags.Append(udpProtoDiags...)

	serverGroups, serverGroupDiags := helpers.FlattenServerGroups(ctx, segment.AppServerGroups)
	diags.Append(serverGroupDiags...)

	inspectionApps, inspectionAppsDiags := flattenInspectionApps(ctx, segment.InspectionAppDto)
	diags.Append(inspectionAppsDiags...)

	// Preserve common_apps_dto from plan/state and update app_id/inspect_app_id from InspectionAppDto (matching SDKv2 setInspectionAppIDsInCommonAppsDto)
	commonApps, commonDiags := setInspectionAppIDsInCommonAppsDto(ctx, existingState.CommonAppsDto, segment.InspectionAppDto)
	diags.Append(commonDiags...)

	state = ApplicationSegmentInspectionModel{
		ID:                        helpers.StringValueOrNull(segment.ID),
		Name:                      helpers.StringValueOrNull(segment.Name),
		Description:               helpers.StringValueOrNull(segment.Description),
		SegmentGroupID:            helpers.StringValueOrNull(segment.SegmentGroupID),
		SegmentGroupName:          helpers.StringValueOrNull(segment.SegmentGroupName),
		BypassType:                helpers.StringValueOrNull(segment.BypassType),
		BypassOnReauth:            types.BoolValue(segment.BypassOnReauth),
		ConfigSpace:               helpers.StringValueOrNull(segment.ConfigSpace),
		DomainNames:               domainNames,
		DoubleEncrypt:             types.BoolValue(segment.DoubleEncrypt),
		Enabled:                   types.BoolValue(segment.Enabled),
		PassiveHealthEnabled:      types.BoolValue(segment.PassiveHealthEnabled),
		HealthCheckType:           helpers.StringValueOrNull(segment.HealthCheckType),
		HealthReporting:           helpers.StringValueOrNull(segment.HealthReporting),
		ICMPAccessType:            helpers.StringValueOrNull(segment.ICMPAccessType),
		IPAnchored:                types.BoolValue(segment.IPAnchored),
		IsCnameEnabled:            types.BoolValue(segment.IsCnameEnabled),
		TCPKeepAlive:              helpers.StringValueOrNull(segment.TCPKeepAlive),
		SelectConnectorCloseToApp: types.BoolValue(segment.SelectConnectorCloseToApp),
		UseInDrMode:               types.BoolValue(segment.UseInDrMode),
		IsIncompleteDRConfig:      types.BoolValue(segment.IsIncompleteDRConfig),
		AdpEnabled:                types.BoolValue(segment.AdpEnabled),
		AutoAppProtectEnabled:     types.BoolValue(segment.AutoAppProtectEnabled),
		MicroTenantID:             helpers.StringValueOrNull(segment.MicroTenantID),
		MicroTenantName:           helpers.StringValueOrNull(segment.MicroTenantName),
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPPortRange:              tcpPorts,
		UDPPortRange:              udpPorts,
		TCPProtocols:              tcpProtocols,
		UDPProtocols:              udpProtocols,
		ServerGroups:              serverGroups,
		CommonAppsDto:             commonApps,
		InspectionApps:            inspectionApps,
	}

	return state, diags
}

func expandInspectionCommonAppsDto(ctx context.Context, value types.List) (applicationsegmentinspection.CommonAppsDto, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := applicationsegmentinspection.CommonAppsDto{}
	if value.IsNull() || value.IsUnknown() {
		return result, diags
	}

	var models []inspectionCommonAppsModel
	diags.Append(value.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return result, diags
	}

	apps, appsDiags := expandInspectionAppsConfig(ctx, models[0].AppsConfig)
	diags.Append(appsDiags...)
	if diags.HasError() {
		return result, diags
	}

	result.AppsConfig = apps
	return result, diags
}

func expandInspectionAppsConfig(ctx context.Context, list types.List) ([]applicationsegmentinspection.AppsConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var models []inspectionAppsConfigModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}

	result := make([]applicationsegmentinspection.AppsConfig, 0, len(models))
	for _, model := range models {
		appTypes, appTypeDiags := helpers.SetValueToStringSlice(ctx, model.AppTypes)
		diags.Append(appTypeDiags...)
		if diags.HasError() {
			return nil, diags
		}

		name := model.Name.ValueString()
		domain := model.Domain.ValueString()
		if strings.TrimSpace(name) == "" && strings.TrimSpace(domain) != "" {
			name = domain
		}

		result = append(result, applicationsegmentinspection.AppsConfig{
			AppID:               model.AppID.ValueString(),
			InspectAppID:        model.InspectAppID.ValueString(),
			Name:                name,
			Description:         model.Description.ValueString(),
			AppTypes:            appTypes,
			ApplicationPort:     model.ApplicationPort.ValueString(),
			ApplicationProtocol: model.ApplicationProtocol.ValueString(),
			CertificateID:       model.CertificateID.ValueString(),
			Domain:              domain,
			TrustUntrustedCert:  helpers.BoolValue(model.TrustUntrustedCert, false),
		})
	}

	return result, diags
}

func flattenInspectionApps(ctx context.Context, apps []applicationsegmentinspection.InspectionAppDto) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":                   types.StringType,
		"app_id":               types.StringType,
		"name":                 types.StringType,
		"description":          types.StringType,
		"enabled":              types.BoolType,
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"certificate_id":       types.StringType,
		"certificate_name":     types.StringType,
		"domain":               types.StringType,
		"protocols":            types.ListType{ElemType: types.StringType},
		"trust_untrusted_cert": types.BoolType,
		"microtenant_id":       types.StringType,
		"microtenant_name":     types.StringType,
	}

	if len(apps) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(apps))
	var diags diag.Diagnostics
	for _, app := range apps {
		protocols, protoDiags := helpers.StringSliceToList(ctx, app.Protocols)
		diags.Append(protoDiags...)

		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                   helpers.StringValueOrNull(app.ID),
			"app_id":               helpers.StringValueOrNull(app.AppID),
			"name":                 helpers.StringValueOrNull(app.Name),
			"description":          helpers.StringValueOrNull(app.Description),
			"enabled":              types.BoolValue(app.Enabled),
			"application_port":     helpers.StringValueOrNull(app.ApplicationPort),
			"application_protocol": helpers.StringValueOrNull(app.ApplicationProtocol),
			"certificate_id":       helpers.StringValueOrNull(app.CertificateID),
			"certificate_name":     helpers.StringValueOrNull(app.CertificateName),
			"domain":               helpers.StringValueOrNull(app.Domain),
			"protocols":            protocols,
			"trust_untrusted_cert": types.BoolValue(app.TrustUntrustedCert),
			"microtenant_id":       helpers.StringValueOrNull(app.MicroTenantID),
			"microtenant_name":     helpers.StringValueOrNull(app.MicroTenantName),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

// setInspectionAppIDsInCommonAppsDto updates app_id and inspect_app_id in common_apps_dto from InspectionAppDto response
// This matches SDKv2's setInspectionAppIDsInCommonAppsDto function
func setInspectionAppIDsInCommonAppsDto(ctx context.Context, commonAppsDto types.List, inspectionApps []applicationsegmentinspection.InspectionAppDto) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if commonAppsDto.IsNull() || commonAppsDto.IsUnknown() {
		return commonAppsDto, diags
	}

	if len(inspectionApps) == 0 {
		return commonAppsDto, diags
	}

	// Extract app_id, inspect_app_id, name, and description from the first Inspection app (matching SDKv2 logic)
	appID := inspectionApps[0].AppID
	inspectAppID := inspectionApps[0].ID
	apiName := inspectionApps[0].Name
	apiDescription := inspectionApps[0].Description

	// Expand the common_apps_dto to get the structure
	type inspectionCommonAppsModel struct {
		AppsConfig types.List `tfsdk:"apps_config"`
	}

	type inspectionAppsConfigModel struct {
		AppID               types.String `tfsdk:"app_id"`
		InspectAppID        types.String `tfsdk:"inspect_app_id"`
		Name                types.String `tfsdk:"name"`
		Description         types.String `tfsdk:"description"`
		AppTypes            types.Set    `tfsdk:"app_types"`
		ApplicationPort     types.String `tfsdk:"application_port"`
		ApplicationProtocol types.String `tfsdk:"application_protocol"`
		CertificateID       types.String `tfsdk:"certificate_id"`
		Domain              types.String `tfsdk:"domain"`
		TrustUntrustedCert  types.Bool   `tfsdk:"trust_untrusted_cert"`
	}

	var containers []inspectionCommonAppsModel
	diags.Append(commonAppsDto.ElementsAs(ctx, &containers, false)...)
	if diags.HasError() || len(containers) == 0 {
		return commonAppsDto, diags
	}

	// Update app_id and inspect_app_id for each apps_config entry
	var updatedAppConfigs []attr.Value
	container := containers[0]
	if !container.AppsConfig.IsNull() && !container.AppsConfig.IsUnknown() {
		var items []inspectionAppsConfigModel
		diags.Append(container.AppsConfig.ElementsAs(ctx, &items, false)...)
		if diags.HasError() {
			return commonAppsDto, diags
		}

		appAttrTypes := map[string]attr.Type{
			"app_id":               types.StringType,
			"inspect_app_id":       types.StringType,
			"name":                 types.StringType,
			"description":          types.StringType,
			"app_types":            types.SetType{ElemType: types.StringType},
			"application_port":     types.StringType,
			"application_protocol": types.StringType,
			"certificate_id":       types.StringType,
			"domain":               types.StringType,
			"trust_untrusted_cert": types.BoolType,
		}

		// Update the first entry with app_id, inspect_app_id, name, and description (matching SDKv2 logic)
		for i, item := range items {
			updatedAppID := item.AppID
			updatedInspectAppID := item.InspectAppID
			updatedName := item.Name
			updatedDescription := item.Description

			// Only update the first entry (matching SDKv2 setInspectionAppIDsInCommonAppsDto)
			if i == 0 {
				updatedAppID = types.StringValue(appID)
				updatedInspectAppID = types.StringValue(inspectAppID)
				// Populate name and description from API if they're null/unknown in plan
				if updatedName.IsNull() || updatedName.IsUnknown() {
					updatedName = helpers.StringValueOrNull(apiName)
				}
				if updatedDescription.IsNull() || updatedDescription.IsUnknown() {
					updatedDescription = helpers.StringValueOrNull(apiDescription)
				}
			}

			obj, objDiags := types.ObjectValue(appAttrTypes, map[string]attr.Value{
				"app_id":               updatedAppID,
				"inspect_app_id":       updatedInspectAppID,
				"name":                 updatedName,
				"description":          updatedDescription,
				"app_types":            item.AppTypes,
				"application_port":     item.ApplicationPort,
				"application_protocol": item.ApplicationProtocol,
				"certificate_id":       item.CertificateID,
				"domain":               item.Domain,
				"trust_untrusted_cert": item.TrustUntrustedCert,
			})
			diags.Append(objDiags...)
			if diags.HasError() {
				return commonAppsDto, diags
			}
			updatedAppConfigs = append(updatedAppConfigs, obj)
		}
	}

	// Rebuild the common_apps_dto structure
	appsConfigList, listDiags := types.ListValue(types.ObjectType{AttrTypes: map[string]attr.Type{
		"app_id":               types.StringType,
		"inspect_app_id":       types.StringType,
		"name":                 types.StringType,
		"description":          types.StringType,
		"app_types":            types.SetType{ElemType: types.StringType},
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"certificate_id":       types.StringType,
		"domain":               types.StringType,
		"trust_untrusted_cert": types.BoolType,
	}}, updatedAppConfigs)
	diags.Append(listDiags...)
	if diags.HasError() {
		return commonAppsDto, diags
	}

	commonAttrTypes := map[string]attr.Type{
		"apps_config": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"app_id":               types.StringType,
			"inspect_app_id":       types.StringType,
			"name":                 types.StringType,
			"description":          types.StringType,
			"app_types":            types.SetType{ElemType: types.StringType},
			"application_port":     types.StringType,
			"application_protocol": types.StringType,
			"certificate_id":       types.StringType,
			"domain":               types.StringType,
			"trust_untrusted_cert": types.BoolType,
		}}},
	}

	commonObj, objDiags := types.ObjectValue(commonAttrTypes, map[string]attr.Value{
		"apps_config": appsConfigList,
	})
	diags.Append(objDiags...)
	if diags.HasError() {
		return commonAppsDto, diags
	}

	updatedList, listDiags := types.ListValue(types.ObjectType{AttrTypes: commonAttrTypes}, []attr.Value{commonObj})
	diags.Append(listDiags...)
	return updatedList, diags
}

func validateInspectionCustomizeDiff(ctx context.Context, plan *ApplicationSegmentInspectionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if helpers.BoolValue(plan.AdpEnabled, false) && helpers.BoolValue(plan.AutoAppProtectEnabled, false) {
		diags.AddError("Invalid configuration", "If 'adp_enabled' is set to true, 'auto_app_protect_enabled' must be false.")
	}

	if plan.CommonAppsDto.IsNull() || plan.CommonAppsDto.IsUnknown() {
		return diags
	}

	var entries []inspectionCommonAppsModel
	diags.Append(plan.CommonAppsDto.ElementsAs(ctx, &entries, false)...)
	if diags.HasError() || len(entries) == 0 {
		return diags
	}

	for _, entry := range entries {
		if entry.AppsConfig.IsNull() || entry.AppsConfig.IsUnknown() {
			continue
		}
		var configs []inspectionAppsConfigModel
		diags.Append(entry.AppsConfig.ElementsAs(ctx, &configs, false)...)
		if diags.HasError() || len(configs) == 0 {
			return diags
		}
		for i, cfg := range configs {
			if strings.EqualFold(cfg.ApplicationProtocol.ValueString(), "HTTP") && !cfg.CertificateID.IsNull() && cfg.CertificateID.ValueString() != "" {
				diags.AddError("Invalid common_apps_dto configuration", fmt.Sprintf("common_apps_dto.apps_config[%d]: certificate_id must not be set when application_protocol is HTTP", i))
			}
		}
	}

	return diags
}
