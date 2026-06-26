package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

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
