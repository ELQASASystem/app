package http

import (
	"github.com/ELQASASystem/app/internal/app"

	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type sign struct{ *app.App }

// Sign 帐号
func Sign() *sign { return &sign{app.AC} }

// in 登录
func (s *sign) in(c *context.Context) {

	res, err := s.DB.Account().ReadAccountsList(c.Params().Get("user"))
	if err != nil {
		log.Error().Err(err).Msg("读取数据库失败")
		c.StatusCode(500)
		return
	}

	if c.URLParam("p") != res.Password {
		c.StatusCode(403)
		return
	}

	c.StatusCode(200)
}
