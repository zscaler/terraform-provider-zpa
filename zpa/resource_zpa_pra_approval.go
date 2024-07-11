package zpa

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praapproval"
)

func resourcePRAPrivilegedApprovalController() *schema.Resource {
	return &schema.Resource{
		Create: resourcePRAPrivilegedApprovalControllerCreate,
		Read:   resourcePRAPrivilegedApprovalControllerRead,
		Update: resourcePRAPrivilegedApprovalControllerUpdate,
		Delete: resourcePRAPrivilegedApprovalControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*Client)
				service := client.PRAApproval

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
					resp, _, err := praapproval.GetByEmailID(service, id)
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the privileged approval",
			},
			"email_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The email address of the user that you are assigning the privileged approval to",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The start date that the user has access to the privileged approval",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The end date that the user no longer has access to the privileged approval",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"INVALID",
					"ACTIVE",
					"FUTURE",
					"EXPIRED",
				}, false),
				Description: "The status of the privileged approval",
			},
			"working_hours": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The days of the week that you want to enable the privileged approval",
						},
						"start_time": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate24HourTimeFormat,
							Description:  "The start time that the user has access to the privileged approval",
						},
						"start_time_cron": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The cron expression provided to configure the privileged approval start time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]",
						},
						"end_time": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate24HourTimeFormat,
							Description:  "The end time that the user no longer has access to the privileged approval",
						},
						"end_time_cron": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The cron expression provided to configure the privileged approval end time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]The cron expression provided to configure the privileged approval end time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]",
						},
						"timezone": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "The time zone for the time window of a privileged approval",
							ValidateFunc: validateTimeZone,
						},
					},
				},
			},
			"applications": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The unique identifier of the pra application segment",
						},
					},
				},
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant.",
			},
		},
	}
}

func resourcePRAPrivilegedApprovalControllerCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PRAApproval

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Convert user-provided RFC 2822 start and end times to epoch format.
	startTimeStr, endTimeStr := d.Get("start_time").(string), d.Get("end_time").(string)

	// Validate start and end times
	if err := validateStartTime(startTimeStr, endTimeStr); err != nil {
		return err
	}

	startTimeEpoch, err := convertToEpoch(startTimeStr)
	if err != nil {
		return fmt.Errorf("start time conversion error: %s", err)
	}
	endTimeEpoch, err := convertToEpoch(endTimeStr)
	if err != nil {
		return fmt.Errorf("end time conversion error: %s", err)
	}

	// Prepare the request object using the converted epoch times.
	req := expandPRAPrivilegedApproval(d)
	req.StartTime = fmt.Sprintf("%d", startTimeEpoch)
	req.EndTime = fmt.Sprintf("%d", endTimeEpoch)

	log.Printf("[INFO] Creating privileged approval with request\n%+v\n", req)

	praApproval, _, err := praapproval.Create(service, &req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created privileged approval request. ID: %v\n", praApproval.ID)

	d.SetId(praApproval.ID)
	return resourcePRAPrivilegedApprovalControllerRead(d, meta)
}

func resourcePRAPrivilegedApprovalControllerRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PRAApproval

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	resp, _, err := praapproval.Get(service, d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing privileged approval %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting privileged approval controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("email_ids", resp.EmailIDs)
	_ = d.Set("status", resp.Status)
	_ = d.Set("microtenant_id", resp.MicroTenantID)

	// Use the existing utility function to convert epoch to RFC1123
	startTimeStr, err := epochToRFC1123(resp.StartTime, false) // Adjust second parameter as needed
	if err != nil {
		return err
	}
	endTimeStr, err := epochToRFC1123(resp.EndTime, false) // Adjust second parameter as needed
	if err != nil {
		return err
	}

	_ = d.Set("start_time", startTimeStr)
	_ = d.Set("end_time", endTimeStr)

	_ = d.Set("applications", flattenPRAApplicationsSimple(resp.Applications))

	_ = d.Set("working_hours", flattenWorkingHours(resp.WorkingHours))

	return nil
}

func resourcePRAPrivilegedApprovalControllerUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PRAApproval

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	id := d.Id()
	log.Printf("[INFO] Updating privileged approval ID: %v\n", id)

	// Convert user-provided RFC 2822 start and end times to epoch format.
	startTimeStr, endTimeStr := d.Get("start_time").(string), d.Get("end_time").(string)

	// Validate start and end times
	if err := validateStartTime(startTimeStr, endTimeStr); err != nil {
		return err
	}

	startTimeEpoch, err := convertToEpoch(startTimeStr)
	if err != nil {
		return fmt.Errorf("start time conversion error: %s", err)
	}
	endTimeEpoch, err := convertToEpoch(endTimeStr)
	if err != nil {
		return fmt.Errorf("end time conversion error: %s", err)
	}

	// Prepare the request object using the converted epoch times.
	req := expandPRAPrivilegedApproval(d)
	req.StartTime = fmt.Sprintf("%d", startTimeEpoch)
	req.EndTime = fmt.Sprintf("%d", endTimeEpoch)

	if _, _, err := praapproval.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := praapproval.Update(service, id, &req); err != nil {
		return err
	}

	return resourcePRAPrivilegedApprovalControllerRead(d, meta)
}

func resourcePRAPrivilegedApprovalControllerDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PRAApproval

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting privileged approval ID: %v\n", d.Id())

	if _, err := praapproval.Delete(service, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] privileged approval deleted")
	return nil
}

func expandPRAPrivilegedApproval(d *schema.ResourceData) praapproval.PrivilegedApproval {
	result := praapproval.PrivilegedApproval{
		ID:           d.Id(),
		EmailIDs:     SetToStringList(d, "email_ids"),
		StartTime:    d.Get("start_time").(string),
		EndTime:      d.Get("end_time").(string),
		Status:       d.Get("status").(string),
		Applications: expandPRAApplications(d),
		WorkingHours: expandWorkingHours(d),
	}
	return result
}

func expandPRAApplications(d *schema.ResourceData) []praapproval.Applications {
	praAppsInterface, ok := d.GetOk("applications")
	if ok {
		praApp := praAppsInterface.(*schema.Set)
		log.Printf("[INFO] pra application data: %+v\n", praApp)
		var praApps []praapproval.Applications
		for _, praApp := range praApp.List() {
			praApp, ok := praApp.(map[string]interface{})
			if ok {
				for _, id := range praApp["id"].([]interface{}) {
					praApps = append(praApps, praapproval.Applications{
						ID: id.(string),
					})
				}
			}
		}
		return praApps
	}

	return []praapproval.Applications{}
}

func expandWorkingHours(d *schema.ResourceData) *praapproval.WorkingHours {
	if v, ok := d.GetOk("working_hours"); ok {
		workingHoursList := v.(*schema.Set).List()
		if len(workingHoursList) > 0 {
			workingHoursMap := workingHoursList[0].(map[string]interface{})
			days := []string{}
			if daysInterface, exists := workingHoursMap["days"].(*schema.Set); exists {
				for _, day := range daysInterface.List() {
					days = append(days, day.(string))
				}
			}

			return &praapproval.WorkingHours{
				Days:          days,
				StartTime:     workingHoursMap["start_time"].(string),
				EndTime:       workingHoursMap["end_time"].(string),
				StartTimeCron: workingHoursMap["start_time_cron"].(string),
				EndTimeCron:   workingHoursMap["end_time_cron"].(string),
				TimeZone:      workingHoursMap["timezone"].(string),
			}
		}
	}
	return nil
}

func flattenPRAApplicationsSimple(apps []praapproval.Applications) []interface{} {
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

func validateTimeZone(v interface{}, k string) (ws []string, errors []error) {
	tzStr := v.(string)
	_, err := time.LoadLocation(tzStr)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid timezone", tzStr))
	}

	return
}

// validate24HourTimeFormat validates that a string is in "HH:MM" 24-hour time format.
func validate24HourTimeFormat(v interface{}, k string) (ws []string, errors []error) {
	timeStr := v.(string)

	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q for %s is not a valid time format (expected HH:MM)", timeStr, k))
	}

	return
}

// This function is working
func convertToEpoch(dateStr string) (int64, error) {
	layouts := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			return t.Unix(), nil
		}
	}

	// If none of the formats match, return the last error
	return 0, fmt.Errorf("unable to parse date: %v", err)
}

// Validation function to be used with the start_time field.
func validateStartTime(startTimeStr interface{}, endTimeStr interface{}) error {
	startTimeEpoch, err := convertToEpoch(startTimeStr.(string))
	if err != nil {
		return fmt.Errorf("start time conversion error: %s", err)
	}
	currentTimeEpoch := time.Now().Unix()
	if startTimeEpoch < currentTimeEpoch-int64(3600) {
		return fmt.Errorf("the approval start time cannot be more than 1 hour in the past")
	}

	// Assuming endTimeStr is also passed as a string to this function.
	endTimeEpoch, err := convertToEpoch(endTimeStr.(string))
	if err != nil {
		return fmt.Errorf("end time conversion error: %s", err)
	}

	// Validate that end_time is within one year of start_time.
	oneYear := int64(31536000) // 365 * 24 * 60 * 60
	if endTimeEpoch > startTimeEpoch+oneYear {
		return fmt.Errorf("the start time should be less than the future end time with a max range of 1 year")
	}

	return nil
}
