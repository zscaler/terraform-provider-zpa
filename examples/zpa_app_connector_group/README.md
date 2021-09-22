# Retrieve App Connector Group

This example will show you how to retrieve an App Connector Group ID to attach to a ZPA Server Group
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/connector-group-controller/getAppConnectorGroup).

To run, configure your ZPA provider as described [Here](https://github.com/willguibr/terraform-provider-zpa/blob/master/website/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```