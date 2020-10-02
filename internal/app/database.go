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

	// accountsListTab 帐号列表表
	accountsListTab struct {
		ID       string `json:"id"`       // 唯一标识符
		Password string `json:"password"` // 密码
	}

	// answerListTab 答题列表表
	answerListTab struct {
		SHA1       string `json:"sha1"`        // 散列值
		QuestionID uint32 `json:"question_id"` // 问题唯一标识符
		AnswererID uint64 `json:"answerer_id"` // 回答者
		Answer     string `json:"answer"`      // 回答内容
		Time       string `json:"time"`        // 回答时间
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

// readAccountsList 读取 accountsListTab 表。
// u => 用户名
func readAccountsList(u string) (data *accountsListTab, err error) {

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
	data = new(accountsListTab)
	err = row.Scan(&data.ID, &data.Password)
	if err != nil {
		return
	}

	return

}

// readQuestionList 读取 questionListTab 表[教师方]。
// u => 用户名
func readQuestionList(u string) (tab []*questionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.creator_id = "%v" ORDER BY question_list.id DESC`,
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

// readQuestion 读取指定问题。
// i => 问题ID
func readQuestion(i uint32) (data *questionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.id = %v`,
		i,
	)

	row, err := db.Query(sq)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		return
	}
	data = new(questionListTab)
	err = row.Scan(&data.ID, &data.Question, &data.CreatorID, &data.Market)
	if err != nil {
		return
	}

	return

}

// readQuestionMarket 读取 questionListTab 表[市场方]
func readQuestionMarket() (tab []*questionListTab, err error) {

	rows, err := db.Query(`SELECT question_list.* FROM question_list WHERE question_list.market = true ORDER BY question_list.id DESC`)
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

// writeQuestionList 写入 questionListTab 表
func writeQuestionList(d *questionListTab) (err error) {

	i, err := db.Prepare(`INSERT INTO question_list (id, question, creator_id, market) VALUES (?, ?, ?, ?)`)
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

// readAnswerList 写入 answerListTab 表
func readAnswerList(i uint32) (tab []*answerListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT answer_list.* FROM answer_list WHERE answer_list.question_id = %v ORDER BY answer_list.time DESC`,
		i,
	)

	rows, err := db.Query(sq)
	if err != nil {
		return
	}
	defer rows.Close()

	var data []*answerListTab
	for rows.Next() {

		data0 := new(answerListTab)
		err = rows.Scan(
			&data0.SHA1, &data0.QuestionID, &data0.AnswererID, &data0.Answer, &data0.Time,
		)
		if err != nil {
			return
		}

		data = append(data, data0)

	}

	tab = data
	return

}

// writeAnswerList 写入 answerListTab 表
func writeAnswerList(d *answerListTab) (err error) {

	i, err := db.Prepare(`INSERT INTO answer_list (sha1, question_id, answerer_id, answer, time) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}
	defer i.Close()

	_, err = i.Exec(d.SHA1, d.QuestionID, d.AnswererID, d.Answer, d.Time)
	if err != nil {
		return
	}

	return

}
