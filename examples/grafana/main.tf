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

resource "mcloud_pki_ca" "test" {
  name = "mimir-test"
}

resource "mcloud_grafana_mimir" "test" {
  name = "mimir-test"
  fqdn = "mimir-test.makrotan.com"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  version = "2.6.0"
}

resource "mcloud_grafana_loki" "test" {
  name = "loki-test"
  fqdn = "loki-test.makrotan.com"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  version = "2.7.4"
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Grafana:
        Access: https://${mcloud_grafana.test.fqdn}/
        User: admin
        Password: ${mcloud_grafana.test.admin_password}

    Mimir:
        Access: https://${mcloud_grafana_mimir.test.fqdn}/
        Prometheus Remote Write URL: https://${mcloud_grafana_mimir.test.fqdn}/api/v1/push
        User: ${mcloud_grafana_mimir.test.basic_auth_user}
        Password: ${mcloud_grafana_mimir.test.basic_auth_password}

    Loki:
        Access: https://${mcloud_grafana_loki.test.fqdn}/
        User: ${mcloud_grafana_loki.test.basic_auth_user}
        Password: ${mcloud_grafana_loki.test.basic_auth_password}

EOT
  sensitive = true
}

