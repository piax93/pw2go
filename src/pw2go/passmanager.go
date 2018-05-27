package main

import (
	//	"crypto/aes"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	_ "github.com/mattn/go-sqlite3"
)

type Password struct {
	service string
	value   string
}

type PasswordManager struct {
	dbname      string
	mastertable string
	passtable   string
	masterhash  string
	db          *sql.DB
}

func (pm *PasswordManager) Init() error {
	var err error
	dbfile := "./" + pm.dbname + ".db"
	pm.db, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}
	query := "CREATE TABLE IF NOT EXISTS " + pm.mastertable + " (masterhash VARCHAR(512) PRIMARY KEY);"
	if _, err := pm.db.Exec(query); err != nil {
		return err
	}
	query = "CREATE TABLE IF NOT EXISTS " + pm.passtable + " (service VARCHAR(128) PRIMARY KEY, value VARCHAR(1024) NOT NULL);"
	if _, err := pm.db.Exec(query); err != nil {
		return err
	}
	rows, err := pm.db.Query("SELECT * FROM " + pm.mastertable + " LIMIT 1")
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&pm.masterhash); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PasswordManager) SetMaster(master string) error {
	hash := sha256.Sum256([]byte(master))
	pm.masterhash = base64.StdEncoding.EncodeToString(hash[:])
	stmt, err := pm.db.Prepare("INSERT INTO " + pm.mastertable + " VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(pm.masterhash); err != nil {
		return err
	}
	return nil
}

func (pm *PasswordManager) Close() {
	pm.db.Close()
}
