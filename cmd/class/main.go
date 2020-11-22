package main

import (
	_ "github.com/ELQASASystem/server/cmd/class/basic" // 全局初始化
	"github.com/ELQASASystem/server/internal"
)

func main() {

	internal.Main()
	select {}

}
