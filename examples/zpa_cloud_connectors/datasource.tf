data "zpa_cloud_connector_group" "aws_cloud_connector" {
  name = "AWS Cloud Connector"
}

output "get_cloud_connector_group" {
  value = data.zpa_cloud_connector_group.aws_cloud_connector.id
}
