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
  instance_count = 2
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
  name = "consul-test"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
  pki_ca_id = mcloud_pki_ca.test.id
  version = "1.11.5"
}

resource "mcloud_vault_cluster" "test" {
  name = "vault-test"
  version = "1.12.2-1"
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
}

resource "mcloud_nomad_cluster" "test" {
  name = "nomad-test"
  version = "1.4.3-1"
  consul_cluster_id = mcloud_consul_cluster.test.id
  master_server_pool_id = mcloud_server_pool_hcloud.test.id
  pki_ca_id = mcloud_pki_ca.test.id
  firewall_whitelist_ipv4 = var.firewall_whitelist_ipv4
  vault_cluster_id = mcloud_vault_cluster.test.id
}

resource "mcloud_nomad_server_pool" "client_a" {
  name = "nomad-client-A"
  nomad_cluster_id = mcloud_nomad_cluster.test.id
  server_pool_id = mcloud_server_pool_hcloud.client_a.id
}

provider "nomad" {
  address = "https://${mcloud_nomad_cluster.test.master_domain}"
  region  = "europe"
  http_auth = "${mcloud_nomad_cluster.test.ui_basic_auth_user}:${mcloud_nomad_cluster.test.ui_basic_auth_password}"
  ca_pem = "${mcloud_nomad_cluster.test.admin_ca}"
  cert_pem = "${mcloud_nomad_cluster.test.admin_cert}"
  key_pem = "${mcloud_nomad_cluster.test.admin_key}"
}

#resource "nomad_job" "app" {
#  jobspec = <<EOT
#job "docs" {
#  datacenters = ["dc1"]
#
#  group "example" {
#    vault {
#      policies  = ["fooo"]
#    }
#
#    network {
#      port "http" {
#        to = "80"
##        static = "8080"
#      }
#    }
#    task "server" {
#      driver = "docker"
#
#      config {
#        image = "nginx"
#        ports = ["http"]
#        volumes = ["html:/usr/share/nginx/html"]
##        args = [
##          "-listen",
##          ":5678",
##          "-text",
##          "hello world",
##        ]
#      }
#
#      template {
#        data   = <<EOF
#my secret: "{{ with secret "kv/foo" }}{{ .Data.data.bar }}{{ end }}"
#EOF
##        destination = "$${NOMAD_TASK_DIR}/usr/share/nginx/html/index.html"
#        destination = "html/index.html"
#
#      }
#    }
#  }
#}
#
#EOT
#}

output "out" {
  value = <<EOT
Resources successfully installed:

    Consul
        Access: https://${mcloud_consul_cluster.test.ui_basic_auth_user}:${mcloud_consul_cluster.test.ui_basic_auth_password}@${mcloud_consul_cluster.test.master_domain}/ui/
        User: ${mcloud_consul_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_consul_cluster.test.ui_basic_auth_password}

    Vault
        Access: https://${mcloud_vault_cluster.test.ui_basic_auth_user}:${mcloud_vault_cluster.test.ui_basic_auth_password}@${mcloud_vault_cluster.test.master_domain}/ui/
        User: ${mcloud_vault_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_vault_cluster.test.ui_basic_auth_password}
        Root-Token: ${mcloud_vault_cluster.test.root_token}

    Nomad
        Access: https://${mcloud_nomad_cluster.test.ui_basic_auth_user}:${mcloud_nomad_cluster.test.ui_basic_auth_password}@${mcloud_nomad_cluster.test.master_domain}/ui/
        User: ${mcloud_nomad_cluster.test.ui_basic_auth_user}
        Password: ${mcloud_nomad_cluster.test.ui_basic_auth_password}
EOT
  sensitive = true
}


#│ Error: error applying jobspec: Put "https://nomad-test.nomad.makrotan.com/v1/jobs?region=europe": x509:
#                                              nomad-test.nomad.makrotan.com
#
#certificate is valid for nomad-nomad-test-admin, localhost, not nomad-test.nomad.makrotan.com
