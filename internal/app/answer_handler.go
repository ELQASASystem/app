package class

import (
	"github.com/ELQASASystem/app/internal/app/qq"
)

// Public 使用 i：问题ID(ID) 发布问答问题
func Public(i string) error {

	// TODO 用i查询数据库

	// 发送群消息

	return nil
}

// StartQA 使用 i：问题ID(ID) 开始监听问题
func StartQA(i string) error {

	return nil
}

// uploadUserAnswer 上报用户答案
func uploadUserAnswer() {
	//TODO
}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *qq.Msg) {}

// StopQA 注销问题, 返回该问题和是否注销成功
func StopQA(qid uint64) (ok bool) {

	// 删除
	return false
}
