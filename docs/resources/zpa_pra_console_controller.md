---
page_title: "zpa_pra_console_controller Resource - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-consoles
  API documentation https://help.zscaler.com/zpa/configuring-privileged-consoles-using-api
  Creates and manages ZPA privileged remote access console
---

# zpa_pra_console_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-consoles)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-consoles-using-api)

The **zpa_pra_console_controller** resource creates a privileged remote access console in the Zscaler Private Access cloud. This resource can then be referenced in an privileged access policy resource and a privileged access portal.

## Example Usage

```terraform
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

data "zpa_application_segment_by_type" "this" {
    application_type = "SECURE_REMOTE_ACCESS"
    name = "rdp_pra"
    depends_on = [zpa_application_segment_pra.this]
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


resource "zpa_pra_console_controller" "ssh_pra" {
  name        = "ssh_console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = data.zpa_application_segment_by_type.this.id
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this.id]
  }
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) The name of the privileged console.

- `pra_application` (Block Set, Max: 1) The Privileged Remote Access application segment resource
    - `id` - (String) The unique identifier of the Privileged Remote Access-enabled application.
    ~> **NOTE** This is the ID for each `apps_config` block within `common_apps_dto`
- `pra_portals` (Block Set) The Privileged Remote Access Portal resource
    - `id` - (List of Strings) The unique identifier of the privileged portal.

### Optional

In addition to all arguments above, the following attributes are exported:

- `description` - (String) The description of the privileged console.
- `enabled` - (Boolean) Whether or not the privileged console is enabled.
- `icon_text` - (String) The privileged console icon. The icon image is converted to base64 encoded text format.
- `microtenant_id` (String) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as `0` when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

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
