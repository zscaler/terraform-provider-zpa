package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/lssconfigcontroller"
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
			"connector_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "App Connector Group(s) to be added to the LSS configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the LSS configuration",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether this LSS configuration is enabled or not. Supported values: true, false",
						},
						"filter": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter for the LSS configuration. Format given by the following API to get status codes: /mgmtconfig/v2/admin/lssConfig/statusCodes",
						},
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Format of the log type. Format given by the following API to get formats: /mgmtconfig/v2/admin/lssConfig/logType/formats",
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Description:  "Name of the LSS configuration",
						},
						"lss_host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Host of the LSS configuration",
						},
						"lss_port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Port of the LSS configuration",
						},
						"source_log_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log type of the LSS configuration",
							ValidateFunc: validation.StringInSlice([]string{
								"zpn_trans_log",
								"zpn_auth_log",
								"zpn_ast_auth_log",
								"zpn_http_trans_log",
								"zpn_audit_log",
								"zpn_sys_auth_log",
								"zpn_http_insp",
							}, false),
						},
						"use_tls": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceLSSConfigControllerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandLSSResource(d)
	log.Printf("[INFO] Creating zpa lss config controller with request\n%+v\n", req)

	resp, _, err := zClient.lssconfigcontroller.Create(&req)
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
	_ = d.Set("config", flattenLSSConfig(resp.LSSConfig))
	_ = d.Set("connector_groups", flattenConnectorGroups(resp.ConnectorGroups))
	return nil

}

func resourceLSSConfigControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating lss config controller ID: %v\n", id)
	req := expandLSSResource(d)

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

func expandLSSResource(d *schema.ResourceData) lssconfigcontroller.LSSResource {
	req := lssconfigcontroller.LSSResource{

		ID:              d.Get("id").(string),
		LSSConfig:       expandLSSConfigController(d),
		ConnectorGroups: expandConnectorGroups(d),
	}
	return req
}

func expandLSSConfigController(d *schema.ResourceData) *lssconfigcontroller.LSSConfig {
	return &lssconfigcontroller.LSSConfig{
		AuditMessage:  d.Get("audit_message").(string),
		Description:   d.Get("description").(string),
		Enabled:       d.Get("enabled").(bool),
		Filter:        SetToStringList(d, "filter"),
		Format:        d.Get("format").(string),
		Name:          d.Get("name").(string),
		LSSHost:       d.Get("lss_host").(string),
		LSSPort:       d.Get("lss_port").(string),
		SourceLogType: d.Get("source_log_type").(string),
		UseTLS:        d.Get("use_tls").(bool),
	}
}

func expandConnectorGroups(d *schema.ResourceData) []lssconfigcontroller.ConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] connector groups data: %+v\n", appConnector)
		var appConnectorGroups []lssconfigcontroller.ConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, ok := appConnectorGroup.(map[string]interface{})
			if ok {
				for _, id := range appConnectorGroup["id"].([]interface{}) {
					appConnectorGroups = append(appConnectorGroups, lssconfigcontroller.ConnectorGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appConnectorGroups
	}

	return []lssconfigcontroller.ConnectorGroups{}
}
