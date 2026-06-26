package models

import (
	"database/sql"
	"fmt"
	"kinerja-pegawai/internal/utils"
)

type PeriodeDB struct {
	IdPeriode   string
	NamaPeriode string
	StartDate   string
	EndDate     string
}

type PeriodeModel struct {
	db *sql.DB
}

func (pm *PeriodeModel) FindAll() *[]PeriodeDB {
	var result []PeriodeDB

	rows, err := pm.db.Query(`
		SELECT 
			id, nama_periode, start_date, end_date
		FROM periode
		ORDER BY id DESC
	`)

	if err != nil {
		utils.Log(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var temp PeriodeDB
		rows.Scan(
			&temp.IdPeriode,
			&temp.NamaPeriode,
			&temp.StartDate,
			&temp.EndDate,
		)

		result = append(result, temp)
	}

	if err = rows.Err(); err != nil {
		utils.Log(err.Error())
	}

	return &result
}

func (pm *PeriodeModel) FindPeriode(limit, offset int, q string) []PeriodeDB {
	q = fmt.Sprintf("%%%s%%", q)
	var result []PeriodeDB

	rows, err := pm.db.Query(`
		SELECT 
			id, nama_periode, start_date, end_date
		FROM periode
		WHERE nama_periode like ?
		LIMIT ? OFFSET ?
	`, q, limit, offset)

	if err != nil {
		utils.Log(err.Error())
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var temp PeriodeDB
		rows.Scan(
			&temp.IdPeriode,
			&temp.NamaPeriode,
			&temp.StartDate,
			&temp.EndDate,
		)

		result = append(result, temp)
	}

	if err = rows.Err(); err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (pm *PeriodeModel) Count(q string) int {
	var result int
	q = fmt.Sprintf("%%%s%%", q)
	err := pm.db.QueryRow("SELECT COUNT(*) as cnt FROM periode WHERE nama_periode LIKE ?", q).Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func (pm *PeriodeModel) Insert(nama, startDate, endDate string) {
	_, err := pm.db.Query("INSERT INTO periode (id, nama_periode, start_date, end_date) VALUES (NULL, ?, ?, ?)", nama, startDate, endDate)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (pm *PeriodeModel) Update(id, nama, startDate, endDate string) {
	_, err := pm.db.Query("UPDATE periode SET nama_periode = ?, start_date = ?, end_date = ? WHERE id = ?", nama, startDate, endDate, id)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (pm *PeriodeModel) Delete(id string) {
	_, err := pm.db.Query("DELETE FROM periode WHERE id = ?", id)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (pm *PeriodeModel) FindById(id string) *PeriodeDB {
	var temp PeriodeDB
	err := pm.db.QueryRow("SELECT id, nama_periode, start_date, end_date FROM periode WHERE id = ?", id).Scan(
		&temp.IdPeriode,
		&temp.NamaPeriode,
		&temp.StartDate,
		&temp.EndDate,
	)

	if err != nil {
		utils.Log(err.Error())
		return nil
	}

	return &temp
}

func (pm *PeriodeModel) CountTotal() int {
	var result int
	err := pm.db.QueryRow("SELECT COUNT(*) as cnt FROM periode").Scan(&result)

	if err != nil {
		utils.Log(err.Error())
	}

	return result
}

func NewPeriodeModel(db *sql.DB) *PeriodeModel {
	return &PeriodeModel{
		db: db,
	}
}
