package controllers

import (
	"fmt"
	"kinerja-pegawai/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func (ap *AdminPage) AdminTupoksiPage(ctx fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		utils.Log(err.Error())
	}

	selectedUser := ctx.Query("user", "")
	query := ctx.Query("q", "")
	limit := 10
	offset := ((page - 1) * limit)

	allParent := ap.anggotaModel.AllParent()

	var result interface{}
	var total int
	if selectedUser != "" {
		result = ap.tupoksiModel.FindTupoksi(selectedUser, limit, offset, query)
		total = ap.tupoksiModel.Count(selectedUser, query)
	}

	return ctx.Render(
		"pages/admin_tupoksi",
		fiber.Map{
			"Title":        "Kelola Tupoksi Pengguna",
			"Level":        "admin",
			"Heading":      "Tupoksi Pengguna",
			"SubHeading":   "Kelola tugas pokok dan fungsi pengguna",
			"Page":         page,
			"Previous":     page - 1,
			"Next":         page + 1,
			"HasNext":      (total - page*limit) > 0,
			"Offset":       offset,
			"Query":        query,
			"Result":       result,
			"AllParent":    allParent,
			"SelectedUser": selectedUser,
		},
		"layouts/main",
	)
}

func (ap *AdminPage) AdminPostTupoksi(ctx fiber.Ctx) error {
	user := ctx.FormValue("user")
	nama := ctx.FormValue("nama")
	deskripsi := ctx.FormValue("deskripsi")

	ap.tupoksiModel.Insert(user, nama, deskripsi)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/tupoksi?user=%s&q=%s&page=%s", user, ctx.FormValue("q"), ctx.FormValue("page")))
}

func (ap *AdminPage) AdminUpdateTupoksi(ctx fiber.Ctx) error {
	user := ctx.FormValue("user")
	id := ctx.FormValue("id")
	nama := ctx.FormValue("nama")
	deskripsi := ctx.FormValue("deskripsi")

	ap.tupoksiModel.Update(id, nama, deskripsi)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/tupoksi?user=%s&q=%s&page=%s", user, ctx.FormValue("q"), ctx.FormValue("page")))
}

func (ap *AdminPage) AdminDeleteTupoksi(ctx fiber.Ctx) error {
	user := ctx.FormValue("user")
	id := ctx.FormValue("id")

	ap.tupoksiModel.Delete(id)

	return ctx.Redirect().To(fmt.Sprintf("/dashboard/admin/tupoksi?user=%s&q=%s&page=%s", user, ctx.FormValue("q"), ctx.FormValue("page")))
}
