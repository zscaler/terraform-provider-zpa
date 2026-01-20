---
page_title: "zpa_zia_cloud_config Resource - terraform-provider-zpa"
subcategory: "Customer Config Controller"
description: |-
  Configure Zscaler Cloud Sandbox Settings
---

# zpa_zia_cloud_config (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-user-portals)
* [API documentation](https://help.zscaler.com/zpa/about-user-portals)

The **zpa_zia_cloud_config** resource configures the Zscaler Cloud Sandbox Settings in the Zscaler Private Access cloud. This resource is required when configuring the policy type resource `zpa_policy_capabilities_rule`.

## Example Usage

```hcl
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zpa_zia_cloud_config" "this" {
  zia_username              = ""
  zia_password              = ""
  zia_cloud_service_api_key = ""
  zia_sandbox_api_token     = ""
  zia_cloud_domain          = ""

}
```

## Schema

### Required

* `zia_username` - (String) The ZIA admin username with permission to use the api key
* `zia_password` - (String) The ZIA admin password with permission to use the api key
* `zia_cloud_service_api_key` - (String) The ZIA Cloud service api key
* `zia_sandbox_api_token` - (String) The ZIA Sandbox API token
* `zia_cloud_domain` - (String) The supported ZIA cloud name. Supported values are: 

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_zia_cloud_config** can be imported by using `<ZIA_CLOUD_CONFIG>` as the import ID.

For example:

```shell
terraform import zpa_zia_cloud_config.example <zia_cloud_config>
```

or

```shell
terraform import zpa_zia_cloud_config.example <zia_cloud_config>
```