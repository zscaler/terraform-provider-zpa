---
page_title: "zpa_microtenant_controller Data Source - terraform-provider-zpa"
subcategory: "Microtenant Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-microtenants
  API documentation https://help.zscaler.com/zpa/configuring-microtenants-using-api
  Get information about Microtenants in Zscaler Private Access cloud.
---

# zpa_microtenant_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-microtenants)
* [API documentation](https://help.zscaler.com/zpa/configuring-microtenants-using-api)

The **zpa_microtenant_controller** data source to get information about a machine group created in the Zscaler Private Access cloud. This data source allows administrators to retrieve a specific microtenant ID, which can be passed to other supported resources via the `microtenant_id` attribute.

⚠️ **WARNING:**: This feature is in limited availability and requires additional license. To learn more, contact Zscaler Support or your local account team.

## Example Usage

```terraform
# ZPA Microtenant Controller Data Source
data "zpa_microtenant_controller" "this" {
  name = "Microtenant_A"
}
```

## Schema

### Required

* `name` - (Required) Name of the microtenant controller.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `criteria_attribute` - (string) Type of authentication criteria for the microtenant
* `criteria_attribute_values` - (string) The domain associated with the respective microtenant controller resource
* `description` (string) Description of the microtenant controller.
* `enabled` (bool) Whether this microtenant resource is enabled or not.
