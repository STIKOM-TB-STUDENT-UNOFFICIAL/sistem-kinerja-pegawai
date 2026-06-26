package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"kinerja-pegawai/internal/app/models"
	"kinerja-pegawai/internal/utils"
	"os"
	"strconv"
	"strings"
	"time"

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

func (db *DashboardPage) PostAktivitas(ctx fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return ctx.Redirect().To(fmt.Sprintf("/dashboard/aktivitas?error=%s", err.Error()))
	}

	var (
		tanggal string
		mulai   string
		selesai any
		status  string
	)
	tupoksi := form.Value["tupoksi"]
	catatan := form.Value["catatan"]

	if len(tupoksi) == 0 || len(catatan) == 0 {
		return ctx.Redirect().To("/dashboard/aktivitas")
	}

	tanggalMulai := ""
	if v, ok := form.Value["startdate"]; ok && len(v) > 0 {
		tanggalMulai = v[0]
	}

	tanggalSelesai := ""
	if v, ok := form.Value["enddate"]; ok && len(v) > 0 {
		tanggalSelesai = v[0]
	}

	if tanggalMulai != "" {
		tanggal = tanggalMulai[:10]
		mulai = tanggalMulai
	} else {
		tanggal = time.Now().Format("2006-01-02")
		mulai = time.Now().Truncate(time.Minute).String()
	}

	if tanggalSelesai != "" {
		selesai = tanggalSelesai
		status = "Diproses"
	} else {
		selesai = nil
		status = "Dalam Pengerjaan"
	}

	aktivitasID, err := db.aktivitasModel.InsertAktivitas(
		ctx.Locals("userid").(string),
		tupoksi[0],
		catatan[0],
		tanggal,
		mulai,
		selesai,
		status,
	)

	if err != nil {
		return ctx.Redirect().To("/dashboard/aktivitas")
	}

	files := form.File["file"]

	if len(files) > 0 {
		for _, file := range files {

			filename := fmt.Sprintf(
				"%d_%s",
				time.Now().UnixNano(),
				file.Filename,
			)

			lokasi := "./public/uploads/" + filename

			if err := ctx.SaveFile(file, lokasi); err != nil {
				utils.Log(err.Error())
				continue
			}

			err := db.fileBuktiModel.Insert(
				aktivitasID,
				file.Filename,
				lokasi,
			)

			if err != nil {
				utils.Log(err.Error())
			}
		}
	}

	return ctx.Redirect().To("/dashboard/aktivitas")
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

func (db *DashboardPage) CheckSelesai(ctx fiber.Ctx) error {
	userid := ctx.Locals("userid").(string)
	aktivitasId := ctx.Params("aktivitasId")

	db.aktivitasModel.UpdateStatusAktivitas(aktivitasId, userid, "Diproses")
	db.aktivitasModel.SetSelesai(aktivitasId, userid)

	return ctx.Redirect().To("/dashboard/aktivitas")
}

func (db *DashboardPage) DeleteTupoksi(ctx fiber.Ctx) error {
	id := ctx.FormValue("id")

	db.tupoksiModel.Delete(id)

	return ctx.Redirect().To("/dashboard/tupoksi")
}

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

func (db *DashboardPage) AktivitasViewPage(ctx fiber.Ctx) error {
	aktivitasId := ctx.Params("aktivitasId")

	aktivitas := db.aktivitasModel.FindAktivitaById(aktivitasId)
	komentar := db.komentarAktivitasModel.FindAllKomentar(aktivitasId)
	fileBukti := db.fileBuktiModel.FindAllByAktivitasId(aktivitasId)

	for i := range *fileBukti {
		(*fileBukti)[i].LokasiFile = strings.Replace((*fileBukti)[i].LokasiFile, "./public", "", 1)
	}

	userid := ctx.Params("userid", "")
	anggotaView := false

	if userid != "" {
		anggotaView = true
	}

	return ctx.Render(
		"pages/user_aktivitas_view",
		fiber.Map{
			"Title":           "Komentar Aktivitas",
			"Level":           ctx.Locals("level"),
			"Heading":         "Komentar Aktivitas",
			"SubHeading":      "Berikan masukan terkait aktivitas",
			"Aktivitas":       aktivitas,
			"Komentar":        komentar,
			"Anggota":         ctx.Params("userid", ""),
			"UserId":          ctx.Locals("userid").(string),
			"AnggotaView":     anggotaView,
			"FileBukti":       fileBukti,
			"FileBuktiLength": len(*fileBukti),
		},
		"layouts/main",
	)
}

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

func (db *DashboardPage) PostFile(ctx fiber.Ctx) error {
	aktivitasId := ctx.Params("aktivitasId")
	form, err := ctx.MultipartForm()

	if err != nil {
		utils.Log(err.Error())
		return ctx.Redirect().To(fmt.Sprintf("/dashboard/aktivitas/%s", aktivitasId))
	}

	files := form.File["file"]

	if len(files) > 0 {
		for _, file := range files {

			filename := fmt.Sprintf(
				"%d_%s",
				time.Now().UnixNano(),
				file.Filename,
			)

			lokasi := "./public/uploads/" + filename

			if err := ctx.SaveFile(file, lokasi); err != nil {
				utils.Log(err.Error())
				continue
			}

			id, err := strconv.Atoi(aktivitasId)

			err = db.fileBuktiModel.Insert(
				int64(id),
				file.Filename,
				lokasi,
			)

			if err != nil {
				utils.Log(err.Error())
			}
		}
	}

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/aktivitas/%s", aktivitasId))
}

func (db *DashboardPage) DeleteFile(ctx fiber.Ctx) error {
	aktivitasId := ctx.Params("aktivitasId")
	fileId := ctx.FormValue("id_file")

	file := db.fileBuktiModel.FindById(fileId)

	if file != nil {
		os.Remove(file.LokasiFile)
	}

	db.fileBuktiModel.DeleteById(fileId)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/aktivitas/%s", aktivitasId))
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
