package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praconsole"
)

func dataSourcePRAConsoleController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePRAConsoleControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the privileged console",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the privileged console",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the privileged console",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the privileged console is enabled",
			},
			"icon_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The privileged console icon. The icon image is converted to base64 encoded text format",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged console is created",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The tenant who modified the privileged console",
			},
			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged console is modified",
			},
			"pra_application": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the Privileged Remote Access-enabled application",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Privileged Remote Access-enabled application",
						},
					},
				},
			},
			"pra_portals": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the privileged portal",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the privileged portal",
						},
					},
				},
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
			"microtenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Microtenant",
			},
		},
	}
}

func dataSourcePRAConsoleControllerRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PRAConsole

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	var resp *praconsole.PRAConsole
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for pra console controller %s\n", id)
		res, _, err := praconsole.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for sra console controller name %s\n", name)
		res, _, err := praconsole.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("icon_text", resp.IconText)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)

		if err := d.Set("pra_portals", flattenPRAPortalIDName(resp.PRAPortals)); err != nil {
			return fmt.Errorf("failed to read pra portals %s", err)
		}
		if err := d.Set("pra_application", flattenPRAApplicationIDName(resp.PRAApplication)); err != nil {
			return fmt.Errorf("failed to read pra applications %s", err)
		}

	} else {
		return fmt.Errorf("couldn't find any sra privileged approval with id '%s'", id)
	}

	return nil
}

func flattenPRAApplicationIDName(application praconsole.PRAApplication) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["id"] = application.ID
	result[0]["name"] = application.Name
	return result
}

func flattenPRAPortalIDName(sraPortal []praconsole.PRAPortals) []interface{} {
	applicationServers := make([]interface{}, len(sraPortal))
	for i, portalItem := range sraPortal {
		applicationServers[i] = map[string]interface{}{
			"id":   portalItem.ID,
			"name": portalItem.Name,
		}
	}
	return applicationServers
}
