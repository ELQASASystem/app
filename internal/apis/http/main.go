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
		g := Group()

		G.Get("/list", auth.auth(g.list))
		G.Get("/{i}/praise", auth.auth(g.praise))
	}
	{
		Q := API.Party("questions")
		q := Question()

		Q.Get("/list", auth.auth(q.list))
		Q.Get("/question/{question_id}", auth.auth(q.detail))

		Q.Post("/", auth.auth(q.new))
		Q.Put("/question/{question_id}", auth.auth(q.edit))
		Q.Put("/question/{question_id}/status", auth.auth(q.status))
		Q.Delete("/question/{question_id}", auth.auth(q.delete))
	}
	{
		M := API.Party("market")
		m := Market()

		M.Get("/{subject}/list", auth.auth(m.list))
		M.Get("/{i}/copy", auth.auth(m.copy))
	}
	{
		U := API.Party("upload")
		u := Upload()

		U.Options("/docx", auth.auth(u.options))
		U.Post("/docx", auth.auth(u.docx))
		U.Get("/docx/{p}/parse", auth.auth(u.parseDocx))

		U.Options("/picture", auth.auth(u.options))
		U.Post("/picture", auth.auth(u.picture))
	}

	if err := app.Listen(":4040"); err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
