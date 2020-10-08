package class

import (
	"github.com/ELQASASystem/app/configs"
	"github.com/ELQASASystem/app/internal/app/database"
	"github.com/ELQASASystem/app/internal/app/qq"

	"github.com/rs/zerolog/log"
)

var Bot *qq.Rina // Bot 机器人对象

// New 新建一个机器人
func New() {

	var (
		c  = configs.GetAllConf()
		ch = make(chan *qq.Msg, 10)
		r  = qq.NewRina(c.QQID, c.QQPassword, &ch)
	)

	Bot = r

	if database.Class.ConnectDB(c.DatabaseUrl) != nil {
		log.Panic().Msg("数据库连接失败")
	}
	r.RegEventHandle()
	go monitorGroup()

}
