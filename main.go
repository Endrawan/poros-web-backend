package main

import (
	"github.com/divisi-developer-poros/poros-web-backend/config"
	"github.com/divisi-developer-poros/poros-web-backend/migrations"
	"github.com/divisi-developer-poros/poros-web-backend/routes"
)

var dbModel config.DBModel

func main() {
	db := dbModel.PostgreConn()
	migrations.Start(db)
	routes.Start()
}
