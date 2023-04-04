---
subcategory: "Application Segment"
layout: "zscaler"
page_title: "ZPA: application_segment"
description: |-
  Creates and manages ZPA Application Segment for Inspection.
---

# Resource: zpa_application_segment_inspection

The **zpa_application_segment_inspection** resource creates an inspection application segment in the Zscaler Private Access cloud. This resource can then be referenced in an access policy inspection rule. This resource supports Inspection for both `HTTP` and `HTTPS`.

## Example Usage

```hcl
data "zpa_ba_certificate" "jenkins" {
  name = "jenkins.securitygeek.io"
}

resource "zpa_application_segment_inspection" "this" {
  name             = "ZPA_Inspection_Example"
  description      = "ZPA_Inspection_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["443", "443"]
  domain_names     = ["jenkins.example.com"]
  segment_group_id = zpa_segment_group.this.id
  server_groups {
    id = [zpa_server_group.this.id]
  }
  common_apps_dto {
    apps_config {
      name                 = "jenkins.example.com"
      domain               = "jenkins.example.com"
      application_protocol = "HTTPS"
      application_port     = "443"
      certificate_id       = data.zpa_ba_certificate.jenkins.id
      enabled              = true
      app_types            = [ "INSPECT" ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the App Connector Group to be exported.
* `domain_names` - (Required) List of domains and IPs.
* `server_groups` - (Required) List of Server Group IDs
* `segment_group_id` - (Required) List of Segment Group IDs
* `common_apps_dto` - (Required) List of applications (e.g., Inspection, Browser Access or Privileged Remote Access)
  * `apps_config:` - (Required) List of applications to be configured
    * `name` - (Required) Name of the Inspection Application Segment.
    * `domain` - (Required) Domain name of the Inspection Application Segment.
    * `application_protocol` - (Required) Protocol for the Inspection Application Segment.. Supported values: `HTTP` and `HTTPS`
    * `application_port` - (Required) Port for the Inspection Application Segment.
    * `app_types` - (Required) Indicates the type of application as inspection. Supported value: `INSPECT`
    * `certificate_id` - (string) - ID of the signing certificate. This field is required if the applicationProtocol is set to `HTTPS`. The certificateId is not supported if the applicationProtocol is set to `HTTP`.
    * `enabled` - (Optional) Whether this application is enabled or not
* `tcp_port_ranges` - (Required) TCP port ranges used to access the app.
* `udp_port_ranges` - (Required) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model:
-> **NOTE:** When removing TCP and/or UDP ports, parameter must be defined but set as empty due to current API behavior.

* `tcp_port_range` - (Required) TCP port ranges used to access the app.
  * `from:`
  * `to:`
* `udp_port_range` - (Required) UDP port ranges used to access the app.
  * `from:`
  * `to:`

-> **NOTE:** Application segments must have unique ports and cannot have overlapping domain names using the same tcp/udp ports across multiple application segments.

## Attributes Reference

* `description` - (Optional) Description of the application.
* `bypass_type` - (Optional) Indicates whether users can bypass ZPA to access applications.
* `config_space` - (Optional)
* `double_encrypt` - (Optional) Whether Double Encryption is enabled or disabled for the app.
* `enabled` - (Optional) Whether this application is enabled or not.
* `health_reporting` - (Optional) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
* `health_check_type` (Optional)
* `icmp_access_type` - (Optional)
* `ip_anchored` - (Optional)
* `is_cname_enabled` - (Optional) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors.
* `tcp_keep_alive` (Optional) Supported values: `true` for Enabled and `false` for Disabled
* `passive_health_enabled` - (Optional) Supported values: `true`, `false`
* `select_connector_close_to_app` - (Optional) Supported values: `true`, `false`
* `use_in_dr_mode` - (Optional) Supported values: `true`, `false`
* `is_incomplete_dr_config` - (Optional) Supported values: `true`, `false`
* `select_connector_close_to_app` - (Optional) Supported values: `true`, `false`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Inspection Application Segment can be imported by using `<APPLICATION SEGMENT ID>` or `<APPLICATION SEGMENT NAME>` as the import ID.

```shell
terraform import zpa_application_segment_inspection.example <application_segment_id>
```

or

```shell
terraform import zpa_application_segment_inspection.example <application_segment_name>
```
