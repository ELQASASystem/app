package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Mysql
)

type (

	/*

		方法组织

	*/

	// Str 数据库包结构体
	Str struct {
		DB       *sql.DB  // DB 数据库指针
		Account  Account  // Account 帐号相关
		Question Question // Question 问题相关
		Answer   Answer   // Answer 回答相关
	}

	Account  struct{} // Account 帐号相关
	Question struct{} // Question 问题相关
	Answer   struct{} // Answer 回答相关

	/*

		数据结构

	*/

	// AccountsListTab 帐号
	AccountsListTab struct {
		ID       string `json:"id"`       // 唯一标识符
		Password string `json:"password"` // 密码
	}

	// QuestionListTab 问题
	QuestionListTab struct {
		ID        uint32 `json:"id"`         // 唯一标识符
		Question  string `json:"question"`   // 问题
		CreatorID string `json:"creator_id"` // 创建者
		Status    uint8  `json:"status"`     // 状态
		Market    bool   `json:"market"`     // 进入市场
	}

	// AnswerListTab 回答
	AnswerListTab struct {
		ID         uint32 `json:"id"`          // 唯一标识符
		QuestionID uint32 `json:"question_id"` // 问题唯一标识符
		AnswererID uint64 `json:"answerer_id"` // 回答者
		Answer     string `json:"answer"`      // 回答内容
		Time       string `json:"time"`        // 回答时间
	}
)

var Class = Str{} // Class 数据库相关

// ConnectDB 连接数据库
func (s *Str) ConnectDB(u string) (err error) {

	s.DB, err = sql.Open("mysql", u)
	if err != nil {
		return
	}

	err = s.DB.Ping()
	if err != nil {
		return
	}

	return

}
