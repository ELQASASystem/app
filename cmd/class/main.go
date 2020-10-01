package main

import (
	_ "github.com/ELQASASystem/app/cmd/class/basic" // 全局初始化
	"github.com/ELQASASystem/app/internal/app"
)

func main() {

	class.New()
	select {}

}
