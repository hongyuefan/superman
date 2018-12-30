package database

import (
	_ "database/sql"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

/*
CREATE TABLE IF NOT EXISTS kline
    (
		id				INTEGER PRIMARY KEY	AUTOINCREMENT,
		open			VARCHAR(32)	NOT NULL,
		high			VARCHAR(32)	NOT NULL,
		low				VARCHAR(32)	NOT NULL,
		close			VARCHAR(32)	NOT NULL,
		deal			VARCHAR(32)	NOT NULL,
		time			VARCHAR(32)	NOT NULL
	);

*/
func RegistDB() {

	orm.RegisterDriver("sqlite", orm.DRSqlite)

	orm.RegisterDataBase("default", "sqlite3", "./storedb.db")

}
