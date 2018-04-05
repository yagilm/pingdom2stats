# What
Pulls data from Pingdom's API and pushes them to a mysql or postgres table.

# How to use
A database is needed with the appropriate rights.
Note: if you are using postgres, change `--mysqlurl` with `--pgurl` with the appropriate DSN format.
## DB configuration
You can use mysql or postgres with the appropriate flags. The DSN for each are:
- Mysql: "username:password@(address)/dbname"
- postgres: "postgres://username:password@address:port/dbname?sslmode=disable"
  - (sslmode can be disable, require, verify-ca, verify-full depending on your server configuration)

## DB table
On first run, the user must create the table. They can do that manually or by running:
`pingdom2mysql --inittable --mysqlurl="username:password@(address)/dbname"` (mysql)
or
`pingdom2mysql --inittable --pgurl="postgres://username:password@address:port/dbname?sslmode=disable"` (postgres)
in order to create the table (and check the DB connection)

## Adding checks
For every check that it's added, run with `--addcheck` in order to add the appropriate columns to the table.
**Attention!** For very big tables this might take some time.
`pingdom2mysql --addcheck --checkid=$YOUR_CHECK_ID --mysqlurl="username:password@(address)/dbname"`
Two new columns will be created with the name of the check and the check result fields. So your table will look like this:
```
mysql> describe summary_performances;
+-----------------------------+----------+------+-----+---------+-------+
| Field                       | Type     | Null | Key | Default | Extra |
+-----------------------------+----------+------+-----+---------+-------+
| timestamp                   | datetime | NO   | PRI | NULL    |       |
| $Name_of_check_avgresponse  | int(11)  | YES  |     | NULL    |       |
| $Name_of_check_downtime     | int(11)  | YES  |     | NULL    |       |
+-----------------------------+----------+------+-----+---------+-------+
7 rows in set (0.00 sec)
```
The program will use $checkid and pull the name of the check which will use for naming the columns(`checkname.go`).
Run the --addcheck multiple times to add multiple checks.

## Populating statistics
Add it to a job scheduler like cron or chronos. I prefer to run it every 20 hours (it pulls the last 24 hours' statistics). Note that if the timestamp exists the program will just update the values. No double timestamps are possible in the datastore.

- The program should be used like this:
```
 pingdom2mysql --appkey=$YOURAPPKEY --checkid=$CHECKID --email=$ACCOUNTMAIL --pass=$ACCOUNTPASSWORD --mysqlurl="$DBUSER:$DBPASS@($DBIP:$DBPORT)/$DBNAME" --output="mysql"
```
You might want to use `--output="console"` first to see the data that will end up in your database.

## Running in docker
For running it inside docker create the docker image by running
`make pingdom2mysql-docker` and then run it like
```
docker run --rm pingdom2mysql --appkey=$YOURAPPKEY --checkid=$CHECKID --email=$ACCOUNTMAIL --pass=$ACCOUNTPASSWORD  --mysqlurl="$DBUSER:$DBPASS@($DBIP:$DBPORT)/$DBNAME" --output="mysql"
```

## Pulling the historical data from Pingdom
In order to pull the historical data and write it in the data store, you need to run `fetch_history $UNIX_TIMESTAMP_OF_CHECK_CREATION`.
You would need to add the configuration variables, lines 4-9 of `fetch_history`.

# Usage information
```
./pingdom2mysql  --help
Using Pingdom's API as described in: https://www.pingdom.com/resources/api
Version: v0.2.2-4-g690a40b
Usage: pingdom2mysql [options]
All options are required (but some have defaults):
  --addcheck
        Add new check into the mysql table, requires --mysqlurl, --checkid
  --appkey string
        Appkey for pingdom's API
  --checkid string
        ID of the check, aka the domain are we checking.
  --email string
        Pingdom's API configured e-mail account
  --from value
        from which (Unix)time we are asking, default 24 hours ago which is  (default 1522841253)
  --inittable
        Initialize the table, requires --mysqlurl
  --mysqlurl string
        mysql connection in DSN, like: username:password@(address)/dbname.
        Cannot use together with --pgurl
  --output string
        Output destination (console, db) (default "console")
  --pass string
        password for pingdom's API
  --pgurl string
        postgres connection in DSN, like: postgres://username:password@address:port/dbname?sslmode=disable.
        Cannot use together with --mysqlurl
  --to value
        until which (Unix)time we are asking, default now which is  (default 1522927653)
```
