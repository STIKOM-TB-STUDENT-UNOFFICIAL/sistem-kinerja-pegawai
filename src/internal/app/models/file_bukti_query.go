package models

import (
	"kinerja-pegawai/internal/utils"
)

func (fbm *FileBuktiModel) FindById(fileId string) *FileBukti {
	var file FileBukti

	err := fbm.db.QueryRow("SELECT * FROM file_bukti WHERE id = ?", fileId).Scan(
		&file.Id,
		&file.IdAktivitas,
		&file.NamaFile,
		&file.LokasiFile,
	)

	if err != nil {
		utils.Log(err.Error())
	}

	return &file
}

func (fbm *FileBuktiModel) FindAllByAktivitasId(aktivitasId string) *[]FileBukti {
	rows, err := fbm.db.Query(
		"SELECT * FROM file_bukti WHERE aktivitas_id = ?",
		aktivitasId,
	)

	if err != nil {
		utils.Log(err.Error())
	}

	defer rows.Close()

	var files []FileBukti = []FileBukti{}

	for rows.Next() {
		var file FileBukti

		if err := rows.Scan(
			&file.Id,
			&file.IdAktivitas,
			&file.NamaFile,
			&file.LokasiFile,
		); err != nil {
			continue
		}

		files = append(files, file)
	}

	return &files
}
