package models

import (
	"kinerja-pegawai/internal/utils"
)

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
