---
page_title: "zpa_app_connector_assistant_schedule Data Source - terraform-provider-zpa"
subcategory: "App Connector Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/deleting-disconnected-app-connectors
  documentation https://help.zscaler.com/zpa/configuring-auto-delete-disconnected-app-connectors-using-api
  Get information about ZPA App Connector Assistant Schedule in Zscaler Private Access cloud.
---

# zpa_app_connector_assistant_schedule (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/deleting-disconnected-app-connectors)
* [API documentation](https://help.zscaler.com/zpa/configuring-auto-delete-disconnected-app-connectors-using-api)

Use the **zpa_app_connector_assistant_schedule** data source to get information about Auto Delete frequency of the App Connector for the specified customer in the Zscaler Private Access cloud.

~> **NOTE** - The `customer_id` attribute is optional and not required during the configuration.

## Example Usage

```terraform
// Retrieve All Assistant Schedules
data "zpa_app_connector_assistant_schedule" "this" {}

// Retrieve A Specific Assistant Schedule by ID
data "zpa_app_connector_assistant_schedule" "this" {
    id = "1"
}

// Retrieve A Specific Assistant Schedule by the Customer ID
data "zpa_app_connector_assistant_schedule" "this" {
    customer_id = "1234567891012"
}
```

## Schema

### Required

The following arguments are supported:

* `id` - (Number) The unique identifier for the App Connector auto deletion configuration for a customer. This field is only required for the PUT request to update the frequency of the App Connector Settings.
* `customer_id` - (Number) The unique identifier of the ZPA tenant.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `enabled` (Boolean) - Indicates if the setting for deleting App Connectors is enabled or disabled.
* `delete_disabled` (Boolean) - Indicates if the App Connectors are included for deletion if they are in a disconnected state based on frequencyInterval and frequency values.
* `frequency` (String) - The scheduled frequency at which the disconnected App Connectors are deleted. Supported value is: `days`
* `frequency_interval` - (String) - The interval for the configured frequency value. The minimum supported value is 5. Supported values are: `5`, `7`, `14`, `30`, `60` and `90`
