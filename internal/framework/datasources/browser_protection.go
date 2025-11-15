package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

var (
	_ datasource.DataSource              = &BrowserProtectionDataSource{}
	_ datasource.DataSourceWithConfigure = &BrowserProtectionDataSource{}
)

func NewBrowserProtectionDataSource() datasource.DataSource {
	return &BrowserProtectionDataSource{}
}

type BrowserProtectionDataSource struct {
	client *client.Client
}

type BrowserProtectionModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	DefaultCSP        types.Bool   `tfsdk:"default_csp"`
	CreationTime      types.String `tfsdk:"creation_time"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
	ModifiedTime      types.String `tfsdk:"modified_time"`
	CriteriaFlagsMask types.String `tfsdk:"criteria_flags_mask"`
	Criteria          types.Set    `tfsdk:"criteria"`
}

func (d *BrowserProtectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browser_protection"
}

func (d *BrowserProtectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	browserBlock := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"browser_eng":      schema.BoolAttribute{Computed: true},
			"browser_eng_ver":  schema.BoolAttribute{Computed: true},
			"browser_name":     schema.BoolAttribute{Computed: true},
			"browser_version":  schema.BoolAttribute{Computed: true},
			"canvas":           schema.BoolAttribute{Computed: true},
			"flash_ver":        schema.BoolAttribute{Computed: true},
			"fp_usr_agent_str": schema.BoolAttribute{Computed: true},
			"is_cookie":        schema.BoolAttribute{Computed: true},
			"is_local_storage": schema.BoolAttribute{Computed: true},
			"is_sess_storage":  schema.BoolAttribute{Computed: true},
			"ja3":              schema.BoolAttribute{Computed: true},
			"mime":             schema.BoolAttribute{Computed: true},
			"plugin":           schema.BoolAttribute{Computed: true},
			"silverlight_ver":  schema.BoolAttribute{Computed: true},
		},
	}

	locationBlock := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"lat": schema.BoolAttribute{Computed: true},
			"lon": schema.BoolAttribute{Computed: true},
		},
	}

	systemBlock := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"avail_screen_resolution": schema.BoolAttribute{Computed: true},
			"cpu_arch":                schema.BoolAttribute{Computed: true},
			"curr_screen_resolution":  schema.BoolAttribute{Computed: true},
			"font":                    schema.BoolAttribute{Computed: true},
			"java_ver":                schema.BoolAttribute{Computed: true},
			"mobile_dev_type":         schema.BoolAttribute{Computed: true},
			"monitor_mobile":          schema.BoolAttribute{Computed: true},
			"os_name":                 schema.BoolAttribute{Computed: true},
			"os_version":              schema.BoolAttribute{Computed: true},
			"sys_lang":                schema.BoolAttribute{Computed: true},
			"tz":                      schema.BoolAttribute{Computed: true},
			"usr_lang":                schema.BoolAttribute{Computed: true},
		},
	}

	fingerPrintBlock := schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"collect_location":    schema.BoolAttribute{Computed: true},
			"fingerprint_timeout": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"browser":  browserBlock,
			"location": locationBlock,
			"system":   systemBlock,
		},
	}

	criteriaBlock := schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"finger_print_criteria": fingerPrintBlock,
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieves a browser protection profile by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the browser protection profile.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Optional name of the browser protection profile.",
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"default_csp": schema.BoolAttribute{
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
			"criteria_flags_mask": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"criteria": criteriaBlock,
		},
	}
}

func (d *BrowserProtectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BrowserProtectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data BrowserProtectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(data.Name.ValueString())

	tflog.Debug(ctx, "Retrieving browser protection profile", map[string]any{"name": name})

	var profile *browser_protection.BrowserProtection
	var err error

	if name != "" {
		profile, _, err = browser_protection.GetBrowserProtectionProfileByName(ctx, d.client.Service, name)
	} else {
		profiles, _, e := browser_protection.GetBrowserProtectionProfile(ctx, d.client.Service)
		err = e
		if len(profiles) > 0 {
			profile = &profiles[0]
		}
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read browser protection profile: %v", err))
		return
	}

	if profile == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Browser protection profile with name %q not found.", name))
		return
	}

	id := profile.ID
	if id == "" {
		id = helpers.GenerateShortID(profile.Name)
	}

	data.ID = types.StringValue(id)
	data.Name = types.StringValue(profile.Name)
	data.Description = types.StringValue(profile.Description)
	data.DefaultCSP = types.BoolValue(profile.DefaultCSP)
	data.CreationTime = types.StringValue(profile.CreationTime)
	data.ModifiedBy = types.StringValue(profile.ModifiedBy)
	data.ModifiedTime = types.StringValue(profile.ModifiedTime)
	data.CriteriaFlagsMask = types.StringValue(profile.CriteriaFlagsMask)

	criteriaSet, diags := flattenBrowserProtectionCriteria(profile.Criteria)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Criteria = criteriaSet

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenBrowserProtectionCriteria(criteria browser_protection.Criteria) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	browserAttrTypes := map[string]attr.Type{
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

	locationAttrTypes := map[string]attr.Type{
		"lat": types.BoolType,
		"lon": types.BoolType,
	}

	systemAttrTypes := map[string]attr.Type{
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

	fingerAttrTypes := map[string]attr.Type{
		"collect_location":    types.BoolType,
		"fingerprint_timeout": types.StringType,
		"browser":             types.ObjectType{AttrTypes: browserAttrTypes},
		"location":            types.ObjectType{AttrTypes: locationAttrTypes},
		"system":              types.ObjectType{AttrTypes: systemAttrTypes},
	}

	criteriaAttrTypes := map[string]attr.Type{
		"finger_print_criteria": types.ObjectType{AttrTypes: fingerAttrTypes},
	}

	fingerValue, fingerDiags := flattenFingerPrintCriteriaObject(criteria.FingerPrintCriteria, browserAttrTypes, locationAttrTypes, systemAttrTypes)
	diags.Append(fingerDiags...)

	objValue, objDiags := types.ObjectValue(criteriaAttrTypes, map[string]attr.Value{
		"finger_print_criteria": fingerValue,
	})
	diags.Append(objDiags...)

	setValue, setDiags := types.SetValue(types.ObjectType{AttrTypes: criteriaAttrTypes}, []attr.Value{objValue})
	diags.Append(setDiags...)

	return setValue, diags
}

func flattenFingerPrintCriteriaObject(criteria browser_protection.FingerPrintCriteria, browserAttrTypes, locationAttrTypes, systemAttrTypes map[string]attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	browserValue, browserDiags := flattenBrowserCriteriaObject(criteria.Browser, browserAttrTypes)
	diags.Append(browserDiags...)

	locationValue, locationDiags := flattenLocationCriteriaObject(criteria.Location, locationAttrTypes)
	diags.Append(locationDiags...)

	systemValue, systemDiags := flattenSystemCriteriaObject(criteria.System, systemAttrTypes)
	diags.Append(systemDiags...)

	value, valueDiags := types.ObjectValue(map[string]attr.Type{
		"collect_location":    types.BoolType,
		"fingerprint_timeout": types.StringType,
		"browser":             types.ObjectType{AttrTypes: browserAttrTypes},
		"location":            types.ObjectType{AttrTypes: locationAttrTypes},
		"system":              types.ObjectType{AttrTypes: systemAttrTypes},
	}, map[string]attr.Value{
		"collect_location":    types.BoolValue(criteria.CollectLocation),
		"fingerprint_timeout": types.StringValue(criteria.FingerprintTimeout),
		"browser":             browserValue,
		"location":            locationValue,
		"system":              systemValue,
	})
	diags.Append(valueDiags...)

	return value, diags
}

func flattenBrowserCriteriaObject(criteria browser_protection.BrowserCriteria, attrTypes map[string]attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	value, valueDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"browser_eng":      types.BoolValue(criteria.BrowserEng),
		"browser_eng_ver":  types.BoolValue(criteria.BrowserEngVer),
		"browser_name":     types.BoolValue(criteria.BrowserName),
		"browser_version":  types.BoolValue(criteria.BrowserVersion),
		"canvas":           types.BoolValue(criteria.Canvas),
		"flash_ver":        types.BoolValue(criteria.FlashVer),
		"fp_usr_agent_str": types.BoolValue(criteria.FpUsrAgentStr),
		"is_cookie":        types.BoolValue(criteria.IsCookie),
		"is_local_storage": types.BoolValue(criteria.IsLocalStorage),
		"is_sess_storage":  types.BoolValue(criteria.IsSessStorage),
		"ja3":              types.BoolValue(criteria.Ja3),
		"mime":             types.BoolValue(criteria.Mime),
		"plugin":           types.BoolValue(criteria.Plugin),
		"silverlight_ver":  types.BoolValue(criteria.SilverlightVer),
	})
	diags.Append(valueDiags...)

	return value, diags
}

func flattenLocationCriteriaObject(criteria browser_protection.LocationCriteria, attrTypes map[string]attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	value, valueDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"lat": types.BoolValue(criteria.Lat),
		"lon": types.BoolValue(criteria.Lon),
	})
	diags.Append(valueDiags...)

	return value, diags
}

func flattenSystemCriteriaObject(criteria browser_protection.SystemCriteria, attrTypes map[string]attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	value, valueDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"avail_screen_resolution": types.BoolValue(criteria.AvailScreenResolution),
		"cpu_arch":                types.BoolValue(criteria.CPUArch),
		"curr_screen_resolution":  types.BoolValue(criteria.CurrScreenResolution),
		"font":                    types.BoolValue(criteria.Font),
		"java_ver":                types.BoolValue(criteria.JavaVer),
		"mobile_dev_type":         types.BoolValue(criteria.MobileDevType),
		"monitor_mobile":          types.BoolValue(criteria.MonitorMobile),
		"os_name":                 types.BoolValue(criteria.OSName),
		"os_version":              types.BoolValue(criteria.OSVersion),
		"sys_lang":                types.BoolValue(criteria.SysLang),
		"tz":                      types.BoolValue(criteria.Tz),
		"usr_lang":                types.BoolValue(criteria.UsrLang),
	})
	diags.Append(valueDiags...)

	return value, diags
}
