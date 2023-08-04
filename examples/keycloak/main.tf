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
  name = "keycloak-test"
  instance_type = "cx11"
  location = "nbg1"
  instance_count = 1
}

resource "mcloud_pki_ca" "test" {
  name = "keycloak-test"
}

resource "mcloud_keycloak" "test" {
  name = "keycloak-test"
  fqdn = "keycloak-test.makrotan.com"
  server_pool_id = mcloud_server_pool_hcloud.test.id
  version = "21.0.1"
  sku = "dev"
  pki_ca_id = mcloud_pki_ca.test.id

#  themes = jsonencode({
#    mytheme = {
#      url = "https://gitlab.com/path/to/theme.git"
#      username = "gitlab+deploy-token-changeme"
#      password = "changeme"
#    }
#  })

  postgres = jsonencode({
    host = "yuga-prod.yuga.makrotan.com"
    port = 5433
    database = "keycloak_mium"
    username = "keycloak_mium"
    password = "changeme"
    ssl_cert = file("yugabytedb-client-mcloud.pub.pem")
    ssl_key = file("yugabytedb-client-mcloud.priv.pem")
    ssl_ca = file("yugabytedb-ca.pem")
  })
}

output "out" {
  value = <<EOT
Resources successfully installed:

    keycloak:
        Access: https://${mcloud_keycloak.test.fqdn}/admin
        User: admin
        Password: ${mcloud_keycloak.test.admin_password}

EOT
  sensitive = true
}

