package listeners

import (
	"github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/kangkangshen/ArcherNetDisk/utils/bufs"
	"github.com/kangkangshen/ArcherNetDisk/utils/common"
	"github.com/kangkangshen/ArcherNetDisk/utils/conc"
	"log"
	"os"
)

type ConfigMgr struct {
	common.LifecycleListenerBase
}

func NewConfigMgr() *ConfigMgr {
	return new(ConfigMgr)
}
func (this *ConfigMgr) OnEvent(event *common.LifecycleEvent) {
	if event.EventType == common.AFTER_INIT_EVENT {
		//cmd args hash been parsed,now apply the config
		applyGlobalConfig()
		applyCustomConfig(event.Host)
	}
	event.Done()
}

func applyGlobalConfig() {
	config.Logger = log.New(os.Stdout, "log-", log.LstdFlags)
	config.BufPool = bufs.NewBufferPool(config.G_Conf.GetInt(config.BUF_POOL_SIZE), config.G_Conf.GetInt(config.BUF_SIZE))
	config.Executors = conc.NewExecutorService()
}

func applyCustomConfig(server Server) {
	//somce code

}
