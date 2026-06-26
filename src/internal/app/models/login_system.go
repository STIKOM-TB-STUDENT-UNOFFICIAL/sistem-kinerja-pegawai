package models

import "database/sql"

type LoginSystemDB struct {
	Userid      string
	Password    string
	NamaLengkap *string
	Foto        *string
	Level       *string
	Blokir      *string
	LastLogin   *string
}

type LoginSystem struct {
	db *sql.DB
}

func (ls *LoginSystem) GetUserByUsernamePassword(username string, password string) (*LoginSystemDB, error) {
	var loginsystem LoginSystemDB

	err := ls.db.QueryRow("SELECT * from login_system where userid = ? and password = ?", username, password).Scan(
		&loginsystem.Userid,
		&loginsystem.Password,
		&loginsystem.NamaLengkap,
		&loginsystem.Foto,
		&loginsystem.Level,
		&loginsystem.Blokir,
		&loginsystem.LastLogin,
	)

	return &loginsystem, err
}

func NewLoginSystem(db *sql.DB) *LoginSystem {
	return &LoginSystem{
		db: db,
	}
}
