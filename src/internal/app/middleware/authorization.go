package middleware

import (
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Authorization(level []string, path string) func(ctx fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		sess := session.FromContext(ctx)

		if sess == nil {
			return ctx.Redirect().Status(302).To("/login")
		}

		if !slices.Contains(level, sess.Get("level").(string)) {
			return ctx.Redirect().Status(302).To(fmt.Sprintf("%s/%s", path, sess.Get("level")))
		}

		return ctx.Next()
	}
}
