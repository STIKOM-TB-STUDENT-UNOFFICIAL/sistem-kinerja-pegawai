package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func (db *DashboardPage) PostKomentar(ctx fiber.Ctx) error {
	aktivitasId := ctx.Params("aktivitasId")
	anggota := ctx.Params("userid", "")
	user := ctx.Locals("userid").(string)
	komentar := ctx.FormValue("komentar", "")

	db.komentarAktivitasModel.InsertKomentar(user, aktivitasId, komentar)

	if anggota != "" {
		return ctx.Redirect().To(fmt.Sprintf("/dashboard/anggota/%s/%s", anggota, aktivitasId))
	}
	return ctx.Redirect().To(fmt.Sprintf("/dashboard/aktivitas/%s", aktivitasId))
}

func (db *DashboardPage) HapusKomentar(ctx fiber.Ctx) error {
	aktivitasId := ctx.Params("aktivitasId")
	anggota := ctx.Params("userid", "")
	komentarId := ctx.FormValue("id_komentar")

	db.komentarAktivitasModel.DeleteKomentar(komentarId)

	if anggota != "" {
		return ctx.Redirect().To(fmt.Sprintf("/dashboard/anggota/%s/%s", anggota, aktivitasId))
	}
	return ctx.Redirect().To(fmt.Sprintf("/dashboard/aktivitas/%s", aktivitasId))
}
