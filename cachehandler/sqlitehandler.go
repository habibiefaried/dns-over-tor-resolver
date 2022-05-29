package cachehandler

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	FileName string
	db       *sql.DB
}

func (s *SqliteHandler) Init() error {
	var err error
	s.db, err = sql.Open("sqlite3", s.FileName)
	if err != nil {
		return err
	}

	const create string = `
	CREATE TABLE IF NOT EXISTS dnscache (
		key TEXT NOT NULL PRIMARY KEY,
		value TEXT NOT NULL
	);
	
	CREATE UNIQUE INDEX IF NOT EXISTS dnskey ON dnscache(key);
	`

	if _, err := s.db.Exec(create); err != nil {
		return err
	}

	return nil
}

func (s *SqliteHandler) Get(key string) (*string, error) {
	var err error
	var keyR, valueR string
	row := s.db.QueryRow(`SELECT key, value FROM dnscache WHERE key=$1;`, key)
	switch err = row.Scan(&keyR, &valueR); err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		fmt.Printf("domain %v is on cache with IP %v\n", keyR, valueR)
		return &valueR, nil
	default:
		return nil, err
	}
}

func (s *SqliteHandler) Put(key string, value string) error {
	isExist := false
	var err error
	var keyR, valueR string

	row := s.db.QueryRow(`SELECT key, value FROM dnscache WHERE key=$1;`, key)
	switch err = row.Scan(&keyR, &valueR); err {
	case sql.ErrNoRows:
	case nil:
		fmt.Printf("domain %v is on cache with IP %v\n", keyR, valueR)
		isExist = true
	default:
		return err
	}

	if isExist {
		_, err = s.db.Exec("UPDATE dnscache set value=? WHERE key=?", value, key)
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec("INSERT INTO dnscache VALUES (?,?);", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SqliteHandler) Close() error {
	s.db.Close()
	return nil
}
