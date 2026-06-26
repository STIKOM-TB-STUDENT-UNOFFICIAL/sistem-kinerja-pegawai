package models

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/utils"
	"time"
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

func (am *AktivitasModel) FindAktivitas(userid string, limit int, offset int, q string) *[]AktivitasDB {
	q = fmt.Sprintf("%%%s%%", q)
	var result []AktivitasDB = []AktivitasDB{}

	rows, err := am.db.Query(`
		SELECT 
			aktivitas.id,
    		aktivitas.userid, 
    		aktivitas.id_tupoksi,
    		tupoksi.nama,
    		tupoksi.deskripsi,
    		aktivitas.catatan,
    		aktivitas.tanggal,
    		aktivitas.status,
			CASE
			    WHEN aktivitas.selesai IS NULL THEN
			        CONCAT(
			            FLOOR(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()) / 1440), ' Hari ',
			            FLOOR(MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()), 1440) / 60), ' Jam ',
			            MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()), 60), ' Menit',
			            ' (Dalam Penyelesaian)'
			        )
			    ELSE
			        CONCAT(
			            FLOOR(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai) / 1440), ' Hari ',
			            FLOOR(MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai), 1440) / 60), ' Jam ',
			            MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai), 60), ' Menit'
			        )
			END AS waktu_pengerjaan,
			aktivitas.mulai,
			aktivitas.selesai
		FROM aktivitas
		INNER JOIN tupoksi ON id_tupoksi = tupoksi.id
		WHERE aktivitas.userid = ? 
		AND (tupoksi.nama LIKE ? OR tupoksi.deskripsi LIKE ? OR aktivitas.catatan LIKE ?)
		ORDER BY aktivitas.id DESC
		LIMIT ? OFFSET ?`,
		userid, q, q, q, limit, offset,
	)

	if err != nil {
		utils.Log(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var aktivitas AktivitasDB

		rows.Scan(
			&aktivitas.Id,
			&aktivitas.UserId,
			&aktivitas.IdTupoksi,
			&aktivitas.Nama,
			&aktivitas.Deskripsi,
			&aktivitas.Catatan,
			&aktivitas.Tanggal,
			&aktivitas.Status,
			&aktivitas.WaktuPengerjaan,
			&aktivitas.Mulai,
			&aktivitas.Selesai,
		)

		result = append(result, aktivitas)
	}

	return &result
}

func (am *AktivitasModel) FindAktivitasByPeriode(
	userid string,
	periodeId string,
	limit int,
	offset int,
	q string,
) *[]AktivitasDB {
	q = fmt.Sprintf("%%%s%%", q)
	var result []AktivitasDB = []AktivitasDB{}

	rows, err := am.db.Query(`
		SELECT 
			aktivitas.id,
    		aktivitas.userid,
    		aktivitas.id_tupoksi,
    		tupoksi.nama,
    		tupoksi.deskripsi,
    		aktivitas.catatan,
    		aktivitas.tanggal,
    		aktivitas.status,
			CASE
			    WHEN aktivitas.selesai IS NULL THEN
			        CONCAT(
			            FLOOR(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()) / 1440), ' Hari ',
			            FLOOR(MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()), 1440) / 60), ' Jam ',
			            MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()), 60), ' Menit',
			            ' (Dalam Penyelesaian)'
			        )
			    ELSE
			        CONCAT(
			            FLOOR(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai) / 1440), ' Hari ',
			            FLOOR(MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai), 1440) / 60), ' Jam ',
			            MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai), 60), ' Menit'
			        )
			END AS waktu_pengerjaan
		FROM aktivitas
		INNER JOIN tupoksi ON id_tupoksi = tupoksi.id
		INNER JOIN periode ON aktivitas.tanggal >= periode.start_date AND aktivitas.tanggal <= periode.end_date
		WHERE aktivitas.userid = ? 
		AND periode.id = ? AND (tupoksi.nama LIKE ? OR tupoksi.deskripsi LIKE ? OR aktivitas.catatan LIKE ?)
		ORDER BY aktivitas.id DESC
		LIMIT ? OFFSET ?`,
		userid, periodeId, q, q, q, limit, offset,
	)

	if err != nil {
		utils.Log(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var aktivitas AktivitasDB

		rows.Scan(
			&aktivitas.Id,
			&aktivitas.UserId,
			&aktivitas.IdTupoksi,
			&aktivitas.Nama,
			&aktivitas.Deskripsi,
			&aktivitas.Catatan,
			&aktivitas.Tanggal,
			&aktivitas.Status,
			&aktivitas.WaktuPengerjaan,
		)

		result = append(result, aktivitas)
	}

	return &result
}

func (am *AktivitasModel) FindAktivitasByPeriodeAll(
	userid string,
	periodeId string,
) *[]AktivitasDB {
	var result []AktivitasDB = []AktivitasDB{}

	rows, err := am.db.Query(`
		SELECT 
			aktivitas.id,
    		aktivitas.userid,
    		aktivitas.id_tupoksi,
    		tupoksi.nama,
    		tupoksi.deskripsi,
    		aktivitas.catatan,
    		aktivitas.tanggal,
    		aktivitas.status,
			CASE
			    WHEN aktivitas.selesai IS NULL THEN
			        CONCAT(
			            FLOOR(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()) / 1440), ' Hari ',
			            FLOOR(MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()), 1440) / 60), ' Jam ',
			            MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, NOW()), 60), ' Menit',
			            ' (Dalam Penyelesaian)'
			        )
			    ELSE
			        CONCAT(
			            FLOOR(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai) / 1440), ' Hari ',
			            FLOOR(MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai), 1440) / 60), ' Jam ',
			            MOD(TIMESTAMPDIFF(MINUTE, aktivitas.mulai, aktivitas.selesai), 60), ' Menit'
			        )
			END AS waktu_pengerjaan
		FROM aktivitas
		INNER JOIN tupoksi ON id_tupoksi = tupoksi.id
		INNER JOIN periode ON aktivitas.tanggal >= periode.start_date AND aktivitas.tanggal <= periode.end_date
		WHERE aktivitas.userid = ? 
		AND periode.id = ?
		ORDER BY aktivitas.id DESC`,
		userid, periodeId,
	)

	if err != nil {
		utils.Log(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var aktivitas AktivitasDB

		rows.Scan(
			&aktivitas.Id,
			&aktivitas.UserId,
			&aktivitas.IdTupoksi,
			&aktivitas.Nama,
			&aktivitas.Deskripsi,
			&aktivitas.Catatan,
			&aktivitas.Tanggal,
			&aktivitas.Status,
			&aktivitas.WaktuPengerjaan,
		)

		result = append(result, aktivitas)
	}

	return &result
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

func (am *AktivitasModel) InsertAktivitas(
	userid string,
	tupoksi string,
	catatan string,
	tanggal string,
	mulai string,
	selesai any,
	status string,
) (int64, error) {

	result, err := am.db.Exec(`
		INSERT INTO aktivitas
			(id, userid, id_tupoksi, tanggal, mulai, selesai, catatan, status)
		VALUES
			(NULL, ?, ?, ?, ?, ?, ?, ?)
	`,
		userid,
		tupoksi,
		tanggal,
		mulai,
		selesai,
		catatan,
		status,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (am *AktivitasModel) UpdateAktivitas(
	id,
	userid,
	tupoksi,
	catatan string,
	mulai,
	selesai *string,
) {
	query := `
		UPDATE aktivitas
		SET
			id_tupoksi = ?,
			catatan = ?
	`

	args := []any{
		tupoksi,
		catatan,
	}

	if selesai != nil {
		t, err := time.Parse("2006-01-02T15:04", *selesai)

		if err == nil {
			query += `,
				selesai = ?,
				status = 'Diproses'
			`

			args = append(
				args,
				t.Format("2006-01-02 15:04:05"),
			)
		}
	} else {
		query += `,
			status = 'Dalam Pengerjaan'
		`
	}

	if mulai != nil {
		t, err := time.Parse("2006-01-02T15:04", *mulai)

		if err == nil {
			query += `,
				tanggal = ?,
				mulai = ?
			`

			args = append(
				args,
				t.Format("2006-01-02"),
				t.Format("2006-01-02 15:04:05"),
			)
		}
	}

	query += ` WHERE id = ? AND userid = ?`

	args = append(args, id, userid)

	_, err := am.db.Exec(query, args...)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (am *AktivitasModel) UpdateStatusAktivitas(id, userid, status string) {
	_, err := am.db.Query("UPDATE aktivitas SET status = ? WHERE id = ? AND userid = ?", status, id, userid)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (am *AktivitasModel) SetSelesai(id, userid string) {
	_, err := am.db.Query("UPDATE aktivitas SET selesai = STR_TO_DATE(DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i'), '%Y-%m-%d %H:%i') WHERE id = ? AND userid = ?", id, userid)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (am *AktivitasModel) DeleteAktivitas(id string) {
	_, err := am.db.Query("DELETE FROM komentar_aktivitas where id_aktivitas = ?", id)

	if err != nil {
		utils.Log(err.Error())
	}

	_, err = am.db.Query("DELETE FROM aktivitas where id = ?", id)

	if err != nil {
		utils.Log(err.Error())
	}
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
