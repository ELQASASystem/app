package class

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Mysql
)

type (
	// questionListTab 问题列表表
	questionListTab struct {
		ID        string `json:"id"`         // 唯一标识符
		Question  string `json:"question"`   // 问题
		CreatorID string `json:"creator_id"` // 创建者
		Market    bool   `json:"market"`     // 进入市场
	}
)

var db *sql.DB // 数据库

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

// readQuestionList 读取 questionList 表[教师方]。
// u => 用户名
func readQuestionList(u string) (tab []*questionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT problem_list.* FROM problem_list WHERE problem_list.creator_id = "%v" ORDER BY problem_list.id DESC`,
		u,
	)

	rows, err := db.Query(sq)
	if err != nil {
		return
	}
	defer rows.Close()

	tab, err = writeQuestionList(rows)
	return

}

// readQuestionMarket 读取 questionList 表[市场方]
func readQuestionMarket() (tab []*questionListTab, err error) {

	rows, err := db.Query(`SELECT problem_list.* FROM problem_list WHERE problem_list.market = true ORDER BY problem_list.id DESC`)
	if err != nil {
		return
	}
	defer rows.Close()

	tab, err = writeQuestionList(rows)
	return

}

// writeQuestionList 复用
func writeQuestionList(rows *sql.Rows) (tab []*questionListTab, err error) {

	data := make([]*questionListTab, 0)
	for rows.Next() {

		data0 := new(questionListTab)
		err = rows.Scan(
			&data0.ID, &data0.Question, &data0.CreatorID, &data0.Market,
		)
		if err != nil {
			return
		}

		data = append(data, data0)

	}

	tab = data
	return

}
