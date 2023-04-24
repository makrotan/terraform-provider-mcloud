resource "mcloud_ip_scope" "consul" {
  name = "consul"
}

resource "mcloud_ip_scope_block_assignment" "private__consul" {
  name = "private__consul"
  block_id = "private"
  scope_id = mcloud_ip_scope.consul.id
}


# vaco server pool
resource "mcloud_ip_scope_block_assignment" "consul__vaco" {
  name = "consul__${mcloud_server_pool_hcloud.vaco.ip_block_id}"
  block_id = mcloud_server_pool_hcloud.vaco.ip_block_id
  scope_id = mcloud_ip_scope.consul.id
}
# nomad server pool
resource "mcloud_ip_scope_block_assignment" "consul__nomad" {
  name = "consul__${mcloud_server_pool_hcloud.test.ip_block_id}"
  block_id = mcloud_server_pool_hcloud.test.ip_block_id
  scope_id = mcloud_ip_scope.consul.id
}
# client A server pool
resource "mcloud_ip_scope_block_assignment" "consul_client_a" {
  name = "consul__${mcloud_server_pool_hcloud.client_a.ip_block_id}"
  block_id = mcloud_server_pool_hcloud.client_a.ip_block_id
  scope_id = mcloud_ip_scope.consul.id
}

