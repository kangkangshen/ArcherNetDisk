package config

import (
	"github.com/kangkangshen/ArcherNetDisk/utils/bufs"
	"github.com/kangkangshen/ArcherNetDisk/utils/conc"
	"log"
)

//全局对象
var (
	IsM        bool    //当前是master还是worker
	G_Conf     *Config //全局配置信息储存地
	Logger     *log.Logger
	BufPool    *bufs.BufferPool
	Executors  conc.ExecutorService
	MasterAddr string      //master位置
	WorkerLoad interface{} //worker节点的负载信息
)
