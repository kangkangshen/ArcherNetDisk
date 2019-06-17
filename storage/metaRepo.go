package storage

import (
	"github.com/etcd-io/etcd/clientv3"
	"time"
)

type MetaRepo interface {
	GetFileMeta(uuid []byte) *FileMeta
	GetFileSplitMeta(uuid []byte) *FileSplitMeta
	GetWorkerLoad(workerid string)
	Update(node string, meta *FileSplitMeta)
}

//simple impl for MetaRepo ,uses cache design ,depends on etcd ,maybe  depends on db in next version
type simpleMetaRepo struct {
	cache    map[string]*FileMeta
	client   *clientv3.Client
	prefix   string
	lastsync time.Time
}

func NewMetaRepo() MetaRepo {
	return nil
}

func Update(node string, meta *FileSplitMeta) {

}
