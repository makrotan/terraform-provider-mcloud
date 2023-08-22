---
page_title: "mcloud_container_app Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  ContainerApp(status, backup_ref, name, created, server_pool, definition, fqdn, app_port)
---

# Resource `mcloud_container_app`

ContainerApp(status, backup_ref, name, created, server_pool, definition, fqdn, app_port)



## Argument Reference

The following arguments are supported:

- `definition` - (Required) [string]  
- `fqdn` - (Required) [string]  
- `name` - (Required) [string]  
- `server_pool_id` - [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
