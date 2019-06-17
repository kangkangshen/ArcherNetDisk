package job

import (
	"github.com/kangkangshen/ArcherNetDisk/storage"
	"github.com/kangkangshen/ArcherNetDisk/utils/common"
	"time"
)

//可暂停-重启的Job模型
type Job interface {
	common.Lifecycle
	Pausable
}

type Pausable interface {
	Pause()
	Resume()
}

//每一个上传/下载/列出请求对应于一个Job ,Job内部进一步由多个文件传输块执行
type baseJob struct {
	jobID       string            `json:"job_id"`       //创建的JobID 使用UUID进行识别，客户端启停JOB时用来定位的字段
	createdTime time.Time         `json:"created_time"` //Job创建时间
	endTime     time.Time         `json:"end_time"`     //job完成时间
	creatorID   string            `json:"creator_id"`   //客户端创建者ID，在未来将启用，由master解析请求获得
	creatorAddr string            `json:"creator_addr"` //上传者的网络地址
	workersID   []string          `json:"workers_id"`   //具体执行该任务的workers节点的ID，由master反填
	fs          *storage.FileMeta `json:"fs"`           //该任务所对应的文件元信息
	speed       int               `json:"speed"`        //该任务完成后统计而得的平均上传/下载速度

}
