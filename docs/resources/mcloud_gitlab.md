---
page_title: "mcloud_gitlab Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Gitlab(status, backup_ref, name, created, server_pool, admin_username, admin_password, fqdn, app_port, shared_runners_registration_token, meta, version)
---

# Resource `mcloud_gitlab`

Gitlab(status, backup_ref, name, created, server_pool, admin_username, admin_password, fqdn, app_port, shared_runners_registration_token, meta, version)



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `meta` - [map]  
- `name` - (Required) [string]  
- `server_pool_id` - [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `16.2.3-ce.0` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `admin_password` - [string] 
- `admin_username` - [string] 
- `shared_runners_registration_token` - [string] 