package zpa

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

func dataSourceApplicationSegmentMultimatchBulk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationSegmentMultimatchBulkRead,
		Schema: map[string]*schema.Schema{
			"domain_names": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of domain names to check for unsupported multimatch references.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"unsupported_references": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of application segments that cannot support multimatch.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application segment ID.",
						},
						"app_segment_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application segment name.",
						},
						"domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of domain names for this segment.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tcp_ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of TCP ports for this segment.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"match_style": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current match style of the segment (EXCLUSIVE or INCLUSIVE).",
						},
						"microtenant_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microtenant name associated with this segment.",
						},
					},
				},
			},
		},
	}
}

func dataSourceApplicationSegmentMultimatchBulkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	domainNamesInterface := d.Get("domain_names").([]interface{})
	if len(domainNamesInterface) == 0 {
		return diag.FromErr(fmt.Errorf("at least one domain name must be provided"))
	}

	// Convert interface slice to string slice
	domainNames := make([]string, len(domainNamesInterface))
	for i, v := range domainNamesInterface {
		domainNames[i] = v.(string)
	}

	log.Printf("[INFO] Getting multimatch unsupported references for %d domain names\n", len(domainNames))

	// Call POST endpoint GetMultiMatchUnsupportedReferences
	domainPayload := applicationsegment.MultiMatchUnsupportedReferencesPayload(domainNames)
	unsupportedRefs, _, err := applicationsegment.GetMultiMatchUnsupportedReferences(ctx, service, domainPayload)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get multimatch unsupported references: %w", err))
	}

	// Flatten the response
	flattened := flattenMultiMatchUnsupportedReferences(unsupportedRefs)

	// Set a synthetic ID based on domain names
	syntheticID := generateSyntheticIDFromDomains(domainNames)
	d.SetId(syntheticID)

	_ = d.Set("unsupported_references", flattened)
	_ = d.Set("domain_names", domainNames)

	log.Printf("[INFO] Found %d unsupported multimatch references\n", len(unsupportedRefs))

	return nil
}

// flattenMultiMatchUnsupportedReferences flattens the response from GetMultiMatchUnsupportedReferences
func flattenMultiMatchUnsupportedReferences(refs []applicationsegment.MultiMatchUnsupportedReferencesResponse) []map[string]interface{} {
	if refs == nil {
		return nil
	}

	result := make([]map[string]interface{}, len(refs))
	for i, ref := range refs {
		result[i] = map[string]interface{}{
			"id":               ref.ID,
			"app_segment_name": ref.AppSegmentName,
			"domains":          ref.Domains,
			"tcp_ports":        ref.TCPPorts,
			"match_style":      ref.MatchStyle,
			"microtenant_name": ref.MicrotenantName,
		}
	}

	return result
}

// generateSyntheticIDFromDomains generates a deterministic ID from domain names
func generateSyntheticIDFromDomains(domainNames []string) string {
	// Sort domains for consistent hashing
	sorted := make([]string, len(domainNames))
	copy(sorted, domainNames)
	sort.Strings(sorted)

	// Create a unique string combining sorted domains
	idString := strings.Join(sorted, ",")

	// Generate MD5 hash for a shorter, deterministic ID
	hash := md5.Sum([]byte(idString))
	return fmt.Sprintf("%x", hash)
}
