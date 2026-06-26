package controllers

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/app/models"
	"kinerja-pegawai/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type AdminPage struct {
	anggotaModel   *models.AnggotaModel
	periodeModel   *models.PeriodeModel
	userModel      *models.UserModel
	aktivitasModel *models.AktivitasModel
}

func (ap *AdminPage) DashboardAdminRingkasan(ctx fiber.Ctx) error {
	totalPengguna := ap.userModel.CountTotal()
	totalPeriode := ap.periodeModel.CountTotal()

	return ctx.Render(
		"pages/admin_home",
		fiber.Map{
			"Title":         "Dashboard Admin",
			"Level":         "admin",
			"Heading":       "Dashboard",
			"SubHeading":    "Ringkasan data sistem",
			"TotalPengguna": totalPengguna,
			"TotalPeriode":  totalPeriode,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) AnggotaAdminPage(ctx fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))

	if err != nil {
		utils.Log(err.Error())
	}

	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	result := ap.anggotaModel.FindParent(query, limit, offset)
	total := ap.anggotaModel.CountParent(query)
	children := ap.anggotaModel.FindParentChildren(query, limit, offset)
	allParent := ap.anggotaModel.AllParent()

	return ctx.Render(
		"pages/admin_dashboard",
		fiber.Map{
			"Title":      "Anggota",
			"Level":      "admin",
			"Heading":    "Anggota",
			"SubHeading": "Kelola aktivitas anda",
			"Page":       page,
			"Previous":   page - 1,
			"Next":       page + 1,
			"HasNext":    (total - page*limit) > 0,
			"Offset":     offset,
			"Query":      query,
			"Parent":     result,
			"Children":   children,
			"AllParent":  allParent,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) InsertAnggota(ctx fiber.Ctx) error {
	parentId := ctx.FormValue("parent_id")
	childrenId := ctx.FormValue("children_id")

	ap.anggotaModel.InsertAnggota(parentId, childrenId)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/anggota?q=%s&page=%s", ctx.FormValue("query"), ctx.FormValue("page")))
}

func (ap *AdminPage) HapusAnggota(ctx fiber.Ctx) error {
	anggotaId := ctx.FormValue("anggota_id")

	ap.anggotaModel.DeleteAnggota(anggotaId)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/anggota?q=%s&page=%s", ctx.FormValue("query"), ctx.FormValue("page")))
}

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

func NewAdminPage(db *sql.DB) *AdminPage {
	anggotaModel := models.NewAnggotaModel(db)
	periodeModel := models.NewPeriodeModel(db)
	userModel := models.NewUserModel(db)
	aktivitasModel := models.NewAktivitasModel(db)

	return &AdminPage{
		anggotaModel:   anggotaModel,
		periodeModel:   periodeModel,
		userModel:      userModel,
		aktivitasModel: aktivitasModel,
	}
}
