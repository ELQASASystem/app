package http

import (
	class "github.com/ELQASASystem/app/internal/app"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type sign struct{}

// Sign 帐号
func Sign() *sign { return new(sign) }

// in 登录
func (s *sign) in(c *context.Context) {

	res, err := class.Database().Account().ReadAccountsList(c.Params().Get("u"))
	if err != nil {
		log.Error().Err(err).Msg("获取用户帐号失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	if c.URLParam("p") != res.Password {
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"message": "yes"})
}
