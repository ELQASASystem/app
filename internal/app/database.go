package class

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Mysql
)

type (
	// questionListTab 问题列表表
	questionListTab struct {
		ID        uint32 `json:"id"`         // 唯一标识符
		Question  string `json:"question"`   // 问题
		CreatorID string `json:"creator_id"` // 创建者
		Market    bool   `json:"market"`     // 进入市场
	}

	// accountsList 帐号列表表
	accountsList struct {
		ID       string `json:"id"`       // 唯一标识符
		Password string `json:"password"` // 密码
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

// readAccountsList 读取 accountsList 表。
// u => 用户名
func readAccountsList(u string) (data *accountsList, err error) {

	sq := fmt.Sprintf(
		`SELECT accounts_list.* FROM accounts_list WHERE accounts_list.id = "%v"`,
		u,
	)

	row, err := db.Query(sq)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		return
	}
	data = new(accountsList)
	err = row.Scan(&data.ID, &data.Password)
	if err != nil {
		return
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

	tab, err = joinQuestionList(rows)
	return

}

// readQuestionMarket 读取 questionList 表[市场方]
func readQuestionMarket() (tab []*questionListTab, err error) {

	rows, err := db.Query(`SELECT problem_list.* FROM problem_list WHERE problem_list.market = true ORDER BY problem_list.id DESC`)
	if err != nil {
		return
	}
	defer rows.Close()

	tab, err = joinQuestionList(rows)
	return

}

// joinQuestionList 复用
func joinQuestionList(rows *sql.Rows) (tab []*questionListTab, err error) {

	var data []*questionListTab
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

// writeQuestionList 写入 questionList 表
func writeQuestionList(d *questionListTab) (err error) {

	i, err := db.Prepare(`INSERT INTO problem_list (id, question, creator_id, market) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return
	}
	defer i.Close()

	_, err = i.Exec(nil, d.Question, d.CreatorID, d.Market)
	if err != nil {
		return
	}

	return

}
