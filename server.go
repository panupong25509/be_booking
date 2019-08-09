package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/route"
)

func main() {
	db.Init()
	e := route.Init()
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":3001"))
}
