package defines

// go-cqhttp 配置文件结构(仅读取需要的部分)
type CQHttpConfig struct {
	Account struct {
		Uin      string `yaml:"uin"`
		Password string `yaml:"password"`
	}
}
