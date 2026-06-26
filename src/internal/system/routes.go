package system

import (
	"database/sql"
	"kinerja-pegawai/internal/app/controllers"
	"kinerja-pegawai/internal/app/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(route fiber.Router, db *sql.DB) {
	dp := controllers.NewDashboardPage(db)
	ap := controllers.NewAdminPage(db)
	pp := controllers.NewPrintPage(db)

	route.Get("/", dp.LandingPage)
	route.Get("/login", dp.LoginPage)
	route.Get("/logout", middleware.Authentication, dp.Logout)
	route.Post("/login-action", dp.LoginAction)

	admin := route.Group(
		"/dashboard/admin",
		middleware.Authentication,
		middleware.Authorization([]string{"admin"}, "/dashboard"),
	)
	admin.Get("/", ap.DashboardAdminRingkasan)
	admin.Get("/anggota", ap.AnggotaAdminPage)
	admin.Get("/pengguna", ap.UserPage)
	admin.Get("/pengguna/:userid", ap.ViewAktivitas)

	admin.Post("/pengguna/post", ap.PostPengguna)
	admin.Post("/pengguna/update", ap.UpdatePengguna)
	admin.Post("/pengguna/delete", ap.DeletePengguna)

	admin.Get("/periode", ap.PeriodeList)
	admin.Get("/periode/create", ap.PeriodeCreatePage)
	admin.Post("/periode/create", ap.PeriodeCreate)
	admin.Get("/periode/edit/:id", ap.PeriodeEditPage)
	admin.Post("/periode/update", ap.PeriodeUpdate)
	admin.Post("/periode/delete", ap.PeriodeDelete)

	admin.Post("/anggota/post", ap.InsertAnggota)
	admin.Post("/anggota/delete", ap.HapusAnggota)

	admin.Get("/tupoksi", ap.AdminTupoksiPage)
	admin.Post("/tupoksi/post", ap.AdminPostTupoksi)
	admin.Post("/tupoksi/update", ap.AdminUpdateTupoksi)
	admin.Post("/tupoksi/delete", ap.AdminDeleteTupoksi)

	user := route.Group(
		"/dashboard",
		middleware.Authentication,
		middleware.Authorization([]string{"dosen", "pegawai"}, "/dashboard"),
	)

	user.Get("/tupoksi", dp.UserTupoksiPage)
	user.Get("/anggota", dp.AnggotaPage)
	user.Get("/anggota/:userid", dp.UserAnggotaAktivitasPage)
	user.Get("/anggota/:userid/:aktivitasId", dp.AktivitasViewPage)

	user.Get("/aktivitas", dp.UserAktivitasPage)
	user.Get("/aktivitas/:aktivitasId", dp.AktivitasViewPage)
	user.Get("/aktivitas/:aktivitasId/selesai", dp.CheckSelesai)

	user.Post("/aktivitas/:aktivitasId/post", dp.PostFile)
	user.Post("/aktivitas/:aktivitasId/hapus", dp.DeleteFile)

	user.Post("/tupoksi/post", dp.PostTupoksi)
	user.Post("/tupoksi/update", dp.UpdateTupoksi)
	user.Post("/tupoksi/delete", dp.DeleteTupoksi)

	user.Post("/anggota/:userid/update", dp.UpdateAktivitasAnggota)
	user.Post("/anggota/:userid/:aktivitasId/post_komentar", dp.PostKomentar)
	user.Post("/anggota/:userid/:aktivitasId/hapus_komentar", dp.HapusKomentar)

	user.Post("/aktivitas/post", dp.PostAktivitas)
	user.Post("/aktivitas/update", dp.UpdateAktivitas)
	user.Post("/aktivitas/delete", dp.DeleteAktivitas)
	user.Post("/aktivitas/:aktivitasId/post_komentar", dp.PostKomentar)
	user.Post("/aktivitas/:aktivitasId/hapus_komentar", dp.HapusKomentar)

	user.Get("/:level", dp.UserHomePage)

	print := route.Group("/print", middleware.Authentication)

	print.Get("/:userid/:periode", pp.Print)
}
