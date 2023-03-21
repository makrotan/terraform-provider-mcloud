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

- `instance_count` - (Required) [number]  
- `instance_type` - (Required) [string]  
- `location` - [string]   (default: `spread`)
- `name` - (Required) [string]  
- `status` - [string] `new`, `running`, `failed`, `deleting` or `deleted`  (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
