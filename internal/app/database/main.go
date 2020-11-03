package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Mysql
)

type (

	/*

		方法组织

	*/

	// Database 数据库
	Database struct {
		conn *sql.DB // conn 数据库连接
	}

	// account 帐号相关
	account struct {
		conn *sql.DB // conn 数据库连接
	}

	// question 问题相关
	question struct {
		conn *sql.DB // conn 数据库连接
	}

	// answer 回答相关
	answer struct {
		conn *sql.DB // conn 数据库连接
	}

	/*

		数据结构

	*/

	// AccountsListTab 帐号
	AccountsListTab struct {
		ID       string `json:"id"`       // ID 唯一标识符
		Password string `json:"password"` // Password 密码
		Class    string `json:"class"`    // Class 班级
	}

	// QuestionListTab 问题
	QuestionListTab struct {
		ID        uint32 `json:"id"`         // ID 唯一标识符
		Type      uint8  `json:"type"`       // Type 类型
		Subject   uint8  `json:"subject"`    // Subject 学科
		Question  string `json:"question"`   // Question 问题
		CreatorID string `json:"creator_id"` // CreatorID 创建者
		Target    uint64 `json:"target"`     // Target 目标
		Status    uint8  `json:"status"`     // Status 状态
		Options   string `json:"options"`    // Options 选项
		Key       string `json:"key"`        // Key 答案
		Market    bool   `json:"market"`     // Market 存在市场
	}

	// AnswerListTab 回答
	AnswerListTab struct {
		ID         uint32 `json:"id"`          // ID 唯一标识符
		QuestionID uint32 `json:"question_id"` // QuestionID 问题唯一标识符
		AnswererID uint64 `json:"answerer_id"` // AnswererID 回答者
		Answer     string `json:"answer"`      // Answer 回答内容
		Time       string `json:"time"`        // Time 回答时间
	}
)

// New 新建一个数据库事务
func New() *Database { return new(Database) }

// ConnectDB 连接数据库
func (d *Database) ConnectDB(u string) (err error) {

	d.conn, err = sql.Open("mysql", u)
	if err != nil {
		return
	}

	return d.conn.Ping()
}

// Account 帐号相关
func (d *Database) Account() *account { return &account{conn: d.conn} }

// Question 问题相关
func (d *Database) Question() *question { return &question{conn: d.conn} }

// Answer 回答相关
func (d *Database) Answer() *answer { return &answer{conn: d.conn} }
