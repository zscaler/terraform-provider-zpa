package zpa

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
)

var detachLock sync.Mutex

func resourceServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerGroupCreate,
		Read:   resourceServerGroupRead,
		Update: resourceServerGroupUpdate,
		Delete: resourceServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.ServerGroup

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
					resp, _, err := servergroup.GetByName(service, id)
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
			"config_space": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"SIEM",
				}, false),
				Default: "DEFAULT",
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"servers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
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
				Computed:    true,
				Description: "List of app-connector IDs.",
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

func resourceServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ServerGroup

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

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
	resp, _, err := servergroup.Create(service, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created server group request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceServerGroupRead(d, meta)
}

func resourceServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ServerGroup

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := servergroup.Get(service, d.Id())
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
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("app_connector_groups", flattenCommonAppConnectorGroups(resp.AppConnectorGroups))
	_ = d.Set("applications", flattenServerGroupApplicationsSimple(resp.Applications))
	_ = d.Set("servers", flattenServers(resp.Servers))

	return nil
}

func resourceServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ServerGroup

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating server group ID: %v\n", id)
	req := expandServerGroup(d)

	if d.HasChanges("servers", "dynamic_discovery") {
		if req.DynamicDiscovery && len(req.Servers) > 0 {
			log.Printf("[ERROR] Can't update the server group: an application server can only be attached to a server when DynamicDiscovery is disabled\n")
			return fmt.Errorf("can't perform the changes: an application server can only be attached to a server when DynamicDiscovery is disabled")
		}
		if !req.DynamicDiscovery && len(req.Servers) == 0 {
			log.Printf("[ERROR] Can't update server group: servers must not be empty when DynamicDiscovery is disabled\n")
			return fmt.Errorf("can't update server group: servers must not be empty when DynamicDiscovery is disabled")
		}
	}

	if _, _, err := servergroup.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := servergroup.Update(service, id, &req); err != nil {
		return err
	}
	return resourceServerGroupRead(d, meta)
}

func resourceServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	applicationSegmentService := zClient.ApplicationSegment
	policySetControllerService := zClient.PolicySetController
	appConnectorGroupService := zClient.AppConnectorGroup
	service := zClient.ServerGroup

	if microTenantID != "" {
		applicationSegmentService = applicationSegmentService.WithMicroTenant(microTenantID)
		policySetControllerService = policySetControllerService.WithMicroTenant(microTenantID)
		appConnectorGroupService = appConnectorGroupService.WithMicroTenant(microTenantID)
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting server group ID: %v\n", d.Id())

	err := detachServerGroupFromAppConnectorGroups(d.Id(), service, appConnectorGroupService)
	if err != nil {
		log.Printf("[ERROR] Detaching server group ID: %v from app connector groups failed: %v\n", d.Id(), err)
	}

	detachServerGroupFromAllAccessPolicyRules(d.Id(), policySetControllerService)
	detachServerGroupFromAllAppSegments(d.Id(), applicationSegmentService)

	if _, err := servergroup.Delete(service, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] Server group deleted")
	return nil
}

func detachServerGroupFromAllAccessPolicyRules(id string, policySetControllerService *services.Service) {
	policyRulesDetchLock.Lock()
	defer policyRulesDetchLock.Unlock()

	accessPolicySet, _, err := policysetcontroller.GetByPolicyType(policySetControllerService, "ACCESS_POLICY")
	if err != nil {
		return
	}
	accessPolicyRules, _, err := policysetcontroller.GetAllByType(policySetControllerService, "ACCESS_POLICY")
	if err != nil {
		return
	}
	for _, accessPolicyRule := range accessPolicyRules {
		ids := []servergroup.ServerGroup{}
		changed := false
		for _, app := range accessPolicyRule.AppServerGroups {
			if app.ID == id {
				changed = true
				continue
			}
			ids = append(ids, servergroup.ServerGroup{
				ID: app.ID,
			})
		}
		accessPolicyRule.AppServerGroups = ids
		if changed {
			if _, err := policysetcontroller.UpdateRule(policySetControllerService, accessPolicySet.ID, accessPolicyRule.ID, &accessPolicyRule); err != nil {
				continue
			}
		}
	}
}

