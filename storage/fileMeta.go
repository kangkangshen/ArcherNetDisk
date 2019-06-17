package storage

import (
	"context"
	"encoding/json"
	"github.com/etcd-io/etcd/clientv3"
	"time"
)

//etcd中文件元信息字段的前缀
const FILES_PREFIX = "files"

type FileMeta struct {
	FileName    string                 `json:"file_name"`     //文件名
	FsMetaUuids [][]byte               `json:"fs_meta_uuids"` //该文件所对应的文件块元信息,此处只获取Uuid,因为调用端并非需要获取所有的文件分块信息，因此这样做可以降低带宽
	FileLength  int64                  `json:"file_length"`   //文件大小
	RefCount    int                    `json:"ref_count"`     //引用数
	CreateTime  time.Time              `json:"create"`        //创建时间
	Uuid        []byte                 `json:"uuid"`          //文件UUID
	Hash        []byte                 `json:"hash"`          //文件内容hash值
	Done        bool                   `json:"done"`          //当前文件是否完成传输，通常是上传
	Options     map[string]interface{} `json:"options"`       //其他可选项
}

//向etcd传送/拉取已经完成的文件元信息
func PullFileMeta(etcdctl *clientv3.Client, uuid []byte) *FileMeta {
	var (
		getResp *clientv3.GetResponse
		err     error
		key     string
		result  *FileMeta
	)
	key = FILES_PREFIX + "/" + string(uuid)
	if getResp, err = etcdctl.Get(context.TODO(), key); err != nil {
		//忽略错误，直接返回nil
		return nil
	} else {
		if len(getResp.Kvs) != 1 {
			//获取了多个？ 此处不会执行到
			return nil
		} else {
			kv := getResp.Kvs[0]
			result = new(FileMeta)
			if err = json.Unmarshal(kv.Value, result); err != nil {
				return nil
			} else {
				return result
			}
		}
	}
}

func UpdateFileMeta(etcdctl *clientv3.Client, fm *FileMeta) error {
	var (
		data []byte
		err  error
		key  string
	)
	if fm == nil {
		//传入的fm为nil 直接返回 ，因此fm的有效性由调用者保证
		return nil
	} else {

		if data, err = json.Marshal(fm); err != nil {
			return err
		} else {
			key = FILES_PREFIX + "/" + string(fm.Uuid)
			if _, err = etcdctl.Put(context.TODO(), key, string(data)); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

func DeleteFileMeta(etcdctl *clientv3.Client, uuid []byte) error {
	var (
		err error
		key string
	)
	if uuid == nil {
		//传入的fm为nil 直接返回 ，因此fm的有效性由调用者保证
		return nil
	} else {
		key = FILES_PREFIX + "/" + string(uuid)
		if _, err = etcdctl.Delete(context.TODO(), key); err != nil {
			return err
		} else {
			return nil
		}
	}
}
