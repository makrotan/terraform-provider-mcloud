---
page_title: "mcloud_yugabytedb_database Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  YugabytedbDatabase(status, backup_ref, name, yugabytedb, password, created)
---

# Resource `mcloud_yugabytedb_database`

YugabytedbDatabase(status, backup_ref, name, yugabytedb, password, created)



## Argument Reference

The following arguments are supported:

- `name` - (Required) [string]  
- `status` - [string]   (default: `running`)
- `yugabytedb_id` - [string]  

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `password` - [string] 
- `username` - [string] 