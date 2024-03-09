package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praapproval"
)

func dataSourcePRAPrivilegedApprovalController() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePRAPrivilegedApprovalControllerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the privileged approval",
			},
			"email_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The email address of the user that you are assigning the privileged approval to",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start date that the user has access to the privileged approval",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The end date that the user no longer has access to the privileged approval",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the privileged approval",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged approval is created",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the tenant who modified the privileged approval",
			},
			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the privileged approval is modified",
			},
			"working_hours": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The days of the week that you want to enable the privileged approval",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start time that the user has access to the privileged approval",
						},
						"start_time_cron": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cron expression provided to configure the privileged approval start time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end time that the user no longer has access to the privileged approval",
						},
						"end_time_cron": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cron expression provided to configure the privileged approval end time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]The cron expression provided to configure the privileged approval end time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]",
						},
						"timezone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time zone for the time window of a privileged approval",
						},
					},
				},
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the pra application segment",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the pra application segment",
						},
					},
				},
			},
			"microtenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.",
			},
		},
	}
}

func dataSourcePRAPrivilegedApprovalControllerRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).praapproval.WithMicroTenant(GetString(d.Get("microtenant_id")))

	var resp *praapproval.PrivilegedApproval
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for pra approval controller %s\n", id)
		res, _, err := service.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	emailID, ok := d.Get("email_ids").(string)
	if id == "" && ok && emailID != "" {
		log.Printf("[INFO] Getting data for pra approval email ID %s\n", emailID)
		res, _, err := service.GetByEmailID(emailID)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("email_ids", resp.EmailIDs)
		_ = d.Set("start_time", resp.StartTime)
		_ = d.Set("end_time", resp.EndTime)
		_ = d.Set("status", resp.Status)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("microtenant_id", resp.MicroTenantID)

		if err := d.Set("working_hours", flattenWorkingHours(resp.WorkingHours)); err != nil {
			return fmt.Errorf("failed to read pra working hours %s", err)
		}
		if err := d.Set("applications", flattenPRAApplications(resp.Applications)); err != nil {
			return fmt.Errorf("failed to read pra applications %s", err)
		}

	} else {
		return fmt.Errorf("couldn't find any pra privileged approval with id '%s'", id)
	}

	return nil
}

func flattenWorkingHours(wh *praapproval.WorkingHours) []interface{} {
	if wh == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	result["days"] = wh.Days
	result["start_time"] = wh.StartTime
	result["end_time"] = wh.EndTime
	result["start_time_cron"] = wh.StartTimeCron
	result["end_time_cron"] = wh.EndTimeCron
	result["timezone"] = wh.TimeZone

	return []interface{}{result}
}

func flattenPRAApplications(applications []praapproval.Applications) []interface{} {
	praAppliations := make([]interface{}, len(applications))
	for i, praApplication := range applications {
		praAppliations[i] = map[string]interface{}{
			"id":   praApplication.ID,
			"name": praApplication.Name,
		}
	}

	return praAppliations
}
