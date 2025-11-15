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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

var (
	_ resource.Resource                = &ApplicationSegmentPRAResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSegmentPRAResource{}
	_ resource.ResourceWithImportState = &ApplicationSegmentPRAResource{}
)

func NewApplicationSegmentPRAResource() resource.Resource {
	return &ApplicationSegmentPRAResource{}
}

type ApplicationSegmentPRAResource struct {
	client *client.Client
}

type ApplicationSegmentPRAModel struct {
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
	IcmpAccessType            types.String `tfsdk:"icmp_access_type"`
	IPAnchored                types.Bool   `tfsdk:"ip_anchored"`
	FqdnDnsCheck              types.Bool   `tfsdk:"fqdn_dns_check"`
	SelectConnectorCloseToApp types.Bool   `tfsdk:"select_connector_close_to_app"`
	UseInDrMode               types.Bool   `tfsdk:"use_in_dr_mode"`
	IsIncompleteDRConfig      types.Bool   `tfsdk:"is_incomplete_dr_config"`
	IsCnameEnabled            types.Bool   `tfsdk:"is_cname_enabled"`
	TCPKeepAlive              types.String `tfsdk:"tcp_keep_alive"`
	MicroTenantID             types.String `tfsdk:"microtenant_id"`
	MicroTenantName           types.String `tfsdk:"microtenant_name"`
	TCPPortRanges             types.List   `tfsdk:"tcp_port_ranges"`
	UDPPortRanges             types.List   `tfsdk:"udp_port_ranges"`
	TCPPortRange              types.List   `tfsdk:"tcp_port_range"`
	UDPPortRange              types.List   `tfsdk:"udp_port_range"`
	ServerGroups              types.List   `tfsdk:"server_groups"`
	ZpnERID                   types.List   `tfsdk:"zpn_er_id"`
	CommonAppsDto             types.List   `tfsdk:"common_apps_dto"`
	PraApps                   types.List   `tfsdk:"pra_apps"`
}

func (r *ApplicationSegmentPRAResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_pra"
}

func (r *ApplicationSegmentPRAResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Privileged Remote Access (PRA) application segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the PRA application segment.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"segment_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Segment group identifier associated with the application.",
			},
			"segment_group_name": schema.StringAttribute{Computed: true},
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
				Description: "List of domains and IPs for the PRA application segment.",
			},
			"double_encrypt":         schema.BoolAttribute{Optional: true, Computed: true},
			"enabled":                schema.BoolAttribute{Optional: true, Computed: true},
			"passive_health_enabled": schema.BoolAttribute{Optional: true, Computed: true},
			"health_check_type":      schema.StringAttribute{Optional: true, Computed: true},
			"health_reporting":       schema.StringAttribute{Optional: true, Computed: true},
			"icmp_access_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("PING_TRACEROUTING", "PING", "NONE"),
				},
			},
			"ip_anchored":                   schema.BoolAttribute{Optional: true, Computed: true},
			"fqdn_dns_check":                schema.BoolAttribute{Optional: true, Computed: true},
			"select_connector_close_to_app": schema.BoolAttribute{Optional: true, Computed: true},
			"use_in_dr_mode":                schema.BoolAttribute{Optional: true, Computed: true},
			"is_incomplete_dr_config":       schema.BoolAttribute{Optional: true, Computed: true},
			"is_cname_enabled":              schema.BoolAttribute{Optional: true, Computed: true},
			"tcp_keep_alive":                schema.StringAttribute{Optional: true, Computed: true},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant identifier used to scope API calls.",
			},
			"microtenant_name": schema.StringAttribute{Computed: true},
			"tcp_port_ranges":  schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"udp_port_ranges":  schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
		},
		Blocks: map[string]schema.Block{
			"zpn_er_id": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
					},
				},
			},
			"pra_apps": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"app_id":               schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"application_port":     schema.StringAttribute{Computed: true},
						"application_protocol": schema.StringAttribute{Computed: true},
						"certificate_id":       schema.StringAttribute{Computed: true},
						"certificate_name":     schema.StringAttribute{Computed: true},
						"connection_security":  schema.StringAttribute{Computed: true},
						"domain":               schema.StringAttribute{Computed: true},
						"enabled":              schema.BoolAttribute{Computed: true},
						"hidden":               schema.BoolAttribute{Computed: true},
						"portal":               schema.BoolAttribute{Computed: true},
						"microtenant_id":       schema.StringAttribute{Computed: true},
						"microtenant_name":     schema.StringAttribute{Computed: true},
					},
				},
			},
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
			// server_groups: TypeList in SDKv2, id is TypeSet (Required)
			"server_groups": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.SetAttribute{ElementType: types.StringType, Required: true},
					},
				},
			},
			// common_apps_dto: TypeList in SDKv2, apps_config is TypeList
			"common_apps_dto": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"apps_config": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"app_id":           schema.StringAttribute{Optional: true, Computed: true},
									"pra_app_id":       schema.StringAttribute{Optional: true, Computed: true},
									"name":             schema.StringAttribute{Optional: true},
									"description":      schema.StringAttribute{Optional: true},
									"app_types":        schema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
									"application_port": schema.StringAttribute{Optional: true, Computed: true},
									"application_protocol": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("RDP", "SSH", "VNC"),
										},
									},
									"connection_security": schema.StringAttribute{Optional: true},
									"domain":              schema.StringAttribute{Required: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ApplicationSegmentPRAResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationSegmentPRAResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ApplicationSegmentPRAModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diags := validatePRAPlan(ctx, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	payload, diags := r.expandApplicationSegment(ctx, plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if portDiags := helpers.ValidateAppPorts(payload.SelectConnectorCloseToApp, payload.UDPAppPortRange, payload.UDPPortRanges); portDiags.HasError() {
		resp.Diagnostics.Append(portDiags...)
		return
	}

	created, _, err := applicationsegmentpra.Create(ctx, service, payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create PRA application segment: %v", err))
		return
	}

	plan.ID = types.StringValue(created.ID)

	state, readDiags := r.readIntoState(ctx, service, created.ID, plan.MicroTenantID, plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentPRAResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ApplicationSegmentPRAModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	newState, diags := r.readIntoState(ctx, service, state.ID.ValueString(), state.MicroTenantID, state)
	if diags.HasError() {
		for _, d := range diags {
			if d.Severity() == diag.SeverityError && strings.Contains(strings.ToLower(d.Detail()), "not found") {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ApplicationSegmentPRAResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan ApplicationSegmentPRAModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diags := validatePRAPlan(ctx, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	existing, _, err := applicationsegmentpra.Get(ctx, service, plan.ID.ValueString())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve PRA application segment: %v", err))
		return
	}

	payload, diags := r.expandApplicationSegment(ctx, plan, existing)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if portDiags := helpers.ValidateAppPorts(payload.SelectConnectorCloseToApp, payload.UDPAppPortRange, payload.UDPPortRanges); portDiags.HasError() {
		resp.Diagnostics.Append(portDiags...)
		return
	}

	if _, err := applicationsegmentpra.Update(ctx, service, plan.ID.ValueString(), &payload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update PRA application segment: %v", err))
		return
	}

	state, readDiags := r.readIntoState(ctx, service, plan.ID.ValueString(), plan.MicroTenantID, plan)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationSegmentPRAResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state ApplicationSegmentPRAModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	if !state.SegmentGroupID.IsNull() && state.SegmentGroupID.ValueString() != "" {
		if err := helpers.DetachSegmentGroup(ctx, r.client, state.ID.ValueString(), state.SegmentGroupID.ValueString()); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to detach segment group: %v", err))
			return
		}
	}

	if _, err := applicationsegmentpra.Delete(ctx, service, state.ID.ValueString()); err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete PRA application segment: %v", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ApplicationSegmentPRAResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

func (r *ApplicationSegmentPRAResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && !microtenantID.IsUnknown() {
		trimmed := strings.TrimSpace(microtenantID.ValueString())
		if trimmed != "" {
			service = service.WithMicroTenant(trimmed)
		}
	}
	return service
}

func (r *ApplicationSegmentPRAResource) expandApplicationSegment(ctx context.Context, plan ApplicationSegmentPRAModel, existing *applicationsegmentpra.AppSegmentPRA) (applicationsegmentpra.AppSegmentPRA, diag.Diagnostics) {
	var diags diag.Diagnostics

	domainNames, domainDiags := helpers.ListValueToStringSlice(ctx, plan.DomainNames)
	diags.Append(domainDiags...)

	tcpRanges, tcpRangesDiags := helpers.ListValueToStringSlice(ctx, plan.TCPPortRanges)
	diags.Append(tcpRangesDiags...)

	udpRanges, udpRangesDiags := helpers.ListValueToStringSlice(ctx, plan.UDPPortRanges)
	diags.Append(udpRangesDiags...)

	tcpPorts, tcpDiags := helpers.ExpandNetworkPorts(ctx, plan.TCPPortRange)
	diags.Append(tcpDiags...)

	udpPorts, udpDiags := helpers.ExpandNetworkPorts(ctx, plan.UDPPortRange)
	diags.Append(udpDiags...)

	serverGroups, sgDiags := helpers.ExpandServerGroups(ctx, plan.ServerGroups)
	diags.Append(sgDiags...)

	zpnER, zpnDiags := helpers.ExpandZPNERID(ctx, plan.ZpnERID)
	diags.Append(zpnDiags...)

	if diags.HasError() {
		return applicationsegmentpra.AppSegmentPRA{}, diags
	}

	commonApps, filteredDomains, commonDiags := r.buildCommonApps(ctx, plan, domainNames, existing)
	diags.Append(commonDiags...)
	if diags.HasError() {
		return applicationsegmentpra.AppSegmentPRA{}, diags
	}

	payload := applicationsegmentpra.AppSegmentPRA{
		ID:                        strings.TrimSpace(plan.ID.ValueString()),
		Name:                      strings.TrimSpace(plan.Name.ValueString()),
		Description:               strings.TrimSpace(plan.Description.ValueString()),
		SegmentGroupID:            strings.TrimSpace(plan.SegmentGroupID.ValueString()),
		SegmentGroupName:          strings.TrimSpace(plan.SegmentGroupName.ValueString()),
		BypassType:                strings.TrimSpace(plan.BypassType.ValueString()),
		BypassOnReauth:            helpers.BoolValue(plan.BypassOnReauth, false),
		ConfigSpace:               strings.TrimSpace(plan.ConfigSpace.ValueString()),
		DomainNames:               filteredDomains,
		DoubleEncrypt:             helpers.BoolValue(plan.DoubleEncrypt, false),
		Enabled:                   helpers.BoolValue(plan.Enabled, false),
		PassiveHealthEnabled:      helpers.BoolValue(plan.PassiveHealthEnabled, false),
		HealthCheckType:           strings.TrimSpace(plan.HealthCheckType.ValueString()),
		HealthReporting:           strings.TrimSpace(plan.HealthReporting.ValueString()),
		IcmpAccessType:            strings.TrimSpace(plan.IcmpAccessType.ValueString()),
		IpAnchored:                helpers.BoolValue(plan.IPAnchored, false),
		FQDNDnsCheck:              helpers.BoolValue(plan.FqdnDnsCheck, false),
		SelectConnectorCloseToApp: helpers.BoolValue(plan.SelectConnectorCloseToApp, false),
		UseInDrMode:               helpers.BoolValue(plan.UseInDrMode, false),
		IsIncompleteDRConfig:      helpers.BoolValue(plan.IsIncompleteDRConfig, false),
		IsCnameEnabled:            helpers.BoolValue(plan.IsCnameEnabled, false),
		TCPKeepAlive:              strings.TrimSpace(plan.TCPKeepAlive.ValueString()),
		MicroTenantID:             strings.TrimSpace(plan.MicroTenantID.ValueString()),
		ServerGroups:              serverGroups,
		ZPNERID:                   zpnER,
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPAppPortRange:           tcpPorts,
		UDPAppPortRange:           udpPorts,
		CommonAppsDto:             commonApps,
	}

	return payload, diags
}

func (r *ApplicationSegmentPRAResource) buildCommonApps(ctx context.Context, plan ApplicationSegmentPRAModel, domainNames []string, existing *applicationsegmentpra.AppSegmentPRA) (applicationsegmentpra.CommonAppsDto, []string, diag.Diagnostics) {
	var diags diag.Diagnostics

	planConfigs, configsDiags := helpers.ExpandPRACommonApps(ctx, plan.CommonAppsDto)
	diags.Append(configsDiags...)
	if diags.HasError() {
		return applicationsegmentpra.CommonAppsDto{}, domainNames, diags
	}

	// Normalize domain names from plan configs
	planDomainMap := make(map[string]applicationsegmentpra.AppsConfig)
	orderedDomains := make([]string, 0, len(planConfigs))
	for _, cfg := range planConfigs {
		domain := strings.TrimSpace(cfg.Domain)
		if domain == "" {
			continue
		}
		cfg.Domain = domain
		cfg.Name = strings.TrimSpace(cfg.Name)
		cfg.Description = strings.TrimSpace(cfg.Description)
		cfg.ApplicationPort = strings.TrimSpace(cfg.ApplicationPort)
		cfg.ApplicationProtocol = strings.TrimSpace(cfg.ApplicationProtocol)
		cfg.ConnectionSecurity = strings.TrimSpace(cfg.ConnectionSecurity)
		planDomainMap[strings.ToLower(domain)] = cfg
		orderedDomains = append(orderedDomains, domain)
	}

	existingByDomain := make(map[string]applicationsegmentpra.PRAApps)
	if existing != nil {
		for _, pra := range existing.PRAApps {
			existingByDomain[strings.ToLower(pra.Domain)] = pra
		}
	}

	resultConfigs := make([]applicationsegmentpra.AppsConfig, 0, len(planDomainMap))
	for _, domain := range orderedDomains {
		cfg := planDomainMap[strings.ToLower(domain)]
		if len(cfg.AppTypes) == 0 {
			cfg.AppTypes = []string{"SECURE_REMOTE_ACCESS"}
		}

		if existingApp, ok := existingByDomain[strings.ToLower(domain)]; ok {
			if strings.TrimSpace(cfg.AppID) == "" {
				cfg.AppID = existingApp.AppID
			}
			if strings.TrimSpace(cfg.PRAAppID) == "" {
				cfg.PRAAppID = existingApp.ID
			}
			if strings.TrimSpace(cfg.Name) == "" {
				cfg.Name = existingApp.Name
			}
			if strings.TrimSpace(cfg.Description) == "" {
				cfg.Description = existingApp.Description
			}
			if strings.TrimSpace(cfg.ApplicationPort) == "" {
				cfg.ApplicationPort = existingApp.ApplicationPort
			}
			if strings.TrimSpace(cfg.ApplicationProtocol) == "" {
				cfg.ApplicationProtocol = existingApp.ApplicationProtocol
			}
			if strings.TrimSpace(cfg.ConnectionSecurity) == "" {
				cfg.ConnectionSecurity = existingApp.ConnectionSecurity
			}
		}

		resultConfigs = append(resultConfigs, cfg)
	}

	deletedPraIDs := make([]string, 0)
	deletedDomains := make(map[string]struct{})
	if existing != nil {
		for _, pra := range existing.PRAApps {
			if _, ok := planDomainMap[strings.ToLower(pra.Domain)]; !ok {
				deletedPraIDs = append(deletedPraIDs, pra.ID)
				deletedDomains[strings.ToLower(pra.Domain)] = struct{}{}
			}
		}
	}

	if len(deletedDomains) > 0 {
		filtered := make([]string, 0, len(domainNames))
		for _, domain := range domainNames {
			if _, isDeleted := deletedDomains[strings.ToLower(strings.TrimSpace(domain))]; isDeleted {
				continue
			}
			filtered = append(filtered, domain)
		}
		domainNames = filtered
	}

	dto := applicationsegmentpra.CommonAppsDto{
		AppsConfig:     resultConfigs,
		DeletedPraApps: deletedPraIDs,
	}

	return dto, domainNames, diags
}

func (r *ApplicationSegmentPRAResource) readIntoState(ctx context.Context, service *zscaler.Service, id string, microTenantID types.String, existingState ApplicationSegmentPRAModel) (ApplicationSegmentPRAModel, diag.Diagnostics) {
	segment, _, err := applicationsegmentpra.Get(ctx, service, id)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			return ApplicationSegmentPRAModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("PRA application segment %s not found", id))}
		}
		return ApplicationSegmentPRAModel{}, diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read PRA application segment: %v", err))}
	}

	// Preserve domain_names order from plan/state (SDKv2 sets it from API but we preserve order to match plan)
	domainNames := existingState.DomainNames
	if domainNames.IsNull() || domainNames.IsUnknown() {
		var domainDiags diag.Diagnostics
		domainNames, domainDiags = helpers.StringSliceToList(ctx, segment.DomainNames)
		if domainDiags.HasError() {
			return ApplicationSegmentPRAModel{}, domainDiags
		}
	}

	tcpRanges, tcpDiags := helpers.StringSliceToList(ctx, segment.TCPPortRanges)
	udpRanges, udpDiags := helpers.StringSliceToList(ctx, segment.UDPPortRanges)
	tcpPorts, tcpPortsDiags := helpers.FlattenNetworkPorts(ctx, segment.TCPAppPortRange)
	udpPorts, udpPortsDiags := helpers.FlattenNetworkPorts(ctx, segment.UDPAppPortRange)
	serverGroups, sgDiags := helpers.FlattenServerGroups(ctx, segment.ServerGroups)
	zpnER, zpnDiags := helpers.FlattenZPNERID(ctx, segment.ZPNERID)
	praApps, praDiags := helpers.FlattenPRAApps(ctx, segment.PRAApps)
	// Preserve common_apps_dto from plan/state and update app_id/pra_app_id from PRA apps (matching SDKv2 setAppIDsInCommonAppsDto)
	commonApps, commonDiags := setAppIDsInCommonAppsDto(ctx, existingState.CommonAppsDto, segment.PRAApps)

	diags := diag.Diagnostics{}
	diags.Append(tcpDiags...)
	diags.Append(commonDiags...)
	diags.Append(udpDiags...)
	diags.Append(tcpPortsDiags...)
	diags.Append(udpPortsDiags...)
	diags.Append(sgDiags...)
	diags.Append(zpnDiags...)
	diags.Append(praDiags...)
	if diags.HasError() {
		return ApplicationSegmentPRAModel{}, diags
	}

	state := ApplicationSegmentPRAModel{
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
		IcmpAccessType:            helpers.StringValueOrNull(segment.IcmpAccessType),
		IPAnchored:                types.BoolValue(segment.IpAnchored),
		FqdnDnsCheck:              types.BoolValue(segment.FQDNDnsCheck),
		SelectConnectorCloseToApp: types.BoolValue(segment.SelectConnectorCloseToApp),
		UseInDrMode:               types.BoolValue(segment.UseInDrMode),
		IsIncompleteDRConfig:      types.BoolValue(segment.IsIncompleteDRConfig),
		IsCnameEnabled:            types.BoolValue(segment.IsCnameEnabled),
		TCPKeepAlive:              helpers.StringValueOrNull(segment.TCPKeepAlive),
		MicroTenantID:             helpers.StringValueOrNull(segment.MicroTenantID),
		MicroTenantName:           helpers.StringValueOrNull(segment.MicroTenantName),
		TCPPortRanges:             tcpRanges,
		UDPPortRanges:             udpRanges,
		TCPPortRange:              tcpPorts,
		UDPPortRange:              udpPorts,
		ServerGroups:              serverGroups,
		ZpnERID:                   zpnER,
		CommonAppsDto:             commonApps,
		PraApps:                   praApps,
	}

	return state, diags
}

// setAppIDsInCommonAppsDto updates app_id and pra_app_id in common_apps_dto from PRA apps response
// This matches SDKv2's setAppIDsInCommonAppsDto function
func setAppIDsInCommonAppsDto(ctx context.Context, commonAppsDto types.List, praApps []applicationsegmentpra.PRAApps) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if commonAppsDto.IsNull() || commonAppsDto.IsUnknown() {
		return commonAppsDto, diags
	}

	if len(praApps) == 0 {
		return commonAppsDto, diags
	}

	// Build a map of existing apps by domain
	existingMap := make(map[string]applicationsegmentpra.PRAApps)
	for _, app := range praApps {
		existingMap[strings.ToLower(strings.TrimSpace(app.Domain))] = app
	}

	// Expand the common_apps_dto to get the structure
	type praCommonAppsModel struct {
		AppsConfig types.List `tfsdk:"apps_config"`
	}

	type praAppConfigModel struct {
		AppID               types.String `tfsdk:"app_id"`
		PRAAppID            types.String `tfsdk:"pra_app_id"`
		Name                types.String `tfsdk:"name"`
		Description         types.String `tfsdk:"description"`
		AppTypes            types.Set    `tfsdk:"app_types"`
		ApplicationPort     types.String `tfsdk:"application_port"`
		ApplicationProtocol types.String `tfsdk:"application_protocol"`
		ConnectionSecurity  types.String `tfsdk:"connection_security"`
		Domain              types.String `tfsdk:"domain"`
	}

	var containers []praCommonAppsModel
	diags.Append(commonAppsDto.ElementsAs(ctx, &containers, false)...)
	if diags.HasError() || len(containers) == 0 {
		return commonAppsDto, diags
	}

	// Update app_id and pra_app_id for each apps_config entry
	var updatedAppConfigs []attr.Value
	container := containers[0]
	if !container.AppsConfig.IsNull() && !container.AppsConfig.IsUnknown() {
		var items []praAppConfigModel
		diags.Append(container.AppsConfig.ElementsAs(ctx, &items, false)...)
		if diags.HasError() {
			return commonAppsDto, diags
		}

		appAttrTypes := map[string]attr.Type{
			"app_id":               types.StringType,
			"pra_app_id":           types.StringType,
			"name":                 types.StringType,
			"description":          types.StringType,
			"app_types":            types.SetType{ElemType: types.StringType},
			"application_port":     types.StringType,
			"application_protocol": types.StringType,
			"connection_security":  types.StringType,
			"domain":               types.StringType,
		}

		for _, item := range items {
			domain := strings.ToLower(strings.TrimSpace(item.Domain.ValueString()))
			name := strings.TrimSpace(item.Name.ValueString())

			appID := item.AppID
			praAppID := item.PRAAppID

			// Match by domain and name (matching SDKv2 logic)
			if existingApp, ok := existingMap[domain]; ok && strings.TrimSpace(existingApp.Name) == name {
				appID = types.StringValue(existingApp.AppID)
				praAppID = types.StringValue(existingApp.ID)
			} else {
				// Clear stale IDs (matching SDKv2 logic)
				appID = types.StringValue("")
				praAppID = types.StringValue("")
			}

			obj, objDiags := types.ObjectValue(appAttrTypes, map[string]attr.Value{
				"app_id":               appID,
				"pra_app_id":           praAppID,
				"name":                 item.Name,
				"description":          item.Description,
				"app_types":            item.AppTypes,
				"application_port":     item.ApplicationPort,
				"application_protocol": item.ApplicationProtocol,
				"connection_security":  item.ConnectionSecurity,
				"domain":               item.Domain,
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
		"pra_app_id":           types.StringType,
		"name":                 types.StringType,
		"description":          types.StringType,
		"app_types":            types.SetType{ElemType: types.StringType},
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"connection_security":  types.StringType,
		"domain":               types.StringType,
	}}, updatedAppConfigs)
	diags.Append(listDiags...)
	if diags.HasError() {
		return commonAppsDto, diags
	}

	commonAttrTypes := map[string]attr.Type{
		"apps_config": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"app_id":               types.StringType,
			"pra_app_id":           types.StringType,
			"name":                 types.StringType,
			"description":          types.StringType,
			"app_types":            types.SetType{ElemType: types.StringType},
			"application_port":     types.StringType,
			"application_protocol": types.StringType,
			"connection_security":  types.StringType,
			"domain":               types.StringType,
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

func validatePRAPlan(ctx context.Context, plan *ApplicationSegmentPRAModel) diag.Diagnostics {
	var diags diag.Diagnostics

	domainNames, domainDiags := helpers.ListValueToStringSlice(ctx, plan.DomainNames)
	diags.Append(domainDiags...)
	if diags.HasError() {
		return diags
	}

	domainSet := make(map[string]struct{}, len(domainNames))
	for _, domain := range domainNames {
		trimmed := strings.TrimSpace(domain)
		if trimmed != "" {
			domainSet[strings.ToLower(trimmed)] = struct{}{}
		}
	}

	configs, configDiags := helpers.ExpandPRACommonApps(ctx, plan.CommonAppsDto)
	diags.Append(configDiags...)
	if diags.HasError() {
		return diags
	}

	for i, cfg := range configs {
		domain := strings.TrimSpace(cfg.Domain)
		if domain == "" {
			diags.AddError("Invalid configuration", fmt.Sprintf("common_apps_dto.apps_config[%d]: domain must be provided", i))
			continue
		}

		if _, exists := domainSet[strings.ToLower(domain)]; !exists {
			diags.AddError("Invalid configuration", fmt.Sprintf("common_apps_dto.apps_config[%d]: domain %q is not defined in domain_names", i, domain))
		}

		protocol := strings.ToUpper(strings.TrimSpace(cfg.ApplicationProtocol))
		connectionSecurity := strings.TrimSpace(cfg.ConnectionSecurity)
		if protocol == "RDP" {
			if connectionSecurity == "" {
				diags.AddError("Invalid configuration", fmt.Sprintf("common_apps_dto.apps_config[%d]: connection_security must be set when application_protocol is RDP", i))
			}
		} else if connectionSecurity != "" {
			diags.AddError("Invalid configuration", fmt.Sprintf("common_apps_dto.apps_config[%d]: connection_security can only be set when application_protocol is RDP", i))
		}
	}

	return diags
}
