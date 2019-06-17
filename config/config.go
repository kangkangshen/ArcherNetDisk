package config

import (
	"time"
)

type Config struct {
	inited     bool              //是否初始化完成了
	etcdAddr   string            //远程全局配置仓库地址 当前版本使用zookeeper进行全局配置管理,推荐使用该种方式
	timeout    time.Duration     //连接到etcd时的超时时间
	configFile string            //配置文件位置,使用properties文档格式
	properties map[string]string //全部配置项
	id         int               //当前节点ID
	peers      []int             //其他节点ID
	server     string            //网盘服务器地址
	local      string            //当前客户端地址
	optionals  map[string]string //可选配置项
	required   map[string]string //必须配置项
}

func NewConfigUseEtcd(etcdAddr string) (*Config, error) {
	config := &Config{etcdAddr: etcdAddr}
	return config, nil
}

func NewConfigUseFile(configFile string) (*Config, error) {
	config := &Config{configFile: configFile}
	if err := parseConfigFile(configFile, config); err != nil {
		return nil, err
	} else {
		return config, nil
	}
}

func parseConfigFile(s string, config *Config) error {
	return nil
}

func CheckConfig() error {
	return nil
}

func (c *Config) GetInt(name string) int {
	return 0
}
func (c *Config) GetBool(name string) bool {
	return false
}
func (c *Config) GetFloat64(name string) float64 {
	return 0.00
}
func (c *Config) Get(name string) string {
	return ""
}

func (c *Config) parseConfigFile(configFile string, config *Config) error {
	return nil
}
func (c *Config) parseConfigEtcd() error {
	return nil
}
