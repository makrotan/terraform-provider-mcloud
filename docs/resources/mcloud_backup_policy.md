---
page_title: "mcloud_backup_policy Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  BackupPolicy(name, status, ref, keep_backup_count, schedule, scheduler_job)
---

# Resource `mcloud_backup_policy`

BackupPolicy(name, status, ref, keep_backup_count, schedule, scheduler_job)



## Argument Reference

The following arguments are supported:

- `keep_backup_count` - [number]   (default: `7`)
- `name` - (Required) [string]  
- `ref` - (Required) [string]  
- `schedule` - [string]  
- `status` - [string]   (default: `running`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `scheduler_job_id` - [string] 