package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Authentication(ctx fiber.Ctx) error {
	sess := session.FromContext(ctx)

	if sess == nil {
		return ctx.Redirect().Status(302).To("/login")
	}

	if sess.Get("is_auth") != true {
		return ctx.Redirect().Status(302).To("/login")
	}

	ctx.Locals("level", sess.Get("level"))
	ctx.Locals("is_auth", sess.Get("is_auth"))
	ctx.Locals("userid", sess.Get("userid"))
	ctx.Locals("nama_lengkap", sess.Get("nama_lengkap"))

	return ctx.Next()
}
