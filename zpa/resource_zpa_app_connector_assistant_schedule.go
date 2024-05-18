package zpa

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorcontroller"
)

func resourceAppConnectorAssistantSchedule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAppConnectorAssistantScheduleCreate,
		Read:     resourceAppConnectorAssistantScheduleRead,
		Update:   resourceAppConnectorAssistantScheduleUpdate,
		Delete:   resourceAppConnectorAssistantScheduleDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true, // Keep this to allow the value to be computed if not set
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"days",
				}, false),
			},
			"frequency_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5",
					"7",
					"14",
					"30",
					"60",
					"90",
				}, false),
			},
		},
	}
}

func resourceAppConnectorAssistantScheduleCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req, err := expandAssistantSchedule(d)
	if err != nil {
		return err
	}

	// Use = instead of := because err is already declared
	_, _, err = zClient.appconnectorcontroller.CreateSchedule(req)
	if err != nil {
		// Assuming err.Error() returns a string representation of the error
		errStr := err.Error()

		// Check if the error string contains the specific message indicating the resource already exists
		if strings.Contains(errStr, "resource.already.exist") {
			log.Printf("[INFO] Resource already exists. Updating instead.")

			// Get the current state of the resource
			resp, _, err := zClient.appconnectorcontroller.GetSchedule()
			if err != nil {
				return fmt.Errorf("failed to retrieve existing resource for update: %v", err)
			}

			// Set the resource ID in the Terraform state
			d.SetId(resp.ID)

			// Proceed to update the resource
			return resourceAppConnectorAssistantScheduleUpdate(d, m)
		}
		return err
	}
	log.Printf("[INFO] Created app connector assistant schedule request. ID: %v\n", req.ID)
	d.SetId(req.ID)

	return resourceAppConnectorAssistantScheduleRead(d, m)
}

func resourceAppConnectorAssistantScheduleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.appconnectorcontroller.GetSchedule()
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing app connector assistant schedule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting application server:\n%+v\n", resp)
	_ = d.Set("customer_id", resp.CustomerID)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("delete_disabled", resp.DeleteDisabled)
	_ = d.Set("frequency", resp.Frequency)
	_ = d.Set("frequency_interval", resp.FrequencyInterval)
	return nil
}

func resourceAppConnectorAssistantScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating app connector group ID: %v\n", id)
	req, err := expandAssistantSchedule(d)
	if err != nil {
		return err
	}

	if _, _, err := zClient.appconnectorcontroller.GetSchedule(); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.appconnectorcontroller.UpdateSchedule(id, &req); err != nil {
		return err
	}

	return resourceAppConnectorAssistantScheduleRead(d, m)
}

func resourceAppConnectorAssistantScheduleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func expandAssistantSchedule(d *schema.ResourceData) (appconnectorcontroller.AssistantSchedule, error) {
	var customerID string
	if id, exists := d.GetOk("customer_id"); exists {
		customerID = id.(string)
	} else if id := os.Getenv("ZPA_CUSTOMER_ID"); id != "" {
		customerID = id
	} else {
		return appconnectorcontroller.AssistantSchedule{}, fmt.Errorf("customer_id must be provided either in the HCL or as an environment variable ZPA_CUSTOMER_ID")
	}

	scheduler := appconnectorcontroller.AssistantSchedule{
		ID:                d.Get("id").(string),
		CustomerID:        customerID, // Now guaranteed to be non-empty
		Enabled:           d.Get("enabled").(bool),
		DeleteDisabled:    d.Get("delete_disabled").(bool),
		FrequencyInterval: d.Get("frequency_interval").(string),
		Frequency:         d.Get("frequency").(string),
	}
	return scheduler, nil
}
