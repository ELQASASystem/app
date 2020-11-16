package http

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

// New 新建一个 API 服务
func New() {

	app := iris.New()
	API := app.Party("apis/")

	// 测试
	API.Party("hello").Get("/", func(c *context.Context) {
		_, _ = c.JSON(iris.Map{"message": "hello"})
	})

	API.Party("sign").Get("/{u}/in", Sign().in)

	{
		G := API.Party("group")

		G.Get("/list", Group().list)
		G.Get("/{i}/praise", Group().praise)
	}

	{
		Q := API.Party("question")

		Q.Get("/{u}/list", Question().list)
		Q.Get("/a/{i}", Question().read)

		Q.Get("/{question_id}/start", Question().start)
		Q.Get("/{question_id}/stop", Question().stop)
		Q.Get("/{question_id}/prepare", Question().prepare)

		Q.Post("/add", Question().add)
		Q.Get("/{question_id}/delete", Question().delete)
	}

	{
		M := API.Party("market")

		M.Get("/{subject}/list", Market().list)
		M.Get("/{i}/copy", Market().copy)
	}

	{
		U := API.Party("upload")

		U.Options("/docx", Upload().options)
		U.Post("/docx", Upload().docx)
		U.Get("/docx/{p}/parse", Upload().parseDocx)

		U.Options("/picture", Upload().options)
		U.Post("/picture", Upload().picture)

		U.Get("/{text}/split", Upload().split)
	}

	if err := app.Listen(":4040"); err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
