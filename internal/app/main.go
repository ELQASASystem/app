package class

import (
	"github.com/ELQASASystem/app/configs"
	"github.com/rs/zerolog/log"
)

var classBot *Rina // classBot 机器人对象

// New 新建一个机器人
func New() {

	var (
		c  = configs.GetAllConf()
		ch = make(chan *QQMsg, 10)
		r  = newRina(c.QQID, c.QQPassword, &ch)
	)

	classBot = r

	if connectDB(c.DatabaseUrl) != nil {
		log.Panic().Msg("数据库连接失败")
	}
	r.regEventHandle()
	go monitorGroup()
	go startAPI()

}
