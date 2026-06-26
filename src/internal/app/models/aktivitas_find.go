package models

import (
	"fmt"
	"kinerja-pegawai/internal/utils"
)

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
