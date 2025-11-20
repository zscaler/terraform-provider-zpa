package datasources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

var (
	_ datasource.DataSource              = &SCIMGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &SCIMGroupDataSource{}
)

func NewSCIMGroupDataSource() datasource.DataSource {
	return &SCIMGroupDataSource{}
}

type SCIMGroupDataSource struct {
	client *client.Client
}

type SCIMGroupModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	IdpGroupID   types.String `tfsdk:"idp_group_id"`
	IdpID        types.String `tfsdk:"idp_id"`
	IdpName      types.String `tfsdk:"idp_name"`
	CreationTime types.Int64  `tfsdk:"creation_time"`
	ModifiedTime types.Int64  `tfsdk:"modified_time"`
}

func (d *SCIMGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_groups"
}

func (d *SCIMGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a SCIM group by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the SCIM group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the SCIM group.",
			},
			"idp_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the IDP associated with the SCIM group.",
			},
			"idp_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the IDP associated with the SCIM group.",
			},
			"idp_group_id": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.Int64Attribute{
				Computed: true,
			},
			"modified_time": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (d *SCIMGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SCIMGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data SCIMGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.TrimSpace(data.ID.ValueString())
	name := strings.TrimSpace(data.Name.ValueString())
	idpID := strings.TrimSpace(data.IdpID.ValueString())
	idpName := strings.TrimSpace(data.IdpName.ValueString())

	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be provided.")
		return
	}

	if id != "" {
		tflog.Debug(ctx, "Retrieving SCIM group by ID", map[string]any{"id": id})
		group, _, err := scimgroup.Get(ctx, d.client.Service, id)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SCIM group: %v", err))
			return
		}
		if group == nil {
			resp.Diagnostics.AddError("Not Found", fmt.Sprintf("SCIM group with id %q not found.", id))
			return
		}
		state, diags := flattenSCIMGroup(group)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	if idpID == "" && idpName == "" {
		resp.Diagnostics.AddError("Missing Required Attribute", "Either 'idp_id' or 'idp_name' must be provided when looking up by name.")
		return
	}

	service := d.client.Service

	var (
		idp *idpcontroller.IdpController
		err error
	)
	if idpID != "" {
		idp, _, err = idpcontroller.Get(ctx, service, idpID)
	} else {
		idp, _, err = idpcontroller.GetByName(ctx, service, idpName)
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read IDP: %v", err))
		return
	}
	if idp == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to locate IDP with id %q or name %q.", idpID, idpName))
		return
	}

	tflog.Debug(ctx, "Retrieving SCIM group by name", map[string]any{
		"name":     name,
		"idp_id":   idp.ID,
		"idp_name": idp.Name,
	})

	group, _, err := scimgroup.GetByName(ctx, service, name, idp.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SCIM group: %v", err))
		return
	}
	if group == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("SCIM group with name %q not found.", name))
		return
	}

	state, diags := flattenSCIMGroup(group)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.IdpID = types.StringValue(idp.ID)
	state.IdpName = types.StringValue(idp.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenSCIMGroup(group *scimgroup.ScimGroup) (SCIMGroupModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	idpID := ""
	if group.IdpID != 0 {
		idpID = strconv.FormatInt(group.IdpID, 10)
	}

	model := SCIMGroupModel{
		ID:         types.StringValue(strconv.FormatInt(group.ID, 10)),
		Name:       types.StringValue(group.Name),
		IdpGroupID: types.StringValue(group.IdpGroupID),
		IdpID:      types.StringValue(idpID),
		IdpName:    types.StringValue(group.IdpName),
	}

	if group.CreationTime != 0 {
		model.CreationTime = types.Int64Value(group.CreationTime)
	} else {
		model.CreationTime = types.Int64Null()
	}

	if group.ModifiedTime != 0 {
		model.ModifiedTime = types.Int64Value(group.ModifiedTime)
	} else {
		model.ModifiedTime = types.Int64Null()
	}

	return model, diags
}
