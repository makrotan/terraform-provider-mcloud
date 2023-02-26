terraform {
  required_providers {
    mcloud = {
      version = "0.2.0"
      source = "makrotan/mcloud"
    }
  }
}

variable "mcloud_token" {}
variable "firewall_whitelist_ipv4" {
  default = ""
}

provider "mcloud" {
  api_token = var.mcloud_token
  host = "http://127.0.0.1:10004/"
}

resource "mcloud_server_pool_hcloud" "test" {
  name = "yuga-test"
  instance_type = "cpx31"
#  location = "spread"
  location = "nbg1"
  instance_count = 3
}

resource "mcloud_pki_ca" "test" {
  name = "yugabytedb-test"
}

resource "mcloud_consul_cluster" "test" {
  name = "consul-test"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
  pki_ca_id = mcloud_pki_ca.test.id
  version = "1.11.5"
}

resource "mcloud_yugabytedb" "test" {
  name = "yuga-test"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  version = "2.17.1.0"
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Consul
        Access: https://${mcloud_consul_cluster.test.ui_basic_auth_user}:${mcloud_consul_cluster.test.ui_basic_auth_password}@${mcloud_consul_cluster.test.master_domain}/ui/
        User: ${mcloud_consul_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_consul_cluster.test.ui_basic_auth_password}

EOT
  sensitive = true
}

