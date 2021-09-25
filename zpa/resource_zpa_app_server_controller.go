package zpa

import (
	"log"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/appservercontroller"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
)

func resourceApplicationServer() *schema.Resource {
	return &schema.Resource{
		Create:   resourceApplicationServerCreate,
		Read:     resourceApplicationServerRead,
		Update:   resourceApplicationServerUpdate,
		Delete:   resourceApplicationServerDelete,
		Importer: &schema.ResourceImporter{},

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
				Description: "This field defines the description of the server.",
			},
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "This field defines the domain or IP address of the server.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This field defines the status of the server.",
			},
			// App Server Group ID can only be attached if Dynamic Server Discovery in Server Group is False
			"app_server_group_ids": {
				Type:        schema.TypeSet,
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
		if err.(*client.ErrorResponse).IsObjectNotFound() {
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

	if d.HasChange("app_server_group_ids") || d.HasChange("name") || d.HasChange("address") {
		log.Println("The AppServerGroupID, name or address has been changed")

		if _, err := zClient.appservercontroller.Update(d.Id(), appservercontroller.ApplicationServer{
			AppServerGroupIds: SetToStringSlice(d.Get("app_server_group_ids").(*schema.Set)),
			Name:              d.Get("name").(string),
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
		return err
	}

	return nil
}

func removeServerFromGroup(zClient *Client, serverID string) error {
	// Remove the reference to this server from server groups.

	resp, _, err := zClient.appservercontroller.Get(serverID)
	if err != nil {
		return err
	}

	if len(resp.AppServerGroupIds) != 0 {
		log.Printf("[INFO] Removing server group ID/s from application server: %s", serverID)
		resp.AppServerGroupIds = make([]string, 0)

		log.Printf("[INFO] Updating server group ID: %s", serverID)
		_, err = zClient.appservercontroller.Update(serverID, *resp)
		if err != nil {
			log.Printf("[ERROR] Failed to update application server ID: %s", serverID)
			return err
		}
	}

	return nil
}

func expandCreateAppServerRequest(d *schema.ResourceData) appservercontroller.ApplicationServer {
	applicationServer := appservercontroller.ApplicationServer{
		Address:           d.Get("address").(string),
		ConfigSpace:       d.Get("config_space").(string),
		AppServerGroupIds: SetToStringSlice(d.Get("app_server_group_ids").(*schema.Set)),
		Description:       d.Get("description").(string),
		Enabled:           d.Get("enabled").(bool),
		Name:              d.Get("name").(string),
	}
	return applicationServer
}
