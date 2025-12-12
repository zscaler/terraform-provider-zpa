# Inspection Custom Control Example

This example will show you how to use Terraform to implement a ZPA application segment resource.
This example codifies [this API](https://help.zscaler.com/zpa/inspection-control-controller#/mgmtconfig/v1/admin/customers/{customerId}/inspectionControls/custom-post).

To run, configure your ZPA provider as described [Here](https://github.com/SecurityGeekIO/terraform-provider-zpa/blob/master/docs/index.md)

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
