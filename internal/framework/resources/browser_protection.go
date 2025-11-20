package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

var (
	_ resource.Resource                = &BrowserProtectionResource{}
	_ resource.ResourceWithConfigure   = &BrowserProtectionResource{}
	_ resource.ResourceWithImportState = &BrowserProtectionResource{}
)

func NewBrowserProtectionResource() resource.Resource {
	return &BrowserProtectionResource{}
}

type BrowserProtectionResource struct {
	client *client.Client
}

type BrowserProtectionModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	CriteriaFlagsMask types.String `tfsdk:"criteria_flags_mask"`
	DefaultCSP        types.Bool   `tfsdk:"default_csp"`
	Criteria          types.List   `tfsdk:"criteria"`
}

type BrowserProtectionCriteriaModel struct {
	FingerPrintCriteria types.List `tfsdk:"finger_print_criteria"`
}

type FingerPrintCriteriaModel struct {
	CollectLocation    types.Bool   `tfsdk:"collect_location"`
	FingerprintTimeout types.String `tfsdk:"fingerprint_timeout"`
	Browser            types.List   `tfsdk:"browser"`
	Location           types.List   `tfsdk:"location"`
	System             types.List   `tfsdk:"system"`
}

type BrowserCriteriaModel struct {
	BrowserEng     types.Bool `tfsdk:"browser_eng"`
	BrowserEngVer  types.Bool `tfsdk:"browser_eng_ver"`
	BrowserName    types.Bool `tfsdk:"browser_name"`
	BrowserVersion types.Bool `tfsdk:"browser_version"`
	Canvas         types.Bool `tfsdk:"canvas"`
	FlashVer       types.Bool `tfsdk:"flash_ver"`
	FpUsrAgentStr  types.Bool `tfsdk:"fp_usr_agent_str"`
	IsCookie       types.Bool `tfsdk:"is_cookie"`
	IsLocalStorage types.Bool `tfsdk:"is_local_storage"`
	IsSessStorage  types.Bool `tfsdk:"is_sess_storage"`
	Ja3            types.Bool `tfsdk:"ja3"`
	Mime           types.Bool `tfsdk:"mime"`
	Plugin         types.Bool `tfsdk:"plugin"`
	SilverlightVer types.Bool `tfsdk:"silverlight_ver"`
}

type LocationCriteriaModel struct {
	Lat types.Bool `tfsdk:"lat"`
	Lon types.Bool `tfsdk:"lon"`
}

type SystemCriteriaModel struct {
	AvailScreenResolution types.Bool `tfsdk:"avail_screen_resolution"`
	CPUArch               types.Bool `tfsdk:"cpu_arch"`
	CurrScreenResolution  types.Bool `tfsdk:"curr_screen_resolution"`
	Font                  types.Bool `tfsdk:"font"`
	JavaVer               types.Bool `tfsdk:"java_ver"`
	MobileDevType         types.Bool `tfsdk:"mobile_dev_type"`
	MonitorMobile         types.Bool `tfsdk:"monitor_mobile"`
	OSName                types.Bool `tfsdk:"os_name"`
	OSVersion             types.Bool `tfsdk:"os_version"`
	SysLang               types.Bool `tfsdk:"sys_lang"`
	Tz                    types.Bool `tfsdk:"tz"`
	UsrLang               types.Bool `tfsdk:"usr_lang"`
}

func (r *BrowserProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browser_protection"
}

