package main

// Configuration options
type Configuration struct {
	usermail      string
	pass          string
	headerXappkey string
	// checkname     string // name of the check, ex summary.average
	checkid   string // id of the check, aka, which domain are we checking
	from      int32
	to        int32
	output    string
	mysqlurl  string // mysql connection in DSN (Data Source Name)
	inittable bool
	addcheck  bool
}

// Check if configuration is invalid
func (conf Configuration) configurationInvalid() bool {
	if conf.inittable {
		return conf.mysqlurl == ""
	}
	if conf.addcheck {
		return conf.mysqlurl == "" ||
			conf.checkid == ""
	}
	return conf.usermail == "" ||
		conf.pass == "" ||
		conf.headerXappkey == "" ||
		// conf.checkname == "" ||
		conf.checkid == "" ||
		conf.output == ""
}
