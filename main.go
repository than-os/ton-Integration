package main

import (
	"github.com/fatih/color"
	"github.com/jasonlvhit/gocron"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/than-os/sent-dante/dbo"
	"github.com/than-os/sent-dante/services"
)

var d = TON{}

func main() {


	s := gocron.NewScheduler()
	go func() {
		s.Every(5).Seconds().Do(service.KeepAlive)
		color.Red("%s", "running the job")
		<-s.Start()
	}()
	e := echo.New()

	//middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE},
	}))

	e.GET("/", service.RootFunc)
	e.GET("/all", service.GetAllTonNodes)
	e.POST("/register", service.RegisterTonNode)
	e.POST("/node", service.GetTonNodeByIP)
	e.DELETE("/node", service.RemoveTonNode)

	//Start the server
	e.Start(":30001")
}

func init() {
	d.Server = "localhost"
	d.Database = "ton"

	d.NewSession()
}
