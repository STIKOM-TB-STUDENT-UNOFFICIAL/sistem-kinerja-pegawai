package models

import (
	"database/sql"
	"kinerja-pegawai/internal/utils"
)

type FileBukti struct {
	Id          string
	IdAktivitas string
	NamaFile    string
	LokasiFile  string
}

type FileBuktiModel struct {
	db *sql.DB
}

func (fbm *FileBuktiModel) Insert(
	aktivitasID int64,
	namaFile string,
	lokasiFile string,
) error {

	_, err := fbm.db.Exec(`
		INSERT INTO file_bukti
			(id, aktivitas_id, nama_file, lokasi_file)
		VALUES
			(NULL, ?, ?, ?)
	`, aktivitasID, namaFile, lokasiFile)

	return err
}

func (fbm *FileBuktiModel) DeleteByAktivitas(aktivitasId string) {
	_, err := fbm.db.Exec(`
		DELETE FROM file_bukti where aktivitas_id = ?
	`, aktivitasId)

	if err != nil {
		utils.Log(err.Error())
	}
}

func (fbm *FileBuktiModel) DeleteById(fileId string) {
	_, err := fbm.db.Exec(`
		DELETE FROM file_bukti where id = ?
	`, fileId)

	if err != nil {
		utils.Log(err.Error())
	}
}

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

func NewFileBuktiModel(db *sql.DB) *FileBuktiModel {
	return &FileBuktiModel{
		db: db,
	}
}
