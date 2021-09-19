# Retrieve SCIM Attribute Header

This example will show you how to retrieve a SCIM Attribute Header ID to attach to a ZPA Access Policy Rule.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/scim-attribute-header-controller/getAllSCIMAttributes).

To run, configure your ZPA provider as described [Here](https://github.com/SecurityGeekIO/terraform-provider-zpa/blob/master/website/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```