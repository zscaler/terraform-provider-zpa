---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Get information about ZPA Application Segment for Privileged Remote Access.
---

# Data Source: zpa_application_segment_pra

Use the **zpa_application_segment_pra** data source to get information about an application segment for Privileged Remote Access in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule, access policy client forwarding rule and inspection policy. This resource supports Privileged Remote Access for both `RDP` and `SSH`.

## Example Usage

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment_pra" "this" {
  name = "PRA_Example"
}
```

```hcl
# ZPA Application Segment Data Source
data "zpa_application_segment_pra" "this" {
  id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PRA Application Segment to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `domain_names` - (string) List of domains and IPs.
* `server_groups` - (string) List of Server Group IDs
  * `id:` - (string) List of Server Group IDs
* `segment_group_id` - (String) Segment Group IDs
* `creation_time` - (String)
* `modified_time` - (String)
* `modifiedby` - (String)
* `tcp_port_ranges` - (string) TCP port ranges used to access the app.
* `udp_port_ranges` - (string) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model:

* `tcp_port_range` - (string) TCP port ranges used to access the app.
  * `from:`
  * `to:`
* `udp_port_range` - (string) UDP port ranges used to access the app.
  * `from:`
  * `to:`

* `description` - (string) Description of the application.
* `bypass_type` - (string) Indicates whether users can bypass ZPA to access applications.
* `config_space` - (string)
* `double_encrypt` - (bool) Whether Double Encryption is enabled or disabled for the app.
* `enabled` - (bool) Whether this application is enabled or not.
* `health_reporting` - (string) Whether health reporting for the app is Continuous or On Access. Supported values: NONE, ON_ACCESS, CONTINUOUS.
* `health_check_type` (string)
* `icmp_access_type` - (string)
* `ip_anchored` - (bool)
* `is_cname_enabled` - (bool) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
* `passive_health_enabled` - (bool)

* `sra_apps` - (string) TCP port ranges used to access the app.
  * `app_id:` - (string)
  * `name:` - (string) Name of the Privileged Remote Access
  * `description:` - (string) Description of the Privileged Remote Access
  * `domain:` - (string) Domain name of the Privileged Remote Access
  * `application_port` - (string) Port for the Privileged Remote Accessvalues: `RDP` and `SSH`
  * `application_protocol` - (string) Protocol for the Privileged Remote Access. Supported values: `RDP` and `SSH`
  * `connection_security` - (string) - Parameter required when `application_protocol` is of type `RDP`
  * `enabled` - (bool) Whether this application is enabled or not
  * `select_connector_close_to_app` (bool)
  * `use_in_dr_mode` (bool)
  * `is_incomplete_dr_config` (bool)
  * `select_connector_close_to_app` (bool)
* `microtenant_id` (Computed)
* `microtenant_name` (Computed)
