---
page_title: "zpa_pra_credential_controller Resource - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-credentials
  API documentation https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api
  Creates and manages ZPA privileged remote access credential
---

# zpa_pra_credential_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-credentials)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api)

The **zpa_pra_credential_controller** resource creates a privileged remote access credential in the Zscaler Private Access cloud. This resource can then be referenced in an privileged access policy resource.

## Example Usage

```terraform
#### PASSWORDS OR RELATED CREDENTIALS ATTRIBUTES IN THIS FILE ARE FOR EXAMPLE ONLY AND NOT USED IN PRODUCTION SYSTEMS ####
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


```terraform
#### PASSWORDS OR RELATED CREDENTIALS ATTRIBUTES IN THIS FILE ARE FOR EXAMPLE ONLY AND NOT USED IN PRODUCTION SYSTEMS ####
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
-----END PRIVATE KEY-----
    EOT
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) The name of the privileged credential.
- `domain` - (String) The description of the privileged credential.
- `credential_type` - (String) The protocol type that was designated for that particular privileged credential. The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements. The supported values are:
    - ``USERNAME_PASSWORD``
    - ``SSH_KEY``
    - ``PASSWORD``

⚠️ **WARNING:**: The resource `credential_type` and associated attributes cannot be updated once created.

- `user_domain` - (String) - The domain name associated with the username. You can also include the domain name as part of the username. The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain.
- `username` - (String) - The username for the login you want to use for the privileged credential.

### Optional

In addition to all arguments above, the following attributes are exported:

- `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

**pra_credential_controller** can be imported by using `<CREDENTIAL ID>` or `<CREDENTIAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_pra_credential_controller.this <credential_id>
```

or

```shell
terraform import zpa_pra_credential_controller.this <credential_name>
```
