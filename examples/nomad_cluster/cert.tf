resource mcloud_pki_cert "test" {
  ca_id = mcloud_pki_ca.test.id
  name = "test"
  common_name = "test"
}

output "cert" {
  value = mcloud_pki_cert.test.key_pub
}
output "key" {
  value = mcloud_pki_cert.test.key_priv
  sensitive = true
}
output "ca" {
  value = mcloud_pki_ca.test.key_pub
}
