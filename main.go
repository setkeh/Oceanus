package main

import (
	"os"

	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	middleware "github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/setkeh/Oceanus/db"
	"github.com/setkeh/Oceanus/handlers"
	"github.com/sirupsen/logrus"
)

const (
	//url = "http://127.0.0.1:1323"
	url = "https://screenshots.local.setkeh.com"
)

func main() {
	// DB Instance
	d := db.Client{
		Path: "images.db",
	}

	d.Open()
	// Echo instance
	e := echo.New()

	// Logger
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.TextFormatter{})
	e.Logger = log.Logger()
	e.Use(middleware.Logger())
	log.Info("Logger enabled!!")

	// Routes
	handlers.DB = &d
	handlers.URL = url
	e.GET("/", handlers.Hello)
	e.POST("/image", handlers.PostImageHandler)
	e.GET("/image", handlers.GetImageHandler)
	e.GET("/images", handlers.GetImageListHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
