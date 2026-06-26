package controllers

import (
	"fmt"
	"kinerja-pegawai/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (ap *AdminPage) PeriodeList(ctx fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))

	if err != nil {
		utils.Log(err.Error())
	}

	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	result := ap.periodeModel.FindPeriode(limit, offset, query)
	total := ap.periodeModel.Count(query)

	return ctx.Render(
		"pages/admin_periode_list",
		fiber.Map{
			"Title":      "Daftar Periode",
			"Level":      "admin",
			"Heading":    "Periode",
			"SubHeading": "Kelola periode",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"Periode":    result,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) PeriodeCreatePage(ctx fiber.Ctx) error {
	return ctx.Render(
		"pages/admin_periode_form",
		fiber.Map{
			"Title":   "Tambah Periode",
			"Level":   "admin",
			"Heading": "Tambah Periode",
			"Periode": nil,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) PeriodeCreate(ctx fiber.Ctx) error {
	nama := ctx.FormValue("nama_periode")
	start := ctx.FormValue("start_date")
	end := ctx.FormValue("end_date")

	ap.periodeModel.Insert(nama, start, end)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/periode?q=%s&page=%s", ctx.FormValue("query"), ctx.FormValue("page")))
}

func (ap *AdminPage) PeriodeEditPage(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	periode := ap.periodeModel.FindById(id)

	return ctx.Render(
		"pages/admin_periode_form",
		fiber.Map{
			"Title":   "Edit Periode",
			"Level":   "admin",
			"Heading": "Edit Periode",
			"Periode": periode,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) PeriodeUpdate(ctx fiber.Ctx) error {
	id := ctx.FormValue("id")
	nama := ctx.FormValue("nama_periode")
	start := ctx.FormValue("start_date")
	end := ctx.FormValue("end_date")

	ap.periodeModel.Update(id, nama, start, end)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/periode?q=%s&page=%s", ctx.FormValue("query"), ctx.FormValue("page")))
}

func (ap *AdminPage) PeriodeDelete(ctx fiber.Ctx) error {
	id := ctx.FormValue("id")
	ap.periodeModel.Delete(id)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/periode?q=%s&page=%s", ctx.FormValue("query"), ctx.FormValue("page")))
}
