package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appservercontroller"
)

func resourceApplicationServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationServerCreate,
		Read:   resourceApplicationServerRead,
		Update: resourceApplicationServerUpdate,
		Delete: resourceApplicationServerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.appservercontroller.GetByName(id)
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "This field defines the name of the server.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This field defines the description of the server.",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This field defines the domain or IP address of the server.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "This field defines the status of the server.",
			},
			// App Server Group ID can only be attached if Dynamic Server Discovery in Server Group is False
			"app_server_group_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This field defines the list of server groups IDs.",
				Optional:    true,
			},
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"SIEM",
				}, false),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationServerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandCreateAppServerRequest(d)
	log.Printf("[INFO] Creating zpa application server with request\n%+v\n", req)

	resp, _, err := zClient.appservercontroller.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created application server request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceApplicationServerRead(d, m)
}

func resourceApplicationServerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appservercontroller.Get(d.Id())
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing application server %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
	_ = d.Set("address", resp.Address)
	_ = d.Set("app_server_group_ids", resp.AppServerGroupIds)
	_ = d.Set("config_space", resp.ConfigSpace)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	return nil

}

func resourceApplicationServerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Println("An updated occurred")

	if d.HasChange("app_server_group_ids") || d.HasChange("name") || d.HasChange("description") || d.HasChange("address") || d.HasChange("enabled") {
		log.Println("The AppServerGroupID, name, description or address has been changed")

		if _, _, err := zClient.appservercontroller.Get(d.Id()); err != nil {
			if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
				d.SetId("")
				return nil
			}
		}

		if _, err := zClient.appservercontroller.Update(d.Id(), &appservercontroller.ApplicationServer{
			AppServerGroupIds: SetToStringSlice(d.Get("app_server_group_ids").(*schema.Set)),
			Name:              d.Get("name").(string),
			Description:       d.Get("description").(string),
			Address:           d.Get("address").(string),
			Enabled:           d.Get("enabled").(bool),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceApplicationServerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting application server ID: %v\n", d.Id())

	err := removeServerFromGroup(zClient, d.Id())
	if err != nil {
		return err
	}

	if _, err = zClient.appservercontroller.Delete(d.Id()); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	return nil
}

func removeServerFromGroup(zClient *Client, serverID string) error {
	// Remove the reference to this server from server groups.

	resp, _, err := zClient.appservercontroller.Get(serverID)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			return nil
		}
		return err
	}

	if len(resp.AppServerGroupIds) != 0 {
		log.Printf("[INFO] Removing server group ID/s from application server: %s", serverID)
		resp.AppServerGroupIds = make([]string, 0)

		log.Printf("[INFO] Updating server group ID: %s", serverID)
		_, err = zClient.appservercontroller.Update(serverID, resp)
		if err != nil {
			log.Printf("[ERROR] Failed to update application server ID: %s", serverID)
			return err
		}
	}

	return nil
}

func expandCreateAppServerRequest(d *schema.ResourceData) appservercontroller.ApplicationServer {
	applicationServer := appservercontroller.ApplicationServer{
		ID:                d.Id(),
		Address:           d.Get("address").(string),
		ConfigSpace:       d.Get("config_space").(string),
		AppServerGroupIds: SetToStringSlice(d.Get("app_server_group_ids").(*schema.Set)),
		Description:       d.Get("description").(string),
		Enabled:           d.Get("enabled").(bool),
		Name:              d.Get("name").(string),
	}
	return applicationServer
}
