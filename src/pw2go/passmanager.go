package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Name of the service and its password
type Password struct {
	service string
	value   string
}

// PasswordManager class with database info
type PasswordManager struct {
	dbname      string
	mastertable string
	passtable   string
	masterhash  string
	db          *sql.DB
}

const (
	// Message for wrong password input
	BADMASTER = "Bad master password"
	// Set master error
	SETMASTERERROR = "Master hash already present"
	// AES key length
	KLEN = 16
)

// Initialize password manager object
func (pm *PasswordManager) Init() error {
	var err error
	pm.db, err = sql.Open("sqlite3", "./"+pm.dbname+".db")
	if err != nil {
		return err
	}
	query := `CREATE TABLE IF NOT EXISTS %s (
		masterhash VARCHAR(512) PRIMARY KEY
	);`
	if _, err := pm.db.Exec(fmt.Sprintf(query, pm.mastertable)); err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS %s (
		service VARCHAR(128) PRIMARY KEY,
		value VARCHAR(1024) NOT NULL,
		nonce VARCHAR(64) NOT NULL UNIQUE
	);`
	if _, err := pm.db.Exec(fmt.Sprintf(query, pm.passtable)); err != nil {
		return err
	}
	rows, err := pm.db.Query("SELECT * FROM " + pm.mastertable + " LIMIT 1")
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&pm.masterhash); err != nil {
			pm.masterhash = ""
			return err
		}
	}
	return nil
}

// Update table containing a row with a single value
func setSingleValueTable(tx *sql.Tx, table string, value string) error {
	if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s;", table)); err != nil {
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s VALUES(?);", table))
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(value); err != nil {
		return err
	}
	return nil
}

// Set master password (if unset)
func (pm *PasswordManager) SetMaster(master string) error {
	if pm.masterhash != "" {
		return errors.New(SETMASTERERROR)
	}
	pm.masterhash = sha256sum(master)
	tx, err := pm.db.Begin()
	if err != nil {
		return err
	}
	if err := setSingleValueTable(tx, pm.mastertable, pm.masterhash); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Add password to database
func (pm *PasswordManager) AddPassword(service string, password string, master string) error {
	if pm.masterhash != sha256sum(master) {
		return errors.New(BADMASTER)
	}
	return errors.New("Not implemented")
}

// Get password from database
func (pm *PasswordManager) GetPassword(service string, master string) (string, error) {
	if pm.masterhash != sha256sum(master) {
		return "", errors.New(BADMASTER)
	}
	return "", errors.New("Not implemented")
}

// Change master password
func (pm *PasswordManager) ChangeMaster(oldmaster string, newmaster string) error {
	// Check old master
	if pm.masterhash != sha256sum(oldmaster) {
		return errors.New(BADMASTER)
	}
	// Set new master
	tx, err := pm.db.Begin()
	if err != nil {
		return err
	}
	hash := sha256sum(newmaster)
	if err := setSingleValueTable(tx, pm.mastertable, hash); err != nil {
		tx.Rollback()
		return err
	}
	rows, err := tx.Query(fmt.Sprintf("SELECT * FROM %s;", pm.passtable))
	if err != nil {
		tx.Rollback()
		return err
	}
	// Update all values in DB
	var service, plainval string
	var cipherobj AESCipher
	stmt, err := tx.Prepare(fmt.Sprintf("UPDATE %s SET value = ?, nonce = ? WHERE service = ?;", pm.passtable))
	if err != nil {
		tx.Rollback()
		return err
	}
	for rows.Next() {
		rows.Scan(service, cipherobj.ciphertext, cipherobj.nonce)
		plainval, err = decryptAESGCM(&cipherobj, oldmaster, service, KLEN)
		if err != nil {
			tx.Rollback()
			return err
		}
		cipherobj, err = encryptAESGCM(plainval, newmaster, service, KLEN)
		if err != nil {
			tx.Rollback()
			return err
		}
		if _, err := stmt.Exec(cipherobj.ciphertext, cipherobj.nonce, service); err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err == nil {
		pm.masterhash = hash
	}
	return err
}

// Close database connection
func (pm *PasswordManager) Close() {
	pm.db.Close()
}
