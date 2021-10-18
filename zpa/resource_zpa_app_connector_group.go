package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appservercontroller"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
)

func resourceAppConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAppConnectorGroupCreate,
		Read:     resourceAppConnectorGroupRead,
		Update:   resourceAppConnectorGroupUpdate,
		Delete:   resourceAppConnectorGroupDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city_country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_query_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lss_app_connector_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"upgrade_day": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"upgrade_time_in_secs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override_version_profile": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"version_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAppConnectorGroupCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandAppConnectorGroup(d)
	log.Printf("[INFO] Creating zpa app connector group with request\n%+v\n", req)

	resp, _, err := zClient.appservercontroller.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created app connector group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceAppConnectorGroupRead(d, m)
}

func resourceAppConnectorGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appconnectorgroup.Create(req)
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

func resourceAppConnectorGroupUpdate(d *schema.ResourceData, m interface{}) error {
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

func resourceAppConnectorGroupDelete(d *schema.ResourceData, m interface{}) error {
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