func (r *BrowserProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	browserCriteriaBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"browser_eng":      schema.BoolAttribute{Optional: true},
				"browser_eng_ver":  schema.BoolAttribute{Optional: true},
				"browser_name":     schema.BoolAttribute{Optional: true},
				"browser_version":  schema.BoolAttribute{Optional: true},
				"canvas":           schema.BoolAttribute{Optional: true},
				"flash_ver":        schema.BoolAttribute{Optional: true},
				"fp_usr_agent_str": schema.BoolAttribute{Optional: true},
				"is_cookie":        schema.BoolAttribute{Optional: true},
				"is_local_storage": schema.BoolAttribute{Optional: true},
				"is_sess_storage":  schema.BoolAttribute{Optional: true},
				"ja3":              schema.BoolAttribute{Optional: true},
				"mime":             schema.BoolAttribute{Optional: true},
				"plugin":           schema.BoolAttribute{Optional: true},
				"silverlight_ver":  schema.BoolAttribute{Optional: true},
			},
		},
	}

	locationCriteriaBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"lat": schema.BoolAttribute{Optional: true},
				"lon": schema.BoolAttribute{Optional: true},
			},
		},
	}

	systemCriteriaBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"avail_screen_resolution": schema.BoolAttribute{Optional: true},
				"cpu_arch":                schema.BoolAttribute{Optional: true},
				"curr_screen_resolution":  schema.BoolAttribute{Optional: true},
				"font":                    schema.BoolAttribute{Optional: true},
				"java_ver":                schema.BoolAttribute{Optional: true},
				"mobile_dev_type":         schema.BoolAttribute{Optional: true},
				"monitor_mobile":          schema.BoolAttribute{Optional: true},
				"os_name":                 schema.BoolAttribute{Optional: true},
				"os_version":              schema.BoolAttribute{Optional: true},
				"sys_lang":                schema.BoolAttribute{Optional: true},
				"tz":                      schema.BoolAttribute{Optional: true},
				"usr_lang":                schema.BoolAttribute{Optional: true},
			},
		},
	}

	fingerPrintCriteriaBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"collect_location":    schema.BoolAttribute{Optional: true},
				"fingerprint_timeout": schema.StringAttribute{Optional: true},
			},
			Blocks: map[string]schema.Block{
				"browser":  browserCriteriaBlock,
				"location": locationCriteriaBlock,
				"system":   systemCriteriaBlock,
			},
		},
	}

	criteriaBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"finger_print_criteria": fingerPrintCriteriaBlock,
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages a browser protection profile. Note: This resource is commented out in SDKv2 and may have limitations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the browser protection profile",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the browser protection profile",
			},
			"criteria_flags_mask": schema.StringAttribute{
				Optional:    true,
				Description: "Criteria flags mask",
			},
			"default_csp": schema.BoolAttribute{
				Optional:    true,
				Description: "Default CSP (Content Security Policy)",
			},
		},
		Blocks: map[string]schema.Block{
			"criteria": criteriaBlock,
		},
	}
}

func (r *BrowserProtectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BrowserProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan BrowserProtectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := r.expandBrowserProtection(&plan)

	// Since there's no Create function in the SDK, we'll use UpdateBrowserProtectionProfile
	// First, try to find the profile by name to get its ID
	profileName := plan.Name.ValueString()
	profile, _, err := browser_protection.GetBrowserProtectionProfileByName(ctx, r.client.Service, profileName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to find browser protection profile with name %s: %v. Note: Browser protection profiles must exist before they can be managed.", profileName, err),
		)
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("Browser protection profile with name %s not found. Note: Browser protection profiles must exist before they can be managed.", profileName),
		)
		return
	}

	payload.ID = profile.ID

	// Use UpdateBrowserProtectionProfile to set the profile as active
	if _, err := browser_protection.UpdateBrowserProtectionProfile(ctx, r.client.Service, payload.ID); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update browser protection profile: %v", err))
		return
	}

	plan.ID = types.StringValue(profile.ID)

	state, readDiags := r.readBrowserProtection(ctx, profile.ID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BrowserProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var state BrowserProtectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readBrowserProtection(ctx, state.ID.ValueString())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	newState.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *BrowserProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before managing resources.")
		return
	}

	var plan BrowserProtectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileID := plan.ID.ValueString()

	// Check if the profile still exists before updating
	allProfiles, _, err := browser_protection.GetBrowserProtectionProfile(ctx, r.client.Service)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to list browser protection profiles: %v", err))
		return
	}

	var profileExists bool
	for _, profile := range allProfiles {
		if profile.ID == profileID {
			profileExists = true
			break
		}
	}

	if !profileExists {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Browser protection profile %s no longer exists", profileID))
		return
	}

	// Since there's no Update function in the SDK, we'll use UpdateBrowserProtectionProfile
	// This sets the profile as active
	if _, err := browser_protection.UpdateBrowserProtectionProfile(ctx, r.client.Service, profileID); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update browser protection profile: %v", err))
		return
	}

	state, readDiags := r.readBrowserProtection(ctx, profileID)
	resp.Diagnostics.Append(readDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = plan.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BrowserProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete is a no-op as per SDKv2 implementation
}

