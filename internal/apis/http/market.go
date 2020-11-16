package http

import (
	"github.com/ELQASASystem/app/internal/app"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type market struct{}

// Market 问题市场
func Market() *market { return new(market) }

// list 问题市场列表
func (m *market) list(c *context.Context) {

	res, err := class.Database().Question().ReadQuestionMarket(c.Params().GetUint8Default("subject", 0))
	if err != nil {
		log.Error().Err(err).Msg("读取问题列表失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(res)
}

// copy 复制问题
func (m *market) copy(c *context.Context) {

	qid, err := c.Params().GetUint32("i")
	if err != nil {
		log.Error().Err(err).Msg("解析问题ID失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	q, err := class.Database().Question().ReadQuestion(qid)
	if err != nil {
		log.Error().Err(err).Msg("读取题目失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	q.CreatorID = c.URLParam("user")
	q.Market = false

	if err := class.Database().Question().WriteQuestionList(q); err != nil {
		log.Error().Err(err).Msg("写入题目失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"message": "yes"})
}
