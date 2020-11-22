package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Conf 配置
type Conf struct {
	QQID        uint64 `yaml:"QQID"`        // QQID QQ 帐号
	QQPassword  string `yaml:"QQPassword"`  // QQPassword QQ 密码
	DatabaseUrl string `yaml:"DatabaseUrl"` // DatabaseUrl 数据库地址
}

var (
	CommitID string // CommitID 提交的短ID
	confs    *Conf  // confs 配置信息
)

// ReadConfigs 读取配置
func ReadConfigs() (err error) {

	f, err := ioutil.ReadFile("configs/configs.yml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(f, &confs)
	if err != nil {
		return
	}

	return

}

// GetAllConf 获取所有配置
func GetAllConf() *Conf { return confs }
