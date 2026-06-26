package models

import (
	"kinerja-pegawai/internal/utils"
)

func (am *AnggotaModel) AnggotaCard(anggotaId string) *UserCard {
	var result UserCard
	err := am.db.QueryRow(`
			SELECT 
				anggota.userid_parent,
				detail_pegawai.nip,
    			login_system.nama_lengkap,
    			detail_pegawai.jabatan,
    			detail_pegawai.departemen
			FROM anggota
			INNER JOIN login_system ON anggota.userid_children = login_system.userid
			LEFT JOIN detail_pegawai ON anggota.userid_children = detail_pegawai.userid
			WHERE anggota.userid_children = ? LIMIT 1
		`, anggotaId).Scan(&result.UserId, &result.NIP, &result.NamaLengkap, &result.Jabatan, &result.Departemen)

	if err != nil {
		utils.Log(err.Error())
	}

	return &result
}

func (am *AnggotaModel) FindAtasan(anggotaId string) *[]UserCard {
	var atasan []UserCard = []UserCard{}
	anggota := anggotaId

	for range 2 {
		var temp UserCard
		err := am.db.QueryRow(`
			SELECT 
				anggota.userid_parent,
				detail_pegawai.nip,
    			login_system.nama_lengkap,
    			detail_pegawai.jabatan,
    			detail_pegawai.departemen
			FROM anggota
			INNER JOIN login_system ON anggota.userid_parent = login_system.userid
			LEFT JOIN detail_pegawai ON anggota.userid_parent = detail_pegawai.userid
			WHERE anggota.userid_children = ? LIMIT 1
		`, anggota).Scan(&temp.UserId, &temp.NIP, &temp.NamaLengkap, &temp.Jabatan, &temp.Departemen)

		if err != nil {
			utils.Log(err.Error())
			break
		}

		atasan = append(atasan, temp)
		anggota = temp.UserId
	}

	return &atasan
}
