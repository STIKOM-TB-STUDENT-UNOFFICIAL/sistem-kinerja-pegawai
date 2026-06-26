package models

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/utils"
)

type TupoksiDB struct {
	Id        string
	UserId    string
	Nama      *string
	Deskripsi *string
}

type TupoksiModel struct {
	db *sql.DB
}

func (tm *TupoksiModel) FindAllTupoksi(userid string) *[]TupoksiDB {
	var result []TupoksiDB

	rows, err := tm.db.Query("SELECT * FROM tupoksi WHERE userid = ?", userid)

	if err != nil {
		utils.Log(err.Error())
	}

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp TupoksiDB
		rows.Scan(
			&temp.Id,
			&temp.UserId,
			&temp.Nama,
			&temp.Deskripsi,
		)
		result = append(result, temp)
	}

	return &result
}

func (tm *TupoksiModel) FindTupoksi(userid string, limit, offset int, query string) *[]TupoksiDB {
	var result []TupoksiDB
	query = fmt.Sprintf("%%%s%%", query)

	rows, err := tm.db.Query("SELECT * FROM tupoksi WHERE userid = ? AND (nama LIKE ? OR deskripsi LIKE ?) LIMIT ? OFFSET ?",
		userid, query, query, limit, offset,
	)

	if err != nil {
		utils.Log(err.Error())
	}

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp TupoksiDB
		rows.Scan(
			&temp.Id,
			&temp.UserId,
			&temp.Nama,
			&temp.Deskripsi,
		)
		result = append(result, temp)
	}

	return &result
}

func (tm *TupoksiModel) Count(userid, query string) int {
	var result int
	query = fmt.Sprintf("%%%s%%", query)

	err := tm.db.QueryRow("SELECT COUNT(*) as cnt FROM tupoksi WHERE userid = ? AND (nama LIKE ? OR deskripsi LIKE ?)", userid, query, query).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

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

func NewTupoksiModel(db *sql.DB) *TupoksiModel {
	return &TupoksiModel{
		db: db,
	}
}
