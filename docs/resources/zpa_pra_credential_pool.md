---
page_title: "zpa_pra_credential_pool Resource - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-credentials
  API documentation https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api
  Creates and manages ZPA privileged remote access credential pool
---

# zpa_pra_credential_pool (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-credential-pools)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api)

The **zpa_pra_credential_pool** resource creates a privileged remote access credential pool in the Zscaler Private Access cloud. This resource can then be referenced in an privileged access policy resource.

## Example Usage

```terraform
# Creates Credential Pool of Type "USERNAME_PASSWORD"

resource "zpa_pra_credential_pool" "this" {
  name            = "PRACredentialPool01"
  credential_type = "USERNAME_PASSWORD"
  credentials {
    id = [zpa_pra_credential_controller.this.id]
  }
}

resource "zpa_pra_credential_controller" "this" {
  name            = "John Doe"
  description     = "Created with Terraform"
  credential_type = "PASSWORD"
  user_domain     = "acme.com"
  password        = ""
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) The name of the privileged credential.
- `domain` - (String) The description of the privileged credential.
- `credential_type` - (String) `USERNAME-PASSWORD` is for `RDP` and `SSH` machines. `SSH Key` is for Private Key based SSH machines. `PASSWORD` is for `VNC` machines. The supported values are:
    - ``USERNAME_PASSWORD``
    - ``SSH_KEY``
    - ``PASSWORD``

⚠️ **WARNING:**: The `credential_type` attribute cannot be updated once created.

* `credentials` - (Required)
  * `id` - (Required) The ID of each individual pra credential user to be associated with the pool.

### Optional

In addition to all arguments above, the following attributes are exported:

- `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_pra_credential_pool** can be imported by using `<POOL ID>` or `<POOL NAME>` as the import ID.

For example:

```shell
terraform import zpa_pra_credential_pool.this <pool_id>
```

or

```shell
terraform import zpa_pra_credential_pool.this <pool_name>
```
