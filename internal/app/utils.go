package class

import (
	"crypto/sha1"
	"fmt"
	"regexp"
)

// chainToString 消息链转文本
func chainToString(chain []Chain) (fullText string) {

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
func parseAnswer(m *QQMsg, aid uint64) (*Answer, bool) {

	s := chainToString(m.Chain)

	if isValidAnswer(s) && !isAnswered(m.Group.ID, m.User.ID) {
		return &Answer{s, m.User.ID, aid}, true
	} else {
		return nil, false
	}

}

// isAnswered 检查对应 QQ 号用户是否已经答题过了
func isAnswered(gid uint64, qid uint64) bool {

	if question, ok := questionPool[gid]; ok {
		for _, user := range question.AnsweredUsers {
			if user.Sender == qid {
				return true
			}
		}
	}
	return false

}

// hashSHA1 将答题数据散列
func hashSHA1(data *Question) string {

	s := sha1.New()

	s.Write([]byte(fmt.Sprintf("%v", data)))

	return fmt.Sprintf("%x", s.Sum(nil))
}
