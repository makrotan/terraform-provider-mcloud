---
page_title: "mcloud_erpnext Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ERPNext Instance.
---

# Resource `mcloud_erpnext`

ERPNext Instance.



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `fqdn` - (Required) [string]  
- `sku` - (Required) [string]  
- `version` - (Required) [string]  
- `server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `admin_password` - [string] 