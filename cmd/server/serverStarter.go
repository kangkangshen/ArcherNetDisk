package main

import (
	"context"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/kangkangshen/ArcherNetDisk/protocol/http11"
	"github.com/kangkangshen/ArcherNetDisk/utils/common"
	"net/http"
	"time"
)

//go's main can not be export
type ServerStarter struct {
	*common.LifecycleListenerBase
	server *Server
}

func NewServerStarter(server *Server) *ServerStarter{
	return &ServerStarter{server:server}
}

func (this *ServerStarter) OnEvent(event *common.LifecycleEvent){
	if event.EventType==common.BEFORE_START_EVENT{
		if runAsMaster(this.server){
			goto RETURN
		}else{
			if !runAsWorker(this.server){
				config.Logger.Fatal("")
			}
		}
	}
	RETURN:
	event.Done()
}

func runAsMaster(server *Server) bool{
	var (
		lkaRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)
	if lkaRespChan=common.RunFor(server.appID,common.CLUSTER_LEADER,server.etcdctl);lkaRespChan!=nil{
		go keepalive(lkaRespChan,server.ctx)
		http.HandleFunc("/storage/",http11.MasterHandler)
		//启动http Server MASTER_LISTEN_CLIENT_PORT必须设置
		if err:=http.ListenAndServe(config.G_Conf.Get(config.MASTER_LISTEN_CLIENT_PORT),nil);err!=nil{
			return false
		}

		//监听，同步worker节点负载信息
		go syncWorkerLoad(server)
		go syncFileMeta(server)
		go syncFileSplitMeta(server)
		return true
	}
	return false
}

//master节点从etcd拉取文件分块元信息
func syncFileSplitMeta(server *Server) {

}

//master节点从etcd拉取文件元信息
func syncFileMeta(server *Server) {

}

//master节点从etcd拉取worker节点负载信息
func syncWorkerLoad(server *Server) {
	var (
		etcdctl *clientv3.Client
		watchChan clientv3.WatchChan
		watchResp *clientv3.WatchResponse
	)
	etcdctl=server.etcdctl
	watchChan=etcdctl.Watch(server.ctx,config.WORKERS_PREFIX)
	for watchResp =range watchChan{
		for _,event :=range watchResp.Events{
			node :=event.PrevKv.Key
			loadInfo:=event.PrevKv.Value
			server.UpdateMetaRepo(string(node),loadInfo)
		}
	}
}

func runAsWorker(server *Server) bool{
	var (
		lkaRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)
	if lkaRespChan=common.RunFor(server.appID,common.CLUSTER_WORKER,server.etcdctl);lkaRespChan!=nil{
		go keepalive(lkaRespChan,server.ctx)
		http.HandleFunc("/storage/",http11.)
		//启动http Server
		if err:=http.ListenAndServe(config.G_Conf.Get(config.WORKER_LISTEN_MASTER_PORT),nil);err!=nil{
			return false
		}

		//监听master节点
		go syncMaster(server)
		go updateWorkerLoad(server)
		return true
	}
	return false
}

func syncMaster(server *Server) {
	var (
		etcdctl *clientv3.Client
		err error
		getResp *clientv3.GetResponse
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
	)
	if server.isMaster{
		return
	}else{
		etcdctl=server.etcdctl
		if getResp,err=etcdctl.Get(server.ctx,common.CLUSTER_LEADER);err==nil{
			server.masterAddr=GetNetAddrFrom(string(getResp.Kvs[0].Value))
		}
		watchChan=etcdctl.Watch(server.ctx,common.CLUSTER_LEADER)
		for watchResp =range watchChan{
			for _,event :=range watchResp.Events{
				masterID:=event.PrevKv.Value

			}
		}
	}
}

func updateWorkerLoad(server *Server) {
	
}




func keepalive(lkaRespChan <-chan *clientv3.LeaseKeepAliveResponse,ctx context.Context){
	for{
		select {
		case <-ctx.Done():
			return
		case lkaResp:=<-lkaRespChan:
			//time.Sleep以nanosecond为单位，lease以秒为单位
			time.Sleep(time.Duration(lkaResp.TTL*time.Second.Nanoseconds()))
		}
	}
}

//v2 may be
func chooseSolutionAndRun(){
	initMode :=config.G_Conf.Get("init_mode")
	switch initMode{
	case config.START_ONLY_HTTP_MODE:startOnlyHttpMode()
	case config.START_ONLY_FAST_HTTP_MODE:startOnlyFastHttpMode()
	case config.START_FULL_MODE:startFullMode()
	case config.START_CUSTOMIZE_MODE:startCustomizeMode()
	default:config.Logger.Fatal("not support")
	}
}

func startOnlyHttpMode() {
	/*
	http.HandleFunc("/storage/",.Handler)
	config.Logger.Fatal(http.ListenAndServe(config.G_Conf.Get()))


	 */
}

func startOnlyFastHttpMode() {
	config.Logger.Fatal("not support")
}

func startFullMode() {
	config.Logger.Fatal("not support")
}

func startCustomizeMode() {
	config.Logger.Fatal("not support")
}


