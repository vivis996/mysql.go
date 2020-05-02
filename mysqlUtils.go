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

func queryRow(query string, rs RowScanner, params ...interface{}) error {
	db := dbConn()
	return rs.ScanRow(db.QueryRow(query, params...))
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

// example
// ulist := new(UserList)
// if err := queryRows(queryString, ulist, arg1, arg2); err != nil {
//     panic(err)
// }

// or
// u := new(User)
// if err := queryRow(queryString, u, arg1, arg2); err != nil {
//     panic(err)
// }
