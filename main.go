package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	flag "github.com/ogier/pflag"
)

// Lasttimeserie is the global store to remember which is the last timeserie we
// sent to statsd. When not assigned, it is '0001-01-01 00:00:00 +0000 UTC',
// thus earlier of any possible timeserie
// --will not need if ES?
// var Lasttimeserie time.Time

// Config keeps the configuration
var Config Configuration
var version = "development"

// Response describes the parts we want from cloudflare's json response
type Response struct {
	Summary struct {
		Hours []struct {
			Avgresponse int       `json:"avgresponse"`
			Downtime    int       `json:"downtime"`
			Starttime   Timestamp `json:"starttime"`
			Uptime      int       `json:"uptime"`
		} `json:"hours"`
	} `json:"summary"`
}

// init() runs before the main function as described in:
// https://golang.org/doc/effective_go.html#init
func init() {
	flag.StringVar(&Config.usermail, "email", "", "Pingdom's API configured e-mail account")
	flag.StringVar(&Config.pass, "pass", "", "password for pingdom's API")
	flag.StringVar(&Config.headerXappkey, "appkey", "", "Appkey for pingdom's API")
	// flag.StringVar(&Config.checkname, "checkname", "", "Name of the check (eg summary.performance)") //multiple checks seperated by comma?
	flag.StringVar(&Config.checkid, "checkid", "", "ID of the check, aka the domain are we checking.")
	flag.Int32Var(&Config.from, "from", int32(time.Now().Add(-24*time.Hour).Unix()), "from which (Unix)time we are asking, default 24 hours ago which is ")
	flag.Int32Var(&Config.to, "to", int32(time.Now().Unix()), "until which (Unix)time we are asking, default now which is ")
	flag.StringVar(&Config.output, "output", "console", "Output destination (console, mysql)")
	flag.StringVar(&Config.mysqlurl, "mysqlurl", "", "mysql connection in DSN, like: username:password@(address)/dbname")
	flag.BoolVar(&Config.inittable, "inittable", false, "Initialize the table, requires --mysqlurl ")

	flag.Usage = func() {
		fmt.Println("Using Pingdom's API as described in: https://www.pingdom.com/resources/api")
		fmt.Printf("Version: %s\nUsage: pingdom2mysql [options]\nAll options are required (but some have defaults):\n", version)
		flag.PrintDefaults()
	}
	flag.Parse()
	if Config.configurationInvalid() {
		flag.Usage()
		os.Exit(1)
	}
}

// Gets data from Pingdom's API
func getPingdomData() (*Response, error) {
	// make the request with the appropriate headers
	req, err := http.NewRequest("GET",
		fmt.Sprintf(
			"https://api.pingdom.com/api/2.0/summary.performance/%s?from=%d&to=%d&includeuptime=true", //TODO Add: ?from=$(date -d '1 minute ago' +"%s")\&includeuptime=true
			// Config.checkname,
			Config.checkid,
			Config.from,
			Config.to),
		nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(Config.usermail, Config.pass)
	req.Header.Set("app-key", Config.headerXappkey)
	req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("API not 200")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	response.Summary.Hours = response.Summary.Hours[:len(response.Summary.Hours)-1]
	// fmt.Println(response.Summary.Hours[0].Starttime.String())
	return &response, err
}

// Send statistics to the data store
// func sendstatistics() {
// I think I will be using olivere/elastic https://github.com/olivere/elastic for this to send it to ES
// Nevertheless, we need to decide on a data store....
// }
func consoleOutput(res *Response) error {
	for _, hour := range res.Summary.Hours {
		fmt.Println(hour.Starttime.Time, hour.Uptime, hour.Avgresponse, hour.Downtime)
	}
	return nil
}
func connectToDB() *sql.DB {
	db, err := sql.Open("mysql", Config.mysqlurl)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}

func sendToMysql(res *Response) error {
	db := connectToDB()
	defer db.Close()

	// cupcake9:bi-db:$TABLE:
	// Here is the structure of the table:
	// id | timestamp | avg_uptime | avg_responcetime |
	// check if scheme is correct?
	return nil
}
func initializeTable() {
	db := connectToDB()
	defer db.Close()
	var err error
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS summary_performances (timestamp DATETIME PRIMARY KEY);`)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	res, err := getPingdomData()
	if err != nil {
		log.Panicln("Something went wrong requesting the json in the API:", err)
	}
	if Config.inittable {
		initializeTable()
		os.Exit(0)
	}
	if Config.output == "console" {
		consoleOutput(res)
	} else if Config.output == "mysql" {
		sendToMysql(res)
	}
}
