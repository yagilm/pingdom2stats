# What
Pulls data from Pingdom's API and pushes them to a data store (Elastic Search?)

# How
All options are required:
```
Usage: pingdom2-- [options]
Required options:
  --appkey string
        Appkey for pingdom's API
  --checkid string
        id of the check, which domain are we checking?
  --checkname string
        Name of the check (eg summary.average)
  --email string
        e-mail account for pingdom's API
  --pass string
        password for pingdom's API
```
