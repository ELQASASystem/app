package class

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Mysql
)

type (
	problemListTab struct {
		id      uint64 // 唯一标识符
		problem string // 问题
	}
)

var db *sql.DB

// connectDB 连接数据库
func connectDB(u string) (err error) {

	db, err = sql.Open("mysql", u)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return

}

// readProblemList 读取 ProblemList 表
func readProblemList() (tab *problemListTab, err error) {

	return
}
