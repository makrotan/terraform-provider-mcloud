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

resource "mcloud_server_pool_hcloud" "foo" {
  name = "cockroachdb-test"
  instance_type = "cpx11"
  instance_count = 1
  location = "spread"
}

resource "mcloud_pki_ca" "test" {
  name = "cockroachdb-test"
}

resource "mcloud_consul_cluster" "test" {
  name = "cockroachdb-test"
  master_server_pool_id = mcloud_server_pool_hcloud.foo.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
  pki_ca_id = mcloud_pki_ca.test.id
  version = "1.15.2"
}

resource "mcloud_cockroachdb" "test" {
  name = "cockroachdb-test"
  version = "22.1.0"
  pki_ca_id = mcloud_pki_ca.test.id
  consul_cluster_id = mcloud_consul_cluster.test.id
  server_pool_id = mcloud_server_pool_hcloud.foo.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Consul
        Access: https://${mcloud_consul_cluster.test.ui_basic_auth_user}:${mcloud_consul_cluster.test.ui_basic_auth_password}@${mcloud_consul_cluster.test.master_domain}/ui/
        User: ${mcloud_consul_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_consul_cluster.test.ui_basic_auth_password}

    CockroachDB:
        Access: https://${mcloud_cockroachdb.test.ui_basic_auth_user}:${mcloud_cockroachdb.test.ui_basic_auth_password}@${mcloud_cockroachdb.test.master_domain}/
        User: ${mcloud_cockroachdb.test.ui_basic_auth_user}
        Password: ${mcloud_cockroachdb.test.ui_basic_auth_password}

EOT
  sensitive = true
}
