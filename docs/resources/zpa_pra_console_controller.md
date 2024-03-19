---
subcategory: "Privileged Remote Access"
layout: "zscaler"
page_title: "ZPA): pra_console_controller"
description: |-
  Creates and manages ZPA privileged remote access console
---

# Resource: zpa_pra_console_controller

The **zpa_pra_console_controller** resource creates a privileged remote access console in the Zscaler Private Access cloud. This resource can then be referenced in an privileged access policy resource and a privileged access portal.

## Example Usage

```hcl
# Creates Privileged Remote Access Application Segment"
resource "zpa_application_segment_pra" "this" {
  name             = "Example"
  description      = "Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["3389", "3389"]
  domain_names     = [ "rdp_pra.example.com"]
  segment_group_id = zpa_segment_group.this.id
  common_apps_dto {
    apps_config {
      name                 = "rdp_pra"
      domain               = "rdp_pra.example.com"
      application_protocol = "RDP"
      connection_security  = "ANY"
      application_port     = "3389"
      enabled              = true
      app_types            = ["SECURE_REMOTE_ACCESS"]
    }
  }
}

# Creates Segment Group for Application Segment"
resource "zpa_segment_group" "this" {
  name        = "Example"
  description = "Example"
  enabled     = true
}

# Retrieves the Browser Access Certificate
data "zpa_ba_certificate" "this" {
  name = "pra01.example.com"
}

# Creates PRA Portal"
resource "zpa_pra_portal_controller" "this1" {
  name                      = "pra01.example.com"
  description               = "pra01.example.com"
  enabled                   = true
  domain                    = "pra01.example.com"
  certificate_id            = data.zpa_ba_certificate.this.id
  user_notification         = "Created with Terraform"
  user_notification_enabled = true
}

locals {
  pra_application_ids = {
    for app_dto in flatten([for common_apps in zpa_application_segment_pra.this.common_apps_dto : common_apps.apps_config]) :
    app_dto.name => app_dto.id
  }
  pra_application_id_rdp_pra = lookup(local.pra_application_ids, "rdp_pra", "")
}

resource "zpa_pra_console_controller" "ssh_pra" {
  name        = "ssh_console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = local.pra_application_id_rdp_pra
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this.id]
  }
}

```

## Attributes Reference

### Required

* `name` - (Required) The name of the privileged console.
* `description` - (Required) The description of the privileged console.
* `pra_application` - The Privileged Remote Access application segment resource
    - `id` - (String) The unique identifier of the Privileged Remote Access-enabled application.
    ~> **NOTE** This is the ID for each `apps_config` block within `common_apps_dto`
* `pra_portals` - The Privileged Remote Access Portal resource
    - `id` - (List) The unique identifier of the privileged portal.
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**pra_credential_controller** can be imported by using `<CONSOLE ID>` or `<CONSOLE NAME>` as the import ID.

For example:

```shell
terraform import zpa_pra_console_controller.this <console_id>
```

or

```shell
terraform import zpa_pra_console_controller.this <console_name>
```
