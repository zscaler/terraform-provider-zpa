---
subcategory: "Privileged Remote Access"
layout: "zscaler"
page_title: "ZPA): pra_credential_controller"
description: |-
  Creates and manages ZPA privileged remote access credential
---

# Resource: zpa_pra_credential_controller

The **zpa_pra_credential_controller** resource creates a privileged remote access credential in the Zscaler Private Access cloud. This resource can then be referenced in an privileged access policy resource.

## Example Usage

```hcl
# Creates Credential of Type "USERNAME_PASSWORD"
resource "zpa_pra_credential_controller" "this" {
    name = "John Doe"
    description = "Created with Terraform"
    credential_type = "USERNAME_PASSWORD"
    user_domain = "acme.com"
    username = "jdoe"
    password = ""
}
```

```hcl
# Creates Credential of Type "SSH_KEY"
resource "zpa_pra_credential_controller" "this" {
    name = "John Doe"
    description = "Created with Terraform"
    credential_type = "SSH_KEY"
    user_domain = "acme.com"
    username = "jdoe"
    private_key = <<-EOT
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDEjc8pPoobS0l6
KjldVtViVwqCTXZZOjHnmEIMn+XKU3sEYlqDKLp6TByIsBlITKd3Ju8qMBNwXcfi
-----END PRIVATE KEY-----
    EOT
}
```

## Attributes Reference

### Required

* `name` - (Required) The name of the privileged credential.
* `domain` - (Required) The description of the privileged credential.
* `credential_type` - (Required) The protocol type that was designated for that particular privileged credential. The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements. The supported values are:
    - ``USERNAME_PASSWORD``
    - ``SSH_KEY``
    - ``PASSWORD``

⚠️ **WARNING:**: The resource `credential_type` and associated attributes cannot be updated once created.

* `user_domain` - (Required) - The domain name associated with the username. You can also include the domain name as part of the username. The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain.
* `username` - (Required) - The username for the login you want to use for the privileged credential.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**pra_credential_controller** can be imported by using `<CREDENTIAL ID>` or `<CREDENTIAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_pra_credential_controller.this <credential_id>
```

or

```shell
terraform import zpa_pra_credential_controller.this <credential_name>
```
