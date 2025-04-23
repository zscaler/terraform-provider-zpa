---
page_title: "zpa_pra_credential_pool Data Source - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-credentials
  API documentation https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api
  Get information about ZPA privileged remote access credential pool in Zscaler Private Access cloud.
---

# zpa_pra_credential_pool (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-credential-pools)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api)

The **zpa_pra_credential_pool** data source to get information about a privileged remote access credential pool created in the Zscaler Private Access cloud.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

# Retrieves PRA Credential Pool By Name
```terraform
data "zpa_pra_credential_pool" "this" {
  name = "PRACredentialPool01"
}
```

# Retrieves PRA Credential Pool By ID
```terraform
data "zpa_pra_credential_pool" "this" {
  id = "5458"
}
```

## Schema

### Required

* `name` - (String) The name of the privileged credential.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the privileged credential pool.
- `credential_type` - (String) `USERNAME-PASSWORD` is for `RDP` and `SSH` machines. `SSH Key` is for Private Key based SSH machines. `PASSWORD` is for `VNC` machines. The supported values are:
    - ``USERNAME_PASSWORD``
    - ``SSH_KEY``
    - ``PASSWORD``

* `credentials` - (List)
  * `id` - (List of String) The ID of each individual pra credential user to be associated with the pool.

- `microtenant_id` (String) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

- `microtenant_name` (String) The name of the Microtenant.
- `credential_mapping_count` (String)