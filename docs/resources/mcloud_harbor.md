---
page_title: "mcloud_harbor Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Harbor Instance.
---

# Resource `mcloud_harbor`

Harbor Instance.



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `server_pool_id` - (Required) [string]  
- `sku` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `admin_password` - [string] 