package zpa

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

func resourceApplicationSegmentMultimatchBulk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationSegmentMultimatchBulkCreate,
		ReadContext:   resourceApplicationSegmentMultimatchBulkRead,
		UpdateContext: resourceApplicationSegmentMultimatchBulkUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"application_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of application segment IDs to update match_style for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"match_style": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Match style to apply to all specified application segments. Valid values: EXCLUSIVE, INCLUSIVE.",
				ValidateFunc: validation.StringInSlice([]string{
					"EXCLUSIVE",
					"INCLUSIVE",
				}, false),
			},
		},
	}
}

func resourceApplicationSegmentMultimatchBulkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationIDsStr := SetToStringList(d, "application_ids")
	if len(applicationIDsStr) == 0 {
		return diag.FromErr(fmt.Errorf("at least one application_id must be provided"))
	}

	// Convert string IDs to integers for the API payload
	applicationIDs, err := convertStringIDsToInts(applicationIDsStr)
	if err != nil {
		return diag.FromErr(err)
	}

	matchStyle := d.Get("match_style").(string)

	// Build payload for bulk update
	payload := applicationsegment.BulkUpdateMultiMatchPayload{
		ApplicationIDs: applicationIDs,
		MatchStyle:     matchStyle,
	}

	log.Printf("[INFO] Creating bulk multimatch update for %d application segments with match_style: %s\n", len(applicationIDs), matchStyle)
	_, err = applicationsegment.UpdatebulkUpdateMultiMatch(ctx, service, payload)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to bulk update multimatch: %w", err))
	}

	// Generate a synthetic ID based on application_ids (as strings) and match_style
	id := generateBulkMultimatchID(applicationIDsStr, matchStyle)
	d.SetId(id)

	log.Printf("[INFO] Successfully created bulk multimatch update. ID: %s\n", id)

	return resourceApplicationSegmentMultimatchBulkRead(ctx, d, meta)
}

func resourceApplicationSegmentMultimatchBulkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationIDs := SetToStringList(d, "application_ids")
	if len(applicationIDs) == 0 {
		// If no application_ids in state, try to reconstruct from ID
		log.Printf("[WARN] No application_ids found in state, resource may need to be recreated")
		d.SetId("")
		return nil
	}

	// Collect all domain names from the application segments
	allDomainNames := make([]string, 0)
	segmentDomainMap := make(map[string][]string) // segment ID -> domain names

	// Fetch each segment to get domain names and current match_style
	for _, appID := range applicationIDs {
		segment, _, err := applicationsegment.Get(ctx, service, appID)
		if err != nil {
			log.Printf("[WARN] Failed to fetch application segment %s: %v", appID, err)
			continue
		}
		segmentDomainMap[appID] = segment.DomainNames
		allDomainNames = append(allDomainNames, segment.DomainNames...)
	}

	// Use POST endpoint GetMultiMatchUnsupportedReferences to maintain state
	// This endpoint returns segments with their current match_style
	matchStyleFromState := ""
	if v, ok := d.GetOk("match_style"); ok {
		matchStyleFromState = v.(string)
	}
	actualMatchStyle := matchStyleFromState

	if len(allDomainNames) > 0 {
		domainPayload := applicationsegment.MultiMatchUnsupportedReferencesPayload(allDomainNames)
		unsupportedRefs, _, err := applicationsegment.GetMultiMatchUnsupportedReferences(ctx, service, domainPayload)
		if err != nil {
			log.Printf("[WARN] Failed to get unsupported references via POST: %v", err)
			// Fallback to individual GET calls
			actualMatchStyle = verifyMatchStyleViaIndividualGETs(ctx, service, applicationIDs, matchStyleFromState)
		} else {
			// Build a map of segment ID to match_style from the POST response
			segmentMatchStyleMap := make(map[string]string)
			for _, ref := range unsupportedRefs {
				if ref.ID != "" && ref.MatchStyle != "" {
					segmentMatchStyleMap[ref.ID] = ref.MatchStyle
				}
			}

			// Determine the actual match_style from the POST response
			// Check if all segments in our list have the same match_style
			for _, appID := range applicationIDs {
				if refMatchStyle, exists := segmentMatchStyleMap[appID]; exists {
					if actualMatchStyle == "" {
						actualMatchStyle = refMatchStyle
					} else if actualMatchStyle != refMatchStyle {
						log.Printf("[WARN] Application segments have different match_style values. Segment %s has '%s', but others have '%s'", appID, refMatchStyle, actualMatchStyle)
					}
				}
			}

			// If we didn't find match_style in POST response, fallback to individual GETs
			if actualMatchStyle == "" {
				actualMatchStyle = verifyMatchStyleViaIndividualGETs(ctx, service, applicationIDs, matchStyleFromState)
			}
		}
	} else {
		// No domain names available, use individual GET calls
		actualMatchStyle = verifyMatchStyleViaIndividualGETs(ctx, service, applicationIDs, matchStyleFromState)
	}

	// Set state with the actual match_style found
	if actualMatchStyle != "" {
		_ = d.Set("match_style", actualMatchStyle)
	}
	_ = d.Set("application_ids", applicationIDs)

	return nil
}

func resourceApplicationSegmentMultimatchBulkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	applicationIDsStr := SetToStringList(d, "application_ids")
	if len(applicationIDsStr) == 0 {
		return diag.FromErr(fmt.Errorf("at least one application_id must be provided"))
	}

	// Convert string IDs to integers for the API payload
	applicationIDs, err := convertStringIDsToInts(applicationIDsStr)
	if err != nil {
		return diag.FromErr(err)
	}

	matchStyle := d.Get("match_style").(string)

	// Build payload for bulk update
	payload := applicationsegment.BulkUpdateMultiMatchPayload{
		ApplicationIDs: applicationIDs,
		MatchStyle:     matchStyle,
	}

	log.Printf("[INFO] Updating bulk multimatch for %d application segments with match_style: %s\n", len(applicationIDs), matchStyle)
	_, err = applicationsegment.UpdatebulkUpdateMultiMatch(ctx, service, payload)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to bulk update multimatch: %w", err))
	}

	// Update ID if application_ids or match_style changed
	newID := generateBulkMultimatchID(applicationIDsStr, matchStyle)
	if d.Id() != newID {
		d.SetId(newID)
	}

	log.Printf("[INFO] Successfully updated bulk multimatch. ID: %s\n", newID)

	return resourceApplicationSegmentMultimatchBulkRead(ctx, d, meta)
}

// verifyMatchStyleViaIndividualGETs fetches each segment individually to verify match_style
func verifyMatchStyleViaIndividualGETs(ctx context.Context, service *zscaler.Service, applicationIDs []string, expectedMatchStyle string) string {
	actualMatchStyle := ""
	for _, appID := range applicationIDs {
		segment, _, err := applicationsegment.Get(ctx, service, appID)
		if err != nil {
			log.Printf("[WARN] Failed to fetch application segment %s for match_style verification: %v", appID, err)
			continue
		}
		if actualMatchStyle == "" {
			actualMatchStyle = segment.MatchStyle
		}
		if segment.MatchStyle != expectedMatchStyle && expectedMatchStyle != "" {
			log.Printf("[WARN] Application segment %s has match_style '%s' (expected '%s')", appID, segment.MatchStyle, expectedMatchStyle)
		}
	}
	return actualMatchStyle
}

// convertStringIDsToInts converts a slice of string IDs to integers
func convertStringIDsToInts(applicationIDsStr []string) ([]int, error) {
	applicationIDs := make([]int, len(applicationIDsStr))
	for i, idStr := range applicationIDsStr {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert application_id '%s' to integer: %w", idStr, err)
		}
		applicationIDs[i] = idInt
	}
	return applicationIDs, nil
}

// generateBulkMultimatchID generates a synthetic ID based on application_ids and match_style
func generateBulkMultimatchID(applicationIDs []string, matchStyle string) string {
	// Sort IDs for consistent hashing
	sortedIDs := make([]string, len(applicationIDs))
	copy(sortedIDs, applicationIDs)
	sort.Strings(sortedIDs)

	// Create a unique string combining sorted IDs and match_style
	idString := fmt.Sprintf("%s|%s", strings.Join(sortedIDs, ","), matchStyle)

	// Generate MD5 hash for a shorter, deterministic ID
	hash := md5.Sum([]byte(idString))
	return fmt.Sprintf("%x", hash)
}
