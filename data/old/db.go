package data

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
mysql -u root -p (quick2listen)
CREATE USER 'comrad'@localhost' IDENTIFIED BY 'xyz!123';
CREATE DATABASE comrad
GRANT ALL PRIVILEGES ON comrad . * TO 'comrad'@'localhost';
*/

var db *sqlx.DB

func Initialize(username, password, ip, port, databasename string) error {
	var err error
	db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, ip, port, databasename))
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return err
	}
	return nil
}

func CreateTables() error {
	sqlStmt := `create table user (id varchar(50) not null primary key, email varchar(50) unique, hashword varchar(100), role integer, since datetime);`

	/*	`create table project (id integer not null primary key, name text, description text, lon text, lat text, start_date text, end_date text, notify_list text, customer_email text, communityId integer, creatorId integer, open integer);
		create table users_project (user_id integer not null primary key, project_id integer, role text);
		create table card (id integer not null primary key, project_id text, description text, start-date text, end-date text);` */
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sqlStmt = `create table post (id varchar(50) not null primary key, authorid varchar(50), title varchar(50), pagestub varchar(50), content varchar(5000), date datetime);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func DropTables() error {
	sqlStmts := []string{`drop table user;`, `drop table post;`}
	for _, stmt := range sqlStmts {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func GetDb() *sqlx.DB {
	return db
}