func (r *BrowserProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "Configure the provider before importing browser protection profile.")
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import requires the browser protection profile ID or name.")
		return
	}

	var profile *browser_protection.BrowserProtection
	var err error

	// Try to get by name first
	profile, _, err = browser_protection.GetBrowserProtectionProfileByName(ctx, r.client.Service, id)
	if err != nil {
		// If not found by name, try to find by ID in all profiles
		allProfiles, _, listErr := browser_protection.GetBrowserProtectionProfile(ctx, r.client.Service)
		if listErr != nil {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate browser protection profile %q: %v", id, listErr))
			return
		}
		for _, p := range allProfiles {
			if p.ID == id {
				profile = &p
				break
			}
		}
	}

	if profile == nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to locate browser protection profile %q", id))
		return
	}

	state, readDiags := r.readBrowserProtection(ctx, profile.ID)
	if readDiags.HasError() {
		resp.Diagnostics.Append(readDiags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BrowserProtectionResource) expandBrowserProtection(plan *BrowserProtectionModel) browser_protection.BrowserProtection {
	result := browser_protection.BrowserProtection{
		ID:                plan.ID.ValueString(),
		Name:              plan.Name.ValueString(),
		Description:       helpers.StringValue(plan.Description),
		CriteriaFlagsMask: helpers.StringValue(plan.CriteriaFlagsMask),
		DefaultCSP:        helpers.BoolValue(plan.DefaultCSP, false),
		Criteria:          r.expandCriteria(plan.Criteria),
	}
	return result
}

func (r *BrowserProtectionResource) expandCriteria(criteria types.List) browser_protection.Criteria {
	if criteria.IsNull() || criteria.IsUnknown() || len(criteria.Elements()) == 0 {
		return browser_protection.Criteria{}
	}

	var criteriaModel BrowserProtectionCriteriaModel
	diags := criteria.ElementsAs(context.Background(), &criteriaModel, false)
	if diags.HasError() {
		return browser_protection.Criteria{}
	}

	return browser_protection.Criteria{
		FingerPrintCriteria: r.expandFingerPrintCriteria(criteriaModel.FingerPrintCriteria),
	}
}

func (r *BrowserProtectionResource) expandFingerPrintCriteria(fpc types.List) browser_protection.FingerPrintCriteria {
	if fpc.IsNull() || fpc.IsUnknown() || len(fpc.Elements()) == 0 {
		return browser_protection.FingerPrintCriteria{}
	}

	var fpcModel FingerPrintCriteriaModel
	diags := fpc.ElementsAs(context.Background(), &fpcModel, false)
	if diags.HasError() {
		return browser_protection.FingerPrintCriteria{}
	}

	return browser_protection.FingerPrintCriteria{
		CollectLocation:    helpers.BoolValue(fpcModel.CollectLocation, false),
		FingerprintTimeout: helpers.StringValue(fpcModel.FingerprintTimeout),
		Browser:            r.expandBrowserCriteria(fpcModel.Browser),
		Location:           r.expandLocationCriteria(fpcModel.Location),
		System:             r.expandSystemCriteria(fpcModel.System),
	}
}

func (r *BrowserProtectionResource) expandBrowserCriteria(bc types.List) browser_protection.BrowserCriteria {
	if bc.IsNull() || bc.IsUnknown() || len(bc.Elements()) == 0 {
		return browser_protection.BrowserCriteria{}
	}

	var bcModel BrowserCriteriaModel
	diags := bc.ElementsAs(context.Background(), &bcModel, false)
	if diags.HasError() {
		return browser_protection.BrowserCriteria{}
	}

	return browser_protection.BrowserCriteria{
		BrowserEng:     helpers.BoolValue(bcModel.BrowserEng, false),
		BrowserEngVer:  helpers.BoolValue(bcModel.BrowserEngVer, false),
		BrowserName:    helpers.BoolValue(bcModel.BrowserName, false),
		BrowserVersion: helpers.BoolValue(bcModel.BrowserVersion, false),
		Canvas:         helpers.BoolValue(bcModel.Canvas, false),
		FlashVer:       helpers.BoolValue(bcModel.FlashVer, false),
		FpUsrAgentStr:  helpers.BoolValue(bcModel.FpUsrAgentStr, false),
		IsCookie:       helpers.BoolValue(bcModel.IsCookie, false),
		IsLocalStorage: helpers.BoolValue(bcModel.IsLocalStorage, false),
		IsSessStorage:  helpers.BoolValue(bcModel.IsSessStorage, false),
		Ja3:            helpers.BoolValue(bcModel.Ja3, false),
		Mime:           helpers.BoolValue(bcModel.Mime, false),
		Plugin:         helpers.BoolValue(bcModel.Plugin, false),
		SilverlightVer: helpers.BoolValue(bcModel.SilverlightVer, false),
	}
}

func (r *BrowserProtectionResource) expandLocationCriteria(lc types.List) browser_protection.LocationCriteria {
	if lc.IsNull() || lc.IsUnknown() || len(lc.Elements()) == 0 {
		return browser_protection.LocationCriteria{}
	}

	var lcModel LocationCriteriaModel
	diags := lc.ElementsAs(context.Background(), &lcModel, false)
	if diags.HasError() {
		return browser_protection.LocationCriteria{}
	}

	return browser_protection.LocationCriteria{
		Lat: helpers.BoolValue(lcModel.Lat, false),
		Lon: helpers.BoolValue(lcModel.Lon, false),
	}
}

func (r *BrowserProtectionResource) expandSystemCriteria(sc types.List) browser_protection.SystemCriteria {
	if sc.IsNull() || sc.IsUnknown() || len(sc.Elements()) == 0 {
		return browser_protection.SystemCriteria{}
	}

	var scModel SystemCriteriaModel
	diags := sc.ElementsAs(context.Background(), &scModel, false)
	if diags.HasError() {
		return browser_protection.SystemCriteria{}
	}

	return browser_protection.SystemCriteria{
		AvailScreenResolution: helpers.BoolValue(scModel.AvailScreenResolution, false),
		CPUArch:               helpers.BoolValue(scModel.CPUArch, false),
		CurrScreenResolution:  helpers.BoolValue(scModel.CurrScreenResolution, false),
		Font:                  helpers.BoolValue(scModel.Font, false),
		JavaVer:               helpers.BoolValue(scModel.JavaVer, false),
		MobileDevType:         helpers.BoolValue(scModel.MobileDevType, false),
		MonitorMobile:         helpers.BoolValue(scModel.MonitorMobile, false),
		OSName:                helpers.BoolValue(scModel.OSName, false),
		OSVersion:             helpers.BoolValue(scModel.OSVersion, false),
		SysLang:               helpers.BoolValue(scModel.SysLang, false),
		Tz:                    helpers.BoolValue(scModel.Tz, false),
		UsrLang:               helpers.BoolValue(scModel.UsrLang, false),
	}
}

func (r *BrowserProtectionResource) readBrowserProtection(ctx context.Context, profileID string) (BrowserProtectionModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	allProfiles, _, err := browser_protection.GetBrowserProtectionProfile(ctx, r.client.Service)
	if err != nil {
		return BrowserProtectionModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to read browser protection profiles: %v", err)),
		}
	}

	var profile *browser_protection.BrowserProtection
	for _, p := range allProfiles {
		if p.ID == profileID {
			profile = &p
			break
		}
	}

	if profile == nil {
		return BrowserProtectionModel{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("Not Found", fmt.Sprintf("Browser protection profile with id '%s' not found", profileID)),
		}
	}

	criteria, criteriaDiags := r.flattenCriteria(ctx, profile.Criteria)
	diags.Append(criteriaDiags...)

	return BrowserProtectionModel{
		ID:                types.StringValue(profile.ID),
		Name:              types.StringValue(profile.Name),
		Description:       helpers.StringValueOrNull(profile.Description),
		CriteriaFlagsMask: helpers.StringValueOrNull(profile.CriteriaFlagsMask),
		DefaultCSP:        types.BoolValue(profile.DefaultCSP),
		Criteria:          criteria,
	}, diags
}

