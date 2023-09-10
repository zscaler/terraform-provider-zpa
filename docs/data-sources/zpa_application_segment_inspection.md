---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Get information about ZPA Application Segment for Inspection.
---

# Data Source: zpa_application_segment_inspection

Use the **zpa_application_segment_inspection** data source to get information about an inspection application segment in the Zscaler Private Access cloud. This resource can then be referenced in a ZPA access inspection policy. This resource supports ZPA Inspection for both `HTTP` and `HTTPS`.

## Example Usage

```hcl
# ZPA Inspection Application Segment Data Source
data "zpa_application_segment_inspection" "this" {
  name = "ZPA_Inspection_Example"
}
```

```hcl
# ZPA Inspection Application Segment Data Source
data "zpa_application_segment_inspection" "this" {
  id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Inspection Application Segment to be exported.
* `id` - (Optional) The ID of the Inspection Application Segment to be exported.

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
* `health_reporting` - (string) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
* `health_check_type` (string)
* `icmp_access_type` - (string)
* `ip_anchored` - (bool)
* `is_cname_enabled` - (bool) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
* `passive_health_enabled` - (bool)
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

* `inspection_apps` - (string) TCP port ranges used to access the app.
  * `app_id:` - (string)
  * `name:` - (string) Name of the Inspection Application
  * `description:` - (string) Description of the Inspection Application
  * `domain:` - (string) Domain name of the inspection application
  * `application_port` - (string) TCP/UDP Port for ZPA Inspection.
  * `application_protocol` - (string) Protocol for the Inspection Application. Supported values: `HTTP` and `HTTPS`
  * `certificate_id` - (string) - ID of the signing certificate. This field is required if the applicationProtocol is set to `HTTPS`. The certificateId is not supported if the applicationProtocol is set to `HTTP`.
  * `certificate_name` - (string) - Parameter required when `application_protocol` is of type `HTTPS`
  * `enabled` - (bool) Whether this application is enabled or not
  * `select_connector_close_to_app` (bool)
  * `use_in_dr_mode` (bool)
  * `is_incomplete_dr_config` (bool)
  * `select_connector_close_to_app` (bool)
  * `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  * `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
