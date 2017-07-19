package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// CheckNameResponse describes the response of
// https://api.pingdom.com/api/2.0/checks/$checkid
type CheckNameResponse struct {
	Check struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"check"`
}

// getCheckName extracts the name of a checkid
// replacing spaces with _
func getCheckName() string {
	req, err := http.NewRequest("GET",
		fmt.Sprintf(
			"https://api.pingdom.com/api/2.0/checks/%s",
			Config.checkid), nil)
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(Config.usermail, Config.pass)
	req.Header.Set("app-key", Config.headerXappkey)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	if res.StatusCode != http.StatusOK {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var responseName CheckNameResponse
	err = json.Unmarshal(body, &responseName)
	if err != nil {
		panic(err.Error())
	}
	checkidName := strings.Replace(responseName.Check.Name, " ", "_", -1)
	return checkidName
}
