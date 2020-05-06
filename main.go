package main

import (
	"os"

	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	middleware "github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

const (
	url = "http://localhost"
)

func main() {
	// DB Instance
	d := db.dbClient{
		Path: "images.db",
	}

	d.dbOpen()

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
	e.GET("/", handlers.hello)
	e.POST("/image", handlers.postImageHandler)
	e.GET("/image", handlers.getImageHandler)
	e.GET("/images", handlers.getImageListHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
