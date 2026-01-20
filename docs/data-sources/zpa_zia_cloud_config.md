---
page_title: "zpa_zia_cloud_config Data Source - terraform-provider-zpa"
subcategory: "Customer Config Controller"
description: |-
  Retrieve Zscaler Cloud Sandbox settings
---

# zpa_zia_cloud_config (Data Source)

The **zpa_zia_cloud_config** retrieve configures the Zscaler Cloud Sandbox Settings in the Zscaler Private Access cloud.

**NOTE** Passwords are not returned in the API response.

## Example Usage

```hcl
data "zpa_zia_cloud_config" "this" {

}
```

## Schema

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