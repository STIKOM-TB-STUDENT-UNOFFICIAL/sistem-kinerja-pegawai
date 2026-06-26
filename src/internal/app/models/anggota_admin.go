package models

import (
	"fmt"
	"kinerja-pegawai/internal/utils"
)

func (am *AnggotaModel) FindParent(query string, limit, offset int) *[]Parent {
	query = fmt.Sprintf("%%%s%%", query)
	var result []Parent

	rows, err := am.db.Query("SELECT userid, nama_lengkap FROM login_system WHERE (login_system.level = 'dosen' OR login_system.level = 'pegawai') AND nama_lengkap LIKE ? LIMIT ? OFFSET ?", query, limit, offset)

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp Parent
		rows.Scan(
			&temp.UserId,
			&temp.Nama,
		)
		result = append(result, temp)
	}

	return &result
}

func (am *AnggotaModel) AllParent() *[]Parent {
	var result []Parent

	rows, err := am.db.Query("SELECT userid, nama_lengkap FROM login_system WHERE (login_system.level = 'dosen' OR login_system.level = 'pegawai')")

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp Parent
		rows.Scan(
			&temp.UserId,
			&temp.Nama,
		)
		result = append(result, temp)
	}

	return &result
}

func (am *AnggotaModel) CountParent(query string) int {
	query = fmt.Sprintf("%%%s%%", query)
	var result int

	err := am.db.QueryRow("SELECT COUNT(userid) as cnt FROM login_system WHERE (login_system.level = 'dosen' OR login_system.level = 'pegawai') AND nama_lengkap LIKE ?", query).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (am *AnggotaModel) FindParentChildren(query string, limit, offset int) *[]ParentAndChildren {
	query = fmt.Sprintf("%%%s%%", query)
	var result []ParentAndChildren

	rows, err := am.db.Query(`
		SELECT
			anggota.id as id,
			parent.userid as parent_userid,
			children.userid as children_userid, 
			children.nama_lengkap 
		FROM anggota
		INNER JOIN login_system parent ON parent.userid = anggota.userid_parent
		INNER JOIN login_system children ON children.userid = anggota.userid_children
		WHERE (parent.level = 'dosen' OR parent.level = 'pegawai') AND parent.nama_lengkap LIKE ? LIMIT ? OFFSET ?;
		`, query, limit, offset)

	if err != nil {
		utils.Log(err.Error())
		return &result
	}
	defer rows.Close()

	for rows.Next() {
		var temp ParentAndChildren
		rows.Scan(
			&temp.Id,
			&temp.ParentUserId,
			&temp.ChildrenUserId,
			&temp.ChildrenName,
		)
		result = append(result, temp)
	}

	return &result
}
