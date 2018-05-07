package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection
var EPG_DB *pop.Connection
var EPG_DB_RO *pop.Connection // read-only EPG database

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}

	EPG_DB, err = pop.Connect("epg_" + env)
	if err != nil {
		log.Fatal(err)
	}

	EPG_DB_RO, err = pop.Connect("epg_ro_" + env)
	if err != nil {
		log.Fatal(err)
	}

	pop.Debug = env == "development"
}
