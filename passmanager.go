package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Password Struct with the name of the service and its password
type Password struct {
	service string
	value   string
}

// PasswordManager Struct with database info
type PasswordManager struct {
	dbname      string
	mastertable string
	passtable   string
	masterhash  string
	services    map[string]bool
	db          *sql.DB
}

const (
	// BADMASTER Message for wrong password input
	BADMASTER = "Bad master password"
	// SETMASTERERROR Message for SetMaster error
	SETMASTERERROR = "Master hash already present"
	// ALREADYPRESENT Message for password already present
	ALREADYPRESENT = "Password already present for this service"
	// NOTFOUND Message for service not found
	NOTFOUND = "Service not found"
	// KLEN AES key length
	KLEN = 16
)

// Init Initialize password manager object
func (pm *PasswordManager) Init() error {
	var err error
	pm.db, err = sql.Open("sqlite3", fmt.Sprintf("./%s.db", pm.dbname))
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
		nonce VARCHAR(64) NOT NULL
	);`
	if _, err := pm.db.Exec(fmt.Sprintf(query, pm.passtable)); err != nil {
		return err
	}
	rows, err := pm.db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 1;", pm.mastertable))
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
	pm.services = make(map[string]bool)
	rows, err = pm.db.Query(fmt.Sprintf("SELECT service FROM %s;", pm.passtable))
	if err != nil {
		return err
	}
	var serv string
	for rows.Next() {
		if err = rows.Scan(&serv); err != nil {
			return err
		}
		pm.services[serv] = true
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

// SetMaster sets master password (if unset)
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

// AddPassword adds a password to database
func (pm *PasswordManager) AddPassword(service string, password string, master string) error {
	if pm.masterhash != sha256sum(master) {
		return errors.New(BADMASTER)
	}
	if pm.services[service] {
		return errors.New(ALREADYPRESENT)
	}
	cipherobj, err := encryptAESGCM(password, master, service, KLEN)
	if err != nil {
		return err
	}
	tx, err := pm.db.Begin()
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s VALUES(?, ?, ?);", pm.passtable))
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := stmt.Exec(service, cipherobj.ciphertext, cipherobj.nonce); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err == nil {
		pm.services[service] = true
	}
	return err
}

// GetPassword retrieves a service's password from database
func (pm *PasswordManager) GetPassword(service string, master string) (string, error) {
	if pm.masterhash != sha256sum(master) {
		return "", errors.New(BADMASTER)
	}
	if !pm.services[service] {
		return "", errors.New(NOTFOUND)
	}
	stmt, err := pm.db.Prepare(fmt.Sprintf("SELECT value, nonce FROM %s WHERE service = ?;", pm.passtable))
	if err != nil {
		return "", err
	}
	rows, err := stmt.Query(service)
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		return "", errors.New(NOTFOUND)
	}
	var cipherobj AESCipher
	if err := rows.Scan(&cipherobj.ciphertext, &cipherobj.nonce); err != nil {
		return "", err
	}
	return decryptAESGCM(&cipherobj, master, service, KLEN)
}

// ChangeMaster changes the master password
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
		rows.Scan(&service, &cipherobj.ciphertext, &cipherobj.nonce)
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
