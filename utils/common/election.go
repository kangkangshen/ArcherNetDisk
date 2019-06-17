package common

import (
	"github.com/etcd-io/etcd/clientv3"
)

//etcd选主实现

const CLUSTER_LEADER = "cluster_leader"
const CLUSTER_WORKER = "cluster_worker"
const TMP_JOB = "tmp_job" //当worker升级为master的时候，用于保存未完成的任务的key
func RunFor(appID string, pos string, cli *clientv3.Client) <-chan *clientv3.LeaseKeepAliveResponse {
	return nil
}
