package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praconsole"
)

func resourcePRAConsoleController() *schema.Resource {
	return &schema.Resource{
		Create: resourcePRAConsoleControllerCreate,
		Read:   resourcePRAConsoleControllerRead,
		Update: resourcePRAConsoleControllerUpdate,
		Delete: resourcePRAConsoleControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := service.GetByName(id)
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the privileged console",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the privileged console",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the privileged console",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether or not the privileged console is enabled",
			},
			"icon_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The privileged console icon. The icon image is converted to base64 encoded text format",
			},
			"pra_portals": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The unique identifier of the privileged portal",
						},
					},
				},
			},
			"pra_application": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier of the Privileged Remote Access-enabled application",
						},
					},
				},
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant.",
			},
		},
	}
}

func resourcePRAConsoleControllerCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	req := expandPRAConsole(d)
	log.Printf("[INFO] Creating pra console with request\n%+v\n", req)

	praConsole, _, err := service.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created pra console request. ID: %v\n", praConsole)

	d.SetId(praConsole.ID)
	return resourcePRAConsoleControllerRead(d, m)
}

func resourcePRAConsoleControllerRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	resp, _, err := service.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing pra console %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting pra console controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("icon_text", resp.IconText)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("pra_application", flattenPRAApplicationIDSimple(resp.PRAApplication))
	_ = d.Set("pra_portals", flattenPRAPortalIDSimple(resp.PRAPortals))

	return nil
}

func resourcePRAConsoleControllerUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	id := d.Id()
	log.Printf("[INFO] Updating pra console ID: %v\n", id)
	req := expandPRAConsole(d)

	if _, _, err := service.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := service.Update(id, &req); err != nil {
		return err
	}

	return resourcePRAConsoleControllerRead(d, m)
}

func resourcePRAConsoleControllerDelete(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	log.Printf("[INFO] Deleting pra console ID: %v\n", d.Id())

	if _, err := service.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] pra console deleted")
	return nil
}

func expandPRAConsole(d *schema.ResourceData) praconsole.PRAConsole {
	result := praconsole.PRAConsole{
		ID:            d.Id(),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		IconText:      d.Get("icon_text").(string),
		MicroTenantID: d.Get("microtenant_id").(string),
		PRAPortals:    expandPRAPortal(d),
	}
	application := expandPRAApplication(d)
	if application != nil {
		result.PRAApplication = *application // TODO: Need to fix pointer to PRAApplication Struct
	}
	return result
}

func expandPRAApplication(d *schema.ResourceData) *praconsole.PRAApplication {
	if v, ok := d.GetOk("pra_application"); ok {
		applicationList := v.([]interface{})
		if len(applicationList) > 0 {
			firstApplication, ok := applicationList[0].(map[string]interface{})
			if !ok || firstApplication == nil {
				// The first application is not a valid map, handle accordingly
				return nil
			}
			id, ok := firstApplication["id"].(string)
			if !ok || id == "" {
				// The ID is not a string or is empty, handle accordingly
				return nil
			}
			return &praconsole.PRAApplication{
				ID: id,
			}
		}
	}
	return nil // No pra_application found or ID is empty
}

func expandPRAPortal(d *schema.ResourceData) []praconsole.PRAPortals {
	praPortalInterface, ok := d.GetOk("pra_portals")
	if ok {
		praPortal := praPortalInterface.(*schema.Set)
		log.Printf("[INFO] pra portal data: %+v\n", praPortal)
		var praPortals []praconsole.PRAPortals
		for _, praPortal := range praPortal.List() {
			praPortal, ok := praPortal.(map[string]interface{})
			if ok {
				for _, id := range praPortal["id"].(*schema.Set).List() {
					praPortals = append(praPortals, praconsole.PRAPortals{
						ID: id.(string),
					})
				}
			}
		}
		return praPortals
	}

	return []praconsole.PRAPortals{}
}

func flattenPRAPortalIDSimple(praPortals []praconsole.PRAPortals) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(praPortals))
	for i, portal := range praPortals {
		ids[i] = portal.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenPRAApplicationIDSimple(praApplication praconsole.PRAApplication) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["id"] = praApplication.ID
	return result
}
