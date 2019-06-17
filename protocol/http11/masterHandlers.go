package http11

import (
	"github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/kangkangshen/ArcherNetDisk/storage"
	"github.com/kangkangshen/ArcherNetDisk/storage/job"
	"net/http"
	"strconv"
)

// parse request layer
func MasterHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPut:
		put(w, r)
	case http.MethodGet:
		get(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//get request
func get(writer http.ResponseWriter, request *http.Request) {
	var (
		cJob *job.DownloadJob //job corresponding to the current request
	)
	switch request.Header.Get(config.UPLOAD_TYPE) {
	case config.DOWNLOAD_NEW:
		size, err := strconv.ParseInt(request.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if cJob = assignJob(size); cJob == nil {
			//There is no enough space to allocate
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		config.Executors.SubmitTask(cJob)
	case config.CONTINUED_TRANSMISSION:
		//There is no data transmission for too long，job has been destroyed
		uuid := []byte(request.Header.Get(config.JOB_ID))
		//here uuid length depends on specified uuid gen lib ,here is xid ,uuid len is 12
		if len(uuid) != 12 {
			//Uuid format does not match
			writer.WriteHeader(http.StatusBadRequest)
		}
		if cJob = locateJob(uuid); cJob == nil {
			//no found job
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		//found job
		config.Executors.SubmitTask(cJob)
	case config.PAUSE_TRANSMISSION:
		//pause not be supported yet
		//writer.WriteHeader(http.StatusForbidden)
		//There is no data transmission for too long，job has been destroyed
		uuid := []byte(request.Header.Get(config.JOB_ID))
		//here uuid length depends on specified uuid gen lib ,here is xid ,uuid len is 12
		if len(uuid) != 12 {
			//Uuid format does not match
			writer.WriteHeader(http.StatusBadRequest)
		}
		if cJob = locateJob(uuid); cJob == nil {
			//no found job
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		//found job
		cJob.Pause()
		return
	}
	//request succeed ,return job.id and 200 ok

	writer.Header().Add(config.JOB_ID, string(cJob.Uuid()))
	writer.WriteHeader(http.StatusOK)

}

//put request
func put(writer http.ResponseWriter, request *http.Request) {
	var (
		cJob *job.UploadJob //job corresponding to the current request
	)
	switch request.Header.Get(config.UPLOAD_TYPE) {
	case config.UPLOAD_NEW:
		size, err := strconv.ParseInt(request.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if cJob = assignJob(size); cJob == nil {
			//There is no enough space to allocate
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		packJobAndResponse(cJob, writer)
	case config.CONTINUED_TRANSMISSION:
		//There is no data transmission for too long，job has been destroyed
		uuid := []byte(request.Header.Get(config.JOB_ID))
		//here uuid length depends on specified uuid gen lib ,here is xid ,uuid len is 12
		if len(uuid) != 12 {
			//Uuid format does not match
			writer.WriteHeader(http.StatusBadRequest)
		}
		if cJob = locateJob(uuid); cJob == nil {
			//no found job
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		//found job ,pack job info and return ,wait client transmite file
		packJobAndResponse(cJob, writer)
	case config.PAUSE_TRANSMISSION:
		//pause not be supported yet
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	//request succeed ,return job.id and 200 ok
	writer.Header().Add(config.JOB_ID, string(cJob.Uuid()))
	writer.WriteHeader(http.StatusOK)
}

func assignJob(size int64) *job.UploadJob {
	fss := storage.AllocateEnoughSpace(size)
	return wrap(fss)
}

func locateJob(uuid []byte) *job.UploadJob {
	fm := config.MetaRepo.GetFileMeta(uuid)
	return wrap(pickUnfinishedFileSplit(fm))
}

func pickUnfinishedFileSplit(fm *storage.FileMeta) []*storage.FileSplit {
	fss := make([]*storage.FileSplit, 0)
	for fsm := range fm.FsMetas {
		if !fsm.Done {
			append(fss, fsm.FileSplit)
		}
	}
	return fss

}

//wrap splits to a job ,to add some other code
func wrap(splits []*storage.FileSplit) *job.UploadJob {
	cJob := storage.NewUploadJob(splits)
	for plugin := range config.Plugins {
		if plugin.match(cJob) {
			plugin.plugin(cJob)
		}
	}
	return cJob
}

/*
func put(writer http.ResponseWriter, request *http.Request) {
	var(
		wg sync.WaitGroup
		fileSize int64
		space []*storage.FileSplit
		err error
	)
	if space,err=AllocateEnoughspace(fileSize);err!=nil{
		//分配空间失败
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}
	wg.Add(len(space))
	for fs:=range space{
		go func(wg *sync.WaitGroup){
			err =StartTransmite(fs)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	writer.WriteHeader(http.StatusOK)
}

func get(writer http.ResponseWriter, request *http.Request) {
	var(
		wg sync.WaitGroup
		fileSplits []*storage.FileSplit
		fileInfo *storage.FileMeta
		err error
	)
	if fileInfo,err=LocateFile(request);err!=nil{
		//此处假设是没有找到对应的文件
		writer.WriteHeader(http.StatusBadRequest)
	}
	fileSplits=ExtractFileSplits(fileInfo)
	for fs:=range fileSplits{
		wg.Add(len(fileSplits))
			go func(wg *sync.WaitGroup){
				err =StartTransmite(fs)
				wg.Done()
			}(&wg)
	}
	wg.Wait()
}

func StartTransmite(i int) error {
	return nil
}

func ExtractFileSplits(meta *storage.FileMeta) []*storage.FileSplit {
	return nil
}

func LocateFile(request *http.Request) (*storage.FileMeta, error) {
	return nil,nil
}

*/
