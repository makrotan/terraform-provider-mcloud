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
  name = "vault-test"
  instance_type = "cpx11"
  instance_count = 1
  location = "spread"
}


resource "mcloud_pki_ca" "test" {
  name = "vault-test"
}

resource "mcloud_vault_cluster" "test" {
  name = "vault-test"
  version = "1.12.2-1"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  ip_scope = mcloud_ip_scope.vault.id
}

output "out" {
  value = <<EOT
Resources successfully installed:

    Vault
        Access: https://${mcloud_vault_cluster.test.ui_basic_auth_user}:${mcloud_vault_cluster.test.ui_basic_auth_password}@${mcloud_vault_cluster.test.master_domain}/ui/
        User: ${mcloud_vault_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_vault_cluster.test.ui_basic_auth_password}
        Root-Token: ${mcloud_vault_cluster.test.root_token}

EOT
  sensitive = true
}


#│ Error: error applying jobspec: Put "https://vaule-test.vaule.makrotan.com/v1/jobs?region=europe": x509:
#                                              vaule-test.vaule.makrotan.com
#
#certificate is valid for vaule-vaule-test-admin, localhost, not vaule-test.vaule.makrotan.com
