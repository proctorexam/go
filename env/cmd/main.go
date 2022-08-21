package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/proctorexam/go/env"
)

var dir = flag.String("dir", "", "dirname to create env var file")

// usage: go run github.com/proctorexam/go/env/cmd main VAR1,Must VAR2,Fetch,default VAR3,Get
func main() {
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pkg := os.Args[1]

	if *dir != "" {
		os.MkdirAll(path.Join(cwd, *dir), os.ModePerm)
	}

	file, err := os.OpenFile(path.Join(cwd, *dir, "vars.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	vars := make([]env.Var, 0)

	for _, a := range os.Args[2:] {
		s := strings.Split(a, ",")
		v := env.Var{Name: s[0], Method: s[1]}
		if len(s) == 3 {
			v.Default = s[2]
		}
		vars = append(vars, v)
	}

	env.Gen(file, pkg, vars)
}
