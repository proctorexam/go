package main

import (
	"fmt"

	_ "github.com/proctorexam/go/env"
)

//go:generate go run gitlab.com/proctorexam1/go/env/cmd main REQUIRED_VAR,Must WITH_DEFAULT,Fetch,defval AS_IS,Get

func main() {
	fmt.Println(REQUIRED_VAR, WITH_DEFAULT, AS_IS)
}
