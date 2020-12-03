package http

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

// New 新建一个 API 服务
func New() {

	app := iris.New()
	auth := NewAuth()
	API := app.Party("apis/")

	// 测试
	API.Party("hello").Get("/", func(c *context.Context) {

		type info struct {
			Message           string `json:"message"`
			CookieUser        string `json:"cookie_user"`
			CookieLoginToken  string `json:"cookie_login_token"`
			CookieOnlineToken string `json:"cookie_online_token"`
		}

		_, _ = c.JSON(info{
			"Hello",
			c.GetCookie("user"),
			c.GetCookie("loginToken"),
			c.GetCookie("onlineToken"),
		})
	})

	{
		sign := Sign(auth)
		API.Post("login/{user}", sign.in)
	}

	{
		G := API.Party("group")
		group := Group(auth)

		G.Get("/list", group.list)
		G.Get("/{i}/praise", group.praise)
	}

	{
		Q := API.Party("questions")
		questions := Question(auth)

		Q.Get("/list", questions.list)
		Q.Get("/question/{question_id}", questions.detail)

		Q.Post("/", questions.new)
		Q.Put("/question/{question_id}", questions.edit)
		Q.Put("/question/{question_id}/status", questions.status)
		Q.Delete("/question/{question_id}", questions.delete)
	}

	{
		M := API.Party("market")
		market := Market(auth)

		M.Get("/{subject}/list", market.list)
		M.Get("/{i}/copy", market.copy)
	}

	{
		U := API.Party("upload")
		upload := Upload(auth)

		U.Options("/docx", upload.options)
		U.Post("/docx", upload.docx)
		U.Get("/docx/{p}/parse", upload.parseDocx)

		U.Options("/picture", upload.options)
		U.Post("/picture", upload.picture)

		U.Get("/{text}/split", upload.split)
	}

	if err := app.Listen(":4040"); err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
