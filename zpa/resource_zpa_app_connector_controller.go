package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/appconnectorcontroller"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/client"
)

func resourceAppConnectorController() *schema.Resource {
	return &schema.Resource{
		// Create: resourceAppConnectorBulkeDelete,
		Create: resourceAppConnectorControllerCreate,
		Read:   resourceAppConnectorControllerRead,
		Update: resourceAppConnectorControllerUpdate,
		Delete: resourceAppConnectorControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.appconnectorcontroller.GetByName(id)
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
			// "ids": {
			// 	Type:        schema.TypeSet,
			// 	Computed:    true,
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// 	Description: "This field defines the list of app connector ids IDs.",
			// 	Optional:    true,
			// },
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the App Connector",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Description of the App Connector",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether this App Connector is enabled or not",
			},
		},
	}
}

/*
func resourceAppConnectorBulkeDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandAppConnectorController(d)
	log.Printf("[INFO] Creating zpa app connector with request\n%+v\n", req)

	resp, _, err := zClient.appconnectorcontroller.BulkDelete(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created app connector request. ID: %v\n", resp)
	d.SetId(resp.IDs)

	return resourceAppConnectorControllerRead(d, m)
}
*/

func resourceAppConnectorControllerCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAppConnectorControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appconnectorcontroller.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing app connector group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	return nil

}

func resourceAppConnectorControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating app connector group ID: %v\n", id)
	req := expandAppConnectorController(d)

	if _, err := zClient.appconnectorcontroller.Update(id, &req); err != nil {
		return err
	}

	return resourceAppConnectorControllerRead(d, m)
}

func resourceAppConnectorControllerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting app connector groupID: %v\n", d.Id())

	if _, err := zClient.appconnectorcontroller.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] app connector group deleted")
	return nil
}

func expandAppConnectorController(d *schema.ResourceData) appconnectorcontroller.AppConnector {
	AppConnectorController := appconnectorcontroller.AppConnector{
		ID: d.Get("id").(string),
		// IDs:         d.Get("ids").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
	}
	return AppConnectorController
}
