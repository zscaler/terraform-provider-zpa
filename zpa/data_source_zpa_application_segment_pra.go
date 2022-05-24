package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/app_segment_sra_apps"
)

func dataSourceSRAPortalAppSegment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSRAPortalAppSegmentRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"segment_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"segment_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bypass_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_space": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the application.",
			},
			"domain_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of domains and IPs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"double_encrypt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Double Encryption is enabled or disabled for the app.",
			},
			"health_check_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"passive_health_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"health_reporting": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.",
			},
			"ip_anchored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_cname_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the application.",
			},
			"common_apps_dto": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_options": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"app_types": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"application_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"application_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ba_app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cname": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"connection_security": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"hidden": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"local_domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"portal": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sra_app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"trust_untrusted_cert": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"sra_apps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_security": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hidden": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"portal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"tcp_port_range": resourceAppSegmentPortRange("tcp port range"),
			"udp_port_range": resourceAppSegmentPortRange("udp port range"),

			"tcp_port_ranges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "TCP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"udp_port_ranges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UDP port ranges used to access the app.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"server_groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of the server group IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Computed: true,
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

func dataSourceSRAPortalAppSegmentRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	var resp *app_segment_sra_apps.AppSegmentSraApps
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for sra application %s\n", id)
		res, _, err := zClient.app_segment_sra_apps.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for sra application name %s\n", name)
		res, _, err := zClient.app_segment_sra_apps.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("segment_group_id", resp.SegmentGroupID)
		_ = d.Set("segment_group_name", resp.SegmentGroupName)
		_ = d.Set("bypass_type", resp.BypassType)
		_ = d.Set("config_space", resp.ConfigSpace)
		_ = d.Set("domain_names", resp.DomainNames)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("passive_health_enabled", resp.PassiveHealthEnabled)
		_ = d.Set("double_encrypt", resp.DoubleEncrypt)
		_ = d.Set("health_check_type", resp.HealthCheckType)
		_ = d.Set("is_cname_enabled", resp.IsCnameEnabled)
		_ = d.Set("ip_anchored", resp.IPAnchored)
		_ = d.Set("health_reporting", resp.HealthReporting)
		_ = d.Set("tcp_port_ranges", resp.TCPPortRanges)
		_ = d.Set("udp_port_ranges", resp.UDPPortRanges)

		if err := d.Set("sra_apps", flattenSRAApps(resp)); err != nil {
			return fmt.Errorf("failed to read clientless apps %s", err)
		}

		if err := d.Set("common_apps_dto", flattenCommonAppDto(resp.CommonApplicationDto)); err != nil {
			return fmt.Errorf("failed to read common apps %s", err)
		}

		if err := d.Set("server_groups", flattenSRAAppServerGroups(resp.AppServerGroups)); err != nil {
			return fmt.Errorf("failed to read app server groups %s", err)
		}

		if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.TCPAppPortRange)); err != nil {
			return err
		}

		if err := d.Set("tcp_port_range", flattenNetworkPorts(resp.UDPAppPortRange)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any browser access application with name '%s' or id '%s'", name, id)
	}

	return nil

}

func flattenSRAAppServerGroups(appServerGroup []app_segment_sra_apps.AppServerGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appServerGroup))
	for i, serverGroup := range appServerGroup {
		ids[i] = serverGroup.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenCommonAppDto(commonAppDto *app_segment_sra_apps.CommonApplicationDto) []map[string]interface{} {
	if commonAppDto == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"apps_config": flattenAppsConfig(commonAppDto.AppsConfig),
		},
	}
}

func flattenAppsConfig(appsConfig []app_segment_sra_apps.AppsConfig) []interface{} {
	appsConfigs := make([]interface{}, len(appsConfig))
	for i, appsConfig := range appsConfig {
		appsConfigs[i] = map[string]interface{}{
			"allow_options":        appsConfig.AllowOptions,
			"app_id":               appsConfig.AppID,
			"app_types":            appsConfig.AppTypes,
			"application_port":     appsConfig.ApplicationPort,
			"application_protocol": appsConfig.ApplicationProtocol,
			"ba_app_id":            appsConfig.BaAppId,
			"certificate_id":       appsConfig.CertificateID,
			"certificate_name":     appsConfig.CertificateName,
			"cname":                appsConfig.Cname,
			"connection_security":  appsConfig.ConnectionSecurity,
			"description":          appsConfig.Description,
			"domain":               appsConfig.Domain,
			"enabled":              appsConfig.Enabled,
			"hidden":               appsConfig.Hidden,
			"local_domain":         appsConfig.LocalDomain,
			"name":                 appsConfig.Name,
			"path":                 appsConfig.Path,
			"portal":               appsConfig.Portal,
			"sra_app_id":           appsConfig.SraAppId,
			"trust_untrusted_cert": appsConfig.TrustUntrustedCert,
		}
	}

	return appsConfigs
}

func flattenSRAApps(sraApp *app_segment_sra_apps.AppSegmentSraApps) []interface{} {
	sraApps := make([]interface{}, len(sraApp.SraApps))
	for i, val := range sraApp.SraApps {
		sraApps[i] = map[string]interface{}{
			"id":                   val.ID,
			"app_id":               val.AppID,
			"application_port":     val.ApplicationPort,
			"application_protocol": val.ApplicationProtocol,
			"certificate_id":       val.CertificateID,
			"certificate_name":     val.CertificateName,
			"connection_security":  val.ConnectionSecurity,
			"description":          val.Description,
			"domain":               val.Domain,
			"enabled":              val.Enabled,
			"hidden":               val.Hidden,
			"name":                 val.Name,
			"portal":               val.Portal,
		}
	}

	return sraApps
}
