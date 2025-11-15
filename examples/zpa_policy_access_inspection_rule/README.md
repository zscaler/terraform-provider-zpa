# Policy Access Inspection Rule Example

This example will show you how to use Terraform to implement a ZPA policy access inspection rule resource.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/policy-set-controller).

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
