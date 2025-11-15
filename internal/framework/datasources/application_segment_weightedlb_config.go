package datasources

import (
	"context"
	"fmt"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

var (
	_ datasource.DataSource              = &ApplicationSegmentWeightedLBConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationSegmentWeightedLBConfigDataSource{}
)

func NewApplicationSegmentWeightedLBConfigDataSource() datasource.DataSource {
	return &ApplicationSegmentWeightedLBConfigDataSource{}
}

type ApplicationSegmentWeightedLBConfigDataSource struct {
	client *client.Client
}

type ApplicationSegmentWeightedLBConfigDataSourceModel struct {
	ID                               types.String `tfsdk:"id"`
	ApplicationID                    types.String `tfsdk:"application_id"`
	ApplicationName                  types.String `tfsdk:"application_name"`
	MicrotenantID                    types.String `tfsdk:"microtenant_id"`
	WeightedLoadBalancing            types.Bool   `tfsdk:"weighted_load_balancing"`
	ApplicationToServerGroupMappings types.Set    `tfsdk:"application_to_server_group_mappings"`
}

type ApplicationToServerGroupMappingModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Passive types.Bool   `tfsdk:"passive"`
	Weight  types.String `tfsdk:"weight"`
}

func (d *ApplicationSegmentWeightedLBConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_segment_weightedlb_config"
}

func (d *ApplicationSegmentWeightedLBConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves weighted load balancer configuration for an application segment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this data source.",
			},
			"application_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Application segment identifier to query. One of application_id or application_name must be provided.",
			},
			"application_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Application segment name to query. One of application_id or application_name must be provided.",
			},
			"microtenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "Optional microtenant identifier.",
			},
			"weighted_load_balancing": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if weighted load balancing is enabled for the application segment.",
			},
		},
		Blocks: map[string]schema.Block{
			"application_to_server_group_mappings": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Server group mapping identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Server group name.",
						},
						"passive": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the server group is passive.",
						},
						"weight": schema.StringAttribute{
							Computed:    true,
							Description: "Assigned weight for the server group.",
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationSegmentWeightedLBConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *ApplicationSegmentWeightedLBConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Provider", "The provider must be configured before reading data sources.")
		return
	}

	var model ApplicationSegmentWeightedLBConfigDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := d.serviceForMicrotenant(model.MicrotenantID)

	applicationID := helpers.StringValue(model.ApplicationID)
	applicationName := helpers.StringValue(model.ApplicationName)

	if applicationID == "" && applicationName == "" {
		resp.Diagnostics.AddError(
			"Missing Application Identifier",
			"Either application_id or application_name must be provided.",
		)
		return
	}

	if applicationID == "" && applicationName != "" {
		app, _, err := applicationsegment.GetByName(ctx, service, applicationName)
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Failed to find application segment named %s: %v", applicationName, err),
			)
			return
		}
		applicationID = app.ID
		model.ApplicationID = types.StringValue(applicationID)
		model.ApplicationName = types.StringValue(app.Name)
	}

	config, _, err := applicationsegment.GetWeightedLoadBalancerConfig(ctx, service, applicationID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			resp.Diagnostics.AddError(
				"Not Found",
				fmt.Sprintf("Weighted load balancer config not found for application %s", applicationID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to retrieve weighted load balancer config for application %s: %v", applicationID, err),
		)
		return
	}

	if config == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("No weighted load balancer config returned for application %s", applicationID),
		)
		return
	}

	model.WeightedLoadBalancing = types.BoolValue(config.WeightedLoadBalancing)
	model.ApplicationID = types.StringValue(applicationID)

	mappings, diags := d.flattenApplicationToServerGroupMappings(ctx, config.ApplicationToServerGroupMaps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.ApplicationToServerGroupMappings = mappings

	stateID := applicationID
	if !model.MicrotenantID.IsNull() && !model.MicrotenantID.IsUnknown() {
		microTenantID := helpers.StringValue(model.MicrotenantID)
		if microTenantID != "" {
			stateID = fmt.Sprintf("%s:%s", microTenantID, applicationID)
		}
	}
	model.ID = types.StringValue(stateID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (d *ApplicationSegmentWeightedLBConfigDataSource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := d.client.Service
	id := helpers.StringValue(microtenantID)
	if id != "" {
		service = service.WithMicroTenant(id)
	}
	return service
}

func (d *ApplicationSegmentWeightedLBConfigDataSource) flattenApplicationToServerGroupMappings(ctx context.Context, mappings []applicationsegment.ApplicationToServerGroupMapping) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		"id":      types.StringType,
		"name":    types.StringType,
		"passive": types.BoolType,
		"weight":  types.StringType,
	}

	if len(mappings) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	values := make([]attr.Value, 0, len(mappings))
	for _, mapping := range mappings {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":      helpers.StringValueOrNull(mapping.ID),
			"name":    helpers.StringValueOrNull(mapping.Name),
			"passive": types.BoolValue(mapping.Passive),
			"weight":  helpers.StringValueOrNull(mapping.Weight),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	set, setDiags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(setDiags...)
	return set, diags
}
