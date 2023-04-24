resource "mcloud_ip_block" "private" {
  name = "private"
}

resource "mcloud_ip" "mine" {
  name = "mine"
  block_id = mcloud_ip_block.private.id
  ip = "109.109.6.152"
}

resource "mcloud_ip_scope" "vault" {
  name = "vault"
}

resource "mcloud_ip_scope_block_assignment" "private__vault" {
  name = "private__vault"
  block_id = mcloud_ip_block.private.id
  scope_id = mcloud_ip_scope.vault.id
}

resource "mcloud_ip_scope_block_assignment" "vault__vault" {
  name = "vault__vault"
  block_id = mcloud_server_pool_hcloud.test.ip_block_id
  scope_id = mcloud_ip_scope.vault.id
}
