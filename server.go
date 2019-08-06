package main

import (
	"github.com/JewlyTwin/echo-restful-api/db"
	"github.com/JewlyTwin/echo-restful-api/route"
)

func main() {
	db.Init()
	e := route.Init()

	e.Logger.Fatal(e.Start(":3001"))
}
