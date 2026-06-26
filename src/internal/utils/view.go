package utils

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
)

func ViewFile(ctx fiber.Ctx, file string) error {
	var viewsPath string

	if os.Getenv("ENVIRONMENT") == "dev" {
		viewsPath = "./internal/app/views"
	} else if os.Getenv("ENVIRONMENT") == "production" {
		viewsPath = "./views"
	} else {
		PanicLog("ENVIRONMENT Invalid Value !")
	}

	return ctx.SendFile(fmt.Sprintf("%s/%s.html", viewsPath, file))
}
