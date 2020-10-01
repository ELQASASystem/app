package class

import "github.com/ELQASASystem/app/configs"

var classBot *Rina // classBot 机器人对象

// New 新建一个机器人
func New() {

	c := configs.GetAllConf()

	var (
		ch = make(chan *QQMsg, 10)
		r  = newRina(c.QQID, c.QQPassword, &ch)
	)

	go monitorGroup()
	r.regEventHandle()

	classBot = r

}
