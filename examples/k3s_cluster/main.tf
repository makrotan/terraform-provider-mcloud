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
  name = "k3s_test"
  instance_type = "cpx11"
  instance_count = 1
  location = "nbg1"
}

resource "mcloud_k3s_cluster" "test" {
  name = "foo"
  sku = "dev"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  version = "v1.23.1+k3s2"
  firewall_whitelist_ipv4 = "109.109.6.152,92.79.101.164"
}

resource "local_file" "test" {
  content  = mcloud_k3s_cluster.test.k3s_config_yaml
  filename =  pathexpand("~/.kube/config-mcloud-dev-${mcloud_k3s_cluster.test.name}.yml")
  file_permission      = "0600"
  directory_permission = "0700"
}

output "out" {
  value = <<EOT
Resources successfully installed:

    K3sCluster

EOT
  sensitive = true
}