func (r *BrowserProtectionResource) flattenCriteria(ctx context.Context, criteria browser_protection.Criteria) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	fpc, fpcDiags := r.flattenFingerPrintCriteria(ctx, criteria.FingerPrintCriteria)
	diags.Append(fpcDiags...)

	attrTypes := map[string]attr.Type{
		"finger_print_criteria": types.ListType{ElemType: types.ObjectType{AttrTypes: fingerPrintCriteriaAttrTypes()}},
	}

	objValue, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"finger_print_criteria": fpc,
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return list, diags
}

func (r *BrowserProtectionResource) flattenFingerPrintCriteria(ctx context.Context, fpc browser_protection.FingerPrintCriteria) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	browser, browserDiags := r.flattenBrowserCriteria(ctx, fpc.Browser)
	diags.Append(browserDiags...)

	location, locationDiags := r.flattenLocationCriteria(ctx, fpc.Location)
	diags.Append(locationDiags...)

	system, systemDiags := r.flattenSystemCriteria(ctx, fpc.System)
	diags.Append(systemDiags...)

	attrTypes := fingerPrintCriteriaAttrTypes()

	objValue, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"collect_location":    types.BoolValue(fpc.CollectLocation),
		"fingerprint_timeout": types.StringValue(fpc.FingerprintTimeout),
		"browser":             browser,
		"location":            location,
		"system":              system,
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return list, diags
}

