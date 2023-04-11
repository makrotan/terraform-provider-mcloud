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

resource "mcloud_ip_block" "private" {
  name = "private"
}
resource "mcloud_ip" "mine" {
  name = "mine"
  block_id = mcloud_ip_block.private.id
  ip = "109.109.6.152"
}

resource "mcloud_ip" "netcupqa1" {
  name = "netcup-qa-1"
  block_id = mcloud_ip_block.private.id
  ip = "37.120.185.180"
}

resource "mcloud_server_pool_hcloud" "test" {
  name = "yuga-test"
  instance_type = "cpx21"
#  location = "spread"
  location = "nbg1"
  instance_count = 3
}

resource "mcloud_pki_ca" "test" {
  name = "yugabytedb-test"
  algorithm = "rsa-2048"
}

resource "mcloud_consul_cluster" "test" {
  name = "consul-test"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
  pki_ca_id = mcloud_pki_ca.test.id
  version = "1.11.5"
}

resource "mcloud_ip_scope" "yuga" {
  name = "yuga"
}

resource "mcloud_ip_scope_block_assignment" "private__yuga" {
  name = "private__yuga"
  block_id = mcloud_ip_block.private.id
  scope_id = mcloud_ip_scope.yuga.id
}

resource "mcloud_yugabytedb" "test" {
  name = "yuga-test"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  version = "2.17.1.0"
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
#  consul_cluster_id = mcloud_consul_cluster.test.id
  ip_scope_client_id = mcloud_ip_scope.yuga.id
  ip_scope_admin_id = mcloud_ip_scope.yuga.id
}

resource "mcloud_yugabytedb_database" "test" {
  name = "test"
  yugabytedb_id = mcloud_yugabytedb.test.id
}

resource "mcloud_pki_cert" "test" {
  ca_id = mcloud_pki_ca.test.id
  common_name = "test"
  name = "test"
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Consul
        Access: https://${mcloud_consul_cluster.test.ui_basic_auth_user}:${mcloud_consul_cluster.test.ui_basic_auth_password}@${mcloud_consul_cluster.test.master_domain}/ui/
        User: ${mcloud_consul_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_consul_cluster.test.ui_basic_auth_password}

    Yugabyte:
        FQDN: ${mcloud_yugabytedb.test.fqdn}
        DB: ${mcloud_yugabytedb_database.test.name}
            User: ${mcloud_yugabytedb_database.test.username}
            Password: ${mcloud_yugabytedb_database.test.password}
            Pubkey:
                ${mcloud_pki_cert.test.key_pub}
            Privkey:
                ${mcloud_pki_cert.test.key_priv}
            CA:
                ${mcloud_pki_ca.test.key_pub}


EOT
  sensitive = true
}

