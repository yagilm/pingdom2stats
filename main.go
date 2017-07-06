package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	flag "github.com/ogier/pflag"
)

// Lasttimeserie is the global store to remember which is the last timeserie we
// sent to statsd. When not assigned, it is '0001-01-01 00:00:00 +0000 UTC',
// thus earlier of any possible timeserie
// --will not need if ES?
// var Lasttimeserie time.Time

// Configuration options
type Configuration struct {
	usermail      string
	pass          string
	headerXappkey string
	checkname     string // name of the check, ex summary.average
	checkid       string // id of the check, aka, which domain are we checking
}

// Config keeps the configuration
var Config Configuration

// Check if configuration is invalid
func (conf Configuration) configurationInvalid() bool {
	return conf.usermail == "" ||
		conf.pass == "" ||
		conf.headerXappkey == "" ||
		conf.checkname == "" ||
		conf.checkid == ""
}

// Response describes the parts we want from cloudflare's json response
type Response struct {
	Summary struct {
		// Depending on the data store see if I need to do smth with the time, if i
		// send it as unix time, i might as well send the int
		// https://stackoverflow.com/questions/24987131/how-to-parse-unix-timestamp-in-golang
		Responsetime struct {
			Avgresponse int `json:"avgresponse"`
			From        int `json:"from"`
			To          int `json:"to"`
		} `json:"responsetime"`
	} `json:"summary"`
	Status struct {
		Totaldown    int `json:"totaldown"`
		Totalunknown int `json:"totalunknown"`
		Totalup      int `json:"totalup"`
	}
}

// init() runs before the main function as described in:
// https://golang.org/doc/effective_go.html#init
func init() {
	flag.StringVar(&Config.usermail, "email", "", "e-mail account for pingdom's API")
	flag.StringVar(&Config.pass, "pass", "", "password for pingdom's API")
	flag.StringVar(&Config.headerXappkey, "appkey", "", "Appkey for pingdom's API")
	flag.StringVar(&Config.checkname, "checkname", "", "Name of the check (eg summary.average)") //multiple checks seperated by comma?
	flag.StringVar(&Config.checkid, "checkid", "", "id of the check, which domain are we checking?")
	flag.Usage = func() {
		fmt.Printf("Usage: pingdom2-- [options]\nRequired options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if Config.configurationInvalid() {
		flag.Usage()
		os.Exit(1)
	}
}

// Gets data from Pingdom's API
func getpingdomdata() (*Response, error) {
	// make the request with the appropriate headers
	req, err := http.NewRequest("GET",
		fmt.Sprintf(
			"https://api.pingdom.com/api/2.0/%s/%s?from=1499330970&includeuptime=true", //TODO Add: ?from=$(date -d '1 minute ago' +"%s")\&includeuptime=true
			Config.checkname,
			Config.checkid),
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
	log.Println(string(body))
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// Send statistics to the data store
// func sendstatistics() {
// I think I will be using olivere/elastic https://github.com/olivere/elastic for this to send it to ES
// Nevertheless, we need to decide on a data store....
// }

func main() {
	timer := time.NewTicker(time.Second * 2)
	// infinite loop
	for {
		res, err := getpingdomdata()
		if err != nil {
			fmt.Print("Something went wrong requesting the json in the API:", err)
		}
		toprint, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(string(toprint))
		<-timer.C
	}
}
