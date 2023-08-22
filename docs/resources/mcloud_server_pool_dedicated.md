---
page_title: "mcloud_server_pool_dedicated Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ServerPoolDedicated(polymorphic_ctype, name, created, status, ip_block, description, consul_cluster, servers, total_memory, total_cpu, total_disk, total_price_per_month, serverpool_ptr)
---

# Resource `mcloud_server_pool_dedicated`

ServerPoolDedicated(polymorphic_ctype, name, created, status, ip_block, description, consul_cluster, servers, total_memory, total_cpu, total_disk, total_price_per_month, serverpool_ptr)



## Argument Reference

The following arguments are supported:

- `consul_cluster_id` - [string]  
- `description` - [string]  
- `name` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `ip_block_id` - [string] 
- `servers` - [number] 
- `total_cpu` - [number] 
- `total_disk` - [number] 
- `total_memory` - [number] 
- `total_price_per_month` - [number] 