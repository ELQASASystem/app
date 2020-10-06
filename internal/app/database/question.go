package database

import (
	"database/sql"
	"fmt"
)

// ReadQuestionList 使用 u：用户名(ID) 查询 QuestionListTab 表。
// 列出所有答题
func (q Question) ReadQuestionList(u string) (tab []*QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.creator_id = "%v" ORDER BY question_list.id DESC`,
		u,
	)

	rows, err := Class.DB.Query(sq)
	if err != nil {
		return
	}

	tab, err = joinQuestionList(rows)
	return

}

// ReadQuestion 使用 i：问题ID(ID) 查询 QuestionListTab 表。
// 答题信息
func (q Question) ReadQuestion(i uint32) (data *QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.id = %v`,
		i,
	)

	row, err := Class.DB.Query(sq)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		return
	}
	data = new(QuestionListTab)
	err = row.Scan(&data.ID, &data.Type, &data.Question, &data.CreatorID, &data.Target, &data.Status, &data.Options,
		&data.Key, &data.Market)
	if err != nil {
		return
	}

	return

}

// ReadQuestionMarket 查询 QuestionListTab 表。
// 答题市场
func (q Question) ReadQuestionMarket() (tab []*QuestionListTab, err error) {

	rows, err := Class.DB.Query(
		`SELECT question_list.* FROM question_list WHERE question_list.market = true ORDER BY question_list.id DESC`,
	)
	if err != nil {
		return
	}

	tab, err = joinQuestionList(rows)
	return

}

// joinQuestionList 复用
func joinQuestionList(rows *sql.Rows) (tab []*QuestionListTab, err error) {

	defer rows.Close()
	var data []*QuestionListTab
	for rows.Next() {

		data0 := new(QuestionListTab)
		err = rows.Scan(
			&data0.ID, &data0.Type, &data0.Question, &data0.CreatorID, &data0.Target, &data0.Status, &data0.Options,
			&data0.Key, &data0.Market)
		if err != nil {
			return
		}

		data = append(data, data0)

	}

	tab = data
	return

}

// WriteQuestionList 写入 QuestionListTab 表。
// 新建答题
func (q Question) WriteQuestionList(tab *QuestionListTab) (err error) {

	i, err := Class.DB.Prepare(
		`INSERT INTO question_list 
(id, type, question, creator_id, target, status, options, key, market) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return
	}
	defer i.Close()

	// ID 自增无需输入
	// Status 默认为 0
	_, err = i.Exec(nil, tab.Type, tab.Question, tab.CreatorID, tab.Target, 0, tab.Options, tab.Key, tab.Market)
	if err != nil {
		return
	}

	return

}
