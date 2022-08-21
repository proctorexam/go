package main

import e "github.com/proctorexam/go/env"

var REQUIRED_VAR string
var WITH_DEFAULT string
var AS_IS string

func init() {
	REQUIRED_VAR = e.Must("REQUIRED_VAR")
	WITH_DEFAULT = e.Fetch("WITH_DEFAULT", "defval")
	AS_IS = e.Get("AS_IS")
}
