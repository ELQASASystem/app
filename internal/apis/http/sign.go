package http

import (
	"net/http"
	"time"

	"github.com/ELQASASystem/server/internal/app"

	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type sign struct {
	*app.App
	Auth *auth
}

// Sign 帐号
func Sign(a *auth) *sign { return &sign{app.AC, a} }

// in 登录
func (s *sign) in(c *context.Context) {

	u := c.Params().Get("user")

	res, err := s.DB.Account().ReadAccountsList(u)
	if err != nil {
		log.Error().Err(err).Msg("读取数据库失败")
		c.StatusCode(500)
		return
	}

	if c.URLParam("p") != res.Password {
		c.StatusCode(403)
		return
	}

	s.Auth.generateLoginToken(u, c)

	c.SetCookie(&http.Cookie{
		Name: "user", Value: u,
		Expires: time.Now().AddDate(0, 1, 0), Secure: true,
	})
	c.StatusCode(200)
}
