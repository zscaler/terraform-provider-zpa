package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func resourceCBIBanners() *schema.Resource {
	return &schema.Resource{
		Create: resourceCBIBannersCreate,
		Read:   resourceCBIBannersRead,
		Update: resourceCBIBannersUpdate,
		Delete: resourceCBIBannersDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.cbibannercontroller.GetByName(id)
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"text_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logo": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"banner": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"persist": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceCBIBannersCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandCBIBanner(d)
	log.Printf("[INFO] Creating cbi banner with request\n%+v\n", req)

	cbiBanner, _, err := zClient.cbibannercontroller.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created cbi banner request. ID: %v\n", cbiBanner)

	d.SetId(cbiBanner.ID)
	return resourceCBIBannersRead(d, m)

}

func resourceCBIBannersRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.cbibannercontroller.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing cbi certificate %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting cbi certificate:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("primary_color", resp.PrimaryColor)
	_ = d.Set("text_color", resp.TextColor)
	_ = d.Set("notification_title", resp.NotificationTitle)
	_ = d.Set("notification_text", resp.NotificationText)
	_ = d.Set("logo", resp.Logo)
	_ = d.Set("banner", resp.Banner)

	return nil
}

func resourceCBIBannersUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating cbi certificate ID: %v\n", id)
	req := expandCBIBanner(d)

	if _, _, err := zClient.cbibannercontroller.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := zClient.cbibannercontroller.Update(id, &req); err != nil {
		return err
	}

	return resourceCBIBannersRead(d, m)
}

func resourceCBIBannersDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting cbi banner ID: %v\n", d.Id())

	if _, err := zClient.cbibannercontroller.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] cbi banner deleted")
	return nil
}

func expandCBIBanner(d *schema.ResourceData) cbibannercontroller.CBIBannerController {
	cbiBanner := cbibannercontroller.CBIBannerController{
		ID:                d.Id(),
		Name:              d.Get("name").(string),
		PrimaryColor:      d.Get("primary_color").(string),
		TextColor:         d.Get("text_color").(string),
		NotificationTitle: d.Get("notification_title").(string),
		NotificationText:  d.Get("notification_text").(string),
		Logo:              d.Get("logo").(string),
		Banner:            d.Get("banner").(bool),
		Persist:           d.Get("persist").(bool),
	}
	return cbiBanner
}
