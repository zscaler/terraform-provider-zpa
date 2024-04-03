---
page_title: "zpa_pra_portal_controller Data Source - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-portals
  API documentation https://help.zscaler.com/zpa/configuring-privileged-portals-using-api
  Get information about ZPA privileged remote access portal in Zscaler Private Access cloud.
---

# zpa_pra_portal_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-portals)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-portals-using-api)

Use the **zpa_pra_portal_controller** data source to get information about a privileged remote access portal created in the Zscaler Private Access cloud. This data source can then be referenced in an privileged remote access console resource.

## Example Usage

```terraform
# ZPA PRA Portal Data Source
data "zpa_pra_portal_controller" "this" {
 name = "Example"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (String) The name of the privileged remote access portal to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the privileged remote access portal to be exported.
* `description` - (string)
* `enabled` (bool) Whether or not the privileged portal is enabled.
* `certificate_id` - (string) The unique identifier of the certificate.
* `user_notification` (string) The notification message displayed in the banner of the privileged portallink, if enabled.
* `user_notification_enabled` (bool) Indicates if the Notification Banner is enabled (true) or disabled (false).
* `microtenant_id` (string) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
* `creation_time` - (string) The time the privileged portal is created.
* `modified_time` - (string) The time the privileged portal is modified.
* `modified_by` - (string) The unique identifier of the tenant who modified the privileged portal.
