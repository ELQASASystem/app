package class

import (
	"github.com/ELQASASystem/app/internal/app/database"
	"github.com/ELQASASystem/app/internal/app/qq"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

// 问答基本服务线程池
var QABasicSrvPoll = map[uint64]database.QuestionListTab{}

// StartQA 使用 i：问题ID(ID) 开始监听问题
func StartQA(i string) error {
	if q, err := getQuestion(i); err != nil {
		return err
	} else {
		QABasicSrvPoll[uint64(q.ID)] = *q

		// 发送群消息
		m := Bot.NewMsg().AddText("问题:\n").AddText(q.Question).AddText("\n选项:\n").AddText(q.Options).AddText("\n直接回复答案即可作答").To(q.Target)
		Bot.SendGroupMsg(m)
	}
	return nil
}

// PrepareQA 使用 i：问题ID(ID) 开始准备作答
func PrepareQA(i string) error {
	q, err := getQuestion(i)

	if err != nil {
		return err
	} else {
		q.Status = 0
		if err := database.Class.Question.WriteQuestionList(q); err != nil {
			return err
		}
	}

	return nil
}

// reportUserAnswer 上报用户答案
func reportUserAnswer(q *database.QuestionListTab, m *qq.Msg) {
	if err := database.Class.Answer.WriteAnswerList(&database.AnswerListTab{
		QuestionID: q.ID,
		AnswererID: m.User.ID,
		Answer:     m.Chain[0].Text,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		log.Warn().Err(err).Msg("写入答案时出现异常")
		return
	}

	// TODO: Websocket 上报
}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *qq.Msg) {
	if q, ok := QABasicSrvPoll[m.Group.ID]; !ok {
		return
	} else {
		if ans, err := database.Class.Answer.ReadAnswerList(q.ID); err != nil {
			log.Warn().Err(err).Msg("读取答案列表时出现异常")
			return
		} else {
			for _, an := range ans {
				if an.AnswererID == m.User.ID {
					return
				}
			}

			if !isValidAnswer(m.Chain[0].Text) {
				return
			}

			reportUserAnswer(&q, m)
		}
	}
}

// StopQA 注销问题, 返回该问题和是否注销成功
func StopQA(qid uint64) (ok bool) {
	// 检查该问题 ID 对应问题是否在服务线程池中
	if _, ok := QABasicSrvPoll[qid]; !ok {
		return false
	}

	// 如果存在则删除
	delete(QABasicSrvPoll, qid)
	return true
}

// getQuestion 通过问题ID(ID) 获取问题 复用
func getQuestion(i string) (*database.QuestionListTab, error) {
	id, err := strconv.ParseUint(i, 10, 32)

	if err != nil {
		return nil, err
	}

	if q, err := database.Class.Question.ReadQuestion(uint32(id)); err != nil {
		return nil, err
	} else {
		return q, nil
	}
}
