---
page_title: "mcloud_nomad_server_pool Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  NomadServerPool(status, backup_ref, name, nomad_cluster, server_pool)
---

# Resource `mcloud_nomad_server_pool`

NomadServerPool(status, backup_ref, name, nomad_cluster, server_pool)



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `nomad_cluster_id` - (Required) [string]  
- `server_pool_id` - (Required) [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
