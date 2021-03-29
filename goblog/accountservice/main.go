package main

import (
	"fmt"

	"github.com/ratnadeepb/goblog/accountservice/dbclient"
	"github.com/ratnadeepb/goblog/accountservice/service"
)

var appName = "accountservice"

func main() {
	fmt.Printf("Starting: %v\n", appName)
	initialiseBoltClient()
	service.StartWebServer("6767")
}

func initialiseBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
