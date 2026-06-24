
# Create Private Cloud Resource

This example will show you how to create a Private Cloud Resource
This example codifies [this API](https://help.zscaler.com/legacy-apis/private-cloud-controller-management-zpa-api-reference#/mgmtconfig/v1/admin/customers/{customerId}/privateCloudController-get).

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
