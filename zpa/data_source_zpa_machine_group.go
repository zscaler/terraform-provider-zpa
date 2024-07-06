package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/machinegroup"
)

func dataSourceMachineGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMachineGroupRead,
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
						"modified_by": {
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
						"microtenant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"microtenant_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"modified_by": {
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
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microtenant_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.MachineGroup

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var resp *machinegroup.MachineGroup
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for machine group  %s\n", id)
		res, _, err := machinegroup.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for machine group name %s\n", name)
		res, _, err := machinegroup.GetByName(service, name)
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
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("microtenant_id", resp.MicroTenantID)
		_ = d.Set("microtenant_name", resp.MicroTenantName)
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
			"modified_by":        machineItem.ModifiedBy,
			"modified_time":      machineItem.ModifiedTime,
			"name":               machineItem.Name,
			"signing_cert":       machineItem.SigningCert,
			"microtenant_id":     machineItem.MicroTenantID,
			"microtenant_name":   machineItem.MicroTenantName,
		}
	}

	return machines
}
