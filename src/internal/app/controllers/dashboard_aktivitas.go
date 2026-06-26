package controllers

import (
	"kinerja-pegawai/internal/app/models"
	"kinerja-pegawai/internal/utils"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (db *DashboardPage) UserAktivitasPage(ctx fiber.Ctx) error {
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
		result = db.aktivitasModel.FindAktivitasByPeriode(ctx.Locals("userid").(string), periode, limit, offset, query)
	} else {
		result = db.aktivitasModel.FindAktivitas(ctx.Locals("userid").(string), limit, offset, query)
	}

	total := db.aktivitasModel.CountAktivitas(ctx.Locals("userid").(string), query)
	allTupoksi := db.tupoksiModel.FindAllTupoksi(ctx.Locals("userid").(string))
	periodeAll := db.periodeModel.FindAll()

	return ctx.Render(
		"pages/user_dashboard",
		fiber.Map{
			"Title":      "Aktivitas",
			"Level":      ctx.Locals("level"),
			"Heading":    "Aktivitas",
			"SubHeading": "Kelola aktivitas anda",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"Result":     result,
			"TupoksiAll": allTupoksi,
			"PeriodeAll": periodeAll,
			"Periode":    periode,
		},
		"layouts/main",
	)
}

func (db *DashboardPage) UpdateAktivitas(ctx fiber.Ctx) error {
	idAktivitas := ctx.FormValue("id_aktivitas")
	tupoksi := ctx.FormValue("tupoksi")
	catatan := ctx.FormValue("catatan")

	var mulai *string
	var selesai *string

	if v := ctx.FormValue("startdate"); v != "" {
		mulai = &v
	}

	if v := ctx.FormValue("enddate"); v != "" {
		selesai = &v
	}

	db.aktivitasModel.UpdateAktivitas(
		idAktivitas,
		ctx.Locals("userid").(string),
		tupoksi,
		catatan,
		mulai,
		selesai,
	)

	return ctx.Redirect().To("/dashboard/aktivitas")
}

func (db *DashboardPage) DeleteAktivitas(ctx fiber.Ctx) error {
	idAktivitas := ctx.FormValue("id_aktivitas")

	files := db.fileBuktiModel.FindAllByAktivitasId(idAktivitas)

	for _, file := range *files {
		_ = os.Remove(file.LokasiFile)
	}

	db.fileBuktiModel.DeleteByAktivitas(idAktivitas)
	db.aktivitasModel.DeleteAktivitas(idAktivitas)

	return ctx.Redirect().To("/dashboard/aktivitas")
}

func (db *DashboardPage) CheckSelesai(ctx fiber.Ctx) error {
	userid := ctx.Locals("userid").(string)
	aktivitasId := ctx.Params("aktivitasId")

	db.aktivitasModel.UpdateStatusAktivitas(aktivitasId, userid, "Diproses")
	db.aktivitasModel.SetSelesai(aktivitasId, userid)

	return ctx.Redirect().To("/dashboard/aktivitas")
}
