package class

import (
	"github.com/ELQASASystem/app/internal/app/qq"

	"github.com/rs/zerolog/log"
)

var (
	// 答题数据储存池
	// 如需获取对应问题, 请使用 getQuestionByID 或 getQuestionByGroup 方法
	questionPool []Question
)

// Question 问题数据
type Question struct {
	QuestionID     uint64   // 问题 ID
	QuestionText   string   // 问题发布时使用的文本
	QuestionAnswer string   // 问题的答案
	TargetGroup    uint64   // 问题发布的目标群聊
	AnsweredUsers  []Answer // 回答的答案
}

// Answer 答题数据
type Answer struct {
	Text     string // 回答者的答案
	Sender   uint64 // 回答者 ID
	AnswerID uint64 // 问题 ID
}

// 注销问题, 返回该问题和是否注销成功
func expiredQuestion(qid uint64) (q *Question, ok bool) {

	if v, i, ok := getQuestionByID(qid); ok && v.QuestionID == qid {
		questionPool = append(questionPool[:i], questionPool[i+1:]...)
		return v, ok
	} else {
		return nil, false
	}

}

// publishQuestion 发布问题开始答题
func publishQuestion(q *Question) bool {
	if q != nil {
		classBot.SendGroupMsg(classBot.NewText(q.QuestionText).To(q.TargetGroup))
		return true
	}
	return false
}

// uploadUserAnswer 上报用户答案
func uploadUserAnswer(groupId uint64, ans *Answer) {

	if v, _, ok := getQuestionByGroup(groupId); ok {
		v.AnsweredUsers = append(v.AnsweredUsers, *ans)

		// 检查是否有客户端正在监听此问题
		if conn, ok := getConnByQID(uint32(v.QuestionID)); ok {
			if err := conn.conn.WriteMessage(conn.mt, []byte(hashSHA1(v.AnsweredUsers))); err != nil {
				log.Warn().Err(err).Msg("上报答案失败")
			}
		}

		// TODO: 记得再上报给 Web 端
	}
}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *qq.Msg) {

	if question, _, ok := getQuestionByGroup(m.Group.ID); ok {
		if ans, ok := parseAnswer(m, question.QuestionID); ok {
			uploadUserAnswer(m.Group.ID, ans)
		}
	}

}
