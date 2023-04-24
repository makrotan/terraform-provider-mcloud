---
page_title: "mcloud_consul_loadbalancer Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ConsulLoadbalancer(name, status, backup_ref, created, server_pool, ip_scope_admin)
---

# Resource `mcloud_consul_loadbalancer`

ConsulLoadbalancer(name, status, backup_ref, created, server_pool, ip_scope_admin)



## Argument Reference

The following arguments are supported:

- `ip_scope_admin_id` - [string]  
- `name` - (Required) [string]  
- `server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
