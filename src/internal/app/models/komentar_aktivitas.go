package models

import (
	"database/sql"
	"kinerja-pegawai/internal/utils"
)

type Komentar struct {
	IdAktivitas string
	UserId      string
	IdKomentar  string
	NamaLengkap *string
	Komentar    *string
}

type KomentarAktivitasModel struct {
	db *sql.DB
}

func (kam *KomentarAktivitasModel) FindAllKomentar(idAktivitas string) *[]Komentar {
	var result []Komentar

	rows, err := kam.db.Query(
		`SELECT 
			komentar_aktivitas.id_aktivitas, 
    		login_system.userid,
			komentar_aktivitas.id,
    		login_system.nama_lengkap, 
    		komentar_aktivitas.komentar 
		FROM komentar_aktivitas
		INNER JOIN login_system ON login_system.userid = komentar_aktivitas.userid
		WHERE komentar_aktivitas.id_aktivitas = ?`, idAktivitas,
	)

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp Komentar

		rows.Scan(
			&temp.IdAktivitas,
			&temp.UserId,
			&temp.IdKomentar,
			&temp.NamaLengkap,
			&temp.Komentar,
		)

		result = append(result, temp)
	}

	return &result
}

func (kam *KomentarAktivitasModel) InsertKomentar(userid, aktivitasId, komentar string) {
	_, err := kam.db.Query("INSERT INTO komentar_aktivitas (id, id_aktivitas, userid, komentar) VALUES (NULL, ?, ?, ?)", aktivitasId, userid, komentar)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (kam *KomentarAktivitasModel) DeleteKomentar(komentarId string) {
	_, err := kam.db.Query("DELETE FROM komentar_aktivitas where id = ?", komentarId)

	if err != nil {
		utils.Log(err.Error())
	}
}

func NewKomentarAktivitasModel(db *sql.DB) *KomentarAktivitasModel {
	return &KomentarAktivitasModel{
		db: db,
	}
}
