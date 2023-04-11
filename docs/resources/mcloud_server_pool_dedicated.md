---
page_title: "mcloud_server_pool_dedicated Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ServerPoolDedicated(polymorphic_ctype, name, created, status, ip_block, serverpool_ptr)
---

# Resource `mcloud_server_pool_dedicated`

ServerPoolDedicated(polymorphic_ctype, name, created, status, ip_block, serverpool_ptr)



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `ip_block_id` - [string] 