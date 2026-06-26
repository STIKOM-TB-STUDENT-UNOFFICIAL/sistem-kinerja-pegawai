package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"kinerja-pegawai/internal/app/models"
	"kinerja-pegawai/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type DashboardPage struct {
	db                     *sql.DB
	loginSystemModel       *models.LoginSystem
	aktivitasModel         *models.AktivitasModel
	tupoksiModel           *models.TupoksiModel
	anggotaModel           *models.AnggotaModel
	komentarAktivitasModel *models.KomentarAktivitasModel
	periodeModel           *models.PeriodeModel
	fileBuktiModel         *models.FileBuktiModel
}

func (dp DashboardPage) LandingPage(ctx fiber.Ctx) error {
	return utils.ViewFile(ctx, "pages/landing")
}

func (db DashboardPage) LoginPage(ctx fiber.Ctx) error {
	return ctx.Render("pages/login", fiber.Map{})
}

func (db DashboardPage) Logout(ctx fiber.Ctx) error {
	sess := session.FromContext(ctx)

	if sess == nil {
		return ctx.Redirect().Status(302).To("/login")
	}

	err := sess.Destroy()

	if err != nil {
		return ctx.Redirect().Status(302).To("/login")
	}

	return ctx.Redirect().Status(302).To("/login")
}

func (db *DashboardPage) LoginAction(ctx fiber.Ctx) error {
	username := ctx.FormValue("username", "")
	pwSum := md5.Sum([]byte(ctx.FormValue("password", "")))
	password := hex.EncodeToString(pwSum[:])
	sess := session.FromContext(ctx)

	result, err := db.loginSystemModel.GetUserByUsernamePassword(username, password)

	if err != nil {
		return ctx.Redirect().To("/login?error=true&msg=Username%%20atau%%20Password%%20anda%%20salah")
	}

	if err := sess.Regenerate(); err != nil {
		return ctx.Redirect().To("/login?error=true&msg=Error%%20Internal%%20Server")
	}

	if *result.Blokir != "N" {
		return ctx.Redirect().To("/login?error=true&msg=Akun%%20anda%%20telah%%20diblokir")
	}

	sess.Set("level", *result.Level)
	sess.Set("is_auth", true)
	sess.Set("userid", result.Userid)
	sess.Set("nama_lengkap", result.NamaLengkap)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/%s", *result.Level))
}

func (db *DashboardPage) UserHomePage(ctx fiber.Ctx) error {
	userid := ctx.Locals("userid").(string)
	totalAktivitas := db.aktivitasModel.CountAktivitasUser(userid)
	totalAnggota := db.anggotaModel.Count(userid, "")

	return ctx.Render(
		"pages/user_home",
		fiber.Map{
			"Title":          "Dashboard",
			"Level":          ctx.Locals("level"),
			"Heading":        "Dashboard",
			"SubHeading":     "Ringkasan aktivitas anda",
			"TotalAktivitas": totalAktivitas,
			"TotalAnggota":   totalAnggota,
		},
		"layouts/main",
	)
}

func NewDashboardPage(db *sql.DB) *DashboardPage {
	loginSystemModel := models.NewLoginSystem(db)
	aktivitasModel := models.NewAktivitasModel(db)
	tupoksiModel := models.NewTupoksiModel(db)
	anggotaModel := models.NewAnggotaModel(db)
	komentarAktivitasModel := models.NewKomentarAktivitasModel(db)
	periodeModel := models.NewPeriodeModel(db)
	fileBuktiModel := models.NewFileBuktiModel(db)

	return &DashboardPage{
		db:                     db,
		loginSystemModel:       loginSystemModel,
		aktivitasModel:         aktivitasModel,
		tupoksiModel:           tupoksiModel,
		anggotaModel:           anggotaModel,
		komentarAktivitasModel: komentarAktivitasModel,
		periodeModel:           periodeModel,
		fileBuktiModel:         fileBuktiModel,
	}
}
