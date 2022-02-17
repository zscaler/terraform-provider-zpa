data "zpa_enrollment_cert" "root" {
    name = "Root"
}
data "zpa_enrollment_cert" "client" {
    name = "Client"
}
data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}
data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}
data "zpa_enrollment_cert" "isolation_client" {
    name = "Isolation Client"
}