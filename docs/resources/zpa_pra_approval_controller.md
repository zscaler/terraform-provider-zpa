---
page_title: "zpa_pra_approval_controller Resource - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-approvals
  API documentation https://help.zscaler.com/zpa/configuring-privileged-approvals-using-api
  Creates and manages ZPA privileged remote access approval
---

# zpa_pra_approval_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-approvals)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-approvals-using-api)

The **zpa_pra_approval_controller** resource creates a privileged remote access approval in the Zscaler Private Access cloud. This resource allows third-party users and contractors to be able to log in to a Privileged Remote Access (PRA) portal. 

## Example Usage

```terraform
# ZPA Application Segment resource
resource "zpa_application_segment" "this" {
    name              = "Example"
    description       = "Example"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
    tcp_port_ranges   = ["8080", "8080"]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.this.id
    server_groups {
        id = [ zpa_server_group.this.id]
    }
    depends_on = [ zpa_server_group.this, zpa_segment_group.this]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "this" {
  name            = "Example"
  description     = "Example"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.this.id ]
  }
  depends_on = [ zpa_app_connector_group.this ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "this" {
  name                          = "Example"
  description                   = "Example"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}

# Create PRA Approval Controller
resource "zpa_pra_approval_controller" "this" {
    email_ids = ["jdoe@acme.com"]
    start_time = "Tue, 07 Mar 2024 11:05:30 PST"
    end_time = "Tue, 07 Jun 2024 11:05:30 PST"
    status = "FUTURE"
    applications {
      id = [zpa_application_segment.this.id]
    }
    working_hours {
      days = ["FRI", "MON", "SAT", "SUN", "THU", "TUE", "WED"]
      start_time = "00:10"
      start_time_cron = "0 0 8 ? * MON,TUE,WED,THU,FRI,SAT"
      end_time = "09:15"
      end_time_cron = "0 15 17 ? * MON,TUE,WED,THU,FRI,SAT"
      timezone = "America/Vancouver"
    }
}
```

## Schema

### Required

The following arguments are supported:

- `email_ids` (List of Strings) The email_id associated with the privileged approval.
    ~> **NOTE**: Although the attribute `email_ids` is a list of strings, the ZPA API only supports a single email address.
- `start_time` (String) The set start time in either `RFC1123Z` i.e `"Mon, 02 Jan 2006 15:04:05 -0700"` or `RFC1123` i.e `"Mon, 02 Jan 2006 15:04:05 MST"` format that the user has access to the Privileged Remote Access portal. 
    ~> **NOTE**: The approval `start_time` cannot be more than 1 hour in the past.
- `end_time` (String) The set end time in either `RFC1123Z` i.e `"Mon, 02 Jan 2006 15:04:05 -0700"` or `RFC1123` i.e `"Mon, 02 Jan 2006 15:04:05 MST"` format that the user has access to the Privileged Remote Access portal.
    ~> **NOTE**: The approval `end_time` cannot be more than 1 year or 365 days.
- `applications` (Block Set) The unique identifier of the application segment.
    - `id` (List of Strings) The unique identifier of the application segment
- `working_hours` - The Privileged Remote Access application segment resource
    - `days` (List of Strings) The days of the week that you want to enable the privileged approval. Supported values are: `"MON"`, `"TUE"`, `"WED"`, `"THU"`, `"FRI"`, `"SAT"`, `"SUN"`
    - `start_time` - (String) The start time that the user has access to the privileged approval.
    - `start_time_cron` - (String) The cron expression provided to configure the privileged approval start time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]. For example, 0 15 10 ? * MON-FRI represents the start time working hours for 10:15 AM every Monday, Tuesday, Wednesday, Thursday and Friday.
    - `end_time` - (String) The end time that the user no longer has access to the privileged approval.
    - `end_time_cron` - (String) The cron expression provided to configure the privileged approval end time working hours. The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year]. For example, 0 15 10 ? * MON-FRI represents the end time working hours for 10:15 AM every Monday, Tuesday, Wednesday, Thursday and Friday.
    - `timezone` - (String) The time zone for the time window of a privileged approval in IANA format `"America/Vancouver"`.[Learn More](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)

### Optional

In addition to all arguments above, the following attributes are exported:

- `status` - (Required) The status of the privileged approval. The supported values are:
    - `INVALID`: The privileged approval is invalid.
    - `ACTIVE`: The privileged approval is currently available for the user.
    - `FUTURE`: The privileged approval is available for a user at a set time in the future.
    - `EXPIRED`: The privileged approval is no longer available for the user.

- `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_pra_approval_controller** can be imported by using `<APPROVAL ID>` or `<APPROVAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_pra_approval_controller.this <approval_id>
```

or

```shell
terraform import zpa_pra_approval_controller.this <approval_name>
```
