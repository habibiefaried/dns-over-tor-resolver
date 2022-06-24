package cachehandler

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	FileName          string
	ExpireAgeinSecond int
	db                *sql.DB
}

func (s *SqliteHandler) Init() error {
	var err error
	s.db, err = sql.Open("sqlite3", s.FileName)
	if err != nil {
		return err
	}

	const create string = `
	CREATE TABLE IF NOT EXISTS dnscache (
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		source TEXT NOT NULL,
		time LONG INTEGER NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS dnskey ON dnscache(key);
	`

	if _, err := s.db.Exec(create); err != nil {
		return err
	}

	if s.ExpireAgeinSecond == 0 {
		s.ExpireAgeinSecond = 3600 // default value is 3600 seconds
	}

	return nil
}

func (s *SqliteHandler) Get(key string) ([]string, error) {
	var value string
	var source string
	ret := []string{}

	rows, err := s.db.Query(fmt.Sprintf(`SELECT value, source FROM dnscache WHERE key=? and time > strftime('%s', dateTIME('now','-%v second'))`, "%s", s.ExpireAgeinSecond), key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&value, &source)
		if err != nil {
			fmt.Printf("error while retrieving data: %v\n", err)
			continue
		} else {
			ret = append(ret, value)
			fmt.Printf("[cache HIT] domain %v is found %v from %v\n", key, value, source)
		}
	}

	return ret, nil
}

func (s *SqliteHandler) Put(key, value, source string, unixtimestamp int64) error {
	var err error

	_, err = s.db.Exec("INSERT INTO dnscache VALUES (?,?,?,?);", key, value, source, unixtimestamp)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteHandler) Close() error {
	s.db.Close()
	return nil
}

func (s *SqliteHandler) CleanUp() (int64, error) {
	stmt, err := s.db.Prepare(fmt.Sprintf(`DELETE FROM dnscache where time <= strftime('%s', dateTIME('now','-%v second'))`, "%s", s.ExpireAgeinSecond))
	if err != nil {
		return 0, fmt.Errorf("fail to build statement: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		return 0, fmt.Errorf("fail to delete, got error: %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to count rows affectd: %v", err)
	}

	return affected, nil
}
