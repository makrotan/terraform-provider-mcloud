```
resource mcloud_ssh_key "foo" {
  name = "foo1"
}

output "public_key" {
  value = mcloud_ssh_key.foo.public_key
}
output "private_key" {
  value = mcloud_ssh_key.foo.private_key
  sensitive = true
}
```

