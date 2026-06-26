package controllers

import (
	"kinerja-pegawai/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (db *DashboardPage) UserTupoksiPage(ctx fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))

	if err != nil {
		utils.Log(err.Error())
	}

	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	result := db.tupoksiModel.FindTupoksi(ctx.Locals("userid").(string), limit, offset, query)
	total := db.tupoksiModel.Count(ctx.Locals("userid").(string), query)

	return ctx.Render(
		"pages/user_tupoksi",
		fiber.Map{
			"Title":      "Tupoksi",
			"Level":      ctx.Locals("level"),
			"Heading":    "Tupoksi",
			"SubHeading": "Kelola tupoksi",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"Result":     result,
		},
		"layouts/main",
	)
}

func (db *DashboardPage) PostTupoksi(ctx fiber.Ctx) error {
	nama := ctx.FormValue("nama")
	deskripsi := ctx.FormValue("deskripsi")

	db.tupoksiModel.Insert(ctx.Locals("userid").(string), nama, deskripsi)

	return ctx.Redirect().To("/dashboard/tupoksi")
}

func (db *DashboardPage) UpdateTupoksi(ctx fiber.Ctx) error {
	id := ctx.FormValue("id")
	nama := ctx.FormValue("nama")
	deskripsi := ctx.FormValue("deskripsi")

	db.tupoksiModel.Update(id, nama, deskripsi)

	return ctx.Redirect().To("/dashboard/tupoksi")
}

func (db *DashboardPage) DeleteTupoksi(ctx fiber.Ctx) error {
	id := ctx.FormValue("id")

	db.tupoksiModel.Delete(id)

	return ctx.Redirect().To("/dashboard/tupoksi")
}
