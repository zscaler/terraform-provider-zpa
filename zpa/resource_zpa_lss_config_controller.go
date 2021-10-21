package zpa

/*
import (
	"log"

	"github.com/willguibr/terraform-provider-zpa/gozscaler/lssconfigcontroller"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
)

func resourceLSSConfigController() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLSSConfigControllerCreate,
		Read:     resourceLSSConfigControllerRead,
		Update:   resourceLSSConfigControllerUpdate,
		Delete:   resourceLSSConfigControllerDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"filter": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"lss_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"lss_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_log_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"use_tls": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"connector_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceLSSConfigControllerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandLSSConfigController(d)
	log.Printf("[INFO] Creating zpa lss config controller with request\n%+v\n", req)

	resp, _, err := zClient.lssconfigcontroller.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created lss config controller request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceLSSConfigControllerRead(d, m)
}

func resourceLSSConfigControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.lssconfigcontroller.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing lss config controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting lss config controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("config", flattenAppServerGroupsSimple(resp))
	return nil

}

func resourceLSSConfigControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating lss config controller ID: %v\n", id)
	req := expandLSSConfigController(d)

	if _, err := zClient.lssconfigcontroller.Update(id, &req); err != nil {
		return err
	}

	return resourceLSSConfigControllerRead(d, m)
}

func resourceLSSConfigControllerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting lss config controller ID: %v\n", d.Id())

	if _, err := zClient.lssconfigcontroller.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] lss config controller deleted")
	return nil
}
*/
