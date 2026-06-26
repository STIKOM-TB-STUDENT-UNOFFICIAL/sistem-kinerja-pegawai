package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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

func (um *UserModel) Insert(
	userId string,
	level string,
	nip string,
	namaLengkap string,
	jabatan string,
	departemen string,
) error {
	hashBytes := md5.Sum([]byte(userId))
	password := hex.EncodeToString(hashBytes[:])

	_, err := um.db.Query(`
		INSERT INTO login_system (userid, password, nama_lengkap, foto, level, blokir)
		VALUES (?, ?, ?, 'kosong.png', ?, 'N')
	`, userId, password, namaLengkap, level)

	if err != nil {
		return err
	}

	_, err = um.db.Query(`
		INSERT INTO detail_pegawai (userid, nip, jabatan, departemen)
		VALUES (?, ?, ?, ?)
	`, userId, nip, jabatan, departemen)

	return err
}

func (um *UserModel) Upsert(
	userId string,
	level string,
	nip string,
	namaLengkap string,
	jabatan string,
	departemen string,
) error {
	_, err := um.db.Query(`
		UPDATE login_system SET nama_lengkap = ?, level = ? WHERE userid = ?
	`, namaLengkap, level, userId)

	if err != nil {
		return err
	}

	_, err = um.db.Query(`
		INSERT INTO detail_pegawai (userid, nip, jabatan, departemen)
		VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE nip = VALUES(nip), jabatan = VALUES(jabatan), departemen = VALUES(departemen)
	`, userId, nip, jabatan, departemen)

	return err
}

func (um *UserModel) Delete(
	userid string,
) error {
	_, err := um.db.Query("DELETE from anggota where userid_parent = ? or userid_children = ?", userid, userid)

	if err != nil {
		return err
	}

	_, err = um.db.Query("DELETE from aktivitas where userid = ?", userid)

	if err != nil {
		return err
	}

	_, err = um.db.Query("DELETE from komentar_aktivitas where userid = ?", userid)

	if err != nil {
		return err
	}

	_, err = um.db.Query("DELETE from tupoksi where userid = ?", userid)

	if err != nil {
		return err
	}

	_, err = um.db.Query("DELETE from detail_pegawai where userid = ?", userid)

	if err != nil {
		return err
	}

	_, err = um.db.Query("DELETE from login_system where userid = ?", userid)

	return err
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
