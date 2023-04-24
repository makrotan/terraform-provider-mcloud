resource "mcloud_ip_scope" "nomad" {
  name = "nomad"
}

resource "mcloud_ip_scope_block_assignment" "private__nomad" {
  name = "private__nomad"
  block_id = "private"
  scope_id = mcloud_ip_scope.nomad.id
}


# nomad server pool
resource "mcloud_ip_scope_block_assignment" "nomad__nomad" {
  name = "${mcloud_server_pool_hcloud.test.ip_block_id}__nomad"
  block_id = mcloud_server_pool_hcloud.test.ip_block_id
  scope_id = mcloud_ip_scope.nomad.id
}

# client A server pool
resource "mcloud_ip_scope_block_assignment" "client_a__nomad" {
  name = "${mcloud_server_pool_hcloud.client_a.ip_block_id}__nomad"
  block_id = mcloud_server_pool_hcloud.client_a.ip_block_id
  scope_id = mcloud_ip_scope.nomad.id
}

