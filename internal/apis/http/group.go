package http

import (
	"github.com/ELQASASystem/app/internal/app"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type group struct{ *app.App }

// Group 群
func Group() *group { return &group{app.AC} }

// list 群列表
func (g *group) list(c *context.Context) {

	type groupList struct {
		ID       uint64 `json:"id"`        // 群号
		Name     string `json:"name"`      // 群名
		MemCount uint16 `json:"mem_count"` // 群成员数
	}

	var data []groupList
	for _, v := range g.Cli.C.GroupList {
		data = append(data, groupList{uint64(v.Uin), v.Name, v.MemberCount})
	}

	_, _ = c.JSON(data)
}

// praise 表扬
func (g *group) praise(c *context.Context) {

	i, err := c.Params().GetUint64("i")
	if err != nil {
		log.Error().Err(err).Msg("解析目标群失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	var ids []uint64
	err = jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(c.URLParam("mem"), &ids)
	if err != nil {
		log.Error().Err(err).Msg("解析目标成员失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	m := g.Cli.NewMsg().AddText("表扬以下答对的同学:\n")
	for _, id := range ids {
		m.AddAt(id)
	}
	g.Cli.SendGroupMsg(m.AddText("\n希望同学们再接再厉!").To(i))

	_, _ = c.JSON(iris.Map{"message": "yes"})
}
