package class

import (
	"github.com/ELQASASystem/app/configs"
	"github.com/ELQASASystem/app/internal/app/database"
	"github.com/ELQASASystem/app/internal/app/qq"

	"github.com/rs/zerolog/log"
)

var Bot *qq.Rina // Bot 机器人对象
var qch chan *Question
var db = database.New()

// New 新建一个机器人
func New(qc chan *Question) {

	var (
		c  = configs.GetAllConf()
		ch = make(chan *qq.Msg, 10)
		r  = qq.NewRina(c.QQID, c.QQPassword, &ch)
	)

	Bot = r
	qch = qc

	if db.ConnectDB(c.DatabaseUrl) != nil {
		log.Panic().Msg("数据库连接失败")
	}
	r.RegEventHandle()
	go monitorGroup()

}

// Database 获取数据库事务实例
func Database() *database.Database { return db }
