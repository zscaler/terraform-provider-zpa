package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/clienttypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccessPolicyClientTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessPolicyClientTypesRead,
		Importer:    &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"zpn_client_type_exporter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_exporter_noauth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_browser_isolation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_machine_tunnel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_ip_anchoring": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_edge_connector": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_zapp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_slogger": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_client_type_branch_connector": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessPolicyClientTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Getting data for all client types set\n")

	resp, _, err := clienttypes.GetAllClientTypes(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved client types:\n%+v\n", resp)
	d.SetId("client_types")
	if err := d.Set("zpn_client_type_exporter", resp.ZPNClientTypeExplorer); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_exporter: %v", err))
	}
	if err := d.Set("zpn_client_type_exporter_noauth", resp.ZPNClientTypeNoAuth); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_exporter_noauth: %v", err))
	}
	if err := d.Set("zpn_client_type_browser_isolation", resp.ZPNClientTypeBrowserIsolation); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_browser_isolation: %v", err))
	}
	if err := d.Set("zpn_client_type_machine_tunnel", resp.ZPNClientTypeMachineTunnel); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_machine_tunnel: %v", err))
	}
	if err := d.Set("zpn_client_type_ip_anchoring", resp.ZPNClientTypeIPAnchoring); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_ip_anchoring: %v", err))
	}
	if err := d.Set("zpn_client_type_edge_connector", resp.ZPNClientTypeEdgeConnector); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_edge_connector: %v", err))
	}
	if err := d.Set("zpn_client_type_zapp", resp.ZPNClientTypeZAPP); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_zapp: %v", err))
	}
	if err := d.Set("zpn_client_type_slogger", resp.ZPNClientTypeSlogger); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_slogger: %v", err))
	}
	if err := d.Set("zpn_client_type_branch_connector", resp.ZPNClientTypeBranchConnector); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set zpn_client_type_branch_connector: %v", err))
	}

	return nil
}
