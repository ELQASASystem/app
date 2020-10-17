package basic

import (
	"fmt"
	"os"

	"github.com/ELQASASystem/app/configs"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {

	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "01-02|15:04:05"}
	writer.FormatCaller = func(i interface{}) string {
		return fmt.Sprintf("\x1b[1m%v\x1b[0m \x1b[36m>\x1b[0m", i.(string)[25:])
	}

	log.Logger = log.Output(writer).With().Caller().Logger()

	log.Info().Msg("Copyright (C) 2020-present  CCServe  AGPL-3.0 License | version：" + configs.CommitId)

	err := configs.ReadConfigs()
	if err != nil {
		log.Panic().Err(err).Msg("读取配置文件失败")
	}

}
