package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
)

func dataSourceServerGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_space": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dynamic_discovery": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"app_connector_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.ServerGroup

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	var resp *servergroup.ServerGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for server group  %s\n", id)
		res, _, err := servergroup.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for server group name %s\n", name)
		res, _, err := servergroup.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("config_space", resp.ConfigSpace)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("dynamic_discovery", resp.DynamicDiscovery)
		_ = d.Set("ip_anchored", resp.IpAnchored)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)

		if err := d.Set("applications", flattenApplicationsSegments(resp.Applications)); err != nil {
			return err
		}

		if err := d.Set("app_connector_groups", flattenAppConnectorGroups(resp.AppConnectorGroups)); err != nil {
			return err
		}

		if err := d.Set("servers", flattenServers(resp.Servers)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any server group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenApplicationsSegments(applications []servergroup.Applications) []interface{} {
	serverGroupApplications := make([]interface{}, len(applications))
	for i, srvApplication := range applications {
		serverGroupApplications[i] = map[string]interface{}{
			"id":   srvApplication.ID,
			"name": srvApplication.Name,
		}
	}

	return serverGroupApplications
}

func flattenAppConnectorGroups(appConnectorGroup []appconnectorgroup.AppConnectorGroup) []interface{} {
	appConnectorGroups := make([]interface{}, len(appConnectorGroup))
	for i, appConnectorGroup := range appConnectorGroup {
		appConnectorGroups[i] = map[string]interface{}{
			"id":   appConnectorGroup.ID,
			"name": appConnectorGroup.Name,
		}
	}

	return appConnectorGroups
}

func flattenServers(applicationServer []appservercontroller.ApplicationServer) []interface{} {
	applicationServers := make([]interface{}, len(applicationServer))
	for i, appServerItem := range applicationServer {
		applicationServers[i] = map[string]interface{}{
			"id":   appServerItem.ID,
			"name": appServerItem.Name,
		}
	}
	return applicationServers
}
