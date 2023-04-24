resource "mcloud_ip_scope" "vault" {
  name = "vault"
}

resource "mcloud_ip_scope_block_assignment" "private__vault" {
  name = "private__vault"
  block_id = "private"
  scope_id = mcloud_ip_scope.vault.id
}


# vaco server pool
resource "mcloud_ip_scope_block_assignment" "vault__vaco" {
  name = "vault__${mcloud_server_pool_hcloud.vaco.ip_block_id}"
  block_id = mcloud_server_pool_hcloud.vaco.ip_block_id
  scope_id = mcloud_ip_scope.vault.id
}
# nomad server pool
resource "mcloud_ip_scope_block_assignment" "vault__nomad" {
  name = "vault__${mcloud_server_pool_hcloud.test.ip_block_id}"
  block_id = mcloud_server_pool_hcloud.test.ip_block_id
  scope_id = mcloud_ip_scope.vault.id
}
# client A server pool
resource "mcloud_ip_scope_block_assignment" "vault_client_a" {
  name = "vault__${mcloud_server_pool_hcloud.client_a.ip_block_id}"
  block_id = mcloud_server_pool_hcloud.client_a.ip_block_id
  scope_id = mcloud_ip_scope.vault.id
}
