package http

import (
	"github.com/ELQASASystem/app/internal/app"
	"github.com/ELQASASystem/app/internal/app/database"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type question struct{ *app.App }

// Question 问题
func Question() *question { return &question{app.AC} }

// list 问题列表
func (q *question) list(c *context.Context) {

	res, err := q.DB.Question().ReadQuestionList(c.Params().Get("u"))
	if err != nil {
		log.Error().Err(err).Msg("读取问题列表失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(res)
}

// read 读取问题
func (q *question) read(c *context.Context) {

	i, err := c.Params().GetUint32("i")
	if err != nil {
		log.Error().Err(err).Msg("解析问题ID失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	res, err := q.ReadQuestion(i)
	if err != nil {
		log.Error().Err(err).Msg("获取答题失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(res)
}

// start 开始问答
func (q *question) start(c *context.Context) {

	qid, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Msg("解析问题 ID 失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	if err = q.StartQA(qid); err != nil {
		log.Error().Err(err).Msg("开启问答失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"message": "yes"})
}

// stop 停止问答
func (q *question) stop(c *context.Context) {

	qid, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Msg("解析问题 ID 失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	if err = q.StopQA(qid); err != nil {
		log.Error().Err(err).Msg("停止答题失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"message": "yes"})
}

// prepare 准备问答
func (q *question) prepare(c *context.Context) {

	qid, err := c.Params().GetUint32("question_id")
	if err != nil {
		log.Error().Err(err).Msg("解析问题 ID 失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	if err = q.PrepareQA(qid); err != nil {
		log.Error().Err(err).Msg("准备答题失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"message": "yes"})
}

// add 新增问题
func (q *question) add(c *context.Context) {

	// 初始化问题列表, 用于解析 JSON 后储存
	que := database.QuestionListTab{}

	if err := c.ReadJSON(&que); err != nil {
		log.Error().Err(err).Msg("解析传入 JSON 失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	err := q.DB.Question().WriteQuestionList(&que)
	if err != nil {
		log.Error().Err(err).Msg("新增答题失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"message": "yes"})
}

// delete 删除问题
func (q *question) delete(c *context.Context) {

	_, err := c.Params().GetUint64("question_id")
	if err != nil {
		log.Error().Err(err).Msg("解析问题ID失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	// TODO: 调用数据库删除 QJNKSM:这个先咕咕
	// class.Database().Question().RemoveQuestion(qid)

	_, _ = c.JSON(iris.Map{"message": "yes"})
}
