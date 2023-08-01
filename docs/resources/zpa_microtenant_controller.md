---
subcategory: "Microtenant"
layout: "zscaler"
page_title: "ZPA: microtenant_controller"
description: |-
  Creates and manages ZPA Microtenants.
---

# Resource: zpa_microtenant_controller

The **zpa_microtenant_controller** resource creates and microtenants in the Zscaler Private Access (ZPA) cloud. A Microtenant is created within a tenant and is used when departments or subsidiaries within an organization want to manage their configurations independently.

~> **NOTE:** This feature is in ``Limited Availability``. To learn more, contact Zscaler Support.

## Example Usage

```hcl
resource "zpa_microtenant_controller" "this" {
   name = "Microtenant_A"
   description = "Microtenant A"
   enabled = true
   criteria_attribute = "AuthDomain"
   criteria_attribute_values = "acme.com"
}

output "zpa_microtenant_controller" {
  value = zpa_microtenant_controller.this
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the microtenant to be created.
* `criteria_attribute` - (Required) The authentication type for this microtenant.
* `criteria_attribute_values` - (Required) The authentication domains from the drop-down menu to authenticate the users to the microtenant.

## Attributes Reference

* `description` - (Optional) This field defines the description of the microtenant.
* `enabled` - (Optional) This field defines the status of the microtenant. Supported values: `true`, `false`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Microtenant configuration can be imported by using `<MICROTENANT ID>` or `<MICROTENANT NAME>` as the import ID

For example:

```shell
terraform import zpa_microtenant_controller.example <microtenant_id>
```

or

```shell
terraform import zpa_microtenant_controller.example <microtenant_name>
```
