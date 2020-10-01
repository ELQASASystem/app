package configs

// Conf 配置
type Conf struct {
	QQID        uint64 // QQID QQ 帐号
	QQPassword  string // QQPassword QQ 密码
	DatabaseUrl string // DatabaseUrl 数据库地址
}

var (
	CommitId string // CommitId 提交的短ID
	confs    *Conf  // confs 配置信息
)

// ReadConfigs 读取配置
func ReadConfigs() { confs = fullConfigs }

// GetAllConf 获取所有配置
func GetAllConf() *Conf { return confs }
