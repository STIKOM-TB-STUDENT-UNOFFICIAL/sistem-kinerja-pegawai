package main

import (
	"fmt"
	"kinerja-pegawai/internal/system"
	"kinerja-pegawai/internal/utils"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
)

func main() {
	system.LoadConfig()

	var viewsPath string
	var reload bool
	var prefork bool

	if os.Getenv("ENVIRONMENT") == "dev" {
		viewsPath = "./internal/app/views"
		reload = true
	} else if os.Getenv("ENVIRONMENT") == "production" {
		viewsPath = "./views"
		reload = false
	} else {
		utils.PanicLog("ENVIRONMENT Invalid Value !")
	}

	if os.Getenv("PREFORK") == "true" {
		prefork = true
	} else {
		prefork = false
	}

	engine := html.New(viewsPath, ".html")
	engine.Reload(reload)
	engine.AddFunc("gt", func(a, b int) bool {
		return a > b
	})

	engine.AddFunc("gte", func(a, b int) bool {
		return a >= b
	})

	engine.AddFunc("lt", func(a, b int) bool {
		return a < b
	})

	engine.AddFunc("lte", func(a, b int) bool {
		return a <= b
	})

	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	db, err := system.OpenDatabase()

	if err != nil {
		utils.PanicLog(err.Error())
	}

	app.Use(session.New())
	app.Use("/*", static.New("./public"))

	system.SetupRoutes(app, db)

	err = app.Listen(
		fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		fiber.ListenConfig{
			EnablePrefork: prefork,
		},
	)

	if err != nil {
		utils.PanicLog(err.Error())
	}
}
