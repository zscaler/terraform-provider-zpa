package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorcontroller"
)

func dataSourceAppConnectorAssistantSchedule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppConnectorAssistantScheduleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"delete_disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frequency_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppConnectorAssistantScheduleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	var resp *appconnectorcontroller.AssistantSchedule
	var err error

	id, idOk := d.GetOk("id")
	customerID, customerIDOk := d.GetOk("customer_id")

	if idOk && id != "" {
		log.Printf("[INFO] Getting data for app connector assistant schedule %s\n", id)
		resp, _, err = zClient.appconnectorcontroller.GetSchedule()
		if err != nil {
			return err
		}
	} else if customerIDOk && customerID != "" {
		log.Printf("[INFO] Getting data for app connector name %s\n", customerID)
		resp, _, err = zClient.appconnectorcontroller.GetSchedule()
		if err != nil {
			return err
		}
	} else {
		log.Printf("[INFO] No specific ID or customer ID provided, fetching default schedule")
		resp, _, err = zClient.appconnectorcontroller.GetSchedule()
		if err != nil {
			return err
		}
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("customer_id", resp.CustomerID)
		_ = d.Set("enabled", resp.Enabled)
		_ = d.Set("delete_disabled", resp.DeleteDisabled)
		_ = d.Set("frequency", resp.Frequency)
		_ = d.Set("frequency_interval", resp.FrequencyInterval)
	} else {
		return fmt.Errorf("couldn't find any app connector assistant schedule")
	}

	return nil
}
