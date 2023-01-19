---
page_title: "mcloud_server_pool_hcloud Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  A collection of dedicated servers.
---

# Resource `mcloud_server_pool_hcloud`

A collection of dedicated servers.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `instance_type` - (Required) [string]  
- `location` - [string]   (default: `spread`)
- `instance_count` - (Required) [number]  
- `status` - [string] `new`, `running`, `failed`, `deleting` or `deleted`  (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
