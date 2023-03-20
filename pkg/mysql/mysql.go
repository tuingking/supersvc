package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL interface {
	Stop() error
	Get() *sql.DB
}

type mySQL struct {
	opt Option
	db  *sql.DB
}

type Option struct {
	DriverName       string
	ConnectionString string
}

func NewMySQL(opt Option) MySQL {
	return &mySQL{
		opt: opt,
		db:  initConnection(opt),
	}
}

func (m *mySQL) Stop() error {
	return m.db.Close()
}

func (m *mySQL) Get() *sql.DB {
	return m.db
}

func initConnection(opt Option) *sql.DB {
	db, err := sql.Open(opt.DriverName, opt.ConnectionString+"?parseTime=true")
	if err != nil {
		log.Fatalf("[mysql] unable connect to database. err=%s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("[mysql] ping failed. err=%s", err)
	}

	log.Print("[mysql] database connection success")

	return db
}
