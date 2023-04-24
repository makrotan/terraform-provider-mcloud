resource "mcloud_vault_cluster" "test" {
  name = "vault-test"
  version = "1.12.2-1"
  master_server_pool_id = mcloud_server_pool_hcloud.vaco.id
  pki_ca_id = mcloud_pki_ca.test.id
  ip_scope_id = mcloud_ip_scope.vault.id
}
