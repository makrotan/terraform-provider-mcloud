---
page_title: "mcloud_server_pool_hcloud Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ServerPoolHcloud(polymorphic_ctype, name, created, status, ip_block, serverpool_ptr, instance_type, instance_count, location, terraform_state)
---

# Resource `mcloud_server_pool_hcloud`

ServerPoolHcloud(polymorphic_ctype, name, created, status, ip_block, serverpool_ptr, instance_type, instance_count, location, terraform_state)



## Argument Reference

The following arguments are supported:

- `instance_count` - (Required) [number]  
- `instance_type` - (Required) [string] Possible values: `cx11`, `cpx11`, `cx21`, `cpx21` 
- `location` - (Required) [string] Possible values: `fsn1`, `nbg1`, `hel1`, `spread` 
- `name` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `ip_block_id` - [string] 