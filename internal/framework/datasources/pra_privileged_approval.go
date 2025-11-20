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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

var (
	_ datasource.DataSource              = &PRAPrivilegedApprovalDataSource{}
	_ datasource.DataSourceWithConfigure = &PRAPrivilegedApprovalDataSource{}
)

func NewPRAPrivilegedApprovalDataSource() datasource.DataSource {
	return &PRAPrivilegedApprovalDataSource{}
}

type PRAPrivilegedApprovalDataSource struct {
	client *client.Client
}

type PRAPrivilegedApprovalModel struct {
	ID            types.String `tfsdk:"id"`
	EmailIDs      types.List   `tfsdk:"email_ids"`
	StartTime     types.String `tfsdk:"start_time"`
	EndTime       types.String `tfsdk:"end_time"`
	Status        types.String `tfsdk:"status"`
	CreationTime  types.String `tfsdk:"creation_time"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	WorkingHours  types.Set    `tfsdk:"working_hours"`
	Applications  types.List   `tfsdk:"applications"`
	MicroTenantID types.String `tfsdk:"microtenant_id"`
}

func (d *PRAPrivilegedApprovalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pra_privileged_approval"
}

func (d *PRAPrivilegedApprovalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a ZPA PRA privileged approval by ID or email ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the PRA privileged approval.",
			},
			"email_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Email address associated with the privileged approval.",
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"start_time": schema.StringAttribute{Computed: true},
			"end_time":   schema.StringAttribute{Computed: true},
			"status":     schema.StringAttribute{Computed: true},
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
			"working_hours": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"days": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"start_time": schema.StringAttribute{Computed: true},
						"start_time_cron": schema.StringAttribute{
							Computed: true,
						},
						"end_time": schema.StringAttribute{Computed: true},
						"end_time_cron": schema.StringAttribute{
							Computed: true,
						},
						"timezone": schema.StringAttribute{Computed: true},
					},
				},
			},
			"applications": schema.ListNestedBlock{
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

func (d *PRAPrivilegedApprovalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PRAPrivilegedApprovalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data PRAPrivilegedApprovalModel
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
	var email string
	if !data.EmailIDs.IsNull() && !data.EmailIDs.IsUnknown() {
		var emails []string
		diags := data.EmailIDs.ElementsAs(ctx, &emails, false)
		resp.Diagnostics.Append(diags...)
		if len(emails) > 0 {
			email = strings.TrimSpace(emails[0])
		}
	}

	if id == "" && email == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Provide either 'id' or at least one 'email_ids' entry.")
		return
	}

	var approval *praapproval.PrivilegedApproval
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving PRA privileged approval by ID", map[string]any{"id": id})
		approval, _, err = praapproval.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving PRA privileged approval by email ID", map[string]any{"email_id": email})
		approval, _, err = praapproval.GetByEmailID(ctx, service, email)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PRA privileged approval: %v", err))
		return
	}

	if approval == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PRA privileged approval not found with id %q and email %q.", id, email))
		return
	}

	workingHours, whDiags := flattenPRAWorkingHours(ctx, approval.WorkingHours)
	resp.Diagnostics.Append(whDiags...)
	applications, appDiags := flattenPRAApplications(ctx, approval.Applications)
	resp.Diagnostics.Append(appDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailList, emailDiags := flattenEmailIDs(ctx, approval.EmailIDs)
	resp.Diagnostics.Append(emailDiags...)

	data.ID = types.StringValue(approval.ID)
	data.EmailIDs = emailList
	data.StartTime = types.StringValue(approval.StartTime)
	data.EndTime = types.StringValue(approval.EndTime)
	data.Status = types.StringValue(approval.Status)
	data.CreationTime = types.StringValue(approval.CreationTime)
	data.ModifiedBy = types.StringValue(approval.ModifiedBy)
	data.ModifiedTime = types.StringValue(approval.ModifiedTime)
	data.WorkingHours = workingHours
	data.Applications = applications

	if !data.MicroTenantID.IsNull() && strings.TrimSpace(data.MicroTenantID.ValueString()) != "" {
		// retain provided value
	} else if approval.MicroTenantID != "" {
		data.MicroTenantID = types.StringValue(approval.MicroTenantID)
	} else {
		data.MicroTenantID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenEmailIDs(ctx context.Context, ids []string) (types.List, diag.Diagnostics) {
	return types.ListValueFrom(ctx, types.StringType, ids)
}

func flattenPRAWorkingHours(ctx context.Context, wh *praapproval.WorkingHours) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"days":            types.SetType{ElemType: types.StringType},
		"start_time":      types.StringType,
		"start_time_cron": types.StringType,
		"end_time":        types.StringType,
		"end_time_cron":   types.StringType,
		"timezone":        types.StringType,
	}

	if wh == nil {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	days, daysDiags := types.SetValueFrom(ctx, types.StringType, wh.Days)
	diags.Append(daysDiags...)

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"days":            days,
		"start_time":      types.StringValue(wh.StartTime),
		"start_time_cron": types.StringValue(wh.StartTimeCron),
		"end_time":        types.StringValue(wh.EndTime),
		"end_time_cron":   types.StringValue(wh.EndTimeCron),
		"timezone":        types.StringValue(wh.TimeZone),
	})
	diags.Append(objDiags...)

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(setDiags...)
	return set, diags
}

func flattenPRAApplications(ctx context.Context, apps []praapproval.Applications) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}

	if len(apps) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(apps))
	for _, app := range apps {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(app.ID),
			"name": types.StringValue(app.Name),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
