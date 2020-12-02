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

	API.Post("login/{user}", Sign(auth).in)

	{
		G := API.Party("group")

		G.Get("/list", Group().list)
		G.Get("/{i}/praise", Group().praise)
	}

	{
		Q := API.Party("questions")

		Q.Get("/list", Question().list)
		Q.Get("/question/{question_id}", Question().detail)

		Q.Post("/", Question().new)
		Q.Put("/question/{question_id}", Question().edit)
		Q.Put("/question/{question_id}/status", Question().status)
		Q.Delete("/question/{question_id}", Question().delete)
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