func (r *BrowserProtectionResource) flattenBrowserCriteria(ctx context.Context, bc browser_protection.BrowserCriteria) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := browserCriteriaAttrTypes()

	objValue, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"browser_eng":      types.BoolValue(bc.BrowserEng),
		"browser_eng_ver":  types.BoolValue(bc.BrowserEngVer),
		"browser_name":     types.BoolValue(bc.BrowserName),
		"browser_version":  types.BoolValue(bc.BrowserVersion),
		"canvas":           types.BoolValue(bc.Canvas),
		"flash_ver":        types.BoolValue(bc.FlashVer),
		"fp_usr_agent_str": types.BoolValue(bc.FpUsrAgentStr),
		"is_cookie":        types.BoolValue(bc.IsCookie),
		"is_local_storage": types.BoolValue(bc.IsLocalStorage),
		"is_sess_storage":  types.BoolValue(bc.IsSessStorage),
		"ja3":              types.BoolValue(bc.Ja3),
		"mime":             types.BoolValue(bc.Mime),
		"plugin":           types.BoolValue(bc.Plugin),
		"silverlight_ver":  types.BoolValue(bc.SilverlightVer),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return list, diags
}

func (r *BrowserProtectionResource) flattenLocationCriteria(ctx context.Context, lc browser_protection.LocationCriteria) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := locationCriteriaAttrTypes()

	objValue, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"lat": types.BoolValue(lc.Lat),
		"lon": types.BoolValue(lc.Lon),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return list, diags
}

func (r *BrowserProtectionResource) flattenSystemCriteria(ctx context.Context, sc browser_protection.SystemCriteria) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := systemCriteriaAttrTypes()

	objValue, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"avail_screen_resolution": types.BoolValue(sc.AvailScreenResolution),
		"cpu_arch":                types.BoolValue(sc.CPUArch),
		"curr_screen_resolution":  types.BoolValue(sc.CurrScreenResolution),
		"font":                    types.BoolValue(sc.Font),
		"java_ver":                types.BoolValue(sc.JavaVer),
		"mobile_dev_type":         types.BoolValue(sc.MobileDevType),
		"monitor_mobile":          types.BoolValue(sc.MonitorMobile),
		"os_name":                 types.BoolValue(sc.OSName),
		"os_version":              types.BoolValue(sc.OSVersion),
		"sys_lang":                types.BoolValue(sc.SysLang),
		"tz":                      types.BoolValue(sc.Tz),
		"usr_lang":                types.BoolValue(sc.UsrLang),
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return list, diags
}

func fingerPrintCriteriaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"collect_location":    types.BoolType,
		"fingerprint_timeout": types.StringType,
		"browser":             types.ListType{ElemType: types.ObjectType{AttrTypes: browserCriteriaAttrTypes()}},
		"location":            types.ListType{ElemType: types.ObjectType{AttrTypes: locationCriteriaAttrTypes()}},
		"system":              types.ListType{ElemType: types.ObjectType{AttrTypes: systemCriteriaAttrTypes()}},
	}
}

func browserCriteriaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"browser_eng":      types.BoolType,
		"browser_eng_ver":  types.BoolType,
		"browser_name":     types.BoolType,
		"browser_version":  types.BoolType,
		"canvas":           types.BoolType,
		"flash_ver":        types.BoolType,
		"fp_usr_agent_str": types.BoolType,
		"is_cookie":        types.BoolType,
		"is_local_storage": types.BoolType,
		"is_sess_storage":  types.BoolType,
		"ja3":              types.BoolType,
		"mime":             types.BoolType,
		"plugin":           types.BoolType,
		"silverlight_ver":  types.BoolType,
	}
}

func locationCriteriaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"lat": types.BoolType,
		"lon": types.BoolType,
	}
}

func systemCriteriaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"avail_screen_resolution": types.BoolType,
		"cpu_arch":                types.BoolType,
		"curr_screen_resolution":  types.BoolType,
		"font":                    types.BoolType,
		"java_ver":                types.BoolType,
		"mobile_dev_type":         types.BoolType,
		"monitor_mobile":          types.BoolType,
		"os_name":                 types.BoolType,
		"os_version":              types.BoolType,
		"sys_lang":                types.BoolType,
		"tz":                      types.BoolType,
		"usr_lang":                types.BoolType,
	}
}
