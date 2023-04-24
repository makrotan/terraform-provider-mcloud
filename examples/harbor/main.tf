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
  name = "harbor-test"
  instance_type = "cpx11"
  instance_count = 1
  location = "spread"
}

resource "mcloud_harbor" "test" {
  name = "harbor-test"
  sku = "dev"
  version = "2.7.0"
  fqdn = "harbor-test.makrotan.com"
  server_pool_id = mcloud_server_pool_hcloud.test.id
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Harbor
        Access: https://${mcloud_harbor.test.fqdn}/
        Username: admin
        Password: ${mcloud_harbor.test.admin_password}
        Docker Login CMD: docker login -u admin -p ${mcloud_harbor.test.admin_password} ${mcloud_harbor.test.fqdn}
        API Usage: curl -u "admin:${mcloud_harbor.test.admin_password}" https://${mcloud_harbor.test.fqdn}/api/v2.0/statistics

EOT
  sensitive = true
}
