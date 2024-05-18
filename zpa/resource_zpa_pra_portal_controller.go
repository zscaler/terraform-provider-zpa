package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praportal"
)

func resourcePRAPortalController() *schema.Resource {
	return &schema.Resource{
		Create: resourcePRAPortalControllerCreate,
		Read:   resourcePRAPortalControllerRead,
		Update: resourcePRAPortalControllerUpdate,
		Delete: resourcePRAPortalControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				service := m.(*Client).praportal.WithMicroTenant(GetString(d.Get("microtenant_id")))

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
				Description: "The unique identifier of the privileged portal",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the privileged portal",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the privileged portal",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not the privileged portal is enabled",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain of the privileged portal",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the certificate",
			},
			"user_notification": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The notification message displayed in the banner of the privileged portallink, if enabled",
			},
			"user_notification_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the Notification Banner is enabled (true) or disabled (false)",
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
		},
	}
}

func resourcePRAPortalControllerCreate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praportal.WithMicroTenant(GetString(d.Get("microtenant_id")))

	req := expandPRAPortalController(d)
	log.Printf("[INFO] Creating pra portal controller with request\n%+v\n", req)

	praPortal, _, err := service.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created pra portal controller request. ID: %v\n", praPortal)

	d.SetId(praPortal.ID)
	return resourcePRAPortalControllerRead(d, m)
}

func resourcePRAPortalControllerRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praportal.WithMicroTenant(GetString(d.Get("microtenant_id")))

	resp, _, err := service.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing pra portal controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting pra portal controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("domain", resp.Domain)
	_ = d.Set("certificate_id", resp.CertificateID)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("user_notification", resp.UserNotification)
	_ = d.Set("user_notification_enabled", resp.UserNotificationEnabled)
	return nil
}

func resourcePRAPortalControllerUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praportal.WithMicroTenant(GetString(d.Get("microtenant_id")))

	id := d.Id()
	log.Printf("[INFO] Updating pra portal controller ID: %v\n", id)
	req := expandPRAPortalController(d)

	if _, _, err := service.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := service.Update(id, &req); err != nil {
		return err
	}

	return resourcePRAPortalControllerRead(d, m)
}

func resourcePRAPortalControllerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	portalID := d.Id()

	// Detach the portal from any consoles before attempting to delete it.
	consoleService := client.praconsole.WithMicroTenant(GetString(d.Get("microtenant_id")))
	if err := detachAndCleanUpPRAPortals(portalID, consoleService); err != nil {
		return fmt.Errorf("error detaching PRAPortal with ID %s from PRAConsoleControllers: %s", portalID, err)
	}

	// Proceed with deletion of the portal after successful detachment.
	service := client.praportal.WithMicroTenant(GetString(d.Get("microtenant_id")))
	log.Printf("[INFO] Deleting PRA Portal Controller with ID: %s", portalID)
	if _, err := service.Delete(portalID); err != nil {
		return fmt.Errorf("error deleting PRA Portal Controller with ID %s: %s", portalID, err)
	}

	log.Printf("[INFO] PRA Portal Controller with ID %s deleted", portalID)
	d.SetId("") // Indicate that the resource was successfully deleted.
	return nil
}

func expandPRAPortalController(d *schema.ResourceData) praportal.PRAPortal {
	praPortal := praportal.PRAPortal{
		ID:                      d.Id(),
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Enabled:                 d.Get("enabled").(bool),
		Domain:                  d.Get("domain").(string),
		CertificateID:           d.Get("certificate_id").(string),
		MicroTenantID:           d.Get("microtenant_id").(string),
		UserNotification:        d.Get("user_notification").(string),
		UserNotificationEnabled: d.Get("user_notification_enabled").(bool),
	}
	return praPortal
}

// Detach and optionally delete PRAPortalControllers from PRAConsoleControllers.
func detachAndCleanUpPRAPortals(portalID string, consoleService *praconsole.Service) error {
	// Fetch all PRAConsoleControllers
	consoles, _, err := consoleService.GetAll()
	if err != nil {
		return fmt.Errorf("failed to list all PRAConsoleControllers: %s", err)
	}

	for _, console := range consoles {
		// Identify if the current console is associated with the portalID
		var portalFound bool
		for _, portal := range console.PRAPortals {
			if portal.ID == portalID {
				portalFound = true
				break // Found the portal in this console
			}
		}

		// Proceed if the portal is found within the console
		if portalFound {
			// Remove the portal from the console's portal list
			updatedPortals := []praconsole.PRAPortals{}
			for _, portal := range console.PRAPortals {
				if portal.ID != portalID {
					updatedPortals = append(updatedPortals, portal)
				}
			}

			if len(updatedPortals) == 0 {
				// Delete the console if it no longer contains any portals
				_, err = consoleService.Delete(console.ID)
				if err != nil {
					return fmt.Errorf("failed to delete PRAConsoleController with ID %s: %s", console.ID, err)
				}
				log.Printf("[INFO] Deleted PRAConsoleController with ID %s as it no longer has any associated PRAPortals", console.ID)
			} else {
				// Update the console with the remaining portals
				console.PRAPortals = updatedPortals
				_, err = consoleService.Update(console.ID, &console)
				if err != nil {
					return fmt.Errorf("failed to update PRAConsoleController with ID %s: %s", console.ID, err)
				}
				log.Printf("[INFO] Updated PRAConsoleController with ID %s after detaching PRAPortal with ID %s", console.ID, portalID)
			}
		}
	}

	return nil
}
