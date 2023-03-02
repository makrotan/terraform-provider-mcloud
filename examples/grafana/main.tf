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
  name = "grafana-test"
  instance_type = "cx11"
  location = "nbg1"
  instance_count = 1
}

resource "mcloud_grafana" "test" {
  name = "grafana-test"
  fqdn = "grafana-test.makrotan.com"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  version = "9.3.6"
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Grafana:
        Access: https://${mcloud_grafana.test.fqdn}/
        User: admin
        Password: ${mcloud_grafana.test.admin_password}

EOT
  sensitive = true
}

