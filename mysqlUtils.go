package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "go_test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func ExecuteQuery(query string, args ...interface{}) {
	db := dbConn()
	dbPrepare, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}

	result, err := dbPrepare.Exec(args...)
	if err != nil {
		panic(err.Error())
	}
	if result != nil {

	}
	defer db.Close()
}

func queryRow(query string, rs RowScanner, args ...interface{}) error {
	db := dbConn()
	err := rs.ScanRow(db.QueryRow(query, args...))
	defer db.Close()
	return err
}

func queryRows(query string, rs RowScanner, args ...interface{}) error {
	db := dbConn()
	rows, err := db.Query(query, args...)
	if err != nil {
		panic(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rs.ScanRow(rows); err != nil {
			panic(err.Error())
			return err
		}
	}
	defer db.Close()
	return rows.Err()
}

type Employee struct {
	Id   int
	Name string
	City string
	// ...
}

type EmployeeList struct {
	Items []*Employee
}

type RowScanner interface {
	ScanRow(Row) error
}
type Row interface {
	Scan(...interface{}) error
}

// Implements RowScanner
func (u *Employee) ScanRow(r Row) error {
	return r.Scan(
		&u.Id,
		&u.Name,
		&u.City,
		// ...
	)
}

// Implements RowScanner
func (list *EmployeeList) ScanRow(r Row) error {
	u := new(Employee)
	if err := u.ScanRow(r); err != nil {
		return err
	}
	list.Items = append(list.Items, u)
	return nil
}