func detachServerGroupFromAllAppSegments(id string, applicationSegmentService *services.Service) {
	apps, _, err := applicationsegment.GetAll(applicationSegmentService)
	if err != nil {
		return
	}
	for _, app := range apps {
		ids := []servergroup.ServerGroup{}
		for _, appServerGroup := range app.ServerGroups {
			if appServerGroup.ID == id {
				continue
			}
			ids = append(ids, servergroup.ServerGroup{
				ID: appServerGroup.ID,
			})
		}
		app.ServerGroups = ids
		if _, err := applicationsegment.Update(applicationSegmentService, app.ID, app); err != nil {
			continue
		}
	}
}

func detachServerGroupFromAppConnectorGroups(serverGroupID string, serverGroupService *services.Service, appConnectorGroupService *services.Service) error {
	log.Printf("[INFO] Detaching Server Group %s from App Connector Groups\n", serverGroupID)
	serverGroup, _, err := servergroup.Get(serverGroupService, serverGroupID)
	if err != nil {
		return err
	}

	// lock to avoid updating app connector group with a deleted server group ID when running in parallel
	detachLock.Lock()
	defer detachLock.Unlock()

	for _, appConnectorGroup := range serverGroup.AppConnectorGroups {
		app, _, err := appconnectorgroup.Get(appConnectorGroupService, appConnectorGroup.ID)
		if err != nil {
			continue
		}
		appServerGroups := []appconnectorgroup.AppServerGroup{}
		for _, s := range app.AppServerGroup {
			if s.ID == serverGroupID {
				continue
			}
			appServerGroups = append(appServerGroups, s)
		}
		app.AppServerGroup = appServerGroups
		_, err = appconnectorgroup.Update(appConnectorGroupService, app.ID, app)
		if err != nil {
			continue
		}
	}
	return nil
}

func expandServerGroup(d *schema.ResourceData) servergroup.ServerGroup {
	result := servergroup.ServerGroup{
		ID:                 d.Id(),
		Enabled:            d.Get("enabled").(bool),
		Description:        d.Get("description").(string),
		IpAnchored:         d.Get("ip_anchored").(bool),
		ConfigSpace:        d.Get("config_space").(string),
		DynamicDiscovery:   d.Get("dynamic_discovery").(bool),
		MicroTenantID:      d.Get("microtenant_id").(string),
		AppConnectorGroups: expandCommonAppConnectorGroups(d),
		Applications:       expandServerGroupApplications(d),
		Servers:            expandApplicationServers(d),
	}
	if d.HasChange("name") {
		result.Name = d.Get("name").(string)
	}
	return result
}

/*
func expandAppConnectorGroups(d *schema.ResourceData) []servergroup.AppConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("app_connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app connector groups data: %+v\n", appConnector)
		var appConnectorGroups []servergroup.AppConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, ok := appConnectorGroup.(map[string]interface{})
			if ok {
				for _, id := range appConnectorGroup["id"].(*schema.Set).List() {
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
*/

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

func expandApplicationServers(d *schema.ResourceData) []appservercontroller.ApplicationServer {
	applicationServersInterface, ok := d.GetOk("servers")
	if ok {
		applicationServer := applicationServersInterface.(*schema.Set)
		log.Printf("[INFO] server group application data: %+v\n", applicationServer)
		var applicationServers []appservercontroller.ApplicationServer
		for _, applicationServer := range applicationServer.List() {
			applicationServer, _ := applicationServer.(map[string]interface{})
			if applicationServer != nil {
				for _, id := range applicationServer["id"].([]interface{}) {
					applicationServers = append(applicationServers, appservercontroller.ApplicationServer{
						ID: id.(string),
					})
				}
			}
		}
		return applicationServers
	}

	return []appservercontroller.ApplicationServer{}
}

func flattenServerGroupApplicationsSimple(apps []servergroup.Applications) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(apps))
	for i, app := range apps {
		ids[i] = app.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

/*
func flattenAppConnectorGroupsSimple(appConnectorGroups []servergroup.AppConnectorGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appConnectorGroups))
	for i, group := range appConnectorGroups {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
*/
