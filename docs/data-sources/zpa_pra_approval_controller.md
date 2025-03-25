---
page_title: "zpa_pra_approval_controller Data Source - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-approvals
  API documentation https://help.zscaler.com/zpa/configuring-privileged-approvals-using-api
  Get information about ZPA privileged remote access approval in Zscaler Private Access cloud.
---

# zpa_pra_approval_controller (Data Source)

Use the **zpa_pra_approval_controller** data source to get information about a privileged remote access approval created in the Zscaler Private Access cloud.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
# ZPA PRA Portal Data Source
data "zpa_pra_approval_controller" "this" {
 email_ids = "jdoe@example.com"
}
```

## Schema

### Required

The following arguments are supported:

* `email_ids` - (Required) The name of the privileged remote access portal to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the privileged remote access portal to be exported.
* `start_time` - (string) The set start time in either `RFC1123Z` i.e `"Mon, 02 Jan 2006 15:04:05 -0700"` or `RFC1123` i.e `"Mon, 02 Jan 2006 15:04:05 MST"` format that the user has access to the Privileged Remote Access portal. 
    ~> **NOTE**: The approval `start_time` cannot be more than 1 hour in the past.
* `end_time` - (string)  The set end time in either `RFC1123Z` i.e `"Mon, 02 Jan 2006 15:04:05 -0700"` or `RFC1123` i.e `"Mon, 02 Jan 2006 15:04:05 MST"` format that the user has access to the Privileged Remote Access portal.
    ~> **NOTE**: The approval `end_time` cannot be more than 1 year or 365 days.
* `pra_application` - The Privileged Remote Access application segment resource
    - `id` - (List) The unique identifier of the Privileged Remote Access-enabled application segment.
* `working_hours` - The Privileged Remote Access application segment resource
    - `days` - (List) The days of the week that you want to enable the privileged approval. Supported values are: `"MON"`, `"TUE"`, `"WED"`, `"THU"`, `"FRI"`, `"SAT"`, `"SUN"`
    - `start_time` - (String) The start time that the user has access to the privileged approval.
    - `start_time_cron` - (string)  The cron expression provided to configure the privileged approval start time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]. For example, 0 15 10 ? * MON-FRI represents the start time working hours for 10:15 AM every Monday, Tuesday, Wednesday, Thursday and Friday.
    - `end_time` - (string) The end time that the user no longer has access to the privileged approval.
    - `end_time_cron` - (string) The cron expression provided to configure the privileged approval end time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]. For example, 0 15 10 ? * MON-FRI represents the end time working hours for 10:15 AM every Monday, Tuesday, Wednesday, Thursday and Friday.
    - `timezone` - (String) The time zone for the time window of a privileged approval in IANA format `"America/Vancouver"`.[Learn More](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - (string) The status of the privileged approval. The supported values are:
    - `INVALID`: The privileged approval is invalid.
    - `ACTIVE`: The privileged approval is currently available for the user.
    - `FUTURE`: The privileged approval is available for a user at a set time in the future.
    - `EXPIRED`: The privileged approval is no longer available for the user.

* `microtenant_id` (string)  The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.