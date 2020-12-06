package app

import (
	"github.com/ELQASASystem/server/configs"
	"github.com/ELQASASystem/server/internal/app/database"
	"github.com/ELQASASystem/server/internal/qq"

	"github.com/rs/zerolog/log"
)

var AC *App

type App struct {
	Cli *qq.Rina           // Cli QQ 客户端
	mch chan *qq.Msg       // mch 消息同步管道
	wch chan *qq.Msg       // wch 词云同步管道
	qch chan *Question     // qch 问答同步管道
	DB  *database.Database // DB 数据库
}

// New 新建一个机器人
func New(qc chan *Question, wc chan *qq.Msg) (app *App) {

	var (
		conf = configs.GetAllConf()
		mch  = make(chan *qq.Msg, 10)
		r    = qq.NewRina(conf.QQID, conf.QQPassword, &mch)
		db   = database.New()
	)

	if db.ConnectDB(conf.DatabaseUrl) != nil {
		log.Panic().Msg("数据库连接失败")
	}

	app = &App{r, mch, wc, qc, db}
	AC = app

	r.RegEventHandle()
	go app.monitorGroup()
	return
}
