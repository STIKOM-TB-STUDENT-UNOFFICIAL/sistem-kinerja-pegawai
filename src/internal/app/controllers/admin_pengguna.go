package controllers

import (
	"kinerja-pegawai/internal/app/models"
	"kinerja-pegawai/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (ap *AdminPage) UserPage(ctx fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))

	if err != nil {
		utils.Log(err.Error())
	}

	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	total := ap.userModel.Count(query)
	user := ap.userModel.FindUser(limit, offset, query)

	return ctx.Render(
		"pages/admin_user_list",
		fiber.Map{
			"Title":      "Daftar Pengguna",
			"Level":      "admin",
			"Heading":    "Daftar Pengguna",
			"SubHeading": "Kelola aktivitas anda",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"User":       user,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) ViewAktivitas(ctx fiber.Ctx) error {
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
		result = ap.aktivitasModel.FindAktivitasByPeriode(ctx.Params("userid"), periode, limit, offset, query)
	} else {
		result = ap.aktivitasModel.FindAktivitas(ctx.Params("userid"), limit, offset, query)
	}

	total := ap.aktivitasModel.CountAktivitas(ctx.Params("userid"), query)
	periodeAll := ap.periodeModel.FindAll()

	return ctx.Render(
		"pages/admin_user_aktivitas",
		fiber.Map{
			"Title":      "Aktivitas Pengguna",
			"Level":      ctx.Locals("level"),
			"Heading":    "Dashboard Aktivitas",
			"SubHeading": "Riwayat aktivitas pengguna",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"Result":     result,
			"Pengguna":   ctx.Params("userid"),
			"PeriodeAll": periodeAll,
			"Periode":    periode,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) PostPengguna(ctx fiber.Ctx) error {
	userid := ctx.FormValue("userid")
	level := ctx.FormValue("level")
	nip := ctx.FormValue("nip")
	namaLengkap := ctx.FormValue("nama_lengkap")
	jabatan := ctx.FormValue("jabatan")
	departemen := ctx.FormValue("departemen")

	err := ap.userModel.Insert(userid, level, nip, namaLengkap, jabatan, departemen)

	if err != nil {
		utils.Log(err.Error())
	}

	return ctx.Redirect().To("/dashboard/admin/pengguna")
}

func (ap *AdminPage) UpdatePengguna(ctx fiber.Ctx) error {
	userid := ctx.FormValue("userid")
	level := ctx.FormValue("level")
	nip := ctx.FormValue("nip")
	namaLengkap := ctx.FormValue("nama_lengkap")
	jabatan := ctx.FormValue("jabatan")
	departemen := ctx.FormValue("departemen")

	err := ap.userModel.Upsert(userid, level, nip, namaLengkap, jabatan, departemen)

	if err != nil {
		utils.Log(err.Error())
	}

	return ctx.Redirect().To("/dashboard/admin/pengguna")
}

func (ap *AdminPage) DeletePengguna(ctx fiber.Ctx) error {
	userid := ctx.FormValue("userid")

	err := ap.userModel.Delete(userid)

	if err != nil {
		utils.Log(err.Error())
	}

	return ctx.Redirect().To("/dashboard/admin/pengguna")
}
