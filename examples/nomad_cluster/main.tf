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
  name = "nomad-test"
  instance_type = "cpx11"
  instance_count = 1
}

resource "mcloud_server_pool_hcloud" "client_a" {
  name = "nomad-client-A"
  instance_type = "cpx11"
  instance_count = 1
}

resource "mcloud_pki_ca" "test" {
  name = "nomad-test"
}

resource "mcloud_consul_cluster" "test" {
  name = "harbor-test"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
  pki_ca_id = mcloud_pki_ca.test.id
  version = "1.11.5"
}


resource "mcloud_nomad_cluster" "test" {
  name = "nomad-test"
  version = "1.4.3-1"
  consul_cluster_id = mcloud_consul_cluster.test.id
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
}

resource "mcloud_nomad_server_pool" "client_a" {
  name = "nomad-client-A"
  nomad_cluster_id = mcloud_nomad_cluster.test.id
  server_pool_id = mcloud_server_pool_hcloud.client_a.id
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Consul
        Access: https://${mcloud_consul_cluster.test.ui_basic_auth_user}:${mcloud_consul_cluster.test.ui_basic_auth_password}@${mcloud_consul_cluster.test.master_domain}/ui/
        User: ${mcloud_consul_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_consul_cluster.test.ui_basic_auth_password}

    Nomad
        Access: https://${mcloud_nomad_cluster.test.ui_basic_auth_user}:${mcloud_nomad_cluster.test.ui_basic_auth_password}@${mcloud_nomad_cluster.test.master_domain}/ui/
        User: ${mcloud_nomad_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_nomad_cluster.test.ui_basic_auth_password}
EOT
  sensitive = true
}
