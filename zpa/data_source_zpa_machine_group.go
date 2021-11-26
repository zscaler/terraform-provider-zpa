package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/machinegroup"
)

func machineGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"machines": {
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
					"fingerprint": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"issued_cert_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"machine_group_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"machine_group_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"machine_token_id": {
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
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
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
	}
}

func dataSourceMachineGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceMachineGroupRead,
		Schema: machineGroupSchema(),
	}
}

func dataSourceMachineGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *machinegroup.MachineGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for machine group  %s\n", id)
		res, _, err := zClient.machinegroup.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for machine group name %s\n", name)
		res, _, err := zClient.machinegroup.GetByName(name)
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
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("machines", flattenMachines(resp))

	} else {
		return fmt.Errorf("couldn't find any machine group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenMachines(machineGroup *machinegroup.MachineGroup) []interface{} {
	machines := make([]interface{}, len(machineGroup.Machines))
	for i, machineItem := range machineGroup.Machines {
		machines[i] = map[string]interface{}{
			"creation_time":      machineItem.CreationTime,
			"description":        machineItem.Description,
			"fingerprint":        machineItem.Fingerprint,
			"id":                 machineItem.ID,
			"issued_cert_id":     machineItem.IssuedCertID,
			"machine_group_id":   machineItem.MachineGroupID,
			"machine_group_name": machineItem.MachineGroupName,
			"machine_token_id":   machineItem.MachineTokenID,
			"modifiedby":         machineItem.ModifiedBy,
			"modified_time":      machineItem.ModifiedTime,
			"name":               machineItem.Name,
			"signing_cert":       machineItem.SigningCert,
		}
	}

	return machines
}
