---
page_title: "mcloud_mattermost Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Mattermost(name, status, backup_ref, sku, fqdn, secret_key, created, version, app_port, server_pool, postgres_username, postgres_password)
---

# Resource `mcloud_mattermost`

Mattermost(name, status, backup_ref, sku, fqdn, secret_key, created, version, app_port, server_pool, postgres_username, postgres_password)



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `postgres_password` - [string]  
- `postgres_username` - [string]  
- `server_pool_id` - [string]  
- `sku` - (Required) [string] Possible values: `dev` 
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `7.1`, `7.7`, `7.9.2` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `secret_key` - [string] 