---
page_title: "zpa_microtenant_controller Resource - terraform-provider-zpa"
subcategory: "Microtenant Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-microtenants/
  API documentation https://help.zscaler.com/zpa/configuring-microtenants-using-api
  Creates and manages ZPA Microtenant resources
---

# zpa_microtenant_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-microtenants)
* [API documentation](https://help.zscaler.com/zpa/configuring-microtenants-using-api)

The **zpa_microtenant_controller** resource creates a microtenant controller in the Zscaler Private Access cloud. This resource allows organizations to delegate responsibilities of admins directly to the acquired or merged company admins so that they can manage their configurations independently

⚠️ **WARNING:**: This feature is in limited availability and requires additional license. To learn more, contact Zscaler Support or your local account team.

## Example Usage

```terraform
# ZPA Microtenant Controller resource
resource "zpa_microtenant_controller" "this" {
  name = "Microtenant_A"
  description = "Microtenant_A"
  enabled = true
  criteria_attribute = "AuthDomain"
  criteria_attribute_values = ["acme.com"]
}

// To output specific Microtenant user information,
// the following output configuration is required.
output "zpa_microtenant_controller1" {
  value = [for u in zpa_microtenant_controller.this.user : {
    microtenant_id = u.microtenant_id
    username       = u.username
    password       = u.password
  }]
}
```

## Schema

### Required

* `name` - (Required) Name of the microtenant controller.
* `criteria_attribute` - (Required) Type of authentication criteria for the microtenant
* `criteria_attribute_values` - (Required) The domain associated with the respective microtenant controller resource

### Optional

In addition to all arguments above, the following attributes are exported:

* `description` (Optional) Description of the microtenant controller.
* `enabled` (Optional) Whether this microtenant resource is enabled or not.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**microtenant_controller** can be imported by using `<MICROTENANT ID>` or `<MICROTENANT NAME>` as the import ID.

For example:

```shell
terraform import zpa_microtenant_controller.example <microtenant_id>
```

or

```shell
terraform import zpa_microtenant_controller.example <microtenant_name>
```
