---
page_title: "mcloud_erpnext Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Erpnext(name, status, backup_ref, sku, fqdn, admin_password, created, version, app_port, server_pool)
---

# Resource `mcloud_erpnext`

Erpnext(name, status, backup_ref, sku, fqdn, admin_password, created, version, app_port, server_pool)



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `server_pool_id` - [string]  
- `sku` - (Required) [string] Possible values: `dev` 
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `14.13.0`, `14.18.1` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `admin_password` - [string] 