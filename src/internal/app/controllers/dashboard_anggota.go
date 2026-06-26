package controllers

import (
	"fmt"
	"kinerja-pegawai/internal/app/models"
	"kinerja-pegawai/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (db *DashboardPage) AnggotaPage(ctx fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))

	if err != nil {
		utils.Log(err.Error())
	}

	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	result := db.anggotaModel.FindAllAnggota(ctx.Locals("userid").(string), query, limit, offset)
	total := db.anggotaModel.Count(ctx.Locals("userid").(string), query)

	return ctx.Render("pages/user_anggota", fiber.Map{
		"Title":      "Anggota",
		"Level":      ctx.Locals("level"),
		"Heading":    "Anggota",
		"SubHeading": "Kelola anggota",
		"Page":       page,
		"Previous":   page - 1,
		"Next":       page + 1,
		"HasNext":    (total - page*limit) > 0,
		"Offset":     offset,
		"Query":      query,
		"Result":     result,
	}, "layouts/main")
}

func (db *DashboardPage) UserAnggotaAktivitasPage(ctx fiber.Ctx) error {
	anggota := ctx.Params("userid", "")
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	periode := ctx.Query("periode", "")

	if err != nil {
		utils.Log(err.Error())
	}

	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	var result *[]models.AktivitasDB

	if periode != "" {
		result = db.aktivitasModel.FindAktivitasByPeriode(anggota, periode, limit, offset, query)
	} else {
		result = db.aktivitasModel.FindAktivitas(anggota, limit, offset, query)
	}

	total := db.aktivitasModel.CountAktivitas(anggota, query)
	dataAnggota := db.anggotaModel.FindAnggotaByUserId(anggota, ctx.Locals("userid").(string))
	periodeAll := db.periodeModel.FindAll()

	return ctx.Render(
		"pages/user_anggota_view",
		fiber.Map{
			"Title":      "Aktivitas Anggota",
			"Level":      ctx.Locals("level"),
			"Heading":    "Aktivitas Anggota",
			"SubHeading": "Kelola aktivitas anggota",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"Result":     result,
			"Anggota":    anggota,
			"Nama":       dataAnggota.Nama,
			"Periode":    periode,
			"PeriodeAll": periodeAll,
		},
		"layouts/main",
	)
}

func (db *DashboardPage) UpdateAktivitasAnggota(ctx fiber.Ctx) error {
	idAktivitas := ctx.FormValue("id_aktivitas")
	status := ctx.FormValue("status")

	db.aktivitasModel.UpdateStatusAktivitas(idAktivitas, ctx.Params("userid", ""), status)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/anggota/%s", ctx.Params("userid", "")))
}
