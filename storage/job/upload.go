package job

import "github.com/kangkangshen/ArcherNetDisk/storage"

type UploadJob struct {
	baseJob
}

func NewUploadJob(fs *storage.FileSplit) Job {
	return nil
}
