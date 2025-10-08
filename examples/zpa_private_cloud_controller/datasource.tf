data "zpa_private_cloud_controller" "foo" {
  name = "DataCenter"
}


## Example Usage - Search by ID

data "zpa_private_cloud_controller" "foo" {
  id = "123456789"
}
