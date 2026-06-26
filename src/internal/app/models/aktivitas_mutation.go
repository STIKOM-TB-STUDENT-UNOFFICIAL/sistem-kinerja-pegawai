package models

import (
	"kinerja-pegawai/internal/utils"
	"time"
)

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
