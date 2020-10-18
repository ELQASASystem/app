package class

import (
	"strings"
	"time"

	"github.com/ELQASASystem/app/internal/app/database"
	"github.com/ELQASASystem/app/internal/app/qq"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

var QABasicSrvPoll = map[uint64]*database.QuestionListTab{} // QABasicSrvPoll 问答基本服务线程池

// StartQA 使用 i：问题ID(ID) 开始监听问题
func StartQA(i uint32) (err error) {

	q, err := getQuestion(i)
	if err != nil {
		return
	}

	var (
		options []struct {
			Type string `json:"type"` // 选项号
			Body string `json:"body"` // 选项内容
		}
		m = Bot.NewMsg().AddText("问题:\n").AddText(q.Question).AddText("\n选项:\n")
	)

	if err = jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(q.Options, &options); err != nil {
		log.Error().Err(err).Msg("解析选项失败")
		return
	}
	for _, v := range options {
		m.AddText(v.Type + ". " + v.Body + "\n")
	}

	if q.Type == 0 {
		m.AddText("\n回复选项即可作答")
	} else {
		m.AddText("\n@+回答内容即可作答")
	}

	Bot.SendGroupMsg(m.To(q.Target))
	QABasicSrvPoll[q.Target] = q

	return
}

// StopQA 使用 i：问题ID(ID) 停止问答
func StopQA(i uint32) (err error) {

	err = deleteQABasicSrvPoll(i)
	err = database.Class.Question.UpdateQuestion(i, 2)
	if err != nil {
		return
	}

	return
}

// PrepareQA 使用 i：问题ID(ID) 开始准备作答
func PrepareQA(i uint32) (err error) {

	err = deleteQABasicSrvPoll(i)
	err = database.Class.Question.UpdateQuestion(i, 0)
	if err != nil {
		return
	}

	return
}

type Question struct {
	*database.QuestionListTab
	Answer []*database.AnswerListTab `json:"answer"`
}

// ReadQuestion 使用 i：问题ID(ID) 读取问答信息
func ReadQuestion(i uint32) (q *Question, err error) {

	res, err := database.Class.Question.ReadQuestion(i)
	if err != nil {
		return
	}

	res2, err := database.Class.Answer.ReadAnswerList(i)
	if err != nil {
		return
	}

	return &Question{res, res2}, nil
}

// reportUserAnswer 上报用户答案
func reportUserAnswer(q *database.QuestionListTab, m *qq.Msg) {

	err := database.Class.Answer.WriteAnswerList(&database.AnswerListTab{
		QuestionID: q.ID,
		AnswererID: m.User.ID,
		Answer:     strings.ToUpper(m.Chain[0].Text),
		Time:       time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		log.Warn().Err(err).Msg("写入答案失败")
		return
	}

	log.Info().Msg("成功写入回答")

	res, err := ReadQuestion(q.ID)
	if err != nil {
		log.Error().Err(err).Msg("读取问题信息失败")
	}

	log.Info().Msg("准备发送回答数据")

	qch <- res

}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *qq.Msg) {

	q, ok := QABasicSrvPoll[m.Group.ID]
	if !ok {
		return
	}

	// TODO 改内存实现
	ans, err := database.Class.Answer.ReadAnswerList(q.ID)
	if err != nil {
		log.Warn().Err(err).Msg("读取答案列表失败")
		return
	}

	for _, v := range ans {
		if v.AnswererID == m.User.ID {
			return
		}
	}

	if isAnswer(m.Chain[0].Text) {
		reportUserAnswer(q, m)
	}

}

// getQuestion 使用 i：问题ID(ID) 获取问题
func getQuestion(i uint32) (q *database.QuestionListTab, err error) {

	q, err = database.Class.Question.ReadQuestion(i)
	if err != nil {
		return
	}

	return
}

// deleteQABasicSrvPoll 使用 i：问题ID(ID) 删除问答基本服务池字段
func deleteQABasicSrvPoll(i uint32) (err error) {

	q, err := getQuestion(i)
	if err != nil {
		return
	}

	delete(QABasicSrvPoll, q.Target)
	return
}
