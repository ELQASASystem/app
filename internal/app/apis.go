package class

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

// startAPI 开启 API 服务
func startAPI() {

	app := iris.New()

	// hello
	Hello := app.Party("hello")
	{

		Hello.Get("/", func(c *context.Context) {
			_, _ = c.JSON(iris.Map{"message": "hello"})
		})

	}

	// question
	Question := app.Party("question")
	{

		// 获取问题列表
		Question.Get("/{u}", func(c *context.Context) {

			res, err := readQuestionList(c.Params().Get("u"))
			if err != nil {
				log.Error().Err(err).Msg("读取问题列表时出错")
				return
			}

			_, _ = c.JSON(res)

		})

		// 获取问题市场
		Question.Get("/market", func(c *context.Context) {

			res, err := readQuestionMarket()
			if err != nil {
				log.Error().Err(err).Msg("读取问题列表时出错")
				return
			}

			_, _ = c.JSON(res)

		})

	}

	err := app.Listen(":8080")
	if err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
