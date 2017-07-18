# What
Pulls data from Pingdom's API and pushes them to a mysql table

# Prerequisites
We need a DB table to write the data into.. Since there would probably be many different checks (--checkid), I am thinking of table columns like `$checkid_avg`, `$checkid_down` etc..

# Usage information
```
Using Pingdom's API as described in: https://www.pingdom.com/resources/api
Version: 0.1
Usage: pingdom2mysql [options]
All options are required (but some have defaults):
  --appkey string
        Appkey for pingdom's API
  --checkid string
        ID of the check, aka the domain are we checking.
        Note: for multiple checks, run pingdom2mysql multiple times!
  --checkname string
        Name of the check (eg summary.performance)
  --email string
        Pingdom's API configured e-mail account
  --from value
        from which (Unix)time we are asking, default 24 hours ago which is (default 1500287566)
  --pass string
        password for pingdom's API
  --to value
        until which (Unix)time we are asking, default now which is  (default 1500373966)
```
