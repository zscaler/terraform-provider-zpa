// Copyright (c) SecurityGeekIO, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

var (
	_ resource.Resource                = &SegmentGroupsResource{}
	_ resource.ResourceWithConfigure   = &SegmentGroupsResource{}
	_ resource.ResourceWithImportState = &SegmentGroupsResource{}
)

var policyRulesDetachLock sync.Mutex

type SegmentGroupsResource struct {
	client *client.Client
}

type SegmentGroupsResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	MicroTenantID types.String `tfsdk:"microtenant_id"`
}

func NewSegmentGroupsResource() resource.Resource {
	return &SegmentGroupsResource{}
}

func (r *SegmentGroupsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment_group"
}

func (r *SegmentGroupsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZPA Segment Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the segment group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the segment group.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the segment group.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this segment group is enabled.",
				Optional:    true,
			},
			"microtenant_id": schema.StringAttribute{
				Description: "Microtenant ID to scope segment group operations.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SegmentGroupsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SegmentGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SegmentGroupsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)
	segmentGroupReq, diags := expandSegmentGroup(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating segment group", map[string]interface{}{"name": segmentGroupReq.Name})

	created, _, err := segmentgroup.Create(ctx, service, &segmentGroupReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create segment group: %s", err))
		return
	}

	if diags := flattenSegmentGroup(ctx, created, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SegmentGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SegmentGroupsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing Segment Group ID", "Segment group ID is required to read the resource")
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	segmentGroup, _, err := segmentgroup.Get(ctx, service, state.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			tflog.Warn(ctx, "Segment group not found, removing from state", map[string]interface{}{"id": state.ID.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read segment group: %s", err))
		return
	}

	if diags := flattenSegmentGroup(ctx, segmentGroup, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SegmentGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SegmentGroupsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing Segment Group ID", "Segment group ID is required to update the resource")
		return
	}

	service := r.serviceForMicrotenant(plan.MicroTenantID)

	// Check if resource still exists before updating
	if _, _, err := segmentgroup.Get(ctx, service, plan.ID.ValueString()); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	segmentGroupReq, diags := expandSegmentGroup(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating segment group", map[string]interface{}{"id": plan.ID.ValueString()})

	if _, err := segmentgroup.Update(ctx, service, plan.ID.ValueString(), &segmentGroupReq); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update segment group: %s", err))
		return
	}

	updated, _, err := segmentgroup.Get(ctx, service, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read segment group after update: %s", err))
		return
	}

	if diags := flattenSegmentGroup(ctx, updated, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SegmentGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SegmentGroupsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing Segment Group ID", "Segment group ID is required to delete the resource")
		return
	}

	service := r.serviceForMicrotenant(state.MicroTenantID)

	tflog.Info(ctx, "Deleting segment group", map[string]interface{}{"id": state.ID.ValueString()})

	if err := detachSegmentGroupFromAllPolicyRules(ctx, state.ID.ValueString(), service); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach segment group from policies: %s", err))
		return
	}

	if _, err := segmentgroup.Delete(ctx, service, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete segment group: %s", err))
		return
	}
}

func (r *SegmentGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	service := r.client.Service

	if _, err := strconv.Atoi(id); err == nil {
		segmentGroup, _, err := segmentgroup.Get(ctx, service, id)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import segment group by ID: %s", err))
			return
		}

		var state SegmentGroupsResourceModel
		if diags := flattenSegmentGroup(ctx, segmentGroup, &state); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	segmentGroup, _, err := segmentgroup.GetByName(ctx, service, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import segment group by name: %s", err))
		return
	}

	var state SegmentGroupsResourceModel
	if diags := flattenSegmentGroup(ctx, segmentGroup, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SegmentGroupsResource) serviceForMicrotenant(microtenantID types.String) *zscaler.Service {
	service := r.client.Service
	if !microtenantID.IsNull() && !microtenantID.IsUnknown() && microtenantID.ValueString() != "" {
		service = service.WithMicroTenant(microtenantID.ValueString())
	}
	return service
}

func expandSegmentGroup(ctx context.Context, data SegmentGroupsResourceModel) (segmentgroup.SegmentGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Expanding segment group", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	return segmentgroup.SegmentGroup{
		ID:            data.ID.ValueString(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueString(),
		Enabled:       !data.Enabled.IsNull() && data.Enabled.ValueBool(),
		MicroTenantID: data.MicroTenantID.ValueString(),
	}, diags
}

func flattenSegmentGroup(ctx context.Context, segmentGroup *segmentgroup.SegmentGroup, data *SegmentGroupsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Flattening segment group", map[string]interface{}{
		"id":   segmentGroup.ID,
		"name": segmentGroup.Name,
	})

	data.ID = types.StringValue(segmentGroup.ID)
	data.Name = types.StringValue(segmentGroup.Name)
	data.Description = types.StringValue(segmentGroup.Description)
	data.Enabled = types.BoolValue(segmentGroup.Enabled)
	data.MicroTenantID = types.StringValue(segmentGroup.MicroTenantID)

	return diags
}

func detachSegmentGroupFromAllPolicyRules(ctx context.Context, id string, service *zscaler.Service) error {
	typesList := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}

	policyRulesDetachLock.Lock()
	defer policyRulesDetachLock.Unlock()

	for _, policyType := range typesList {
		policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, policyType)
		if err != nil {
			tflog.Warn(ctx, "Failed to fetch policy set", map[string]interface{}{"type": policyType, "error": err.Error()})
			continue
		}

		rules, _, err := policysetcontroller.GetAllByType(ctx, service, policyType)
		if err != nil {
			tflog.Warn(ctx, "Failed to fetch policy rules", map[string]interface{}{"type": policyType, "error": err.Error()})
			continue
		}

		for _, rule := range rules {
			changed := false
			for i, condition := range rule.Conditions {
				operands := make([]policysetcontroller.Operands, 0, len(condition.Operands))
				for _, operand := range condition.Operands {
					if operand.ObjectType == "APP_GROUP" && operand.LHS == "id" && operand.RHS == id {
						changed = true
						continue
					}
					operands = append(operands, operand)
				}
				rule.Conditions[i].Operands = operands
			}

			if !changed {
				continue
			}

			if _, err := policysetcontroller.UpdateRule(ctx, service, policySet.ID, rule.ID, &rule); err != nil {
				return fmt.Errorf("failed to update policy rule %s: %w", rule.ID, err)
			}
		}
	}

	return nil
}
