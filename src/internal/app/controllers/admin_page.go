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
	tupoksiModel   *models.TupoksiModel
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

func NewAdminPage(db *sql.DB) *AdminPage {
	anggotaModel := models.NewAnggotaModel(db)
	periodeModel := models.NewPeriodeModel(db)
	userModel := models.NewUserModel(db)
	aktivitasModel := models.NewAktivitasModel(db)
	tupoksiModel := models.NewTupoksiModel(db)

	return &AdminPage{
		anggotaModel:   anggotaModel,
		periodeModel:   periodeModel,
		userModel:      userModel,
		aktivitasModel: aktivitasModel,
		tupoksiModel:   tupoksiModel,
	}
}
