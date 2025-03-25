---
page_title: "zpa_pra_credential_controller Data Source - terraform-provider-zpa"
subcategory: "Privileged Remote Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-credentials
  API documentation https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api
  Get information about ZPA privileged remote access credential in Zscaler Private Access cloud.
---

# zpa_pra_credential_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-credentials)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-credentials-using-api)

The **zpa_pra_credential_controller** resource creates a privileged remote access credential in the Zscaler Private Access cloud. This resource can then be referenced in an privileged access policy resource.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
# Retrieves PRA Credential By Name
resource "zpa_pra_credential_controller" "this" {
    name = "John Doe"
}

# Retrieves PRA Credential By ID
resource "zpa_pra_credential_controller" "this" {
    name = "1234567890"
}
```

## Schema

### Required

* `name` - (String) The name of the privileged credential.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the privileged credential.
* `domain` - (Required) The description of the privileged credential.
* `credential_type` - (Required) The protocol type that was designated for that particular privileged credential. The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements. The supported values are:
    - ``USERNAME_PASSWORD``
    - ``SSH_KEY``
    - ``PASSWORD``
    
* `user_domain` - (string) - The domain name associated with the username. You can also include the domain name as part of the username. The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain.
* `username` - (string) - The username for the login you want to use for the privileged credential.
* `microtenant_id` (string) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.
