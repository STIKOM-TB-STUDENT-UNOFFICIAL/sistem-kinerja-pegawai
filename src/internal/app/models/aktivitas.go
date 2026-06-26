package models

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/utils"
)

type AktivitasDB struct {
	Id              string
	IdTupoksi       string
	UserId          string
	Nama            *string
	Deskripsi       *string
	Catatan         *string
	Tanggal         *string
	Mulai           *string
	Selesai         *string
	Status          string
	WaktuPengerjaan string
}

type AktivitasModel struct {
	db *sql.DB
}

func (am *AktivitasModel) CountAktivitas(userid string, q string) int {
	q = fmt.Sprintf("%%%s%%", q)
	var result int

	err := am.db.QueryRow(`
		SELECT COUNT(aktivitas.id) as cnt
		FROM aktivitas
		INNER JOIN tupoksi ON id_tupoksi = tupoksi.id
		WHERE aktivitas.userid = ? 
		AND (tupoksi.nama LIKE ? OR tupoksi.deskripsi LIKE ? OR aktivitas.catatan LIKE ?)`,
		userid, q, q, q,
	).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (am *AktivitasModel) CountRataRata(userid string, periode string) *string {
	var result *string

	err := am.db.QueryRow(`
		SELECT CONCAT(
		    FLOOR(avg_menit / 1440), ' Hari ',
		    FLOOR(MOD(avg_menit, 1440) / 60), ' Jam ',
		    MOD(avg_menit, 60), ' Menit'
		) AS rata_rata_waktu
		FROM (
		    SELECT ROUND(
		        AVG(
		            TIMESTAMPDIFF(MINUTE, aktivitas.mulai, COALESCE(aktivitas.selesai, NOW()))
		        )
		    ) AS avg_menit
		    FROM aktivitas
		    INNER JOIN periode
		        ON aktivitas.tanggal BETWEEN periode.start_date AND periode.end_date
		    WHERE aktivitas.userid = ?
		      AND periode.id = ?
		) x`,
		userid, periode,
	).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (am *AktivitasModel) FindAktivitaById(id string) *AktivitasDB {
	var result AktivitasDB

	err := am.db.QueryRow(`
		SELECT 
			aktivitas.id,
    		aktivitas.userid, 
    		aktivitas.id_tupoksi,
    		tupoksi.nama,
    		tupoksi.deskripsi,
    		aktivitas.catatan,
    		aktivitas.tanggal,
    		aktivitas.status
		FROM aktivitas
		INNER JOIN tupoksi ON id_tupoksi = tupoksi.id
		WHERE aktivitas.id = ?`, id,
	).Scan(
		&result.Id,
		&result.UserId,
		&result.IdTupoksi,
		&result.Nama,
		&result.Deskripsi,
		&result.Catatan,
		&result.Tanggal,
		&result.Status,
	)

	if err != nil {
		utils.Log(err.Error())
	}

	return &result
}

func (am *AktivitasModel) CountAktivitasUser(userid string) int {
	var result int

	err := am.db.QueryRow(`
		SELECT COUNT(id) as cnt
		FROM aktivitas
		WHERE userid = ?`,
		userid,
	).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func NewAktivitasModel(db *sql.DB) *AktivitasModel {
	return &AktivitasModel{
		db: db,
	}
}
