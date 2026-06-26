package models

import (
	"crypto/md5"
	"encoding/hex"
)

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
