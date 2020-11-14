package database

import (
	"database/sql"
	"fmt"
)

// ReadQuestionList 使用 u：用户名(ID) 查询 QuestionListTab 表。
// 列出所有答题
func (q *question) ReadQuestionList(u string) (tab []*QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.creator_id = "%v" ORDER BY question_list.id DESC`,
		u,
	)

	rows, err := q.conn.Query(sq)
	if err != nil {
		return
	}

	return joinQuestionList(rows)
}

// ReadQuestion 使用 i：问题ID(ID) 查询 QuestionListTab 表。
// 答题信息
func (q *question) ReadQuestion(i uint32) (data *QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.id = %v`,
		i,
	)

	row, err := q.conn.Query(sq)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		return
	}
	data = new(QuestionListTab)
	if row.Scan(
		&data.ID, &data.Type, &data.Subject, &data.Question, &data.CreatorID, &data.Target, &data.Status, &data.Options,
		&data.Key, &data.Market,
	) != nil {
		return
	}

	return
}

// ReadQuestionMarket 使用 s：学科(Subject) 查询 QuestionListTab 表。
// 答题市场
func (q *question) ReadQuestionMarket(s uint8) (tab []*QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.market = TRUE AND question_list.subject = %v ORDER BY question_list.id DESC`,
		s,
	)

	rows, err := q.conn.Query(sq)
	if err != nil {
		return
	}

	return joinQuestionList(rows)
}

// joinQuestionList 复用
func joinQuestionList(rows *sql.Rows) (tab []*QuestionListTab, err error) {

	defer rows.Close()
	var data []*QuestionListTab
	for rows.Next() {

		data0 := new(QuestionListTab)
		if rows.Scan(
			&data0.ID, &data0.Type, &data0.Subject, &data0.Question, &data0.CreatorID, &data0.Target, &data0.Status,
			&data0.Options, &data0.Key, &data0.Market,
		) != nil {
			return
		}

		data = append(data, data0)

	}

	tab = data
	return
}

// WriteQuestionList 写入 QuestionListTab 表。
// 新建答题
func (q *question) WriteQuestionList(tab *QuestionListTab) (err error) {

	i, err := q.conn.Prepare(
		"INSERT INTO question_list (id, type, subject, question, creator_id, target, status, options, `key`, market) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return
	}
	defer i.Close()

	// Tips： ID 自增无需输入、Status 默认为 0

	if _, err = i.Exec(
		nil, tab.Type, tab.Subject, tab.Question, tab.CreatorID, tab.Target, 0, tab.Options, tab.Key, tab.Market,
	); err != nil {
		return
	}

	return
}

// UpdateQuestion 使用 i：问题ID(ID) 、 s：状态码(Status) 更新问题 status 字段。
// 更新状态
func (q *question) UpdateQuestion(i uint32, s uint8) (err error) {

	l, err := q.conn.Prepare(`UPDATE question_list SET question_list.status = ? WHERE question_list.id = ?`)
	if err != nil {
		return
	}
	defer l.Close()

	if _, err = l.Exec(s, i); err != nil {
		return
	}

	return
}
