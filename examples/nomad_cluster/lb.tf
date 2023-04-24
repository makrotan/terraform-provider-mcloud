resource "mcloud_ip_scope" "lb_admin" {
  name = "lb_admin"
}

resource "mcloud_ip_scope_block_assignment" "private__lb" {
  name = "private__lb"
  block_id = "private"
  scope_id = mcloud_ip_scope.lb_admin.id
}


resource "mcloud_consul_loadbalancer" "lb" {
  name = "lb"
  server_pool_id = mcloud_server_pool_hcloud.client_a.id
  ip_scope_admin_id = mcloud_ip_scope.lb_admin.id
}
