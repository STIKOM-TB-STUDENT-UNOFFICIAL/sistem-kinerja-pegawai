package controllers

import (
	"fmt"
	"kinerja-pegawai/internal/utils"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

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
