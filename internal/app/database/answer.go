package database

import "fmt"

// ReadAnswerList 使用 i：问题ID(ID) 查询 AnswerListTab 表。
// 回答信息
func (a *answer) ReadAnswerList(i uint32) (tab []*AnswerListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT answer_list.* FROM answer_list WHERE answer_list.question_id = %v ORDER BY answer_list.time DESC`,
		i,
	)

	rows, err := a.conn.Query(sq)
	if err != nil {
		return
	}
	defer rows.Close()

	var data []*AnswerListTab
	for rows.Next() {

		data0 := new(AnswerListTab)
		err = rows.Scan(
			&data0.ID, &data0.QuestionID, &data0.AnswererID, &data0.Answer, &data0.Time,
		)
		if err != nil {
			return
		}

		data = append(data, data0)

	}

	tab = data
	return
}

// WriteAnswerList 写入 AnswerListTab 表。
// 写入回答
func (a *answer) WriteAnswerList(d *AnswerListTab) (err error) {

	i, err := a.conn.Prepare(
		`INSERT INTO answer_list (id, question_id, answerer_id, answer, time) VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return
	}
	defer i.Close()

	_, err = i.Exec(nil, d.QuestionID, d.AnswererID, d.Answer, d.Time)
	if err != nil {
		return
	}

	return
}
