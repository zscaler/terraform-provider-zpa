package zpa

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgeschedule"
)

func resourceServiceEdgeAssistantSchedule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceServiceEdgeAssistantScheduleCreate,
		Read:     resourceServiceEdgeAssistantScheduleRead,
		Update:   resourceServiceEdgeAssistantScheduleUpdate,
		Delete:   resourceServiceEdgeAssistantScheduleDelete,
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

func resourceServiceEdgeAssistantScheduleCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ServiceEdgeSchedule

	req, err := expandServiceEdgeAssistantSchedule(d)
	if err != nil {
		return err
	}

	// Use = instead of := because err is already declared
	_, _, err = serviceedgeschedule.CreateSchedule(service, req)
	if err != nil {
		// Assuming err.Error() returns a string representation of the error
		errStr := err.Error()

		// Check if the error string contains the specific message indicating the resource already exists
		if strings.Contains(errStr, "resource.already.exist") {
			log.Printf("[INFO] Resource already exists. Updating instead.")

			// Get the current state of the resource
			resp, _, err := serviceedgeschedule.GetSchedule(service)
			if err != nil {
				return fmt.Errorf("failed to retrieve existing resource for update: %v", err)
			}

			// Set the resource ID in the Terraform state
			d.SetId(resp.ID)

			// Proceed to update the resource
			return resourceServiceEdgeAssistantScheduleUpdate(d, m)
		}
		return err
	}
	log.Printf("[INFO] Created service edge assistant schedule request. ID: %v\n", req.ID)
	d.SetId(req.ID)

	return resourceServiceEdgeAssistantScheduleRead(d, m)
}

func resourceServiceEdgeAssistantScheduleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ServiceEdgeSchedule

	resp, _, err := serviceedgeschedule.GetSchedule(service)
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing service edge assistant schedule %s from state because it no longer exists in ZPA", d.Id())
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

func resourceServiceEdgeAssistantScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ServiceEdgeSchedule

	id := d.Id()
	log.Printf("[INFO] Updating service edge ID: %v\n", id)
	req, err := expandServiceEdgeAssistantSchedule(d)
	if err != nil {
		return err
	}

	if _, _, err := serviceedgeschedule.GetSchedule(service); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := serviceedgeschedule.UpdateSchedule(service, id, &req); err != nil {
		return err
	}

	return resourceServiceEdgeAssistantScheduleRead(d, m)
}

func resourceServiceEdgeAssistantScheduleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func expandServiceEdgeAssistantSchedule(d *schema.ResourceData) (serviceedgeschedule.AssistantSchedule, error) {
	var customerID string
	if id, exists := d.GetOk("customer_id"); exists {
		customerID = id.(string)
	} else if id := os.Getenv("ZPA_CUSTOMER_ID"); id != "" {
		customerID = id
	} else {
		return serviceedgeschedule.AssistantSchedule{}, fmt.Errorf("customer_id must be provided either in the HCL or as an environment variable ZPA_CUSTOMER_ID")
	}

	scheduler := serviceedgeschedule.AssistantSchedule{
		ID:                d.Get("id").(string),
		CustomerID:        customerID, // Now guaranteed to be non-empty
		Enabled:           d.Get("enabled").(bool),
		DeleteDisabled:    d.Get("delete_disabled").(bool),
		FrequencyInterval: d.Get("frequency_interval").(string),
		Frequency:         d.Get("frequency").(string),
	}
	return scheduler, nil
}
