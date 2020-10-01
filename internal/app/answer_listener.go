package class

var (
	// 答题数据储存池, key 为群号, value 为问题 ID
	// 一群对应一个问题 ID
	questionPool = make(map[uint64]Question)
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

// 注销问题, 返回是否注销成功
func expiredQuestion(groupID uint64, aid uint64) bool {

	if v, ok := questionPool[groupID]; ok && v.QuestionID == aid {
		delete(questionPool, groupID)
		return true
	} else {
		return false
	}

}

func publishQuestion(q *Question) bool {
	if q.TargetGroup != 0 {
		classBot.SendGroupMsg(NewText(q.QuestionText).To(q.TargetGroup))
		return true
	}
	return false
}

// uploadUserAnswer 上报用户答案
func uploadUserAnswer(groupId uint64, ans *Answer) {

	if v, ok := questionPool[groupId]; ok {
		v.AnsweredUsers = append(v.AnsweredUsers, *ans)
		// TODO: 记得再上报给 Web 端
	}
}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *QQMsg) {

	groupId := m.Group.ID

	if question, ok := questionPool[groupId]; ok {
		if ans, ok := parseAnswer(m, question.QuestionID); ok {
			uploadUserAnswer(m.Group.ID, ans)
		}
	}

}
