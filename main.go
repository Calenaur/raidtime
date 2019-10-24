package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/foolin/echo-template"
	"github.com/calenaur/raidtime/db"
	"github.com/calenaur/raidtime/store"
	"github.com/calenaur/raidtime/config"
	"github.com/calenaur/raidtime/handler"
)

func main() {
	//Load config
	cfg, err := config.Load("config.json")
    if err != nil {
		panic(err)
    }

	//Connect to database
	d, err := db.New(cfg.Database.Username, cfg.Database.Password, cfg.Database.Database)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	//Setup stores
	userStore := store.NewUserStore(d, cfg)
	eventStore := store.NewEventStore(d, cfg)
	discordStore := store.NewDiscordStore(cfg)

	//Setup handler
	handler := handler.New(userStore, eventStore, discordStore, cfg)

	//Setup echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = echotemplate.Default()

	//Routes
	handler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}