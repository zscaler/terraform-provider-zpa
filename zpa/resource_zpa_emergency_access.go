package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/emergencyaccess"
)

func resourceEmergencyAccess() *schema.Resource {
	return &schema.Resource{
		Create:   resourceEmergencyAccessCreate,
		Read:     resourceEmergencyAccessRead,
		Update:   resourceEmergencyAccessUpdate,
		Delete:   resourceEmergencyAccessDeactivated,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: "The unique identifier of the emergency access user",
				Computed:    true,
				Optional:    true,
			},
			"email_id": {
				Type:        schema.TypeString,
				Description: "The email address of the emergency access user, as provided by the admin",
				Optional:    true,
				ForceNew:    true,
			},
			"first_name": {
				Type:        schema.TypeString,
				Description: "The first name of the emergency access user",
				Optional:    true,
				Computed:    true,
			},
			"last_name": {
				Type:        schema.TypeString,
				Description: "The last name of the emergency access user, as provided by the admin",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceEmergencyAccessCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandEmergencyAccess(d)
	log.Printf("[INFO] Creating emergency access user with request\n%+v\n", req)

	emgAccess, _, err := zClient.emergencyaccess.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created emergency access user request. ID: %v\n", emgAccess)

	d.SetId(emgAccess.UserID)
	return resourceEmergencyAccessRead(d, m)
}

func resourceEmergencyAccessRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.emergencyaccess.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing emergency access user %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting emergency access user:\n%+v\n", resp)
	d.SetId(resp.UserID)
	_ = d.Set("first_name", resp.FirstName)
	_ = d.Set("last_name", resp.LastName)
	_ = d.Set("email_id", resp.EmailID)
	return nil
}

func resourceEmergencyAccessUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating emergency access user ID: %v\n", id)
	req := expandEmergencyAccess(d)

	if _, _, err := zClient.emergencyaccess.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.emergencyaccess.Update(id, &req); err != nil {
		return err
	}

	return resourceEmergencyAccessRead(d, m)
}

func resourceEmergencyAccessDeactivated(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deactivated Emergency Access User ID: %v\n", d.Id())

	if _, err := zClient.emergencyaccess.Deactivate(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] Emergency Access User ID deactivated")
	return nil
}

func expandEmergencyAccess(d *schema.ResourceData) emergencyaccess.EmergencyAccess {
	emgAccessUser := emergencyaccess.EmergencyAccess{
		UserID:    d.Id(),
		EmailID:   d.Get("email_id").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
	}
	return emgAccessUser
}
