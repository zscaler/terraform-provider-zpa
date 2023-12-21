---
subcategory: "App Connector Controller"
layout: "zscaler"
page_title: "ZPA: app_connector_assistant_schedule"
description: |-
  Get information about ZPA App Connector Assistant Schedule in Zscaler Private Access cloud.
---

# Data Source: app_connector_assistant_schedule

Use the **zpa_app_connector_assistant_schedule** data source to get information about Auto Delete frequency of the App Connector for the specified customer in the Zscaler Private Access cloud.

~> **NOTE** - The `customer_id` attribute is optional and not required during the configuration.

## Example Usage

```hcl
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

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique identifier for the App Connector auto deletion configuration for a customer. This field is only required for the PUT request to update the frequency of the App Connector Settings.
* `customer_id` - (Optional) The unique identifier of the ZPA tenant.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` (Boolean) - Indicates if the setting for deleting App Connectors is enabled or disabled.
* `delete_disabled` (Boolean) - Indicates if the App Connectors are included for deletion if they are in a disconnected state based on frequencyInterval and frequency values.
* `frequency` (String) - The scheduled frequency at which the disconnected App Connectors are deleted. Supported value is: `days`
* `frequency_interval` - (String) - The interval for the configured frequency value. The minimum supported value is 5. Supported values are: `5`, `7`, `14`, `30`, `60` and `90`
