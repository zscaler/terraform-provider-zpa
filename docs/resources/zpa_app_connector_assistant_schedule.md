---
subcategory: "App Connector Controller"
layout: "zscaler"
page_title: "ZPA: app_connector_assistant_schedule"
description: |-
  Configures Auto Delete for the specified disconnected App Connector.
---

# Data Source: app_connector_assistant_schedule

Use the **zpa_app_connector_assistant_schedule** resource sets the scheduled frequency at which the disconnected App Connectors are eligible for deletion. The supported value for frequency is days. The frequencyInterval field is the number of days after an App Connector disconnects for it to become eligible for deletion. The minimum supported value for frequencyInterval is 5.

~> **NOTE** - When enabling the Assistant Schedule for the first time, you must provide the `customer_id` information. If you authenticated using environment variables and used `ZPA_CUSTOMER_ID` environment variable, you don't have to define the customer_id attribute in the HCL configuration, and the provider will automatically use the value from the environment variable `ZPA_CUSTOMER_ID`

## Example Usage - Defined Customer ID Value

```hcl
resource "zpa_app_connector_assistant_schedule" "this" {
  customer_id = "123456789101112"
  frequency = "days"
  frequency_interval = "5"
  enabled = true
  delete_disabled = true
}
```

## Example Usage - Customer ID Via Environment Variable

```hcl
resource "zpa_app_connector_assistant_schedule" "this" {
  frequency = "days"
  frequency_interval = "5"
  enabled = true
  delete_disabled = true
}
```

## Argument Reference

The following arguments are supported:

* `customer_id` - (Optional) - When enabling the Assistant Schedule for the first time, you must provide the `customer_id` information. If you authenticated using environment variables and used `ZPA_CUSTOMER_ID` environment variable, you don't have to define the customer_id attribute in the HCL configuration, and the provider will automatically use the value from the environment variable `ZPA_CUSTOMER_ID`
* `frequency_interval` - (Required) - The interval for the configured frequency value. The minimum supported value is 5. Supported values are: `5`, `7`, `14`, `30`, `60` and `90`
* `frequency` (Optional) - The scheduled frequency at which the disconnected App Connectors are deleted. Supported value is: `days`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` (Optional) - Indicates if the setting for deleting App Connectors is enabled or disabled. Supported values are: `true` or `false`
* `delete_disabled` (Optional) - Indicates if the App Connectors are included for deletion if they are in a disconnected state based on frequencyInterval and frequency values. Supported values are: `true` or `false`
