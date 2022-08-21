package env

import (
	"io"
	"text/template"
)

// Var is composite data to define env variable configuration behaviour
// Name is the name of the variable
// Method is one of Must, Fetch or Get
// Default is the fallback value when variable not exists, only applies to Fetch method
type Var struct {
	Name    string
	Method  string
	Default string
}

type input struct {
	Pkg  string
	Vars []Var
}

var templ = template.Must(template.New("qr").Parse(templateStr))

// Gen generates a go code which defines env vars as exported variables in a package for given package name
// example;
//
// Gen(os.Stdout, "main", []Var{{Name: "MY_VARIABLE", Method: "FETCH", Default: "fallback value"}})
func Gen(w io.Writer, pkg string, vars []Var) {
	templ.Execute(w, input{Pkg: pkg, Vars: vars})
}

const templateStr = `package {{.Pkg}}

import e "github.com/proctorexam/go/env"
{{range .Vars}}
var {{.Name}} string
{{- end}}

func init() {
{{- range .Vars -}}
{{- if .Default}}
	{{.Name}} = e.{{.Method}}("{{.Name}}","{{.Default}}")
{{- else}}
	{{.Name}} = e.{{.Method}}("{{.Name}}")
{{- end -}}
{{- end}}
}
`
