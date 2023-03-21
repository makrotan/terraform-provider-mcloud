---
page_title: "mcloud_grafana Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  Grafana(backup_ref, name, created, server_pool, admin_password, fqdn, app_port, status, version)
---

# Resource `mcloud_grafana`

Grafana(backup_ref, name, created, server_pool, admin_password, fqdn, app_port, status, version)



## Argument Reference

The following arguments are supported:

- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `server_pool_id` - [string]  
- `status` - [string]   (default: `running`)
- `version` - (Required) [string] Possible values: `9.3.6` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `admin_password` - [string] 