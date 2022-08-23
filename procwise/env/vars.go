package env

import e "github.com/proctorexam/go/env"

var ENV = e.Fetch("ENV", "development")
var PE_USER string
var PE_PASSWORD string

func init() {
	switch ENV {
	case "production":
		PE_USER = e.Must("PE_USER")
		PE_PASSWORD = e.Must("PE_PASSWORD")
	default:
		PE_USER = e.Fetch("PE_USER", "super@proctorexam.com")
		PE_PASSWORD = e.Fetch("PE_PASSWORD", "_Def4ultP4ssw0rd_")
	}
}
