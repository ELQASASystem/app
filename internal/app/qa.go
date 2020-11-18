package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/ELQASASystem/app/internal/app/database"
	"github.com/ELQASASystem/app/internal/qq"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

// Question 问题
type Question struct {
	*database.QuestionListTab
	Mems      []*memInfo                `json:"mems"`
	GroupName string                    `json:"group_name"`
	Answer    []*database.AnswerListTab `json:"answer"`
}

// memInfo 成员信息
type memInfo struct {
	ID   uint64 `json:"id"`   // 群员帐号
	Name string `json:"name"` // 群员名片
}

var QABasicSrvPoll = map[uint64]*Question{} // QABasicSrvPoll 问答基本服务线程池

// StartQA 使用 i：问题ID(ID) 开始作答
func (a *App) StartQA(i uint32) (err error) {

	q, err := a.ReadQuestion(i)
	if err != nil {
		log.Error().Err(err).Msg("读取问题失败")
		return
	}

	if err = a.DB.Question().UpdateQuestion(i, 1); err != nil {
		log.Error().Err(err).Msg("更新问答状态字段失败")
		return
	}

	log.Info().Msg("问题 " + strconv.Itoa(int(i)) + " 开始监听")
	if err = a.sendQuestionMsg(q); err != nil {
		return
	}

	q.Status = 1
	QABasicSrvPoll[q.Target] = q
	return
}

// sendQuestionMsg 发送问答问题消息
func (a *App) sendQuestionMsg(q *Question) (err error) {

	var (
		question []struct {
			Type string `json:"type"` // 类型
			Text string `json:"text"` // 文本
			Path string `json:"path"` // 路径[图片]
		}
		options []struct {
			Type string `json:"type"` // 选项号
			Body string `json:"body"` // 选项内容
		}
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	if err = json.UnmarshalFromString(q.Question, &question); err != nil {
		log.Error().Err(err).Msg("解析问题失败")
		return
	}
	if err = json.UnmarshalFromString(q.Options, &options); err != nil {
		log.Error().Err(err).Msg("解析选项失败")
		return
	}

	m := a.Cli.NewText("问题:\n")
	for _, v := range question {
		if v.Type == "img" {
			m.AddImage("web/assets/question/pictures/" + v.Path).AddText("\n")
			continue
		}

		m.AddText(v.Text + "\n")
	}

	m.AddText("选项:\n")
	for _, v := range options {
		m.AddText(v.Type + ". " + v.Body + "\n")
	}

	if q.Type == 0 {
		m.AddText("\n回复选项即可作答")
	} else {
		m.AddText("\n@+回答内容即可作答")
	}

	a.Cli.SendGroupMsg(m.To(q.Target))
	return
}

// StopQA 使用 i：问题ID(ID) 停止问答
func (a *App) StopQA(i uint32) (err error) {

	if err = a.deleteQABasicSrvPoll(i); err != nil {
		log.Error().Err(err).Msg("删除问答基本服务监听失败")
		return
	}
	if err = a.DB.Question().UpdateQuestion(i, 2); err != nil {
		log.Error().Err(err).Msg("更新问答状态字段失败")
		return
	}

	log.Info().Msg("问题 " + strconv.Itoa(int(i)) + " 已停止答题")
	return
}

// PrepareQA 使用 i：问题ID(ID) 准备作答
func (a *App) PrepareQA(i uint32) (err error) {

	if err = a.deleteQABasicSrvPoll(i); err != nil {
		log.Error().Err(err).Msg("删除问答基本服务监听失败")
		return
	}
	if err = a.DB.Question().UpdateQuestion(i, 0); err != nil {
		log.Error().Err(err).Msg("更新问答状态字段失败")
		return
	}

	return
}

// ReadQuestion 使用 i：问题ID(ID) 读取问答信息
func (a *App) ReadQuestion(i uint32) (q *Question, err error) {

	res, err := a.DB.Question().ReadQuestion(i)
	if err != nil {
		return
	}

	answer, err := a.DB.Answer().ReadAnswerList(i)
	if err != nil {
		return
	}

	groupInfo := a.Cli.C.FindGroupByUin(int64(res.Target))
	mems := a.ReadMemInfo(uint64(groupInfo.Uin))

	return &Question{res, mems, groupInfo.Name, answer}, nil
}

// ReadMemInfo 使用 i：群ID(ID) 读取群成员信息
func (a *App) ReadMemInfo(i uint64) []*memInfo {

	var data []*memInfo
	for _, v := range a.Cli.C.FindGroupByUin(int64(i)).Members {
		data = append(data, &memInfo{uint64(v.Uin), v.DisplayName()})
	}

	return data
}

// writeAnswer 写入回答答案
func (a *App) writeAnswer(q *Question, stu uint64, ans string) {

	answer := &database.AnswerListTab{
		QuestionID: q.ID,
		AnswererID: stu,
		Answer:     ans,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
	}

	err := a.DB.Answer().WriteAnswerList(answer)
	if err != nil {
		log.Warn().Err(err).Msg("写入答案失败")
		return
	}

	q.Answer = append(q.Answer, answer)

	log.Info().Msg("成功写入回答")
	a.qch <- q
}

// writeAnswerOverFill 写入答案
func (a *App) writeAnswerOverFill() {
	//
}

// handleAnswer 处理消息中可能存在的答案
func (a *App) handleAnswer(m *qq.Msg) {

	q, ok := QABasicSrvPoll[m.Group.ID]
	if !ok {
		return
	}

	for _, v := range q.Answer {
		if v.AnswererID == m.User.ID {
			return
		}
	}

	switch q.Type {
	// 选择题
	case 0:
		if checkAnswerForSelect(m.Chain[0].Text) {
			a.writeAnswer(q, m.User.ID, strings.ToUpper(m.Chain[0].Text))
		}
	// 简答题
	case 1:
		if checkAnswerForFill(m.Chain[0].Text) {
			a.writeAnswer(q, m.User.ID, strings.TrimPrefix(m.Chain[0].Text, "#"))
		}
	// 多选题
	case 2:

	// 填空题
	case 3:

	}

}

// deleteQABasicSrvPoll 使用 i：问题ID(ID) 删除问答基本服务池字段
func (a *App) deleteQABasicSrvPoll(i uint32) (err error) {

	q, err := a.DB.Question().ReadQuestion(i)
	if err != nil {
		return
	}

	delete(QABasicSrvPoll, q.Target)
	return
}
