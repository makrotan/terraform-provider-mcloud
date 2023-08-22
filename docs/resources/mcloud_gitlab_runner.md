---
page_title: "mcloud_gitlab_runner Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  GitlabRunner(status, backup_ref, name, created, server_pool, gitlab, version, meta, tags)
---

# Resource `mcloud_gitlab_runner`

GitlabRunner(status, backup_ref, name, created, server_pool, gitlab, version, meta, tags)



## Argument Reference

The following arguments are supported:

- `gitlab_id` - [string]  
- `meta` - [map]  
- `name` - (Required) [string]  
- `server_pool_id` - [string]  
- `status` - [string]   (default: `running`)
- `tags` - [string]  
- `version` - (Required) [string] Possible values: `16.1.1` 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
