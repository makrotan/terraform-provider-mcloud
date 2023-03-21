---
page_title: "mcloud_pki_ca Resource - terraform-provider-mcloud"
subcategory: ""
description: |-
  PKI Certificate Authority.
---

# Resource `mcloud_pki_ca`

PKI Certificate Authority.



## Argument Reference

The following arguments are supported:

- `city` - [string]   (default: `Leipzig`)
- `country` - [string]   (default: `DE`)
- `email` - [string]   (default: `info@makrotan.com`)
- `name` - (Required) [string]  
- `organisation` - [string]   (default: `Makrotan`)
- `state` - [string]   (default: `DE`)
- `unit` - [string]   (default: `IT`)
- `valid_days` - [number]   (default: `7300`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `key_priv` - [string] 
- `key_pub` - [string] 