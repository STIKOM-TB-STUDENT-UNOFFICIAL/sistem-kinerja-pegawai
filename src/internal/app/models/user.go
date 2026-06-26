package models

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/utils"
)

type UserDB struct {
	UserId      string
	NamaLengkap string
	Level       string
	Nip         *string
	Jabatan     *string
	Departemen  *string
}

type UserModel struct {
	db *sql.DB
}

func (um *UserModel) FindUser(limit, offset int, q string) *[]UserDB {
	q = fmt.Sprintf("%%%s%%", q)
	var result []UserDB

	rows, err := um.db.Query(
		`SELECT
			login_system.userid,
    		login_system.nama_lengkap,
			login_system.level,
    		detail_pegawai.nip,
    		detail_pegawai.jabatan,
    		detail_pegawai.departemen
		FROM login_system
		LEFT JOIN detail_pegawai ON detail_pegawai.userid = login_system.userid
		WHERE 
			(login_system.nama_lengkap like ?
			OR detail_pegawai.nip like ?
			OR detail_pegawai.jabatan like ?
			OR detail_pegawai.departemen like ?)
			AND login_system.level != 'admin'
		LIMIT ? OFFSET ?
		`, q, q, q, q, limit, offset,
	)

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp UserDB
		rows.Scan(
			&temp.UserId,
			&temp.NamaLengkap,
			&temp.Level,
			&temp.Nip,
			&temp.Jabatan,
			&temp.Departemen,
		)

		result = append(result, temp)
	}

	return &result
}

func (um *UserModel) Count(q string) int {
	q = fmt.Sprintf("%%%s%%", q)
	var result int

	err := um.db.QueryRow(
		`SELECT
			COUNT(*) as cnt
		FROM login_system
		LEFT JOIN detail_pegawai ON detail_pegawai.userid = login_system.userid
		WHERE 
			(login_system.nama_lengkap like ?
			OR detail_pegawai.nip like ?
			OR detail_pegawai.jabatan like ?
			OR detail_pegawai.departemen like ?)
			AND login_system.level != 'admin'
		`, q, q, q, q,
	).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (um *UserModel) CountTotal() int {
	var result int

	err := um.db.QueryRow(
		`SELECT COUNT(*) as cnt FROM login_system WHERE level != 'admin'`,
	).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}
