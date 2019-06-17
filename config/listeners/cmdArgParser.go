package listeners

import (
	"flag"
	conf "github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/kangkangshen/ArcherNetDisk/utils/common"
	"log"
)

//命令行解析器

type CmdArgParser struct {
	*common.LifecycleListenerBase
	cmdArgs map[string]interface{}
}

func NewCmdArgParser() *CmdArgParser {
	return &CmdArgParser{cmdArgs: make(map[string]interface{})}
}
func (this *CmdArgParser) OnEvent(event *common.LifecycleEvent) {
	if event.EventType == common.BEFORE_INIT_EVENT {
		//init event to parse cmd args
		var (
			etcdAddr   string
			configFile string
			err        error
			logger     = conf.Logger
		)
		flag.StringVar(&etcdAddr, "etcdAddr", "", "specify the etcd-server location")
		flag.StringVar(&configFile, "configFile", "", "specify the configFile location")
		if etcdAddr != "" && configFile != "" {
			log.Fatal("cannot be specified simultaneously etcd-server and configFile")
		}
		if etcdAddr != "" {
			if conf.G_Conf, err = conf.NewConfigUseEtcd(etcdAddr); err != nil {
				logger.Fatal(err)
			}
		} else {
			if conf.G_Conf, err = conf.NewConfigUseFile(configFile); err != nil {
				logger.Fatal(err)
			}
		}
		if err = conf.CheckConfig(); err != nil {
			logger.Fatal(err)
		}
	}
	event.Done()
}
