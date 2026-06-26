package models

import (
	"kinerja-pegawai/internal/utils"
)

func (tm *TupoksiModel) Insert(userid, nama, deskripsi string) {
	_, err := tm.db.Query("INSERT INTO tupoksi (id, userid, nama, deskripsi) VALUES (NULL, ?, ?, ?)", userid, nama, deskripsi)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (tm *TupoksiModel) Update(id, nama, deskripsi string) {
	_, err := tm.db.Query("UPDATE tupoksi SET nama = ?, deskripsi = ? WHERE id = ?", nama, deskripsi, id)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (tm *TupoksiModel) Delete(id string) {
	_, err := tm.db.Query("DELETE FROM komentar_aktivitas WHERE id_aktivitas IN (SELECT id FROM aktivitas WHERE id_tupoksi = ?)", id)

	if err != nil {
		utils.Log(err.Error())
	}

	_, err = tm.db.Query("DELETE FROM aktivitas WHERE id_tupoksi = ?", id)

	if err != nil {
		utils.Log(err.Error())
	}

	_, err = tm.db.Query("DELETE FROM tupoksi WHERE id = ?", id)

	if err != nil {
		utils.Log(err.Error())
	}
}
