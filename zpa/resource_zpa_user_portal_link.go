package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

func resourceUserPortalLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserPortalLinkCreate,
		ReadContext:   resourceUserPortalLinkRead,
		UpdateContext: resourceUserPortalLinkUpdate,
		DeleteContext: resourceUserPortalLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.Service

				microTenantID := GetString(d.Get("microtenant_id"))
				if microTenantID != "" {
					service = service.WithMicroTenant(microTenantID)
				}

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := portal_link.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the User Portal Link",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the User Portal Link",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this User Portal Link is enabled or not",
			},
			"icon_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Icon text for the User Portal Link",
			},
			"link": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Link URL for the User Portal Link",
			},
			"link_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Link path for the User Portal Link",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol for the User Portal Link",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Microtenant ID for the User Portal Link",
			},
			"user_portals": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of User Portals",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceUserPortalLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	req := expandUserPortalLink(d)
	log.Printf("[INFO] Creating zpa user portal link with request\n%+v\n", req)

	resp, _, err := portal_link.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created user portal link request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceUserPortalLinkRead(ctx, d, meta)
}

func resourceUserPortalLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := portal_link.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing user portal link %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting user portal link:\n%+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("icon_text", resp.IconText)
	_ = d.Set("link", resp.Link)
	_ = d.Set("link_path", resp.LinkPath)
	_ = d.Set("protocol", resp.Protocol)
	_ = d.Set("microtenant_id", resp.MicrotenantID)
	_ = d.Set("user_portals", flattenUserPortalsSimple(resp.UserPortals))
	return nil
}

func resourceUserPortalLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating user portal link ID: %v\n", id)

	req := expandUserPortalLink(d)
	log.Printf("[DEBUG] Expanding user portal link request: %+v", req)

	if _, _, err := portal_link.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := portal_link.Update(ctx, service, id, &req); err != nil {
		log.Printf("[ERROR] Failed to update portal link %s: %v", id, err)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully updated portal link %s", id)
	return resourceUserPortalLinkRead(ctx, d, meta)
}

func resourceUserPortalLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting user portal link with id %v\n", d.Id())

	if _, err := portal_link.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandUserPortalLink(d *schema.ResourceData) portal_link.UserPortalLink {
	return portal_link.UserPortalLink{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		IconText:      d.Get("icon_text").(string),
		Link:          d.Get("link").(string),
		LinkPath:      d.Get("link_path").(string),
		Protocol:      d.Get("protocol").(string),
		MicrotenantID: d.Get("microtenant_id").(string),
		UserPortals:   expandUserPortals(d),
	}
}

func expandUserPortals(d *schema.ResourceData) []portal_controller.UserPortalController {
	raw, ok := d.GetOk("user_portals")
	if !ok || raw == nil {
		log.Printf("[DEBUG] No user_portals configured, returning nil")
		return nil
	}

	blocks := raw.([]interface{})
	if len(blocks) == 0 {
		log.Printf("[DEBUG] Empty user_portals blocks, returning nil")
		return nil
	}

	block, ok := blocks[0].(map[string]interface{})
	if !ok {
		log.Printf("[DEBUG] Invalid user_portals block structure")
		return nil
	}

	idRaw, ok := block["id"]
	if !ok || idRaw == nil {
		log.Printf("[DEBUG] No user portal IDs found in configuration")
		return nil
	}

	idSet, ok := idRaw.(*schema.Set)
	if !ok {
		log.Printf("[DEBUG] Invalid user portal ID set structure")
		return nil
	}

	var portals []portal_controller.UserPortalController
	portalIDs := idSet.List()
	log.Printf("[DEBUG] Expanding %d user portals: %v", len(portalIDs), portalIDs)

	for _, id := range portalIDs {
		portalID := id.(string)
		log.Printf("[DEBUG] Adding user portal ID: %s", portalID)
		portals = append(portals, portal_controller.UserPortalController{
			ID: portalID,
		})
	}

	log.Printf("[DEBUG] Expanded %d user portals for portal link", len(portals))
	return portals
}

func flattenUserPortalsSimple(portals []portal_controller.UserPortalController) []interface{} {
	if len(portals) == 0 {
		return nil
	}

	ids := make([]interface{}, 0, len(portals))
	for _, portal := range portals {
		if portal.ID != "" {
			ids = append(ids, portal.ID)
		}
	}

	if len(ids) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id": schema.NewSet(schema.HashString, ids),
		},
	}
}
