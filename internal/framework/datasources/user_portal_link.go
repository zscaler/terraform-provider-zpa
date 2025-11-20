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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

var (
	_ datasource.DataSource              = &UserPortalLinkDataSource{}
	_ datasource.DataSourceWithConfigure = &UserPortalLinkDataSource{}
)

func NewUserPortalLinkDataSource() datasource.DataSource {
	return &UserPortalLinkDataSource{}
}

type UserPortalLinkDataSource struct {
	client *client.Client
}

type UserPortalLinkModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ApplicationID types.String `tfsdk:"application_id"`
	CreationTime  types.String `tfsdk:"creation_time"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	IconText      types.String `tfsdk:"icon_text"`
	Link          types.String `tfsdk:"link"`
	LinkPath      types.String `tfsdk:"link_path"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	Protocol      types.String `tfsdk:"protocol"`
	MicroTenantID types.String `tfsdk:"microtenant_id"`
	UserPortalID  types.String `tfsdk:"user_portal_id"`
	UserPortals   types.List   `tfsdk:"user_portals"`
}

func (d *UserPortalLinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_portal_link"
}

func (d *UserPortalLinkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a user portal link by ID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifier of the user portal link.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the user portal link.",
			},
			"application_id": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"icon_text": schema.StringAttribute{
				Computed: true,
			},
			"link": schema.StringAttribute{
				Computed: true,
			},
			"link_path": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_time": schema.StringAttribute{
				Computed: true,
			},
			"protocol": schema.StringAttribute{
				Computed: true,
			},
			"microtenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Micro-tenant ID used to scope the lookup.",
			},
			"user_portal_id": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"user_portals": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true},
						"name":    schema.StringAttribute{Computed: true},
						"enabled": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *UserPortalLinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserPortalLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Provider",
			"The ZPA provider was not configured. Configure the provider before using this data source.",
		)
		return
	}

	var data UserPortalLinkModel
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

	var link *portal_link.UserPortalLink
	var err error

	if id != "" {
		tflog.Debug(ctx, "Retrieving user portal link by ID", map[string]any{"id": id})
		link, _, err = portal_link.Get(ctx, service, id)
	} else {
		tflog.Debug(ctx, "Retrieving user portal link by name", map[string]any{"name": name})
		link, _, err = portal_link.GetByName(ctx, service, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user portal link: %v", err))
		return
	}

	if link == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("User portal link with id %q or name %q not found.", id, name))
		return
	}

	state, diags := flattenUserPortalLink(ctx, link)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenUserPortalLink(ctx context.Context, link *portal_link.UserPortalLink) (UserPortalLinkModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	userPortalAttrTypes := map[string]attr.Type{
		"id":      types.StringType,
		"name":    types.StringType,
		"enabled": types.BoolType,
	}

	userPortalValues := make([]attr.Value, 0, len(link.UserPortals))
	for _, portal := range link.UserPortals {
		obj, objDiags := types.ObjectValue(userPortalAttrTypes, map[string]attr.Value{
			"id":      types.StringValue(portal.ID),
			"name":    types.StringValue(portal.Name),
			"enabled": types.BoolValue(portal.Enabled),
		})
		diags.Append(objDiags...)
		userPortalValues = append(userPortalValues, obj)
	}

	userPortalList, listDiags := types.ListValue(types.ObjectType{AttrTypes: userPortalAttrTypes}, userPortalValues)
	diags.Append(listDiags...)

	model := UserPortalLinkModel{
		ID:            types.StringValue(link.ID),
		Name:          types.StringValue(link.Name),
		ApplicationID: types.StringValue(link.ApplicationID),
		CreationTime:  types.StringValue(link.CreationTime),
		Description:   types.StringValue(link.Description),
		Enabled:       types.BoolValue(link.Enabled),
		IconText:      types.StringValue(link.IconText),
		Link:          types.StringValue(link.Link),
		LinkPath:      types.StringValue(link.LinkPath),
		ModifiedBy:    types.StringValue(link.ModifiedBy),
		ModifiedTime:  types.StringValue(link.ModifiedTime),
		Protocol:      types.StringValue(link.Protocol),
		MicroTenantID: types.StringValue(link.MicrotenantID),
		UserPortalID:  types.StringValue(link.UserPortalID),
		UserPortals:   userPortalList,
	}

	return model, diags
}
