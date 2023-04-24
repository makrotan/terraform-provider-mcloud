terraform {
  required_providers {
    mcloud = {
      version = "0.1.0"
      source = "makrotan.com/cloud/mcloud"
    }
  }
}

variable "mcloud_username" {}
variable "mcloud_password" {}

provider "mcloud" {
  username = var.mcloud_username
  password = var.mcloud_password
  host = "http://127.0.0.1:8000/"
}

resource mcloud_ssh_key "foo" {
  name = "foo2"
}

resource "mcloud_server_pool_hcloud" "foo" {
  name = "foo4"
  instance_type = "cpx21"
  instance_count = 2
  location = "spread"
}

resource "mcloud_k3s_cluster" "foo" {
  name = "foo"
  sku = "dev"
  master_server_pool_id = mcloud_server_pool_hcloud.foo.id
  k3s_version = "v1.23.1+k3s2"
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
}
#
#
output "pubkey" {
  value = mcloud_ssh_key.foo.public_key
}
#output "private_key" {
#  value = mcloud_ssh_key.foo.private_key
#  sensitive = true
#}
