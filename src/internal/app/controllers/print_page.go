package controllers

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/app/models"
	"time"

	"github.com/gofiber/fiber/v3"
)

type PrintPage struct {
	anggotaModel   *models.AnggotaModel
	aktivitasModel *models.AktivitasModel
	periodeModel   *models.PeriodeModel
}

var namaBulan []string = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func (pp *PrintPage) Print(ctx fiber.Ctx) error {
	userId := ctx.Params("userid")
	periode := ctx.Params("periode")

	anggotaCard := pp.anggotaModel.AnggotaCard(userId)
	atasanCards := pp.anggotaModel.FindAtasan(userId)

	aktivitasList := pp.aktivitasModel.FindAktivitasByPeriodeAll(userId, periode)

	waktu := time.Now()

	printDate := fmt.Sprintf("%02d %s %d",
		waktu.Day(),
		namaBulan[waktu.Month()-1],
		waktu.Year(),
	)

	totalAktivitas := len(*aktivitasList)

	averageCompletionTime := pp.aktivitasModel.CountRataRata(userId, periode)

	periodeName := pp.periodeModel.FindById(periode)

	return ctx.Render("pages/print", fiber.Map{
		"Pengguna":              anggotaCard,
		"Atasan":                atasanCards,
		"Aktivitas":             aktivitasList,
		"PrintDate":             printDate,
		"Periode":               periodeName.NamaPeriode,
		"TotalAktivitas":        totalAktivitas,
		"AverageCompletionTime": averageCompletionTime,
	})
}

func NewPrintPage(db *sql.DB) *PrintPage {
	return &PrintPage{
		anggotaModel:   models.NewAnggotaModel(db),
		aktivitasModel: models.NewAktivitasModel(db),
		periodeModel:   models.NewPeriodeModel(db),
	}
}
