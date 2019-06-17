package config

const (
	MASTER_LISTEN_CLIENT_PORT = "master_listen_client_port"
	WORKER_LISTEN_MASTER_PORT = "worker_listen_master_port"
	MASTER_ADDR               = "master_addr"
	START_ONLY_HTTP_MODE      = "http11"
	START_ONLY_FAST_HTTP_MODE = "http2"
	START_CUSTOMIZE_MODE      = "customize"
	START_FULL_MODE           = "full"
	DLOCKER_PREFIX            = "lock-"
)

//buf pool
const (
	BUF_SIZE      = "buf_size"
	BUF_POOL_SIZE = "buf_pool_size"
)

const (
	SIGNAL_UNLOCK              = "signal_unlock"
	DLOCK_LEASE_TIME           = 3 //分布式锁每次续租时常，默认3秒
	DLOCK_LEASE_KEEPALIVE_TIME = 2 //隔？？时间续租一次
)

const (
	DLOCKER_NEW = iota
	DLOCER_LOCK
	DLOCKER_UNLOCK
)

//支持的hash算法列表
const (
	SHA256 = "SHA256" //DEFAULT
)

//custom field name  in request
const (
	JOB_ID = "job-id"
)

//bad flags
const ()

const (
	WORKERS_PREFIX = "workers"
)

//req type
const (
	UPLOAD_TYPE            = "upload_type"
	UPLOAD_NEW             = "upload_new"             //开始上传
	DOWNLOAD_NEW           = "download_new"           //开始下载
	CONTINUED_TRANSMISSION = "continued_transmission" //继续上传/下载
	PAUSE_TRANSMISSION     = "pause_transmission"     //暂停上传/下载

)
