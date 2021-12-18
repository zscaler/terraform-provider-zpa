package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/cloudconnectorgroup"
)

func dataSourceCloudConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudConnectorGroupRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_connectors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipacl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issued_cert_id": {
							Type:     schema.TypeString,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signing_cert": {
							Type:     schema.TypeMap,
							Elem:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"geolocation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zia_cloud": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zia_org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudConnectorGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *cloudconnectorgroup.CloudConnectorGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for cloud connector group  %s\n", id)
		res, _, err := zClient.cloudconnectorgroup.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for cloud connector group name %s\n", name)
		res, _, err := zClient.cloudconnectorgroup.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {

		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("geolocation_id", resp.GeolocationID)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("zia_cloud", resp.ZiaCloud)
		_ = d.Set("zia_org_id", resp.ZiaOrgid)
		_ = d.Set("cloud_connectors", flattenCloudConnectors(resp))

	} else {
		return fmt.Errorf("couldn't find any cloud connector group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenCloudConnectors(cloudConnectors *cloudconnectorgroup.CloudConnectorGroup) []interface{} {
	connectorItems := make([]interface{}, len(cloudConnectors.CloudConnectors))
	for i, connectorItem := range cloudConnectors.CloudConnectors {
		connectorItems[i] = map[string]interface{}{
			"creation_time":  connectorItem.CreationTime,
			"description":    connectorItem.Description,
			"enabled":        connectorItem.Enabled,
			"fingerprint":    connectorItem.Fingerprint,
			"id":             connectorItem.ID,
			"ipacl":          connectorItem.IPACL,
			"issued_cert_id": connectorItem.IssuedCertID,
			"modifiedby":     connectorItem.ModifiedBy,
			"modified_time":  connectorItem.ModifiedTime,
			"name":           connectorItem.Name,
		}
	}

	return connectorItems
}
