# App Connector Group Example

This example will show you how to use Terraform to implement the ZPA App Connector Group resource
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/connector-group-controller/getAppConnectorGroup).

To run, configure your ZPA provider as described [Here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```

## Destroy 💥

```bash
terraform destroy
```
