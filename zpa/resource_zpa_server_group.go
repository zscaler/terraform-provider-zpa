package zscaler

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/servergroup"
)

func resourceServerGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceServerGroupCreate,
		Read:     resourceServerGroupRead,
		Update:   resourceServerGroupUpdate,
		Delete:   resourceServerGroupDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"SIEM",
				}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field is the description of the server group.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This field defines if the server group is enabled or disabled.",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dynamic_discovery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This field controls dynamic discovery of the servers.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This field defines the name of the server group.",
			},
			"servers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "This field is a list of servers that are applicable only when dynamic discovery is disabled. Server name is required only in cases where the new servers need to be created in this API. For existing servers, pass only the serverId.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"applications": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "This field is a json array of app-connector-id only.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"app_connector_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of app-connector IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
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

func resourceServerGroupCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandServerGroup(d)
	log.Printf("[INFO] Creating zpa server group with request\n%+v\n", req)
	if len(req.Servers) > 0 && req.DynamicDiscovery {
		log.Printf("[ERROR] An application server can only be attached to a server when DynamicDiscovery is disabled\n")
		return fmt.Errorf("an application server can only be attached to a server when DynamicDiscovery is disabled")
	}
	if !req.DynamicDiscovery && len(req.Servers) == 0 {
		log.Printf("[ERROR] Servers must not be empty when DynamicDiscovery is disabled\n")
		return fmt.Errorf("servers must not be empty when DynamicDiscovery is disabled")
	}
	resp, _, err := zClient.servergroup.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created server group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceServerGroupRead(d, m)
}

func resourceServerGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.servergroup.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing server group %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting server group:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("config_space", resp.ConfigSpace)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("ip_anchored", resp.IpAnchored)
	_ = d.Set("dynamic_discovery", resp.DynamicDiscovery)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("name", resp.Name)
	_ = d.Set("app_connector_groups", flattenAppConnectorGroups(resp.AppConnectorGroups))
	_ = d.Set("applications", flattenServerGroupApplications(resp.Applications))
	_ = d.Set("servers", flattenServers(resp.Servers))

	return nil

}

func resourceServerGroupUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id := d.Id()
	log.Printf("[INFO] Updating server group ID: %v\n", id)
	req := expandServerGroup(d)
	if (d.HasChange("servers") || d.HasChange("dynamic_discovery")) && req.DynamicDiscovery && len(req.Servers) > 0 {
		log.Printf("[ERROR] Can't update the server group: an application server can only be attached to a server when DynamicDiscovery is disabled\n")
		return fmt.Errorf("can't perform the changes: an application server can only be attached to a server when DynamicDiscovery is disabled")
	}
	if (d.HasChange("servers") || d.HasChange("dynamic_discovery")) && !req.DynamicDiscovery && len(req.Servers) == 0 {
		log.Printf("[ERROR] Can't update server group: servers must not be empty when DynamicDiscovery is disabled\n")
		return fmt.Errorf("can't update server group: servers must not be empty when DynamicDiscovery is disabled")
	}
	if _, err := zClient.servergroup.Update(id, &req); err != nil {
		return err
	}
	return resourceServerGroupRead(d, m)
}

func resourceServerGroupDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting server group ID: %v\n", d.Id())

	if _, err := zClient.servergroup.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] server group deleted")
	return nil
}

func expandServerGroup(d *schema.ResourceData) servergroup.ServerGroup {
	result := servergroup.ServerGroup{
		Enabled:            d.Get("enabled").(bool),
		Description:        d.Get("description").(string),
		IpAnchored:         d.Get("ip_anchored").(bool),
		ConfigSpace:        d.Get("config_space").(string),
		DynamicDiscovery:   d.Get("dynamic_discovery").(bool),
		AppConnectorGroups: expandAppConnectorGroups(d),
		Applications:       expandServerGroupApplications(d),
		Servers:            expandApplicationServers(d),
	}
	if d.HasChange("name") {
		result.Name = d.Get("name").(string)
	}
	return result

}

func expandAppConnectorGroups(d *schema.ResourceData) []servergroup.AppConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("app_connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app connector groups data: %+v\n", appConnector)
		var appConnectorGroups []servergroup.AppConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, ok := appConnectorGroup.(map[string]interface{})
			if ok {
				for _, id := range appConnectorGroup["id"].([]interface{}) {
					appConnectorGroups = append(appConnectorGroups, servergroup.AppConnectorGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appConnectorGroups
	}

	return []servergroup.AppConnectorGroups{}
}

func expandServerGroupApplications(d *schema.ResourceData) []servergroup.Applications {
	serverGroupAppsInterface, ok := d.GetOk("applications")
	if ok {
		serverGroupApp := serverGroupAppsInterface.(*schema.Set)
		log.Printf("[INFO] server group application data: %+v\n", serverGroupApp)
		var serverGroupApps []servergroup.Applications
		for _, serverGroupApp := range serverGroupApp.List() {
			serverGroupApp, ok := serverGroupApp.(map[string]interface{})
			if ok {
				for _, id := range serverGroupApp["id"].([]interface{}) {
					serverGroupApps = append(serverGroupApps, servergroup.Applications{
						ID: id.(string),
					})
				}
			}
		}
		return serverGroupApps
	}

	return []servergroup.Applications{}
}

func expandApplicationServers(d *schema.ResourceData) []servergroup.ApplicationServer {
	applicationServersInterface, ok := d.GetOk("servers")
	if ok {
		applicationServer := applicationServersInterface.(*schema.Set)
		log.Printf("[INFO] server group application data: %+v\n", applicationServer)
		var applicationServers []servergroup.ApplicationServer
		for _, applicationServer := range applicationServer.List() {
			applicationServer, _ := applicationServer.(map[string]interface{})
			if applicationServer != nil {
				for _, id := range applicationServer["id"].([]interface{}) {
					applicationServers = append(applicationServers, servergroup.ApplicationServer{
						ID: id.(string),
					})
				}
			}
		}
		return applicationServers
	}

	return []servergroup.ApplicationServer{}
}
