package main

import (
	"context"
	"github.com/etcd-io/etcd/clientv3"
	jsoniter "github.com/json-iterator/go"
	"github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/kangkangshen/ArcherNetDisk/config/listeners"
	"github.com/kangkangshen/ArcherNetDisk/storage"
	"github.com/kangkangshen/ArcherNetDisk/utils/common"
	"github.com/kangkangshen/ArcherNetDisk/utils/conc"
	"github.com/rs/xid"
)


func main(){
	server:=NewServer()
	var err error
	if err=server.Init();err!=nil{
		config.Logger.Println(err)
		server.Stop()
	}
	if err=server.Start();err!=nil{
		config.Logger.Println(err)
		server.Stop()
	}
	//wait a stop signal



}

type Server struct {
	isMaster bool
	*common.LifecycleBase
	appID string	//server标识，使用IP:PORT来做ID，能够保证唯一性，该ID用来保证其他server节点找到其他节点，并能够识别是master还是worker
	dLock conc.DLocker
	metaRepo storage.MetaRepo
	etcdctl *clientv3.Client
	ctx context.Context		//选主context
	serverAddr string
	masterAddr string
}

func (server *Server) UpdateMetaRepo(node string,loadinfo []byte) {
	fsm:=new(storage.FileSplitMeta)
	if err:=jsoniter.Unmarshal(loadinfo,fsm);err!=nil{
		server.metaRepo.Update(node,nil)
		return
		//todo
	}else {
		server.metaRepo.Update(node,fsm)
	}
}

func NewServer() *Server{

	server:=&Server{appID:xid.New().String()}
	server.AddLifecycleListener(listeners.NewCmdArgParser())
	server.AddLifecycleListener(listeners.NewConfigMgr())
	server.AddLifecycleListener(NewServerStarter(server))
	return server
}



func GetNetAddrFrom(serverID string) string{
	return serverID
}


