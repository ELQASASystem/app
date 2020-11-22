package http

import (
	"github.com/ELQASASystem/app/internal/app"
	"github.com/ELQASASystem/app/internal/app/database"

	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type question struct{ *app.App }

// Question 问题
func Question() *question { return &question{app.AC} }

// list 列表
func (q *question) list(c *context.Context) {

	type list []*database.QuestionListTab
	type group struct {
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
		MemCount uint16 `json:"mem_count"`
	}

	type res struct {
		List   list    `json:"questions"`
		Groups []group `json:"groups"`
	}

	l, err := q.DB.Question().ReadQuestionList(c.URLParam("u"))
	if err != nil {
		log.Error().Err(err).Msg("读取问题列表失败")
		c.StatusCode(500)
		return
	}

	var g []group
	for _, v := range q.Cli.C.GroupList {
		g = append(g, group{uint64(v.Uin), v.Name, v.MemberCount})
	}

	_, _ = c.JSON(res{l, g})
}

// detail 详情
func (q *question) detail(c *context.Context) {

	queID, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Str("字段", "Question ID").Msg("读取字段失败")
		c.StatusCode(400)
		return
	}

	que, err := q.ReadQuestion(queID)
	if err != nil {
		log.Error().Err(err).Msg("读取问答失败")
		c.StatusCode(500)
		return
	}

	_, _ = c.JSON(que)
}

// new 新建
func (q *question) new(c *context.Context) {

	que := database.QuestionListTab{}

	if err := c.ReadJSON(&que); err != nil {
		log.Error().Err(err).Msg("读取数据失败")
		c.StatusCode(400)
		return
	}

	err := q.DB.Question().WriteQuestionList(&que)
	if err != nil {
		log.Error().Err(err).Msg("新增答题失败")
		c.StatusCode(500)
		return
	}

	c.StatusCode(201)
}

// edit 编辑问题
func (q *question) edit(c *context.Context) {

	_, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Str("字段", "Question ID").Msg("读取字段失败")
		c.StatusCode(400)
		return
	}

	// TODO: 更新数据库 QJNKSM:咕咕
	// q.DB.Question().UpdateQuestion(queID)

	c.StatusCode(200)
}

// status 状态
func (q *question) status(c *context.Context) {

	queID, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Str("字段", "Question ID").Msg("读取字段失败")
		c.StatusCode(400)
		return
	}

	code, err := c.URLParamInt("c")
	if err != nil {
		log.Error().Err(err).Str("字段", "Code").Msg("读取字段失败")
		c.StatusCode(400)
		return
	}

	switch code {
	case 0:
		err = q.PrepareQA(queID)
	case 1:
		err = q.StartQA(queID)
	case 2:
		err = q.StopQA(queID)
	}
	if err != nil {
		log.Error().Err(err).
			Uint32("目标问题", queID).
			Int("目标状态", code).
			Msg("更新问题状态失败")
		c.StatusCode(500)
		return
	}

	c.StatusCode(200)
}

// delete 删除问题
func (q *question) delete(c *context.Context) {

	_, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Str("字段", "Question ID").Msg("读取字段失败")
		c.StatusCode(400)
		return
	}

	// TODO: 调用数据库删除 QJNKSM:这个先咕咕
	//q.DB.Question().RemoveQuestion(queID)

	c.StatusCode(200)
}
