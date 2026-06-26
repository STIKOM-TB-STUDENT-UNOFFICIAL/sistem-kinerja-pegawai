package models

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/utils"
)

type AnggotaMinimalDB struct {
	UserId string
	Nama   *string
}

type Parent struct {
	UserId string
	Nama   *string
}

type UserCard struct {
	UserId      string
	NIP         *string
	NamaLengkap *string
	Jabatan     *string
	Departemen  *string
}

type ParentAndChildren struct {
	Id             string
	ParentUserId   string
	ChildrenUserId string
	ChildrenName   string
}

type AnggotaModel struct {
	db *sql.DB
}

func (am *AnggotaModel) FindAllAnggota(userid string, query string, limit, offset int) *[]AnggotaMinimalDB {
	query = fmt.Sprintf("%%%s%%", query)
	var result []AnggotaMinimalDB

	rows, err := am.db.Query(
		`SELECT login_system.userid, login_system.nama_lengkap FROM anggota
		INNER JOIN login_system ON login_system.userid = anggota.userid_children
		WHERE anggota.userid_parent = ? AND login_system.nama_lengkap LIKE ? LIMIT ? OFFSET ?`, userid, query, limit, offset,
	)

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp AnggotaMinimalDB

		rows.Scan(
			&temp.UserId,
			&temp.Nama,
		)

		result = append(result, temp)
	}

	return &result
}

func (am *AnggotaModel) Count(userid string, query string) int {
	query = fmt.Sprintf("%%%s%%", query)
	var result int

	err := am.db.QueryRow(
		`SELECT COUNT(login_system.userid) as cnt FROM anggota
		INNER JOIN login_system ON login_system.userid = anggota.userid_children
		WHERE anggota.userid_parent = ? AND login_system.nama_lengkap LIKE ?`, userid, query,
	).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (am *AnggotaModel) FindAnggotaByUserId(userid, parent string) AnggotaMinimalDB {
	var result AnggotaMinimalDB

	err := am.db.QueryRow(
		`SELECT login_system.userid, login_system.nama_lengkap FROM anggota
		INNER JOIN login_system ON login_system.userid = anggota.userid_children
		WHERE anggota.userid_parent = ? AND anggota.userid_children = ?`, parent, userid,
	).Scan(&result.UserId, &result.Nama)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (am *AnggotaModel) InsertAnggota(parent_id, children_id string) {
	_, err := am.db.Query("INSERT INTO anggota (id, userid_parent, userid_children) VALUES (NULL, ?, ?)", parent_id, children_id)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (am *AnggotaModel) DeleteAnggota(anggotaId string) {
	_, err := am.db.Query("DELETE FROM anggota WHERE id = ?", anggotaId)

	if err != nil {
		utils.Log(err.Error())
	}
}

func NewAnggotaModel(db *sql.DB) *AnggotaModel {
	return &AnggotaModel{
		db: db,
	}
}
