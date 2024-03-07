package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praconsole"
)

func resourcePRAConsoleControllerBulkBulk() *schema.Resource {
	return &schema.Resource{
		Create: resourcePRAConsoleControllerBulkCreate,
		Read:   resourcePRAConsoleControllerBulkRead,
		Update: resourcePRAConsoleControllerBulkUpdate,
		Delete: resourcePRAConsoleControllerBulkDelete,
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
			"pra_consoles": {
				Type:     schema.TypeList, // Use TypeList if the order is important; otherwise, use TypeSet
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the privileged console",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the privileged console",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the privileged console is enabled",
						},
						"icon_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The privileged console icon text",
						},
						"pra_application": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString, // Ensure this matches the data type expected by your API
										Required:    true,
										Description: "The unique identifier of the Privileged Remote Access-enabled application",
									},
								},
							},
						},
						"pra_portals": {
							Type:     schema.TypeList, // Changed from TypeSet to TypeList
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString, // Corrected to TypeString
										Required: true,
									},
								},
							},
						},
						"microtenant_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique identifier of the Microtenant for the ZPA tenant",
						},
					},
				},
			},
		},
	}
}

func resourcePRAConsoleControllerBulkCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	praConsoles, err := expandPRAConsolesBulk(d)
	if err != nil {
		return err
	}

	praConsolesCreated, _, err := service.CreatePraBulk(praConsoles)
	if err != nil {
		return err
	}

	// Example: setting to the first console's ID. Adjust as needed.
	if len(praConsolesCreated) > 0 {
		d.SetId(praConsolesCreated[0].ID)
	}

	return resourcePRAConsoleControllerBulkRead(d, m)
}

func resourcePRAConsoleControllerBulkRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Fetch all console data. The assumption here is GetAll() returns a slice of praconsole.PRAConsole.
	praConsoles, _, err := service.GetAll() // Correctly use the returned slice.
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] PRA consoles resource not found, removing from state")
			d.SetId("")
			return nil
		}
		return err
	}

	// Directly map the slice of PRAConsole structs to the Terraform schema format for pra_consoles.
	var consoles []interface{}
	for _, item := range praConsoles {
		consoleMap := make(map[string]interface{})
		consoleMap["name"] = item.Name
		consoleMap["description"] = item.Description
		consoleMap["enabled"] = item.Enabled
		consoleMap["icon_text"] = item.IconText

		// Map pra_application as a list with a single map item.
		praApplication := []interface{}{
			map[string]interface{}{
				"id": item.PRAApplication.ID,
			},
		}
		consoleMap["pra_application"] = praApplication

		// Map pra_portals as a list of map items.
		var portals []interface{}
		for _, portal := range item.PRAPortals {
			portals = append(portals, map[string]interface{}{
				"id": portal.ID,
			})
		}
		consoleMap["pra_portals"] = portals

		consoles = append(consoles, consoleMap)
	}

	// Set the pra_consoles attribute with the constructed slice of maps.
	if err := d.Set("pra_consoles", consoles); err != nil {
		return fmt.Errorf("failed to set 'pra_consoles': %s", err)
	}

	return nil
}

func resourcePRAConsoleControllerBulkUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Assuming the ID of the PRA console to update is stored in the resource ID.
	id := d.Id()

	// Expanding the single console that needs to be updated.
	var praConsole praconsole.PRAConsole
	if praConsoles, err := expandPRAConsolesBulk(d); err != nil {
		return err
	} else {
		// Find the console to update based on ID or another unique identifier.
		// This example assumes you have the logic to find the specific console from the slice.
		// If the ID directly maps to a single console, adjust as necessary.
		for _, console := range praConsoles {
			if console.ID == id {
				praConsole = console
				break
			}
		}
	}

	if _, _, err := service.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	// Performing the update. Adjust the method call as necessary based on your SDK.
	_, err := service.Update(id, &praConsole)
	if err != nil {
		return fmt.Errorf("failed to update PRA console: %s", err)
	}

	// After updating, read the console again to refresh the state.
	return resourcePRAConsoleControllerBulkRead(d, m)
}

func resourcePRAConsoleControllerBulkDelete(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))

	log.Printf("[INFO] Deleting pra console ID: %v\n", d.Id())

	if _, err := service.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] pra console deleted")
	return nil
}

func expandPRAConsolesBulk(d *schema.ResourceData) ([]praconsole.PRAConsole, error) {
	var praConsoles []praconsole.PRAConsole

	if v, ok := d.GetOk("pra_consoles"); ok {
		consolesList := v.([]interface{})

		for _, consoleItem := range consolesList {
			consoleMap := consoleItem.(map[string]interface{})

			praConsole := praconsole.PRAConsole{
				Name:          consoleMap["name"].(string),
				Description:   consoleMap["description"].(string),
				Enabled:       consoleMap["enabled"].(bool),
				IconText:      consoleMap["icon_text"].(string),
				MicroTenantID: consoleMap["microtenant_id"].(string),
			}

			// Adjusted processing of pra_application
			if praAppSlice, exists := consoleMap["pra_application"].([]interface{}); exists && len(praAppSlice) > 0 {
				if praAppMap, valid := praAppSlice[0].(map[string]interface{}); valid {
					if appID, validID := praAppMap["id"].(string); validID && appID != "" {
						praConsole.PRAApplication = praconsole.PRAApplication{ID: appID}
					} else {
						// Instead of immediate failure, log and continue to allow for potential partial success or retry
						log.Printf("Warning: pra_application ID is missing or invalid for console")
					}
				} else {
					log.Printf("Warning: pra_application format is invalid for console")
				}
			} else {
				log.Printf("Warning: pra_application is missing for console")
			}

			// Processing pra_portals, ensuring we retain the original handling
			praPortals := make([]praconsole.PRAPortals, 0)
			if portals, ok := consoleMap["pra_portals"].([]interface{}); ok {
				for _, portal := range portals {
					portalMap := portal.(map[string]interface{})
					portalID := portalMap["id"].(string)
					praPortals = append(praPortals, praconsole.PRAPortals{ID: portalID})
				}
			}
			praConsole.PRAPortals = praPortals

			praConsoles = append(praConsoles, praConsole)
		}
	}

	return praConsoles, nil
}
