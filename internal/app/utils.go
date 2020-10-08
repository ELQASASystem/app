package class

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/ELQASASystem/app/internal/app/qq"
)

// chainToString 消息链转文本
func chainToString(chain []qq.Chain) (fullText string) {

	for _, element := range chain {
		fullText += element.Text
	}
	return

}

// isValidAnswer 是否为合法答案 [选择题]
func isValidAnswer(answer string) (ok bool) {

	ok, _ = regexp.MatchString("[a-zA-Z]", answer)
	return

}

// parseAnswer 解析消息中的答案, 并返回 Answer 结构体
func parseAnswer(m *qq.Msg, aid uint64) (*Answer, bool) {

	s := chainToString(m.Chain)

	if isValidAnswer(s) && !isAnswered(m.Group.ID, m.User.ID) {
		return &Answer{s, m.User.ID, aid}, true
	} else {
		return nil, false
	}

}

// isAnswered 检查对应 QQ 号用户是否已经答题过了
func isAnswered(aid uint64, qid uint64) bool {

	if question, _, ok := getQuestionByID(aid); ok {
		for _, user := range question.AnsweredUsers {
			if user.Sender == qid {
				return true
			}
		}
	}
	return false

}

// getQuestionByID 通过问题 ID 获取问题实体
func getQuestionByID(aid uint64) (*Question, int, bool) {
	for i, question := range questionPool {
		if question.QuestionID == aid {
			return &question, i, true
		}
	}
	return nil, 0, false
}

// getQuestionByGroup 通过群号获取问题实体
func getQuestionByGroup(qid uint64) (*Question, int, bool) {
	for index, question := range questionPool {
		if question.TargetGroup == qid {
			return &question, index, true
		}
	}
	return nil, 0, false
}

// hashSHA1 将答题数据散列
func hashSHA1(data interface{}) string {

	h := sha1.New()

	h.Write([]byte(fmt.Sprintf("%v", data)))

	return hex.EncodeToString(h.Sum(nil))
}

// makeQuestion 制作一个 Question 结构类
func makeQuestion(qid uint64, gid uint64, text string) *Question {
	return &Question{
		QuestionID:   qid,
		QuestionText: text,
		TargetGroup:  gid,
	}
}
