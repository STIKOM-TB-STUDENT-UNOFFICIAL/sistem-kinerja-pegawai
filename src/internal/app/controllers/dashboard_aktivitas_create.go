package controllers

import (
	"fmt"
	"kinerja-pegawai/internal/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

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
